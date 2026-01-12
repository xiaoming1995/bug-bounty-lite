package migrate

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migrator 数据库迁移器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移器实例
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Run 执行数据库迁移
func (m *Migrator) Run() error {
	startTime := time.Now()
	fmt.Println("[INFO] Running Database Migrations...")

	// 获取迁移前的表信息（用于日志）
	beforeTables := m.getTableNames()

	// 执行迁移
	// 注意：AutoMigrate 只会添加缺失的列、索引，不会删除或修改现有列
	err := m.db.AutoMigrate(
		&domain.User{},
		&domain.Organization{},
		&domain.UserUpdateLog{},
		&domain.Report{},
		&domain.UserInfoChangeRequest{},
		&domain.Project{},
		&domain.SystemConfig{},
	)

	if err != nil {
		log.Printf("[ERROR] Migration failed: %v", err)
		return err
	}

	// 确保 deleted_at 列存在（手动添加，防止 AutoMigrate 没有正确添加）
	m.ensureDeletedAtColumns()

	// 删除 reports 表的外键约束（改用代码逻辑验证）
	m.dropForeignKeys()

	// 添加表注释 (MySQL)
	m.addTableComments()

	// 添加字段注释 (MySQL)
	m.addColumnComments()

	// 获取迁移后的表信息
	afterTables := m.getTableNames()

	// 打印迁移结果
	duration := time.Since(startTime)
	fmt.Printf("[OK] Database Migrations completed in %v\n", duration)

	// 打印新增的表
	newTables := m.diffTables(beforeTables, afterTables)
	if len(newTables) > 0 {
		fmt.Printf("[INFO] New tables created: %v\n", newTables)
	}

	// 打印当前所有表
	fmt.Printf("[INFO] Current tables: %v\n", afterTables)

	return nil
}

