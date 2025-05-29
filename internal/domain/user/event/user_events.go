package event

import (
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/user/valueobject"
)

// 用户相关事件类型常量
const (
	UserCreatedEventType          = "user.created"
	UserUpdatedEventType          = "user.updated"
	UserVerifiedEventType         = "user.verified"
	UserCompanyVerifiedEventType  = "user.company_verified"
	UserPasswordChangedEventType  = "user.password_changed"
	UserLoggedInEventType         = "user.logged_in"
	UserBehaviorRecordedEventType = "user.behavior_recorded"
)

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	eventID     string
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

func NewBaseDomainEvent(eventType string, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     uuid.New().String(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
	}
}

func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseDomainEvent) OccurredTime() time.Time {
	return e.occurredAt
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	BaseDomainEvent
	UserID   int64
	Username string
	Email    string
}

func NewUserCreatedEvent(userID valueobject.UserID, username valueobject.Username) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserCreatedEventType,
			userID.String(),
		),
		UserID:   userID.Value(),
		Username: username.Value(),
	}
}

// UserVerifiedEvent 用户验证事件
type UserVerifiedEvent struct {
	BaseDomainEvent
	UserID int64
}

func NewUserVerifiedEvent(userID valueobject.UserID) *UserVerifiedEvent {
	return &UserVerifiedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserVerifiedEventType,
			userID.String(),
		),
		UserID: userID.Value(),
	}
}

// UserCompanyVerifiedEvent 用户企业验证事件
type UserCompanyVerifiedEvent struct {
	BaseDomainEvent
	UserID int64
}

func NewUserCompanyVerifiedEvent(userID valueobject.UserID) *UserCompanyVerifiedEvent {
	return &UserCompanyVerifiedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserCompanyVerifiedEventType,
			userID.String(),
		),
		UserID: userID.Value(),
	}
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	BaseDomainEvent
	UserID int64
}

func NewUserPasswordChangedEvent(userID valueobject.UserID) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserPasswordChangedEventType,
			userID.String(),
		),
		UserID: userID.Value(),
	}
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	BaseDomainEvent
	UserID    int64
	IP        string
	UserAgent string
}

func NewUserLoggedInEvent(userID valueobject.UserID, ip string, userAgent string) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserLoggedInEventType,
			userID.String(),
		),
		UserID:    userID.Value(),
		IP:        ip,
		UserAgent: userAgent,
	}
}

// UserBehaviorRecordedEvent 用户行为记录事件
type UserBehaviorRecordedEvent struct {
	BaseDomainEvent
	UserID       int64
	Action       string
	ResourceType string
	ResourceID   int64
}

func NewUserBehaviorRecordedEvent(
	userID valueobject.UserID,
	action string,
	resourceType string,
	resourceID valueobject.UserID,
) *UserBehaviorRecordedEvent {
	return &UserBehaviorRecordedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			UserBehaviorRecordedEventType,
			userID.String(),
		),
		UserID:       userID.Value(),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID.Value(),
	}
}
