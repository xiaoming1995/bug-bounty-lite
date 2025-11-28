package main

import (
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"bug-bounty-lite/internal/domain" // 引入 domain 包
	"bug-bounty-lite/internal/router"
	"fmt"
	"log"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化数据库
	db := database.InitDB(cfg)

	// 3. 自动迁移数据库 (建表)
	fmt.Println("🔄 Running Database Migrations...")
	err := db.AutoMigrate(&domain.User{}, &domain.Report{})
	if err != nil {
		log.Fatalf(" Migration failed: %v", err)
	}
	fmt.Println("✅ Database Migrations executed successfully")

	// 4. 初始化路由 (核心修复点！)
	// 这一步会将 Repo, Service, Handler 全部组装起来
	r := router.SetupRouter(db)

	// 5. 启动 HTTP 服务
	serverAddr := cfg.Server.Port
	fmt.Println("--------------------------------")
	fmt.Printf(" Bug Bounty Platform starting on %s ...\n", serverAddr)
	fmt.Println("--------------------------------")

	// r.Run() 会阻塞在这里监听端口，直到程序被关闭
	// 如果端口被占用或启动失败，会返回 error
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf(" Failed to start server: %v", err)
	}
}