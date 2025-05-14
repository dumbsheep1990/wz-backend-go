package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"wz-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// 服务配置
type ServiceConfig struct {
	Name        string
	URL         string
	RequireAuth bool
}

// 服务路由配置
var serviceRoutes = []ServiceConfig{
	{
		Name:        "site-service",
		URL:         "http://localhost:8081",
		RequireAuth: true,
	},
	{
		Name:        "page-service",
		URL:         "http://localhost:8082",
		RequireAuth: true,
	},
	{
		Name:        "component-service",
		URL:         "http://localhost:8083",
		RequireAuth: true,
	},
	{
		Name:        "render-service",
		URL:         "http://localhost:8084",
		RequireAuth: false,
	},
}

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

	// 设置反向代理路由
	setupReverseProxy(r)

	// 获取服务端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务
	log.Printf("API网关启动在端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

// 设置反向代理路由
func setupReverseProxy(r *gin.Engine) {
	apiGroup := r.Group("/api/v1")

	for _, service := range serviceRoutes {
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

		// 根据服务配置确定是否需要认证中间件
		if service.RequireAuth {
			// 公开API仍可访问
			if service.Name == "site-service" {
				apiGroup.GET("/site-templates", proxyHandler)
				apiGroup.GET("/site-templates/:id", proxyHandler)
			}
			if service.Name == "component-service" {
				apiGroup.GET("/components/categories", proxyHandler)
				apiGroup.GET("/components/:type", proxyHandler)
			}
			if service.Name == "render-service" {
				apiGroup.GET("/render/site", proxyHandler)
				apiGroup.GET("/render/sites/:siteId/:slug", proxyHandler)
			}

			// 需要认证的API路由
			authGroup := apiGroup.Group("")
			authGroup.Use(middleware.Auth())
			{
				authGroup.Any("/*path", proxyHandler)
			}
		} else {
			// 不需要认证的服务直接代理
			apiGroup.Any("/*path", proxyHandler)
		}
	}
}

// 代理请求体
func proxyRequestBody(src io.ReadCloser, dst io.Writer) error {
	defer src.Close()
	_, err := io.Copy(dst, src)
	return err
}
