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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	_ "github.com/go-sql-driver/mysql"
	
	"wz-backend-go/internal/delivery/http/handler/content"
	"wz-backend-go/internal/delivery/http/router"
	"wz-backend-go/internal/repository"
)

func main() {
	// 设置数据库连接
	dbConn, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// 创建仓储
	contentRepo := repository.NewSqlContentRepository(dbConn)

	// 创建HTTP处理器
	categoryHandler := content.NewCategoryHandler(contentRepo)
	postHandler := content.NewPostHandler(contentRepo)

	// 设置Gin路由
	ginRouter := gin.Default()
	ginRouter.Use(gin.Recovery())

	// 注册路由
	router.RegisterContentRoutes(ginRouter, categoryHandler, postHandler)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":8002",
		Handler: ginRouter,
	}

	// 在goroutine中运行服务器
	go func() {
		log.Printf("Starting content microservice on port 8002")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down content microservice...")

	// 创建优雅关闭的截止时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Content microservice exited")
}

// setupDatabase 建立MySQL数据库连接
func setupDatabase() (sqlx.SqlConn, error) {
	// 从环境变量获取数据库连接详情
	host := getEnvWithDefault("DB_HOST", "localhost")
	port := getEnvWithDefault("DB_PORT", "3306")
	user := getEnvWithDefault("DB_USER", "root")
	password := getEnvWithDefault("DB_PASSWORD", "password")
	dbname := getEnvWithDefault("DB_NAME", "wanzhi")

	// 创建DSN字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	// 打开数据库连接
	conn := sqlx.NewMysql(dsn)

	log.Println("Successfully connected to the database")
	return conn, nil
}

// getEnvWithDefault 返回环境变量的值，如果未设置则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 