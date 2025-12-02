package main

import (
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"bug-bounty-lite/pkg/migrate"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 命令行参数
	statusFlag := flag.Bool("status", false, "Show migration status")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - Database Migration Tool")
	fmt.Println("============================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 创建迁移器
	migrator := migrate.NewMigrator(db)

	// 根据参数执行不同操作
	if *statusFlag {
		// 只显示状态
		migrator.Status()
	} else {
		// 执行迁移
		if err := migrator.Run(); err != nil {
			fmt.Printf("[ERROR] Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\n[OK] Migration completed successfully!")
	}
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Database Migration Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -status    Show current migration status")
	fmt.Println("  -help      Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go           # Run migrations")
	fmt.Println("  go run cmd/migrate/main.go -status   # Check status")
	fmt.Println("")
	fmt.Println("Or use Makefile:")
	fmt.Println("  make migrate         # Run migrations")
	fmt.Println("  make migrate-status  # Check status")
}
