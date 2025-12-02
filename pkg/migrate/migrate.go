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
		&domain.Report{},
		&domain.UserInfoChangeRequest{},
	)

	if err != nil {
		log.Printf("[ERROR] Migration failed: %v", err)
		return err
	}

	// 添加表注释 (MySQL)
	m.addTableComments()

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

// addTableComments 添加表级别注释 (MySQL)
func (m *Migrator) addTableComments() {
	tableComments := map[string]string{
		"users":                     "用户表 - 存储平台用户信息(白帽子/厂商/管理员)",
		"reports":                   "漏洞报告表 - 存储白帽子提交的漏洞报告",
		"user_info_change_requests": "用户信息变更申请表 - 存储用户信息变更申请，需后台审核",
	}

	for table, comment := range tableComments {
		sql := fmt.Sprintf("ALTER TABLE `%s` COMMENT '%s'", table, comment)
		if err := m.db.Exec(sql).Error; err != nil {
			log.Printf("[WARN] Failed to add comment for table %s: %v", table, err)
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
