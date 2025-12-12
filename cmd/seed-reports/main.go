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

	fmt.Println("Bug Bounty Lite - Reports Test Data Seeder")
	fmt.Println("==========================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 执行数据填充
	s := seeder.NewReportSeeder(db)
	if err := s.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Reports test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Reports Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-reports/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will skip existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-reports/main.go           # Seed reports test data")
	fmt.Println("  go run cmd/seed-reports/main.go -force    # Force seed (skip existing)")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make seed-reports         # Seed reports test data")
	fmt.Println("  make seed-reports-force   # Force seed")
	fmt.Println("")
	fmt.Println("Prerequisites:")
	fmt.Println("  1. Run 'make migrate' to create database tables")
	fmt.Println("  2. Run 'make seed-projects' to create test projects")
	fmt.Println("  3. Run 'make init' to initialize severity levels")
	fmt.Println("  4. Run 'make seed-users' to create test users")
}
