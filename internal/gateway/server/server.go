package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/gateway/auth"
	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/gateway/docs"
	"wz-backend-go/internal/gateway/handler"
	"wz-backend-go/internal/gateway/middleware"
	"wz-backend-go/internal/gateway/router"
	grpcClient "wz-backend-go/internal/gateway/grpc"
)

// Server API网关服务器
type Server struct {
	engine         *gin.Engine
	config         config.Config
	httpServer     *http.Server
	authManager    *auth.AuthManager
	dynamicRouter  *router.DynamicRouter
	grpcHandler    *handler.GrpcHandler
	swaggerHandler *handler.SwaggerHandler
}

// NewServer 创建新的API网关服务器
func NewServer(configPath string) (*Server, error) {
	// 加载配置
	var conf config.Config
	if err := conf.Load(configPath); err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建Gin引擎
	engine := gin.New()

	// 创建认证管理器
	authManager := auth.NewAuthManager(conf.Security)

	// 创建动态路由管理器
	dynamicRouter := router.NewDynamicRouter(engine)

	// 创建gRPC处理器
	grpcHandler := handler.NewGrpcHandler(conf.GRPC)

	// 创建Swagger处理器
	docsPath := filepath.Join("public", "docs")
	swaggerHandler := handler.NewSwaggerHandler(conf, docsPath)

	return &Server{
		engine:         engine,
		config:         conf,
		authManager:    authManager,
		dynamicRouter:  dynamicRouter,
		grpcHandler:    grpcHandler,
		swaggerHandler: swaggerHandler,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	// 设置全局中间件
	s.setupMiddleware()

	// 注册路由
	s.registerRoutes()

	// 配置HTTP服务器
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler:        s.engine,
		ReadTimeout:    time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动HTTP服务器
	fmt.Printf("API网关服务器启动: %s\n", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP服务器启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("关闭服务器...")

	// 给当前请求留出时间完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭gRPC连接
	s.grpcHandler.Stop()

	// 关闭HTTP服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭错误: %w", err)
	}

	return nil
}

// setupMiddleware 设置全局中间件
func (s *Server) setupMiddleware() {
	// 恢复中间件
	s.engine.Use(gin.Recovery())

	// 日志中间件
	s.engine.Use(middleware.Logger())

	// CORS中间件
	if s.config.Cors.Enabled {
		s.engine.Use(middleware.CORS(s.config.Cors))
	}

	// 安全中间件
	if s.config.Security.XSSProtection {
		s.engine.Use(middleware.Security())
	}

	// 超时中间件
	s.engine.Use(middleware.Timeout(30))

	// 请求ID中间件
	s.engine.Use(middleware.RequestID())
}

// registerRoutes 注册路由
func (s *Server) registerRoutes() {
	// 注册Swagger文档路由
	s.swaggerHandler.RegisterSwaggerRoutes(s.engine)

	// 注册路由管理API
	routeAPI := router.NewRouteAPI(s.dynamicRouter)
	routeAPI.RegisterRouteAPI(s.engine)

	// 注册健康检查
	s.engine.GET("/health", handler.HealthCheck)

	// 注册所有服务路由
	for _, service := range s.config.Services {
		// 创建路由组
		group := s.engine.Group(service.Prefix)

		// 添加认证中间件（如果需要）
		if service.Authentication {
			group.Use(middleware.EnhancedAuthentication(s.authManager, s.config.Security))
		}

		// 根据服务类型注册路由
		switch service.Type {
		case "http":
			// 注册HTTP路由
			s.dynamicRouter.RegisterService(service, handler.RegisterHTTPRoutes)
		case "grpc":
			// 注册gRPC路由
			s.dynamicRouter.RegisterService(service, s.grpcHandler.RegisterGRPCServiceRoutes)
		default:
			fmt.Printf("未知的服务类型: %s\n", service.Type)
		}
	}
}
