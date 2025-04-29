package middleware

import (
	"net/http"
	"sync"
	"time"

	"wz-backend-go/internal/gateway/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 用于存储不同限流器的映射
var (
	limiters   = make(map[string]*rate.Limiter)
	limitersMu sync.RWMutex
)

// RateLimiter 返回一个请求限流中间件
func RateLimiter(config config.RateConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据配置的策略获取限流键
		var key string
		switch config.Strategy {
		case "ip":
			key = c.ClientIP()
		case "user":
			// 尝试从上下文中获取用户ID，如果不存在则使用IP
			userID, exists := c.Get("userID")
			if exists {
				key = userID.(string)
			} else {
				key = c.ClientIP()
			}
		case "tenant":
			// 尝试从上下文中获取租户ID，如果不存在则使用IP
			tenantID, exists := c.Get("tenantID")
			if exists {
				key = tenantID.(string)
			} else {
				key = c.ClientIP()
			}
		default:
			// 默认使用IP
			key = c.ClientIP()
		}

		// 获取或创建限流器
		limiter := getLimiter(key, config)

		// 尝试获取令牌
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求频率过高，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

// getLimiter 获取指定键的限流器，如果不存在则创建
func getLimiter(key string, config config.RateConfig) *rate.Limiter {
	limitersMu.RLock()
	limiter, exists := limiters[key]
	limitersMu.RUnlock()

	if !exists {
		// 计算限流速率
		r := rate.Limit(float64(config.MaxRequests) / float64(config.IntervalSecs))
		// 创建新的限流器，桶容量为最大请求数
		limiter = rate.NewLimiter(r, config.MaxRequests)

		limitersMu.Lock()
		limiters[key] = limiter
		limitersMu.Unlock()

		// 启动一个协程，定期清理不活跃的限流器
		go cleanupLimiter(key)
	}

	return limiter
}

// cleanupLimiter 定期清理不活跃的限流器
func cleanupLimiter(key string) {
	// 一小时后清理
	time.Sleep(1 * time.Hour)

	limitersMu.Lock()
	delete(limiters, key)
	limitersMu.Unlock()
}
