package service

import (
	"context"
	"errors"
	"time"
	"wz-backend-go/internal/application/order/dto"
	orderdomain "wz-backend-go/internal/domain/order/entity"
	orderrepo "wz-backend-go/internal/domain/order/repository"
	orderservice "wz-backend-go/internal/domain/order/service"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
	productvo "wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"

	"github.com/google/uuid"
)

// OrderApplicationService 订单应用服务接口
type OrderApplicationService interface {
	// 命令方法
	CreateOrder(ctx context.Context, cmd dto.CreateOrderCommand) (*dto.OrderDTO, error)
	UpdateOrder(ctx context.Context, cmd dto.UpdateOrderCommand) (*dto.OrderDTO, error)
	AddOrderItem(ctx context.Context, cmd dto.AddOrderItemCommand) (*dto.OrderDTO, error)
	RemoveOrderItem(ctx context.Context, cmd dto.RemoveOrderItemCommand) (*dto.OrderDTO, error)
	UpdateOrderItemQuantity(ctx context.Context, cmd dto.UpdateOrderItemQuantityCommand) (*dto.OrderDTO, error)
	AddOrderDiscount(ctx context.Context, cmd dto.AddOrderDiscountCommand) (*dto.OrderDTO, error)
	RemoveOrderDiscount(ctx context.Context, cmd dto.RemoveOrderDiscountCommand) (*dto.OrderDTO, error)
	SetOrderShippingFee(ctx context.Context, cmd dto.SetOrderShippingFeeCommand) (*dto.OrderDTO, error)
	SetOrderTax(ctx context.Context, cmd dto.SetOrderTaxCommand) (*dto.OrderDTO, error)
	SetOrderPaymentMethod(ctx context.Context, cmd dto.SetOrderPaymentMethodCommand) (*dto.OrderDTO, error)
	SetOrderShippingMethod(ctx context.Context, cmd dto.SetOrderShippingMethodCommand) (*dto.OrderDTO, error)
	SetOrderNote(ctx context.Context, cmd dto.SetOrderNoteCommand) (*dto.OrderDTO, error)
	SetOrderTrackingNumber(ctx context.Context, cmd dto.SetOrderTrackingNumberCommand) (*dto.OrderDTO, error)
	SubmitOrder(ctx context.Context, cmd dto.SubmitOrderCommand) (*dto.OrderDTO, error)
	PayOrder(ctx context.Context, cmd dto.PayOrderCommand) (*dto.OrderDTO, error)
	ShipOrder(ctx context.Context, cmd dto.ShipOrderCommand) (*dto.OrderDTO, error)
	DeliverOrder(ctx context.Context, cmd dto.DeliverOrderCommand) (*dto.OrderDTO, error)
	CompleteOrder(ctx context.Context, cmd dto.CompleteOrderCommand) (*dto.OrderDTO, error)
	CancelOrder(ctx context.Context, cmd dto.CancelOrderCommand) (*dto.OrderDTO, error)
	RequestRefund(ctx context.Context, cmd dto.RequestRefundCommand) (*dto.OrderDTO, error)
	RefundOrder(ctx context.Context, cmd dto.RefundOrderCommand) (*dto.OrderDTO, error)
	
	// 查询方法
	GetOrder(ctx context.Context, query dto.GetOrderQuery) (*dto.OrderDTO, error)
	ListOrders(ctx context.Context, query dto.ListOrdersQuery) (*dto.OrderListDTO, error)
	SearchOrders(ctx context.Context, query dto.SearchOrdersQuery) (*dto.OrderListDTO, error)
}

// orderApplicationService 订单应用服务实现
type orderApplicationService struct {
	orderRepository    orderrepo.OrderRepository
	orderDomainService orderservice.OrderDomainService
}

// 错误定义
var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderCannotBeModified = errors.New("order cannot be modified")
	ErrInvalidOrderStatus    = errors.New("invalid order status")
	ErrInvalidParameter      = errors.New("invalid parameter")
)

// NewOrderApplicationService 创建订单应用服务
func NewOrderApplicationService(
	orderRepository orderrepo.OrderRepository,
	orderDomainService orderservice.OrderDomainService,
) OrderApplicationService {
	return &orderApplicationService{
		orderRepository:    orderRepository,
		orderDomainService: orderDomainService,
	}
}

