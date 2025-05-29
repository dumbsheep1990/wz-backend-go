package event

import "time"

// Event 领域事件接口
type Event interface {
	// EventID 获取事件ID
	EventID() string

	// EventType 获取事件类型
	EventType() string

	// AggregateID 获取聚合根ID
	AggregateID() string

	// OccurredTime 获取事件发生时间
	OccurredTime() time.Time
}
