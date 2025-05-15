package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"wz-backend-go/services/trade-service/config"
	"wz-backend-go/services/trade-service/internal/handler"
	"wz-backend-go/services/trade-service/internal/repository"
	"wz-backend-go/services/trade-service/internal/service"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// 初始化配置
	cfg, err := config.Load("configs/trade.yaml")
	if err != nil {
		log.Printf("加载配置文件失败: %v，使用默认配置", err)
		cfg = config.DefaultConfig()
	}

	// 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建依赖注入容器
	container := initContainer(db, cfg)

	// 创建并配置Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	// 注册API路由
	registerRoutes(r, container)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		log.Printf("交易服务已启动，监听端口 %d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭交易服务...")

	// 设置5秒超时以关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	log.Println("交易服务已安全关闭")
}

// 初始化数据库
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 这里应该使用配置的DSN建立数据库连接
	// 为了简化，返回一个假的连接
	fmt.Println("数据库连接已建立")
	return nil, nil
}

// 中间件
var middleware = struct {
	Cors   func() gin.HandlerFunc
	Logger func() gin.HandlerFunc
}{
	Cors: func() gin.HandlerFunc {
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
	},
	Logger: func() gin.HandlerFunc {
		return func(c *gin.Context) {
			start := time.Now()
			c.Next()
			latency := time.Since(start)
			log.Printf("请求: %s %s | 状态: %d | 耗时: %s",
				c.Request.Method, c.Request.URL.Path,
				c.Writer.Status(), latency)
		}
	},
}

// 初始化依赖注入容器
func initContainer(db *gorm.DB, cfg *config.Config) *Container {
	// 创建仓库
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// 创建服务
	orderSvc := service.NewOrderService(orderRepo)
	cartSvc := service.NewCartService(cartRepo)
	paymentSvc := service.NewPaymentService(paymentRepo)

	// 创建处理器
	orderHandler := handler.NewOrderHandler(orderSvc)
	cartHandler := handler.NewCartHandler(cartSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	return &Container{
		// 仓库
		orderRepository:   orderRepo,
		cartRepository:    cartRepo,
		paymentRepository: paymentRepo,

		// 服务
		orderService:   orderSvc,
		cartService:    cartSvc,
		paymentService: paymentSvc,

		// 处理器
		orderHandler:   orderHandler,
		cartHandler:    cartHandler,
		paymentHandler: paymentHandler,
	}
}

// 注册API路由
func registerRoutes(r *gin.Engine, container *Container) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 添加API版本前缀
	apiV1 := r.Group("/api/v1")

	// 注册路由
	container.orderHandler.RegisterRoutes(apiV1)
	container.cartHandler.RegisterRoutes(apiV1)
	container.paymentHandler.RegisterRoutes(apiV1)
}

// Container 依赖注入容器
type Container struct {
	// 存储库
	orderRepository   repository.OrderRepository
	cartRepository    repository.CartRepository
	paymentRepository repository.PaymentRepository

	// 服务
	orderService   service.OrderService
	cartService    service.CartService
	paymentService service.PaymentService

	// 处理器
	orderHandler   *handler.OrderHandler
	cartHandler    *handler.CartHandler
	paymentHandler *handler.PaymentHandler
}
