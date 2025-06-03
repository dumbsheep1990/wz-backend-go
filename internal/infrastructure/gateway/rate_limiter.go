package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"wz-backend-go/internal/domain/gateway/entity"
)

// RateLimiterImpl 基于Redis的限流器实现
type RateLimiterImpl struct {
	redisClient *redis.Client
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(redisClient *redis.Client) *RateLimiterImpl {
	return &RateLimiterImpl{
		redisClient: redisClient,
	}
}

// ShouldLimit 检查请求是否应该被限流
func (r *RateLimiterImpl) ShouldLimit(ctx context.Context, limiter *entity.RateLimiter, identifier string) (bool, error) {
	// 如果没有提供标识符，使用随机UUID（基于IP或API密钥等唯一标识符）
	if identifier == "" {
		identifier = uuid.New().String()
	}
	
	// 使用Lua脚本实现原子计数器操作
	// 这是一个简单的滑动窗口限流算法
	script := `
	local key = KEYS[1]
	local window = tonumber(ARGV[1])
	local limit = tonumber(ARGV[2])
	local current_time = tonumber(ARGV[3])
	
	-- 移除过期的请求
	redis.call('ZREMRANGEBYSCORE', key, 0, current_time - window)
	
	-- 获取窗口内请求数
	local count = redis.call('ZCARD', key)
	
	-- 如果未达到限制，添加新请求
	if count < limit then
		redis.call('ZADD', key, current_time, current_time .. '-' .. math.random())
		redis.call('EXPIRE', key, window)
		return 0
	end
	
	-- 已达到限制
	return 1
	`
	
	// 构建Redis键
	key := fmt.Sprintf("rate_limit:%s:%s", limiter.Type, identifier)
	
	// 执行Lua脚本
	now := time.Now().Unix()
	window := int64(limiter.WindowSeconds)
	limit := int64(limiter.MaxRequests)
	
	result, err := r.redisClient.Eval(ctx, script, []string{key}, window, limit, now).Result()
	if err != nil {
		return true, fmt.Errorf("执行限流脚本失败: %w", err)
	}
	
	// 返回是否应该限流
	shouldLimit := result.(int64) == 1
	return shouldLimit, nil
}

// ResetLimiter 重置特定标识符的限流计数器
func (r *RateLimiterImpl) ResetLimiter(ctx context.Context, limiterType string, identifier string) error {
	key := fmt.Sprintf("rate_limit:%s:%s", limiterType, identifier)
	_, err := r.redisClient.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("重置限流计数器失败: %w", err)
	}
	return nil
}

// GetRemainingRequests 获取剩余请求数
func (r *RateLimiterImpl) GetRemainingRequests(ctx context.Context, limiter *entity.RateLimiter, identifier string) (int, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", limiter.Type, identifier)
	
	// 移除过期的请求
	now := time.Now().Unix()
	windowStart := now - int64(limiter.WindowSeconds)
	_, err := r.redisClient.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart)).Result()
	if err != nil {
		return 0, fmt.Errorf("移除过期请求失败: %w", err)
	}
	
	// 获取当前计数
	count, err := r.redisClient.ZCard(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取请求计数失败: %w", err)
	}
	
	// 计算剩余请求数
	remaining := limiter.MaxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}
	
	return remaining, nil
}
