package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"wz-backend-go/internal/gateway/config"
)

// ClientManager 管理gRPC客户端连接
type ClientManager struct {
	clients     map[string]*Client
	config      config.GRPCConfig
	mu          sync.RWMutex
	healthCheck time.Duration
	stopCh      chan struct{}
}

// Client 封装gRPC客户端连接及其元数据
type Client struct {
	Conn         *grpc.ClientConn
	ServiceName  string
	Methods      map[string]MethodInfo
	LastRefresh  time.Time
	Healthy      bool
	ReflectionCl grpc_reflection_v1alpha.ServerReflectionClient
}

// MethodInfo 存储gRPC方法元数据
type MethodInfo struct {
	Name           string
	InputType      string
	OutputType     string
	IsServerStream bool
	IsClientStream bool
}

// NewClientManager 创建新的gRPC客户端管理器
func NewClientManager(conf config.GRPCConfig) *ClientManager {
	cm := &ClientManager{
		clients:     make(map[string]*Client),
		config:      conf,
		healthCheck: time.Duration(conf.HealthCheckInterval) * time.Second,
		stopCh:      make(chan struct{}),
	}

	// 启动健康检查
	go cm.healthCheckLoop()

	return cm
}

// GetClient 获取或创建gRPC客户端
func (cm *ClientManager) GetClient(ctx context.Context, serviceName, target string) (*Client, error) {
	// 先检查缓存中是否已有连接
	cm.mu.RLock()
	client, exists := cm.clients[serviceName]
	cm.mu.RUnlock()

	// 如果连接存在且健康，直接返回
	if exists && client.Healthy {
		return client, nil
	}

	// 创建新连接
	return cm.createClient(ctx, serviceName, target)
}

// createClient 创建新的gRPC客户端
func (cm *ClientManager) createClient(ctx context.Context, serviceName, target string) (*Client, error) {
	// 设置gRPC选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(cm.config.KeepAliveTime) * time.Second,
			Timeout:             time.Duration(cm.config.KeepAliveTimeout) * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// 设置消息大小限制
	if cm.config.MaxRecvMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cm.config.MaxRecvMsgSize),
		))
	}
	if cm.config.MaxSendMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(cm.config.MaxSendMsgSize),
		))
	}

	// 创建连接超时上下文
	dialCtx, cancel := context.WithTimeout(ctx, time.Duration(cm.config.DialTimeout)*time.Second)
	defer cancel()

	// 创建gRPC连接
	conn, err := grpc.DialContext(dialCtx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("无法连接gRPC服务 %s: %w", serviceName, err)
	}

	// 创建反射客户端，用于获取服务定义
	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)

	// 创建客户端对象
	client := &Client{
		Conn:         conn,
		ServiceName:  serviceName,
		Methods:      make(map[string]MethodInfo),
		LastRefresh:  time.Now(),
		Healthy:      true,
		ReflectionCl: reflectionClient,
	}

	// 获取服务方法信息
	if err := cm.fetchServiceInfo(client); err != nil {
		conn.Close()
		return nil, fmt.Errorf("无法获取服务 %s 信息: %w", serviceName, err)
	}

	// 保存到缓存
	cm.mu.Lock()
	cm.clients[serviceName] = client
	cm.mu.Unlock()

	return client, nil
}

// fetchServiceInfo 获取服务定义信息
func (cm *ClientManager) fetchServiceInfo(client *Client) error {
	// 这里应使用gRPC反射获取方法信息
	// 简单实现：可以先不填充具体的方法元数据，后续扩展
	// 实际实现应当使用反射客户端来获取服务描述符
	return nil
}

// healthCheckLoop 执行周期性健康检查
func (cm *ClientManager) healthCheckLoop() {
	ticker := time.NewTicker(cm.healthCheck)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.checkAllClients()
		case <-cm.stopCh:
			return
		}
	}
}

// checkAllClients 检查所有客户端连接的健康状态
func (cm *ClientManager) checkAllClients() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for name, client := range cm.clients {
		state := client.Conn.GetState()
		wasHealthy := client.Healthy

		// 更新健康状态
		client.Healthy = (state == connectivity.Ready || state == connectivity.Idle)

		// 如果状态从健康变为不健康，尝试重连
		if wasHealthy && !client.Healthy {
			// 记录连接问题
			fmt.Printf("gRPC客户端 %s 连接状态: %s\n", name, state.String())
			// 连接状态为不可用，可以在这里实现重连逻辑
		}

		// 更新最后检查时间
		client.LastRefresh = time.Now()
	}
}

// Stop 停止客户端管理器
func (cm *ClientManager) Stop() {
	close(cm.stopCh)

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 关闭所有连接
	for _, client := range cm.clients {
		client.Conn.Close()
	}
}
