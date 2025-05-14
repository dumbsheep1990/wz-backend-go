package main

import (
	"log"
	"os"
	"wz-backend-go/middleware"
	"wz-backend-go/services/component-service/handlers"

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

	// 组件库API - 不需要认证
	apiGroup.GET("/components/categories", handlers.ListComponentCategories)
	apiGroup.GET("/components/:type", handlers.GetComponentDefinition)

	// 组件实例API - 需要认证
	componentGroup := apiGroup.Group("")
	componentGroup.Use(middleware.Auth())
	{
		instanceGroup := componentGroup.Group("/sites/:siteId/pages/:pageId/sections/:sectionId/components")
		{
			instanceGroup.POST("", handlers.AddComponent)
			instanceGroup.PUT("/:id", handlers.UpdateComponent)
			instanceGroup.DELETE("/:id", handlers.DeleteComponent)
			instanceGroup.PUT("/reorder", handlers.ReorderComponents)
		}
	}

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// 启动服务
	log.Printf("组件服务启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
