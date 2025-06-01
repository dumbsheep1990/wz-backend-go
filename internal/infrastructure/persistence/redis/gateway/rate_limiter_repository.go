package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

const (
	rateLimitKeyPrefix  = "gateway:ratelimit:"
	ipRateLimitPrefix   = "ip:"
	tokenRateLimitPrefix = "token:"
	serviceRateLimitPrefix = "service:"
)

// RedisRateLimiterRepository implements a distributed rate limiter using Redis
type RedisRateLimiterRepository struct {
	client redis.Client
}

// NewRedisRateLimiterRepository creates a new RedisRateLimiterRepository
func NewRedisRateLimiterRepository(client redis.Client) *RedisRateLimiterRepository {
	return &RedisRateLimiterRepository{
		client: client,
	}
}

// CheckAndIncrement checks if a request is allowed based on rate limits and increments the counter
func (r *RedisRateLimiterRepository) CheckAndIncrement(
	ctx context.Context, 
	identifier string, 
	limiterType valueobject.RateLimitType, 
	maxRequests int, 
	duration time.Duration,
) (bool, error) {
	// Create Redis key based on limit type and identifier
	var key string
	switch limiterType {
	case valueobject.RateLimitTypeIP:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, ipRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeToken:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, tokenRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeService:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, serviceRateLimitPrefix, identifier)
	default:
		return false, fmt.Errorf("unsupported rate limit type: %s", limiterType)
	}

	// Use Redis to implement the rate limiting logic with a Lua script for atomicity
	script := `
		local key = KEYS[1]
		local max = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		
		-- Get the current count
		local count = redis.call('GET', key)
		count = count and tonumber(count) or 0
		
		-- Check if the limit is exceeded
		if count >= max then
			return 0
		end
		
		-- Increment the counter
		redis.call('INCR', key)
		
		-- Set expiration if it's a new key
		if count == 0 then
			redis.call('EXPIRE', key, window)
		end
		
		return 1
	`

	// Execute the script
	result, err := r.client.Eval(ctx, script, []string{key}, maxRequests, int(duration.Seconds())).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Result is 1 if allowed, 0 if rate limited
	allowed := result.(int64) == 1
	return allowed, nil
}

// GetRemainingRequests gets the remaining requests allowed for a specific identifier
func (r *RedisRateLimiterRepository) GetRemainingRequests(
	ctx context.Context, 
	identifier string, 
	limiterType valueobject.RateLimitType, 
	maxRequests int,
) (int, error) {
	// Create Redis key based on limit type and identifier
	var key string
	switch limiterType {
	case valueobject.RateLimitTypeIP:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, ipRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeToken:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, tokenRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeService:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, serviceRateLimitPrefix, identifier)
	default:
		return 0, fmt.Errorf("unsupported rate limit type: %s", limiterType)
	}

	// Get the current count
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		// No existing count means all requests are available
		return maxRequests, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}

	// Calculate remaining requests
	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// GetExpirationTime gets the time until the rate limit window resets
func (r *RedisRateLimiterRepository) GetExpirationTime(
	ctx context.Context, 
	identifier string, 
	limiterType valueobject.RateLimitType,
) (time.Duration, error) {
	// Create Redis key based on limit type and identifier
	var key string
	switch limiterType {
	case valueobject.RateLimitTypeIP:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, ipRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeToken:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, tokenRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeService:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, serviceRateLimitPrefix, identifier)
	default:
		return 0, fmt.Errorf("unsupported rate limit type: %s", limiterType)
	}

	// Get the TTL (time to live) for the key
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit expiration: %w", err)
	}

	// If the key doesn't exist or has no expiration, return 0
	if ttl == -2 || ttl == -1 {
		return 0, nil
	}

	return ttl, nil
}

// Reset resets the rate limit counter for a specific identifier
func (r *RedisRateLimiterRepository) Reset(
	ctx context.Context, 
	identifier string, 
	limiterType valueobject.RateLimitType,
) error {
	// Create Redis key based on limit type and identifier
	var key string
	switch limiterType {
	case valueobject.RateLimitTypeIP:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, ipRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeToken:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, tokenRateLimitPrefix, identifier)
	case valueobject.RateLimitTypeService:
		key = fmt.Sprintf("%s%s%s", rateLimitKeyPrefix, serviceRateLimitPrefix, identifier)
	default:
		return fmt.Errorf("unsupported rate limit type: %s", limiterType)
	}

	// Delete the key
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}

	return nil
}
