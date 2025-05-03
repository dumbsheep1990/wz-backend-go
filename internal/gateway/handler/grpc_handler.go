package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	
	"wz-backend-go/internal/gateway/config"
	grpcClient "wz-backend-go/internal/gateway/grpc"
)

// GrpcHandler 管理gRPC路由处理
type GrpcHandler struct {
	clientManager *grpcClient.ClientManager
}

// NewGrpcHandler 创建新的gRPC处理器
func NewGrpcHandler(conf config.GRPCConfig) *GrpcHandler {
	return &GrpcHandler{
		clientManager: grpcClient.NewClientManager(conf),
	}
}

// RegisterGRPCServiceRoutes 注册单个gRPC服务的路由
func (h *GrpcHandler) RegisterGRPCServiceRoutes(group *gin.RouterGroup, service config.ServiceConfig) {
	// 处理配置中的gRPC方法
	if len(service.GrpcOptions.Methods) > 0 {
		for _, method := range service.GrpcOptions.Methods {
			// 注册HTTP路由到gRPC方法
			httpMethod := method.HTTPMethod
			path := method.Path
			methodName := method.Name

			// 根据HTTP方法注册路由
			switch httpMethod {
			case "GET":
				group.GET(path, h.createGRPCHandler(service, methodName))
			case "POST":
				group.POST(path, h.createGRPCHandler(service, methodName))
			case "PUT":
				group.PUT(path, h.createGRPCHandler(service, methodName))
			case "DELETE":
				group.DELETE(path, h.createGRPCHandler(service, methodName))
			case "PATCH":
				group.PATCH(path, h.createGRPCHandler(service, methodName))
			default:
				// 默认使用POST方法
				group.POST(path, h.createGRPCHandler(service, methodName))
			}
		}
	} else if service.GrpcOptions.EnableReflection {
		// 如果启用了反射但没有定义具体方法，提供通用端点
		group.POST("/:method", func(c *gin.Context) {
			methodName := c.Param("method")
			h.handleGRPCRequest(c, service, methodName)
		})
	}
}

// createGRPCHandler 创建gRPC处理器
func (h *GrpcHandler) createGRPCHandler(service config.ServiceConfig, methodName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.handleGRPCRequest(c, service, methodName)
	}
}

// handleGRPCRequest 处理gRPC请求
func (h *GrpcHandler) handleGRPCRequest(c *gin.Context, service config.ServiceConfig, methodName string) {
	// 获取或创建gRPC客户端
	client, err := h.clientManager.GetClient(context.Background(), service.Name, service.Target)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "无法连接gRPC服务",
			"error":   err.Error(),
		})
		return
	}

	// 处理请求
	grpcClient.HTTPToGRPC(c, client, methodName, service.Timeout)
}

// Stop 停止gRPC处理器
func (h *GrpcHandler) Stop() {
	h.clientManager.Stop()
}
