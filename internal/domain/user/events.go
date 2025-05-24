package user

import (
	"time"
)

// 定义领域事件类型常量
const (
	EventTypeUserRegistered   = "user.registered"
	EventTypeUserLoggedIn     = "user.logged_in"
	EventTypeUserVerified     = "user.verified"
	EventTypePasswordChanged  = "user.password_changed"
	EventTypeProfileUpdated   = "user.profile_updated"
	EventTypeRoleChanged      = "user.role_changed"
	EventTypeUserDisabled     = "user.disabled"
	EventTypeUserEnabled      = "user.enabled"
	EventTypeCompanyVerified  = "user.company_verified"
)

// Event 用户领域事件接口
type Event interface {
	EventType() string
	OccurredAt() time.Time
	UserID() UserID
}

// BaseEvent 基础事件结构，所有事件的基类
type BaseEvent struct {
	eventType  string
	occurredAt time.Time
	userID     UserID
}

// EventType 获取事件类型
func (e BaseEvent) EventType() string {
	return e.eventType
}

// OccurredAt 获取事件发生时间
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// UserID 获取关联的用户ID
func (e BaseEvent) UserID() UserID {
	return e.userID
}

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	BaseEvent
	username string
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID UserID, username string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeUserRegistered,
			occurredAt: time.Now(),
			userID:     userID,
		},
		username: username,
	}
}

// Username 获取用户名
func (e UserRegisteredEvent) Username() string {
	return e.username
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	BaseEvent
	loginTime time.Time
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID UserID, loginTime time.Time) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeUserLoggedIn,
			occurredAt: time.Now(),
			userID:     userID,
		},
		loginTime: loginTime,
	}
}

// LoginTime 获取登录时间
func (e UserLoggedInEvent) LoginTime() time.Time {
	return e.loginTime
}

// UserVerifiedEvent 用户验证事件
type UserVerifiedEvent struct {
	BaseEvent
}

// NewUserVerifiedEvent 创建用户验证事件
func NewUserVerifiedEvent(userID UserID) *UserVerifiedEvent {
	return &UserVerifiedEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeUserVerified,
			occurredAt: time.Now(),
			userID:     userID,
		},
	}
}

// PasswordChangedEvent 密码修改事件
type PasswordChangedEvent struct {
	BaseEvent
}

// NewPasswordChangedEvent 创建密码修改事件
func NewPasswordChangedEvent(userID UserID) *PasswordChangedEvent {
	return &PasswordChangedEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypePasswordChanged,
			occurredAt: time.Now(),
			userID:     userID,
		},
	}
}

// ProfileUpdatedEvent 资料更新事件
type ProfileUpdatedEvent struct {
	BaseEvent
}

// NewProfileUpdatedEvent 创建资料更新事件
func NewProfileUpdatedEvent(userID UserID) *ProfileUpdatedEvent {
	return &ProfileUpdatedEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeProfileUpdated,
			occurredAt: time.Now(),
			userID:     userID,
		},
	}
}

// RoleChangedEvent 角色变更事件
type RoleChangedEvent struct {
	BaseEvent
	newRole UserRole
}

// NewRoleChangedEvent 创建角色变更事件
func NewRoleChangedEvent(userID UserID, newRole UserRole) *RoleChangedEvent {
	return &RoleChangedEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeRoleChanged,
			occurredAt: time.Now(),
			userID:     userID,
		},
		newRole: newRole,
	}
}

// NewRole 获取新角色
func (e RoleChangedEvent) NewRole() UserRole {
	return e.newRole
}

// UserDisabledEvent 用户禁用事件
type UserDisabledEvent struct {
	BaseEvent
	reason     string
	operatorID UserID
}

// NewUserDisabledEvent 创建用户禁用事件
func NewUserDisabledEvent(userID UserID, reason string, operatorID UserID) *UserDisabledEvent {
	return &UserDisabledEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeUserDisabled,
			occurredAt: time.Now(),
			userID:     userID,
		},
		reason:     reason,
		operatorID: operatorID,
	}
}

// Reason 获取禁用原因
func (e UserDisabledEvent) Reason() string {
	return e.reason
}

// OperatorID 获取操作者ID
func (e UserDisabledEvent) OperatorID() UserID {
	return e.operatorID
}

// UserEnabledEvent 用户启用事件
type UserEnabledEvent struct {
	BaseEvent
	operatorID UserID
}

// NewUserEnabledEvent 创建用户启用事件
func NewUserEnabledEvent(userID UserID, operatorID UserID) *UserEnabledEvent {
	return &UserEnabledEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeUserEnabled,
			occurredAt: time.Now(),
			userID:     userID,
		},
		operatorID: operatorID,
	}
}

// OperatorID 获取操作者ID
func (e UserEnabledEvent) OperatorID() UserID {
	return e.operatorID
}

// CompanyVerifiedEvent 企业认证事件
type CompanyVerifiedEvent struct {
	BaseEvent
	companyName string
}

// NewCompanyVerifiedEvent 创建企业认证事件
func NewCompanyVerifiedEvent(userID UserID, companyName string) *CompanyVerifiedEvent {
	return &CompanyVerifiedEvent{
		BaseEvent: BaseEvent{
			eventType:  EventTypeCompanyVerified,
			occurredAt: time.Now(),
			userID:     userID,
		},
		companyName: companyName,
	}
}

// CompanyName 获取企业名称
func (e CompanyVerifiedEvent) CompanyName() string {
	return e.companyName
}
