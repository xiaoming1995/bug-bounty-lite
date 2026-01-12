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
	// 1. 静态文件服务（用于访问上传的文件）
	// ===========================
	r.Static("/uploads", "./uploads")

	// ===========================
	// 2. 全局中间件
	// ===========================
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// 注入高级请求日志系统 (根据配置)
	if cfg.Server.EnableHttpLog {
		r.Use(middleware.HttpLogger())
	}

	// ===========================
	// 3. 初始化 JWT 管理器
	// ===========================
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expire)

	// ===========================
	// 4. 依赖注入 (组装层)
	// ===========================

	// User 模块
	userRepo := repository.NewUserRepo(db)
	orgRepo := repository.NewOrganizationRepo(db)
	userUpdateLogRepo := repository.NewUserUpdateLogRepo(db)

	organizationService := service.NewOrganizationService(orgRepo)
	organizationHandler := handler.NewOrganizationHandler(organizationService)

	userService := service.NewUserService(userRepo, orgRepo, userUpdateLogRepo, jwtManager)
	userHandler := handler.NewUserHandler(userService)

	// SystemConfig 模块（需要在 Report 之前初始化，因为 Report 依赖它）
	systemConfigRepo := repository.NewSystemConfigRepo(db)
	systemConfigService := service.NewSystemConfigService(systemConfigRepo)
	systemConfigHandler := handler.NewSystemConfigHandler(systemConfigService)

	// Report 模块
	reportRepo := repository.NewReportRepo(db)
	reportService := service.NewReportService(reportRepo, systemConfigRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// UserInfoChange 模块
	userInfoChangeRepo := repository.NewUserInfoChangeRepo(db)
	userInfoChangeService := service.NewUserInfoChangeService(userInfoChangeRepo)
	userInfoChangeHandler := handler.NewUserInfoChangeHandler(userInfoChangeService)

	// Project 模块
	projectRepo := repository.NewProjectRepo(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// Upload 模块
	uploadHandler := handler.NewUploadHandler()

	// ===========================
	// 5. 注册路由
	// ===========================

	api := r.Group("/api/v1")
	{
		// 公开路由 - Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 用户个人管理路由
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(jwtManager))
		{
			user.POST("/profile", userHandler.UpdateProfile)
			user.POST("/bind-org", userHandler.BindOrganization)
		}

		// 组织管理路由（限管理员可用逻辑待后续细化，目前先挂载）
		orgs := api.Group("/organizations")
		orgs.Use(middleware.AuthMiddleware(jwtManager))
		{
			orgs.POST("", organizationHandler.Create)
			orgs.GET("", organizationHandler.List)
			orgs.PUT("/:id", organizationHandler.Update)
			orgs.DELETE("/:id", organizationHandler.Delete)
		}

		// 需要认证的路由 - Reports
		reports := api.Group("/reports")
		reports.Use(middleware.AuthMiddleware(jwtManager))
		{
			reports.POST("", reportHandler.CreateHandler)              // 提交
			reports.GET("", reportHandler.ListHandler)                 // 列表
			reports.GET("/:id", reportHandler.GetHandler)              // 详情
			reports.PUT("/:id", reportHandler.UpdateHandler)           // 更新
			reports.DELETE("/:id", reportHandler.DeleteHandler)        // 软删除
			reports.POST("/:id/restore", reportHandler.RestoreHandler) // 恢复已删除
		}

		// 需要认证的路由 - User Info Change
		userInfo := api.Group("/user/info")
		userInfo.Use(middleware.AuthMiddleware(jwtManager))
		{
			userInfo.POST("/change", userInfoChangeHandler.SubmitChangeRequest)   // 提交变更申请
			userInfo.GET("/changes", userInfoChangeHandler.GetUserChangeRequests) // 获取变更申请列表
			userInfo.GET("/changes/:id", userInfoChangeHandler.GetChangeRequest)  // 获取变更申请详情
		}

		// 需要认证的路由 - Projects
		projects := api.Group("/projects")
		projects.Use(middleware.AuthMiddleware(jwtManager))
		{
			projects.POST("", projectHandler.CreateHandler)              // 创建项目（仅admin）
			projects.GET("", projectHandler.ListHandler)                 // 获取项目列表
			projects.GET("/:id", projectHandler.GetHandler)              // 获取项目详情
			projects.PUT("/:id", projectHandler.UpdateHandler)           // 更新项目（仅admin）
			projects.DELETE("/:id", projectHandler.DeleteHandler)        // 软删除项目（仅admin）
			projects.POST("/:id/restore", projectHandler.RestoreHandler) // 恢复已删除项目（仅admin）
		}

		// 需要认证的路由 - System Configs
		configs := api.Group("/configs")
		configs.Use(middleware.AuthMiddleware(jwtManager))
		{
			configs.GET("/:type", systemConfigHandler.GetConfigsByTypeHandler)    // 获取配置列表
			configs.GET("/:type/:id", systemConfigHandler.GetConfigHandler)       // 获取配置详情
			configs.POST("/:type", systemConfigHandler.CreateConfigHandler)       // 创建配置（仅admin）
			configs.PUT("/:type/:id", systemConfigHandler.UpdateConfigHandler)    // 更新配置（仅admin）
			configs.DELETE("/:type/:id", systemConfigHandler.DeleteConfigHandler) // 删除配置（仅admin）
		}

		// 需要认证的路由 - Upload
		upload := api.Group("/upload")
		upload.Use(middleware.AuthMiddleware(jwtManager))
		{
			upload.POST("", uploadHandler.UploadFileHandler) // 上传文件
		}
	}

	return r
}
