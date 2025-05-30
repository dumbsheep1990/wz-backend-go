package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// Payment 支付聚合根
type Payment struct {
	id            valueobject.PaymentID
	orderID       string // 订单ID，应该是OrderID值对象，但为简化暂用string
	userID        string // 用户ID
	amount        valueobject.Money
	method        valueobject.PaymentMethod
	status        valueobject.PaymentStatus
	transactionID string    // 第三方支付平台交易ID
	paymentURL    string    // 支付URL（第三方支付）
	qrCode        string    // 支付二维码
	clientIP      string    // 客户端IP
	returnURL     string    // 支付成功返回URL
	notifyURL     string    // 支付回调URL
	failureReason string    // 失败原因
	metadata      string    // 扩展信息（JSON格式）
	paymentTime   *time.Time // 支付成功时间
	expiredAt     *time.Time // 过期时间
	createdAt     time.Time
	updatedAt     time.Time
	
	// 领域事件
	domainEvents []event.DomainEvent
}

// NewPayment 创建新支付
func NewPayment(
	id valueobject.PaymentID,
	orderID string,
	userID string,
	amount valueobject.Money,
	method valueobject.PaymentMethod,
	clientIP string,
) (*Payment, error) {
	if id.IsEmpty() {
		return nil, errors.New("支付ID不能为空")
	}
	if orderID == "" {
		return nil, errors.New("订单ID不能为空")
	}
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if !amount.IsPositive() {
		return nil, errors.New("支付金额必须大于零")
	}
	
	now := time.Now()
	// 设置支付过期时间（30分钟）
	expiredAt := now.Add(30 * time.Minute)
	
	payment := &Payment{
		id:           id,
		orderID:      orderID,
		userID:       userID,
		amount:       amount,
		method:       method,
		status:       valueobject.NewPendingStatus(),
		clientIP:     clientIP,
		expiredAt:    &expiredAt,
		createdAt:    now,
		updatedAt:    now,
		domainEvents: make([]event.DomainEvent, 0),
	}
	
	// 添加支付创建事件
	payment.addDomainEvent(NewPaymentCreatedEvent(payment))
	
	return payment, nil
}

// Getters
func (p *Payment) ID() valueobject.PaymentID {
	return p.id
}

func (p *Payment) OrderID() string {
	return p.orderID
}

func (p *Payment) UserID() string {
	return p.userID
}

func (p *Payment) Amount() valueobject.Money {
	return p.amount
}

func (p *Payment) Method() valueobject.PaymentMethod {
	return p.method
}

func (p *Payment) Status() valueobject.PaymentStatus {
	return p.status
}

func (p *Payment) TransactionID() string {
	return p.transactionID
}

func (p *Payment) PaymentURL() string {
	return p.paymentURL
}

func (p *Payment) QRCode() string {
	return p.qrCode
}

func (p *Payment) ClientIP() string {
	return p.clientIP
}

func (p *Payment) ReturnURL() string {
	return p.returnURL
}

func (p *Payment) NotifyURL() string {
	return p.notifyURL
}

func (p *Payment) FailureReason() string {
	return p.failureReason
}

func (p *Payment) Metadata() string {
	return p.metadata
}

func (p *Payment) PaymentTime() *time.Time {
	return p.paymentTime
}

