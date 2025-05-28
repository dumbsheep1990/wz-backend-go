package event

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 是所有领域事件的接口
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	OccurredAt() time.Time
	Payload() interface{}
}

// BaseDomainEvent 是领域事件的基础实现
type BaseDomainEvent struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	payload       interface{}
}

// NewBaseDomainEvent 创建一个新的基础领域事件
func NewBaseDomainEvent(
	eventType string,
	aggregateID string,
	aggregateType string,
	payload interface{},
) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:       uuid.New().String(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		payload:       payload,
	}
}

// EventID 返回事件的唯一标识符
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType 返回事件的类型
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID 返回事件适用的聚合根标识符
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType 返回事件适用的聚合根类型
func (e BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// OccurredAt 返回事件发生的时间
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Payload 返回事件载荷
func (e BaseDomainEvent) Payload() interface{} {
	return e.payload
}
