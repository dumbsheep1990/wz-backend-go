package entity

import (
	"time"
	"wz-backend-go/internal/domain/shared/event"
)

// SiteCreatedEvent 站点创建事件
type SiteCreatedEvent struct {
	SiteID      string    `json:"site_id"`
	SiteName    string    `json:"site_name"`
	TenantID    string    `json:"tenant_id"`
	Domain      string    `json:"domain"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e SiteCreatedEvent) GetAggregateID() string {
	return e.SiteID
}

// GetEventType 获取事件类型
func (e SiteCreatedEvent) GetEventType() string {
	return "site.created"
}

// GetOccurredAt 获取发生时间
func (e SiteCreatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e SiteCreatedEvent) GetEventData() interface{} {
	return e
}

// NewSiteCreatedEvent 创建站点创建事件
func NewSiteCreatedEvent(site *Site) event.DomainEvent {
	return SiteCreatedEvent{
		SiteID:     site.ID().Value(),
		SiteName:   site.Name().Value(),
		TenantID:   site.TenantID(),
		Domain:     site.Domain().Value(),
		OccurredAt: time.Now(),
	}
}

// SiteUpdatedEvent 站点更新事件
type SiteUpdatedEvent struct {
	SiteID     string    `json:"site_id"`
	TenantID   string    `json:"tenant_id"`
	Field      string    `json:"field"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e SiteUpdatedEvent) GetAggregateID() string {
	return e.SiteID
}

// GetEventType 获取事件类型
func (e SiteUpdatedEvent) GetEventType() string {
	return "site.updated"
}

// GetOccurredAt 获取发生时间
func (e SiteUpdatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e SiteUpdatedEvent) GetEventData() interface{} {
	return e
}

// NewSiteUpdatedEvent 创建站点更新事件
func NewSiteUpdatedEvent(site *Site, field, oldValue, newValue string) event.DomainEvent {
	return SiteUpdatedEvent{
		SiteID:     site.ID().Value(),
		TenantID:   site.TenantID(),
		Field:      field,
		OldValue:   oldValue,
		NewValue:   newValue,
		OccurredAt: time.Now(),
	}
}

// SitePublishedEvent 站点发布事件
type SitePublishedEvent struct {
	SiteID       string    `json:"site_id"`
	SiteName     string    `json:"site_name"`
	TenantID     string    `json:"tenant_id"`
	Domain       string    `json:"domain"`
	PublishedAt  time.Time `json:"published_at"`
	OccurredAt   time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e SitePublishedEvent) GetAggregateID() string {
	return e.SiteID
}

// GetEventType 获取事件类型
func (e SitePublishedEvent) GetEventType() string {
	return "site.published"
}

// GetOccurredAt 获取发生时间
func (e SitePublishedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e SitePublishedEvent) GetEventData() interface{} {
	return e
}

// NewSitePublishedEvent 创建站点发布事件
func NewSitePublishedEvent(site *Site) event.DomainEvent {
	return SitePublishedEvent{
		SiteID:      site.ID().Value(),
		SiteName:    site.Name().Value(),
		TenantID:    site.TenantID(),
		Domain:      site.Domain().Value(),
		PublishedAt: *site.PublishedAt(),
		OccurredAt:  time.Now(),
	}
}

// SiteArchivedEvent 站点归档事件
type SiteArchivedEvent struct {
	SiteID     string    `json:"site_id"`
	SiteName   string    `json:"site_name"`
	TenantID   string    `json:"tenant_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e SiteArchivedEvent) GetAggregateID() string {
	return e.SiteID
}

// GetEventType 获取事件类型
func (e SiteArchivedEvent) GetEventType() string {
	return "site.archived"
}

// GetOccurredAt 获取发生时间
func (e SiteArchivedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e SiteArchivedEvent) GetEventData() interface{} {
	return e
}

// NewSiteArchivedEvent 创建站点归档事件
func NewSiteArchivedEvent(site *Site) event.DomainEvent {
	return SiteArchivedEvent{
		SiteID:     site.ID().Value(),
		SiteName:   site.Name().Value(),
		TenantID:   site.TenantID(),
		OccurredAt: time.Now(),
	}
} 