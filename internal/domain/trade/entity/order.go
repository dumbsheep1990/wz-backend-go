package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// OrderItem 订单项
type OrderItem struct {
	productID valueobject.ProductID
	name      string
	price     valueobject.Money
	quantity  valueobject.Quantity
}

// NewOrderItem 创建一个新的订单项
func NewOrderItem(
	productID valueobject.ProductID,
	name string,
	price valueobject.Money,
	quantity valueobject.Quantity,
) (*OrderItem, error) {
	if productID.IsEmpty() {
		return nil, errors.New("商品ID不能为空")
	}
	if name == "" {
		return nil, errors.New("商品名称不能为空")
	}

	return &OrderItem{
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
	}, nil
}

// ProductID 获取商品ID
func (i *OrderItem) ProductID() valueobject.ProductID {
	return i.productID
}

// Name 获取商品名称
func (i *OrderItem) Name() string {
	return i.name
}

// Price 获取商品单价
func (i *OrderItem) Price() valueobject.Money {
	return i.price
}

// Quantity 获取商品数量
func (i *OrderItem) Quantity() valueobject.Quantity {
	return i.quantity
}

// Subtotal 计算小计金额
func (i *OrderItem) Subtotal() (valueobject.Money, error) {
	return i.price.Multiply(i.quantity.Value())
}

// Order 订单实体
type Order struct {
	id            valueobject.OrderID
	userID        valueobject.UserID
	items         []*OrderItem
	totalAmount   valueobject.Money
	status        valueobject.OrderStatus
	shippingAddress valueobject.Address
	paymentID     valueobject.PaymentID
	createdAt     time.Time
	updatedAt     time.Time
	
	domainEvents []event.DomainEvent
}

// NewOrder 创建一个新订单
func NewOrder(
	userID valueobject.UserID,
	items []*OrderItem,
	shippingAddress valueobject.Address,
) (*Order, error) {
	if userID.IsEmpty() {
		return nil, errors.New("用户ID不能为空")
	}
	
	if len(items) == 0 {
		return nil, errors.New("订单项不能为空")
	}
	
	// 计算订单总金额
	var totalAmount valueobject.Money
	var err error
	
	// 假设第一个订单项的货币类型为订单的货币类型
	totalAmount, err = valueobject.NewMoney(0, items[0].Price().Currency())
	if err != nil {
		return nil, err
	}
	
	for _, item := range items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return nil, err
		}
		
		totalAmount, err = totalAmount.Add(subtotal)
		if err != nil {
			return nil, err
		}
	}
	
	now := time.Now()
	order := &Order{
		id:            valueobject.NewOrderID(uuid.New().String()),
		userID:        userID,
		items:         items,
		totalAmount:   totalAmount,
		status:        valueobject.StatusPending,
		shippingAddress: shippingAddress,
		createdAt:     now,
		updatedAt:     now,
		domainEvents:  []event.DomainEvent{},
	}
	
	// 添加订单创建事件
	order.addDomainEvent(NewOrderCreatedEvent(order))
	
	return order, nil
}

// ID 获取订单ID
func (o *Order) ID() valueobject.OrderID {
	return o.id
}

// UserID 获取用户ID
func (o *Order) UserID() valueobject.UserID {
	return o.userID
}

// Items 获取订单项
func (o *Order) Items() []*OrderItem {
	return o.items
}

// TotalAmount 获取订单总金额
func (o *Order) TotalAmount() valueobject.Money {
	return o.totalAmount
}

// Status 获取订单状态
func (o *Order) Status() valueobject.OrderStatus {
	return o.status
}

// ShippingAddress 获取配送地址
func (o *Order) ShippingAddress() valueobject.Address {
	return o.shippingAddress
}

// PaymentID 获取支付ID
func (o *Order) PaymentID() valueobject.PaymentID {
	return o.paymentID
}

// CreatedAt 获取创建时间
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt 获取更新时间
func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// BindPayment 绑定支付信息
func (o *Order) BindPayment(paymentID valueobject.PaymentID) error {
	if paymentID.IsEmpty() {
		return errors.New("支付ID不能为空")
	}
	
	o.paymentID = paymentID
	o.updatedAt = time.Now()
	
	return nil
}

