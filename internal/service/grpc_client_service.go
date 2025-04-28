package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"github.com/go-redis/redis/v8"
)

// GRPCClientService gRPC客户端服务接口
type GRPCClientService interface {
	// 创建gRPC连接
	CreateConnection(serviceName string, target string) (*grpc.ClientConn, error)
	// 使用带有追踪信息的上下文创建请求
	ContextWithTrace(ctx context.Context, tenantID string) context.Context
	// 关闭所有连接
	CloseAll() error
}

type grpcClientService struct {
	connections map[string]*grpc.ClientConn
	redis       *redis.Client
}

// ClientOptions gRPC客户端选项
type ClientOptions struct {
	DialTimeout     time.Duration
	KeepAliveTime   time.Duration
	KeepAliveTimeout time.Duration
}

// DefaultClientOptions 默认客户端选项
var DefaultClientOptions = ClientOptions{
	DialTimeout:      5 * time.Second,
	KeepAliveTime:    30 * time.Second,
	KeepAliveTimeout: 10 * time.Second,
}

// NewGRPCClientService 创建gRPC客户端服务
func NewGRPCClientService(redis *redis.Client) GRPCClientService {
	return &grpcClientService{
		connections: make(map[string]*grpc.ClientConn),
		redis:       redis,
	}
}

// CreateConnection 创建gRPC连接
func (s *grpcClientService) CreateConnection(serviceName string, target string) (*grpc.ClientConn, error) {
	// 检查是否已存在连接
	if conn, ok := s.connections[serviceName]; ok {
		return conn, nil
	}

	// 创建新连接
	opts := DefaultClientOptions
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 开发环境使用不安全连接，生产环境应使用TLS
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                opts.KeepAliveTime,
			Timeout:             opts.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
		grpc.WithTimeout(opts.DialTimeout),
		grpc.WithUnaryInterceptor(s.unaryClientInterceptor),
		grpc.WithStreamInterceptor(s.streamClientInterceptor),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", serviceName, err)
	}

	// 保存连接
	s.connections[serviceName] = conn
	return conn, nil
}

// ContextWithTrace 使用带有追踪信息的上下文创建请求
func (s *grpcClientService) ContextWithTrace(ctx context.Context, tenantID string) context.Context {
	md := metadata.Pairs(
		"x-tenant-id", tenantID,
		"x-request-id", generateRequestID(),
		"x-timestamp", fmt.Sprintf("%d", time.Now().UnixNano()),
	)
	return metadata.NewOutgoingContext(ctx, md)
}

// CloseAll 关闭所有连接
func (s *grpcClientService) CloseAll() error {
	var lastErr error
	for name, conn := range s.connections {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close connection to %s: %v", name, err)
		}
	}
	return lastErr
}

// unaryClientInterceptor 一元RPC客户端拦截器
func (s *grpcClientService) unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// 开始计时
	startTime := time.Now()
	
	// 执行调用
	err := invoker(ctx, method, req, reply, cc, opts...)
	
	// 计算耗时
	duration := time.Since(startTime)
	
	// 记录调用信息到Redis，可用于监控和统计
	logKey := fmt.Sprintf("grpc:log:%s:%d", method, startTime.Unix())
	logData := fmt.Sprintf(
		`{"method":"%s","duration_ms":%d,"error":%t,"timestamp":"%s"}`,
		method,
		duration.Milliseconds(),
		err != nil,
		startTime.Format(time.RFC3339),
	)
	
	// 异步记录日志，不影响正常调用
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		s.redis.Set(ctx, logKey, logData, 48*time.Hour)
	}()
	
	return err
}

// streamClientInterceptor 流式RPC客户端拦截器
func (s *grpcClientService) streamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	// 开始计时
	startTime := time.Now()
	
	// 执行调用
	stream, err := streamer(ctx, desc, cc, method, opts...)
	
	// 记录调用信息到Redis，可用于监控和统计
	logKey := fmt.Sprintf("grpc:stream:%s:%d", method, startTime.Unix())
	logData := fmt.Sprintf(
		`{"method":"%s","type":"stream","error":%t,"timestamp":"%s"}`,
		method,
		err != nil,
		startTime.Format(time.RFC3339),
	)
	
	// 异步记录日志，不影响正常调用
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		s.redis.Set(ctx, logKey, logData, 48*time.Hour)
	}()
	
	return stream, err
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}
