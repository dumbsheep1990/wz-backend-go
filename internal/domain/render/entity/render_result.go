package entity

import (
	"time"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderResult 表示渲染操作的结果
type RenderResult struct {
	id            valueobject.RenderID
	content       string
	format        valueobject.RenderFormat
	cacheStrategy valueobject.CacheStrategy
	context       valueobject.TemplateContext
	createdAt     time.Time
	expiresAt     time.Time
}

// NewRenderResult 创建一个新的渲染结果
func NewRenderResult(
	id valueobject.RenderID,
	content string,
	format valueobject.RenderFormat,
	cacheStrategy valueobject.CacheStrategy,
	context valueobject.TemplateContext,
) *RenderResult {
	now := time.Now()
	var expiresAt time.Time
	if cacheStrategy.Enabled() {
		expiresAt = now.Add(cacheStrategy.TTL())
	}

	return &RenderResult{
		id:            id,
		content:       content,
		format:        format,
		cacheStrategy: cacheStrategy,
		context:       context,
		createdAt:     now,
		expiresAt:     expiresAt,
	}
}

// ID 返回渲染结果ID
func (r *RenderResult) ID() valueobject.RenderID {
	return r.id
}

// Content 返回渲染内容
func (r *RenderResult) Content() string {
	return r.content
}

// Format 返回渲染格式
func (r *RenderResult) Format() valueobject.RenderFormat {
	return r.format
}

// CacheStrategy 返回缓存策略
func (r *RenderResult) CacheStrategy() valueobject.CacheStrategy {
	return r.cacheStrategy
}

// Context 返回模板上下文
func (r *RenderResult) Context() valueobject.TemplateContext {
	return r.context
}

// CreatedAt 返回创建时间
func (r *RenderResult) CreatedAt() time.Time {
	return r.createdAt
}

// ExpiresAt 返回过期时间
func (r *RenderResult) ExpiresAt() time.Time {
	return r.expiresAt
}

// IsExpired 检查渲染结果是否已过期
func (r *RenderResult) IsExpired() bool {
	if !r.cacheStrategy.Enabled() {
		return true
	}
	return time.Now().After(r.expiresAt)
}

// WithContent 更新渲染内容
func (r *RenderResult) WithContent(content string) *RenderResult {
	r.content = content
	return r
}

// WithCacheStrategy 更新缓存策略
func (r *RenderResult) WithCacheStrategy(strategy valueobject.CacheStrategy) *RenderResult {
	r.cacheStrategy = strategy
	if strategy.Enabled() {
		r.expiresAt = r.createdAt.Add(strategy.TTL())
	}
	return r
}

// ContentType 返回内容类型
func (r *RenderResult) ContentType() string {
	return r.format.ContentType()
}
