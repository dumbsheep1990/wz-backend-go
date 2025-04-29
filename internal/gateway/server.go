package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/gateway/handler"
	"wz-backend-go/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// Server 表示API网关服务器
type Server struct {
	config     config.Config
	router     *gin.Engine
	httpServer *http.Server
}

// NewServer 创建一个新的网关服务器
func NewServer(config config.Config) (*Server, error) {
	// 设置gin模式
	if config.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建gin路由器
	router := gin.New()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// 添加跨域中间件
	if config.Cors.Enabled {
		router.Use(middleware.Cors(config.Cors))
	}

	// 添加安全相关中间件
	if config.Security.XSSProtection {
		router.Use(middleware.SecurityHeaders())
	}

	// 添加租户识别中间件
	router.Use(middleware.TenantIdentifier())

	// 添加限流中间件
	if config.Rate.Enabled {
		router.Use(middleware.RateLimiter(config.Rate))
	}

	// 注册路由
	err := setupRoutes(router, config)
	if err != nil {
		return nil, fmt.Errorf("设置路由失败: %w", err)
	}

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		config:     config,
		router:     router,
		httpServer: httpServer,
	}, nil
}

// Start 启动API网关服务器
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown 优雅关闭API网关服务器
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("服务器关闭出错: %v\n", err)
	}

	// 关闭所有gRPC连接
	handler.CloseGRPCConnections()
}

// setupRoutes 设置API路由
func setupRoutes(router *gin.Engine, config config.Config) error {
	// 健康检查路由
	router.GET("/health", handler.HealthCheck())

	// 添加API文档路由
	router.GET("/swagger/*any", handler.Swagger())

	// 添加状态监控路由
	router.GET("/metrics", handler.Metrics())

	// 添加服务路由
	for _, service := range config.Services {
		// 为每个服务创建路由组
		group := router.Group(service.Prefix)

		// 根据服务配置添加中间件
		if service.Authentication {
			group.Use(middleware.Authentication(config.Security))
		}

		// 添加请求追踪中间件
		group.Use(middleware.Tracing(service.Name))

		// 添加超时中间件
		if service.Timeout > 0 {
			group.Use(middleware.Timeout(time.Duration(service.Timeout) * time.Second))
		}

		// 根据服务类型设置不同的代理处理方式
		switch service.Type {
		case "http":
			handler.RegisterHTTPRoutes(group, service)
		case "grpc":
			handler.RegisterGRPCRoutes(group, service)
		default:
			return fmt.Errorf("不支持的服务类型: %s", service.Type)
		}
	}

	return nil
}
