package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCServerMiddleware 返回一组 gRPC 服务器中间件，用于集成 OpenTelemetry
func GRPCServerMiddleware() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(otlpUnaryServerInterceptor()),
		grpc.StreamInterceptor(otlpStreamServerInterceptor()),
	}
}

// GRPCClientMiddleware 返回一组 gRPC 客户端中间件，用于集成 OpenTelemetry
func GRPCClientMiddleware() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(otlpUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otlpStreamClientInterceptor()),
	}
}

// otlpUnaryServerInterceptor 扩展的 gRPC 一元服务器拦截器，添加更多自定义属性
func otlpUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	baseInterceptor := otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 添加服务器端拦截器的自定义逻辑
		startTime := time.Now()
		
		// 从元数据中提取租户信息
		md, ok := metadata.FromIncomingContext(ctx)
		tenantID := "unknown"
		if ok {
			if values := md.Get("X-Tenant-ID"); len(values) > 0 {
				tenantID = values[0]
			}
		}

		// 获取当前的 span
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("tenant.id", tenantID))
		
		// 调用基础拦截器
		resp, err := baseInterceptor(ctx, req, info, handler)
		
		// 添加响应信息
		span.SetAttributes(attribute.String("grpc.duration", time.Since(startTime).String()))
		if err != nil {
			s, _ := status.FromError(err)
			span.SetAttributes(attribute.String("grpc.error_code", s.Code().String()))
			span.SetAttributes(attribute.String("grpc.error_message", s.Message()))
		}
		
		return resp, err
	}
}

// otlpStreamServerInterceptor 扩展的 gRPC 流服务器拦截器，添加更多自定义属性
func otlpStreamServerInterceptor() grpc.StreamServerInterceptor {
	baseInterceptor := otelgrpc.StreamServerInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 添加服务器端拦截器的自定义逻辑
		startTime := time.Now()
		
		// 从元数据中提取租户信息
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		tenantID := "unknown"
		if ok {
			if values := md.Get("X-Tenant-ID"); len(values) > 0 {
				tenantID = values[0]
			}
		}

		// 获取当前的 span
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("tenant.id", tenantID))
		
		// 调用基础拦截器
		err := baseInterceptor(srv, ss, info, handler)
		
		// 添加响应信息
		span.SetAttributes(attribute.String("grpc.duration", time.Since(startTime).String()))
		if err != nil {
			s, _ := status.FromError(err)
			span.SetAttributes(attribute.String("grpc.error_code", s.Code().String()))
			span.SetAttributes(attribute.String("grpc.error_message", s.Message()))
		}
		
		return err
	}
}

// otlpUnaryClientInterceptor 扩展的 gRPC 一元客户端拦截器，添加更多自定义属性
func otlpUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	baseInterceptor := otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 添加客户端拦截器的自定义逻辑
		startTime := time.Now()
		
		// 确保上下文中有租户信息传递
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		
		// 从原上下文获取租户ID并传递
		incomingMD, hasIncoming := metadata.FromIncomingContext(ctx)
		if hasIncoming {
			if values := incomingMD.Get("X-Tenant-ID"); len(values) > 0 {
				md.Set("X-Tenant-ID", values[0])
			}
		}
		
		// 设置更新后的元数据
		ctx = metadata.NewOutgoingContext(ctx, md)
		
		// 获取当前的 span
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			attribute.String("grpc.client_method", method),
			attribute.String("grpc.target", cc.Target()),
		)
		
		// 调用基础拦截器
		err := baseInterceptor(ctx, method, req, reply, cc, invoker, opts...)
		
		// 添加响应信息
		span.SetAttributes(attribute.String("grpc.duration", time.Since(startTime).String()))
		if err != nil {
			s, _ := status.FromError(err)
			span.SetAttributes(attribute.String("grpc.error_code", s.Code().String()))
			span.SetAttributes(attribute.String("grpc.error_message", s.Message()))
		}
		
		return err
	}
}

// otlpStreamClientInterceptor 扩展的 gRPC 流客户端拦截器，添加更多自定义属性
func otlpStreamClientInterceptor() grpc.StreamClientInterceptor {
	baseInterceptor := otelgrpc.StreamClientInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// 添加客户端拦截器的自定义逻辑
		
		// 确保上下文中有租户信息传递
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		
		// 从原上下文获取租户ID并传递
		incomingMD, hasIncoming := metadata.FromIncomingContext(ctx)
		if hasIncoming {
			if values := incomingMD.Get("X-Tenant-ID"); len(values) > 0 {
				md.Set("X-Tenant-ID", values[0])
			}
		}
		
		// 设置更新后的元数据
		ctx = metadata.NewOutgoingContext(ctx, md)
		
		// 获取当前的 span
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			attribute.String("grpc.client_method", method),
			attribute.String("grpc.target", cc.Target()),
		)
		
		// 调用基础拦截器
		return baseInterceptor(ctx, desc, cc, method, streamer, opts...)
	}
}
