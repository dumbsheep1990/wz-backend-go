package telemetry

import (
	"time"
)

// Config 表示OpenTelemetry的配置
type Config struct {
	ServiceName    string            // 服务名称
	ServiceVersion string            // 服务版本
	Environment    string            // 环境(dev, staging, prod)
	Attributes     map[string]string // 服务级别的属性
	ExporterConfig ExporterConfig    // 数据导出器配置
	SamplerConfig  SamplerConfig     // 采样配置
}

// ExporterConfig 表示OpenTelemetry导出器的配置
type ExporterConfig struct {
	// 支持多种导出方式
	Type             string        // 导出器类型: "jaeger", "otlp", "zipkin"
	Endpoint         string        // 导出器端点地址
	Headers          map[string]string // 请求头（用于OTLP HTTP）
	Insecure         bool          // 是否不安全连接（用于OTLP gRPC）
	Timeout          time.Duration // 连接超时时间
	BatchTimeout     time.Duration // 批处理超时
	BatchSize        int           // 批处理大小
	MaxExportBatchSize int        // 最大导出批处理大小
	MaxQueueSize     int           // 最大队列大小
}

// SamplerConfig 表示OpenTelemetry采样器的配置
type SamplerConfig struct {
	Type        string  // 采样器类型: "always_on", "always_off", "trace_id_ratio"
	Ratio       float64 // 采样比例，仅当Type="trace_id_ratio"时使用
	ParentBased bool    // 是否基于父级决定采样
}

// DefaultConfig 返回默认的OpenTelemetry配置
func DefaultConfig(serviceName string) *Config {
	return &Config{
		ServiceName:    serviceName,
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Attributes:     make(map[string]string),
		ExporterConfig: ExporterConfig{
			Type:               "otlp",
			Endpoint:           "localhost:4317", // 默认OTLP gRPC端点
			Insecure:           true,
			Timeout:            30 * time.Second,
			BatchTimeout:       5 * time.Second,
			BatchSize:          512,
			MaxExportBatchSize: 512,
			MaxQueueSize:       2048,
		},
		SamplerConfig: SamplerConfig{
			Type:        "always_on",
			Ratio:       1.0,
			ParentBased: true,
		},
	}
}
