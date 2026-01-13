package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AvatarSeeder 头像测试数据填充器
type AvatarSeeder struct {
	db *gorm.DB
}

func NewAvatarSeeder(db *gorm.DB) *AvatarSeeder {
	return &AvatarSeeder{db: db}
}

// Seed 填充头像测试数据
func (s *AvatarSeeder) Seed(force bool) error {
	// 如果 force 为 true，则先清空现有数据
	if force {
		fmt.Println("[INFO] Force mode: clearing existing avatars...")
		s.db.Exec("DELETE FROM avatars")
	}

	// 检查是否已存在数据
	var count int64
	s.db.Model(&domain.Avatar{}).Count(&count)
	if count > 0 && !force {
		fmt.Printf("[INFO] Avatars already exist (%d), skipping seed. Use force=true to reseed.\n", count)
		return nil
	}

	fmt.Println("[INFO] Seeding avatar test data...")

	// 使用多种头像服务作为示例
	// 这些 URL 使用免费的头像占位符服务
	avatarData := []struct {
		Name string
		URL  string
	}{
		// 使用 DiceBear Avatars API (https://www.dicebear.com/)
		{"机器人头像 1", "https://api.dicebear.com/7.x/bottts/svg?seed=robot1"},
		{"机器人头像 2", "https://api.dicebear.com/7.x/bottts/svg?seed=robot2"},
		{"机器人头像 3", "https://api.dicebear.com/7.x/bottts/svg?seed=robot3"},
		{"机器人头像 4", "https://api.dicebear.com/7.x/bottts/svg?seed=robot4"},
		{"机器人头像 5", "https://api.dicebear.com/7.x/bottts/svg?seed=robot5"},

		// 像素风格头像
		{"像素头像 1", "https://api.dicebear.com/7.x/pixel-art/svg?seed=pixel1"},
		{"像素头像 2", "https://api.dicebear.com/7.x/pixel-art/svg?seed=pixel2"},
		{"像素头像 3", "https://api.dicebear.com/7.x/pixel-art/svg?seed=pixel3"},
		{"像素头像 4", "https://api.dicebear.com/7.x/pixel-art/svg?seed=pixel4"},
		{"像素头像 5", "https://api.dicebear.com/7.x/pixel-art/svg?seed=pixel5"},

		// 抽象头像
		{"抽象头像 1", "https://api.dicebear.com/7.x/shapes/svg?seed=shape1"},
		{"抽象头像 2", "https://api.dicebear.com/7.x/shapes/svg?seed=shape2"},
		{"抽象头像 3", "https://api.dicebear.com/7.x/shapes/svg?seed=shape3"},

		// 卡通人物头像
		{"卡通头像 1", "https://api.dicebear.com/7.x/adventurer/svg?seed=user1"},
		{"卡通头像 2", "https://api.dicebear.com/7.x/adventurer/svg?seed=user2"},
		{"卡通头像 3", "https://api.dicebear.com/7.x/adventurer/svg?seed=user3"},
		{"卡通头像 4", "https://api.dicebear.com/7.x/adventurer/svg?seed=user4"},
		{"卡通头像 5", "https://api.dicebear.com/7.x/adventurer/svg?seed=user5"},

		// 大头贴风格
		{"大头贴 1", "https://api.dicebear.com/7.x/big-smile/svg?seed=smile1"},
		{"大头贴 2", "https://api.dicebear.com/7.x/big-smile/svg?seed=smile2"},
		{"大头贴 3", "https://api.dicebear.com/7.x/big-smile/svg?seed=smile3"},

		// 黑客风格头像
		{"黑客头像 1", "https://api.dicebear.com/7.x/identicon/svg?seed=hacker1"},
		{"黑客头像 2", "https://api.dicebear.com/7.x/identicon/svg?seed=hacker2"},
		{"黑客头像 3", "https://api.dicebear.com/7.x/identicon/svg?seed=hacker3"},
		{"黑客头像 4", "https://api.dicebear.com/7.x/identicon/svg?seed=hacker4"},
		{"黑客头像 5", "https://api.dicebear.com/7.x/identicon/svg?seed=hacker5"},
	}

	successCount := 0
	for i, data := range avatarData {
		avatar := domain.Avatar{
			Name:      data.Name,
			URL:       data.URL,
			IsActive:  true,
			SortOrder: i + 1,
		}

		if err := s.db.Create(&avatar).Error; err != nil {
			log.Printf("[WARN] Failed to create avatar %s: %v", data.Name, err)
			continue
		}

		successCount++
		fmt.Printf("[OK] Created avatar: %s (ID: %d)\n", avatar.Name, avatar.ID)
	}

	fmt.Printf("\n[INFO] Seeded %d/%d avatars successfully\n", successCount, len(avatarData))

	// 打印统计
	s.printStatistics()

	return nil
}

// printStatistics 打印统计信息
func (s *AvatarSeeder) printStatistics() {
	fmt.Println("\n========== 头像统计 ==========")

	var total int64
	s.db.Model(&domain.Avatar{}).Count(&total)
	fmt.Printf("   总计: %d 个头像\n", total)

	var active int64
	s.db.Model(&domain.Avatar{}).Where("is_active = ?", true).Count(&active)
	fmt.Printf("   启用: %d 个\n", active)
	fmt.Printf("   禁用: %d 个\n", total-active)

	fmt.Println("===============================")
}
