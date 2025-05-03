# wz-backend-go 更新日志 (2025-05-01)

## 本次更新内容：基于OpenTelemetry的分布式链路追踪系统实现

### 实现概述

今日更新主要针对以下功能：

- 分布式链路追踪系统
- 追踪上下文传递机制
- HTTP和gRPC请求跟踪
- 性能指标收集基础设施

### 功能实现清单

1. **OpenTelemetry核心配置模块**
   
   - 多导出器类型支持 (OTLP, Jaeger, Zipkin等)
   - 可配置的采样策略
   - 批处理与队列管理
   - 资源标签定制

2. **追踪器管理系统**
   
   - 追踪器初始化与全局注册
   - 资源标签配置
   - 上下文传播器设置
   - 优雅关闭机制

3. **HTTP与gRPC中间件集成**
   
   - HTTP服务器追踪中间件
   - gRPC服务器一元与流拦截器
   - gRPC客户端一元与流拦截器
   - 追踪上下文的提取与注入工具

4. **实用工具集**
   
   - Span状态和属性管理
   - 错误记录与处理
   - 业务属性添加工具
   - 函数跟踪包装器

### 代码结构

```
internal/
└── telemetry/
    ├── config.go     # OpenTelemetry配置定义
    ├── tracer.go     # 追踪器初始化与管理
    ├── middleware.go # HTTP与gRPC中间件
    ├── utils.go      # 工具函数集合
    └── example.go    # 使用示例
```

### 模块详细说明

#### 1. 配置模块 (config.go)

- **Config**: 定义OpenTelemetry的配置结构
- **ExporterConfig**: 定义数据导出器的配置
- **SamplerConfig**: 定义采样策略的配置
- **DefaultConfig**: 提供合理的默认配置

关键结构：

```go
// Config 表示OpenTelemetry的配置
type Config struct {
    ServiceName    string            // 服务名称
    ServiceVersion string            // 服务版本
    Environment    string            // 环境(dev, staging, prod)
    Attributes     map[string]string // 服务级别的属性
    ExporterConfig ExporterConfig    // 数据导出器配置
    SamplerConfig  SamplerConfig     // 采样配置
}
```

#### 2. 追踪器管理模块 (tracer.go)

- **TracerProvider**: 封装SDK追踪提供者的包装
- **资源创建**: 基于配置创建OpenTelemetry资源
- **导出器创建**: 根据配置创建对应类型的导出器
- **采样器创建**: 支持多种采样策略
- **函数选项模式**: 提供灵活的配置API

关键接口：

```go
// TracerProvider 是OpenTelemetry追踪提供者的包装
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    exporter sdktrace.SpanExporter
}

// InitTracer 初始化全局追踪器
func InitTracer(serviceName string, opts ...Option) (*TracerProvider, error)
```

#### 3. 中间件模块 (middleware.go)

- **HTTPMiddleware**: 提供HTTP请求的追踪中间件
- **gRPC拦截器**: 提供服务端和客户端的拦截器
- **上下文工具**: 提供追踪上下文的提取和注入
- **追踪ID获取**: 从上下文中获取追踪ID和SpanID

关键函数：

```go
// HTTPMiddleware 是用于HTTP请求的OpenTelemetry中间件
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler

// UnaryServerInterceptor 返回一个gRPC服务器一元拦截器
func UnaryServerInterceptor() grpc.UnaryServerInterceptor

// ExtractTraceInfoFromRequest 从HTTP请求中提取追踪信息
func ExtractTraceInfoFromRequest(req *http.Request) context.Context
```

#### 4. 实用工具模块 (utils.go)

- **Span状态工具**: 统一设置Span状态和错误记录
- **属性管理**: 批量添加Span属性的工具
- **业务属性工具**: 添加用户ID、租户ID等业务属性
- **函数追踪**: 自动追踪函数调用及错误处理

关键函数：

```go
// SetSpanStatus 设置Span的状态
func SetSpanStatus(span trace.Span, err error)

// AddUserIDToSpan 添加用户ID到Span
func AddUserIDToSpan(span trace.Span, userID string)

// WithTraceContext 创建一个带有追踪上下文的函数包装器
func WithTraceContext(operation string, f func(ctx context.Context) error) func(ctx context.Context) error
```
