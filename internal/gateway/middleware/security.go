package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 返回添加安全相关HTTP头的中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置XSS保护头
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 设置内容类型选项
		c.Header("X-Content-Type-Options", "nosniff")
		
		// 设置框架选项，防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		
		// 设置内容安全策略
		// 根据实际需求可调整CSP策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
		
		// 设置引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// 设置特性策略
		c.Header("Feature-Policy", "camera 'none'; microphone 'none'; geolocation 'none'")
		
		// 严格传输安全 (HSTS)
		// 仅在HTTPS下使用，开发环境可注释掉
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		c.Next()
	}
}
