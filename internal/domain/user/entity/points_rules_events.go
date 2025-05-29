package entity

import (
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/shared/event"
)

// PointsRulesCreatedEvent 表示创建积分规则的领域事件
type PointsRulesCreatedEvent struct {
	eventID           string
	aggregateID       string
	tenantID          int64
	signInPoints      int
	commentPoints     int
	sharePoints       int
	articlePoints     int
	invitePoints      int
	purchaseRate      int
	maxDailyPoints    int
	enableExchange    bool
	exchangeRate      int
	minExchangePoints int
	occurredTime      time.Time
}

// NewPointsRulesCreatedEvent 创建一个新的积分规则创建事件
func NewPointsRulesCreatedEvent(rules *PointsRules) event.DomainEvent {
	return &PointsRulesCreatedEvent{
		eventID:           uuid.New().String(),
		aggregateID:       rules.ID().String(),
		tenantID:          rules.TenantID().Value(),
		signInPoints:      rules.SignInPoints(),
		commentPoints:     rules.CommentPoints(),
		sharePoints:       rules.SharePoints(),
		articlePoints:     rules.ArticlePoints(),
		invitePoints:      rules.InvitePoints(),
		purchaseRate:      rules.PurchaseRate(),
		maxDailyPoints:    rules.MaxDailyPoints(),
		enableExchange:    rules.EnableExchange(),
		exchangeRate:      rules.ExchangeRate(),
		minExchangePoints: rules.MinExchangePoints(),
		occurredTime:      time.Now(),
	}
}

// EventID 获取事件ID
func (e *PointsRulesCreatedEvent) EventID() string {
	return e.eventID
}

// AggregateID 获取聚合根ID
func (e *PointsRulesCreatedEvent) AggregateID() string {
	return e.aggregateID
}

// EventType 获取事件类型
func (e *PointsRulesCreatedEvent) EventType() string {
	return "user.points.rules.created"
}

// OccurredTime 获取事件发生时间
func (e *PointsRulesCreatedEvent) OccurredTime() time.Time {
	return e.occurredTime
}

// TenantID 获取租户ID
func (e *PointsRulesCreatedEvent) TenantID() int64 {
	return e.tenantID
}

// PointsRulesUpdatedEvent 表示更新积分规则的领域事件
type PointsRulesUpdatedEvent struct {
	eventID           string
	aggregateID       string
	tenantID          int64
	signInPoints      int
	commentPoints     int
	sharePoints       int
	articlePoints     int
	invitePoints      int
	purchaseRate      int
	maxDailyPoints    int
	enableExchange    bool
	exchangeRate      int
	minExchangePoints int
	occurredTime      time.Time
}

// NewPointsRulesUpdatedEvent 创建一个新的积分规则更新事件
func NewPointsRulesUpdatedEvent(rules *PointsRules) event.DomainEvent {
	return &PointsRulesUpdatedEvent{
		eventID:           uuid.New().String(),
		aggregateID:       rules.ID().String(),
		tenantID:          rules.TenantID().Value(),
		signInPoints:      rules.SignInPoints(),
		commentPoints:     rules.CommentPoints(),
		sharePoints:       rules.SharePoints(),
		articlePoints:     rules.ArticlePoints(),
		invitePoints:      rules.InvitePoints(),
		purchaseRate:      rules.PurchaseRate(),
		maxDailyPoints:    rules.MaxDailyPoints(),
		enableExchange:    rules.EnableExchange(),
		exchangeRate:      rules.ExchangeRate(),
		minExchangePoints: rules.MinExchangePoints(),
		occurredTime:      time.Now(),
	}
}

// EventID 获取事件ID
func (e *PointsRulesUpdatedEvent) EventID() string {
	return e.eventID
}

// AggregateID 获取聚合根ID
func (e *PointsRulesUpdatedEvent) AggregateID() string {
	return e.aggregateID
}

// EventType 获取事件类型
func (e *PointsRulesUpdatedEvent) EventType() string {
	return "user.points.rules.updated"
}

// OccurredTime 获取事件发生时间
func (e *PointsRulesUpdatedEvent) OccurredTime() time.Time {
	return e.occurredTime
}

// TenantID 获取租户ID
func (e *PointsRulesUpdatedEvent) TenantID() int64 {
	return e.tenantID
}
