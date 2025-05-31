package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/order/entity"
	"wz-backend-go/internal/domain/order/event"
	"wz-backend-go/internal/domain/order/repository"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
	productvo "wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// OrderDomainService 订单领域服务
type OrderDomainService struct {
	orderRepository repository.OrderRepository
	eventPublisher  repository.EventPublisher
}

// NewOrderDomainService 创建订单领域服务
func NewOrderDomainService(
	orderRepository repository.OrderRepository,
	eventPublisher repository.EventPublisher,
) *OrderDomainService {
	return &OrderDomainService{
		orderRepository: orderRepository,
		eventPublisher:  eventPublisher,
	}
}

// CreateOrder 创建订单
func (s *OrderDomainService) CreateOrder(
	customerID uservo.UserID,
	shippingAddress ordervo.Address,
	billingAddress ordervo.Address,
	shippingMethod ordervo.ShippingMethod,
) (*entity.Order, error) {
	// 创建订单
	order, err := entity.NewOrder(customerID, shippingAddress, billingAddress, shippingMethod)
	if err != nil {
		return nil, err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return nil, err
	}

	// 发布订单创建事件
	s.eventPublisher.Publish(event.NewOrderCreatedEvent(order))

	return order, nil
}

// AddOrderItem 添加订单项
func (s *OrderDomainService) AddOrderItem(
	orderID ordervo.OrderID,
	productID productvo.ProductID,
	productName string,
	productSKU string,
	quantity int32,
	unitPrice ordervo.Money,
	attributes map[string]string,
) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 添加订单项
	err = order.AddItem(productID, productName, productSKU, quantity, unitPrice, attributes)
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 查找刚添加的订单项
	var addedItem *entity.OrderItem
	for _, item := range order.Items() {
		if item.ProductID().Value() == productID.Value() {
			addedItem = item
			break
		}
	}

	if addedItem != nil {
		// 发布订单项添加事件
		s.eventPublisher.Publish(event.NewOrderItemAddedEvent(order, addedItem))
	}

	return nil
}

// RemoveOrderItem 移除订单项
func (s *OrderDomainService) RemoveOrderItem(orderID ordervo.OrderID, itemID ordervo.OrderItemID) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 查找要移除的订单项
	var removedItem *entity.OrderItem
	for _, item := range order.Items() {
		if item.ID().Value() == itemID.Value() {
			removedItem = item
			break
		}
	}

	if removedItem == nil {
		return errors.New("订单项不存在")
	}

	// 移除订单项
	err = order.RemoveItem(itemID)
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单项移除事件
	s.eventPublisher.Publish(event.NewOrderItemRemovedEvent(order, removedItem))

	return nil
}

// UpdateOrderItemQuantity 更新订单项数量
func (s *OrderDomainService) UpdateOrderItemQuantity(orderID ordervo.OrderID, itemID ordervo.OrderItemID, quantity int32) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 更新订单项数量
	err = order.UpdateItemQuantity(itemID, quantity)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// AddDiscount 添加折扣
func (s *OrderDomainService) AddDiscount(
	orderID ordervo.OrderID,
	discountName string,
	discountType entity.DiscountType,
	discountValue float64,
	discountCode string,
	description string,
	minOrderValue ordervo.Money,
	startTime time.Time,
	endTime time.Time,
) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 创建折扣
	discount, err := entity.NewOrderDiscount(
		uuid.New().String(),
		discountName,
		discountType,
		discountValue,
		discountCode,
		description,
		minOrderValue,
		startTime,
		endTime,
	)
	if err != nil {
		return err
	}

	// 添加折扣
	err = order.AddDiscount(discount)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// RemoveDiscount 移除折扣
func (s *OrderDomainService) RemoveDiscount(orderID ordervo.OrderID, discountID string) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 移除折扣
	err = order.RemoveDiscount(discountID)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// SetShippingFee 设置运费
func (s *OrderDomainService) SetShippingFee(orderID ordervo.OrderID, fee ordervo.Money) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 设置运费
	err = order.SetShippingFee(fee)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// SetShippingMethod 设置配送方式
func (s *OrderDomainService) SetShippingMethod(orderID ordervo.OrderID, method ordervo.ShippingMethod) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 设置配送方式
	err = order.SetShippingMethod(method)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// SetPaymentMethod 设置支付方式
func (s *OrderDomainService) SetPaymentMethod(orderID ordervo.OrderID, method ordervo.PaymentMethod) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 设置支付方式
	err = order.SetPaymentMethod(method)
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// SubmitOrder 提交订单
func (s *OrderDomainService) SubmitOrder(orderID ordervo.OrderID) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 提交订单
	err = order.SubmitOrder()
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// PayOrder 支付订单
func (s *OrderDomainService) PayOrder(orderID ordervo.OrderID, paymentMethod ordervo.PaymentMethod) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 支付订单
	err = order.Pay(paymentMethod)
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单支付事件
	s.eventPublisher.Publish(event.NewOrderPaidEvent(order))

	return nil
}

