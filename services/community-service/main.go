package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/yourusername/wz-backend-go/api/community"
	"github.com/yourusername/wz-backend-go/services/community-service/handlers"
	"github.com/yourusername/wz-backend-go/services/community-service/middleware"
	"github.com/yourusername/wz-backend-go/services/community-service/service"
)

const (
	grpcPort = ":50054"  // Choose an available port for gRPC
	httpPort = ":8084"   // Choose an available port for HTTP
)

// setupGRPCServer 设置gRPC服务器
func setupGRPCServer(service *service.CommunityServiceServer) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterCommunityServiceServer(grpcServer, service)
	
	log.Printf("gRPC server started on port %s", grpcPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}

// setupRESTServer 使用Gin设置REST服务器
func setupRESTServer(service *service.CommunityServiceServer) {
	router := gin.Default()
	
	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID", "X-User-Name"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// 设置静态文件服务
	router.Static("/static", "./static")
	
	// 注册API路由
	handlers.RegisterRoutes(router, service)
	
	// 注册认证路由
	handlers.RegisterAuthRoutes(router)
	
	// 注册健康检查路由
	handlers.RegisterHealthRoutes(router)
	
	// 前端页面路由，单页应用支持
	router.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})
	
	log.Printf("REST服务器已在端口 %s 启动", httpPort)
	go func() {
		if err := router.Run(httpPort); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()
}

func main() {
	service := service.NewCommunityServiceServer()
	
	// 启动gRPC服务器
	setupGRPCServer(service)
	
	// 启动REST服务器
	setupRESTServer(service)
	
	// 保持主goroutine不退出
	select {}
}
