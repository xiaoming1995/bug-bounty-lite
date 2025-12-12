package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// ProjectSeeder 项目测试数据填充器
type ProjectSeeder struct {
	db *gorm.DB
}

func NewProjectSeeder(db *gorm.DB) *ProjectSeeder {
	return &ProjectSeeder{db: db}
}

// Seed 填充项目测试数据（追加模式）
func (s *ProjectSeeder) Seed(force bool) error {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 项目模板
	projectTemplates := []struct {
		Name        string
		Description string
		Note        string
		Status      string
	}{
		{
			Name:        "科技公司官网",
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
	}

	// 生成 5-8 个新项目（每次执行都生成新的）
	numProjects := rand.Intn(4) + 5
	fmt.Printf("[INFO] Generating %d new projects...\n", numProjects)

	successCount := 0
	for i := 0; i < numProjects; i++ {
		template := projectTemplates[rand.Intn(len(projectTemplates))]
		timestamp := time.Now().UnixNano()

		project := domain.Project{
			Name:        fmt.Sprintf("%s_%d", template.Name, timestamp/1000000+int64(i)),
			Description: template.Description,
			Note:        template.Note,
			Status:      template.Status,
		}

		if err := s.db.Create(&project).Error; err != nil {
			log.Printf("[WARN] Failed to create project %s: %v", project.Name, err)
		} else {
			successCount++
			fmt.Printf("[OK] Created project: %s (ID: %d)\n", project.Name, project.ID)
		}

		// 短暂延迟确保时间戳不同
		time.Sleep(time.Millisecond)
	}

	fmt.Printf("[INFO] Seeded %d/%d projects successfully\n", successCount, numProjects)
	return nil
}
