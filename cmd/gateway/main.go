package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	
	appService "wz-backend-go/internal/application/gateway/service"
	domainService "wz-backend-go/internal/domain/gateway/service"
	infraConfig "wz-backend-go/internal/infrastructure/gateway/config"
	infraGateway "wz-backend-go/internal/infrastructure/gateway"
	infraMiddleware "wz-backend-go/internal/infrastructure/gateway/middleware"
	gatewayRepo "wz-backend-go/internal/infrastructure/gateway/repository"
	controller "wz-backend-go/internal/interfaces/http/controller/gateway"
)

var (
	configPath      = flag.String("config", "etc/gateway.yaml", "配置文件路径")
	shutdownTimeout = flag.Duration("shutdown-timeout", 5*time.Second, "优雅关闭超时时间")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 从环境变量获取配置路径
	if envConfigPath := os.Getenv("GATEWAY_CONFIG_PATH"); envConfigPath != "" {
		*configPath = envConfigPath
	}

	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// 创建上下文
	ctx := context.Background()

	// 加载配置
	gatewayConfig, err := infraConfig.NewConfigLoader(*configPath).Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化基础设施
	dbConn, redisClient, err := setupInfrastructure(ctx, gatewayConfig)
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v", err)
	}
	defer dbConn.Close()
	defer redisClient.Close()

	// 初始化仓库
	routeRepo := gatewayRepo.NewRouteRepository(dbConn)
	serviceRepo := gatewayRepo.NewServiceRepository(dbConn)
	rateLimiterRepo := gatewayRepo.NewRateLimiterRepository(dbConn)

	// 初始化基础设施工厂
	infraFactory := infraGateway.NewFactory(
		dbConn,
		redisClient,
		gatewayConfig.Security.JWTSecret,
		"wanzhi.gateway", // 项目固定的JWT签发者
		gatewayConfig.Telemetry.ServiceVersion,
	)

	// 初始化基础设施组件
	router := infraFactory.CreateRouter(routeRepo)
	serviceRegistry := infraFactory.CreateServiceRegistry(serviceRepo)
	rateLimiter := infraFactory.CreateRateLimiter()
	authHandler := infraFactory.CreateAuthHandler(gatewayConfig.Security.JWTExpiration * 60)

	// 初始化领域服务
	gatewayDomainService := domainService.NewGatewayDomainService(routeRepo, serviceRepo)

	// 初始化应用服务
	gatewayApplicationService := appService.NewGatewayApplicationService(gatewayDomainService, routeRepo, serviceRepo)

	// 初始化HTTP控制器
	serviceController := controller.NewServiceController(gatewayApplicationService)
	routeController := controller.NewRouteController(gatewayApplicationService)

	// 初始化中间件管理器
	middlewareManager := infraFactory.CreateMiddlewareManager(gatewayDomainService)

	// 创建Gin引擎
	r := gin.Default()

	// 应用全局中间件
	infraFactory.ApplyGlobalMiddlewares(r, middlewareManager)

	// API版本组
	apiV1 := r.Group("/api/v1")
	
	// 注册管理API
	registerManagementAPI(apiV1, serviceController, routeController)
	
	// API代理路由组
	proxyGroup := r.Group("/")
	
	// 应用代理中间件链
	infraFactory.ApplyRouteMiddlewares(proxyGroup, middlewareManager)
	
	// 存储服务注册表到Gin上下文，供健康检查使用
	r.Use(func(c *gin.Context) {
		c.Set("service_registry", serviceRegistry)
		c.Next()
	})
	
	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", gatewayConfig.Server.Host, gatewayConfig.Server.Port),
		Handler: r,
		ReadTimeout:  time.Duration(gatewayConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(gatewayConfig.Server.WriteTimeout) * time.Second,
	}

	// 启动健康检查服务
	startHealthCheckService(ctx, serviceRegistry, 60*time.Second)

	// 启动服务器
	go func() {
		log.Printf("API网关启动在 %s:%d...\n", gatewayConfig.Server.Host, gatewayConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭API网关服务...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("服务器关闭错误: %v", err)
	}

	fmt.Println("API网关服务已安全关闭")
}

