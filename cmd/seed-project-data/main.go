package main

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"flag"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// 测试项目数据
var testProjects = []struct {
	Name        string
	Description string
	Difficulty  string
	Deadline    time.Time
}{
	{
		Name:        "某金融系统安全测试",
		Description: "## 项目背景\n\n某大型金融机构需要对其核心业务系统进行全面的安全评估，包括但不限于网银系统、移动APP、API接口等进行渗透测试和漏洞挖掘。\n\n## 测试范围\n\n1. 网上银行系统（Web端）\n2. 手机银行APP（iOS/Android）\n3. 开放API接口（OAuth2.0认证）\n4. 后台管理系统",
		Difficulty:  "hard",
		Deadline:    time.Now().AddDate(0, 2, 0),
	},
	{
		Name:        "电商平台渗透测试",
		Description: "## 项目背景\n\n某知名电商平台需要进行年度安全评估，重点关注支付流程、用户数据保护和商家系统的安全性。\n\n## 测试范围\n\n1. 用户端网站和APP\n2. 商家后台管理系统\n3. 支付系统接口\n4. 物流对接系统",
		Difficulty:  "medium",
		Deadline:    time.Now().AddDate(0, 1, 15),
	},
	{
		Name:        "政务系统安全评估",
		Description: "## 项目背景\n\n某省级政务云平台需要进行等保测评前的安全加固，需要发现并修复潜在的安全漏洞。\n\n## 测试范围\n\n1. 政务服务门户\n2. 统一认证平台\n3. 数据交换中心\n4. 移动政务APP",
		Difficulty:  "expert",
		Deadline:    time.Now().AddDate(0, 3, 0),
	},
	{
		Name:        "移动APP安全审计",
		Description: "## 项目背景\n\n一款日活超百万的社交APP需要进行安全审计，重点关注用户隐私保护和通信安全。\n\n## 测试范围\n\n1. Android客户端\n2. iOS客户端\n3. 后端API接口\n4. 即时通信系统",
		Difficulty:  "easy",
		Deadline:    time.Now().AddDate(0, 1, 0),
	},
	{
		Name:        "物联网设备漏洞挖掘",
		Description: "## 项目背景\n\n某智能家居厂商需要对其全系列IoT设备进行安全测试，包括智能门锁、摄像头、网关等设备。\n\n## 测试范围\n\n1. 智能门锁固件\n2. 智能摄像头\n3. 智能网关\n4. 云端控制平台",
		Difficulty:  "hard",
		Deadline:    time.Now().AddDate(0, 2, 15),
	},
}

func main() {
	// 命令行参数
	userID := flag.Uint("user", 0, "指定用户ID（必填）")
	username := flag.String("username", "", "指定用户名（与user二选一）")
	clean := flag.Bool("clean", false, "清理该用户的所有测试数据")
	flag.Parse()

	// 验证参数
	if *userID == 0 && *username == "" {
		log.Fatal("请指定用户ID (-user) 或用户名 (-username)")
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	db := database.InitDB(cfg)

	// 获取用户
	var user domain.User
	var err error

	if *userID > 0 {
		err = db.First(&user, *userID).Error
	} else {
		err = db.Where("username = ?", *username).First(&user).Error
	}

	if err != nil {
		log.Fatalf("用户不存在: %v", err)
	}

	fmt.Printf("目标用户: %s (ID: %d)\n", user.Username, user.ID)

	if *clean {
		cleanTestData(db, user.ID)
		return
	}

	// 生成测试数据
	generateTestData(db, user.ID)
}

// generateTestData 生成测试数据
func generateTestData(db *gorm.DB, userID uint) {
	fmt.Println("\n开始生成测试数据...")

	createdProjects := 0
	createdAssignments := 0
	createdAttachments := 0

	// 测试附件模板
	testAttachments := []struct {
		Name string
		Size string
		Type string
		URL  string
	}{
		{Name: "项目授权书.pdf", Size: "256KB", Type: "pdf", URL: "#"},
		{Name: "测试范围说明.docx", Size: "128KB", Type: "doc", URL: "#"},
		{Name: "系统架构图.png", Size: "1.2MB", Type: "image", URL: "#"},
	}

	for _, tp := range testProjects {
		// 检查项目是否已存在
		var existingProject domain.Project
		result := db.Where("name = ?", tp.Name).First(&existingProject)

		var project domain.Project

		if result.Error == gorm.ErrRecordNotFound {
			// 创建新项目
			deadline := tp.Deadline
			project = domain.Project{
				Name:        tp.Name,
				Description: tp.Description,
				Difficulty:  tp.Difficulty,
				Deadline:    &deadline,
				Status:      "recruiting",
			}

			if err := db.Create(&project).Error; err != nil {
				log.Printf("创建项目失败 [%s]: %v", tp.Name, err)
				continue
			}
			createdProjects++
			fmt.Printf("  ✓ 创建项目: %s (ID: %d)\n", project.Name, project.ID)

			// 为新项目创建附件
			for i, att := range testAttachments {
				attachment := domain.ProjectAttachment{
					ProjectID: project.ID,
					Name:      att.Name,
					URL:       att.URL,
					Size:      att.Size,
					Type:      att.Type,
					SortOrder: i + 1,
				}
				if err := db.Create(&attachment).Error; err != nil {
					log.Printf("创建附件失败: %v", err)
				} else {
					createdAttachments++
				}
			}
			fmt.Printf("    ✓ 添加附件: %d 个\n", len(testAttachments))
		} else if result.Error != nil {
			log.Printf("查询项目失败 [%s]: %v", tp.Name, result.Error)
			continue
		} else {
			project = existingProject
			fmt.Printf("  - 项目已存在: %s (ID: %d)\n", project.Name, project.ID)
		}

		// 检查指派是否已存在
		var existingAssignment domain.ProjectAssignment
		result = db.Where("project_id = ? AND user_id = ?", project.ID, userID).First(&existingAssignment)

		if result.Error == gorm.ErrRecordNotFound {
			// 创建指派记录
			assignment := domain.ProjectAssignment{
				ProjectID: project.ID,
				UserID:    userID,
			}

			if err := db.Create(&assignment).Error; err != nil {
				log.Printf("创建指派失败 [项目ID: %d, 用户ID: %d]: %v", project.ID, userID, err)
				continue
			}
			createdAssignments++
			fmt.Printf("    ✓ 指派给用户 (指派ID: %d)\n", assignment.ID)
		} else if result.Error != nil {
			log.Printf("查询指派失败: %v", result.Error)
		} else {
			fmt.Printf("    - 已指派给该用户\n")
		}
	}

	fmt.Printf("\n生成完成!\n")
	fmt.Printf("  - 新建项目: %d 个\n", createdProjects)
	fmt.Printf("  - 新建附件: %d 个\n", createdAttachments)
	fmt.Printf("  - 新建指派: %d 个\n", createdAssignments)
	fmt.Printf("\n用户现在可以在项目大厅查看这些项目并接受任务。\n")
}

// cleanTestData 清理测试数据
func cleanTestData(db *gorm.DB, userID uint) {
	fmt.Println("\n开始清理测试数据...")

	// 删除该用户的所有任务记录
	result := db.Where("user_id = ?", userID).Delete(&domain.ProjectTask{})
	fmt.Printf("  ✓ 删除任务记录: %d 条\n", result.RowsAffected)

	// 删除该用户的所有指派记录
	result = db.Where("user_id = ?", userID).Delete(&domain.ProjectAssignment{})
	fmt.Printf("  ✓ 删除指派记录: %d 条\n", result.RowsAffected)

	fmt.Println("\n清理完成!")
}
