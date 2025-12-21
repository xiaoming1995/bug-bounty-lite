package main

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

func main() {
	// 命令行参数
	helpFlag := flag.Bool("help", false, "Show help message")
	forceFlag := flag.Bool("force", false, "Force init even if data exists")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - System Data Initializer")
	fmt.Println("==========================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 执行数据初始化
	initializer := NewInitializer(db)
	if err := initializer.Init(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Init failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] System data initialized successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - System Data Initializer")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/init/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force init even if data exists (will skip existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/init/main.go           # Initialize system data")
	fmt.Println("  go run cmd/init/main.go -force    # Force init (skip existing)")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make init         # Initialize system data")
	fmt.Println("  make init-force   # Force init")
}

// Initializer 系统数据初始化器
type Initializer struct {
	db *gorm.DB
}

func NewInitializer(db *gorm.DB) *Initializer {
	return &Initializer{db: db}
}

// Init 初始化系统必需数据
func (i *Initializer) Init(force bool) error {
	// 初始化系统配置数据
	if err := i.initSystemConfigs(force); err != nil {
		return fmt.Errorf("failed to init system configs: %w", err)
	}

	return nil
}

// initSystemConfigs 初始化系统配置数据（危害等级、漏洞类型等）
func (i *Initializer) initSystemConfigs(force bool) error {
	// 初始化危害等级
	if err := i.initSeverityLevels(force); err != nil {
		return err
	}

	// 初始化漏洞类型
	if err := i.initVulnerabilityTypes(force); err != nil {
		return err
	}

	return nil
}

// initSeverityLevels 初始化危害等级
func (i *Initializer) initSeverityLevels(force bool) error {
	var count int64
	i.db.Model(&domain.SystemConfig{}).Where("config_type = ?", "severity_level").Count(&count)

	if count > 0 && !force {
		fmt.Println("[INFO] System configs (severity_level) already exist, skipping init (use -force to override)")
		return nil
	}

	// 危害等级配置
	severityLevels := []domain.SystemConfig{
		{
			ConfigType:  "severity_level",
			ConfigKey:   "NONE",
			ConfigValue: "无危害",
			Description: "无安全危害",
			SortOrder:   1,
			Status:      "active",
		},
		{
			ConfigType:  "severity_level",
			ConfigKey:   "LOW",
			ConfigValue: "低危",
			Description: "低风险漏洞，影响较小",
			SortOrder:   2,
			Status:      "active",
		},
		{
			ConfigType:  "severity_level",
			ConfigKey:   "MEDIUM",
			ConfigValue: "中危",
			Description: "中等风险漏洞，有一定影响",
			SortOrder:   3,
			Status:      "active",
		},
		{
			ConfigType:  "severity_level",
			ConfigKey:   "HIGH",
			ConfigValue: "高危",
			Description: "高风险漏洞，影响较大",
			SortOrder:   4,
			Status:      "active",
		},
		{
			ConfigType:  "severity_level",
			ConfigKey:   "CRITICAL",
			ConfigValue: "严重",
			Description: "严重漏洞，影响极大",
			SortOrder:   5,
			Status:      "active",
		},
	}

	successCount := 0
	for _, config := range severityLevels {
		// 检查是否已存在（根据类型和键）
		var existing domain.SystemConfig
		if err := i.db.Where("config_type = ? AND config_key = ?", config.ConfigType, config.ConfigKey).First(&existing).Error; err == nil {
			if !force {
				fmt.Printf("[SKIP] System config '%s:%s' already exists\n", config.ConfigType, config.ConfigKey)
				continue
			}
		}

		if err := i.db.Create(&config).Error; err != nil {
			log.Printf("[WARN] Failed to create system config %s:%s: %v", config.ConfigType, config.ConfigKey, err)
		} else {
			successCount++
			fmt.Printf("[OK] Created system config: %s - %s\n", config.ConfigKey, config.ConfigValue)
		}
	}

	fmt.Printf("[INFO] Initialized %d/%d severity levels successfully\n", successCount, len(severityLevels))
	return nil
}

// initVulnerabilityTypes 初始化漏洞类型
func (i *Initializer) initVulnerabilityTypes(force bool) error {
	var count int64
	i.db.Model(&domain.SystemConfig{}).Where("config_type = ?", "vulnerability_type").Count(&count)

	if count > 0 && !force {
		fmt.Println("[INFO] System configs (vulnerability_type) already exist, skipping init (use -force to override)")
		return nil
	}

	// 漏洞类型配置
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

	successCount := 0
	for _, config := range vulnerabilityTypes {
		// 检查是否已存在（根据类型和键）
		var existing domain.SystemConfig
		if err := i.db.Where("config_type = ? AND config_key = ?", config.ConfigType, config.ConfigKey).First(&existing).Error; err == nil {
			if !force {
				fmt.Printf("[SKIP] System config '%s:%s' already exists\n", config.ConfigType, config.ConfigKey)
				continue
			}
		}

		if err := i.db.Create(&config).Error; err != nil {
			log.Printf("[WARN] Failed to create system config %s:%s: %v", config.ConfigType, config.ConfigKey, err)
		} else {
			successCount++
			fmt.Printf("[OK] Created system config: %s - %s\n", config.ConfigKey, config.ConfigValue)
		}
	}

	fmt.Printf("[INFO] Initialized %d/%d vulnerability types successfully\n", successCount, len(vulnerabilityTypes))
	return nil
}
