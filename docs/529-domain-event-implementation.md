# 领域事件实现和修复总结

## 1. 问题概述

在DDD架构重构过程中，我们发现了与领域事件相关的几个关键问题：

1. 领域事件基类定义不完整，缺少必要的接口和方法
2. 各业务领域事件未正确实现领域事件接口
3. 事件发布接口与事件对象类型不匹配
4. 没有统一的事件处理机制

这些问题导致编译错误，主要体现在订单服务的领域事件实现中，如：

```
cannot use event.NewOrderCreatedEvent(order) as "wz-backend-go/internal/domain/shared/event".DomainEvent value in argument to s.eventPublisher.Publish: *"wz-backend-go/internal/domain/order/event".OrderCreatedEvent does not implement "wz-backend-go/internal/domain/shared/event".DomainEvent (missing method AggregateID)
```

## 2. 解决方案

### 2.1 完善领域事件基础设施

修改了`internal/domain/shared/event/domain_event.go`文件，完善了领域事件接口定义：

```go
package event

import "time"

// DomainEvent 领域事件接口
type DomainEvent interface {
	// EventID 获取事件ID
	EventID() string
	
	// EventType 获取事件类型
	EventType() string
	
	// OccurredAt 获取事件发生时间
	OccurredAt() time.Time
	
	// AggregateID 获取聚合根ID
	AggregateID() string
}

// BaseEvent 领域事件基类
type BaseEvent struct {
	eventID    string    // 事件ID
	eventType  string    // 事件类型
	occurredAt time.Time // 事件发生时间
	aggregateID string   // 聚合根ID
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventID, eventType, aggregateID string, occurredAt time.Time) BaseEvent {
	return BaseEvent{
		eventID:    eventID,
		eventType:  eventType,
		occurredAt: occurredAt,
		aggregateID: aggregateID,
	}
}

// EventID 获取事件ID
func (e BaseEvent) EventID() string {
	return e.eventID
}

// EventType 获取事件类型
func (e BaseEvent) EventType() string {
	return e.eventType
}

// OccurredAt 获取事件发生时间
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID 获取聚合根ID
func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}
```

### 2.2 修复订单领域事件实现

修改了`internal/domain/order/event/order_events.go`文件，确保所有事件正确实现`DomainEvent`接口：

```go
// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	CustomerID  int64
	TotalAmount int64
	ItemCount   int
	CreatedAt   time.Time
}

// NewOrderCreatedEvent 创建订单创建事件
func NewOrderCreatedEvent(order *entity.Order) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent: event.NewBaseEvent(
			uuid.New().String(),
			"order.created",
			order.ID().Value(),  // 聚合根ID
			time.Now(),
		),
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		TotalAmount: order.TotalAmount().Amount(),
		ItemCount:   len(order.Items()),
		CreatedAt:   order.CreatedAt(),
	}
}
```

对所有订单相关事件都进行了类似修改，确保正确实现`AggregateID()`方法。

### 2.3 修复仓储接口定义

修改了`internal/domain/order/repository/order_repository.go`文件，修复了类型导入问题：

```go
package repository

import (
	"wz-backend-go/internal/domain/order/entity"
	"wz-backend-go/internal/domain/order/valueobject"
	"wz-backend-go/internal/domain/shared/event"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// FindByID 根据ID查找订单
	FindByID(id valueobject.OrderID) (*entity.Order, error)
	
	// 其他方法...
}
```

### 2.4 实现事件总线

创建了`internal/infrastructure/messaging/eventbus/event_bus.go`文件，实现了简单的事件总线：

```go
package eventbus

import (
	"sync"
	"wz-backend-go/internal/domain/shared/event"
)

// 定义事件处理器类型
type EventHandler func(event event.DomainEvent) error

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = []EventHandler{}
	}
	
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish 发布事件
func (b *EventBus) Publish(event event.DomainEvent) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	
	if handlers, exists := b.handlers[event.EventType()]; exists {
		for _, handler := range handlers {
			if err := handler(event); err != nil {
				return err
			}
		}
	}
	
	return nil
}
```

## 3. 修复后的领域事件流程

1. **事件创建**: 在领域层中，当领域状态发生变化时，创建相应的领域事件
2. **事件发布**: 通过`EventPublisher`接口发布事件
3. **事件订阅**: 应用层或其他服务可以订阅特定类型的事件
4. **事件处理**: 当事件发布时，所有订阅该事件类型的处理器会被调用

示例代码（订单支付场景）：

```go
// 在领域服务中
func (s *OrderDomainService) PayOrder(orderID valueobject.OrderID, paymentMethod valueobject.PaymentMethod) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	// 支付订单
	err = order.Pay(paymentMethod)
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单支付事件
	s.eventPublisher.Publish(event.NewOrderPaidEvent(order))

	return nil
}

// 在应用层中订阅事件
eventBus.Subscribe("order.paid", func(e event.DomainEvent) error {
	paidEvent := e.(*order.OrderPaidEvent) // 类型断言
	
	// 执行订单支付后的业务逻辑，如更新库存、创建发票等
	
	return nil
})
```

## 4. 后续工作

1. **完善异步事件处理**: 实现基于消息队列的异步事件处理机制
2. **事件持久化**: 实现事件存储，支持事件回放和事件溯源
3. **事件监控**: 添加事件处理监控和日志记录
4. **跨服务事件**: 实现跨服务的事件传递机制
5. **单元测试**: 为事件相关代码添加单元测试 