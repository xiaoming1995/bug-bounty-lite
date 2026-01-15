package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

// Cleaner 测试数据清理器
type Cleaner struct {
	db *gorm.DB
}

// NewCleaner 创建清理器实例
func NewCleaner(db *gorm.DB) *Cleaner {
	return &Cleaner{db: db}
}

// CleanAll 清理所有测试数据（保留系统配置和管理员账户）
func (c *Cleaner) CleanAll() error {
	fmt.Println("[INFO] Cleaning all test data...")

	// 按依赖顺序清理（先清理有外键依赖的表）
	if err := c.CleanReportComments(); err != nil {
		return err
	}
	if err := c.CleanReports(); err != nil {
		return err
	}
	if err := c.CleanArticleComments(); err != nil {
		return err
	}
	if err := c.CleanArticleLikes(); err != nil {
		return err
	}
	if err := c.CleanArticleViews(); err != nil {
		return err
	}
	if err := c.CleanArticles(); err != nil {
		return err
	}
	if err := c.CleanProjectAttachments(); err != nil {
		return err
	}
	if err := c.CleanProjectTasks(); err != nil {
		return err
	}
	if err := c.CleanProjectAssignments(); err != nil {
		return err
	}
	if err := c.CleanProjects(); err != nil {
		return err
	}
	if err := c.CleanUserUpdateLogs(); err != nil {
		return err
	}
	if err := c.CleanUserInfoChanges(); err != nil {
		return err
	}
	if err := c.CleanUsers(); err != nil {
		return err
	}
	if err := c.CleanAvatars(); err != nil {
		return err
	}
	if err := c.CleanOrganizations(); err != nil {
		return err
	}

	fmt.Println("[OK] All test data cleaned successfully!")
	return nil
}

// CleanUsers 清理用户数据（保留 admin 角色用户）
func (c *Cleaner) CleanUsers() error {
	var count int64
	c.db.Model(&domain.User{}).Where("role != ?", "admin").Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No test users to clean")
		return nil
	}

	// 先清理关联数据：用户更新日志中引用这些用户的记录
	result := c.db.Unscoped().Where("role != ?", "admin").Delete(&domain.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean users: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d users (kept admin users)\n", result.RowsAffected)
	return nil
}

// CleanProjects 清理项目数据（硬删除）
func (c *Cleaner) CleanProjects() error {
	var count int64
	c.db.Unscoped().Model(&domain.Project{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No projects to clean")
		return nil
	}

	result := c.db.Unscoped().Where("1 = 1").Delete(&domain.Project{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean projects: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d projects\n", result.RowsAffected)
	return nil
}

// CleanReports 清理漏洞报告数据（硬删除）
func (c *Cleaner) CleanReports() error {
	var count int64
	c.db.Unscoped().Model(&domain.Report{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No reports to clean")
		return nil
	}

	result := c.db.Unscoped().Where("1 = 1").Delete(&domain.Report{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean reports: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d reports\n", result.RowsAffected)
	return nil
}

// CleanReportComments 清理报告评论数据
func (c *Cleaner) CleanReportComments() error {
	var count int64
	c.db.Model(&domain.ReportComment{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No report comments to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ReportComment{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean report comments: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d report comments\n", result.RowsAffected)
	return nil
}

// CleanAvatars 清理头像数据
func (c *Cleaner) CleanAvatars() error {
	var count int64
	c.db.Model(&domain.Avatar{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No avatars to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.Avatar{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean avatars: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d avatars\n", result.RowsAffected)
	return nil
}

// CleanOrganizations 清理组织数据
func (c *Cleaner) CleanOrganizations() error {
	var count int64
	c.db.Model(&domain.Organization{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No organizations to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.Organization{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean organizations: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d organizations\n", result.RowsAffected)
	return nil
}

// CleanArticles 清理文章数据
func (c *Cleaner) CleanArticles() error {
	var count int64
	c.db.Model(&domain.Article{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No articles to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.Article{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean articles: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d articles\n", result.RowsAffected)
	return nil
}

// CleanArticleViews 清理文章访问记录
func (c *Cleaner) CleanArticleViews() error {
	var count int64
	c.db.Model(&domain.ArticleView{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No article views to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ArticleView{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean article views: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d article views\n", result.RowsAffected)
	return nil
}

// CleanArticleLikes 清理文章点赞记录
func (c *Cleaner) CleanArticleLikes() error {
	var count int64
	c.db.Model(&domain.ArticleLike{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No article likes to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ArticleLike{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean article likes: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d article likes\n", result.RowsAffected)
	return nil
}

// CleanArticleComments 清理文章评论
func (c *Cleaner) CleanArticleComments() error {
	var count int64
	c.db.Model(&domain.ArticleComment{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No article comments to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ArticleComment{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean article comments: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d article comments\n", result.RowsAffected)
	return nil
}

// CleanUserInfoChanges 清理用户信息变更申请
func (c *Cleaner) CleanUserInfoChanges() error {
	var count int64
	c.db.Model(&domain.UserInfoChangeRequest{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No user info change requests to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.UserInfoChangeRequest{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean user info changes: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d user info change requests\n", result.RowsAffected)
	return nil
}

// CleanUserUpdateLogs 清理用户更新日志
func (c *Cleaner) CleanUserUpdateLogs() error {
	var count int64
	c.db.Model(&domain.UserUpdateLog{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No user update logs to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.UserUpdateLog{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean user update logs: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d user update logs\n", result.RowsAffected)
	return nil
}

// CleanProjectAssignments 清理项目指派记录
func (c *Cleaner) CleanProjectAssignments() error {
	var count int64
	c.db.Model(&domain.ProjectAssignment{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No project assignments to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ProjectAssignment{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean project assignments: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d project assignments\n", result.RowsAffected)
	return nil
}

// CleanProjectTasks 清理项目任务记录
func (c *Cleaner) CleanProjectTasks() error {
	var count int64
	c.db.Model(&domain.ProjectTask{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No project tasks to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ProjectTask{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean project tasks: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d project tasks\n", result.RowsAffected)
	return nil
}

// CleanProjectAttachments 清理项目附件
func (c *Cleaner) CleanProjectAttachments() error {
	var count int64
	c.db.Model(&domain.ProjectAttachment{}).Count(&count)

	if count == 0 {
		fmt.Println("[INFO] No project attachments to clean")
		return nil
	}

	result := c.db.Where("1 = 1").Delete(&domain.ProjectAttachment{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean project attachments: %w", result.Error)
	}

	fmt.Printf("[OK] Cleaned %d project attachments\n", result.RowsAffected)
	return nil
}

// PrintStatistics 打印当前数据统计
func (c *Cleaner) PrintStatistics() {
	fmt.Println("\n========== 数据统计 ==========")

	tables := []struct {
		name  string
		model interface{}
	}{
		{"用户 (users)", &domain.User{}},
		{"项目 (projects)", &domain.Project{}},
		{"报告 (reports)", &domain.Report{}},
		{"文章 (articles)", &domain.Article{}},
		{"组织 (organizations)", &domain.Organization{}},
		{"头像 (avatars)", &domain.Avatar{}},
	}

	for _, t := range tables {
		var count int64
		c.db.Model(t.model).Count(&count)
		fmt.Printf("   %s: %d 条\n", t.name, count)
	}

	fmt.Println("==============================")
}
