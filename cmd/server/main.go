package main

import (
	"bug-bounty-lite/internal/router"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"bug-bounty-lite/pkg/migrate"
	"flag"
	"fmt"
	"log"
)

func main() {
	// 命令行参数
	migrateFlag := flag.Bool("migrate", false, "Run database migrations before starting server")
	flag.Parse()

	fmt.Println("Bug Bounty Platform")
	fmt.Println("======================")

	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化数据库
	db := database.InitDB(cfg)

	// 3. 可选：执行数据库迁移
	if *migrateFlag {
		migrator := migrate.NewMigrator(db)
		if err := migrator.Run(); err != nil {
			log.Fatalf("[ERROR] Migration failed: %v", err)
		}
	} else {
		fmt.Println("[INFO] Skipping migrations (use --migrate to run)")
	}

	// 4. 初始化路由
	// 这一步会将 Repo, Service, Handler, Middleware 全部组装起来
	r := router.SetupRouter(db, cfg)

	// 5. 启动 HTTP 服务
	serverAddr := cfg.Server.Port
	fmt.Println("--------------------------------")
	fmt.Printf("[INFO] Server starting on %s ...\n", serverAddr)
	fmt.Println("--------------------------------")

	// r.Run() 会阻塞在这里监听端口，直到程序被关闭
	// 如果端口被占用或启动失败，会返回 error
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}
