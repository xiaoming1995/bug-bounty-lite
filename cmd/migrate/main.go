package main

import (
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"bug-bounty-lite/pkg/migrate"
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== Bug Bounty Lite Database Migration Tool ===")

	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化数据库连接
	db := database.InitDB(cfg)

	// 3. 创建迁移器
	migrator := migrate.NewMigrator(db)

	// 4. 执行迁移
	fmt.Println("[STEP] Running Migrations...")
	if err := migrator.Run(); err != nil {
		log.Fatalf("[FATAL] Migration failed: %v", err)
	}

	// 5. 打印状态
	fmt.Println("[STEP] Verifying Status...")
	migrator.Status()

	fmt.Println("\n[SUCCESS] All migrations completed successfully!")
}
