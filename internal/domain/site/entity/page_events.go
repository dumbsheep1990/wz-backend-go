package entity

import (
	"time"
	"wz-backend-go/internal/domain/shared/event"
)

// PageCreatedEvent 页面创建事件
type PageCreatedEvent struct {
	PageID     string    `json:"page_id"`
	SiteID     string    `json:"site_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PageCreatedEvent) GetAggregateID() string {
	return e.PageID
}

// GetEventType 获取事件类型
func (e PageCreatedEvent) GetEventType() string {
	return "page.created"
}

// GetOccurredAt 获取发生时间
func (e PageCreatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PageCreatedEvent) GetEventData() interface{} {
	return e
}

// NewPageCreatedEvent 创建页面创建事件
func NewPageCreatedEvent(page *Page) event.DomainEvent {
	return PageCreatedEvent{
		PageID:     page.ID().Value(),
		SiteID:     page.SiteID().Value(),
		Title:      page.Title().Value(),
		Slug:       page.Slug().Value(),
		OccurredAt: time.Now(),
	}
}

// PageUpdatedEvent 页面更新事件
type PageUpdatedEvent struct {
	PageID     string    `json:"page_id"`
	SiteID     string    `json:"site_id"`
	Field      string    `json:"field"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PageUpdatedEvent) GetAggregateID() string {
	return e.PageID
}

// GetEventType 获取事件类型
func (e PageUpdatedEvent) GetEventType() string {
	return "page.updated"
}

// GetOccurredAt 获取发生时间
func (e PageUpdatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PageUpdatedEvent) GetEventData() interface{} {
	return e
}

// NewPageUpdatedEvent 创建页面更新事件
func NewPageUpdatedEvent(page *Page, field, oldValue, newValue string) event.DomainEvent {
	return PageUpdatedEvent{
		PageID:     page.ID().Value(),
		SiteID:     page.SiteID().Value(),
		Field:      field,
		OldValue:   oldValue,
		NewValue:   newValue,
		OccurredAt: time.Now(),
	}
}

// PageSetAsHomepageEvent 页面设置为首页事件
type PageSetAsHomepageEvent struct {
	PageID     string    `json:"page_id"`
	SiteID     string    `json:"site_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PageSetAsHomepageEvent) GetAggregateID() string {
	return e.PageID
}

// GetEventType 获取事件类型
func (e PageSetAsHomepageEvent) GetEventType() string {
	return "page.set_as_homepage"
}

// GetOccurredAt 获取发生时间
func (e PageSetAsHomepageEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PageSetAsHomepageEvent) GetEventData() interface{} {
	return e
}

// NewPageSetAsHomepageEvent 创建页面设置为首页事件
func NewPageSetAsHomepageEvent(page *Page) event.DomainEvent {
	return PageSetAsHomepageEvent{
		PageID:     page.ID().Value(),
		SiteID:     page.SiteID().Value(),
		Title:      page.Title().Value(),
		Slug:       page.Slug().Value(),
		OccurredAt: time.Now(),
	}
}

// PagePublishedEvent 页面发布事件
type PagePublishedEvent struct {
	PageID     string    `json:"page_id"`
	SiteID     string    `json:"site_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PagePublishedEvent) GetAggregateID() string {
	return e.PageID
}

// GetEventType 获取事件类型
func (e PagePublishedEvent) GetEventType() string {
	return "page.published"
}

// GetOccurredAt 获取发生时间
func (e PagePublishedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PagePublishedEvent) GetEventData() interface{} {
	return e
}

// NewPagePublishedEvent 创建页面发布事件
func NewPagePublishedEvent(page *Page) event.DomainEvent {
	return PagePublishedEvent{
		PageID:     page.ID().Value(),
		SiteID:     page.SiteID().Value(),
		Title:      page.Title().Value(),
		Slug:       page.Slug().Value(),
		OccurredAt: time.Now(),
	}
}

// PageUnpublishedEvent 页面取消发布事件
type PageUnpublishedEvent struct {
	PageID     string    `json:"page_id"`
	SiteID     string    `json:"site_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PageUnpublishedEvent) GetAggregateID() string {
	return e.PageID
}

// GetEventType 获取事件类型
func (e PageUnpublishedEvent) GetEventType() string {
	return "page.unpublished"
}

// GetOccurredAt 获取发生时间
func (e PageUnpublishedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PageUnpublishedEvent) GetEventData() interface{} {
	return e
}

// NewPageUnpublishedEvent 创建页面取消发布事件
func NewPageUnpublishedEvent(page *Page) event.DomainEvent {
	return PageUnpublishedEvent{
		PageID:     page.ID().Value(),
		SiteID:     page.SiteID().Value(),
		Title:      page.Title().Value(),
		Slug:       page.Slug().Value(),
		OccurredAt: time.Now(),
	}
} 