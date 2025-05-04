package middleware

import (
	"fmt"
	"time"
	
	"wz-backend-go/internal/telemetry"
	
	"github.com/gin-gonic/gin"
)

// Monitoring 返回监控中间件
func Monitoring() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		
		// 获取租户ID，如果没有则用"unknown"标记
		tenantID, exists := c.Get("tenantID")
		tenantIDStr := "unknown"
		if exists {
			tenantIDStr = tenantID.(string)
		}
		
		// 执行请求处理
		c.Next()
		
		// 记录请求指标
		status := c.Writer.Status()
		latency := time.Since(start).Seconds()
		
		// 记录请求计数
		telemetry.RequestsTotal.WithLabelValues(
			method,
			path,
			fmt.Sprintf("%d", status),
			tenantIDStr,
		).Inc()
		
		// 记录请求延迟
		telemetry.RequestDuration.WithLabelValues(
			method,
			path,
			tenantIDStr,
		).Observe(latency)
		
		// 如果是限流响应，增加限流计数
		if status == 429 {
			telemetry.RateLimitedRequestsTotal.WithLabelValues(
				path,
				tenantIDStr,
			).Inc()
		}
	}
}
