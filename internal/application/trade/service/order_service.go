package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/wz-backend-go/internal/application/trade/dto"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/repository"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/database"
)

// OrderService 订单应用服务
type OrderService struct {
	orderRepository  repository.OrderRepository
	cartRepository   repository.CartRepository
	paymentRepository repository.PaymentRepository
	eventBus         event.EventBus
	validator        *validator.Validate
	unitOfWork       database.UnitOfWork
}

// NewOrderService 创建订单应用服务
func NewOrderService(
	orderRepository repository.OrderRepository,
	cartRepository repository.CartRepository,
	paymentRepository repository.PaymentRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *OrderService {
	return &OrderService{
		orderRepository:  orderRepository,
		cartRepository:   cartRepository,
		paymentRepository: paymentRepository,
		eventBus:         eventBus,
		validator:        validator.New(),
		unitOfWork:       unitOfWork,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证创建订单请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 转换订单项
	orderItems := make([]*entity.OrderItem, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		productID, err := valueobject.NewProductID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("无效的商品ID: %w", err)
		}

		price, err := valueobject.NewMoney(item.Price, item.Currency)
		if err != nil {
			return nil, fmt.Errorf("无效的商品价格: %w", err)
		}

		quantity, err := valueobject.NewQuantity(item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("无效的商品数量: %w", err)
		}

		orderItem, err := entity.NewOrderItem(productID, item.Name, price, quantity)
		if err != nil {
			return nil, fmt.Errorf("创建订单项失败: %w", err)
		}

		orderItems = append(orderItems, orderItem)
	}

	// 转换收货地址
	address, err := valueobject.NewAddress(
		req.ShippingAddress.Province,
		req.ShippingAddress.City,
		req.ShippingAddress.District,
		req.ShippingAddress.Detail,
		req.ShippingAddress.Receiver,
		req.ShippingAddress.PhoneNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("无效的收货地址: %w", err)
	}

	// 创建订单实体
	order, err := entity.NewOrder(userID, orderItems, address)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 使用工作单元保存订单并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 发布领域事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toOrderDTO(order), nil
}

// CreateOrderFromCart 从购物车创建订单
func (s *OrderService) CreateOrderFromCart(ctx context.Context, req dto.CreateOrderFromCartRequest) (*dto.OrderDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证从购物车创建订单请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 获取用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %w", err)
	}

	if cart == nil {
		return nil, errors.New("未找到购物车")
	}

	if cart.IsEmpty() {
		return nil, errors.New("购物车为空，无法创建订单")
	}

	// 转换收货地址
	address, err := valueobject.NewAddress(
		req.ShippingAddress.Province,
		req.ShippingAddress.City,
		req.ShippingAddress.District,
		req.ShippingAddress.Detail,
		req.ShippingAddress.Receiver,
		req.ShippingAddress.PhoneNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("无效的收货地址: %w", err)
	}

	// 从购物车创建订单
	order, err := cart.ToOrder(address)
	if err != nil {
		return nil, fmt.Errorf("从购物车创建订单失败: %w", err)
	}

	// 使用工作单元保存订单、更新购物车并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		// 保存订单
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 清空购物车
		cart.Clear()
		if err := s.cartRepository.Save(ctx, cart); err != nil {
			return fmt.Errorf("更新购物车失败: %w", err)
		}

		// 发布购物车事件
		for _, event := range cart.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布购物车事件失败: %w", err)
			}
		}

		// 发布订单事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		cart.ClearDomainEvents()
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toOrderDTO(order), nil
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	// 验证订单ID
	id, err := valueobject.NewOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}

	// 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	if order == nil {
		return nil, errors.New("未找到订单")
	}

	// 转换为DTO
	return s.toOrderDTO(order), nil
}

