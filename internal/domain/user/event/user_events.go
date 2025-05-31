package event

import (
	"time"

	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	*event.BaseEvent
	UserID   valueobject.UserID
	Username string
	Email    string
	Phone    string
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(user *entity.User) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: event.NewBaseEvent("user.registered", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		Email:     user.Email().Value(),
		Phone:     user.Phone().Value(),
	}
}

// UserLoginEvent 用户登录事件
type UserLoginEvent struct {
	*event.BaseEvent
	UserID    valueobject.UserID
	Username  string
	LoginTime time.Time
	ClientIP  string
	UserAgent string
}

// NewUserLoginEvent 创建用户登录事件
func NewUserLoginEvent(user *entity.User, clientIP, userAgent string) *UserLoginEvent {
	return &UserLoginEvent{
		BaseEvent: event.NewBaseEvent("user.login", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		LoginTime: time.Now(),
		ClientIP:  clientIP,
		UserAgent: userAgent,
	}
}

// UserLoginFailedEvent 用户登录失败事件
type UserLoginFailedEvent struct {
	*event.BaseEvent
	UserID       valueobject.UserID
	Username     string
	FailReason   string
	FailCount    int32
	ClientIP     string
	UserAgent    string
	IsLocked     bool
}

// NewUserLoginFailedEvent 创建用户登录失败事件
func NewUserLoginFailedEvent(user *entity.User, reason, clientIP, userAgent string) *UserLoginFailedEvent {
	return &UserLoginFailedEvent{
		BaseEvent:  event.NewBaseEvent("user.login.failed", user.ID().Value(), time.Now()),
		UserID:     user.ID(),
		Username:   user.Username().Value(),
		FailReason: reason,
		FailCount:  user.LoginFailCount(),
		ClientIP:   clientIP,
		UserAgent:  userAgent,
		IsLocked:   user.Status().IsLocked(),
	}
}

// UserStatusChangedEvent 用户状态变更事件
type UserStatusChangedEvent struct {
	*event.BaseEvent
	UserID    valueobject.UserID
	Username  string
	OldStatus valueobject.UserStatus
	NewStatus valueobject.UserStatus
	Reason    string
	Operator  string
}

// NewUserStatusChangedEvent 创建用户状态变更事件
func NewUserStatusChangedEvent(user *entity.User, oldStatus valueobject.UserStatus, reason, operator string) *UserStatusChangedEvent {
	return &UserStatusChangedEvent{
		BaseEvent: event.NewBaseEvent("user.status.changed", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		OldStatus: oldStatus,
		NewStatus: user.Status(),
		Reason:    reason,
		Operator:  operator,
	}
}

// UserActivatedEvent 用户激活事件
type UserActivatedEvent struct {
	*event.BaseEvent
	UserID   valueobject.UserID
	Username string
	Email    string
}

// NewUserActivatedEvent 创建用户激活事件
func NewUserActivatedEvent(user *entity.User) *UserActivatedEvent {
	return &UserActivatedEvent{
		BaseEvent: event.NewBaseEvent("user.activated", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		Email:     user.Email().Value(),
	}
}

// UserVerifiedEvent 用户验证事件
type UserVerifiedEvent struct {
	*event.BaseEvent
	UserID         valueobject.UserID
	Username       string
	Email          string
	VerifiedType   string // "email" 或 "company"
}

// NewUserVerifiedEvent 创建用户验证事件
func NewUserVerifiedEvent(user *entity.User, verifiedType string) *UserVerifiedEvent {
	return &UserVerifiedEvent{
		BaseEvent:    event.NewBaseEvent("user.verified", user.ID().Value(), time.Now()),
		UserID:       user.ID(),
		Username:     user.Username().Value(),
		Email:        user.Email().Value(),
		VerifiedType: verifiedType,
	}
}

// UserSuspendedEvent 用户暂停事件
type UserSuspendedEvent struct {
	*event.BaseEvent
	UserID   valueobject.UserID
	Username string
	Reason   string
	Operator string
	Duration *time.Duration // 暂停时长，nil表示无限期
}

// NewUserSuspendedEvent 创建用户暂停事件
func NewUserSuspendedEvent(user *entity.User, reason, operator string, duration *time.Duration) *UserSuspendedEvent {
	return &UserSuspendedEvent{
		BaseEvent: event.NewBaseEvent("user.suspended", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		Reason:    reason,
		Operator:  operator,
		Duration:  duration,
	}
}

// UserBannedEvent 用户封禁事件
type UserBannedEvent struct {
	*event.BaseEvent
	UserID   valueobject.UserID
	Username string
	Reason   string
	Operator string
	Duration *time.Duration // 封禁时长，nil表示永久封禁
}

// NewUserBannedEvent 创建用户封禁事件
func NewUserBannedEvent(user *entity.User, reason, operator string, duration *time.Duration) *UserBannedEvent {
	return &UserBannedEvent{
		BaseEvent: event.NewBaseEvent("user.banned", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		Reason:    reason,
		Operator:  operator,
		Duration:  duration,
	}
}

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
	*event.BaseEvent
	UserID      valueobject.UserID
	Username    string
	UpdatedFields map[string]interface{}
}

// NewUserProfileUpdatedEvent 创建用户资料更新事件
func NewUserProfileUpdatedEvent(user *entity.User, updatedFields map[string]interface{}) *UserProfileUpdatedEvent {
	return &UserProfileUpdatedEvent{
		BaseEvent:     event.NewBaseEvent("user.profile.updated", user.ID().Value(), time.Now()),
		UserID:        user.ID(),
		Username:      user.Username().Value(),
		UpdatedFields: updatedFields,
	}
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	*event.BaseEvent
	UserID    valueobject.UserID
	Username  string
	ClientIP  string
	UserAgent string
}

// NewUserPasswordChangedEvent 创建用户密码修改事件
func NewUserPasswordChangedEvent(user *entity.User, clientIP, userAgent string) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseEvent: event.NewBaseEvent("user.password.changed", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		ClientIP:  clientIP,
		UserAgent: userAgent,
	}
}

// UserDeletedEvent 用户删除事件
type UserDeletedEvent struct {
	*event.BaseEvent
	UserID   valueobject.UserID
	Username string
	Email    string
	Reason   string
	Operator string
}

// NewUserDeletedEvent 创建用户删除事件
func NewUserDeletedEvent(user *entity.User, reason, operator string) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: event.NewBaseEvent("user.deleted", user.ID().Value(), time.Now()),
		UserID:    user.ID(),
		Username:  user.Username().Value(),
		Email:     user.Email().Value(),
		Reason:    reason,
		Operator:  operator,
	}
}

// UserBehaviorRecordedEvent 用户行为记录事件
type UserBehaviorRecordedEvent struct {
	*event.BaseEvent
	UserID       valueobject.UserID
	Username     string
	Action       string
	ResourceType string
	ResourceID   int64
	ClientIP     string
	UserAgent    string
}

// NewUserBehaviorRecordedEvent 创建用户行为记录事件
func NewUserBehaviorRecordedEvent(
	user *entity.User,
	action, resourceType string,
	resourceID int64,
	clientIP, userAgent string,
) *UserBehaviorRecordedEvent {
	return &UserBehaviorRecordedEvent{
		BaseEvent:    event.NewBaseEvent("user.behavior.recorded", user.ID().Value(), time.Now()),
		UserID:       user.ID(),
		Username:     user.Username().Value(),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
	}
}
