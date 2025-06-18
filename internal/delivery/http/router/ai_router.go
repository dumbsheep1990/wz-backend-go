package router

import (
	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/delivery/http/handler/ai"
	"wz-backend-go/internal/middleware"
)

// RegisterAIRoutes 注册AI服务的所有路由
func RegisterAIRoutes(
	router *gin.Engine,
	aiHandler *ai.AIHandler,
) {
	api := router.Group("/api/v1")
	
	// 应用JWT中间件
	api.Use(middleware.JWTAuthMiddleware())
	
	// AI服务路由
	aiGroup := api.Group("/ai")
	{
		// 推荐相关
		aiGroup.POST("/recommend", aiHandler.Recommend)
		aiGroup.GET("/recommend", aiHandler.GetRecommendByScene)
		
		// 内容审核
		aiGroup.POST("/content/review", aiHandler.ContentReview)
		
		// 客服对话
		aiGroup.POST("/chat", aiHandler.Chat)
	}
	
	// 健康检查路由（不需要认证）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "ai",
			"version": "1.0.0",
		})
	})
} 