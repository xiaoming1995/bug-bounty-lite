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

	fmt.Println("Bug Bounty Lite - Users Test Data Seeder")
	fmt.Println("=========================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 执行数据填充
	s := seeder.NewUserSeeder(db)
	if err := s.Seed(*forceFlag); err != nil {
		fmt.Printf("[ERROR] Seed failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Users test data seeded successfully!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Users Test Data Seeder")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed-users/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force    Force seed even if data exists (will skip existing users)")
	fmt.Println("  -help     Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed-users/main.go           # Seed users test data")
	fmt.Println("  go run cmd/seed-users/main.go -force    # Force seed (skip existing)")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make seed-users         # Seed users test data")
	fmt.Println("  make seed-users-force   # Force seed")
	fmt.Println("")
	fmt.Println("Test Users Created:")
	fmt.Println("  - whitehat_zhang / password123 (白帽子)")
	fmt.Println("  - whitehat_li / password123 (白帽子)")
	fmt.Println("  - whitehat_wang / password123 (白帽子)")
	fmt.Println("  - whitehat_zhao / password123 (白帽子)")
	fmt.Println("  - whitehat_chen / password123 (白帽子)")
	fmt.Println("  - vendor_test / password123 (厂商)")
	fmt.Println("  - admin_test / admin123 (管理员)")
}
