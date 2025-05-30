package entity

import (
	"time"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/trade/valueobject"
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
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	ClientIP      string    `json:"client_ip"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentCreatedEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentCreatedEvent) GetEventType() string {
	return "payment.created"
}

// GetOccurredAt 获取发生时间
func (e PaymentCreatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentCreatedEvent) GetEventData() interface{} {
	return e
}

// NewPaymentCreatedEvent 创建支付创建事件
func NewPaymentCreatedEvent(payment *Payment) event.DomainEvent {
	return PaymentCreatedEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		Amount:        payment.Amount().Amount(),
		Currency:      payment.Amount().Currency(),
		PaymentMethod: payment.Method().Value(),
		ClientIP:      payment.ClientIP(),
		OccurredAt:    time.Now(),
	}
}

// PaymentInfoSetEvent 支付信息设置事件
type PaymentInfoSetEvent struct {
	PaymentID     string    `json:"payment_id"`
	PaymentURL    string    `json:"payment_url"`
	QRCode        string    `json:"qr_code"`
	TransactionID string    `json:"transaction_id"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentInfoSetEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentInfoSetEvent) GetEventType() string {
	return "payment.info_set"
}

// GetOccurredAt 获取发生时间
func (e PaymentInfoSetEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentInfoSetEvent) GetEventData() interface{} {
	return e
}

// NewPaymentInfoSetEvent 创建支付信息设置事件
func NewPaymentInfoSetEvent(payment *Payment, paymentURL, qrCode, transactionID string) event.DomainEvent {
	return PaymentInfoSetEvent{
		PaymentID:     payment.ID().Value(),
		PaymentURL:    paymentURL,
		QRCode:        qrCode,
		TransactionID: transactionID,
		OccurredAt:    time.Now(),
	}
}

// PaymentStatusChangedEvent 支付状态变更事件
type PaymentStatusChangedEvent struct {
	PaymentID string    `json:"payment_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	OccurredAt time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentStatusChangedEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentStatusChangedEvent) GetEventType() string {
	return "payment.status_changed"
}

// GetOccurredAt 获取发生时间
func (e PaymentStatusChangedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentStatusChangedEvent) GetEventData() interface{} {
	return e
}

// NewPaymentStatusChangedEvent 创建支付状态变更事件
func NewPaymentStatusChangedEvent(payment *Payment, oldStatus, newStatus valueobject.PaymentStatus) event.DomainEvent {
	return PaymentStatusChangedEvent{
		PaymentID: payment.ID().Value(),
		OldStatus: oldStatus.Value(),
		NewStatus: newStatus.Value(),
		OccurredAt: time.Now(),
	}
}

// PaymentSuccessEvent 支付成功事件
type PaymentSuccessEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	PaymentTime   time.Time `json:"payment_time"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentSuccessEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentSuccessEvent) GetEventType() string {
	return "payment.success"
}

// GetOccurredAt 获取发生时间
func (e PaymentSuccessEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentSuccessEvent) GetEventData() interface{} {
	return e
}

// NewPaymentSuccessEvent 创建支付成功事件
func NewPaymentSuccessEvent(payment *Payment) event.DomainEvent {
	return PaymentSuccessEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		Amount:        payment.Amount().Amount(),
		Currency:      payment.Amount().Currency(),
		PaymentMethod: payment.Method().Value(),
		TransactionID: payment.TransactionID(),
		PaymentTime:   *payment.PaymentTime(),
		OccurredAt:    time.Now(),
	}
}

// PaymentFailedEvent 支付失败事件
type PaymentFailedEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	FailureReason string    `json:"failure_reason"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentFailedEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentFailedEvent) GetEventType() string {
	return "payment.failed"
}

// GetOccurredAt 获取发生时间
func (e PaymentFailedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentFailedEvent) GetEventData() interface{} {
	return e
}

// NewPaymentFailedEvent 创建支付失败事件
func NewPaymentFailedEvent(payment *Payment, reason string) event.DomainEvent {
	return PaymentFailedEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		PaymentMethod: payment.Method().Value(),
		FailureReason: reason,
		OccurredAt:    time.Now(),
	}
}

// PaymentCancelledEvent 支付取消事件
type PaymentCancelledEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	CancelReason  string    `json:"cancel_reason"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentCancelledEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentCancelledEvent) GetEventType() string {
	return "payment.cancelled"
}

// GetOccurredAt 获取发生时间
func (e PaymentCancelledEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentCancelledEvent) GetEventData() interface{} {
	return e
}

// NewPaymentCancelledEvent 创建支付取消事件
func NewPaymentCancelledEvent(payment *Payment, reason string) event.DomainEvent {
	return PaymentCancelledEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		PaymentMethod: payment.Method().Value(),
		CancelReason:  reason,
		OccurredAt:    time.Now(),
	}
}

// PaymentRefundedEvent 支付退款事件
type PaymentRefundedEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentRefundedEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentRefundedEvent) GetEventType() string {
	return "payment.refunded"
}

// GetOccurredAt 获取发生时间
func (e PaymentRefundedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentRefundedEvent) GetEventData() interface{} {
	return e
}

// NewPaymentRefundedEvent 创建支付退款事件
func NewPaymentRefundedEvent(payment *Payment) event.DomainEvent {
	return PaymentRefundedEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		Amount:        payment.Amount().Amount(),
		Currency:      payment.Amount().Currency(),
		PaymentMethod: payment.Method().Value(),
		OccurredAt:    time.Now(),
	}
}

// PaymentExpiredEvent 支付过期事件
type PaymentExpiredEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	ExpiredAt     time.Time `json:"expired_at"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// GetAggregateID 获取聚合ID
func (e PaymentExpiredEvent) GetAggregateID() string {
	return e.PaymentID
}

// GetEventType 获取事件类型
func (e PaymentExpiredEvent) GetEventType() string {
	return "payment.expired"
}

// GetOccurredAt 获取发生时间
func (e PaymentExpiredEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e PaymentExpiredEvent) GetEventData() interface{} {
	return e
}

// NewPaymentExpiredEvent 创建支付过期事件
func NewPaymentExpiredEvent(payment *Payment) event.DomainEvent {
	return PaymentExpiredEvent{
		PaymentID:     payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		PaymentMethod: payment.Method().Value(),
		ExpiredAt:     *payment.ExpiredAt(),
		OccurredAt:    time.Now(),
	}
}
