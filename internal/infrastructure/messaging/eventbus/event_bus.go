package eventbus

import (
	"context"
	"sync"

	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

// EventHandler defines a function that handles a domain event
type EventHandler func(ctx context.Context, event event.DomainEvent) error

// EventBus 是事件总线接口的内存实现
type EventBus struct {
	handlers map[string][]EventHandler // 事件类型到处理函数的映射
	mu       sync.RWMutex              // 读写锁，保证并发安全
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 为特定事件类型注册处理函数
func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = []EventHandler{}
	}
	
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish 将事件发布给所有注册的处理函数
func (b *EventBus) Publish(ctx context.Context, event event.DomainEvent) error {
	b.mu.RLock() // 加读锁，允许并发读取
	handlers, exists := b.handlers[event.EventType()]
	b.mu.RUnlock() // 及时释放读锁
	
	if !exists {
		return nil // 没有为该事件类型注册的处理函数
	}
	
	// 执行所有处理函数
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			// 在实际应用中，可能需要不同的错误处理方式
			// 目前我们只是继续处理其他处理函数
			// 可以考虑记录错误日志或实现重试机制
		}
	}
	
	return nil
}
