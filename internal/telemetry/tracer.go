package telemetry

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TracerProvider 是OpenTelemetry追踪提供者的包装
type TracerProvider struct {
	provider *sdktrace.TracerProvider
	exporter sdktrace.SpanExporter
}

// NewTracerProvider 创建并配置一个新的TracerProvider
func NewTracerProvider(config *Config) (*TracerProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	// 创建资源
	res, err := createResource(config)
	if err != nil {
		return nil, fmt.Errorf("创建资源失败: %w", err)
	}

	// 创建导出器
	exporter, err := createExporter(config)
	if err != nil {
		return nil, fmt.Errorf("创建导出器失败: %w", err)
	}

	// 创建批处理器
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithMaxExportBatchSize(config.ExporterConfig.MaxExportBatchSize),
		sdktrace.WithBatchTimeout(config.ExporterConfig.BatchTimeout),
		sdktrace.WithMaxQueueSize(config.ExporterConfig.MaxQueueSize),
	)

	// 创建采样器
	sampler, err := createSampler(config)
	if err != nil {
		return nil, fmt.Errorf("创建采样器失败: %w", err)
	}

	// 创建追踪提供者
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)

	// 设置全局追踪提供者
	otel.SetTracerProvider(tracerProvider)

	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: tracerProvider,
		exporter: exporter,
	}, nil
}

// Tracer 返回一个追踪器
func (tp *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return tp.provider.Tracer(name, opts...)
}

// Shutdown 关闭追踪提供者并刷新所有待处理的Spans
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	return tp.provider.Shutdown(ctx)
}

// createResource 创建OpenTelemetry资源
func createResource(config *Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	}

	// 添加自定义属性
	for k, v := range config.Attributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	), nil
}

// createExporter 根据配置创建导出器
func createExporter(config *Config) (sdktrace.SpanExporter, error) {
	switch config.ExporterConfig.Type {
	case "stdout":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	case "otlp":
		ctx, cancel := context.WithTimeout(context.Background(), config.ExporterConfig.Timeout)
		defer cancel()

		if config.ExporterConfig.Endpoint == "" {
			return nil, fmt.Errorf("OTLP导出器需要指定端点")
		}

		// OTLP导出器支持gRPC和HTTP两种传输方式
		if config.ExporterConfig.Headers != nil && len(config.ExporterConfig.Headers) > 0 {
			// 使用HTTP传输
			opts := []otlptracehttp.Option{
				otlptracehttp.WithEndpoint(config.ExporterConfig.Endpoint),
				otlptracehttp.WithTimeout(config.ExporterConfig.Timeout),
			}
			
			// 添加自定义请求头
			headers := make(map[string]string)
			for k, v := range config.ExporterConfig.Headers {
				headers[k] = v
			}
			if len(headers) > 0 {
				opts = append(opts, otlptracehttp.WithHeaders(headers))
			}
			
			// 是否使用不安全连接
			if config.ExporterConfig.Insecure {
				opts = append(opts, otlptracehttp.WithInsecure())
			}
			
			return otlptrace.New(ctx, otlptracehttp.NewClient(opts...))
		} else {
			// 使用gRPC传输
			opts := []otlptracegrpc.Option{
				otlptracegrpc.WithEndpoint(config.ExporterConfig.Endpoint),
				otlptracegrpc.WithTimeout(config.ExporterConfig.Timeout),
			}
			
			// 是否使用不安全连接
			if config.ExporterConfig.Insecure {
				opts = append(opts, otlptracegrpc.WithInsecure())
			} else {
				// 使用TLS连接
				opts = append(opts, otlptracegrpc.WithTLSCredentials(grpc.WithTransportCredentials(insecure.NewCredentials())))
			}
			
			return otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
		}
	case "zipkin":
		return zipkin.New(config.ExporterConfig.Endpoint)
	default:
		return nil, fmt.Errorf("不支持的导出器类型: %s", config.ExporterConfig.Type)
	}
}

// createSampler 创建采样器
func createSampler(config *Config) (sdktrace.Sampler, error) {
	var sampler sdktrace.Sampler

	switch config.SamplerConfig.Type {
	case "always_on":
		sampler = sdktrace.AlwaysSample()
	case "always_off":
		sampler = sdktrace.NeverSample()
	case "trace_id_ratio":
		if config.SamplerConfig.Ratio < 0 || config.SamplerConfig.Ratio > 1 {
			return nil, fmt.Errorf("采样比例必须在 0.0 到 1.0 之间")
		}
		sampler = sdktrace.TraceIDRatioBased(config.SamplerConfig.Ratio)
	default:
		return nil, fmt.Errorf("不支持的采样器类型: %s", config.SamplerConfig.Type)
	}

	// 如果开启了父级采样，则使用父级采样器
	if config.SamplerConfig.ParentBased {
		return sdktrace.ParentBased(sampler), nil
	}

	return sampler, nil
}

// InitTracer 初始化全局追踪器
func InitTracer(serviceName string, opts ...Option) (*TracerProvider, error) {
	config := DefaultConfig(serviceName)
	
	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	tp, err := NewTracerProvider(config)
	if err != nil {
		return nil, err
	}

	log.Printf("[Telemetry] 已初始化OpenTelemetry追踪器: 服务名=%s, 导出器=%s, 端点=%s",
		config.ServiceName, config.ExporterConfig.Type, config.ExporterConfig.Endpoint)
	
	return tp, nil
}

// Option 是配置函数选项
type Option func(*Config)

// WithServiceVersion 设置服务版本
func WithServiceVersion(version string) Option {
	return func(c *Config) {
		c.ServiceVersion = version
	}
}

// WithEnvironment 设置环境
func WithEnvironment(env string) Option {
	return func(c *Config) {
		c.Environment = env
	}
}

// WithExporterType 设置导出器类型
func WithExporterType(exporterType string) Option {
	return func(c *Config) {
		c.ExporterConfig.Type = exporterType
	}
}

// WithExporterEndpoint 设置导出器端点
func WithExporterEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.ExporterConfig.Endpoint = endpoint
	}
}

// WithAttribute 添加自定义属性
func WithAttribute(key, value string) Option {
	return func(c *Config) {
		if c.Attributes == nil {
			c.Attributes = make(map[string]string)
		}
		c.Attributes[key] = value
	}
}

// WithInsecure 设置是否使用不安全连接
func WithInsecure(insecure bool) Option {
	return func(c *Config) {
		c.ExporterConfig.Insecure = insecure
	}
}

// WithSamplerType 设置采样器类型
func WithSamplerType(samplerType string) Option {
	return func(c *Config) {
		c.SamplerConfig.Type = samplerType
	}
}

// WithSamplerRatio 设置采样比例
func WithSamplerRatio(ratio float64) Option {
	return func(c *Config) {
		c.SamplerConfig.Ratio = ratio
	}
}

// WithParentBasedSampler 设置是否基于父级决定采样
func WithParentBasedSampler(enabled bool) Option {
	return func(c *Config) {
		c.SamplerConfig.ParentBased = enabled
	}
}
