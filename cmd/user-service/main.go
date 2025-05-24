package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "wz-backend-go/services/user-service/api/proto"
	appuser "wz-backend-go/internal/application/user"
	domainuser "wz-backend-go/internal/domain/user"
	"wz-backend-go/internal/infrastructure/eventbus"
	"wz-backend-go/internal/infrastructure/repository"
	grpcuser "wz-backend-go/internal/interfaces/grpc/user"
)

func main() {
	// 配置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("启动用户服务...")

	// 获取数据库连接信息（实际应用中应从配置中读取）
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "password")
	dbName := getEnv("DB_NAME", "wz_backend")
	
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	
	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}
	defer db.Close()
	
	// 设置数据库连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 创建事件总线
	eventBus := eventbus.NewEventBus()
	
	// 注册事件处理器（实际应用中可能需要更多处理器）
	eventBus.SubscribeFunc("UserRegisteredEvent", func(event interface{}) error {
		log.Printf("处理用户注册事件: %v", event)
		return nil
	})
	
	eventBus.SubscribeFunc("UserLoggedInEvent", func(event interface{}) error {
		log.Printf("处理用户登录事件: %v", event)
		return nil
	})

	// 创建仓储
	userRepo := repository.NewUserRepository(db)
	
	// 创建领域服务
	userDomainService := domainuser.NewUserService()
	
	// 创建应用服务
	userAppService := appuser.NewUserApplicationService(userRepo, userDomainService, eventBus)
	
	// 创建gRPC处理器
	userHandler := grpcuser.NewUserHandler(userAppService)
	
	// 获取gRPC服务器端口
	grpcPort := getEnv("GRPC_PORT", "50051")
	
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()
	
	// 注册服务
	pb.RegisterUserServer(grpcServer, userHandler)
	
	// 启用反射服务，便于gRPC客户端发现服务
	reflection.Register(grpcServer)
	
	// 启动gRPC服务器
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}
	
	go func() {
		log.Printf("gRPC服务器正在监听端口: %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC服务器运行失败: %v", err)
		}
	}()
	
	// 等待终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("正在关闭服务器...")
	
	// 优雅关闭gRPC服务器
	grpcServer.GracefulStop()
	
	log.Println("服务器已关闭")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
