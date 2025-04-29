package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Tracing 返回请求追踪中间件
func Tracing(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成或获取请求跟踪ID
		var traceID string
		existingTraceID := c.GetHeader("X-Trace-ID")
		if existingTraceID != "" {
			traceID = existingTraceID
		} else {
			traceID = uuid.New().String()
		}
		
		// 生成唯一的跨度ID
		spanID := uuid.New().String()
		
		// 设置跟踪信息到请求上下文
		c.Set("traceID", traceID)
		c.Set("spanID", spanID)
		c.Set("service", serviceName)
		c.Set("startTime", time.Now())
		
		// 设置追踪头信息，传递给下游服务
		c.Request.Header.Set("X-Trace-ID", traceID)
		c.Request.Header.Set("X-Span-ID", spanID)
		c.Request.Header.Set("X-Parent-Service", serviceName)
		
		// 为响应头设置追踪ID
		c.Header("X-Trace-ID", traceID)
		
		// 创建带有追踪信息的上下文
		ctx := context.WithValue(c.Request.Context(), "traceID", traceID)
		ctx = context.WithValue(ctx, "spanID", spanID)
		ctx = context.WithValue(ctx, "service", serviceName)
		
		// 使用新的上下文更新请求
		c.Request = c.Request.WithContext(ctx)
		
		// 处理请求
		c.Next()
		
		// 计算请求处理耗时
		if startTime, exists := c.Get("startTime"); exists {
			duration := time.Since(startTime.(time.Time))
			// 这里可以发送追踪信息到分布式追踪系统（例如Jaeger或Zipkin）
			// 目前只是输出日志
			fmt.Printf("[TRACE] Service: %s, TraceID: %s, SpanID: %s, Duration: %v, Status: %d\n",
				serviceName, traceID, spanID, duration, c.Writer.Status())
		}
	}
}
