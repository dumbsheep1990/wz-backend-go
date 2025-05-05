package telemetry

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClientFactory 是一个创建集成了OpenTelemetry的gRPC客户端的工厂
type GrpcClientFactory struct {
	dialOptions []grpc.DialOption
}

// NewGrpcClientFactory 创建一个新的gRPC客户端工厂
func NewGrpcClientFactory() *GrpcClientFactory {
	// 默认的gRPC客户端选项
	dialOptions := []grpc.DialOption{
		// 添加OpenTelemetry拦截器
		grpc.WithUnaryInterceptor(otlpUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otlpStreamClientInterceptor()),
	}

	return &GrpcClientFactory{
		dialOptions: dialOptions,
	}
}

// WithInsecure 添加不安全的连接选项（仅用于开发环境）
func (f *GrpcClientFactory) WithInsecure() *GrpcClientFactory {
	f.dialOptions = append(f.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return f
}

// WithTimeout 添加超时选项
func (f *GrpcClientFactory) WithTimeout(timeout time.Duration) *GrpcClientFactory {
	f.dialOptions = append(f.dialOptions, grpc.WithTimeout(timeout))
	return f
}

// WithBlock 添加阻塞连接选项
func (f *GrpcClientFactory) WithBlock() *GrpcClientFactory {
	f.dialOptions = append(f.dialOptions, grpc.WithBlock())
	return f
}

// WithOptions 添加自定义选项
func (f *GrpcClientFactory) WithOptions(opts ...grpc.DialOption) *GrpcClientFactory {
	f.dialOptions = append(f.dialOptions, opts...)
	return f
}

// Dial 连接到gRPC服务器
func (f *GrpcClientFactory) Dial(target string) (*grpc.ClientConn, error) {
	return grpc.Dial(target, f.dialOptions...)
}

// DialContext 使用上下文连接到gRPC服务器
func (f *GrpcClientFactory) DialContext(ctx context.Context, target string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target, f.dialOptions...)
}

// 便捷函数，用于快速创建具有OpenTelemetry支持的gRPC客户端连接
func NewGrpcClient(target string) (*grpc.ClientConn, error) {
	factory := NewGrpcClientFactory().WithInsecure()
	return factory.Dial(target)
}

// 使用上下文创建具有OpenTelemetry支持的gRPC客户端连接
func NewGrpcClientWithContext(ctx context.Context, target string) (*grpc.ClientConn, error) {
	factory := NewGrpcClientFactory().WithInsecure()
	return factory.DialContext(ctx, target)
}
