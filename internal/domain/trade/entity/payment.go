package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// Payment 支付实体
type Payment struct {
	id             valueobject.PaymentID
	orderID        valueobject.OrderID
	userID         valueobject.UserID
	amount         valueobject.Money
	method         valueobject.PaymentMethod
	status         valueobject.PaymentStatus
	transactionID  string // 第三方支付平台的交易ID
	paymentTime    *time.Time // 支付完成时间，未支付时为nil
	refundTime     *time.Time // 退款时间，未退款时为nil
	createdAt      time.Time
	updatedAt      time.Time
	
	domainEvents []event.DomainEvent
}

// NewPayment 创建一个新的支付实体
func NewPayment(
	orderID valueobject.OrderID,
	userID valueobject.UserID,
	amount valueobject.Money,
	method valueobject.PaymentMethod,
) (*Payment, error) {
	if orderID.IsEmpty() {
		return nil, errors.New("订单ID不能为空")
	}
	if userID.IsEmpty() {
		return nil, errors.New("用户ID不能为空")
	}
	
	now := time.Now()
	payment := &Payment{
		id:            valueobject.NewPaymentID(uuid.New().String()),
		orderID:       orderID,
		userID:        userID,
		amount:        amount,
		method:        method,
		status:        valueobject.PaymentStatusPending,
		transactionID: "",
		paymentTime:   nil,
		refundTime:    nil,
		createdAt:     now,
		updatedAt:     now,
		domainEvents:  []event.DomainEvent{},
	}
	
	// 添加支付创建事件
	payment.addDomainEvent(NewPaymentCreatedEvent(payment))
	
	return payment, nil
}

// ID 获取支付ID
func (p *Payment) ID() valueobject.PaymentID {
	return p.id
}

// OrderID 获取订单ID
func (p *Payment) OrderID() valueobject.OrderID {
	return p.orderID
}

// UserID 获取用户ID
func (p *Payment) UserID() valueobject.UserID {
	return p.userID
}

// Amount 获取支付金额
func (p *Payment) Amount() valueobject.Money {
	return p.amount
}

// Method 获取支付方式
func (p *Payment) Method() valueobject.PaymentMethod {
	return p.method
}

// Status 获取支付状态
func (p *Payment) Status() valueobject.PaymentStatus {
	return p.status
}

// TransactionID 获取交易ID
func (p *Payment) TransactionID() string {
	return p.transactionID
}

// PaymentTime 获取支付时间
func (p *Payment) PaymentTime() *time.Time {
	return p.paymentTime
}

// RefundTime 获取退款时间
func (p *Payment) RefundTime() *time.Time {
	return p.refundTime
}

// CreatedAt 获取创建时间
func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt 获取更新时间
func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

// Pay 完成支付
func (p *Payment) Pay(transactionID string) error {
	if p.status != valueobject.PaymentStatusPending {
		return errors.New("只有待支付状态的支付才能标记为已支付")
	}
	
	if transactionID == "" {
		return errors.New("交易ID不能为空")
	}
	
	now := time.Now()
	p.status = valueobject.PaymentStatusSuccess
	p.transactionID = transactionID
	p.paymentTime = &now
	p.updatedAt = now
	
	// 添加支付成功事件
	p.addDomainEvent(NewPaymentSucceededEvent(p))
	
	return nil
}

// Fail 标记支付失败
func (p *Payment) Fail(reason string) error {
	if p.status != valueobject.PaymentStatusPending {
		return errors.New("只有待支付状态的支付才能标记为失败")
	}
	
	p.status = valueobject.PaymentStatusFailed
	p.updatedAt = time.Now()
	
	// 添加支付失败事件
	p.addDomainEvent(NewPaymentFailedEvent(p, reason))
	
	return nil
}

// RequestRefund 申请退款
func (p *Payment) RequestRefund() error {
	if p.status != valueobject.PaymentStatusSuccess {
		return errors.New("只有支付成功的支付才能申请退款")
	}
	
	p.status = valueobject.PaymentStatusRefunding
	p.updatedAt = time.Now()
	
	// 添加申请退款事件
	p.addDomainEvent(NewPaymentRefundRequestedEvent(p))
	
	return nil
}

// CompleteRefund 完成退款
func (p *Payment) CompleteRefund() error {
	if p.status != valueobject.PaymentStatusRefunding {
		return errors.New("只有退款中状态的支付才能标记为已退款")
	}
	
	now := time.Now()
	p.status = valueobject.PaymentStatusRefunded
	p.refundTime = &now
	p.updatedAt = now
	
	// 添加退款完成事件
	p.addDomainEvent(NewPaymentRefundCompletedEvent(p))
	
	return nil
}

// IsPaid 检查是否已支付
func (p *Payment) IsPaid() bool {
	return p.status == valueobject.PaymentStatusSuccess ||
		p.status == valueobject.PaymentStatusRefunding ||
		p.status == valueobject.PaymentStatusRefunded
}

// IsRefunded 检查是否已退款
func (p *Payment) IsRefunded() bool {
	return p.status == valueobject.PaymentStatusRefunded
}

// IsRefunding 检查是否正在退款中
func (p *Payment) IsRefunding() bool {
	return p.status == valueobject.PaymentStatusRefunding
}

// 添加领域事件
func (p *Payment) addDomainEvent(event event.DomainEvent) {
	p.domainEvents = append(p.domainEvents, event)
}

// GetDomainEvents 获取所有领域事件
func (p *Payment) GetDomainEvents() []event.DomainEvent {
	return p.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (p *Payment) ClearDomainEvents() {
	p.domainEvents = []event.DomainEvent{}
}
