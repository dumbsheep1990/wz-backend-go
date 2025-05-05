package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wz-backend-go/internal/gateway"
	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/telemetry"
)

var configFile = flag.String("f", "configs/gateway.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	var c config.Config
	err := c.Load(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	
	// 初始化OpenTelemetry
	tp, err := telemetry.InitTracer(
		"api-gateway",               // 服务名称 
		"1.0.0",                     // 服务版本
		c.Server.Environment,        // 环境
		"otlp",                      // 导出器类型
		c.Telemetry.CollectorURL,    // 收集器地址
	)
	if err != nil {
		log.Printf("警告: OpenTelemetry初始化失败: %v", err)
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), c.Server.ShutdownTimeout)
			defer cancel()
			telemetry.Shutdown(ctx, tp)
		}()
	}

	// 启动网关服务
	server, err := gateway.NewServer(c)
	if err != nil {
		log.Fatalf("创建网关服务失败: %v", err)
	}
	
	// 启动服务
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("启动网关服务失败: %v", err)
		}
	}()
	
	fmt.Printf("API网关服务已启动，监听端口: %d\n", c.Server.Port)
	
	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("正在关闭API网关服务...")
	server.Shutdown()
	fmt.Println("API网关服务已安全关闭")
}
