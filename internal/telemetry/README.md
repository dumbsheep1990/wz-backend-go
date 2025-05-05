# OpenTelemetry 分布式追踪集成使用指南

本文档介绍如何在微服务架构中使用 OpenTelemetry 进行分布式链路追踪，以实现对请求流程的完整监控与可观测性。

## 一、背景介绍

OpenTelemetry 是一个开源的可观测性框架，用于生成、收集和导出遥测数据（包括追踪、度量和日志）。在微服务架构中，它可以帮助我们：

- 跟踪请求在多个服务间的流转
- 识别性能瓶颈和延迟问题
- 快速定位错误和异常
- 提供系统行为的全面视图

## 二、集成架构

我们的 OpenTelemetry 集成架构如下：

1. **API网关层**：捕获所有入站请求，创建初始追踪 Span
2. **服务间通信层**：确保追踪上下文在服务间正确传递
3. **服务实现层**：为每个服务添加自定义 Span 以记录业务操作
4. **持久化层**：追踪数据库和缓存操作
5. **收集与可视化**：通过 Collector 收集数据并发送到后端存储

## 三、快速开始

### 1. 微服务初始化

每个微服务都应在启动时初始化 OpenTelemetry：

```go
// 初始化 OpenTelemetry
tp, err := telemetry.InitTracer(
    "your-service-name",        // 服务名称 
    "1.0.0",                    // 服务版本
    "development",              // 环境
    "otlp",                     // 导出器类型
    "localhost:4317",           // 收集器地址
)
if err != nil {
    log.Printf("警告: OpenTelemetry初始化失败: %v", err)
}

// 确保程序退出时关闭追踪提供者
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    telemetry.Shutdown(ctx, tp)
}()
```

### 2. API 网关中的请求追踪

API 网关已集成了 OpenTelemetry 中间件，会自动为每个传入请求创建一个 Span：

```go
// 在 gin 路由中添加追踪中间件
router.Use(middleware.Tracing("api-gateway"))
```

### 3. gRPC 服务集成

对于 gRPC 服务，使用提供的拦截器：

```go
// 创建 gRPC 服务器
s := grpc.NewServer(
    telemetry.GRPCServerMiddleware()...
)
```

### 4. gRPC 客户端集成

确保服务间调用传递追踪上下文：

```go
// 创建带有追踪功能的 gRPC 客户端
conn, err := telemetry.NewGrpcClientWithContext(ctx, "localhost:50051")
if err != nil {
    return nil, err
}
defer conn.Close()

// 使用客户端
client := pb.NewYourServiceClient(conn)
```

### 5. 手动创建 Span

对于特定业务逻辑，可以手动创建 Span：

```go
// 使用当前上下文创建子 Span
ctx, span := telemetry.StartSpan(ctx, "your-operation-name")
defer span.End()

// 添加属性
span.SetAttributes(attribute.String("key", "value"))

// 添加事件
span.AddEvent("event-name", trace.WithAttributes(
    attribute.String("event.key", "value"),
))

// 标记错误
if err != nil {
    telemetry.RecordError(span, err, "操作失败")
}
```

## 四、最佳实践

1. **命名规范**：
   - 服务名称：使用小写字母和连字符，如 `user-service`
   - Span 名称：使用动词短语，如 `get-user-profile`，`process-payment`

2. **添加合适的属性**：
   - 为每个 Span 添加具有业务意义的属性
   - 对于 HTTP 请求，添加方法、URL、状态码等
   - 对于数据库操作，添加查询类型、表名等

3. **错误处理**：
   - 对于错误情况，记录详细的错误信息
   - 使用 `telemetry.RecordError` 确保一致的错误记录

4. **传播上下文**：
   - 在异步操作中，确保正确传递追踪上下文
   - 使用 `telemetry.ExtractTraceInfoFromRequest` 和 `telemetry.InjectTraceInfoToRequest` 在非标准场景中传递上下文

5. **采样策略**：
   - 开发环境：使用全量采样 (1.0)
   - 生产环境：根据流量调整采样率 (0.1 - 0.5)

## 五、故障排查

1. **追踪数据没有出现**：
   - 检查收集器是否运行正常
   - 确认服务配置中的 CollectorURL 正确
   - 验证服务间是否正确传递上下文

2. **Span 不完整**：
   - 确保所有 Span 正确结束（调用 `span.End()`）
   - 检查是否有未捕获的异常导致 Span 提前终止

3. **性能问题**：
   - 检查 BatchSpanProcessor 的配置
   - 考虑调整 MaxQueueSize 和 BatchTimeout

## 六、扩展阅读

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [Go 语言 OpenTelemetry SDK 文档](https://pkg.go.dev/go.opentelemetry.io/otel)
- [Jaeger UI 使用指南](https://www.jaegertracing.io/docs/latest/getting-started/)
