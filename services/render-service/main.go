package main

import (
	"log"
	"os"
	"wz-backend-go/middleware"
	"wz-backend-go/services/render-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// 创建Gin引擎
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由
	apiGroup := r.Group("/api/v1")

	// 预览路由 - 需要认证
	previewGroup := apiGroup.Group("/preview")
	previewGroup.Use(middleware.Auth())
	{
		// 预览整个站点
		previewGroup.GET("/sites/:siteId", handlers.PreviewSite)
		// 预览单个页面
		previewGroup.GET("/sites/:siteId/pages/:pageId", handlers.PreviewPage)
	}

	// 公开访问路由 - 不需要认证
	renderGroup := r.Group("/render")
	{
		// 根据域名渲染站点
		renderGroup.GET("/site", handlers.RenderSiteByDomain)
		// 渲染特定站点的页面
		renderGroup.GET("/sites/:siteId/:slug", handlers.RenderPageBySlug)
	}

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// 启动服务
	log.Printf("渲染服务启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
