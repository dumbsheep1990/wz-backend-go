package main

import (
	"log"
	"os"
	"wz-backend-go/middleware"
	"wz-backend-go/services/component-service/handlers"
	"wz-backend-go/services/component-service/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// 初始化数据库连接
	dsn := "root:password@tcp(127.0.0.1:3306)/component_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := service.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 启动服务
	log.Printf("组件服务启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
