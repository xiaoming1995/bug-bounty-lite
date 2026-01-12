package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// OrganizationSeeder 组织测试数据填充器
type OrganizationSeeder struct {
	db *gorm.DB
}

func NewOrganizationSeeder(db *gorm.DB) *OrganizationSeeder {
	return &OrganizationSeeder{db: db}
}

// Seed 填充组织测试数据
func (s *OrganizationSeeder) Seed(force bool) error {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 组织模板
	orgTemplates := []struct {
		Name        string
		Description string
	}{
		{Name: "核心研发部", Description: "负责公司核心产品的开发与维护"},
		{Name: "网络安全部", Description: "负责公司业务的安全审计、漏洞扫描与应急响应"},
		{Name: "大数据中心", Description: "负责海量业务数据的清洗、存储与分析挖据"},
		{Name: "运维保障部", Description: "负责公司基础设施的日常维护与高可用架构调优"},
		{Name: "财务管理部", Description: "负责公司资金核算、财报发布与预算执行"},
		{Name: "人力资源部", Description: "负责人才招聘、员工培训与组织文化建设"},
		{Name: "行政后勤部", Description: "负责公司日常行政事务与后勤物力保障"},
		{Name: "市场营销中心", Description: "负责产品推广、品牌建设与市场渠道拓展"},
	}

	// 随机生成 4-6 个组织
	numOrgs := rand.Intn(3) + 4
	fmt.Printf("[INFO] Generating %d new organizations...\n", numOrgs)

	successCount := 0
	for i := 0; i < numOrgs; i++ {
		template := orgTemplates[rand.Intn(len(orgTemplates))]
		timestamp := time.Now().UnixNano()

		org := domain.Organization{
			Name:        fmt.Sprintf("%s_%d", template.Name, timestamp/1000000+int64(i)),
			Description: template.Description,
		}

		if err := s.db.Create(&org).Error; err != nil {
			log.Printf("[WARN] Failed to create organization %s: %v", org.Name, err)
		} else {
			successCount++
			fmt.Printf("[OK] Created organization: %s (ID: %d)\n", org.Name, org.ID)
		}

		// 短暂延迟
		time.Sleep(time.Millisecond)
	}

	fmt.Printf("[INFO] Seeded %d/%d organizations successfully\n", successCount, numOrgs)
	return nil
}
