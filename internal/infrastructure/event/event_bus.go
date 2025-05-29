package event

import (
	"sync"

	"wz-backend-go/internal/domain/shared/event"
)

// SimpleEventBus 简单的事件总线实现
type SimpleEventBus struct {
	handlers map[string][]func(event.DomainEvent) error
	mu       sync.RWMutex
}

// NewSimpleEventBus 创建一个简单的事件总线
func NewSimpleEventBus() event.EventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]func(event.DomainEvent) error),
	}
}

// Publish 发布事件
func (b *SimpleEventBus) Publish(ctx interface{}, event event.DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[event.EventType()]; ok {
		for _, handler := range handlers {
			if err := handler(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// Subscribe 订阅事件
func (b *SimpleEventBus) Subscribe(eventType string, handler func(event.DomainEvent) error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.handlers[eventType]; !ok {
		b.handlers[eventType] = []func(event.DomainEvent) error{}
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
