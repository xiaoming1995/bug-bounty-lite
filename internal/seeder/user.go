package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder 用户测试数据填充器
type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Seed 填充用户测试数据
func (s *UserSeeder) Seed(force bool) error {
	// 检查是否已存在测试用户
	var count int64
	s.db.Model(&domain.User{}).Where("username LIKE ?", "whitehat_%").Count(&count)

	if count > 0 && !force {
		fmt.Printf("[INFO] Test users already exist (%d found), skipping seed (use -force to override)\n", count)
		return nil
	}

	testUsers := []struct {
		Username string
		Password string
		Role     string
		Name     string
		Email    string
		Phone    string
	}{
		{"whitehat_zhang", "password123", "whitehat", "张三", "zhang@example.com", "13800138001"},
		{"whitehat_li", "password123", "whitehat", "李四", "li@example.com", "13800138002"},
		{"whitehat_wang", "password123", "whitehat", "王五", "wang@example.com", "13800138003"},
		{"whitehat_zhao", "password123", "whitehat", "赵六", "zhao@example.com", "13800138004"},
		{"whitehat_chen", "password123", "whitehat", "陈七", "chen@example.com", "13800138005"},
		{"vendor_test", "password123", "vendor", "测试厂商", "vendor@example.com", "13900139001"},
		{"admin_test", "admin123", "admin", "测试管理员", "admin@example.com", "13900139002"},
	}

	successCount := 0
	for _, u := range testUsers {
		var user domain.User
		if err := s.db.Where("username = ?", u.Username).First(&user).Error; err != nil {
			// 用户不存在，创建
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			user = domain.User{
				Username: u.Username,
				Password: string(hashedPassword),
				Role:     u.Role,
				Name:     u.Name,
				Email:    u.Email,
				Phone:    u.Phone,
			}
			if err := s.db.Create(&user).Error; err != nil {
				log.Printf("[WARN] Failed to create user %s: %v", u.Username, err)
				continue
			}
			successCount++
			fmt.Printf("[OK] Created: %s (%s) - %s | Password: %s\n", u.Username, u.Name, u.Role, u.Password)
		} else {
			fmt.Printf("[SKIP] User '%s' already exists (ID: %d)\n", u.Username, user.ID)
		}
	}

	fmt.Printf("\n[INFO] Seeded %d/%d users successfully\n", successCount, len(testUsers))

	// 打印统计
	s.printStatistics()

	return nil
}

// printStatistics 打印统计信息
func (s *UserSeeder) printStatistics() {
	fmt.Println("\n========== 用户统计 ==========")

	type roleStat struct {
		Role  string
		Count int64
	}
	var roleStats []roleStat
	s.db.Table("users").
		Select("role, count(*) as count").
		Group("role").
		Scan(&roleStats)

	roleMap := map[string]string{
		"whitehat": "白帽子",
		"vendor":   "厂商",
		"admin":    "管理员",
	}

	for _, stat := range roleStats {
		name := roleMap[stat.Role]
		if name == "" {
			name = stat.Role
		}
		fmt.Printf("   %s (%s): %d 个\n", stat.Role, name, stat.Count)
	}

	var total int64
	s.db.Model(&domain.User{}).Count(&total)
	fmt.Printf("\n   总计: %d 个用户\n", total)
	fmt.Println("===============================")
}
