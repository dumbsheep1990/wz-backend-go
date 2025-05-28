package entity

import (
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

const (
	EventTypeCommunityCreated   = "community.created"
	EventTypeCommunityUpdated   = "community.updated"
	EventTypeCommunityDeleted   = "community.deleted"
	EventTypeCommunityActivated = "community.activated"
)

// CommunityCreatedEvent 表示社区创建时的事件
type CommunityCreatedEvent struct {
	event.BaseDomainEvent
}

// NewCommunityCreatedEvent 创建一个新的社区创建事件
func NewCommunityCreatedEvent(community *Community) CommunityCreatedEvent {
	return CommunityCreatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCommunityCreated,
			community.ID().String(),
			"Community",
			map[string]interface{}{
				"id":          community.ID().String(),
				"name":        community.Name().String(),
				"description": community.Description().String(),
				"ownerID":     community.OwnerID().String(),
				"status":      string(community.Status()),
				"createdAt":   community.CreatedAt(),
			},
		),
	}
}

// CommunityUpdatedEvent 表示社区更新时的事件
type CommunityUpdatedEvent struct {
	event.BaseDomainEvent
}

// NewCommunityUpdatedEvent 创建一个新的社区更新事件
func NewCommunityUpdatedEvent(community *Community) CommunityUpdatedEvent {
	return CommunityUpdatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCommunityUpdated,
			community.ID().String(),
			"Community",
			map[string]interface{}{
				"id":          community.ID().String(),
				"name":        community.Name().String(),
				"description": community.Description().String(),
				"status":      string(community.Status()),
				"updatedAt":   community.UpdatedAt(),
			},
		),
	}
}

// CommunityDeletedEvent 表示社区删除时的事件
type CommunityDeletedEvent struct {
	event.BaseDomainEvent
}

// NewCommunityDeletedEvent 创建一个新的社区删除事件
func NewCommunityDeletedEvent(community *Community) CommunityDeletedEvent {
	return CommunityDeletedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCommunityDeleted,
			community.ID().String(),
			"Community",
			map[string]interface{}{
				"id":        community.ID().String(),
				"deletedAt": community.UpdatedAt(),
			},
		),
	}
}

// CommunityActivatedEvent 表示社区激活时的事件
type CommunityActivatedEvent struct {
	event.BaseDomainEvent
}

// NewCommunityActivatedEvent 创建一个新的社区激活事件
func NewCommunityActivatedEvent(community *Community) CommunityActivatedEvent {
	return CommunityActivatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCommunityActivated,
			community.ID().String(),
			"Community",
			map[string]interface{}{
				"id":          community.ID().String(),
				"activatedAt": community.UpdatedAt(),
			},
		),
	}
}
