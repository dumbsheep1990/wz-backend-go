package entity

import (
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserPointsCreatedEvent 表示创建积分记录的领域事件
type UserPointsCreatedEvent struct {
	eventID      string
	aggregateID  string
	userID       int64
	points       int
	pointsType   int
	source       string
	description  string
	operatorID   int64
	totalPoints  int
	occurredTime time.Time
}

// NewUserPointsCreatedEvent 创建一个新的积分记录创建事件
func NewUserPointsCreatedEvent(userPoints *UserPoints) event.DomainEvent {
	return &UserPointsCreatedEvent{
		eventID:      uuid.New().String(),
		aggregateID:  userPoints.ID().String(),
		userID:       userPoints.UserID().Value(),
		points:       userPoints.Points().Value(),
		pointsType:   userPoints.PointsType().Value(),
		source:       userPoints.Source().String(),
		description:  userPoints.Description().String(),
		operatorID:   userPoints.OperatorID().Value(),
		totalPoints:  userPoints.TotalPoints().Value(),
		occurredTime: time.Now(),
	}
}

// EventID 获取事件ID
func (e *UserPointsCreatedEvent) EventID() string {
	return e.eventID
}

// AggregateID 获取聚合根ID
func (e *UserPointsCreatedEvent) AggregateID() string {
	return e.aggregateID
}

// EventType 获取事件类型
func (e *UserPointsCreatedEvent) EventType() string {
	return "user.points.created"
}

// OccurredTime 获取事件发生时间
func (e *UserPointsCreatedEvent) OccurredTime() time.Time {
	return e.occurredTime
}

// UserID 获取用户ID
func (e *UserPointsCreatedEvent) UserID() int64 {
	return e.userID
}

// Points 获取积分值
func (e *UserPointsCreatedEvent) Points() int {
	return e.points
}

// PointsType 获取积分类型
func (e *UserPointsCreatedEvent) PointsType() int {
	return e.pointsType
}

// Source 获取积分来源
func (e *UserPointsCreatedEvent) Source() string {
	return e.source
}

// Description 获取描述
func (e *UserPointsCreatedEvent) Description() string {
	return e.description
}

// OperatorID 获取操作者ID
func (e *UserPointsCreatedEvent) OperatorID() int64 {
	return e.operatorID
}

// TotalPoints 获取总积分
func (e *UserPointsCreatedEvent) TotalPoints() int {
	return e.totalPoints
}

// UserPointsRevokedEvent 表示积分记录撤销的领域事件
type UserPointsRevokedEvent struct {
	eventID      string
	aggregateID  string
	userID       int64
	pointsID     string
	operatorID   int64
	occurredTime time.Time
}

// NewUserPointsRevokedEvent 创建一个新的积分记录撤销事件
func NewUserPointsRevokedEvent(userPoints *UserPoints, operatorID valueobject.UserID) event.DomainEvent {
	return &UserPointsRevokedEvent{
		eventID:      uuid.New().String(),
		aggregateID:  userPoints.ID().String(),
		userID:       userPoints.UserID().Value(),
		pointsID:     userPoints.ID().String(),
		operatorID:   operatorID.Value(),
		occurredTime: time.Now(),
	}
}

// EventID 获取事件ID
func (e *UserPointsRevokedEvent) EventID() string {
	return e.eventID
}

// AggregateID 获取聚合根ID
func (e *UserPointsRevokedEvent) AggregateID() string {
	return e.aggregateID
}

// EventType 获取事件类型
func (e *UserPointsRevokedEvent) EventType() string {
	return "user.points.revoked"
}

// OccurredTime 获取事件发生时间
func (e *UserPointsRevokedEvent) OccurredTime() time.Time {
	return e.occurredTime
}

// UserID 获取用户ID
func (e *UserPointsRevokedEvent) UserID() int64 {
	return e.userID
}

// PointsID 获取积分记录ID
func (e *UserPointsRevokedEvent) PointsID() string {
	return e.pointsID
}

// OperatorID 获取操作者ID
func (e *UserPointsRevokedEvent) OperatorID() int64 {
	return e.operatorID
}
