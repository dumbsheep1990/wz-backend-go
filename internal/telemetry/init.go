package telemetry

import (
	"context"
	"log"
	"time"
)

// InitTracer 初始化并返回 OpenTelemetry 追踪提供者
// 这是一个便捷函数，用于在服务启动时快速初始化 OpenTelemetry
func InitTracer(serviceName, serviceVersion, environment string, exporterType string, exporterEndpoint string) (*TracerProvider, error) {
	// 创建默认配置
	config := DefaultConfig()
	
	// 设置服务信息
	config.ServiceName = serviceName
	config.ServiceVersion = serviceVersion
	config.Environment = environment
	
	// 设置导出器
	config.ExporterConfig.Type = exporterType
	config.ExporterConfig.Endpoint = exporterEndpoint
	config.ExporterConfig.Insecure = true // 开发环境通常使用不安全连接
	config.ExporterConfig.Timeout = 5 * time.Second
	
	// 设置采样器
	config.SamplerConfig.Type = "always_on" // 开发环境使用全量采样
	config.SamplerConfig.ParentBased = true // 基于父级采样
	
	// 创建并初始化追踪提供者
	tp, err := NewTracerProvider(config)
	if err != nil {
		return nil, err
	}
	
	log.Printf("[OpenTelemetry] 初始化完成: 服务=%s, 版本=%s, 环境=%s, 导出器=%s", 
		serviceName, serviceVersion, environment, exporterType)
	
	return tp, nil
}

// Shutdown 优雅关闭 OpenTelemetry 追踪提供者
func Shutdown(ctx context.Context, tp *TracerProvider) {
	if tp != nil {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("[OpenTelemetry] 关闭追踪提供者时出错: %v", err)
		} else {
			log.Println("[OpenTelemetry] 追踪提供者已关闭")
		}
	}
}
