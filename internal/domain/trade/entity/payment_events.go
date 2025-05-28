package entity

import (
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

const (
	EventTypePaymentCreated          = "payment.created"
	EventTypePaymentSucceeded        = "payment.succeeded"
	EventTypePaymentFailed           = "payment.failed"
	EventTypePaymentRefundRequested  = "payment.refund_requested"
	EventTypePaymentRefundCompleted  = "payment.refund_completed"
)

// PaymentCreatedEvent 支付创建事件
type PaymentCreatedEvent struct {
	event.BaseDomainEvent
}

// NewPaymentCreatedEvent 创建新的支付创建事件
func NewPaymentCreatedEvent(payment *Payment) PaymentCreatedEvent {
	return PaymentCreatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypePaymentCreated,
			payment.ID().String(),
			"Payment",
			map[string]interface{}{
				"id":        payment.ID().String(),
				"orderID":   payment.OrderID().String(),
				"userID":    payment.UserID().String(),
				"amount":    payment.Amount().Amount(),
				"currency":  payment.Amount().Currency(),
				"method":    string(payment.Method()),
				"status":    string(payment.Status()),
				"createdAt": payment.CreatedAt(),
			},
		),
	}
}

// PaymentSucceededEvent 支付成功事件
type PaymentSucceededEvent struct {
	event.BaseDomainEvent
}

// NewPaymentSucceededEvent 创建新的支付成功事件
func NewPaymentSucceededEvent(payment *Payment) PaymentSucceededEvent {
	return PaymentSucceededEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypePaymentSucceeded,
			payment.ID().String(),
			"Payment",
			map[string]interface{}{
				"id":            payment.ID().String(),
				"orderID":       payment.OrderID().String(),
				"userID":        payment.UserID().String(),
				"amount":        payment.Amount().Amount(),
				"currency":      payment.Amount().Currency(),
				"method":        string(payment.Method()),
				"transactionID": payment.TransactionID(),
				"paymentTime":   payment.PaymentTime(),
			},
		),
	}
}

// PaymentFailedEvent 支付失败事件
type PaymentFailedEvent struct {
	event.BaseDomainEvent
}

// NewPaymentFailedEvent 创建新的支付失败事件
func NewPaymentFailedEvent(payment *Payment, reason string) PaymentFailedEvent {
	return PaymentFailedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypePaymentFailed,
			payment.ID().String(),
			"Payment",
			map[string]interface{}{
				"id":        payment.ID().String(),
				"orderID":   payment.OrderID().String(),
				"userID":    payment.UserID().String(),
				"amount":    payment.Amount().Amount(),
				"currency":  payment.Amount().Currency(),
				"method":    string(payment.Method()),
				"reason":    reason,
				"failedAt":  payment.UpdatedAt(),
			},
		),
	}
}

// PaymentRefundRequestedEvent 支付申请退款事件
type PaymentRefundRequestedEvent struct {
	event.BaseDomainEvent
}

// NewPaymentRefundRequestedEvent 创建新的支付申请退款事件
func NewPaymentRefundRequestedEvent(payment *Payment) PaymentRefundRequestedEvent {
	return PaymentRefundRequestedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypePaymentRefundRequested,
			payment.ID().String(),
			"Payment",
			map[string]interface{}{
				"id":                payment.ID().String(),
				"orderID":           payment.OrderID().String(),
				"userID":            payment.UserID().String(),
				"amount":            payment.Amount().Amount(),
				"currency":          payment.Amount().Currency(),
				"transactionID":     payment.TransactionID(),
				"refundRequestedAt": payment.UpdatedAt(),
			},
		),
	}
}

// PaymentRefundCompletedEvent 支付退款完成事件
type PaymentRefundCompletedEvent struct {
	event.BaseDomainEvent
}

// NewPaymentRefundCompletedEvent 创建新的支付退款完成事件
func NewPaymentRefundCompletedEvent(payment *Payment) PaymentRefundCompletedEvent {
	return PaymentRefundCompletedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypePaymentRefundCompleted,
			payment.ID().String(),
			"Payment",
			map[string]interface{}{
				"id":            payment.ID().String(),
				"orderID":       payment.OrderID().String(),
				"userID":        payment.UserID().String(),
				"amount":        payment.Amount().Amount(),
				"currency":      payment.Amount().Currency(),
				"transactionID": payment.TransactionID(),
				"refundTime":    payment.RefundTime(),
			},
		),
	}
}
