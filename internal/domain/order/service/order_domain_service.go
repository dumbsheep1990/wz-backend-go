package service

import (
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
