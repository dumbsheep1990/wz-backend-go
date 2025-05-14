package main

import (
	"log"
	"os"
	"wz-backend-go/middleware"
	"wz-backend-go/services/site-service/handlers"

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

	// 公开路由
	apiGroup.GET("/site-templates", handlers.ListTemplates)
	apiGroup.GET("/site-templates/:id", handlers.GetTemplate)

	// 需要认证的路由
	authGroup := apiGroup.Group("/sites")
	authGroup.Use(middleware.Auth())
	{
		authGroup.GET("", handlers.ListSites)
		authGroup.GET("/:id", handlers.GetSite)
		authGroup.POST("", handlers.CreateSite)
		authGroup.PUT("/:id", handlers.UpdateSite)
		authGroup.DELETE("/:id", handlers.DeleteSite)
		authGroup.PUT("/:id/publish", handlers.PublishSite)
	}

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// 启动服务
	log.Printf("站点服务启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
