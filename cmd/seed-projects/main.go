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
	forceFlag := flag.Bool("force", false, "Force seed even if data exists")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - Projects Test Data Seeder")
	fmt.Println("=============================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 执行数据填充
	seeder := NewProjectSeeder(db)
	if err := seeder.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Projects test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Projects Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-projects/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will skip existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-projects/main.go           # Seed projects test data")
	fmt.Println("  go run cmd/seed-projects/main.go -force    # Force seed (skip existing)")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make seed-projects         # Seed projects test data")
	fmt.Println("  make seed-projects-force   # Force seed")
}

// ProjectSeeder 项目测试数据填充器
type ProjectSeeder struct {
	db *gorm.DB
}

func NewProjectSeeder(db *gorm.DB) *ProjectSeeder {
	return &ProjectSeeder{db: db}
}

// Seed 填充项目测试数据
func (s *ProjectSeeder) Seed(force bool) error {
	var count int64
	s.db.Model(&domain.Project{}).Count(&count)

	if count > 0 && !force {
		fmt.Println("[INFO] Projects already exist, skipping seed (use -force to override)")
		return nil
	}

	projects := []domain.Project{
		{
			Name:        "某科技公司官网",
			Description: "公司官方网站，包含用户注册、登录、产品展示等功能",
			Note:        "重要项目，需要重点关注安全漏洞",
			Status:      "active",
		},
		{
			Name:        "电商平台系统",
			Description: "在线购物平台，包含商品管理、订单处理、支付系统等核心功能",
			Note:        "涉及支付功能，安全要求高",
			Status:      "active",
		},
		{
			Name:        "移动应用后端API",
			Description: "为移动应用提供数据接口服务，包含用户管理、内容管理等",
			Note:        "API接口较多，需要全面测试",
			Status:      "active",
		},
		{
			Name:        "企业管理系统",
			Description: "企业内部管理系统，包含员工管理、财务管理、审批流程等",
			Note:        "内部系统，数据敏感",
			Status:      "active",
		},
		{
			Name:        "在线教育平台",
			Description: "在线学习平台，包含课程管理、视频播放、在线考试等功能",
			Note:        "用户量大，需要关注性能和安全",
			Status:      "active",
		},
		{
			Name:        "社交网络平台",
			Description: "社交网络应用，包含用户动态、消息推送、好友关系等功能",
			Note:        "用户生成内容较多，需要关注XSS等漏洞",
			Status:      "active",
		},
		{
			Name:        "金融支付系统",
			Description: "金融支付平台，包含账户管理、交易处理、风控系统等",
			Note:        "涉及资金安全，安全要求极高",
			Status:      "active",
		},
		{
			Name:        "内容管理系统",
			Description: "内容发布和管理系统，包含文章编辑、媒体管理、权限控制等",
			Note:        "需要关注文件上传和权限控制",
			Status:      "active",
		},
		{
			Name:        "API网关服务",
			Description: "统一API网关，提供接口路由、限流、认证等功能",
			Note:        "核心基础设施，需要高可用性",
			Status:      "active",
		},
		{
			Name:        "测试项目（已归档）",
			Description: "这是一个已归档的测试项目",
			Note:        "测试用项目，已不再使用",
			Status:      "inactive",
		},
	}

	successCount := 0
	for _, project := range projects {
		// 检查是否已存在（根据名称）
		var existing domain.Project
		if err := s.db.Where("name = ?", project.Name).First(&existing).Error; err == nil {
			if !force {
				fmt.Printf("[SKIP] Project '%s' already exists\n", project.Name)
				continue
			}
		}

		if err := s.db.Create(&project).Error; err != nil {
			log.Printf("[WARN] Failed to create project %s: %v", project.Name, err)
		} else {
			successCount++
			fmt.Printf("[OK] Created project: %s\n", project.Name)
		}
	}

	fmt.Printf("[INFO] Seeded %d/%d projects successfully\n", successCount, len(projects))
	return nil
}

