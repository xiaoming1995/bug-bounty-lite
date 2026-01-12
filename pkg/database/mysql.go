package database

import (
	"bug-bounty-lite/pkg/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
// 参数: cfg 是我们上一节加载好的配置对象
// 返回: *gorm.DB 数据库连接句柄
func InitDB(cfg *config.Config) *gorm.DB {
	// 1. 构造 DSN (Data Source Name)
	// 这里直接使用配置中的 DSN 字符串
	dsn := cfg.Database.DSN

	// 2. 配置 Gorm 的日志模式
	// Debug 模式下会打印所有 SQL 语句，方便调试
	// 生产环境建议改为 logger.Error 只打印错误
	var gormLogger logger.Interface
	if cfg.Server.Mode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// 3. 打开连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		// 如果连不上数据库，整个平台无法运行，直接 Panic
		panic(fmt.Errorf("Fatal error connecting to database: %w", err))
	}

	// 4. 获取底层的 sql.DB 对象，用于设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("Failed to get sql.DB: %w", err))
	}

	// 5. 设置连接池参数 (非常重要！)
	// SetMaxIdleConns: 空闲连接池中保留的最大连接数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)

	// SetMaxOpenConns: 数据库打开的最大连接数
	// 如果请求超过这个数量，新的请求会阻塞等待
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpen)

	// SetConnMaxLifetime: 连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")

	return db
}
