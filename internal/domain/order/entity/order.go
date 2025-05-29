package entity

import (
	"errors"
	"time"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
	productvo "wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"

	"github.com/google/uuid"
)

// Order 订单实体，是订单聚合的根
type Order struct {
	id              ordervo.OrderID
	orderNumber     ordervo.OrderNumber
	customerID      uservo.UserID
	status          ordervo.OrderStatus
	items           []*OrderItem
	discounts       []*OrderDiscount
	shippingAddress ordervo.Address
	billingAddress  ordervo.Address
	paymentMethod   ordervo.PaymentMethod
	shippingMethod  ordervo.ShippingMethod

	subtotal       ordervo.Money // 商品总金额（未含运费和折扣）
	shippingFee    ordervo.Money // 运费
	tax            ordervo.Money // 税费
	discountAmount ordervo.Money // 折扣总金额
	totalAmount    ordervo.Money // 订单总金额

	note           string
	trackingNumber string

	paidAt      *time.Time
	shippedAt   *time.Time
	deliveredAt *time.Time
	completedAt *time.Time
	cancelledAt *time.Time
	refundedAt  *time.Time

	createdAt time.Time
	updatedAt time.Time
}

// NewOrder 创建订单
func NewOrder(
	customerID uservo.UserID,
	shippingAddress ordervo.Address,
	billingAddress ordervo.Address,
	shippingMethod ordervo.ShippingMethod,
) (*Order, error) {
	// 创建订单ID
	orderID := ordervo.NewOrderID(uuid.New().String())

	// 生成订单编号
	orderNumber := ordervo.GenerateOrderNumber()

	// 初始化金额
	zeroMoney, _ := ordervo.NewMoney(0, "CNY")

	// 创建初始状态为已创建的订单
	initialStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusCreated))

	now := time.Now()

	order := &Order{
		id:              orderID,
		orderNumber:     orderNumber,
		customerID:      customerID,
		status:          initialStatus,
		items:           []*OrderItem{},
		discounts:       []*OrderDiscount{},
		shippingAddress: shippingAddress,
		billingAddress:  billingAddress,
		shippingMethod:  shippingMethod,
		paymentMethod:   ordervo.PaymentMethodUnknown,

		subtotal:       zeroMoney,
		shippingFee:    zeroMoney,
		tax:            zeroMoney,
		discountAmount: zeroMoney,
		totalAmount:    zeroMoney,

		note:           "",
		trackingNumber: "",

		createdAt: now,
		updatedAt: now,
	}

	return order, nil
}

