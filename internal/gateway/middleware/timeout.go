package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 返回请求超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 更新请求的上下文
		c.Request = c.Request.WithContext(ctx)

		// 创建一个通道，用于通知请求处理完成
		done := make(chan struct{})
		
		// 使用goroutine处理请求，以便可以捕获超时
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		// 等待请求完成或超时
		select {
		case <-done:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 请求超时
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"code":    504,
					"message": "请求处理超时",
				})
			}
		}
	}
}
