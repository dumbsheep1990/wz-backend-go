package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	
	"wz-backend-go/internal/telemetry"
)

// Tracing 返回请求追踪中间件
func Tracing(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取追踪上下文
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		
		// 如果没有有效的追踪上下文，生成新的追踪ID
		spanCtx := trace.SpanContextFromContext(ctx)
		var traceID string
		if !spanCtx.IsValid() {
			traceID = uuid.New().String()
		} else {
			traceID = spanCtx.TraceID().String()
		}
		
		// 开始一个新的追踪Span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		tracer := otel.Tracer("gateway")
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		
		// 添加请求属性到Span
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.host", c.Request.Host),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.path", c.FullPath()),
			attribute.String("service.name", serviceName),
		)
		
		// 添加请求头和查询参数作为追踪属性
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				span.SetAttributes(attribute.String("http.header."+key, values[0]))
			}
		}
		
		// 设置追踪ID和Span ID到上下文和响应头中
		c.Set("traceID", traceID)
		c.Set("spanID", span.SpanContext().SpanID().String())
		c.Set("service", serviceName)
		c.Set("startTime", time.Now())
		c.Header("X-Trace-ID", traceID)
		
		// 将追踪上下文注入请求中，供下游服务使用
		propagator.Inject(ctx, propagation.HeaderCarrier(c.Request.Header))
		
		// 更新请求的上下文
		c.Request = c.Request.WithContext(ctx)
		
		// 处理请求
		c.Next()
		
		// 根据响应状态码设置Span状态
		status := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", status))
		
		// 计算处理时间并记录日志
		if startTime, exists := c.Get("startTime"); exists {
			duration := time.Since(startTime.(time.Time))
			span.SetAttributes(attribute.String("http.duration", duration.String()))
			
			// 记录到控制台日志
			fmt.Printf("[TRACE] Service: %s, TraceID: %s, SpanID: %s, Duration: %v, Status: %d\n",
				serviceName, traceID, span.SpanContext().SpanID().String(), duration, status)
		}
	}
}