// ReconstructOrder 从存储中重建订单实体
func ReconstructOrder(
	id ordervo.OrderID,
	orderNumber ordervo.OrderNumber,
	customerID uservo.UserID,
	status ordervo.OrderStatus,
	items []*OrderItem,
	discounts []*OrderDiscount,
	shippingAddress ordervo.Address,
	billingAddress ordervo.Address,
	paymentMethod ordervo.PaymentMethod,
	shippingMethod ordervo.ShippingMethod,
	subtotal ordervo.Money,
	shippingFee ordervo.Money,
	tax ordervo.Money,
	discountAmount ordervo.Money,
	totalAmount ordervo.Money,
	note string,
	trackingNumber string,
	paidAt *time.Time,
	shippedAt *time.Time,
	deliveredAt *time.Time,
	completedAt *time.Time,
	cancelledAt *time.Time,
	refundedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Order {
	return &Order{
		id:              id,
		orderNumber:     orderNumber,
		customerID:      customerID,
		status:          status,
		items:           items,
		discounts:       discounts,
		shippingAddress: shippingAddress,
		billingAddress:  billingAddress,
		paymentMethod:   paymentMethod,
		shippingMethod:  shippingMethod,

		subtotal:       subtotal,
		shippingFee:    shippingFee,
		tax:            tax,
		discountAmount: discountAmount,
		totalAmount:    totalAmount,

		note:           note,
		trackingNumber: trackingNumber,

		paidAt:      paidAt,
		shippedAt:   shippedAt,
		deliveredAt: deliveredAt,
		completedAt: completedAt,
		cancelledAt: cancelledAt,
		refundedAt:  refundedAt,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID 获取订单ID
func (o *Order) ID() ordervo.OrderID {
	return o.id
}

// OrderNumber 获取订单编号
func (o *Order) OrderNumber() ordervo.OrderNumber {
	return o.orderNumber
}

// CustomerID 获取客户ID
func (o *Order) CustomerID() uservo.UserID {
	return o.customerID
}

// Status 获取订单状态
func (o *Order) Status() ordervo.OrderStatus {
	return o.status
}

// Items 获取订单项
func (o *Order) Items() []*OrderItem {
	return o.items
}

// Discounts 获取订单折扣
func (o *Order) Discounts() []*OrderDiscount {
	return o.discounts
}

// ShippingAddress 获取配送地址
func (o *Order) ShippingAddress() ordervo.Address {
	return o.shippingAddress
}

// BillingAddress 获取账单地址
func (o *Order) BillingAddress() ordervo.Address {
	return o.billingAddress
}

// PaymentMethod 获取支付方式
func (o *Order) PaymentMethod() ordervo.PaymentMethod {
	return o.paymentMethod
}

// ShippingMethod 获取配送方式
func (o *Order) ShippingMethod() ordervo.ShippingMethod {
	return o.shippingMethod
}

// Subtotal 获取商品总金额
func (o *Order) Subtotal() ordervo.Money {
	return o.subtotal
}

// ShippingFee 获取运费
func (o *Order) ShippingFee() ordervo.Money {
	return o.shippingFee
}

// Tax 获取税费
func (o *Order) Tax() ordervo.Money {
	return o.tax
}

// DiscountAmount 获取折扣总金额
func (o *Order) DiscountAmount() ordervo.Money {
	return o.discountAmount
}

// TotalAmount 获取订单总金额
func (o *Order) TotalAmount() ordervo.Money {
	return o.totalAmount
}

// Note 获取订单备注
func (o *Order) Note() string {
	return o.note
}

// TrackingNumber 获取物流单号
func (o *Order) TrackingNumber() string {
	return o.trackingNumber
}

// PaidAt 获取支付时间
func (o *Order) PaidAt() *time.Time {
	return o.paidAt
}

// ShippedAt 获取发货时间
func (o *Order) ShippedAt() *time.Time {
	return o.shippedAt
}

// DeliveredAt 获取送达时间
func (o *Order) DeliveredAt() *time.Time {
	return o.deliveredAt
}

// CompletedAt 获取完成时间
func (o *Order) CompletedAt() *time.Time {
	return o.completedAt
}

// CancelledAt 获取取消时间
func (o *Order) CancelledAt() *time.Time {
	return o.cancelledAt
}

// RefundedAt 获取退款时间
func (o *Order) RefundedAt() *time.Time {
	return o.refundedAt
}

// CreatedAt 获取创建时间
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt 获取更新时间
func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// AddItem 添加订单项
func (o *Order) AddItem(
	productID productvo.ProductID,
	productName string,
	productSKU string,
	quantity int32,
	unitPrice ordervo.Money,
	attributes map[string]string,
) error {
	// 检查订单状态是否允许添加商品
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许添加商品")
	}

	// 生成订单项ID
	itemID := ordervo.NewOrderItemID(uuid.New().String())

	// 创建新的订单项
	item, err := NewOrderItem(
		itemID,
		productID,
		productName,
		productSKU,
		quantity,
		unitPrice,
		attributes,
	)
	if err != nil {
		return err
	}

	// 添加到订单项列表
	o.items = append(o.items, item)

	// 重新计算订单金额
	err = o.recalculateAmounts()
	if err != nil {
		// 如果计算失败，移除刚添加的订单项
		o.items = o.items[:len(o.items)-1]
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// RemoveItem 移除订单项
func (o *Order) RemoveItem(itemID ordervo.OrderItemID) error {
	// 检查订单状态是否允许移除商品
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许移除商品")
	}

	// 查找订单项
	index := -1
	for i, item := range o.items {
		if item.ID().Value() == itemID.Value() {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("订单项不存在")
	}

	// 移除订单项
	o.items = append(o.items[:index], o.items[index+1:]...)

	// 重新计算订单金额
	err := o.recalculateAmounts()
	if err != nil {
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// UpdateItemQuantity 更新订单项数量
func (o *Order) UpdateItemQuantity(itemID ordervo.OrderItemID, quantity int32) error {
	// 检查订单状态是否允许更新商品数量
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许更新商品数量")
	}

	// 查找订单项
	var targetItem *OrderItem
	for _, item := range o.items {
		if item.ID().Value() == itemID.Value() {
			targetItem = item
			break
		}
	}

	if targetItem == nil {
		return errors.New("订单项不存在")
	}

	// 更新数量
	err := targetItem.UpdateQuantity(quantity)
	if err != nil {
		return err
	}

	// 重新计算订单金额
	err = o.recalculateAmounts()
	if err != nil {
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// AddDiscount 添加折扣
func (o *Order) AddDiscount(discount *OrderDiscount) error {
	// 检查订单状态是否允许添加折扣
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许添加折扣")
	}

	// 检查折扣是否已存在
	for _, d := range o.discounts {
		if d.ID() == discount.ID() {
			return errors.New("折扣已存在")
		}
	}

	// 添加折扣
	o.discounts = append(o.discounts, discount)

	// 重新计算订单金额
	err := o.recalculateAmounts()
	if err != nil {
		// 如果计算失败，移除刚添加的折扣
		o.discounts = o.discounts[:len(o.discounts)-1]
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// RemoveDiscount 移除折扣
func (o *Order) RemoveDiscount(discountID string) error {
	// 检查订单状态是否允许移除折扣
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许移除折扣")
	}

	// 查找折扣
	index := -1
	for i, discount := range o.discounts {
		if discount.ID() == discountID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("折扣不存在")
	}

	// 移除折扣
	o.discounts = append(o.discounts[:index], o.discounts[index+1:]...)

	// 重新计算订单金额
	err := o.recalculateAmounts()
	if err != nil {
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// SetShippingFee 设置运费
func (o *Order) SetShippingFee(fee ordervo.Money) error {
	// 检查订单状态是否允许修改运费
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许修改运费")
	}

	o.shippingFee = fee

	// 重新计算订单金额
	err := o.recalculateAmounts()
	if err != nil {
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// SetTax 设置税费
func (o *Order) SetTax(tax ordervo.Money) error {
	// 检查订单状态是否允许修改税费
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许修改税费")
	}

	o.tax = tax

	// 重新计算订单金额
	err := o.recalculateAmounts()
	if err != nil {
		return err
	}

	o.updatedAt = time.Now()
	return nil
}

// SetPaymentMethod 设置支付方式
func (o *Order) SetPaymentMethod(method ordervo.PaymentMethod) error {
	// 检查订单状态是否允许修改支付方式
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许修改支付方式")
	}

	o.paymentMethod = method
	o.updatedAt = time.Now()
	return nil
}

// SetShippingMethod 设置配送方式
func (o *Order) SetShippingMethod(method ordervo.ShippingMethod) error {
	// 检查订单状态是否允许修改配送方式
	if !o.status.CanTransitionTo(ordervo.OrderStatusPending) {
		return errors.New("当前订单状态不允许修改配送方式")
	}

	o.shippingMethod = method
	o.updatedAt = time.Now()
	return nil
}

// SetNote 设置订单备注
func (o *Order) SetNote(note string) {
	o.note = note
	o.updatedAt = time.Now()
}

// SetTrackingNumber 设置物流单号
func (o *Order) SetTrackingNumber(trackingNumber string) error {
	// 只有已发货的订单才能设置物流单号
	if o.status != ordervo.OrderStatusShipped && o.status != ordervo.OrderStatusDelivered {
		return errors.New("只有已发货的订单才能设置物流单号")
	}

	o.trackingNumber = trackingNumber
	o.updatedAt = time.Now()
	return nil
}

// recalculateAmounts 重新计算订单金额
func (o *Order) recalculateAmounts() error {
	// 计算商品总金额
	zeroMoney, _ := ordervo.NewMoney(0, "CNY")
	subtotal := zeroMoney

	for _, item := range o.items {
		var err error
		subtotal, err = subtotal.Add(item.TotalPrice())
		if err != nil {
			return err
		}
	}

	o.subtotal = subtotal

	// 计算折扣总金额
	discountTotal := zeroMoney
	for _, discount := range o.discounts {
		// 计算折扣金额
		discountAmount, err := discount.CalculateDiscount(o.subtotal)
		if err != nil {
			continue // 跳过无效折扣
		}

		// 累加折扣金额
		discountTotal, err = discountTotal.Add(discountAmount)
		if err != nil {
			return err
		}
	}

	// 确保折扣总额不超过商品总额
	if discountTotal.Amount() > subtotal.Amount() {
		discountTotal, _ = ordervo.NewMoney(subtotal.Amount(), subtotal.Currency())
	}

	o.discountAmount = discountTotal

	// 计算订单总金额：商品总金额 + 运费 + 税费 - 折扣总金额
	totalAmount := subtotal

	// 添加运费
	totalAmount, _ = totalAmount.Add(o.shippingFee)

	// 添加税费
	totalAmount, _ = totalAmount.Add(o.tax)

	// 减去折扣
	totalAmount, _ = totalAmount.Subtract(discountTotal)

	o.totalAmount = totalAmount

	return nil
}

// SubmitOrder 提交订单
func (o *Order) SubmitOrder() error {
	// 检查订单状态
	if o.status != ordervo.OrderStatusCreated {
		return errors.New("只有处于已创建状态的订单才能提交")
	}

	// 检查订单是否有商品
	if len(o.items) == 0 {
		return errors.New("订单中没有商品")
	}

	// 更新状态为待支付
	pendingStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusPending))
	o.status = pendingStatus

	o.updatedAt = time.Now()
	return nil
}

// Pay 支付订单
func (o *Order) Pay(paymentMethod ordervo.PaymentMethod) error {
	// 检查订单状态
	if !o.status.CanPay() {
		return errors.New("当前订单状态不允许支付")
	}

	// 设置支付方式
	o.paymentMethod = paymentMethod

	// 更新状态为已支付
	paidStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusPaid))
	o.status = paidStatus

	// 记录支付时间
	now := time.Now()
	o.paidAt = &now

	o.updatedAt = now
	return nil
}

// Ship 发货
func (o *Order) Ship(trackingNumber string) error {
	// 检查订单状态
	if !o.status.CanShip() {
		return errors.New("当前订单状态不允许发货")
	}

	// 设置物流单号
	o.trackingNumber = trackingNumber

	// 更新状态为已发货
	shippedStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusShipped))
	o.status = shippedStatus

	// 记录发货时间
	now := time.Now()
	o.shippedAt = &now

	o.updatedAt = now
	return nil
}

// Deliver 标记为已送达
func (o *Order) Deliver() error {
	// 检查订单状态
	if !o.status.CanDeliver() {
		return errors.New("当前订单状态不允许标记为已送达")
	}

	// 更新状态为已送达
	deliveredStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusDelivered))
	o.status = deliveredStatus

	// 记录送达时间
	now := time.Now()
	o.deliveredAt = &now

	o.updatedAt = now
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	// 检查订单状态
	if !o.status.CanComplete() {
		return errors.New("当前订单状态不允许标记为已完成")
	}

	// 更新状态为已完成
	completedStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusCompleted))
	o.status = completedStatus

	// 记录完成时间
	now := time.Now()
	o.completedAt = &now

	o.updatedAt = now
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	// 检查订单状态
	if !o.status.CanCancel() {
		return errors.New("当前订单状态不允许取消")
	}

	// 更新状态为已取消
	cancelledStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusCancelled))
	o.status = cancelledStatus

	// 记录取消时间
	now := time.Now()
	o.cancelledAt = &now

	o.updatedAt = now
	return nil
}

