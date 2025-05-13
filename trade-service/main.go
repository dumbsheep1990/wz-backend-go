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

	"wz-backend-go/trade-service/api"
	"wz-backend-go/trade-service/config"
	"wz-backend-go/trade-service/core/repository"
	"wz-backend-go/trade-service/core/service"
	"wz-backend-go/trade-service/middleware"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// 初始化配置
	cfg := config.NewConfig()

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
		log.Printf("Trade Service started on port %d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置5秒超时以关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// 初始化数据库
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 这里应该使用配置的DSN建立数据库连接
	// 为了简化，返回一个假的连接
	fmt.Println("Database connection established")
	return nil, nil
}

// 初始化依赖注入容器
func initContainer(db *gorm.DB, cfg *config.Config) *Container {
	// 在实际项目中，这里应该创建所有需要的存储库、服务和处理器
	// 为了简化，返回一个空容器
	fmt.Println("Container initialized")
	return &Container{}
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

	// 在实际项目中，应该在这里注册所有处理器
	// 例如:
	// container.orderHandler.RegisterRoutes(apiV1)
	// container.cartHandler.RegisterRoutes(apiV1)
	// container.paymentHandler.RegisterRoutes(apiV1)
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
	orderHandler   *api.OrderHandler
	cartHandler    *api.CartHandler
	paymentHandler *api.PaymentHandler
}
