package event

import (
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/order/entity"
	"wz-backend-go/internal/domain/shared/event"
)

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
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.created",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		TotalAmount: order.TotalAmount().Amount(),
		ItemCount:   len(order.Items()),
		CreatedAt:   order.CreatedAt(),
	}
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	event.BaseEvent
	OrderID       string
	OrderNumber   string
	CustomerID    int64
	TotalAmount   int64
	PaymentMethod int32
	PaidAt        time.Time
}

// NewOrderPaidEvent 创建订单支付事件
func NewOrderPaidEvent(order *entity.Order) *OrderPaidEvent {
	return &OrderPaidEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.paid",
			OccurredAt: time.Now(),
		},
		OrderID:       order.ID().Value(),
		OrderNumber:   order.OrderNumber().Value(),
		CustomerID:    order.CustomerID().Value(),
		TotalAmount:   order.TotalAmount().Amount(),
		PaymentMethod: order.PaymentMethod().Value(),
		PaidAt:        *order.PaidAt(),
	}
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	event.BaseEvent
	OrderID        string
	OrderNumber    string
	CustomerID     int64
	TrackingNumber string
	ShippingMethod int32
	ShippedAt      time.Time
}

// NewOrderShippedEvent 创建订单发货事件
func NewOrderShippedEvent(order *entity.Order) *OrderShippedEvent {
	return &OrderShippedEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.shipped",
			OccurredAt: time.Now(),
		},
		OrderID:        order.ID().Value(),
		OrderNumber:    order.OrderNumber().Value(),
		CustomerID:     order.CustomerID().Value(),
		TrackingNumber: order.TrackingNumber(),
		ShippingMethod: order.ShippingMethod().Value(),
		ShippedAt:      *order.ShippedAt(),
	}
}

// OrderDeliveredEvent 订单送达事件
type OrderDeliveredEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	CustomerID  int64
	DeliveredAt time.Time
}

// NewOrderDeliveredEvent 创建订单送达事件
func NewOrderDeliveredEvent(order *entity.Order) *OrderDeliveredEvent {
	return &OrderDeliveredEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.delivered",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		DeliveredAt: *order.DeliveredAt(),
	}
}

// OrderCompletedEvent 订单完成事件
type OrderCompletedEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	CustomerID  int64
	TotalAmount int64
	CompletedAt time.Time
}

// NewOrderCompletedEvent 创建订单完成事件
func NewOrderCompletedEvent(order *entity.Order) *OrderCompletedEvent {
	return &OrderCompletedEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.completed",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		TotalAmount: order.TotalAmount().Amount(),
		CompletedAt: *order.CompletedAt(),
	}
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	CustomerID  int64
	CancelledAt time.Time
	Reason      string
}

// NewOrderCancelledEvent 创建订单取消事件
func NewOrderCancelledEvent(order *entity.Order, reason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.cancelled",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		CancelledAt: *order.CancelledAt(),
		Reason:      reason,
	}
}

// OrderRefundedEvent 订单退款事件
type OrderRefundedEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	CustomerID  int64
	TotalAmount int64
	RefundedAt  time.Time
	Reason      string
}

// NewOrderRefundedEvent 创建订单退款事件
func NewOrderRefundedEvent(order *entity.Order, reason string) *OrderRefundedEvent {
	return &OrderRefundedEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.refunded",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		TotalAmount: order.TotalAmount().Amount(),
		RefundedAt:  *order.RefundedAt(),
		Reason:      reason,
	}
}

// OrderItemAddedEvent 订单项添加事件
type OrderItemAddedEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	ItemID      string
	ProductID   int64
	ProductName string
	Quantity    int32
	UnitPrice   int64
	TotalPrice  int64
}

// NewOrderItemAddedEvent 创建订单项添加事件
func NewOrderItemAddedEvent(order *entity.Order, item *entity.OrderItem) *OrderItemAddedEvent {
	return &OrderItemAddedEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.item_added",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		ItemID:      item.ID().Value(),
		ProductID:   item.ProductID().Value(),
		ProductName: item.ProductName(),
		Quantity:    item.Quantity(),
		UnitPrice:   item.UnitPrice().Amount(),
		TotalPrice:  item.TotalPrice().Amount(),
	}
}

// OrderItemRemovedEvent 订单项移除事件
type OrderItemRemovedEvent struct {
	event.BaseEvent
	OrderID     string
	OrderNumber string
	ItemID      string
	ProductID   int64
	ProductName string
}

// NewOrderItemRemovedEvent 创建订单项移除事件
func NewOrderItemRemovedEvent(order *entity.Order, item *entity.OrderItem) *OrderItemRemovedEvent {
	return &OrderItemRemovedEvent{
		BaseEvent: event.BaseEvent{
			EventID:    uuid.New().String(),
			EventType:  "order.item_removed",
			OccurredAt: time.Now(),
		},
		OrderID:     order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		ItemID:      item.ID().Value(),
		ProductID:   item.ProductID().Value(),
		ProductName: item.ProductName(),
	}
}
