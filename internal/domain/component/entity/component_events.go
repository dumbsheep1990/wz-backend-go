package entity

import (
	"time"
	"wz-backend-go/internal/domain/shared/event"
)

// ComponentCreatedEvent 组件创建事件
type ComponentCreatedEvent struct {
	ComponentID   string    `json:"component_id"`
	Name          string    `json:"name"`
	ComponentType string    `json:"component_type"`
	TenantID      string    `json:"tenant_id"`
	IsPublic      bool      `json:"is_public"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e ComponentCreatedEvent) GetAggregateID() string {
	return e.ComponentID
}

// GetEventType 获取事件类型
func (e ComponentCreatedEvent) GetEventType() string {
	return "component.created"
}

// GetOccurredAt 获取发生时间
func (e ComponentCreatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e ComponentCreatedEvent) GetEventData() interface{} {
	return e
}

// NewComponentCreatedEvent 创建组件创建事件
func NewComponentCreatedEvent(component *Component) event.DomainEvent {
	return ComponentCreatedEvent{
		ComponentID:   component.ID().Value(),
		Name:          component.Name(),
		ComponentType: component.ComponentType().Value(),
		TenantID:      component.TenantID(),
		IsPublic:      component.IsPublic(),
		OccurredAt:    time.Now(),
	}
}

// ComponentUpdatedEvent 组件更新事件
type ComponentUpdatedEvent struct {
	ComponentID string    `json:"component_id"`
	TenantID    string    `json:"tenant_id"`
	Field       string    `json:"field"`
	OldValue    string    `json:"old_value"`
	NewValue    string    `json:"new_value"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e ComponentUpdatedEvent) GetAggregateID() string {
	return e.ComponentID
}

// GetEventType 获取事件类型
func (e ComponentUpdatedEvent) GetEventType() string {
	return "component.updated"
}

// GetOccurredAt 获取发生时间
func (e ComponentUpdatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e ComponentUpdatedEvent) GetEventData() interface{} {
	return e
}

// NewComponentUpdatedEvent 创建组件更新事件
func NewComponentUpdatedEvent(component *Component, field, oldValue, newValue string) event.DomainEvent {
	return ComponentUpdatedEvent{
		ComponentID: component.ID().Value(),
		TenantID:    component.TenantID(),
		Field:       field,
		OldValue:    oldValue,
		NewValue:    newValue,
		OccurredAt:  time.Now(),
	}
}

// ComponentMadePublicEvent 组件设为公开事件
type ComponentMadePublicEvent struct {
	ComponentID   string    `json:"component_id"`
	Name          string    `json:"name"`
	ComponentType string    `json:"component_type"`
	TenantID      string    `json:"tenant_id"`
	Version       string    `json:"version"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e ComponentMadePublicEvent) GetAggregateID() string {
	return e.ComponentID
}

// GetEventType 获取事件类型
func (e ComponentMadePublicEvent) GetEventType() string {
	return "component.made_public"
}

// GetOccurredAt 获取发生时间
func (e ComponentMadePublicEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e ComponentMadePublicEvent) GetEventData() interface{} {
	return e
}

// NewComponentMadePublicEvent 创建组件设为公开事件
func NewComponentMadePublicEvent(component *Component) event.DomainEvent {
	return ComponentMadePublicEvent{
		ComponentID:   component.ID().Value(),
		Name:          component.Name(),
		ComponentType: component.ComponentType().Value(),
		TenantID:      component.TenantID(),
		Version:       component.Version(),
		OccurredAt:    time.Now(),
	}
}

// ComponentVersionUpdatedEvent 组件版本更新事件
type ComponentVersionUpdatedEvent struct {
	ComponentID string    `json:"component_id"`
	TenantID    string    `json:"tenant_id"`
	OldVersion  string    `json:"old_version"`
	NewVersion  string    `json:"new_version"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e ComponentVersionUpdatedEvent) GetAggregateID() string {
	return e.ComponentID
}

// GetEventType 获取事件类型
func (e ComponentVersionUpdatedEvent) GetEventType() string {
	return "component.version_updated"
}

// GetOccurredAt 获取发生时间
func (e ComponentVersionUpdatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e ComponentVersionUpdatedEvent) GetEventData() interface{} {
	return e
}

// NewComponentVersionUpdatedEvent 创建组件版本更新事件
func NewComponentVersionUpdatedEvent(component *Component, oldVersion, newVersion string) event.DomainEvent {
	return ComponentVersionUpdatedEvent{
		ComponentID: component.ID().Value(),
		TenantID:    component.TenantID(),
		OldVersion:  oldVersion,
		NewVersion:  newVersion,
		OccurredAt:  time.Now(),
	}
}

// ComponentDeletedEvent 组件删除事件
type ComponentDeletedEvent struct {
	ComponentID   string    `json:"component_id"`
	Name          string    `json:"name"`
	ComponentType string    `json:"component_type"`
	TenantID      string    `json:"tenant_id"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e ComponentDeletedEvent) GetAggregateID() string {
	return e.ComponentID
}

// GetEventType 获取事件类型
func (e ComponentDeletedEvent) GetEventType() string {
	return "component.deleted"
}

// GetOccurredAt 获取发生时间
func (e ComponentDeletedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e ComponentDeletedEvent) GetEventData() interface{} {
	return e
}

// NewComponentDeletedEvent 创建组件删除事件
func NewComponentDeletedEvent(component *Component) event.DomainEvent {
	return ComponentDeletedEvent{
		ComponentID:   component.ID().Value(),
		Name:          component.Name(),
		ComponentType: component.ComponentType().Value(),
		TenantID:      component.TenantID(),
		OccurredAt:    time.Now(),
	}
} 