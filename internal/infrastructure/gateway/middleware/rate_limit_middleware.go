package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"wz-backend-go/internal/domain/gateway/entity"
)

// RateLimitMiddleware 限流中间件，防止API被过度使用
type RateLimitMiddleware struct {
	redisClient *redis.Client
}

// NewRateLimitMiddleware 创建新的限流中间件
func NewRateLimitMiddleware(redisClient *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
	}
}

// Handle 处理请求限流
func (m *RateLimitMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取路由信息
		routeInterface, exists := c.Get("route")
		if !exists {
			// 如果没有路由信息，继续处理请求
			c.Next()
			return
		}

		route, ok := routeInterface.(*entity.Route)
		if !ok || route.RateLimiterID == "" {
			// 如果路由无效或没有关联限流器，继续处理请求
			c.Next()
			return
		}

		// 获取限流器
		limiterInterface, exists := c.Get("rate_limiter_" + route.RateLimiterID)
		if !exists {
			// 如果没有找到限流器，继续处理请求
			log.Printf("找不到ID为 %s 的限流器", route.RateLimiterID)
			c.Next()
			return
		}

		limiter, ok := limiterInterface.(*entity.RateLimiter)
		if !ok || !limiter.IsActive {
			// 如果限流器无效或未激活，继续处理请求
			c.Next()
			return
		}

		// 获取客户端标识符（可以是IP、API密钥、用户ID等）
		identifier := getClientIdentifier(c)

		// 执行限流检查
		shouldLimit, remaining, err := checkRateLimit(c, m.redisClient, limiter, identifier)
		if err != nil {
			log.Printf("限流检查失败: %v", err)
			// 错误情况下继续处理请求
			c.Next()
			return
		}

		// 添加限流头部
		c.Header("X-RateLimit-Limit", string(rune(limiter.MaxRequests)))
		c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		c.Header("X-RateLimit-Reset", string(rune(time.Now().Unix()+int64(limiter.WindowSeconds))))

		if shouldLimit {
			// 超过限流阈值，返回错误
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// 获取客户端标识符
func getClientIdentifier(c *gin.Context) string {
	// 首先尝试从会话中获取用户ID
	if userSession, exists := c.Get("user_session"); exists {
		if session, ok := userSession.(*entity.UserSession); ok && session.UserID != "" {
			return "user:" + session.UserID
		}
	}

	// 尝试从请求头中获取API密钥
	apiKey := c.GetHeader("X-API-Key")
	if apiKey != "" {
		return "api:" + apiKey
	}

	// 最后使用客户端IP
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = c.Request.RemoteAddr
	}
	return "ip:" + clientIP
}

// 检查限流
func checkRateLimit(c *gin.Context, client *redis.Client, limiter *entity.RateLimiter, identifier string) (bool, int, error) {
	// 构建Redis键
	key := "rate_limit:" + limiter.Type + ":" + identifier

	// 获取当前时间
	now := time.Now().Unix()
	
	// 清理过期的请求记录
	expireTime := now - int64(limiter.WindowSeconds)
	client.ZRemRangeByScore(c, key, "0", string(rune(expireTime)))

	// 获取当前窗口内的请求数量
	count, err := client.ZCard(c, key).Result()
	if err != nil {
		return false, 0, err
	}

	// 计算剩余请求数
	remaining := limiter.MaxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// 检查是否超过限制
	if count >= int64(limiter.MaxRequests) {
		return true, remaining, nil
	}

	// 添加新请求记录
	_, err = client.ZAdd(c, key, &redis.Z{
		Score:  float64(now),
		Member: now,
	}).Result()
	if err != nil {
		return false, remaining, err
	}

	// 设置过期时间
	client.Expire(c, key, time.Duration(limiter.WindowSeconds)*time.Second)

	return false, remaining, nil
}
