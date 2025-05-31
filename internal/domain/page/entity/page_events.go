package entity

import (
	"time"

	"wz-backend-go/internal/domain/page/valueobject"
	"wz-backend-go/internal/domain/shared/event"
)

// PageCreatedEvent 页面创建事件
type PageCreatedEvent struct {
	event.BaseEvent
	PageID    valueobject.PageID `json:"page_id"`
	SiteID    string             `json:"site_id"`
	Name      string             `json:"name"`
	Title     string             `json:"title"`
	Slug      string             `json:"slug"`
	Layout    string             `json:"layout"`
	CreatedAt time.Time          `json:"created_at"`
}

// NewPageCreatedEvent 创建页面创建事件
func NewPageCreatedEvent(page *Page) *PageCreatedEvent {
	return &PageCreatedEvent{
		BaseEvent: event.NewBaseEvent("page.created", page.ID().Value()),
		PageID:    page.ID(),
		SiteID:    page.SiteID(),
		Name:      page.Name(),
		Title:     page.Title().GetUnescaped(),
		Slug:      page.Slug().Value(),
		Layout:    page.Layout(),
		CreatedAt: page.CreatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageCreatedEvent) EventType() string {
	return "page.created"
}

// PageUpdatedEvent 页面更新事件
type PageUpdatedEvent struct {
	event.BaseEvent
	PageID    valueobject.PageID `json:"page_id"`
	SiteID    string             `json:"site_id"`
	Field     string             `json:"field"`
	OldValue  interface{}        `json:"old_value"`
	NewValue  interface{}        `json:"new_value"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// NewPageUpdatedEvent 创建页面更新事件
func NewPageUpdatedEvent(page *Page, field string, oldValue, newValue interface{}) *PageUpdatedEvent {
	return &PageUpdatedEvent{
		BaseEvent: event.NewBaseEvent("page.updated", page.ID().Value()),
		PageID:    page.ID(),
		SiteID:    page.SiteID(),
		Field:     field,
		OldValue:  oldValue,
		NewValue:  newValue,
		UpdatedAt: page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageUpdatedEvent) EventType() string {
	return "page.updated"
}

// PageStatusChangedEvent 页面状态变更事件
type PageStatusChangedEvent struct {
	event.BaseEvent
	PageID      valueobject.PageID     `json:"page_id"`
	SiteID      string                 `json:"site_id"`
	OldStatus   valueobject.PageStatus `json:"old_status"`
	NewStatus   valueobject.PageStatus `json:"new_status"`
	Reason      string                 `json:"reason"`
	IsHomepage  bool                   `json:"is_homepage"`
	PublishedAt *time.Time             `json:"published_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NewPageStatusChangedEvent 创建页面状态变更事件
func NewPageStatusChangedEvent(page *Page, oldStatus, newStatus valueobject.PageStatus, reason string) *PageStatusChangedEvent {
	return &PageStatusChangedEvent{
		BaseEvent:   event.NewBaseEvent("page.status_changed", page.ID().Value()),
		PageID:      page.ID(),
		SiteID:      page.SiteID(),
		OldStatus:   oldStatus,
		NewStatus:   newStatus,
		Reason:      reason,
		IsHomepage:  page.IsHomepage(),
		PublishedAt: page.PublishedAt(),
		UpdatedAt:   page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageStatusChangedEvent) EventType() string {
	return "page.status_changed"
}

// PageDeletedEvent 页面删除事件
type PageDeletedEvent struct {
	event.BaseEvent
	PageID    valueobject.PageID `json:"page_id"`
	SiteID    string             `json:"site_id"`
	Name      string             `json:"name"`
	Title     string             `json:"title"`
	Slug      string             `json:"slug"`
	DeletedAt time.Time          `json:"deleted_at"`
}

// NewPageDeletedEvent 创建页面删除事件
func NewPageDeletedEvent(page *Page) *PageDeletedEvent {
	return &PageDeletedEvent{
		BaseEvent: event.NewBaseEvent("page.deleted", page.ID().Value()),
		PageID:    page.ID(),
		SiteID:    page.SiteID(),
		Name:      page.Name(),
		Title:     page.Title().GetUnescaped(),
		Slug:      page.Slug().Value(),
		DeletedAt: page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageDeletedEvent) EventType() string {
	return "page.deleted"
}

// PageHomepageSetEvent 页面设为首页事件
type PageHomepageSetEvent struct {
	event.BaseEvent
	PageID    valueobject.PageID `json:"page_id"`
	SiteID    string             `json:"site_id"`
	Name      string             `json:"name"`
	Title     string             `json:"title"`
	Slug      string             `json:"slug"`
	URL       string             `json:"url"`
	SetAt     time.Time          `json:"set_at"`
}

// NewPageHomepageSetEvent 创建页面设为首页事件
func NewPageHomepageSetEvent(page *Page) *PageHomepageSetEvent {
	return &PageHomepageSetEvent{
		BaseEvent: event.NewBaseEvent("page.homepage_set", page.ID().Value()),
		PageID:    page.ID(),
		SiteID:    page.SiteID(),
		Name:      page.Name(),
		Title:     page.Title().GetUnescaped(),
		Slug:      page.Slug().Value(),
		URL:       page.GetURL(),
		SetAt:     page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageHomepageSetEvent) EventType() string {
	return "page.homepage_set"
}

// PageHomepageUnsetEvent 页面取消首页事件
type PageHomepageUnsetEvent struct {
	event.BaseEvent
	PageID   valueobject.PageID `json:"page_id"`
	SiteID   string             `json:"site_id"`
	Name     string             `json:"name"`
	Title    string             `json:"title"`
	Slug     string             `json:"slug"`
	UnsetAt  time.Time          `json:"unset_at"`
}

// NewPageHomepageUnsetEvent 创建页面取消首页事件
func NewPageHomepageUnsetEvent(page *Page) *PageHomepageUnsetEvent {
	return &PageHomepageUnsetEvent{
		BaseEvent: event.NewBaseEvent("page.homepage_unset", page.ID().Value()),
		PageID:    page.ID(),
		SiteID:    page.SiteID(),
		Name:      page.Name(),
		Title:     page.Title().GetUnescaped(),
		Slug:      page.Slug().Value(),
		UnsetAt:   page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageHomepageUnsetEvent) EventType() string {
	return "page.homepage_unset"
}

// PageSortOrderChangedEvent 页面排序变更事件
type PageSortOrderChangedEvent struct {
	event.BaseEvent
	PageID       valueobject.PageID `json:"page_id"`
	SiteID       string             `json:"site_id"`
	OldSortOrder int32              `json:"old_sort_order"`
	NewSortOrder int32              `json:"new_sort_order"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// NewPageSortOrderChangedEvent 创建页面排序变更事件
func NewPageSortOrderChangedEvent(page *Page, oldSortOrder, newSortOrder int32) *PageSortOrderChangedEvent {
	return &PageSortOrderChangedEvent{
		BaseEvent:    event.NewBaseEvent("page.sort_order_changed", page.ID().Value()),
		PageID:       page.ID(),
		SiteID:       page.SiteID(),
		OldSortOrder: oldSortOrder,
		NewSortOrder: newSortOrder,
		UpdatedAt:    page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageSortOrderChangedEvent) EventType() string {
	return "page.sort_order_changed"
}

// PagePublishedEvent 页面发布事件
type PagePublishedEvent struct {
	event.BaseEvent
	PageID      valueobject.PageID `json:"page_id"`
	SiteID      string             `json:"site_id"`
	Name        string             `json:"name"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	URL         string             `json:"url"`
	SEOScore    int                `json:"seo_score"`
	PublishedAt time.Time          `json:"published_at"`
}

// NewPagePublishedEvent 创建页面发布事件
func NewPagePublishedEvent(page *Page) *PagePublishedEvent {
	return &PagePublishedEvent{
		BaseEvent:   event.NewBaseEvent("page.published", page.ID().Value()),
		PageID:      page.ID(),
		SiteID:      page.SiteID(),
		Name:        page.Name(),
		Title:       page.Title().GetUnescaped(),
		Slug:        page.Slug().Value(),
		URL:         page.GetURL(),
		SEOScore:    page.GetSEOScore(),
		PublishedAt: *page.PublishedAt(),
	}
}

// EventType 返回事件类型
func (e *PagePublishedEvent) EventType() string {
	return "page.published"
}

// PageUnpublishedEvent 页面取消发布事件
type PageUnpublishedEvent struct {
	event.BaseEvent
	PageID        valueobject.PageID `json:"page_id"`
	SiteID        string             `json:"site_id"`
	Name          string             `json:"name"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	UnpublishedAt time.Time          `json:"unpublished_at"`
}

// NewPageUnpublishedEvent 创建页面取消发布事件
func NewPageUnpublishedEvent(page *Page) *PageUnpublishedEvent {
	return &PageUnpublishedEvent{
		BaseEvent:     event.NewBaseEvent("page.unpublished", page.ID().Value()),
		PageID:        page.ID(),
		SiteID:        page.SiteID(),
		Name:          page.Name(),
		Title:         page.Title().GetUnescaped(),
		Slug:          page.Slug().Value(),
		UnpublishedAt: page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageUnpublishedEvent) EventType() string {
	return "page.unpublished"
}

// PageArchivedEvent 页面归档事件
type PageArchivedEvent struct {
	event.BaseEvent
	PageID     valueobject.PageID `json:"page_id"`
	SiteID     string             `json:"site_id"`
	Name       string             `json:"name"`
	Title      string             `json:"title"`
	Slug       string             `json:"slug"`
	ArchivedAt time.Time          `json:"archived_at"`
}

// NewPageArchivedEvent 创建页面归档事件
func NewPageArchivedEvent(page *Page) *PageArchivedEvent {
	return &PageArchivedEvent{
		BaseEvent:  event.NewBaseEvent("page.archived", page.ID().Value()),
		PageID:     page.ID(),
		SiteID:     page.SiteID(),
		Name:       page.Name(),
		Title:      page.Title().GetUnescaped(),
		Slug:       page.Slug().Value(),
		ArchivedAt: page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageArchivedEvent) EventType() string {
	return "page.archived"
}

// PageRestoredEvent 页面恢复事件
type PageRestoredEvent struct {
	event.BaseEvent
	PageID     valueobject.PageID `json:"page_id"`
	SiteID     string             `json:"site_id"`
	Name       string             `json:"name"`
	Title      string             `json:"title"`
	Slug       string             `json:"slug"`
	RestoredAt time.Time          `json:"restored_at"`
}

// NewPageRestoredEvent 创建页面恢复事件
func NewPageRestoredEvent(page *Page) *PageRestoredEvent {
	return &PageRestoredEvent{
		BaseEvent:  event.NewBaseEvent("page.restored", page.ID().Value()),
		PageID:     page.ID(),
		SiteID:     page.SiteID(),
		Name:       page.Name(),
		Title:      page.Title().GetUnescaped(),
		Slug:       page.Slug().Value(),
		RestoredAt: page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageRestoredEvent) EventType() string {
	return "page.restored"
}

// PageSEOUpdatedEvent 页面SEO更新事件
type PageSEOUpdatedEvent struct {
	event.BaseEvent
	PageID      valueobject.PageID `json:"page_id"`
	SiteID      string             `json:"site_id"`
	OldSEOScore int                `json:"old_seo_score"`
	NewSEOScore int                `json:"new_seo_score"`
	Description string             `json:"description"`
	Keywords    []string           `json:"keywords"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// NewPageSEOUpdatedEvent 创建页面SEO更新事件
func NewPageSEOUpdatedEvent(page *Page, oldSEOScore int) *PageSEOUpdatedEvent {
	return &PageSEOUpdatedEvent{
		BaseEvent:   event.NewBaseEvent("page.seo_updated", page.ID().Value()),
		PageID:      page.ID(),
		SiteID:      page.SiteID(),
		OldSEOScore: oldSEOScore,
		NewSEOScore: page.GetSEOScore(),
		Description: page.SEOMeta().GetUnescapedDescription(),
		Keywords:    page.SEOMeta().Keywords(),
		UpdatedAt:   page.UpdatedAt(),
	}
}

// EventType 返回事件类型
func (e *PageSEOUpdatedEvent) EventType() string {
	return "page.seo_updated"
} 