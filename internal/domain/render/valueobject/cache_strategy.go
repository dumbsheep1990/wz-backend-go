package valueobject

import (
	"errors"
	"time"
)

// CacheStrategy 表示页面渲染的缓存策略
type CacheStrategy struct {
	enabled     bool
	ttl         time.Duration
	cacheLevel  CacheLevel
	cacheKey    string
	cacheGroups []string
}

// CacheLevel 表示缓存级别
type CacheLevel string

const (
	CacheLevelNone     CacheLevel = "none"     // 不缓存
	CacheLevelMemory   CacheLevel = "memory"   // 内存缓存
	CacheLevelRedis    CacheLevel = "redis"    // Redis缓存
	CacheLevelCDN      CacheLevel = "cdn"      // CDN缓存
	CacheLevelMultiple CacheLevel = "multiple" // 多级缓存
)

// NewCacheStrategy 创建一个新的缓存策略
func NewCacheStrategy(enabled bool, ttl time.Duration, level CacheLevel, key string, groups []string) (CacheStrategy, error) {
	if enabled && ttl <= 0 {
		return CacheStrategy{}, errors.New("缓存TTL必须大于0")
	}

	if enabled && level == "" {
		level = CacheLevelMemory // 默认使用内存缓存
	}

	if enabled && level != CacheLevelNone && key == "" {
		return CacheStrategy{}, errors.New("缓存键不能为空")
	}

	return CacheStrategy{
		enabled:     enabled,
		ttl:         ttl,
		cacheLevel:  level,
		cacheKey:    key,
		cacheGroups: groups,
	}, nil
}

// Enabled 返回缓存是否启用
func (c CacheStrategy) Enabled() bool {
	return c.enabled
}

// TTL 返回缓存的生存时间
func (c CacheStrategy) TTL() time.Duration {
	return c.ttl
}

// CacheLevel 返回缓存级别
func (c CacheStrategy) CacheLevel() CacheLevel {
	return c.cacheLevel
}

// CacheKey 返回缓存键
func (c CacheStrategy) CacheKey() string {
	return c.cacheKey
}

// CacheGroups 返回缓存组
func (c CacheStrategy) CacheGroups() []string {
	return c.cacheGroups
}

// Equals 比较两个缓存策略是否相等
func (c CacheStrategy) Equals(other CacheStrategy) bool {
	if c.enabled != other.enabled ||
		c.ttl != other.ttl ||
		c.cacheLevel != other.cacheLevel ||
		c.cacheKey != other.cacheKey {
		return false
	}

	if len(c.cacheGroups) != len(other.cacheGroups) {
		return false
	}

	for i, group := range c.cacheGroups {
		if group != other.cacheGroups[i] {
			return false
		}
	}

	return true
}