// 初始化数据库和Redis连接
func setupInfrastructure(ctx context.Context, config *infraConfig.GatewayConfig) (*sqlx.DB, *redis.Client, error) {
	// 从环境变量获取数据库配置
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	
	dbPort := 3306
	if dbPortEnv := os.Getenv("DB_PORT"); dbPortEnv != "" {
		if port, err := strconv.Atoi(dbPortEnv); err == nil {
			dbPort = port
		}
	}
	
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	
	dbPass := os.Getenv("DB_PASS")
	
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "wanzhidb"
	}
	
	// 初始化MySQL连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	
	dbConn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("连接MySQL失败: %w", err)
	}
	
	// 设置连接池参数
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(3600 * time.Second)
	
	// 测试数据库连接
	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("测试数据库连接失败: %w", err)
	}

	// 从环境变量获取Redis配置
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	
	redisPort := 6379
	if redisPortEnv := os.Getenv("REDIS_PORT"); redisPortEnv != "" {
		if port, err := strconv.Atoi(redisPortEnv); err == nil {
			redisPort = port
		}
	}
	
	redisPass := os.Getenv("REDIS_PASS")
	
	redisDB := 0
	if redisDBEnv := os.Getenv("REDIS_DB"); redisDBEnv != "" {
		if db, err := strconv.Atoi(redisDBEnv); err == nil {
			redisDB = db
		}
	}
	
	// 初始化Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: redisPass,
		DB:       redisDB,
	})

	// 测试Redis连接
	pingCmd := redisClient.Ping(ctx)
	if err := pingCmd.Err(); err != nil {
		return nil, nil, fmt.Errorf("测试Redis连接失败: %w", err)
	}

	log.Println("基础设施初始化成功")
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
func registerManagementAPI(apiGroup *gin.RouterGroup, serviceController *controller.ServiceController, routeController *controller.RouteController) {
	// 服务管理
	services := apiGroup.Group("/services")
	services.GET("", serviceController.ListServices)
	services.GET("/:id", serviceController.GetService)
	services.POST("", serviceController.CreateService)
	services.PUT("/:id", serviceController.UpdateService)
	services.DELETE("/:id", serviceController.DeleteService)
	services.PUT("/:id/active", serviceController.ToggleServiceActive)
	
	// 路由管理
	routes := apiGroup.Group("/routes")
	routes.GET("", routeController.ListRoutes)
	routes.GET("/:id", routeController.GetRoute)
	routes.POST("", routeController.CreateRoute)
	routes.PUT("/:id", routeController.UpdateRoute)
	routes.DELETE("/:id", routeController.DeleteRoute)
	routes.PUT("/:id/active", routeController.ToggleRouteActive)
	routes.GET("/service/:serviceId", routeController.ListRoutesByService)
}

// 设置网关中间件链
func setupGatewayMiddleware(r *gin.Engine, domainService *domainService.GatewayDomainService, rateLimiterRepo *redis.Client) {
	// 创建中间件
	routingMiddleware := middleware.NewRoutingMiddleware(domainService)
	authMiddleware := middleware.NewAuthMiddleware(domainService)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimiterRepo)
	
	// 通配路由处理网关逻辑，需要放在所有特定路由后面
	r.NoRoute(
		rateLimitMiddleware.Handle(),   // 限流
		routingMiddleware.Handle(),     // 路由解析
		authMiddleware.Handle(),        // 身份验证
		middleware.ProxyMiddleware(),      // 代理请求
	)
}

// 启动健康检查服务
func startHealthCheckService(ctx context.Context, serviceRegistry interface{}, interval time.Duration) {
	// 健康检查服务，定期检查所有后端服务的健康状态
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// 使用服务注册表检查所有服务健康状态
				if registry, ok := serviceRegistry.(infraGateway.ServiceRegistry); ok {
					registry.CheckAllServicesHealth(ctx)
				}
			case <-ctx.Done():
				return
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
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}
