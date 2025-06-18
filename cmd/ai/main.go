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
	"github.com/zeromicro/go-zero/zrpc"
	
	"wz-backend-go/internal/delivery/http/handler/ai"
	"wz-backend-go/internal/delivery/http/router"
	"wz-backend-go/internal/delivery/rpc/aiclient"
)

func main() {
	// 设置AI RPC客户端连接
	aiRpcClient, err := setupAIRpcClient()
	if err != nil {
		log.Fatalf("Failed to setup AI RPC client: %v", err)
	}

	// 创建AI客户端
	aiClient := aiclient.NewAI(aiRpcClient)

	// 创建HTTP处理器
	aiHandler := ai.NewAIHandler(aiClient)

	// 设置Gin路由
	ginRouter := gin.Default()
	ginRouter.Use(gin.Recovery())

	// 注册路由
	router.RegisterAIRoutes(ginRouter, aiHandler)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":8003",
		Handler: ginRouter,
	}

	// 在goroutine中运行服务器
	go func() {
		log.Printf("Starting AI microservice on port 8003")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down AI microservice...")

	// 创建优雅关闭的截止时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("AI microservice exited")
}

// setupAIRpcClient 设置AI RPC客户端连接
func setupAIRpcClient() (zrpc.Client, error) {
	// 从环境变量获取AI RPC服务地址
	aiRpcTarget := getEnvWithDefault("AI_RPC_TARGET", "localhost:50051")
	
	// 创建RPC客户端配置
	clientConfig := zrpc.RpcClientConf{
		Endpoints: []string{aiRpcTarget},
		Timeout:   5000, // 5秒超时
	}

	// 创建RPC客户端
	client, err := zrpc.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI RPC client: %w", err)
	}

	log.Printf("Successfully connected to AI RPC service at %s", aiRpcTarget)
	return client, nil
}

// getEnvWithDefault 返回环境变量的值，如果未设置则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 