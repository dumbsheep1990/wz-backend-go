package main

import (
	"log"
	"os"
	"wz-backend-go/middleware"
	"wz-backend-go/services/page-service/handlers"

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
	apiGroup.Use(middleware.Auth())

	// 页面相关路由
	pageGroup := apiGroup.Group("/sites/:siteId/pages")
	{
		pageGroup.GET("", handlers.ListPages)
		pageGroup.GET("/:id", handlers.GetPage)
		pageGroup.POST("", handlers.CreatePage)
		pageGroup.PUT("/:id", handlers.UpdatePage)
		pageGroup.DELETE("/:id", handlers.DeletePage)
		pageGroup.PUT("/:id/homepage", handlers.SetHomepage)
		pageGroup.PUT("/reorder", handlers.ReorderPages)
	}

	// 区块相关路由
	sectionGroup := apiGroup.Group("/sites/:siteId/pages/:pageId/sections")
	{
		sectionGroup.GET("", handlers.ListSections)
		sectionGroup.POST("", handlers.AddSection)
		sectionGroup.PUT("/:id", handlers.UpdateSection)
		sectionGroup.DELETE("/:id", handlers.DeleteSection)
		sectionGroup.PUT("/reorder", handlers.ReorderSections)
	}

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// 启动服务
	log.Printf("页面服务启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
