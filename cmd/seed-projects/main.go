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

	fmt.Println("Bug Bounty Lite - Projects Test Data Seeder")
	fmt.Println("=============================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 执行数据填充
	s := seeder.NewProjectSeeder(db)
	if err := s.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Projects test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Projects Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-projects/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will skip existing data)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-projects/main.go           # Seed projects test data")
	fmt.Println("  go run cmd/seed-projects/main.go -force    # Force seed (skip existing)")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make seed-projects         # Seed projects test data")
	fmt.Println("  make seed-projects-force   # Force seed")
}