// CreateOrder 创建订单
func (s *orderApplicationService) CreateOrder(ctx context.Context, cmd dto.CreateOrderCommand) (*dto.OrderDTO, error) {
	// 解析客户ID
	customerID, err := uservo.NewUserID(cmd.CustomerID)
	if err != nil {
		return nil, err
	}

	// 解析地址
	shippingAddress, err := parseAddressDTO(cmd.ShippingAddress)
	if err != nil {
		return nil, err
	}

	// 如果未提供账单地址，使用配送地址
	billingAddress := shippingAddress
	if cmd.BillingAddress.Name != "" {
		billingAddress, err = parseAddressDTO(cmd.BillingAddress)
		if err != nil {
			return nil, err
		}
	}

	// 解析配送方式
	shippingMethod, err := ordervo.ParseShippingMethod(cmd.ShippingMethod)
	if err != nil {
		return nil, err
	}

	// 创建订单
	order, err := orderdomain.NewOrder(
		customerID,
		shippingAddress,
		billingAddress,
		shippingMethod,
	)
	if err != nil {
		return nil, err
	}

	// 添加订单项
	for _, item := range cmd.Items {
		productID, err := productvo.NewProductID(item.ProductID)
		if err != nil {
			return nil, err
		}

		unitPrice, err := ordervo.NewMoney(float64(item.UnitPrice.Amount), item.UnitPrice.Currency)
		if err != nil {
			return nil, err
		}

		err = order.AddItem(
			productID,
			item.ProductName,
			item.ProductSKU,
			item.Quantity,
			unitPrice,
			item.Attributes,
		)
		if err != nil {
			return nil, err
		}
	}

	// 设置支付方式（如果提供）
	if cmd.PaymentMethod != "" {
		paymentMethod, err := ordervo.ParsePaymentMethod(cmd.PaymentMethod)
		if err != nil {
			return nil, err
		}
		if err := order.SetPaymentMethod(paymentMethod); err != nil {
			return nil, err
		}
	}

	// 设置订单备注（如果提供）
	if cmd.Note != "" {
		order.SetNote(cmd.Note)
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// UpdateOrder 更新订单
func (s *orderApplicationService) UpdateOrder(ctx context.Context, cmd dto.UpdateOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 更新配送地址（如果提供）
	if cmd.ShippingAddress != nil {
		shippingAddress, err := parseAddressDTO(*cmd.ShippingAddress)
		if err != nil {
			return nil, err
		}
		order.SetShippingAddress(shippingAddress)
	}

	// 更新账单地址（如果提供）
	if cmd.BillingAddress != nil {
		billingAddress, err := parseAddressDTO(*cmd.BillingAddress)
		if err != nil {
			return nil, err
		}
		order.SetBillingAddress(billingAddress)
	}

	// 更新配送方式（如果提供）
	if cmd.ShippingMethod != "" {
		shippingMethod, err := ordervo.ParseShippingMethod(cmd.ShippingMethod)
		if err != nil {
			return nil, err
		}
		if err := order.SetShippingMethod(shippingMethod); err != nil {
			return nil, err
		}
	}

	// 更新订单备注（如果提供）
	if cmd.Note != "" {
		order.SetNote(cmd.Note)
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// AddOrderItem 添加订单项
func (s *orderApplicationService) AddOrderItem(ctx context.Context, cmd dto.AddOrderItemCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析产品ID
	productID, err := productvo.NewProductID(cmd.ProductID)
	if err != nil {
		return nil, err
	}

	// 解析单价
	unitPrice, err := ordervo.NewMoney(cmd.UnitPrice.Amount, cmd.UnitPrice.Currency)
	if err != nil {
		return nil, err
	}

	// 添加订单项
	err = order.AddItem(
		productID,
		cmd.ProductName,
		cmd.ProductSKU,
		cmd.Quantity,
		unitPrice,
		cmd.Attributes,
	)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// RemoveOrderItem 移除订单项
func (s *orderApplicationService) RemoveOrderItem(ctx context.Context, cmd dto.RemoveOrderItemCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析订单项ID
	itemID, err := ordervo.NewOrderItemID(cmd.ItemID)
	if err != nil {
		return nil, err
	}

	// 移除订单项
	err = order.RemoveItem(itemID)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// UpdateOrderItemQuantity 更新订单项数量
func (s *orderApplicationService) UpdateOrderItemQuantity(ctx context.Context, cmd dto.UpdateOrderItemQuantityCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析订单项ID
	itemID, err := ordervo.NewOrderItemID(cmd.ItemID)
	if err != nil {
		return nil, err
	}

	// 更新订单项数量
	err = order.UpdateItemQuantity(itemID, cmd.Quantity)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// AddOrderDiscount 添加订单折扣
func (s *orderApplicationService) AddOrderDiscount(ctx context.Context, cmd dto.AddOrderDiscountCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 创建折扣金额
	amount, err := ordervo.NewMoney(cmd.Amount.Amount, cmd.Amount.Currency)
	if err != nil {
		return nil, err
	}

	// 创建折扣
	discount := orderdomain.NewOrderDiscount(
		uuid.New().String(),
		cmd.DiscountType,
		cmd.DiscountName,
		amount,
	)

	// 添加折扣
	err = order.AddDiscount(discount)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// RemoveOrderDiscount 移除订单折扣
func (s *orderApplicationService) RemoveOrderDiscount(ctx context.Context, cmd dto.RemoveOrderDiscountCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 移除折扣
	err = order.RemoveDiscount(cmd.DiscountID)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderShippingFee 设置订单运费
func (s *orderApplicationService) SetOrderShippingFee(ctx context.Context, cmd dto.SetOrderShippingFeeCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析运费
	fee, err := ordervo.NewMoney(cmd.Fee.Amount, cmd.Fee.Currency)
	if err != nil {
		return nil, err
	}

	// 设置运费
	err = order.SetShippingFee(fee)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderTax 设置订单税费
func (s *orderApplicationService) SetOrderTax(ctx context.Context, cmd dto.SetOrderTaxCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析税费
	tax, err := ordervo.NewMoney(cmd.Tax.Amount, cmd.Tax.Currency)
	if err != nil {
		return nil, err
	}

	// 设置税费
	err = order.SetTax(tax)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderPaymentMethod 设置订单支付方式
func (s *orderApplicationService) SetOrderPaymentMethod(ctx context.Context, cmd dto.SetOrderPaymentMethodCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 解析支付方式
	paymentMethod, err := ordervo.ParsePaymentMethod(cmd.PaymentMethod)
	if err != nil {
		return nil, err
	}

	// 设置支付方式
	err = order.SetPaymentMethod(paymentMethod)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderShippingMethod 设置订单配送方式
func (s *orderApplicationService) SetOrderShippingMethod(ctx context.Context, cmd dto.SetOrderShippingMethodCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 解析配送方式
	shippingMethod, err := ordervo.ParseShippingMethod(cmd.ShippingMethod)
	if err != nil {
		return nil, err
	}

	// 设置配送方式
	err = order.SetShippingMethod(shippingMethod)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderNote 设置订单备注
func (s *orderApplicationService) SetOrderNote(ctx context.Context, cmd dto.SetOrderNoteCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 检查订单是否可以修改
	if !order.CanModify() {
		return nil, ErrOrderCannotBeModified
	}

	// 设置备注
	order.SetNote(cmd.Note)

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SetOrderTrackingNumber 设置订单物流单号
func (s *orderApplicationService) SetOrderTrackingNumber(ctx context.Context, cmd dto.SetOrderTrackingNumberCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 设置物流单号
	err = order.SetTrackingNumber(cmd.TrackingNumber)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// SubmitOrder 提交订单
func (s *orderApplicationService) SubmitOrder(ctx context.Context, cmd dto.SubmitOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 提交订单
	err = order.SubmitOrder()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// PayOrder 支付订单
func (s *orderApplicationService) PayOrder(ctx context.Context, cmd dto.PayOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 解析支付方式
	paymentMethod, err := ordervo.ParsePaymentMethod(cmd.PaymentMethod)
	if err != nil {
		return nil, err
	}

	// 支付订单
	err = order.Pay(paymentMethod)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// ShipOrder 发货
func (s *orderApplicationService) ShipOrder(ctx context.Context, cmd dto.ShipOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 发货
	err = order.Ship(cmd.TrackingNumber)
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// DeliverOrder 标记订单为已送达
func (s *orderApplicationService) DeliverOrder(ctx context.Context, cmd dto.DeliverOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 标记为已送达
	err = order.Deliver()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// CompleteOrder 完成订单
func (s *orderApplicationService) CompleteOrder(ctx context.Context, cmd dto.CompleteOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 完成订单
	err = order.Complete()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// CancelOrder 取消订单
func (s *orderApplicationService) CancelOrder(ctx context.Context, cmd dto.CancelOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 取消订单
	err = order.Cancel()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// RequestRefund 申请退款
func (s *orderApplicationService) RequestRefund(ctx context.Context, cmd dto.RequestRefundCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 申请退款
	err = order.RequestRefund()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// RefundOrder 退款订单
func (s *orderApplicationService) RefundOrder(ctx context.Context, cmd dto.RefundOrderCommand) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(cmd.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 退款
	err = order.Refund()
	if err != nil {
		return nil, err
	}

	// 保存订单
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// GetOrder 获取订单
func (s *orderApplicationService) GetOrder(ctx context.Context, query dto.GetOrderQuery) (*dto.OrderDTO, error) {
	// 解析订单ID
	orderID, err := ordervo.NewOrderID(query.OrderID)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// 转换为DTO并返回
	return convertOrderToDTO(order), nil
}

// ListOrders 获取订单列表
func (s *orderApplicationService) ListOrders(ctx context.Context, query dto.ListOrdersQuery) (*dto.OrderListDTO, error) {
	// 初始化过滤条件
	filters := make(map[string]interface{})

	// 添加客户ID过滤（如果提供）
	if query.CustomerID != "" {
		customerID, err := uservo.NewUserID(query.CustomerID)
		if err != nil {
			return nil, err
		}
		filters["customer_id"] = customerID
	}

	// 添加状态过滤（如果提供）
	if len(query.Status) > 0 {
		var statusCodes []int32
		for _, status := range query.Status {
			statusCode, err := ordervo.ParseOrderStatusCode(status)
			if err != nil {
				return nil, err
			}
			statusCodes = append(statusCodes, statusCode)
		}
		filters["status"] = statusCodes
	}

	// 添加日期范围过滤
	if !query.StartDate.IsZero() {
		filters["start_date"] = query.StartDate
	}
	if !query.EndDate.IsZero() {
		filters["end_date"] = query.EndDate
	}

	// 设置分页和排序
	pagination := map[string]interface{}{
		"page":      query.Page,
		"page_size": query.PageSize,
		"sort_by":   query.SortBy,
		"sort_order": query.SortOrder,
	}

	// 查询订单
	orders, total, err := s.orderRepository.FindByFilters(ctx, filters, pagination)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	orderDTOs := make([]dto.OrderDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = *convertOrderToDTO(order)
	}

	// 返回结果
	return &dto.OrderListDTO{
		Total:  total,
		Orders: orderDTOs,
	}, nil
}

// SearchOrders 搜索订单
func (s *orderApplicationService) SearchOrders(ctx context.Context, query dto.SearchOrdersQuery) (*dto.OrderListDTO, error) {
	// 构建搜索条件
	searchCriteria := map[string]interface{}{
		"keyword": query.Keyword,
	}

	// 添加状态过滤（如果提供）
	if len(query.Status) > 0 {
		var statusCodes []int32
		for _, status := range query.Status {
			statusCode, err := ordervo.ParseOrderStatusCode(status)
			if err != nil {
				return nil, err
			}
			statusCodes = append(statusCodes, statusCode)
		}
		searchCriteria["status"] = statusCodes
	}

	// 添加日期范围过滤
	if !query.StartDate.IsZero() {
		searchCriteria["start_date"] = query.StartDate
	}
	if !query.EndDate.IsZero() {
		searchCriteria["end_date"] = query.EndDate
	}

	// 设置分页和排序
	pagination := map[string]interface{}{
		"page":      query.Page,
		"page_size": query.PageSize,
		"sort_by":   query.SortBy,
		"sort_order": query.SortOrder,
	}

	// 搜索订单
	orders, total, err := s.orderRepository.Search(ctx, searchCriteria, pagination)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	orderDTOs := make([]dto.OrderDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = *convertOrderToDTO(order)
	}

	// 返回结果
	return &dto.OrderListDTO{
		Total:  total,
		Orders: orderDTOs,
	}, nil
}

// 辅助函数部分

// convertOrderToDTO 将订单实体转换为DTO
func convertOrderToDTO(order *orderdomain.Order) *dto.OrderDTO {
	// 初始化DTO
	orderDTO := &dto.OrderDTO{
		ID:           order.ID().String(),
		OrderNumber:  order.OrderNumber().String(),
		CustomerID:   order.CustomerID().String(),
		Status:       order.Status().String(),
		StatusCode:   order.Status().Code(),
		ShippingFee:  convertMoneyToDTO(order.ShippingFee()),
		Tax:          convertMoneyToDTO(order.Tax()),
		DiscountAmount: convertMoneyToDTO(order.DiscountAmount()),
		Subtotal:     convertMoneyToDTO(order.Subtotal()),
		TotalAmount:  convertMoneyToDTO(order.TotalAmount()),
		Note:         order.Note(),
		TrackingNumber: order.TrackingNumber(),
		PaymentMethod: order.PaymentMethod().String(),
		ShippingMethod: order.ShippingMethod().String(),
		ShippingAddress: convertAddressToDTO(order.ShippingAddress()),
		BillingAddress: convertAddressToDTO(order.BillingAddress()),
		PaidAt:       order.PaidAt(),
		ShippedAt:    order.ShippedAt(),
		DeliveredAt:  order.DeliveredAt(),
		CompletedAt:  order.CompletedAt(),
		CancelledAt:  order.CancelledAt(),
		RefundedAt:   order.RefundedAt(),
		CreatedAt:    order.CreatedAt(),
		UpdatedAt:    order.UpdatedAt(),
	}

	// 转换订单项
	orderItems := order.Items()
	itemDTOs := make([]dto.OrderItemDTO, len(orderItems))
	for i, item := range orderItems {
		itemDTOs[i] = convertOrderItemToDTO(item)
	}
	orderDTO.Items = itemDTOs

	// 转换折扣
	orderDiscounts := order.Discounts()
	discountDTOs := make([]dto.OrderDiscountDTO, len(orderDiscounts))
	for i, discount := range orderDiscounts {
		discountDTOs[i] = convertOrderDiscountToDTO(discount)
	}
	orderDTO.Discounts = discountDTOs

	return orderDTO
}

// convertOrderItemToDTO 将订单项实体转换为DTO
func convertOrderItemToDTO(item *orderdomain.OrderItem) dto.OrderItemDTO {
	return dto.OrderItemDTO{
		ID:         item.ID().String(),
		ProductID:  item.ProductID().String(),
		ProductName: item.ProductName(),
		ProductSKU: item.ProductSKU(),
		Quantity:   item.Quantity(),
		UnitPrice:  convertMoneyToDTO(item.UnitPrice()),
		TotalPrice: convertMoneyToDTO(item.TotalPrice()),
		Attributes: item.Attributes(),
	}
}

// convertOrderDiscountToDTO 将订单折扣实体转换为DTO
func convertOrderDiscountToDTO(discount *orderdomain.OrderDiscount) dto.OrderDiscountDTO {
	return dto.OrderDiscountDTO{
		ID:          discount.ID(),
		DiscountType: discount.DiscountType(),
		DiscountName: discount.DiscountName(),
		Amount:      convertMoneyToDTO(discount.Amount()),
	}
}

// convertMoneyToDTO 将金额值对象转换为DTO
func convertMoneyToDTO(money ordervo.Money) dto.MoneyDTO {
	return dto.MoneyDTO{
		Amount:   money.Amount(),
		Currency: money.Currency(),
	}
}

// convertAddressToDTO 将地址值对象转换为DTO
func convertAddressToDTO(address ordervo.Address) dto.AddressDTO {
	return dto.AddressDTO{
		Name:          address.Name(),
		Phone:         address.Phone(),
		Province:      address.Province(),
		City:          address.City(),
		District:      address.District(),
		DetailAddress: address.DetailAddress(),
		PostalCode:    address.PostalCode(),
	}
}

// parseAddressDTO 解析地址DTO为值对象
func parseAddressDTO(addressDTO dto.AddressDTO) (ordervo.Address, error) {
	return ordervo.NewAddress(
		addressDTO.Name,
		addressDTO.Phone,
		addressDTO.Province,
		addressDTO.City,
		addressDTO.District,
		addressDTO.DetailAddress,
		addressDTO.PostalCode,
	)
}
