package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	
	"wz-backend-go/internal/application/gateway/dto"
	appService "wz-backend-go/internal/application/gateway/service"
	"wz-backend-go/internal/domain/gateway/service"
	"wz-backend-go/internal/domain/gateway/valueobject"
	"wz-backend-go/internal/infrastructure/factory/gateway"
	"wz-backend-go/internal/interfaces/http/controller/gateway"
	"wz-backend-go/internal/interfaces/http/middleware/gateway"
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

	// 创建上下文
	ctx := context.Background()

	// 初始化基础设施
	dbConn, redisClient, err := setupInfrastructure(ctx, &cfg)
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v", err)
	}
	defer dbConn.Close()
	defer redisClient.Close()

	// 初始化仓库
	repoFactory := gateway.NewRepositoryFactory(dbConn, redisClient)
	routeRepo := repoFactory.GetRouteRepository()
	serviceRepo := repoFactory.GetServiceRepository()
	rateLimiterRepo := repoFactory.GetRateLimiterRepository()

	// 确保数据库表存在
	if err := setupDatabase(ctx, repoFactory); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化领域服务
	domainService := service.NewGatewayDomainService(routeRepo, serviceRepo)

	// 初始化应用服务
	applicationService := appService.NewGatewayApplicationService(domainService, routeRepo, serviceRepo)

	// 初始化HTTP控制器
	serviceController := gateway.NewServiceController(applicationService)
	routeController := gateway.NewRouteController(applicationService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册通用中间件
	r.Use(corsMiddleware())

	// API版本组
	apiV1 := r.Group("/api/v1")
	
	// 注册管理API
	registerManagementAPI(apiV1, serviceController, routeController)
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "gateway",
			"version": "1.0.0",
		})
	})

	// 为所有其他路由设置网关中间件链
	setupGatewayMiddleware(r, domainService, rateLimiterRepo)
	
	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 启动健康检查服务（每60秒检查一次所有服务健康状态）
	startHealthCheckService(ctx, applicationService, 60*time.Second)

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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("服务器关闭错误: %v", err)
	}

	fmt.Println("API网关服务已安全关闭")
}

// 初始化数据库和Redis连接
func setupInfrastructure(ctx context.Context, cfg *config.Config) (*sql.DB, *redis.Client, error) {
	// 初始化MySQL连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("连接MySQL失败: %w", err)
	}
	
	// 设置连接池参数
	dbConn.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	dbConn.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	dbConn.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	
	// 测试数据库连接
	if err := dbConn.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("MySQL连接测试失败: %w", err)
	}
	
	// 初始化Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	
	// 测试Redis连接
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, nil, fmt.Errorf("Redis连接测试失败: %w", err)
	}
	
	return dbConn, redisClient, nil
}

// 确保数据库表存在
func setupDatabase(ctx context.Context, repoFactory *gateway.RepositoryFactory) error {
	// 获取MySQL仓库实例
	routeRepo := repoFactory.GetMySQLRouteRepository()
	serviceRepo := repoFactory.GetMySQLServiceRepository()
	
	// 创建路由表
	if err := routeRepo.EnsureTable(ctx); err != nil {
		return fmt.Errorf("创建路由表失败: %w", err)
	}
	
	// 创建服务表
	if err := serviceRepo.EnsureTable(ctx); err != nil {
		return fmt.Errorf("创建服务表失败: %w", err)
	}
	
	return nil
}

// 注册管理API路由
func registerManagementAPI(apiGroup *gin.RouterGroup, serviceController *gateway.ServiceController, routeController *gateway.RouteController) {
	// 管理API组
	managementGroup := apiGroup.Group("/management")
	
	// 服务管理API
	serviceGroup := managementGroup.Group("/services")
	{
		serviceGroup.POST("", serviceController.RegisterService)
		serviceGroup.GET("", serviceController.ListServices)
		serviceGroup.GET("/:name", serviceController.GetService)
		serviceGroup.PUT("/:name", serviceController.UpdateService)
		serviceGroup.DELETE("/:name", serviceController.DeleteService)
		serviceGroup.GET("/:name/health", serviceController.CheckServiceHealth)
		serviceGroup.POST("/:name/routes", serviceController.AddRouteToService)
		serviceGroup.GET("/:name/routes", serviceController.GetServiceRoutes)
		serviceGroup.POST("/wz-categories", serviceController.CreateWanzhiCategoryRoutes)
	}
	
	// 路由管理API
	routeGroup := managementGroup.Group("/routes")
	{
		routeGroup.POST("", routeController.RegisterRoute)
		routeGroup.GET("", routeController.ListRoutes)
		routeGroup.GET("/:id", routeController.GetRoute)
		routeGroup.PUT("/:id", routeController.UpdateRoute)
		routeGroup.DELETE("/:id", routeController.DeleteRoute)
	}
}

// 设置网关中间件链
func setupGatewayMiddleware(r *gin.Engine, domainService *service.GatewayDomainService, rateLimiterRepo *redis.Client) {
	// 设置路由查找中间件（应用于所有非管理API路由）
	routeFinderMiddleware := gateway.NewRouteFinderMiddleware(domainService)
	
	// 设置认证中间件
	authMiddleware := gateway.NewAuthMiddleware()
	
	// 设置限流中间件
	rateLimitMiddleware := gateway.NewRateLimitMiddleware()
	
	// 设置代理中间件
	proxyMiddleware := gateway.NewProxyMiddleware()
	
	// 应用中间件链到除了API和健康检查外的所有路由
	r.NoRoute(
		routeFinderMiddleware.Handle(),   // 1. 查找匹配路由
		authMiddleware.Handle(),          // 2. 认证
		rateLimitMiddleware.Handle(),     // 3. 限流
		proxyMiddleware.Handle(),         // 4. 代理到目标服务
	)
}

// 启动健康检查服务
func startHealthCheckService(ctx context.Context, applicationService *appService.GatewayApplicationService, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				log.Println("健康检查服务停止")
				return
			case <-ticker.C:
				services, err := applicationService.ListServices(ctx, 0, 1000)
				if err != nil {
					log.Printf("获取服务列表失败: %v", err)
					continue
				}
				
				for _, service := range services.Items {
					// 对每个服务进行健康检查
					go func(serviceDTO dto.ServiceDTO) {
						serviceName, err := valueobject.NewServiceName(serviceDTO.Name)
						if err != nil {
							log.Printf("服务名称无效: %v", err)
							return
						}
						
						isHealthy, err := applicationService.CheckServiceHealth(ctx, serviceName)
						if err != nil {
							log.Printf("健康检查失败 [%s]: %v", serviceDTO.Name, err)
							return
						}
						
						if !isHealthy {
							log.Printf("服务不健康 [%s]", serviceDTO.Name)
						}
					}(service)
				}
			}
		}
	}()
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