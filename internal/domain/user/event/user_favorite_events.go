package event

import (
	"time"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/entity"
)

// UserFavoriteEventType 用户收藏事件类型
const (
	UserFavoriteCreatedEventType = "user.favorite.created"
	UserFavoriteDeletedEventType = "user.favorite.deleted"
)

// UserFavoriteCreatedEvent 用户收藏创建事件
type UserFavoriteCreatedEvent struct {
	event.BaseDomainEvent
	Favorite *entity.UserFavorite
}

// NewUserFavoriteCreatedEvent 创建用户收藏创建事件
func NewUserFavoriteCreatedEvent(favorite *entity.UserFavorite) *UserFavoriteCreatedEvent {
	return &UserFavoriteCreatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			UserFavoriteCreatedEventType,
			favorite.ID().String(),
			"UserFavorite",
			favorite,
		),
		Favorite: favorite,
	}
}

// OccurredTime 实现DomainEvent接口
func (e UserFavoriteCreatedEvent) OccurredTime() time.Time {
	return e.OccurredAt()
}

// UserFavoriteDeletedEvent 用户收藏删除事件
type UserFavoriteDeletedEvent struct {
	event.BaseDomainEvent
	Favorite *entity.UserFavorite
}

// NewUserFavoriteDeletedEvent 创建用户收藏删除事件
func NewUserFavoriteDeletedEvent(favorite *entity.UserFavorite) *UserFavoriteDeletedEvent {
	return &UserFavoriteDeletedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			UserFavoriteDeletedEventType,
			favorite.ID().String(),
			"UserFavorite",
			favorite,
		),
		Favorite: favorite,
	}
}

// OccurredTime 实现DomainEvent接口
func (e UserFavoriteDeletedEvent) OccurredTime() time.Time {
	return e.OccurredAt()
}
