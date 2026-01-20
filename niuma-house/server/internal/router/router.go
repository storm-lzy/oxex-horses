package router

import (
	"niuma-house/internal/handler"
	"niuma-house/internal/middleware"
	"niuma-house/internal/ws"
	"niuma-house/pkg/config"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 公开 API
	api := r.Group("/api")
	{
		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// 职业分类
		api.GET("/occupations", handler.GetOccupations)

		// 需要认证的 API
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 用户
			protected.GET("/user/profile", handler.GetProfile)
			protected.PUT("/user/profile", handler.UpdateProfile)

			// 帖子
			protected.GET("/posts", handler.GetPosts)
			protected.GET("/posts/:id", handler.GetPost)
			protected.POST("/posts", handler.CreatePost)
			protected.PUT("/posts/:id", handler.UpdatePost)
			protected.DELETE("/posts/:id", handler.DeletePost)
			protected.POST("/posts/:id/like", handler.LikePost)
			protected.DELETE("/posts/:id/like", handler.UnlikePost)
			protected.POST("/posts/:id/favorite", handler.FavoritePost)
			protected.DELETE("/posts/:id/favorite", handler.UnfavoritePost)

			// 评论
			protected.GET("/posts/:id/comments", handler.GetComments)
			protected.POST("/posts/:id/comments", handler.CreateComment)
			protected.DELETE("/comments/:id", handler.DeleteComment)

			// 公司
			protected.GET("/companies", handler.GetCompanies)
			protected.GET("/companies/search", handler.SearchCompanies)
			protected.GET("/companies/:id", handler.GetCompany)
			protected.POST("/companies", handler.CreateCompany)

			// 上传
			protected.POST("/upload/presign", handler.GetPresignedURL)

			// 私信
			protected.GET("/messages", handler.GetMessages)
			protected.GET("/messages/unread", handler.GetUnreadCount)
			protected.POST("/messages/read", handler.MarkAsRead)
		}

		// WebSocket
		api.GET("/ws/chat", middleware.JWTAuth(), ws.HandleWebSocket)
	}

	// 管理后台 API
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		// 数据统计
		admin.GET("/dashboard/stats", handler.GetDashboardStats)

		// 用户管理
		admin.GET("/users", handler.AdminGetUsers)
		admin.POST("/users/:id/ban", handler.BanUser)
		admin.POST("/users/:id/unban", handler.UnbanUser)

		// 帖子管理
		admin.GET("/posts", handler.AdminGetPosts)
		admin.DELETE("/posts/:id", handler.AdminDeletePost)
		admin.POST("/posts/:id/top", handler.TopPost)

		// 公司管理
		admin.GET("/companies", handler.AdminGetCompanies)
		admin.DELETE("/companies/:id", handler.AdminDeleteCompany)
	}

	return r
}
