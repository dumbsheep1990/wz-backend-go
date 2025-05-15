package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	
	"wz-backend-go/services/gateway-service/config"
)

// 配置文件路径
var configFile = flag.String("f", "configs/gateway.yaml", "配置文件路径")

func main() {
	// 解析命令行参数
	flag.Parse()

	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// 加载配置
	var cfg config.Config
	err := cfg.Load(*configFile)
	if err != nil {
		log.Printf("加载配置文件失败: %v，使用默认配置", err)
		cfg = config.DefaultConfig()
	}

	// 创建Gin引擎
	r := gin.Default()

	// 注册中间件
	r.Use(corsMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 设置反向代理路由
	setupReverseProxy(r, cfg.Services)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("API网关启动在端口 %d...\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭API网关服务...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭错误: %v", err)
	}

	fmt.Println("API网关服务已安全关闭")
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// 设置反向代理路由
func setupReverseProxy(r *gin.Engine, services []config.ServiceConfig) {
	apiGroup := r.Group("/api/v1")

	for _, service := range services {
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.Fatalf("解析服务URL失败: %v", err)
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 设置代理处理器
		proxyHandler := func(c *gin.Context) {
			// 日志记录
			log.Printf("代理请求到 %s: %s", service.Name, c.Request.URL.Path)

			// 设置原始主机头
			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				req.Host = targetURL.Host
			}

			// 处理响应
			proxy.ModifyResponse = func(resp *http.Response) error {
				// 日志记录响应状态
				log.Printf("来自 %s 的响应状态: %d", service.Name, resp.StatusCode)
				return nil
			}

			// 错误处理
			proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
				log.Printf("代理错误: %v", err)
				c.JSON(http.StatusBadGateway, gin.H{
					"error": "服务暂时不可用",
				})
			}

			// 执行代理
			proxy.ServeHTTP(c.Writer, c.Request)
		}

		// 配置服务特定的路由
		servicePath := fmt.Sprintf("/%s", service.Name)
		
		if service.RequireAuth {
			// 需要认证的API路由
			authGroup := apiGroup.Group(servicePath)
			authGroup.Use(authMiddleware())
			authGroup.Any("/*path", proxyHandler)
		} else {
			// 不需要认证的服务直接代理
			apiGroup.Group(servicePath).Any("/*path", proxyHandler)
		}
	}
}

// 认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权访问",
			})
			c.Abort()
			return
		}
		
		// TODO: 实现真实的TOKEN验证
		
		c.Next()
	}
} 