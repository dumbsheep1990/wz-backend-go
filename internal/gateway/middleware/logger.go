package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger 返回一个日志记录中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		startTime := time.Now()
		
		// 生成请求ID
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		
		// 记录请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		
		// 处理请求
		c.Next()
		
		// 请求结束时间
		endTime := time.Now()
		// 执行时间
		latency := endTime.Sub(startTime)
		
		// 状态码
		statusCode := c.Writer.Status()
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		
		// 记录访问日志
		if errorMessage != "" {
			fmt.Printf("[API-GATEWAY] %v | %3d | %13v | %15s | %s | %s | %s\n",
				endTime.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				ip,
				method,
				path,
				errorMessage,
			)
		} else {
			fmt.Printf("[API-GATEWAY] %v | %3d | %13v | %15s | %s | %s\n",
				endTime.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				ip,
				method,
				path,
			)
		}
	}
}
