package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/infrastructure/factory"
	"wz-backend-go/internal/infrastructure/mock"
)

func main() {
	// 设置Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin路由
	router := gin.Default()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
	// 创建Mock服务客户端，实际应用中应该使用真实的服务客户端
	siteService := mock.NewMockSiteService()
	pageService := mock.NewMockPageService()

	// 创建渲染服务工厂
	renderServiceFactory := factory.NewRenderServiceFactory(
		siteService,
		pageService,
	)

	// 注册路由
	renderServiceFactory.RegisterRoutes(router)

	// 启动缓存清理任务
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	renderServiceFactory.StartCacheCleaner(ctx)

	// 设置端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // 默认端口
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 优雅关闭
	go func() {
		// 监听中断信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭服务器...")

		// 设置关闭超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("服务器关闭错误:", err)
		}
	}()

	// 启动服务器
	log.Printf("渲染服务启动在 http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("服务器启动错误:", err)
	}
}
