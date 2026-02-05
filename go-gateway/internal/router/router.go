package router

import (
	"github.com/gin-gonic/gin"
	"go-gateway/internal/config"
	"go-gateway/internal/grpcclient"
	"go-gateway/internal/handler"
	"go-gateway/internal/middleware"
	"go-gateway/internal/service"
	"go-gateway/pkg/logger"
)

func Register(r *gin.Engine, cfg *config.Config, log *logger.Logger, authService *service.AuthService) {
	// 全局中间件
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(log))
	r.Use(middleware.RateLimit(cfg.RateLimit))
	r.Use(middleware.CORS())

	// Handler 初始化
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	projectService := service.NewProjectService()
	projectHandler := handler.NewProjectHandler(projectService)
	chapterService := service.NewChapterService()
	chapterHandler := handler.NewChapterHandler(chapterService, projectService)
	versionService := service.NewVersionService()
	versionHandler := handler.NewVersionHandler(versionService)

	// AI 服务
	aiSvc := service.NewAIService(cfg.GRPC, log)
	streamHandler := handler.NewStreamHandler(aiSvc)
	chapterGenerateHandler := handler.NewChapterGenerateHandler(aiSvc, chapterService, projectService)

	// AI gRPC 客户端 (用于质量评估)
	aiClient := grpcclient.NewAIClient(cfg.GRPC, log)
	qualityHandler := handler.NewQualityHandler(aiClient, log.Logger)
	constraintHandler := handler.NewConstraintHandler(aiClient, log.Logger)

	// 世界观服务
	characterService := service.NewCharacterService()
	characterHandler := handler.NewCharacterHandler(characterService)
	locationService := service.NewLocationService()
	locationHandler := handler.NewLocationHandler(locationService)
	orgService := service.NewOrganizationService()
	orgHandler := handler.NewOrganizationHandler(orgService)
	worldSettingService := service.NewWorldSettingService()
	worldSettingHandler := handler.NewWorldSettingHandler(worldSettingService)

	// 伏笔服务
	foreshadowService := service.NewForeshadowService()
	foreshadowHandler := handler.NewForeshadowHandler(foreshadowService)

	// 用户布局配置服务
	userLayoutService := service.NewUserLayoutService()
	userLayoutHandler := handler.NewUserLayoutHandler(userLayoutService)

	// 文件上传服务
	fileUploadService := service.NewFileUploadService()
	fileUploadHandler := handler.NewFileUploadHandler(fileUploadService)

	// 批量操作服务
	batchOperationService := service.NewBatchOperationService()
	batchOperationHandler := handler.NewBatchOperationHandler(batchOperationService, projectService)

	// 搜索服务
	searchService := service.NewSearchService()
	searchHandler := handler.NewSearchHandler(searchService)

	// OAuth 服务
	oauthService := service.NewOAuthService(&cfg.OAuth, authService)
	oauthHandler := handler.NewOAuthHandler(oauthService)

	// 健康检查
	r.GET("/healthz", healthHandler.Healthz)

	// 公开接口
	public := r.Group("/api/v1")
	{
		// 认证
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.Refresh)
		public.POST("/auth/password-reset/request", authHandler.RequestPasswordReset)
		public.POST("/auth/password-reset/verify", authHandler.ResetPassword)

		// OAuth
		public.GET("/oauth/providers", oauthHandler.GetProviders)
		public.GET("/oauth/:provider/authorize", oauthHandler.Authorize)
		public.GET("/oauth/:provider/callback", oauthHandler.Callback)
		public.POST("/oauth/:provider/callback", oauthHandler.CallbackPost)
	}

	// 需认证接口
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth(cfg.Auth))
	{
		// 当前用户
		protected.GET("/auth/me", authHandler.GetCurrentUser)

		// 搜索
		protected.GET("/search", searchHandler.GlobalSearch)

		// 用户布局配置
		protected.GET("/user/layout-configs", userLayoutHandler.GetAllLayoutConfigs)
		protected.GET("/user/layout-config/:layoutType", userLayoutHandler.GetLayoutConfig)
		protected.POST("/user/layout-config/:layoutType", userLayoutHandler.SaveLayoutConfig)
		protected.DELETE("/user/layout-config/:layoutType", userLayoutHandler.DeleteLayoutConfig)
		protected.POST("/user/layout-config/:layoutType/reset", userLayoutHandler.ResetLayoutConfig)

		// 文件上传
		protected.POST("/upload", fileUploadHandler.Upload)
		protected.GET("/files", fileUploadHandler.ListFiles)
		protected.GET("/files/stats", fileUploadHandler.GetStorageStats)
		protected.GET("/files/:id", fileUploadHandler.GetFile)
		protected.DELETE("/files/:id", fileUploadHandler.DeleteFile)
		protected.GET("/files/:id/download", fileUploadHandler.DownloadFile)
		protected.GET("/files/:id/thumbnail", fileUploadHandler.GetThumbnail)

		// 项目
		protected.GET("/projects", projectHandler.List)
		protected.POST("/projects", projectHandler.Create)
		protected.GET("/projects/:id", projectHandler.Get)
		protected.PUT("/projects/:id", projectHandler.Update)
		protected.DELETE("/projects/:id", projectHandler.Delete)
		protected.POST("/projects/:id/restore", projectHandler.Restore)
		protected.PUT("/projects/:id/metadata", projectHandler.UpdateMetadata)
		protected.GET("/projects/:id/statistics", projectHandler.GetStatistics)
		protected.GET("/projects/:id/search", searchHandler.ProjectSearch)
		protected.POST("/projects/:id/advanced-filter", searchHandler.AdvancedFilter)

		// 章节 (项目下)
		protected.GET("/projects/:id/chapters", chapterHandler.List)
		protected.POST("/projects/:id/chapters", chapterHandler.Create)
		protected.POST("/projects/:id/chapters/reorder", chapterHandler.Reorder)
		protected.POST("/projects/:id/chapters/batch-update", batchOperationHandler.BatchUpdate)
		protected.POST("/projects/:id/chapters/batch-delete", batchOperationHandler.BatchDelete)
		protected.POST("/projects/:id/chapters/batch-status", batchOperationHandler.BatchStatusUpdate)

		// 章节 (独立)
		protected.GET("/chapters/:id", chapterHandler.Get)
		protected.PUT("/chapters/:id", chapterHandler.Update)
		protected.DELETE("/chapters/:id", chapterHandler.Delete)

		// 版本管理
		protected.GET("/chapters/:id/versions", versionHandler.ListVersions)
		protected.GET("/chapters/:id/versions/diff", versionHandler.DiffVersions)
		protected.GET("/chapters/:id/versions/:n", versionHandler.GetVersion)
		protected.POST("/chapters/:id/versions/:n/restore", versionHandler.RestoreVersion)

		// 章节 AI 生成 (流式)
		protected.POST("/chapters/:id/generate-stream", chapterGenerateHandler.GenerateStream)
		protected.POST("/chapters/:id/partial-regenerate-stream", chapterGenerateHandler.PartialRegenerateStream)
		protected.POST("/chapters/:id/polish-stream", chapterGenerateHandler.PolishStream)

		// 章节质量评估
		protected.POST("/chapters/:id/analyze-reading-power", qualityHandler.AnalyzeReadingPower)
		protected.POST("/chapters/:id/check-consistency", qualityHandler.CheckConsistency)
		protected.POST("/chapters/:id/multi-agent-review", qualityHandler.MultiAgentReview)
		protected.POST("/chapters/:id/evaluate-quality", qualityHandler.EvaluateQuality)

		// 约束检查与豁免
		protected.POST("/projects/:id/chapters/:chapterId/check-constraints", constraintHandler.CheckConstraints)
		protected.POST("/projects/:id/chapters/:chapterId/exemptions", constraintHandler.RequestExemption)
		protected.DELETE("/exemptions/:exemptionId", constraintHandler.RevokeExemption)

		// 角色 (项目下)
		protected.GET("/projects/:id/characters", characterHandler.List)
		protected.POST("/projects/:id/characters", characterHandler.Create)

		// 角色 (独立)
		protected.GET("/characters/:id", characterHandler.Get)
		protected.PUT("/characters/:id", characterHandler.Update)
		protected.DELETE("/characters/:id", characterHandler.Delete)

		// 角色关系
		protected.GET("/characters/:id/relationships", characterHandler.GetRelationships)
		protected.POST("/characters/:id/relationships", characterHandler.CreateRelationship)
		protected.DELETE("/relationships/:id", characterHandler.DeleteRelationship)

		// 角色经历
		protected.GET("/characters/:id/experiences", characterHandler.GetExperiences)
		protected.POST("/characters/:id/experiences", characterHandler.CreateExperience)
		protected.DELETE("/experiences/:id", characterHandler.DeleteExperience)

		// 地点 (项目下)
		protected.GET("/projects/:id/locations", locationHandler.List)
		protected.GET("/projects/:id/locations/tree", locationHandler.GetTree)
		protected.POST("/projects/:id/locations", locationHandler.Create)

		// 地点 (独立)
		protected.GET("/locations/:id", locationHandler.Get)
		protected.PUT("/locations/:id", locationHandler.Update)
		protected.DELETE("/locations/:id", locationHandler.Delete)

		// 组织 (项目下)
		protected.GET("/projects/:id/organizations", orgHandler.List)
		protected.POST("/projects/:id/organizations", orgHandler.Create)

		// 组织 (独立)
		protected.GET("/organizations/:id", orgHandler.Get)
		protected.PUT("/organizations/:id", orgHandler.Update)
		protected.DELETE("/organizations/:id", orgHandler.Delete)

		// 组织成员
		protected.GET("/organizations/:id/members", orgHandler.GetMembers)
		protected.POST("/organizations/:id/members", orgHandler.AddMember)
		protected.DELETE("/members/:id", orgHandler.RemoveMember)

		// 世界设定 (项目下)
		protected.GET("/projects/:id/world-settings", worldSettingHandler.List)
		protected.GET("/projects/:id/world-settings/tree", worldSettingHandler.GetTree)
		protected.GET("/projects/:id/world-settings/category/:category", worldSettingHandler.GetByCategory)
		protected.POST("/projects/:id/world-settings", worldSettingHandler.Create)

		// 世界设定 (独立)
		protected.GET("/world-settings/:id", worldSettingHandler.Get)
		protected.PUT("/world-settings/:id", worldSettingHandler.Update)
		protected.DELETE("/world-settings/:id", worldSettingHandler.Delete)

		// 伏笔 (项目下)
		protected.GET("/projects/:id/foreshadows", foreshadowHandler.List)
		protected.POST("/projects/:id/foreshadows", foreshadowHandler.Create)
		protected.GET("/projects/:id/foreshadows/stats", foreshadowHandler.GetStats)
		protected.GET("/projects/:id/foreshadows/pending", foreshadowHandler.GetPendingReminders)
		protected.GET("/projects/:id/foreshadows/overdue", foreshadowHandler.GetOverdueForeshadows)
		protected.POST("/projects/:id/foreshadows/check-reminders", foreshadowHandler.CheckAndCreateReminders)

		// 伏笔提醒 (项目下)
		protected.GET("/projects/:id/foreshadow-reminders/unread", foreshadowHandler.GetUnreadReminders)
		protected.POST("/projects/:id/foreshadow-reminders/mark-all-read", foreshadowHandler.MarkAllRemindersAsRead)

		// 伏笔 (独立)
		protected.GET("/foreshadows/:id", foreshadowHandler.Get)
		protected.PUT("/foreshadows/:id", foreshadowHandler.Update)
		protected.DELETE("/foreshadows/:id", foreshadowHandler.Delete)
		protected.POST("/foreshadows/:id/resolve", foreshadowHandler.Resolve)

		// 伏笔暗示
		protected.GET("/foreshadows/:id/hints", foreshadowHandler.GetHints)
		protected.POST("/foreshadows/:id/hints", foreshadowHandler.AddHint)
		protected.DELETE("/hints/:id", foreshadowHandler.DeleteHint)

		// 伏笔提醒
		protected.POST("/foreshadows/:id/reminders", foreshadowHandler.CreateReminder)
		protected.POST("/reminders/:id/read", foreshadowHandler.MarkReminderAsRead)
	}

	// 流式接口
	stream := r.Group("/api/v1/stream")
	stream.Use(middleware.JWTAuth(cfg.Auth))
	{
		stream.GET("/generate", streamHandler.Generate)
	}
}
