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
	forceFlag := flag.Bool("force", false, "Force seed even if data exists")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - All Test Data Seeder")
	fmt.Println("======================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 1. 填充项目数据
	fmt.Println("\n>>> Seeding Projects...")
	projectSeeder := seeder.NewProjectSeeder(db)
	if err := projectSeeder.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Project seed failed: %v\n", err)
		os.Exit(1)
	}

	// 2. 填充用户数据
	fmt.Println("\n>>> Seeding Users...")
	userSeeder := seeder.NewUserSeeder(db)
	if err := userSeeder.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] User seed failed: %v\n", err)
		os.Exit(1)
	}

	// 3. 填充报告数据
	fmt.Println("\n>>> Seeding Reports...")
	reportSeeder := seeder.NewReportSeeder(db)
	if err := reportSeeder.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Report seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] All test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - All Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-all/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will skip existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-all/main.go           # Seed all test data")
	fmt.Println("  go run cmd/seed-all/main.go -force    # Force seed (skip existing)")
}
