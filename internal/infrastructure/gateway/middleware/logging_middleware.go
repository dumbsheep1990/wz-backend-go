package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 生成请求ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器以捕获响应体
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 从上下文中获取路由和服务信息
		routeName := "unknown"
		serviceName := "unknown"

		if routeInterface, exists := c.Get("route"); exists {
			if route, ok := routeInterface.(interface{ GetName() string }); ok {
				routeName = route.GetName()
			}
		}

		if serviceInterface, exists := c.Get("service"); exists {
			if service, ok := serviceInterface.(interface{ GetName() string }); ok {
				serviceName = service.GetName()
			}
		}

		// 获取用户信息
		userID := "anonymous"
		if userInterface, exists := c.Get("user_session"); exists {
			if user, ok := userInterface.(interface{ GetUserID() string }); ok {
				userID = user.GetUserID()
			}
		}

		// 过滤掉健康检查请求的日志
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			return
		}

		// 记录日志
		log.Printf(
			"REQUEST: [%s] %s %s | STATUS: %d | LATENCY: %v | SERVICE: %s | ROUTE: %s | USER: %s | IP: %s | USER-AGENT: %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			serviceName,
			routeName,
			userID,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		// 对于特定的情况（如错误），记录请求和响应体
		// 在生产环境中应该避免记录敏感信息
		if c.Writer.Status() >= 400 {
			// 限制请求和响应体的大小，避免日志过大
			maxBodySize := 1024 // 1KB
			
			requestBodyStr := string(requestBody)
			if len(requestBodyStr) > maxBodySize {
				requestBodyStr = requestBodyStr[:maxBodySize] + "... [truncated]"
			}
			
			responseBodyStr := w.body.String()
			if len(responseBodyStr) > maxBodySize {
				responseBodyStr = responseBodyStr[:maxBodySize] + "... [truncated]"
			}
			
			log.Printf("REQUEST BODY [%s]: %s", requestID, requestBodyStr)
			log.Printf("RESPONSE BODY [%s]: %s", requestID, responseBodyStr)
		}
	}
}
