package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"wz-backend-go/internal/gateway/config"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// grpcClients 用于缓存gRPC连接
var grpcClients = make(map[string]*grpc.ClientConn)

// RegisterGRPCRoutes 注册gRPC服务路由
func RegisterGRPCRoutes(group *gin.RouterGroup, service config.ServiceConfig) {
	// 确保gRPC客户端连接已建立
	conn, err := getGRPCConnection(service)
	if err != nil {
		// 记录错误，但不中断启动过程
		fmt.Printf("无法建立gRPC连接到服务 %s: %v\n", service.Name, err)
	}

	// 处理配置中的gRPC方法
	if len(service.GrpcOptions.Methods) > 0 {
		for _, method := range service.GrpcOptions.Methods {
			// 注册HTTP路由到gRPC方法
			httpMethod := strings.ToUpper(method.HTTPMethod)
			path := method.Path

			switch httpMethod {
			case "GET":
				group.GET(path, createGRPCHandler(service, method.Name, conn))
			case "POST":
				group.POST(path, createGRPCHandler(service, method.Name, conn))
			case "PUT":
				group.PUT(path, createGRPCHandler(service, method.Name, conn))
			case "DELETE":
				group.DELETE(path, createGRPCHandler(service, method.Name, conn))
			case "PATCH":
				group.PATCH(path, createGRPCHandler(service, method.Name, conn))
			default:
				// 默认使用POST方法
				group.POST(path, createGRPCHandler(service, method.Name, conn))
			}
		}
	} else {
		// 如果没有定义具体方法，提供一个通用的gRPC调用端点
		group.POST("/:method", func(c *gin.Context) {
			methodName := c.Param("method")
			handleGRPCRequest(c, service, methodName, conn)
		})
	}
}

// getGRPCConnection 获取或创建gRPC连接
func getGRPCConnection(service config.ServiceConfig) (*grpc.ClientConn, error) {
	// 检查缓存中是否已有连接
	if conn, ok := grpcClients[service.Name]; ok {
		return conn, nil
	}

	// 设置gRPC选项
	var opts []grpc.DialOption

	// 设置消息大小限制
	if service.GrpcOptions.MaxRecvMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(service.GrpcOptions.MaxRecvMsgSize),
		))
	}
	if service.GrpcOptions.MaxSendMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(service.GrpcOptions.MaxSendMsgSize),
		))
	}

	// 设置TLS
	if service.GrpcOptions.TLS {
		// TODO: 实现TLS证书加载
		// creds := credentials.NewTLS(&tls.Config{...})
		// opts = append(opts, grpc.WithTransportCredentials(creds))
		return nil, fmt.Errorf("TLS功能尚未实现")
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 创建gRPC连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, service.Target, opts...)
	if err != nil {
		return nil, err
	}

	// 缓存连接
	grpcClients[service.Name] = conn

	return conn, nil
}

// createGRPCHandler 创建gRPC处理器
func createGRPCHandler(service config.ServiceConfig, methodName string, conn *grpc.ClientConn) gin.HandlerFunc {
	return func(c *gin.Context) {
		handleGRPCRequest(c, service, methodName, conn)
	}
}

// handleGRPCRequest 处理gRPC请求
func handleGRPCRequest(c *gin.Context, service config.ServiceConfig, methodName string, conn *grpc.ClientConn) {
	// 如果连接不可用，返回错误
	if conn == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "gRPC服务不可用",
		})
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无法读取请求体",
			"error":   err.Error(),
		})
		return
	}

	// 构建gRPC请求上下文
	fullMethod := fmt.Sprintf("/%s.%s/%s",
		service.GrpcOptions.PackageName,
		service.GrpcOptions.ServiceName,
		methodName)

	// 创建带超时的上下文
	timeout := 30 * time.Second
	if service.Timeout > 0 {
		timeout = time.Duration(service.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	defer cancel()

	// 添加请求元数据
	md := metadata.New(nil)
	// 从HTTP头传递关键信息到gRPC元数据
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			md.Set(k, v[0])
		}
	}
	// 特别处理租户ID
	if tenantID, exists := c.Get("tenantID"); exists {
		md.Set("X-Tenant-ID", tenantID.(string))
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 调用通用gRPC方法
	var respHeaders metadata.MD
	var respTrailers metadata.MD
	
	// 创建用于接收响应的缓冲区
	var resp []byte
	
	// 调用gRPC方法
	err = conn.Invoke(ctx, fullMethod, body, &resp, grpc.Header(&respHeaders), grpc.Trailer(&respTrailers))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "gRPC调用失败",
			"error":   err.Error(),
		})
		return
	}

	// 将gRPC响应头传递到HTTP头
	for k, v := range respHeaders {
		for _, value := range v {
			c.Header(k, value)
		}
	}

	// 返回gRPC响应
	c.Data(http.StatusOK, "application/json", resp)
}

// 关闭所有gRPC连接
func CloseGRPCConnections() {
	for name, conn := range grpcClients {
		if err := conn.Close(); err != nil {
			fmt.Printf("关闭gRPC连接 %s 失败: %v\n", name, err)
		}
	}
}
