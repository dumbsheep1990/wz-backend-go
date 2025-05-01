package telemetry

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// 以下是OpenTelemetry使用示例

// 示例: 在HTTP处理器中使用OpenTelemetry
func ExampleHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求中提取追踪上下文
		ctx := ExtractTraceInfoFromRequest(r)
		
		// 创建新的Span
		ctx, span := StartSpan(ctx, "example-http-handler")
		defer span.End()
		
		// 添加请求相关属性
		AddHTTPRequestAttributesToSpan(span, r)
		
		// 添加自定义属性
		span.SetAttributes(
			attribute.String("custom.attribute", "value"),
			attribute.Int64("request.received.time", time.Now().UnixNano()),
		)
		
		// 执行业务逻辑
		err := processRequest(ctx)
		
		// 设置Span状态
		SetSpanStatus(span, err)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		
		// 向响应头中添加追踪ID
		traceID := GetTraceIDFromContext(ctx)
		w.Header().Set("X-Trace-ID", traceID)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("请求处理成功"))
	}
}

// 示例: 在gRPC服务器中使用OpenTelemetry
func ExampleGRPCServerInterceptorUsage() {
	// 在创建gRPC服务器时添加拦截器
	// 下面是伪代码，展示如何在实际中使用
	/*
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
		grpc.StreamInterceptor(StreamServerInterceptor()),
	)
	*/
	
	// 在gRPC实现方法中使用Span
	/*
	func (s *myService) MyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
		ctx, span := StartSpan(ctx, "my-method")
		defer span.End()
		
		// 添加请求相关属性
		span.SetAttributes(attribute.String("request.id", req.GetRequestId()))
		
		// 执行业务逻辑...
		result, err := s.businessLogic(ctx, req)
		
		// 设置Span状态
		SetSpanStatus(span, err)
		
		return result, err
	}
	*/
}

// 示例: 在gRPC客户端中使用OpenTelemetry
func ExampleGRPCClientInterceptorUsage() {
	// 在创建gRPC客户端连接时添加拦截器
	// 下面是伪代码，展示如何在实际中使用
	/*
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithUnaryInterceptor(UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(StreamClientInterceptor()),
	)
	*/
}

// 示例: 初始化OpenTelemetry追踪器
func ExampleInitializeTracer() {
	// 初始化追踪器
	tp, err := InitTracer(
		"example-service",
		WithServiceVersion("1.0.0"),
		WithEnvironment("dev"),
		WithExporterType("otlp"),
		WithExporterEndpoint("localhost:4317"),
		WithInsecure(true),
		WithSamplerType("trace_id_ratio"),
		WithSamplerRatio(0.5),
	)
	if err != nil {
		log.Fatalf("初始化追踪器失败: %v", err)
	}
	
	// 确保程序结束前关闭追踪器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer tp.Shutdown(ctx)
}

// 示例: 使用追踪Context包装函数
func ExampleWithTraceContext() {
	// 使用WithTraceContext包装函数
	wrappedFunc := WithTraceContext("process-data", func(ctx context.Context) error {
		// 这里是被追踪的函数逻辑
		return processData(ctx)
	})
	
	// 调用包装后的函数
	ctx := context.Background()
	err := wrappedFunc(ctx)
	if err != nil {
		log.Printf("处理数据失败: %v", err)
	}
}

// 示例: 手动创建和管理Span
func ExampleManualSpanManagement(ctx context.Context) {
	// 创建一个子Span
	ctx, span := StartSpan(ctx, "manual-span-example")
	defer span.End()
	
	// 添加属性
	span.SetAttributes(
		attribute.String("operation", "manual-example"),
		attribute.Int("iteration", 42),
	)
	
	// 记录事件
	span.AddEvent("开始处理")
	
	// 执行一些工作
	time.Sleep(100 * time.Millisecond)
	
	// 创建另一个嵌套的子Span
	childCtx, childSpan := StartSpan(ctx, "child-operation")
	// 模拟子操作
	time.Sleep(50 * time.Millisecond)
	// 结束子Span
	childSpan.End()
	
	// 继续主Span的工作
	time.Sleep(100 * time.Millisecond)
	
	// 记录另一个事件
	span.AddEvent("处理完成")
}

// 模拟业务逻辑函数
func processRequest(ctx context.Context) error {
	// 创建子Span
	ctx, span := StartSpan(ctx, "process-request")
	defer span.End()
	
	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)
	
	// 调用数据库
	if err := simulateDatabaseQuery(ctx); err != nil {
		return err
	}
	
	// 调用外部服务
	return simulateExternalServiceCall(ctx)
}

// 模拟数据库查询
func simulateDatabaseQuery(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "database-query")
	defer span.End()
	
	// 添加数据库相关属性
	span.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.name", "users"),
		attribute.String("db.operation", "select"),
	)
	
	// 模拟数据库查询时间
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// 模拟外部服务调用
func simulateExternalServiceCall(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "external-service-call")
	defer span.End()
	
	// 添加服务调用相关属性
	span.SetAttributes(
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.service", "payment-service"),
		attribute.String("rpc.method", "ProcessPayment"),
	)
	
	// 模拟服务调用时间
	time.Sleep(75 * time.Millisecond)
	
	return nil
}

// 模拟数据处理
func processData(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "data-processing")
	defer span.End()
	
	// 添加处理相关属性
	span.SetAttributes(
		attribute.String("processing.type", "batch"),
		attribute.Int64("processing.items", 100),
	)
	
	// 模拟处理时间
	time.Sleep(150 * time.Millisecond)
	
	return nil
}
