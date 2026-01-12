package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"
	"math/rand"
	"time"

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

// Seed 填充用户测试数据（追加模式）
func (s *UserSeeder) Seed(force bool) error {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 用户模板
	userTemplates := []struct {
		RolePrefix string
		Role       string
		NamePrefix string
	}{
		{"whitehat", "whitehat", "白帽子"},
		{"whitehat", "whitehat", "安全研究员"},
		{"whitehat", "whitehat", "渗透测试"},
		{"vendor", "vendor", "厂商"},
		{"admin", "admin", "管理员"},
	}

	// 获取所有可用的组织 ID（用于绑定）
	var orgIDs []uint
	s.db.Model(&domain.Organization{}).Pluck("id", &orgIDs)
	if len(orgIDs) == 0 {
		fmt.Println("[WARN] No organizations found. Users will be created without organization binding.")
	}

	// 生成 5-10 个新用户（每次执行都生成新的）
	numUsers := rand.Intn(6) + 5
	fmt.Printf("[INFO] Generating %d new users...\n", numUsers)

	successCount := 0
	for i := 0; i < numUsers; i++ {
		template := userTemplates[rand.Intn(len(userTemplates))]
		timestamp := time.Now().UnixNano()
		suffix := fmt.Sprintf("%d", timestamp/1000000+int64(i))

		username := fmt.Sprintf("%s_%s", template.RolePrefix, suffix)
		password := "password123"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		bios := []string{
			"热爱安全研究，专注于 Web 漏洞挖掘。",
			"资深渗透测试工程师，CTF 爱好者。",
			"崇尚极简主义，代码拯救世界。",
			"关注业务安全，保护用户隐私。",
			"路漫漫其修远兮，吾将上下而求索。",
		}

		user := domain.User{
			Username: username,
			Password: string(hashedPassword),
			Role:     template.Role,
			Name:     fmt.Sprintf("%s%d", template.NamePrefix, i+1),
			Email:    fmt.Sprintf("%s@example.com", username),
			Phone:    fmt.Sprintf("138%08d", rand.Intn(100000000)),
			Bio:      bios[rand.Intn(len(bios))],
		}

		// 随机分配一个组织 (如果存在组织数据)
		if len(orgIDs) > 0 {
			user.OrgID = orgIDs[rand.Intn(len(orgIDs))]
		}

		if err := s.db.Create(&user).Error; err != nil {
			log.Printf("[WARN] Failed to create user %s: %v", user.Username, err)
			continue
		}

		successCount++
		fmt.Printf("[OK] Created: %s (%s) - %s | Password: %s\n", user.Username, user.Name, user.Role, password)

		// 短暂延迟确保时间戳不同
		time.Sleep(time.Millisecond)
	}

	fmt.Printf("\n[INFO] Seeded %d/%d users successfully\n", successCount, numUsers)

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