// ShipOrder 发货
func (s *OrderDomainService) ShipOrder(orderID ordervo.OrderID, trackingNumber string) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 发货
	err = order.Ship(trackingNumber)
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单发货事件
	s.eventPublisher.Publish(event.NewOrderShippedEvent(order))

	return nil
}

// DeliverOrder 标记为已送达
func (s *OrderDomainService) DeliverOrder(orderID ordervo.OrderID) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 标记为已送达
	err = order.Deliver()
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单送达事件
	s.eventPublisher.Publish(event.NewOrderDeliveredEvent(order))

	return nil
}

// CompleteOrder 完成订单
func (s *OrderDomainService) CompleteOrder(orderID ordervo.OrderID) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 完成订单
	err = order.Complete()
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单完成事件
	s.eventPublisher.Publish(event.NewOrderCompletedEvent(order))

	return nil
}

// CancelOrder 取消订单
func (s *OrderDomainService) CancelOrder(orderID ordervo.OrderID, reason string) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 取消订单
	err = order.Cancel()
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单取消事件
	s.eventPublisher.Publish(event.NewOrderCancelledEvent(order, reason))

	return nil
}

// RequestRefund 申请退款
func (s *OrderDomainService) RequestRefund(orderID ordervo.OrderID) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 申请退款
	err = order.RequestRefund()
	if err != nil {
		return err
	}

	// 保存订单
	return s.orderRepository.Save(order)
}

// RefundOrder 退款完成
func (s *OrderDomainService) RefundOrder(orderID ordervo.OrderID, reason string) error {
	// 查找订单
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("订单不存在")
	}

	// 退款完成
	err = order.Refund()
	if err != nil {
		return err
	}

	// 保存订单
	err = s.orderRepository.Save(order)
	if err != nil {
		return err
	}

	// 发布订单退款事件
	s.eventPublisher.Publish(event.NewOrderRefundedEvent(order, reason))

	return nil
}

// GetOrderByID 根据ID获取订单
func (s *OrderDomainService) GetOrderByID(orderID ordervo.OrderID) (*entity.Order, error) {
	return s.orderRepository.FindByID(orderID)
}

// GetOrderByOrderNumber 根据订单号获取订单
func (s *OrderDomainService) GetOrderByOrderNumber(orderNumber ordervo.OrderNumber) (*entity.Order, error) {
	return s.orderRepository.FindByOrderNumber(orderNumber)
}

// GetOrdersByCustomerID 获取客户的订单
func (s *OrderDomainService) GetOrdersByCustomerID(customerID uservo.UserID, page, pageSize int) ([]*entity.Order, int64, error) {
	return s.orderRepository.FindByCustomerID(customerID, page, pageSize)
}

// GetOrdersByStatus 根据状态获取订单
func (s *OrderDomainService) GetOrdersByStatus(status ordervo.OrderStatus, page, pageSize int) ([]*entity.Order, int64, error) {
	return s.orderRepository.FindByStatus(status, page, pageSize)
}

// GetOrders 分页获取所有订单
func (s *OrderDomainService) GetOrders(page, pageSize int) ([]*entity.Order, int64, error) {
	return s.orderRepository.FindAll(page, pageSize)
}

// GetActiveOrders 获取活跃订单
func (s *OrderDomainService) GetActiveOrders(customerID uservo.UserID, page, pageSize int) ([]*entity.Order, int64, error) {
	return s.orderRepository.FindActiveOrders(customerID, page, pageSize)
}

// GetRecentOrders 获取最近的订单
func (s *OrderDomainService) GetRecentOrders(limit int) ([]*entity.Order, error) {
	return s.orderRepository.FindRecentOrders(limit)
}

// SearchOrders 搜索订单
func (s *OrderDomainService) SearchOrders(keyword string, page, pageSize int) ([]*entity.Order, int64, error) {
	return s.orderRepository.Search(keyword, page, pageSize)
}

// DeleteOrder 删除订单
func (s *OrderDomainService) DeleteOrder(orderID ordervo.OrderID) error {
	return s.orderRepository.Delete(orderID)
}

// CanUserAccessOrder 验证用户是否可以访问订单
func (s *OrderDomainService) CanUserAccessOrder(userID uservo.UserID, order *entity.Order) bool {
	// 用户只能访问自己的订单
	return order.CustomerID().Value() == userID.Value()
}

