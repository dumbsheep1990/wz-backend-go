package entity

import (
	"time"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderEvent 渲染事件接口
type RenderEvent interface {
	EventType() string
	OccurredAt() time.Time
	EntityID() string
}

// BaseRenderEvent 渲染事件基类
type BaseRenderEvent struct {
	eventType  string
	occurredAt time.Time
	entityID   string
}

// EventType 返回事件类型
func (e BaseRenderEvent) EventType() string {
	return e.eventType
}

// OccurredAt 返回事件发生时间
func (e BaseRenderEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EntityID 返回实体ID
func (e BaseRenderEvent) EntityID() string {
	return e.entityID
}

// 渲染事件类型常量
const (
	EventTypeRenderCompleted    = "render.completed"
	EventTypeRenderFailed       = "render.failed"
	EventTypeCacheHit           = "render.cache.hit"
	EventTypeCacheMiss          = "render.cache.miss"
	EventTypeCacheExpired       = "render.cache.expired"
	EventTypeTemplateUpdated    = "render.template.updated"
	EventTypeTemplateCreated    = "render.template.created"
	EventTypeTemplateDeleted    = "render.template.deleted"
)

// RenderCompletedEvent 渲染完成事件
type RenderCompletedEvent struct {
	BaseRenderEvent
	renderID        valueobject.RenderID
	format          valueobject.RenderFormat
	cacheEnabled    bool
	executionTimeMs int64
	contentSize     int
}

// NewRenderCompletedEvent 创建一个新的渲染完成事件
func NewRenderCompletedEvent(
	entityID string,
	renderID valueobject.RenderID,
	format valueobject.RenderFormat,
	cacheEnabled bool,
	executionTimeMs int64,
	contentSize int,
) *RenderCompletedEvent {
	return &RenderCompletedEvent{
		BaseRenderEvent: BaseRenderEvent{
			eventType:  EventTypeRenderCompleted,
			occurredAt: time.Now(),
			entityID:   entityID,
		},
		renderID:        renderID,
		format:          format,
		cacheEnabled:    cacheEnabled,
		executionTimeMs: executionTimeMs,
		contentSize:     contentSize,
	}
}

// RenderID 返回渲染ID
func (e *RenderCompletedEvent) RenderID() valueobject.RenderID {
	return e.renderID
}

// Format 返回渲染格式
func (e *RenderCompletedEvent) Format() valueobject.RenderFormat {
	return e.format
}

// CacheEnabled 返回缓存是否启用
func (e *RenderCompletedEvent) CacheEnabled() bool {
	return e.cacheEnabled
}

// ExecutionTimeMs 返回执行时间（毫秒）
func (e *RenderCompletedEvent) ExecutionTimeMs() int64 {
	return e.executionTimeMs
}

// ContentSize 返回内容大小
func (e *RenderCompletedEvent) ContentSize() int {
	return e.contentSize
}

// RenderFailedEvent 渲染失败事件
type RenderFailedEvent struct {
	BaseRenderEvent
	error          string
	templateID     string
	templateType   string
	executionTimeMs int64
}

// NewRenderFailedEvent 创建一个新的渲染失败事件
func NewRenderFailedEvent(
	entityID string,
	error string,
	templateID string,
	templateType string,
	executionTimeMs int64,
) *RenderFailedEvent {
	return &RenderFailedEvent{
		BaseRenderEvent: BaseRenderEvent{
			eventType:  EventTypeRenderFailed,
			occurredAt: time.Now(),
			entityID:   entityID,
		},
		error:          error,
		templateID:     templateID,
		templateType:   templateType,
		executionTimeMs: executionTimeMs,
	}
}

// Error 返回错误信息
func (e *RenderFailedEvent) Error() string {
	return e.error
}

// TemplateID 返回模板ID
func (e *RenderFailedEvent) TemplateID() string {
	return e.templateID
}

// TemplateType 返回模板类型
func (e *RenderFailedEvent) TemplateType() string {
	return e.templateType
}

// ExecutionTimeMs 返回执行时间（毫秒）
func (e *RenderFailedEvent) ExecutionTimeMs() int64 {
	return e.executionTimeMs
}

// CacheHitEvent 缓存命中事件
type CacheHitEvent struct {
	BaseRenderEvent
	cacheKey      string
	cacheLevel    valueobject.CacheLevel
	responseTimeMs int64
}

// NewCacheHitEvent 创建一个新的缓存命中事件
func NewCacheHitEvent(
	entityID string,
	cacheKey string,
	cacheLevel valueobject.CacheLevel,
	responseTimeMs int64,
) *CacheHitEvent {
	return &CacheHitEvent{
		BaseRenderEvent: BaseRenderEvent{
			eventType:  EventTypeCacheHit,
			occurredAt: time.Now(),
			entityID:   entityID,
		},
		cacheKey:      cacheKey,
		cacheLevel:    cacheLevel,
		responseTimeMs: responseTimeMs,
	}
}

// TemplateUpdatedEvent 模板更新事件
type TemplateUpdatedEvent struct {
	BaseRenderEvent
	templateID    string
	templateName  string
	templateType  string
	siteID        string
}

// NewTemplateUpdatedEvent 创建一个新的模板更新事件
func NewTemplateUpdatedEvent(
	templateID string,
	templateName string,
	templateType string,
	siteID string,
) *TemplateUpdatedEvent {
	return &TemplateUpdatedEvent{
		BaseRenderEvent: BaseRenderEvent{
			eventType:  EventTypeTemplateUpdated,
			occurredAt: time.Now(),
			entityID:   templateID,
		},
		templateID:    templateID,
		templateName:  templateName,
		templateType:  templateType,
		siteID:        siteID,
	}
}

// TemplateID 返回模板ID
func (e *TemplateUpdatedEvent) TemplateID() string {
	return e.templateID
}

// TemplateName 返回模板名称
func (e *TemplateUpdatedEvent) TemplateName() string {
	return e.templateName
}

// TemplateType 返回模板类型
func (e *TemplateUpdatedEvent) TemplateType() string {
	return e.templateType
}

// SiteID 返回站点ID
func (e *TemplateUpdatedEvent) SiteID() string {
	return e.siteID
}
