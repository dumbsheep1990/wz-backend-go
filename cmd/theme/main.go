package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	
	"github.com/wz-backend-go/internal/delivery/http/handler/theme"
	themeRPC "github.com/wz-backend-go/internal/delivery/rpc/theme"
	"github.com/wz-backend-go/internal/middleware"
	pb "wz-backend-go/api/rpc/theme"
)

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 创建HTTP处理器
	themeHandler := theme.NewThemeHandler()

	// 注册HTTP路由
	setupRoutes(r, themeHandler)

	// 启动HTTP服务器
	httpPort := getEnv("THEME_HTTP_PORT", "8011")
	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: r,
	}

	// 启动gRPC服务器
	grpcPort := getEnv("THEME_GRPC_PORT", "50011")
	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Theme gRPC服务监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	themeRPCServer := themeRPC.NewThemeServer()
	pb.RegisterThemeServiceServer(grpcServer, themeRPCServer)

	// 启动HTTP服务
	go func() {
		log.Printf("Theme HTTP服务启动在端口 %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Theme HTTP服务启动失败: %v", err)
		}
	}()

	// 启动gRPC服务
	go func() {
		log.Printf("Theme gRPC服务启动在端口 %s", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Theme gRPC服务启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭Theme服务...")

	// 5秒内优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Theme HTTP服务关闭失败: %v", err)
	}

	// 关闭gRPC服务
	grpcServer.GracefulStop()

	log.Println("Theme服务已退出")
}

// setupRoutes 设置路由
func setupRoutes(r *gin.Engine, themeHandler *theme.ThemeHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "theme-service",
			"time":    time.Now().Unix(),
		})
	})

	// API路由
	apiV1 := r.Group("/api/v1")
	{
		// 主题管理路由
		themeRoutes := apiV1.Group("/themes")
		{
			themeRoutes.GET("", themeHandler.GetThemeList)
			themeRoutes.POST("", themeHandler.CreateTheme)
			themeRoutes.GET("/:id", themeHandler.GetThemeDetail)
			themeRoutes.PUT("/:id", themeHandler.UpdateTheme)
			themeRoutes.DELETE("/:id", themeHandler.DeleteTheme)
			themeRoutes.POST("/:id/apply", themeHandler.ApplyTheme)
			themeRoutes.GET("/:id/preview", themeHandler.PreviewTheme)
			themeRoutes.GET("/:id/export", themeHandler.ExportTheme)
			themeRoutes.POST("/import", themeHandler.ImportTheme)
			themeRoutes.GET("/current", themeHandler.GetCurrentTheme)
		}
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 