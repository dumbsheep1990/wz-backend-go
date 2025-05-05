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
	"wz-backend-go/internal/delivery/rpc/user"
	"wz-backend-go/internal/telemetry"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	
	// 初始化OpenTelemetry
	tp, err := telemetry.InitTracer(
		"user-service",              // 服务名称 
		"1.0.0",                     // 服务版本
		c.Mode,                      // 环境
		"otlp",                      // 导出器类型
		c.Telemetry.CollectorURL,    // 收集器地址
	)
	if err != nil {
		log.Printf("警告: OpenTelemetry初始化失败: %v", err)
	}
	
	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 添加OpenTelemetry gRPC拦截器
	var opts []grpc.ServerOption
	if tp != nil {
		opts = telemetry.GRPCServerMiddleware()
	}
	
	// 创建服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	}, opts...)
	
	// 设置服务关闭时的清理操作
	s.AddShutdownHooks(func() {
		if tp != nil {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			telemetry.Shutdown(shutdownCtx, tp)
		}
	})
	
	// 创建程序终止信号处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// 启动服务
	fmt.Printf("Starting user rpc server at %s...\n", c.ListenOn)
	go s.Start()
	
	// 等待终止信号
	<-quit
	fmt.Println("正在关闭用户服务...")
	s.Stop()
	fmt.Println("用户服务已安全关闭")
}