// ensureDeletedAtColumns 确保 deleted_at 列存在于 projects 和 reports 表
func (m *Migrator) ensureDeletedAtColumns() {
	fmt.Println("[INFO] Ensuring deleted_at columns exist...")

	tables := []string{"projects", "reports"}

	for _, table := range tables {
		// 检查列是否存在
		var count int64
		query := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_SCHEMA = DATABASE() 
			AND TABLE_NAME = '%s' 
			AND COLUMN_NAME = 'deleted_at'
		`, table)

		if err := m.db.Raw(query).Scan(&count).Error; err != nil {
			log.Printf("[WARN] Failed to check deleted_at column for %s: %v", table, err)
			continue
		}

		if count == 0 {
			// 列不存在，添加它
			sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL", table)
			if err := m.db.Exec(sql).Error; err != nil {
				log.Printf("[WARN] Failed to add deleted_at column to %s: %v", table, err)
			} else {
				fmt.Printf("[OK] Added deleted_at column to %s\n", table)
			}

			// 添加索引
			indexSQL := fmt.Sprintf("CREATE INDEX `idx_%s_deleted_at` ON `%s` (`deleted_at`)", table, table)
			if err := m.db.Exec(indexSQL).Error; err != nil {
				// 索引可能已存在，忽略错误
				log.Printf("[INFO] Index on deleted_at for %s may already exist: %v", table, err)
			} else {
				fmt.Printf("[OK] Added index on deleted_at for %s\n", table)
			}
		} else {
			fmt.Printf("[INFO] deleted_at column already exists in %s\n", table)
		}
	}
}

// dropForeignKeys 删除数据库外键约束（改用代码逻辑验证数据一致性）
func (m *Migrator) dropForeignKeys() {
	fmt.Println("[INFO] Checking and removing foreign key constraints...")

	// 需要删除外键的表列表
	tables := []string{"reports", "users"}

	for _, table := range tables {
		// 查询该表的所有外键
		type ForeignKey struct {
			ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
		}

		var foreignKeys []ForeignKey
		query := fmt.Sprintf(`
			SELECT CONSTRAINT_NAME 
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS 
			WHERE TABLE_SCHEMA = DATABASE() 
			AND TABLE_NAME = '%s' 
			AND CONSTRAINT_TYPE = 'FOREIGN KEY'
		`, table)

		if err := m.db.Raw(query).Scan(&foreignKeys).Error; err != nil {
			log.Printf("[WARN] Failed to query foreign keys for %s: %v", table, err)
			continue
		}

		if len(foreignKeys) == 0 {
			fmt.Printf("[INFO] No foreign keys found on %s table\n", table)
			continue
		}

		fmt.Printf("[INFO] Found %d foreign key(s) on %s to remove\n", len(foreignKeys), table)

		for _, fk := range foreignKeys {
			sql := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", table, fk.ConstraintName)
			if err := m.db.Exec(sql).Error; err != nil {
				log.Printf("[WARN] Failed to drop foreign key %s from %s: %v", fk.ConstraintName, table, err)
			} else {
				fmt.Printf("[OK] Dropped foreign key: %s from %s\n", fk.ConstraintName, table)
			}
		}
	}
}

// addTableComments 添加表级别注释 (MySQL)
func (m *Migrator) addTableComments() {
	tableComments := map[string]string{
		"users":                     "用户表 - 存储平台用户信息(白帽子/厂商/管理员)",
		"reports":                   "漏洞报告表 - 存储白帽子提交的漏洞报告，包含项目关联、漏洞类型、详情等信息",
		"user_info_change_requests": "用户信息变更申请表 - 存储用户信息变更申请，需后台审核",
		"projects":                  "项目表 - 存储平台项目资料信息",
		"system_configs":            "系统配置表 - 存储各类系统配置信息（漏洞类型、危害等级等），支持通过config_type区分不同类型的配置",
		"organizations":             "组织管理表 - 存储机构、部门等组织架构信息",
		"user_update_logs":          "用户修改记录表 - 存储用户关键信息（如简介、组织绑定）的变更审计日志",
	}

	for table, comment := range tableComments {
		sql := fmt.Sprintf("ALTER TABLE `%s` COMMENT '%s'", table, comment)
		if err := m.db.Exec(sql).Error; err != nil {
			log.Printf("[WARN] Failed to add comment for table %s: %v", table, err)
		}
	}
}

// addColumnComments 添加字段级别注释 (MySQL)
func (m *Migrator) addColumnComments() {
	// 字段注释列表
	columnComments := []struct {
		table   string
		column  string
		comment string
	}{
		// projects 表
		{
			table:   "projects",
			column:  "status",
			comment: "项目状态(active:活跃/inactive:非活跃)",
		},
		// system_configs 表
		{
			table:   "system_configs",
			column:  "config_type",
			comment: "配置类型(vulnerability_type:漏洞类型/severity_level:危害等级/project_category:项目分类等)",
		},
		{
			table:   "system_configs",
			column:  "config_key",
			comment: "配置键(如:SQL_INJECTION/XSS/CSRF等，用于程序内部识别)",
		},
		{
			table:   "system_configs",
			column:  "config_value",
			comment: "配置值(显示名称，如:SQL注入/XSS跨站脚本，用于前端显示)",
		},
		{
			table:   "system_configs",
			column:  "sort_order",
			comment: "排序顺序(数字越小越靠前)",
		},
		{
			table:   "system_configs",
			column:  "status",
			comment: "配置状态(active:启用/inactive:禁用)",
		},
		{
			table:   "system_configs",
			column:  "extra_data",
			comment: "扩展数据(JSON格式，存储额外信息如图标、颜色等)",
		},
		// reports 表新增字段
		{
			table:   "reports",
			column:  "project_id",
			comment: "关联项目ID(必填，关联projects表)",
		},
		{
			table:   "reports",
			column:  "vulnerability_name",
			comment: "漏洞名称(必填，文本输入)",
		},
		{
			table:   "reports",
			column:  "vulnerability_type_id",
			comment: "关联漏洞类型配置ID(必填，关联system_configs表，config_type='vulnerability_type')",
		},
		{
			table:   "reports",
			column:  "vulnerability_impact",
			comment: "漏洞的危害(文本输入，描述漏洞可能造成的危害)",
		},
		{
			table:   "reports",
			column:  "self_assessment_id",
			comment: "危害自评配置ID(关联system_configs表，config_type=severity_level)",
		},
		{
			table:   "reports",
			column:  "vulnerability_url",
			comment: "漏洞链接(URL格式，指向漏洞相关页面)",
		},
		{
			table:   "reports",
			column:  "vulnerability_detail",
			comment: "漏洞详情(文本输入，详细描述漏洞情况)",
		},
		{
			table:   "reports",
			column:  "attachment_url",
			comment: "附件地址(文件上传后的URL，单个文件，后续可扩展为多个)",
		},
		// users 表新增字段
		{
			table:   "users",
			column:  "bio",
			comment: "个人简介(文本输入)",
		},
		{
			table:   "users",
			column:  "org_id",
			comment: "所属组织ID(关联organizations表)",
		},
		{
			table:   "users",
			column:  "last_login_at",
			comment: "最后登录时间",
		},
	}

	for _, cc := range columnComments {
		// 查询字段的当前定义
		type ColumnInfo struct {
			ColumnType    string  `gorm:"column:column_type"`
			IsNullable    string  `gorm:"column:is_nullable"`
			ColumnDefault *string `gorm:"column:column_default"`
		}

		var colInfo ColumnInfo
		query := fmt.Sprintf(`
			SELECT 
				COLUMN_TYPE as column_type,
				IS_NULLABLE as is_nullable,
				COLUMN_DEFAULT as column_default
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			AND TABLE_NAME = '%s'
			AND COLUMN_NAME = '%s'`, cc.table, cc.column)

		err := m.db.Raw(query).Scan(&colInfo).Error
		if err != nil || colInfo.ColumnType == "" {
			log.Printf("[WARN] Failed to get column info for %s.%s: %v", cc.table, cc.column, err)
			continue
		}

		// 构建 MODIFY COLUMN 语句
		nullClause := "NULL"
		if colInfo.IsNullable == "NO" {
			nullClause = "NOT NULL"
		}

		defaultClause := ""
		if colInfo.ColumnDefault != nil && *colInfo.ColumnDefault != "" && *colInfo.ColumnDefault != "NULL" {
			defaultClause = fmt.Sprintf("DEFAULT '%s'", *colInfo.ColumnDefault)
		}

		sql := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s %s %s COMMENT '%s'",
			cc.table, cc.column, colInfo.ColumnType, nullClause, defaultClause, cc.comment)

		if err := m.db.Exec(sql).Error; err != nil {
			log.Printf("[WARN] Failed to add comment for column %s.%s: %v", cc.table, cc.column, err)
		} else {
			log.Printf("[INFO] Added comment for column %s.%s", cc.table, cc.column)
		}
	}
}

// getTableNames 获取当前数据库中的所有表名 (MySQL)
func (m *Migrator) getTableNames() []string {
	var tables []string
	m.db.Raw("SHOW TABLES").Scan(&tables)
	return tables
}

// diffTables 计算新增的表
func (m *Migrator) diffTables(before, after []string) []string {
	beforeMap := make(map[string]bool)
	for _, t := range before {
		beforeMap[t] = true
	}

	var newTables []string
	for _, t := range after {
		if !beforeMap[t] {
			newTables = append(newTables, t)
		}
	}
	return newTables
}

// Status 打印迁移状态
func (m *Migrator) Status() {
	fmt.Println("Migration Status")
	fmt.Println("-------------------")

	tables := m.getTableNames()
	fmt.Printf("Tables in database: %d\n", len(tables))
	for i, t := range tables {
		fmt.Printf("  %d. %s\n", i+1, t)
	}

	// 检查表结构
	fmt.Println("\nTable Details:")
	m.printTableInfo("users")
	m.printTableInfo("reports")
	m.printTableInfo("user_info_change_requests")
	m.printTableInfo("projects")
	m.printTableInfo("system_configs")
}

// seedInitialData 添加初始数据
func (m *Migrator) seedInitialData() error {
	// 1. 插入漏洞类型初始数据
	var count int64
	m.db.Model(&domain.SystemConfig{}).Where("config_type = ?", "vulnerability_type").Count(&count)
	if count == 0 {
		vulnerabilityTypes := []domain.SystemConfig{
			{ConfigType: "vulnerability_type", ConfigKey: "SQL_INJECTION", ConfigValue: "SQL注入", Description: "SQL注入漏洞", SortOrder: 1, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "XSS", ConfigValue: "XSS跨站脚本", Description: "跨站脚本攻击", SortOrder: 2, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "CSRF", ConfigValue: "CSRF跨站请求伪造", Description: "跨站请求伪造", SortOrder: 3, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "FILE_UPLOAD", ConfigValue: "文件上传漏洞", Description: "文件上传漏洞", SortOrder: 4, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "COMMAND_INJECTION", ConfigValue: "命令执行", Description: "命令注入漏洞", SortOrder: 5, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "INFORMATION_DISCLOSURE", ConfigValue: "信息泄露", Description: "敏感信息泄露", SortOrder: 6, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "PRIVILEGE_ESCALATION", ConfigValue: "权限绕过", Description: "权限提升/绕过", SortOrder: 7, Status: "active"},
			{ConfigType: "vulnerability_type", ConfigKey: "OTHER", ConfigValue: "其他", Description: "其他类型漏洞", SortOrder: 99, Status: "active"},
		}

		for _, vt := range vulnerabilityTypes {
			if err := m.db.Create(&vt).Error; err != nil {
				log.Printf("[WARN] Failed to create vulnerability type %s: %v", vt.ConfigKey, err)
			}
		}
		log.Println("[INFO] Seeded vulnerability types successfully")
	} else {
		log.Println("[INFO] Vulnerability types already exist, skipping seed")
	}

	return nil
}

// printTableInfo 打印表结构信息 (MySQL)
func (m *Migrator) printTableInfo(tableName string) {
	type ColumnInfo struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Null    string `gorm:"column:Null"`
		Key     string `gorm:"column:Key"`
		Default string `gorm:"column:Default"`
	}

	var columns []ColumnInfo
	m.db.Raw("DESCRIBE " + tableName).Scan(&columns)

	if len(columns) == 0 {
		fmt.Printf("\n  [%s] - Table not found\n", tableName)
		return
	}

	fmt.Printf("\n  [%s]\n", tableName)
	for _, col := range columns {
		nullable := "NOT NULL"
		if col.Null == "YES" {
			nullable = "NULL"
		}
		fmt.Printf("    - %-20s %-25s %s\n", col.Field, col.Type, nullable)
	}
}