// GetUserOrders 获取用户订单列表
func (s *OrderService) GetUserOrders(ctx context.Context, req dto.OrderQueryRequest) (*dto.OrderQueryResponse, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证订单查询请求失败: %w", err)
	}

	var userID valueobject.UserID
	var err error
	if req.UserID != "" {
		userID, err = valueobject.NewUserID(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("无效的用户ID: %w", err)
		}
	}

	var status valueobject.OrderStatus
	if req.Status != "" {
		status, err = valueobject.NewOrderStatus(req.Status)
		if err != nil {
			return nil, fmt.Errorf("无效的订单状态: %w", err)
		}
	}

	var orders []*entity.Order
	var totalCount int64

	// 根据条件查询
	if req.UserID != "" && req.Status != "" {
		// 查询指定用户和状态的订单
		orders, totalCount, err = s.orderRepository.FindByUserIDAndStatus(ctx, userID, status, req.Page, req.PerPage)
	} else if req.UserID != "" {
		// 查询指定用户的所有订单
		orders, totalCount, err = s.orderRepository.FindByUserID(ctx, userID, req.Page, req.PerPage)
	} else if req.Status != "" {
		// 查询指定状态的所有订单
		orders, totalCount, err = s.orderRepository.FindByStatus(ctx, status, req.Page, req.PerPage)
	} else {
		// 不应该出现的情况，至少需要指定用户ID或状态
		return nil, errors.New("查询参数不足，至少需要指定用户ID或订单状态")
	}

	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	// 转换为DTO
	orderDTOs := make([]dto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, *s.toOrderDTO(order))
	}

	return &dto.OrderQueryResponse{
		Orders:     orderDTOs,
		TotalCount: totalCount,
		Page:       req.Page,
		PerPage:    req.PerPage,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, req dto.UpdateOrderStatusRequest) (*dto.OrderDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证更新订单状态请求失败: %w", err)
	}

	// 转换订单ID
	orderID, err := valueobject.NewOrderID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}

	// 转换订单状态
	status, err := valueobject.NewOrderStatus(req.Status)
	if err != nil {
		return nil, fmt.Errorf("无效的订单状态: %w", err)
	}

	// 查询订单
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	if order == nil {
		return nil, errors.New("未找到订单")
	}

	// 根据目标状态调用对应的订单状态转换方法
	var updateErr error
	switch status {
	case valueobject.OrderStatusPaid:
		// 查找对应的支付记录
		payment, err := s.paymentRepository.FindByOrderID(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("查询支付记录失败: %w", err)
		}
		if payment == nil {
			return nil, errors.New("未找到支付记录，无法标记订单为已支付")
		}
		if !payment.IsPaid() {
			return nil, errors.New("支付尚未完成，无法标记订单为已支付")
		}
		updateErr = order.MarkAsPaid(payment.ID())
	case valueobject.OrderStatusShipped:
		updateErr = order.MarkAsShipped()
	case valueobject.OrderStatusDelivered:
		updateErr = order.MarkAsDelivered()
	case valueobject.OrderStatusCompleted:
		updateErr = order.MarkAsCompleted()
	case valueobject.OrderStatusCancelled:
		updateErr = order.Cancel()
	case valueobject.OrderStatusRefunding:
		updateErr = order.RequestRefund()
	case valueobject.OrderStatusRefunded:
		updateErr = order.MarkAsRefunded()
	default:
		updateErr = fmt.Errorf("不支持的状态转换: %s", status)
	}

	if updateErr != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", updateErr)
	}

	// 使用工作单元保存订单并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 发布领域事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toOrderDTO(order), nil
}

// 将订单实体转换为DTO
func (s *OrderService) toOrderDTO(order *entity.Order) *dto.OrderDTO {
	// 转换订单项
	items := make([]dto.OrderItemDTO, 0, len(order.Items()))
	for _, item := range order.Items() {
		items = append(items, dto.OrderItemDTO{
			ProductID: item.ProductID().String(),
			Name:      item.Name(),
			Price:     item.Price().Amount(),
			Currency:  item.Price().Currency(),
			Quantity:  item.Quantity().Value(),
		})
	}

	// 转换收货地址
	address := dto.AddressDTO{
		Province:    order.ShippingAddress().Province(),
		City:        order.ShippingAddress().City(),
		District:    order.ShippingAddress().District(),
		Detail:      order.ShippingAddress().Detail(),
		Receiver:    order.ShippingAddress().Receiver(),
		PhoneNumber: order.ShippingAddress().PhoneNumber(),
	}

	// 转换支付ID
	var paymentID string
	if !order.PaymentID().IsEmpty() {
		paymentID = order.PaymentID().String()
	}

	// 构造DTO
	return &dto.OrderDTO{
		ID:              order.ID().String(),
		UserID:          order.UserID().String(),
		OrderItems:      items,
		TotalAmount:     order.TotalAmount().Amount(),
		Currency:        order.TotalAmount().Currency(),
		Status:          string(order.Status()),
		ShippingAddress: address,
		PaymentID:       paymentID,
		CreatedAt:       order.CreatedAt(),
		UpdatedAt:       order.UpdatedAt(),
	}
}
