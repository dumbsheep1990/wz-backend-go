package eventbus

import (
	"log"
	"sync"
)

// EventBus 实现了一个简单的事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event interface{}) error
}

// EventHandlerFunc 函数类型的事件处理器
type EventHandlerFunc func(event interface{}) error

// Handle 实现EventHandler接口
func (f EventHandlerFunc) Handle(event interface{}) error {
	return f(event)
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeFunc 使用函数订阅事件
func (b *EventBus) SubscribeFunc(eventType string, handlerFunc func(event interface{}) error) {
	b.Subscribe(eventType, EventHandlerFunc(handlerFunc))
}

// Publish 发布事件
func (b *EventBus) Publish(event interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 获取事件类型名称
	eventType := getEventType(event)
	
	// 获取事件处理器
	handlers, exists := b.handlers[eventType]
	if !exists {
		// 如果没有处理器，记录并返回
		log.Printf("没有处理器注册处理事件: %s", eventType)
		return nil
	}

	// 异步调用所有处理器
	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler, e interface{}) {
			defer wg.Done()
			if err := h.Handle(e); err != nil {
				log.Printf("处理事件 %s 失败: %v", eventType, err)
			}
		}(handler, event)
	}

	// 等待所有处理器完成
	wg.Wait()
	return nil
}

// getEventType 获取事件类型名称
func getEventType(event interface{}) string {
	// 根据事件类型获取类型名称，这里简化处理
	// 实际应用中可以使用反射获取类型名称
	switch e := event.(type) {
	default:
		return getTypeName(e)
	}
}

// getTypeName 获取类型名称
func getTypeName(v interface{}) string {
	if v == nil {
		return "nil"
	}
	// 这里简化处理，实际应用中可以使用反射获取完整类型名称
	return "event"
}
