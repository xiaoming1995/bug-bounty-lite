package router

import (
	"bug-bounty-lite/internal/handler"
	"bug-bounty-lite/internal/repository"
	"bug-bounty-lite/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// ===========================
	// 1. 依赖注入 (组装层)
	// ===========================

	// User 模块
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Report 模块
	reportRepo := repository.NewReportRepo(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// ===========================
	// 2. 注册路由
	// ===========================

	api := r.Group("/api/v1")
	{
		// Auth 路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Report 路由
		reports := api.Group("/reports")
		{
			reports.POST("", reportHandler.CreateHandler)      // 提交
			reports.GET("", reportHandler.ListHandler)         // 列表
			reports.GET("/:id", reportHandler.GetHandler)      // 详情
		}
	}

	return r
}