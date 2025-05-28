package entity

import (
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

const (
	EventTypeOrderCreated         = "order.created"
	EventTypeOrderPaid            = "order.paid"
	EventTypeOrderShipped         = "order.shipped"
	EventTypeOrderDelivered       = "order.delivered"
	EventTypeOrderCompleted       = "order.completed"
	EventTypeOrderCancelled       = "order.cancelled"
	EventTypeOrderRefundRequested = "order.refund_requested"
	EventTypeOrderRefunded        = "order.refunded"
)

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	event.BaseDomainEvent
}

// NewOrderCreatedEvent 创建新的订单创建事件
func NewOrderCreatedEvent(order *Order) OrderCreatedEvent {
	return OrderCreatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderCreated,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":             order.ID().String(),
				"userID":         order.UserID().String(),
				"totalAmount":    order.TotalAmount().Amount(),
				"currency":       order.TotalAmount().Currency(),
				"status":         string(order.Status()),
				"createdAt":      order.CreatedAt(),
			},
		),
	}
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	event.BaseDomainEvent
}

// NewOrderPaidEvent 创建新的订单支付事件
func NewOrderPaidEvent(order *Order) OrderPaidEvent {
	return OrderPaidEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderPaid,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":          order.ID().String(),
				"userID":      order.UserID().String(),
				"paymentID":   order.PaymentID().String(),
				"totalAmount": order.TotalAmount().Amount(),
				"currency":    order.TotalAmount().Currency(),
				"paidAt":      order.UpdatedAt(),
			},
		),
	}
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	event.BaseDomainEvent
}

// NewOrderShippedEvent 创建新的订单发货事件
func NewOrderShippedEvent(order *Order) OrderShippedEvent {
	return OrderShippedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderShipped,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":         order.ID().String(),
				"userID":     order.UserID().String(),
				"shippingAddress": map[string]string{
					"province":    order.ShippingAddress().Province(),
					"city":        order.ShippingAddress().City(),
					"district":    order.ShippingAddress().District(),
					"detail":      order.ShippingAddress().Detail(),
					"receiver":    order.ShippingAddress().Receiver(),
					"phoneNumber": order.ShippingAddress().PhoneNumber(),
				},
				"shippedAt": order.UpdatedAt(),
			},
		),
	}
}

// OrderDeliveredEvent 订单送达事件
type OrderDeliveredEvent struct {
	event.BaseDomainEvent
}

// NewOrderDeliveredEvent 创建新的订单送达事件
func NewOrderDeliveredEvent(order *Order) OrderDeliveredEvent {
	return OrderDeliveredEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderDelivered,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":          order.ID().String(),
				"userID":      order.UserID().String(),
				"deliveredAt": order.UpdatedAt(),
			},
		),
	}
}

// OrderCompletedEvent 订单完成事件
type OrderCompletedEvent struct {
	event.BaseDomainEvent
}

// NewOrderCompletedEvent 创建新的订单完成事件
func NewOrderCompletedEvent(order *Order) OrderCompletedEvent {
	return OrderCompletedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderCompleted,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":          order.ID().String(),
				"userID":      order.UserID().String(),
				"completedAt": order.UpdatedAt(),
			},
		),
	}
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	event.BaseDomainEvent
}

// NewOrderCancelledEvent 创建新的订单取消事件
func NewOrderCancelledEvent(order *Order) OrderCancelledEvent {
	return OrderCancelledEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderCancelled,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":          order.ID().String(),
				"userID":      order.UserID().String(),
				"cancelledAt": order.UpdatedAt(),
			},
		),
	}
}

// OrderRefundRequestedEvent 订单申请退款事件
type OrderRefundRequestedEvent struct {
	event.BaseDomainEvent
}

// NewOrderRefundRequestedEvent 创建新的订单申请退款事件
func NewOrderRefundRequestedEvent(order *Order) OrderRefundRequestedEvent {
	return OrderRefundRequestedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderRefundRequested,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":                 order.ID().String(),
				"userID":             order.UserID().String(),
				"totalAmount":        order.TotalAmount().Amount(),
				"currency":           order.TotalAmount().Currency(),
				"refundRequestedAt":  order.UpdatedAt(),
			},
		),
	}
}

// OrderRefundedEvent 订单退款事件
type OrderRefundedEvent struct {
	event.BaseDomainEvent
}

// NewOrderRefundedEvent 创建新的订单退款事件
func NewOrderRefundedEvent(order *Order) OrderRefundedEvent {
	return OrderRefundedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeOrderRefunded,
			order.ID().String(),
			"Order",
			map[string]interface{}{
				"id":          order.ID().String(),
				"userID":      order.UserID().String(),
				"totalAmount": order.TotalAmount().Amount(),
				"currency":    order.TotalAmount().Currency(),
				"refundedAt":  order.UpdatedAt(),
			},
		),
	}
}
