package router

import (
	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/delivery/http/handler/content"
	"wz-backend-go/internal/middleware"
)

// RegisterContentRoutes 注册内容服务的所有路由
func RegisterContentRoutes(
	router *gin.Engine,
	categoryHandler *content.CategoryHandler,
	postHandler *content.PostHandler,
) {
	api := router.Group("/api/v1")
	
	// 应用JWT中间件
	api.Use(middleware.JWTAuthMiddleware())
	
	// 分类路由
	categories := api.Group("/categories")
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:category_id", categoryHandler.UpdateCategory)
		categories.DELETE("/:category_id", categoryHandler.DeleteCategory)
		categories.GET("/:category_id", categoryHandler.GetCategory)
		categories.GET("", categoryHandler.ListCategories)
	}
	
	// 帖子路由
	posts := api.Group("/posts")
	{
		posts.POST("", postHandler.CreatePost)
		posts.PUT("/:post_id", postHandler.UpdatePost)
		posts.DELETE("/:post_id", postHandler.DeletePost)
		posts.GET("/:post_id", postHandler.GetPost)
		posts.GET("", postHandler.ListPosts)
	}
	
	// 健康检查路由（不需要认证）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "content",
			"version": "1.0.0",
		})
	})
} 