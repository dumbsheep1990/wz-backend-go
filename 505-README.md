# wz-backend-go 更新日志 (2025-05-05)

## 本次更新内容：基于OpenTelemetry的分布式链路追踪系统实现

### 实现概述

今日更新主要针对以下功能：

- 对每个微服务添加OpenTelemetry初始化代码
- 在API网关添加传入请求的追踪
- 在服务间调用添加上下文传递

### 功能实现清单

1. **OpenTelemetry基础架构集成**
   
   - 追踪提供者（TracerProvider）初始化机制
   - 可配置导出器（支持OTLP、Zipkin和控制台输出）
   - 采样器配置与全局传播器设置
   - 优雅关闭的处理机制

2. **API网关请求追踪**
   
   - 基于Gin的OpenTelemetry中间件
   - 自动捕获和记录请求属性
   - 追踪上下文的提取和注入
   - HTTP请求参数和响应状态的记录

3. **微服务间调用链路追踪**
   
   - gRPC服务器和客户端拦截器
   - 上下文传递机制（Trace ID和Span ID）
   - 租户信息在服务间的传递
   - 自定义属性和事件记录

4. **统一错误处理与记录**
   
   - 标准化的错误记录函数
   - 自动捕获堆栈跟踪
   - 错误属性提取与分类
   - 错误状态与事件记录

### 代码结构

```
internal/
├── telemetry/
│   ├── config.go            # OpenTelemetry配置结构
│   ├── tracer.go            # 追踪提供者实现
│   ├── middleware.go        # 基础HTTP和gRPC中间件
│   ├── grpc_middleware.go   # 增强的gRPC拦截器
│   ├── grpc_client.go       # 具有追踪功能的gRPC客户端
│   ├── utils.go             # 工具函数和错误处理
│   ├── init.go              # 初始化便捷函数
│   └── README.md            # 使用指南和最佳实践
├── gateway/
│   ├── middleware/
│   │   └── tracing.go       # 网关追踪中间件
│   └── main.go              # 网关初始化代码
└── cmd/
    └── rpc/
        ├── user/
        │   └── main.go      # 用户服务初始化代码
        └── [其他微服务]
            └── main.go      # 其他微服务初始化代码
```

### 模块详细说明

#### 1. OpenTelemetry基础架构 (internal/telemetry/)

- **TracerProvider**: 追踪提供者的创建和配置
- **导出器选择**: 支持多种遥测数据导出方式
- **资源和采样器**: 自定义服务信息和采样策略
- **批处理器**: 优化遥测数据传输效率

关键结构：

```go
// TracerProvider 是OpenTelemetry追踪提供者的包装
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    exporter sdktrace.SpanExporter
}

// Config 表示OpenTelemetry的配置
type Config struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    Attributes     map[string]string
    ExporterConfig ExporterConfig
    SamplerConfig  SamplerConfig
}

// InitTracer 初始化并返回 OpenTelemetry 追踪提供者
func InitTracer(serviceName, serviceVersion, environment string, 
                exporterType string, exporterEndpoint string) (*TracerProvider, error)
```

#### 2. API网关追踪中间件 (gateway/middleware/)

- **Tracing中间件**: 处理所有入站HTTP请求
- **上下文处理**: 提取和注入追踪信息
- **属性记录**: 捕获请求和响应信息
- **状态管理**: 基于HTTP状态码设置Span状态

关键函数：

```go
// Tracing 返回请求追踪中间件
func Tracing(serviceName string) gin.HandlerFunc {
    // 从请求头中提取追踪上下文
    // 开始一个新的追踪Span
    // 添加请求属性到Span
    // 将追踪上下文注入请求中
    // 处理响应状态和计算处理时间
}
```

#### 3. 微服务间调用追踪 (internal/telemetry/)

- **gRPC拦截器**: 服务器和客户端追踪拦截器
- **上下文传递**: 确保追踪上下文在服务间正确传递
- **元数据处理**: 通过gRPC元数据传递追踪信息
- **客户端工厂**: 便捷创建具有追踪功能的客户端

关键结构与接口：

```go
// GRPCServerMiddleware 返回一组 gRPC 服务器中间件，用于集成 OpenTelemetry
func GRPCServerMiddleware() []grpc.ServerOption

// GRPCClientMiddleware 返回一组 gRPC 客户端中间件，用于集成 OpenTelemetry
func GRPCClientMiddleware() []grpc.DialOption

// GrpcClientFactory 是一个创建集成了OpenTelemetry的gRPC客户端的工厂
type GrpcClientFactory struct {
    dialOptions []grpc.DialOption
}

// NewGrpcClientWithContext 使用上下文创建具有OpenTelemetry支持的gRPC客户端连接
func NewGrpcClientWithContext(ctx context.Context, target string) (*grpc.ClientConn, error)
```

#### 4. 错误处理与记录 (internal/telemetry/utils.go)

- **错误记录**: 统一的错误记录机制
- **堆栈捕获**: 自动捕获并记录堆栈信息
- **属性提取**: 从错误中提取有用信息
- **错误级别**: 支持不同严重级别的错误记录

关键函数：

```go
// RecordError 记录错误信息到Span
func RecordError(span trace.Span, err error, msg string, attrs ...attribute.KeyValue) {
    // 设置Span状态为错误
    // 记录错误事件
    // 添加堆栈跟踪
    // 添加自定义属性
}

// ExtractErrorDetails 从错误中提取详细信息，并将其转换为属性列表
func ExtractErrorDetails(err error) []attribute.KeyValue
```

### 优化效果

1. **系统可观测性提升**:
   
   - 全链路请求追踪，从网关到数据库
   - 微服务间调用关系的可视化
   - 性能瓶颈和异常点快速定位

2. **问题排查效率提高**:
   
   - 详细的请求处理链路记录
   - 统一的错误处理与堆栈捕获
   - 跨服务问题的根因快速分析

3. **开发与运维协作增强**:
   
   - 标准化的追踪数据结构
   - 多种导出器支持与第三方工具集成
   - 便于开发人员和运维人员共同分析系统

### 后续计划

- 实现服务网格集成（Istio/Linkerd）
- 添加业务度量指标监控
- 实现自定义采样策略
- 完善告警机制与异常检测
- 集成日志系统，实现三大可观测性支柱的统一