// CanModifyOrder 验证订单是否可以修改
func (s *OrderDomainService) CanModifyOrder(order *entity.Order) bool {
	// 只有创建状态的订单可以修改
	return order.CanModify()
}

// CanCancelOrder 验证订单是否可以取消
func (s *OrderDomainService) CanCancelOrder(userID uservo.UserID, order *entity.Order) bool {
	// 验证用户权限
	if !s.CanUserAccessOrder(userID, order) {
		return false
	}

	// 已支付但未发货的订单可以取消
	return order.IsPaid() && !order.IsShipped() && !order.IsCompleted() && !order.IsCancelled()
}

// CanRefundOrder 验证订单是否可以退款
func (s *OrderDomainService) CanRefundOrder(userID uservo.UserID, order *entity.Order) bool {
	// 验证用户权限
	if !s.CanUserAccessOrder(userID, order) {
		return false
	}

	// 已支付且未退款的订单可以申请退款
	if !order.IsPaid() || order.IsRefunded() {
		return false
	}

	// 已取消的订单不能退款
	if order.IsCancelled() {
		return false
	}

	// 检查退款时限（完成后30天内可以退款）
	if order.IsCompleted() && order.CompletedAt() != nil {
		refundDeadline := order.CompletedAt().Add(30 * 24 * time.Hour)
		if time.Now().After(refundDeadline) {
			return false
		}
	}

	return true
}

// ValidateOrderCreation 验证订单创建
func (s *OrderDomainService) ValidateOrderCreation(ctx context.Context, order *entity.Order) error {
	// 验证订单项不为空
	if len(order.Items()) == 0 {
		return ordervo.ErrOrderItemsEmpty
	}

	// 验证订单金额
	if order.TotalAmount().Amount() <= 0 {
		return ordervo.ErrInvalidOrderAmount
	}

	// 验证地址信息
	if err := s.validateAddress(order.ShippingAddress()); err != nil {
		return err
	}

	if err := s.validateAddress(order.BillingAddress()); err != nil {
		return err
	}

	return nil
}

// ValidateOrderPayment 验证订单支付
func (s *OrderDomainService) ValidateOrderPayment(order *entity.Order, paymentMethod ordervo.PaymentMethod) error {
	// 验证订单状态
	if order.IsPaid() {
		return ordervo.ErrOrderAlreadyPaid
	}

	if order.IsCancelled() {
		return ordervo.ErrOrderCancelled
	}

	// 验证支付方式
	if !paymentMethod.IsValid() {
		return ordervo.ErrInvalidPaymentMethod
	}

	// 验证订单金额
	if order.TotalAmount().Amount() <= 0 {
		return ordervo.ErrInvalidOrderAmount
	}

	return nil
}

// ValidateOrderShipment 验证订单发货
func (s *OrderDomainService) ValidateOrderShipment(order *entity.Order, trackingNumber string) error {
	// 验证订单已支付
	if !order.IsPaid() {
		return ordervo.ErrOrderNotPaid
	}

	// 验证订单未取消
	if order.IsCancelled() {
		return ordervo.ErrOrderCancelled
	}

	// 验证订单未发货
	if order.IsShipped() {
		return ordervo.ErrOrderAlreadyShipped
	}

	// 验证物流单号
	if trackingNumber == "" {
		return ordervo.ErrInvalidTrackingNumber
	}

	return nil
}

// CalculateRefundAmount 计算退款金额
func (s *OrderDomainService) CalculateRefundAmount(order *entity.Order) (ordervo.Money, error) {
	// 如果订单未发货，全额退款
	if !order.IsShipped() {
		return order.TotalAmount(), nil
	}

	// 如果订单已发货但未完成，扣除配送费
	if order.IsShipped() && !order.IsCompleted() {
		refundAmount := order.TotalAmount().Amount() - order.ShippingFee().Amount()
		return ordervo.NewMoney(refundAmount, order.TotalAmount().Currency())
	}

	// 如果订单已完成，根据完成时间计算退款比例
	if order.IsCompleted() && order.CompletedAt() != nil {
		daysSinceCompletion := int(time.Since(*order.CompletedAt()).Hours() / 24)
		
		var refundRatio float64
		switch {
		case daysSinceCompletion <= 7:
			refundRatio = 1.0 // 7天内全额退款
		case daysSinceCompletion <= 15:
			refundRatio = 0.8 // 15天内80%退款
		case daysSinceCompletion <= 30:
			refundRatio = 0.5 // 30天内50%退款
		default:
			refundRatio = 0.0 // 超过30天不退款
		}

		refundAmount := int64(float64(order.TotalAmount().Amount()) * refundRatio)
		return ordervo.NewMoney(refundAmount, order.TotalAmount().Currency())
	}

	return ordervo.NewMoney(0, order.TotalAmount().Currency())
}

