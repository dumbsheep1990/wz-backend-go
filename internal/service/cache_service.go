package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService Redis缓存服务
type CacheService interface {
	// 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// 获取缓存
	Get(ctx context.Context, key string, value interface{}) error
	// 删除缓存
	Delete(ctx context.Context, key string) error
	// 批量删除缓存（通过前缀）
	DeleteByPrefix(ctx context.Context, prefix string) error
	// 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// 设置缓存（仅当缓存不存在时）
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	// 使用缓存（优先使用缓存，缓存不存在时调用源数据获取函数并更新缓存）
	UseCache(ctx context.Context, key string, expiration time.Duration, value interface{}, getSourceFn func() (interface{}, error)) error
}

type cacheService struct {
	redis        *redis.Client
	localCache   map[string]interface{} // 本地内存缓存
	localExpires map[string]time.Time   // 本地缓存过期时间
}

// NewCacheService 创建缓存服务
func NewCacheService(redis *redis.Client) CacheService {
	cache := &cacheService{
		redis:        redis,
		localCache:   make(map[string]interface{}),
		localExpires: make(map[string]time.Time),
	}

	// 定期清理过期的本地缓存
	go cache.cleanExpiredLocalCache()

	return cache
}

// Set 设置缓存
func (s *cacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 序列化值
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Redis缓存
	err = s.redis.Set(ctx, key, jsonBytes, expiration).Err()
	if err != nil {
		return err
	}

	// 本地缓存（如果过期时间大于0）
	if expiration > 0 {
		s.localCache[key] = value
		s.localExpires[key] = time.Now().Add(expiration)
	}

	return nil
}

// Get 获取缓存
func (s *cacheService) Get(ctx context.Context, key string, value interface{}) error {
	// 先尝试从本地缓存获取
	if localValue, ok := s.localCache[key]; ok {
		if expire, exists := s.localExpires[key]; exists && time.Now().Before(expire) {
			// 本地缓存未过期，直接使用
			bytes, err := json.Marshal(localValue)
			if err == nil {
				return json.Unmarshal(bytes, value)
			}
		}
	}

	// 本地缓存不存在或已过期，从Redis获取
	jsonBytes, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("cache miss")
		}
		return err
	}

	// 解析数据
	err = json.Unmarshal(jsonBytes, value)
	if err != nil {
		return err
	}

	// 更新本地缓存
	ttl, err := s.redis.TTL(ctx, key).Result()
	if err == nil && ttl > 0 {
		s.localCache[key] = value
		s.localExpires[key] = time.Now().Add(ttl)
	}

	return nil
}

// Delete 删除缓存
func (s *cacheService) Delete(ctx context.Context, key string) error {
	// 删除Redis缓存
	err := s.redis.Del(ctx, key).Err()

	// 删除本地缓存
	delete(s.localCache, key)
	delete(s.localExpires, key)

	return err
}

// DeleteByPrefix 批量删除缓存（通过前缀）
func (s *cacheService) DeleteByPrefix(ctx context.Context, prefix string) error {
	// 使用SCAN命令查找所有符合前缀的key
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = s.redis.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	// 批量删除Redis中的keys
	if len(keys) > 0 {
		err := s.redis.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}
	}

	// 删除本地缓存中符合前缀的键
	for k := range s.localCache {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(s.localCache, k)
			delete(s.localExpires, k)
		}
	}

	return nil
}

// Exists 检查缓存是否存在
func (s *cacheService) Exists(ctx context.Context, key string) (bool, error) {
	// 先检查本地缓存
	if val, ok := s.localCache[key]; ok {
		if expire, exists := s.localExpires[key]; exists && time.Now().Before(expire) {
			return val != nil, nil
		}
	}

	// 检查Redis缓存
	exists, err := s.redis.Exists(ctx, key).Result()
	return exists > 0, err
}

// SetNX 设置缓存（仅当缓存不存在时）
func (s *cacheService) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// 序列化值
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	// 在Redis中设置
	result, err := s.redis.SetNX(ctx, key, jsonBytes, expiration).Result()
	if err != nil {
		return false, err
	}

	// 如果成功设置，更新本地缓存
	if result && expiration > 0 {
		s.localCache[key] = value
		s.localExpires[key] = time.Now().Add(expiration)
	}

	return result, nil
}

// UseCache 使用缓存（优先使用缓存，缓存不存在时调用源数据获取函数并更新缓存）
func (s *cacheService) UseCache(ctx context.Context, key string, expiration time.Duration, value interface{}, getSourceFn func() (interface{}, error)) error {
	// 尝试从缓存获取
	err := s.Get(ctx, key, value)
	if err == nil {
		// 缓存命中
		return nil
	}

	// 缓存未命中，获取源数据
	sourceData, err := getSourceFn()
	if err != nil {
		return err
	}

	// 更新值
	sourceDataBytes, err := json.Marshal(sourceData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(sourceDataBytes, value)
	if err != nil {
		return err
	}

	// 存入缓存
	return s.Set(ctx, key, sourceData, expiration)
}

// cleanExpiredLocalCache 清理过期的本地缓存
func (s *cacheService) cleanExpiredLocalCache() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for k, expire := range s.localExpires {
			if now.After(expire) {
				delete(s.localCache, k)
				delete(s.localExpires, k)
			}
		}
	}
}
