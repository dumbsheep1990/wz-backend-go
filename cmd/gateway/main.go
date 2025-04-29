package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wz-backend-go/internal/gateway"
	"wz-backend-go/internal/gateway/config"
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
