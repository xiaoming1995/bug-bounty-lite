package router

import (
	"bug-bounty-lite/internal/handler"
	"bug-bounty-lite/internal/middleware"
	"bug-bounty-lite/internal/repository"
	"bug-bounty-lite/internal/service"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 设置可信代理（消除安全警告）
	// 本地开发设为 nil，生产环境设置为实际的代理 IP
	r.SetTrustedProxies(nil)

	// ===========================
	// 1. 全局中间件
	// ===========================
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// ===========================
	// 2. 初始化 JWT 管理器
	// ===========================
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expire)

	// ===========================
	// 3. 依赖注入 (组装层)
	// ===========================

	// User 模块
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, jwtManager)
	userHandler := handler.NewUserHandler(userService)

	// Report 模块
	reportRepo := repository.NewReportRepo(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// UserInfoChange 模块
	userInfoChangeRepo := repository.NewUserInfoChangeRepo(db)
	userInfoChangeService := service.NewUserInfoChangeService(userInfoChangeRepo)
	userInfoChangeHandler := handler.NewUserInfoChangeHandler(userInfoChangeService)

	// ===========================
	// 4. 注册路由
	// ===========================

	api := r.Group("/api/v1")
	{
		// 公开路由 - Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 需要认证的路由 - Reports
		reports := api.Group("/reports")
		reports.Use(middleware.AuthMiddleware(jwtManager))
		{
			reports.POST("", reportHandler.CreateHandler)    // 提交
			reports.GET("", reportHandler.ListHandler)       // 列表
			reports.GET("/:id", reportHandler.GetHandler)    // 详情
			reports.PUT("/:id", reportHandler.UpdateHandler) // 更新
		}

		// 需要认证的路由 - User Info Change
		userInfo := api.Group("/user/info")
		userInfo.Use(middleware.AuthMiddleware(jwtManager))
		{
			userInfo.POST("/change", userInfoChangeHandler.SubmitChangeRequest)   // 提交变更申请
			userInfo.GET("/changes", userInfoChangeHandler.GetUserChangeRequests) // 获取变更申请列表
			userInfo.GET("/changes/:id", userInfoChangeHandler.GetChangeRequest)  // 获取变更申请详情
		}
	}

	return r
}