// RequestRefund 申请退款
func (o *Order) RequestRefund() error {
	// 检查订单状态
	if !o.status.CanRefund() {
		return errors.New("当前订单状态不允许申请退款")
	}

	// 检查支付方式是否支持退款
	if !o.paymentMethod.SupportRefund() {
		return errors.New("当前支付方式不支持退款")
	}

	// 更新状态为退款中
	refundingStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusRefunding))
	o.status = refundingStatus

	o.updatedAt = time.Now()
	return nil
}

// Refund 退款完成
func (o *Order) Refund() error {
	// 检查订单状态
	if o.status != ordervo.OrderStatusRefunding {
		return errors.New("只有处于退款中状态的订单才能标记为已退款")
	}

	// 更新状态为已退款
	refundedStatus, _ := ordervo.NewOrderStatus(int32(ordervo.OrderStatusRefunded))
	o.status = refundedStatus

	// 记录退款时间
	now := time.Now()
	o.refundedAt = &now

	o.updatedAt = now
	return nil
}

// CanModify 检查订单是否可以修改
func (o *Order) CanModify() bool {
	return o.status == ordervo.OrderStatusCreated || o.status == ordervo.OrderStatusPending
}

// IsActive 检查订单是否处于活跃状态（非终态）
func (o *Order) IsActive() bool {
	return o.status.IsActive()
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.status == ordervo.OrderStatusPaid ||
		o.status == ordervo.OrderStatusShipped ||
		o.status == ordervo.OrderStatusDelivered ||
		o.status == ordervo.OrderStatusCompleted ||
		o.status == ordervo.OrderStatusRefunding ||
		o.status == ordervo.OrderStatusRefunded
}

// IsShipped 检查订单是否已发货
func (o *Order) IsShipped() bool {
	return o.status == ordervo.OrderStatusShipped ||
		o.status == ordervo.OrderStatusDelivered ||
		o.status == ordervo.OrderStatusCompleted
}

// IsCompleted 检查订单是否已完成
func (o *Order) IsCompleted() bool {
	return o.status == ordervo.OrderStatusCompleted
}

// IsCancelled 检查订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.status == ordervo.OrderStatusCancelled
}

// IsRefunded 检查订单是否已退款
func (o *Order) IsRefunded() bool {
	return o.status == ordervo.OrderStatusRefunded
}
