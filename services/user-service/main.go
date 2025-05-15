package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"wz-backend-go/services/user-service/config"
	"wz-backend-go/services/user-service/internal/server"
	"wz-backend-go/services/user-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "configs/user.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	var cfg config.Config
	err := cfg.Load(*configFile)
	if err != nil {
		log.Printf("加载配置文件失败: %v，使用默认配置", err)
		cfg = config.DefaultConfig()
	}

	// 创建服务对象
	userService := service.NewUserService()

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	server.RegisterUserServer(grpcServer, userService)

	// 如果是开发模式，启用反射服务
	if cfg.Server.Environment == "development" || cfg.Server.Environment == "testing" {
		reflection.Register(grpcServer)
	}

	// 启动服务
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("启动用户服务，监听地址: %s\n", addr)

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("监听失败: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭用户服务...")

	// 优雅停止gRPC服务器
	grpcServer.GracefulStop()

	log.Println("用户服务已安全关闭")
}