func (p *Payment) ExpiredAt() *time.Time {
	return p.expiredAt
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

// SetPaymentInfo 设置支付信息（第三方支付）
func (p *Payment) SetPaymentInfo(paymentURL, qrCode, transactionID string) error {
	if !p.status.IsPending() {
		return errors.New("只有待支付状态才能设置支付信息")
	}
	
	p.paymentURL = paymentURL
	p.qrCode = qrCode
	p.transactionID = transactionID
	p.updatedAt = time.Now()
	
	// 添加支付信息设置事件
	p.addDomainEvent(NewPaymentInfoSetEvent(p, paymentURL, qrCode, transactionID))
	
	return nil
}

// SetCallbackURLs 设置回调URL
func (p *Payment) SetCallbackURLs(returnURL, notifyURL string) {
	p.returnURL = returnURL
	p.notifyURL = notifyURL
	p.updatedAt = time.Now()
}

// SetMetadata 设置扩展信息
func (p *Payment) SetMetadata(metadata string) {
	p.metadata = metadata
	p.updatedAt = time.Now()
}

// StartProcessing 开始处理支付
func (p *Payment) StartProcessing() error {
	if !p.status.CanTransitionTo(valueobject.NewProcessingStatus()) {
		return errors.New("当前状态不能转为处理中")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewProcessingStatus()
	p.updatedAt = time.Now()
	
	// 添加状态变更事件
	p.addDomainEvent(NewPaymentStatusChangedEvent(p, oldStatus, p.status))
	
	return nil
}

// MarkAsSuccess 标记为支付成功
func (p *Payment) MarkAsSuccess(transactionID string) error {
	if !p.status.CanTransitionTo(valueobject.NewSuccessStatus()) {
		return errors.New("当前状态不能转为支付成功")
	}
	
	oldStatus := p.status
	now := time.Now()
	
	p.status = valueobject.NewSuccessStatus()
	p.transactionID = transactionID
	p.paymentTime = &now
	p.updatedAt = now
	
	// 添加支付成功事件
	p.addDomainEvent(NewPaymentSuccessEvent(p))
	
	return nil
}

// MarkAsFailed 标记为支付失败
func (p *Payment) MarkAsFailed(reason string) error {
	if !p.status.CanTransitionTo(valueobject.NewFailedStatus()) {
		return errors.New("当前状态不能转为支付失败")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewFailedStatus()
	p.failureReason = reason
	p.updatedAt = time.Now()
	
	// 添加支付失败事件
	p.addDomainEvent(NewPaymentFailedEvent(p, reason))
	
	return nil
}

// Cancel 取消支付
func (p *Payment) Cancel(reason string) error {
	if !p.status.CanTransitionTo(valueobject.NewCancelledStatus()) {
		return errors.New("当前状态不能取消")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewCancelledStatus()
	p.failureReason = reason
	p.updatedAt = time.Now()
	
	// 添加支付取消事件
	p.addDomainEvent(NewPaymentCancelledEvent(p, reason))
	
	return nil
}

// MarkAsRefunded 标记为已退款
func (p *Payment) MarkAsRefunded() error {
	if !p.status.CanTransitionTo(valueobject.NewRefundedStatus()) {
		return errors.New("只有成功的支付才能退款")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewRefundedStatus()
	p.updatedAt = time.Now()
	
	// 添加退款事件
	p.addDomainEvent(NewPaymentRefundedEvent(p))
	
	return nil
}

// CheckExpiration 检查是否过期
func (p *Payment) CheckExpiration() error {
	if p.expiredAt == nil {
		return nil // 没有设置过期时间
	}
	
	if time.Now().After(*p.expiredAt) && p.status.IsPending() {
		oldStatus := p.status
		p.status = valueobject.NewExpiredStatus()
		p.updatedAt = time.Now()
		
		// 添加过期事件
		p.addDomainEvent(NewPaymentExpiredEvent(p))
		
		return nil
	}
	
	return nil
}

// IsExpired 是否已过期
func (p *Payment) IsExpired() bool {
	if p.expiredAt == nil {
		return false
	}
	return time.Now().After(*p.expiredAt)
}

// CanRetry 是否可以重试
func (p *Payment) CanRetry() bool {
	return p.status.IsFailed() || p.status.IsExpired()
}

// IsOwnedBy 检查是否属于指定用户
func (p *Payment) IsOwnedBy(userID string) bool {
	return p.userID == userID
}

// GetDisplayAmount 获取显示金额
func (p *Payment) GetDisplayAmount() string {
	return p.amount.DisplayString()
}

// GetDomainEvents 获取领域事件
func (p *Payment) GetDomainEvents() []event.DomainEvent {
	return p.domainEvents
}

// ClearDomainEvents 清除领域事件
func (p *Payment) ClearDomainEvents() {
	p.domainEvents = make([]event.DomainEvent, 0)
}

// addDomainEvent 添加领域事件
func (p *Payment) addDomainEvent(domainEvent event.DomainEvent) {
	p.domainEvents = append(p.domainEvents, domainEvent)
}