// Pay 支付订单
func (o *Order) Pay() error {
	if o.status != valueobject.StatusPending {
		return errors.New("只有待付款状态的订单才能支付")
	}
	
	if o.paymentID.IsEmpty() {
		return errors.New("订单未绑定支付信息")
	}
	
	o.status = valueobject.StatusPaid
	o.updatedAt = time.Now()
	
	// 添加订单支付事件
	o.addDomainEvent(NewOrderPaidEvent(o))
	
	return nil
}

// Ship 发货
func (o *Order) Ship() error {
	if o.status != valueobject.StatusPaid {
		return errors.New("只有已付款状态的订单才能发货")
	}
	
	o.status = valueobject.StatusShipping
	o.updatedAt = time.Now()
	
	// 添加订单发货事件
	o.addDomainEvent(NewOrderShippedEvent(o))
	
	return nil
}

// Deliver 送达
func (o *Order) Deliver() error {
	if o.status != valueobject.StatusShipping {
		return errors.New("只有配送中状态的订单才能标记为已送达")
	}
	
	o.status = valueobject.StatusDelivered
	o.updatedAt = time.Now()
	
	// 添加订单送达事件
	o.addDomainEvent(NewOrderDeliveredEvent(o))
	
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	if o.status != valueobject.StatusDelivered {
		return errors.New("只有已送达状态的订单才能标记为已完成")
	}
	
	o.status = valueobject.StatusCompleted
	o.updatedAt = time.Now()
	
	// 添加订单完成事件
	o.addDomainEvent(NewOrderCompletedEvent(o))
	
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.status != valueobject.StatusPending {
		return errors.New("只有待付款状态的订单才能取消")
	}
	
	o.status = valueobject.StatusCancelled
	o.updatedAt = time.Now()
	
	// 添加订单取消事件
	o.addDomainEvent(NewOrderCancelledEvent(o))
	
	return nil
}

// RequestRefund 申请退款
func (o *Order) RequestRefund() error {
	if o.status != valueobject.StatusPaid && o.status != valueobject.StatusShipping && o.status != valueobject.StatusDelivered {
		return errors.New("只有已付款、配送中或已送达状态的订单才能申请退款")
	}
	
	o.status = valueobject.StatusRefunding
	o.updatedAt = time.Now()
	
	// 添加订单申请退款事件
	o.addDomainEvent(NewOrderRefundRequestedEvent(o))
	
	return nil
}

// Refund 确认退款
func (o *Order) Refund() error {
	if o.status != valueobject.StatusRefunding {
		return errors.New("只有退款中状态的订单才能确认退款")
	}
	
	o.status = valueobject.StatusRefunded
	o.updatedAt = time.Now()
	
	// 添加订单退款事件
	o.addDomainEvent(NewOrderRefundedEvent(o))
	
	return nil
}

// AddItem 添加订单项
func (o *Order) AddItem(item *OrderItem) error {
	if o.status != valueobject.StatusPending {
		return errors.New("只有待付款状态的订单才能修改订单项")
	}
	
	// 添加订单项
	o.items = append(o.items, item)
	
	// 重新计算订单总金额
	var totalAmount valueobject.Money
	var err error
	
	totalAmount, err = valueobject.NewMoney(0, o.totalAmount.Currency())
	if err != nil {
		return err
	}
	
	for _, item := range o.items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return err
		}
		
		totalAmount, err = totalAmount.Add(subtotal)
		if err != nil {
			return err
		}
	}
	
	o.totalAmount = totalAmount
	o.updatedAt = time.Now()
	
	return nil
}

// UpdateShippingAddress 更新配送地址
func (o *Order) UpdateShippingAddress(address valueobject.Address) error {
	if o.status != valueobject.StatusPending {
		return errors.New("只有待付款状态的订单才能修改配送地址")
	}
	
	o.shippingAddress = address
	o.updatedAt = time.Now()
	
	return nil
}

// 添加领域事件
func (o *Order) addDomainEvent(event event.DomainEvent) {
	o.domainEvents = append(o.domainEvents, event)
}

// GetDomainEvents 获取所有领域事件
func (o *Order) GetDomainEvents() []event.DomainEvent {
	return o.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (o *Order) ClearDomainEvents() {
	o.domainEvents = []event.DomainEvent{}
}
