package main

import (
	"bug-bounty-lite/internal/seeder"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 命令行参数
	helpFlag := flag.Bool("help", false, "Show help message")
	forceFlag := flag.Bool("force", false, "Force seed even if data exists (will clear existing data)")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - Avatar Test Data Seeder")
	fmt.Println("==========================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 填充头像数据
	fmt.Println("\n>>> Seeding Avatars...")
	avatarSeeder := seeder.NewAvatarSeeder(db)
	if err := avatarSeeder.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Avatar seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Avatar test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Avatar Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-avatars/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will clear existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-avatars/main.go           # Seed avatar test data")
	fmt.Println("  go run cmd/seed-avatars/main.go -force    # Force reseed (clear and seed)")
}
