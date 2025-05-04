package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/gateway/middleware"
	"wz-backend-go/internal/telemetry"
)

// SetupMiddlewares 设置所有中间件
func SetupMiddlewares(router *gin.Engine, redisClient *redis.Client, serviceConfig config.ServiceConfig, rateConfig config.RateConfig) {
	// 注册监控中间件（最先执行，以确保所有请求都被监控）
	router.Use(middleware.Monitoring())
	
	// 租户识别与上下文中间件
	router.Use(middleware.TenantContext())
	
	// 分布式限流中间件
	router.Use(middleware.DistributedRateLimiter(redisClient, rateConfig))
	
	// 熔断中间件
	router.Use(middleware.CircuitBreaker(serviceConfig))
	
	// 其他现有中间件...
	
	// 暴露Prometheus指标端点
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// SetupMonitoring 设置监控和告警系统
func SetupMonitoring() *telemetry.AlertManager {
	// 创建告警管理器
	alertManager := telemetry.NewAlertManager()
	
	// 添加高错误率告警规则
	alertManager.AddRule(&telemetry.AlertRule{
		Name:        "HighErrorRate",
		Description: "服务错误率超过阈值",
		MetricQuery: "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m])) > 0.05",
		Config: telemetry.AlertConfig{
			Enabled:      true,
			Threshold:    0.05,  // 5%错误率
			WindowSize:   5 * time.Minute,
			CooldownTime: 15 * time.Minute,
		},
	})
	
	// 添加高延迟告警规则
	alertManager.AddRule(&telemetry.AlertRule{
		Name:        "HighLatency",
		Description: "服务延迟超过阈值",
		MetricQuery: "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 1",
		Config: telemetry.AlertConfig{
			Enabled:      true,
			Threshold:    1.0,  // 95%请求超过1秒
			WindowSize:   5 * time.Minute,
			CooldownTime: 15 * time.Minute,
		},
	})
	
	// 添加熔断器触发告警规则
	alertManager.AddRule(&telemetry.AlertRule{
		Name:        "CircuitBreakerTripped",
		Description: "服务熔断器被触发",
		MetricQuery: "increase(circuit_breaker_trips_total[5m]) > 0",
		Config: telemetry.AlertConfig{
			Enabled:      true,
			Threshold:    0,  // 任何触发都告警
			WindowSize:   5 * time.Minute,
			CooldownTime: 30 * time.Minute,
		},
	})
	
	return alertManager
}
