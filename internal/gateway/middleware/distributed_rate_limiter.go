package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"wz-backend-go/internal/gateway/config"
	"wz-backend-go/internal/telemetry"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// DistributedRateLimiter 返回一个基于Redis的分布式限流中间件
func DistributedRateLimiter(redisClient *redis.Client, config config.RateConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据配置的策略获取限流键
		key := getLimiterKey(c, config.Strategy)
		
		// 在Redis中实现令牌桶算法
		allowed, err := checkLimit(c.Request.Context(), redisClient, key, config)
		if err != nil {
			// Redis错误时默认允许请求，但记录错误日志
			log.Printf("限流器错误: %v", err)
			c.Next()
			return
		}
		
		if !allowed {
			// 记录被限流的请求
			path := c.Request.URL.Path
			tenantID, exists := c.Get("tenantID")
			tenantIDStr := "unknown"
			if exists {
				tenantIDStr = tenantID.(string)
			}
			
			telemetry.RateLimitedRequestsTotal.WithLabelValues(
				path,
				tenantIDStr,
			).Inc()
			
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求频率过高，请稍后再试",
			})
			return
		}
		
		c.Next()
	}
}

// 获取限流器的键
func getLimiterKey(c *gin.Context, strategy string) string {
	// 构建限流键前缀
	prefix := "rate_limit:"
	
	var key string
	switch strategy {
	case "ip":
		key = c.ClientIP()
	case "user":
		userID, exists := c.Get("userID")
		if exists {
			key = userID.(string)
		} else {
			key = c.ClientIP()
		}
	case "tenant":
		tenantID, exists := c.Get("tenantID")
		if exists {
			key = tenantID.(string)
		} else {
			key = c.ClientIP()
		}
	case "path":
		// 按API路径限流
		key = c.Request.URL.Path
	default:
		key = c.ClientIP()
	}
	
	return prefix + strategy + ":" + key
}

// 检查是否允许请求
func checkLimit(ctx context.Context, client *redis.Client, key string, config config.RateConfig) (bool, error) {
	// 使用Redis的Lua脚本实现令牌桶算法
	script := `
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])
	
	-- 获取当前计数和上次刷新时间
	local currentTokens = tonumber(redis.call('get', key) or limit)
	local lastRefresh = tonumber(redis.call('get', key .. ':ts') or 0)
	
	-- 计算应该恢复的令牌数
	local elapsed = math.max(0, now - lastRefresh)
	local tokensToAdd = math.floor(elapsed / (1000 / (limit / window)))
	local newTokens = math.min(limit, currentTokens + tokensToAdd)
	
	-- 获取令牌
	local allowed = newTokens >= 1
	if allowed then
		newTokens = newTokens - 1
	end
	
	-- 更新令牌数和时间戳
	redis.call('setex', key, window, newTokens)
	redis.call('setex', key .. ':ts', window, now)
	
	return allowed and 1 or 0
	`
	
	// 获取当前时间戳（毫秒）
	now := time.Now().UnixNano() / int64(time.Millisecond)
	
	// 执行Lua脚本
	result, err := client.Eval(ctx, script, []string{key}, 
		now, 
		config.IntervalSecs,
		config.MaxRequests,
	).Int()
	
	if err != nil {
		return false, err
	}
	
	return result == 1, nil
}
