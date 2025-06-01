package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"wz-backend-go/internal/infrastructure/event"
	"wz-backend-go/internal/infrastructure/factory"
	"wz-backend-go/internal/infrastructure/service"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 无法加载 .env 文件，将使用环境变量")
	}

	// 初始化数据库连接
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 初始化Redis连接
	redisClient, err := initRedis()
	if err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}
	defer redisClient.Close()

	// 创建事件总线
	eventBus := event.NewSimpleEventBus()

	// 初始化Casbin服务
	casbinService, err := initCasbinService()
	if err != nil {
		log.Fatalf("Casbin初始化失败: %v", err)
	}

	// 初始化JWT配置
	jwtConfig := initJWTConfig()

	// 创建管理员服务工厂
	adminFactory := factory.NewAdminServiceFactory(
		db,
		redisClient,
		eventBus,
		casbinService.GetEnforcer(),
		jwtConfig,
	)

	// 设置事件处理器
	adminFactory.SetupEventHandlers()

	// 创建Gin路由
	router := gin.Default()

	// 注册CORS和其他中间件
	setupMiddlewares(router)

	// 注册路由
	adminFactory.RegisterRoutes(router)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    getServerAddress(),
		Handler: router,
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		log.Printf("管理员服务启动于 %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待关闭信号
	<-quit
	log.Println("关闭服务器...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已关闭")
}

// 初始化数据库连接
func initDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/wz_admin?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("数据库连接成功")
	return db, nil
}

// 初始化Redis连接
func initRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := 0 // 使用默认数据库

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis连接成功")
	return client, nil
}

// 初始化Casbin服务
func initCasbinService() (*service.CasbinService, error) {
	modelPath := os.Getenv("CASBIN_MODEL_PATH")
	policyPath := os.Getenv("CASBIN_POLICY_PATH")

	return service.NewCasbinService(modelPath, policyPath)
}

// 初始化JWT配置
func initJWTConfig() service.JWTConfig {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "万知网站后台默认密钥，请在生产环境中修改" // 默认密钥，生产环境应修改
	}

	// 解析过期时间（分钟）
	expirationMin := 60 * 24 // 默认24小时
	issuer := "万知网站管理后台"

	return service.JWTConfig{
		SecretKey:     secretKey,
		ExpirationMin: expirationMin,
		Issuer:        issuer,
	}
}

// 设置中间件
func setupMiddlewares(router *gin.Engine) {
	// CORS中间件
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

	// 日志中间件
	router.Use(gin.Logger())

	// 错误恢复中间件
	router.Use(gin.Recovery())
}

// 获取服务器地址
func getServerAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
