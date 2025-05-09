# OpenTelemetry分布式链路追踪实现文档

## 实现概述

本文档详细描述了在wz-backend-go项目中基于OpenTelemetry实现的分布式链路追踪系统。该系统旨在提供端到端的请求追踪、性能分析和问题诊断能力，以便开发和运维团队能够更好地理解和优化微服务架构。

## 实现组件

### 1. 配置管理 (`internal/telemetry/config.go`)

提供了灵活的OpenTelemetry配置选项，支持以下功能：

- 多种导出器类型（OTLP、Jaeger、Zipkin等）
- 可配置的采样策略（全采样、比例采样等）
- 可定制的资源属性（服务名称、环境标识等）
- 批处理和队列参数配置

配置结构设计如下：
```go
type Config struct {
    ServiceName    string            // 服务名称
    ServiceVersion string            // 服务版本
    Environment    string            // 环境(dev, staging, prod)
    Attributes     map[string]string // 服务级别的属性
    ExporterConfig ExporterConfig    // 数据导出器配置
    SamplerConfig  SamplerConfig     // 采样配置
}
```

### 2. 追踪器管理 (`internal/telemetry/tracer.go`)

实现了OpenTelemetry追踪提供者的初始化和管理：

- 基于配置创建资源标识
- 支持多种导出器的创建和配置
- 提供采样器的创建和配置
- 功能选项模式的API设计，便于使用和扩展
- 全局追踪器的初始化和关闭管理

核心结构和函数：
```go
// TracerProvider 是OpenTelemetry追踪提供者的包装
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    exporter sdktrace.SpanExporter
}

// InitTracer 初始化全局追踪器
func InitTracer(serviceName string, opts ...Option) (*TracerProvider, error)
```

### 3. 中间件集成 (`internal/telemetry/middleware.go`)

提供了与HTTP和gRPC集成的中间件：

- HTTP服务器中间件
- gRPC服务器一元和流拦截器
- gRPC客户端一元和流拦截器
- 追踪上下文的提取和注入工具
- 追踪ID和Span ID的获取工具

主要函数：
```go
// HTTPMiddleware 是用于HTTP请求的OpenTelemetry中间件
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler

// UnaryServerInterceptor 返回一个gRPC服务器一元拦截器
func UnaryServerInterceptor() grpc.UnaryServerInterceptor

// StreamServerInterceptor 返回一个gRPC服务器流拦截器
func StreamServerInterceptor() grpc.StreamServerInterceptor
```

### 4. 实用工具 (`internal/telemetry/utils.go`)

提供了一系列辅助函数，简化Span的使用：

- Span状态和属性的设置工具
- 异常记录和处理工具
- 业务相关属性（用户ID、租户ID、请求ID等）的添加工具
- 函数跟踪的包装器
- 追踪信息的格式化和调试工具

主要函数：
```go
// SetSpanStatus 设置Span的状态
func SetSpanStatus(span trace.Span, err error)

// AddUserIDToSpan 添加用户ID到Span
func AddUserIDToSpan(span trace.Span, userID string)

// WithTraceContext 创建一个带有追踪上下文的函数包装器
func WithTraceContext(operation string, f func(ctx context.Context) error) func(ctx context.Context) error
```

### 5. 使用示例 (`internal/telemetry/example.go`)

提供了完整的使用示例和最佳实践：

- HTTP处理器中的追踪示例
- gRPC服务器和客户端的追踪集成示例
- 追踪器初始化和关闭示例
- 手动Span创建和管理示例
- 业务场景的模拟示例

## 技术规范

### 依赖库
- `go.opentelemetry.io/otel`: OpenTelemetry API
- `go.opentelemetry.io/otel/sdk`: OpenTelemetry SDK
- `go.opentelemetry.io/otel/trace`: 追踪API
- `go.opentelemetry.io/otel/exporters`: 各种导出器
- `go.opentelemetry.io/contrib/instrumentation`: HTTP和gRPC集成

### 数据流向
1. 应用程序创建Span并添加属性
2. Span由TracerProvider收集和处理
3. 批处理器将Span批量发送给导出器
4. 导出器将数据发送到后端系统（Jaeger、Zipkin等）
5. 后端系统存储、处理并可视化追踪数据

### 采样策略
- 开发环境使用全采样（`"always_on"`）
- 测试环境使用按比例采样（`"trace_id_ratio"，比例为0.5`）
- 生产环境使用按比例采样（`"trace_id_ratio"，比例为0.1`）
- 支持基于父级的采样决策，保证完整调用链

## 集成指南

### 在微服务中初始化追踪器

在每个微服务的`main.go`中添加如下代码：

```go
// 初始化追踪器
tp, err := telemetry.InitTracer(
    "service-name", // 替换为您的服务名称
    telemetry.WithServiceVersion("1.0.0"),
    telemetry.WithEnvironment("dev"), // dev, staging, prod
    telemetry.WithExporterType("otlp"),
    telemetry.WithExporterEndpoint("localhost:4317"), // 替换为您的OTLP端点
)
if err != nil {
    log.Fatalf("初始化追踪器失败: %v", err)
}

// 确保程序结束前关闭追踪器
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := tp.Shutdown(ctx); err != nil {
        log.Printf("关闭追踪器时发生错误: %v", err)
    }
}()
```

### 添加HTTP中间件

```go
router := http.NewServeMux()
// 注册路由...
return telemetry.HTTPMiddleware("api-gateway")(router)
```

### 添加gRPC拦截器

服务端：
```go
server := grpc.NewServer(
    grpc.UnaryInterceptor(telemetry.UnaryServerInterceptor()),
    grpc.StreamInterceptor(telemetry.StreamServerInterceptor()),
)
```

客户端：
```go
conn, err := grpc.Dial(
    target,
    grpc.WithInsecure(),
    grpc.WithUnaryInterceptor(telemetry.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(telemetry.StreamClientInterceptor()),
)
```

### 手动添加Span

```go
func businessLogic(ctx context.Context, data interface{}) error {
    ctx, span := telemetry.StartSpan(ctx, "business-operation")
    defer span.End()
    
    // 添加业务相关属性
    telemetry.SetSpanAttributes(span, map[string]string{
        "business.operation": "process-data",
        "business.data.type": "user-profile",
    })
    
    // 执行业务逻辑...
    result, err := processData(ctx, data)
    
    // 设置Span状态
    telemetry.SetSpanStatus(span, err)
    
    return err
}
```

## 后端部署

### Jaeger（开发环境）

使用Docker启动Jaeger：
```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

Jaeger UI将在 http://localhost:16686 可用。

### OpenTelemetry Collector（生产环境）

推荐使用 OpenTelemetry Collector 作为生产环境中的中央收集点，可以灵活配置数据流向多个后端系统。

## 最佳实践

1. **命名约定**：使用统一的服务和操作命名约定
   - 服务名格式：`[项目]-[模块]-[服务]`
   - 操作名格式：`[动词]-[资源]`

2. **属性添加**：为每个业务关键Span添加以下属性
   - `tenant.id`：租户ID
   - `user.id`：用户ID（如适用）
   - `request.id`：请求ID
   - `version`：API版本

3. **错误处理**：
   - 始终记录错误信息到Span
   - 使用合适的状态码
   - 添加足够的错误上下文

4. **采样策略**：
   - 开发环境：100%采样
   - 测试环境：50%采样
   - 生产环境：10%采样，特定关键路径100%采样

5. **性能注意事项**：
   - 避免添加过多属性到高频Span
   - 使用批处理机制
   - 监控追踪系统自身的性能影响
