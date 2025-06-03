package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wz-backend-go/internal/delivery/rpc/internal/config"
	"wz-backend-go/internal/delivery/rpc/internal/server"
	"wz-backend-go/internal/delivery/rpc/internal/svc"
	"wz-backend-go/internal/delivery/rpc/gateway"
	"wz-backend-go/internal/telemetry"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	// 初始化日志
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c)
	
	// 初始化OpenTelemetry
	cleanup, err := telemetry.SetupOTelSDK(context.Background(), "gateway-rpc-service", c.Telemetry.Endpoint)
	if err != nil {
		log.Fatalf("Failed to setup OpenTelemetry: %v", err)
	}
	defer cleanup()
	
	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	
	// 初始化gRPC服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 注册网关服务
		gateway.RegisterGatewayServer(grpcServer, server.NewGatewayServer(ctx))
		
		// 在开发和测试模式下注册反射服务
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	
	// 添加优雅关闭处理
	s.AddShutdownCallback(func() {
		log.Println("Shutting down gateway RPC service...")
	})
	
	// 启动gRPC服务器
	go func() {
		fmt.Printf("Starting gateway RPC server at %s...\n", c.ListenOn)
		s.Start()
	}()
	
	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 停止服务
	s.Stop()
	log.Println("Gateway RPC service stopped gracefully")
}