// GetAvailableOrderActions 获取订单可用操作
func (s *OrderDomainService) GetAvailableOrderActions(userID uservo.UserID, order *entity.Order) []string {
	actions := make([]string, 0)

	// 验证用户权限
	if !s.CanUserAccessOrder(userID, order) {
		return actions
	}

	// 根据订单状态判断可用操作
	switch {
	case order.Status().IsCreated():
		actions = append(actions, "modify", "submit", "cancel")
	case order.Status().IsSubmitted():
		actions = append(actions, "pay", "cancel")
	case order.Status().IsPaid():
		if !order.IsShipped() {
			actions = append(actions, "cancel")
		}
		actions = append(actions, "refund")
	case order.Status().IsShipped():
		actions = append(actions, "confirm_delivery", "refund")
	case order.Status().IsDelivered():
		actions = append(actions, "complete", "refund")
	case order.Status().IsCompleted():
		if s.CanRefundOrder(userID, order) {
			actions = append(actions, "refund")
		}
	}

	return actions
}

// EstimateDeliveryTime 估算配送时间
func (s *OrderDomainService) EstimateDeliveryTime(order *entity.Order) time.Time {
	baseTime := time.Now()
	
	// 根据配送方式估算时间
	switch order.ShippingMethod().Value() {
	case int32(ordervo.ShippingMethodExpress):
		return baseTime.Add(1 * 24 * time.Hour) // 快递1天
	case int32(ordervo.ShippingMethodStandard):
		return baseTime.Add(3 * 24 * time.Hour) // 标准配送3天
	case int32(ordervo.ShippingMethodEconomy):
		return baseTime.Add(7 * 24 * time.Hour) // 经济配送7天
	default:
		return baseTime.Add(3 * 24 * time.Hour) // 默认3天
	}
}

// CheckOrderTimeout 检查订单超时
func (s *OrderDomainService) CheckOrderTimeout(ctx context.Context, order *entity.Order) error {
	now := time.Now()

	// 检查提交超时（创建后24小时未提交自动取消）
	if order.Status().IsCreated() {
		submitDeadline := order.CreatedAt().Add(24 * time.Hour)
		if now.After(submitDeadline) {
			return order.Cancel()
		}
	}

	// 检查支付超时（提交后24小时未支付自动取消）
	if order.Status().IsSubmitted() {
		paymentDeadline := order.CreatedAt().Add(48 * time.Hour)
		if now.After(paymentDeadline) {
			return order.Cancel()
		}
	}

	// 检查发货超时（支付后72小时未发货提醒）
	if order.Status().IsPaid() && !order.IsShipped() {
		shipmentDeadline := order.PaidAt().Add(72 * time.Hour)
		if now.After(shipmentDeadline) {
			// 这里可以发送提醒通知，但不自动取消订单
			// 实际业务中可能需要人工处理
		}
	}

	return nil
}

// GenerateOrderReport 生成订单报告
func (s *OrderDomainService) GenerateOrderReport(ctx context.Context, orders []*entity.Order) *OrderReport {
	report := &OrderReport{
		TotalOrders:    len(orders),
		TotalRevenue:   0,
		StatusBreakdown: make(map[string]int),
		PaymentBreakdown: make(map[string]int),
	}

	for _, order := range orders {
		// 统计收入（只计算已支付订单）
		if order.IsPaid() {
			report.TotalRevenue += order.TotalAmount().Amount()
		}

		// 统计状态分布
		statusName := order.Status().String()
		report.StatusBreakdown[statusName]++

		// 统计支付方式分布
		if order.IsPaid() {
			paymentName := order.PaymentMethod().String()
			report.PaymentBreakdown[paymentName]++
		}
	}

	return report
}

// validateAddress 验证地址信息
func (s *OrderDomainService) validateAddress(address ordervo.Address) error {
	if address.Country() == "" {
		return ordervo.ErrInvalidAddress
	}
	if address.Province() == "" {
		return ordervo.ErrInvalidAddress
	}
	if address.City() == "" {
		return ordervo.ErrInvalidAddress
	}
	if address.DetailAddress() == "" {
		return ordervo.ErrInvalidAddress
	}
	if address.ContactName() == "" {
		return ordervo.ErrInvalidAddress
	}
	if address.ContactPhone() == "" {
		return ordervo.ErrInvalidAddress
	}
	return nil
}

// OrderReport 订单报告
type OrderReport struct {
	TotalOrders      int            `json:"total_orders"`
	TotalRevenue     int64          `json:"total_revenue"`
	StatusBreakdown  map[string]int `json:"status_breakdown"`
	PaymentBreakdown map[string]int `json:"payment_breakdown"`
}
