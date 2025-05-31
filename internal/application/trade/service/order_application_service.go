package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/application/trade/dto"
	"wz-backend-go/internal/domain/order/entity"
	orderEvent "wz-backend-go/internal/domain/order/event"
	orderRepo "wz-backend-go/internal/domain/order/repository"
	orderService "wz-backend-go/internal/domain/order/service"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
	productvo "wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
	"wz-backend-go/internal/infrastructure/transaction"
)

// OrderApplicationService 订单应用服务
type OrderApplicationService struct {
	orderRepo       orderRepo.OrderRepository
	orderDomain     *orderService.OrderDomainService
	eventPublisher  orderRepo.EventPublisher
	transactionMgr  transaction.Manager
}

// NewOrderApplicationService 创建订单应用服务
func NewOrderApplicationService(
	orderRepo orderRepo.OrderRepository,
	orderDomain *orderService.OrderDomainService,
	eventPublisher orderRepo.EventPublisher,
	transactionMgr transaction.Manager,
) *OrderApplicationService {
	return &OrderApplicationService{
		orderRepo:      orderRepo,
		orderDomain:    orderDomain,
		eventPublisher: eventPublisher,
		transactionMgr: transactionMgr,
	}
}

// CreateOrder 创建订单
func (s *OrderApplicationService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// 验证客户ID
	customerID := uservo.NewUserID(req.CustomerID)

	// 构建地址
	shippingAddr, err := ordervo.NewAddress(
		req.ShippingAddress.Country,
		req.ShippingAddress.Province,
		req.ShippingAddress.City,
		req.ShippingAddress.District,
		req.ShippingAddress.DetailAddress,
		req.ShippingAddress.PostalCode,
		req.ShippingAddress.ContactName,
		req.ShippingAddress.ContactPhone,
	)
	if err != nil {
		return nil, err
	}

	billingAddr, err := ordervo.NewAddress(
		req.BillingAddress.Country,
		req.BillingAddress.Province,
		req.BillingAddress.City,
		req.BillingAddress.District,
		req.BillingAddress.DetailAddress,
		req.BillingAddress.PostalCode,
		req.BillingAddress.ContactName,
		req.BillingAddress.ContactPhone,
	)
	if err != nil {
		return nil, err
	}

	// 构建配送方式
	shippingMethod, err := ordervo.NewShippingMethod(req.ShippingMethod)
	if err != nil {
		return nil, err
	}

	var result *dto.OrderResponse
	// 开启事务
	err = s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 创建订单
		order, err := entity.NewOrder(customerID, shippingAddr, billingAddr, shippingMethod)
		if err != nil {
			return err
		}

		// 添加订单项
		for _, item := range req.Items {
			productID := productvo.NewProductID(item.ProductID)
			unitPrice, err := ordervo.NewMoney(item.UnitPrice, req.Currency)
			if err != nil {
				return err
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
				return err
			}
		}

		// 设置配送费
		if req.ShippingFee > 0 {
			shippingFee, err := ordervo.NewMoney(req.ShippingFee, req.Currency)
			if err != nil {
				return err
			}
			err = order.SetShippingFee(shippingFee)
			if err != nil {
				return err
			}
		}

		// 设置税费
		if req.Tax > 0 {
			tax, err := ordervo.NewMoney(req.Tax, req.Currency)
			if err != nil {
				return err
			}
			err = order.SetTax(tax)
			if err != nil {
				return err
			}
		}

		// 设置备注
		if req.Note != "" {
			order.SetNote(req.Note)
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布订单创建事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderCreatedEvent(order))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// GetOrder 获取订单详情
func (s *OrderApplicationService) GetOrder(ctx context.Context, req *dto.GetOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)
	
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 验证权限
	if !s.orderDomain.CanUserAccessOrder(uservo.NewUserID(req.UserID), order) {
		return nil, errors.New("无权访问该订单")
	}

	return s.buildOrderResponse(order), nil
}

// GetOrderByNumber 根据订单号获取订单
func (s *OrderApplicationService) GetOrderByNumber(ctx context.Context, req *dto.GetOrderByNumberRequest) (*dto.OrderResponse, error) {
	orderNumber := ordervo.NewOrderNumber(req.OrderNumber)
	
	order, err := s.orderRepo.FindByOrderNumber(orderNumber)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 验证权限
	if !s.orderDomain.CanUserAccessOrder(uservo.NewUserID(req.UserID), order) {
		return nil, errors.New("无权访问该订单")
	}

	return s.buildOrderResponse(order), nil
}

// ListUserOrders 获取用户订单列表
func (s *OrderApplicationService) ListUserOrders(ctx context.Context, req *dto.ListUserOrdersRequest) (*dto.ListOrdersResponse, error) {
	customerID := uservo.NewUserID(req.CustomerID)
	
	orders, total, err := s.orderRepo.FindByCustomerID(customerID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, s.buildOrderResponse(order))
	}

	return &dto.ListOrdersResponse{
		Orders:   orderResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// PayOrder 支付订单
func (s *OrderApplicationService) PayOrder(ctx context.Context, req *dto.PayOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)
	
	// 构建支付方式
	paymentMethod, err := ordervo.NewPaymentMethod(req.PaymentMethod)
	if err != nil {
		return nil, err
	}

	var result *dto.OrderResponse
	err = s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 验证权限
		if !s.orderDomain.CanUserAccessOrder(uservo.NewUserID(req.UserID), order) {
			return errors.New("无权操作该订单")
		}

		// 执行支付
		err = order.Pay(paymentMethod)
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布支付事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderPaidEvent(order))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// ShipOrder 发货
func (s *OrderApplicationService) ShipOrder(ctx context.Context, req *dto.ShipOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)

	var result *dto.OrderResponse
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 执行发货
		err = order.Ship(req.TrackingNumber)
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布发货事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderShippedEvent(order))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// DeliverOrder 确认送达
func (s *OrderApplicationService) DeliverOrder(ctx context.Context, req *dto.DeliverOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)

	var result *dto.OrderResponse
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 确认送达
		err = order.Deliver()
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布送达事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderDeliveredEvent(order))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// CompleteOrder 完成订单
func (s *OrderApplicationService) CompleteOrder(ctx context.Context, req *dto.CompleteOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)

	var result *dto.OrderResponse
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 验证权限
		if !s.orderDomain.CanUserAccessOrder(uservo.NewUserID(req.UserID), order) {
			return errors.New("无权操作该订单")
		}

		// 完成订单
		err = order.Complete()
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布完成事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderCompletedEvent(order))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// CancelOrder 取消订单
func (s *OrderApplicationService) CancelOrder(ctx context.Context, req *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)

	var result *dto.OrderResponse
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 验证权限
		if !s.orderDomain.CanUserAccessOrder(uservo.NewUserID(req.UserID), order) {
			return errors.New("无权操作该订单")
		}

		// 取消订单
		err = order.Cancel()
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布取消事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderCancelledEvent(order, req.Reason))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// RefundOrder 退款
func (s *OrderApplicationService) RefundOrder(ctx context.Context, req *dto.RefundOrderRequest) (*dto.OrderResponse, error) {
	orderID := ordervo.NewOrderID(req.OrderID)

	var result *dto.OrderResponse
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取订单
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.New("订单不存在")
		}

		// 验证权限和退款资格
		if !s.orderDomain.CanRefundOrder(uservo.NewUserID(req.UserID), order) {
			return errors.New("订单不符合退款条件")
		}

		// 执行退款
		err = order.Refund()
		if err != nil {
			return err
		}

		// 保存订单
		err = s.orderRepo.Save(order)
		if err != nil {
			return err
		}

		// 发布退款事件
		err = s.eventPublisher.Publish(orderEvent.NewOrderRefundedEvent(order, req.Reason))
		if err != nil {
			return err
		}

		result = s.buildOrderResponse(order)
		return nil
	})

	return result, err
}

// SearchOrders 搜索订单
func (s *OrderApplicationService) SearchOrders(ctx context.Context, req *dto.SearchOrdersRequest) (*dto.ListOrdersResponse, error) {
	orders, total, err := s.orderRepo.Search(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, s.buildOrderResponse(order))
	}

	return &dto.ListOrdersResponse{
		Orders:   orderResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// buildOrderResponse 构建订单响应对象
func (s *OrderApplicationService) buildOrderResponse(order *entity.Order) *dto.OrderResponse {
	// 构建订单项响应
	items := make([]*dto.OrderItemResponse, 0, len(order.Items()))
	for _, item := range order.Items() {
		items = append(items, &dto.OrderItemResponse{
			ID:          item.ID().Value(),
			ProductID:   item.ProductID().Value(),
			ProductName: item.ProductName(),
			ProductSKU:  item.ProductSKU(),
			Quantity:    item.Quantity(),
			UnitPrice:   item.UnitPrice().Amount(),
			TotalPrice:  item.TotalPrice().Amount(),
			Attributes:  item.Attributes(),
		})
	}

	// 构建折扣响应
	discounts := make([]*dto.OrderDiscountResponse, 0, len(order.Discounts()))
	for _, discount := range order.Discounts() {
		discounts = append(discounts, &dto.OrderDiscountResponse{
			ID:             discount.ID(),
			Type:           discount.Type().Value(),
			Name:           discount.Name(),
			DiscountAmount: discount.DiscountAmount().Amount(),
			Description:    discount.Description(),
		})
	}

	return &dto.OrderResponse{
		ID:          order.ID().Value(),
		OrderNumber: order.OrderNumber().Value(),
		CustomerID:  order.CustomerID().Value(),
		Status:      order.Status().Value(),
		Items:       items,
		Discounts:   discounts,
		ShippingAddress: &dto.AddressResponse{
			Country:       order.ShippingAddress().Country(),
			Province:      order.ShippingAddress().Province(),
			City:          order.ShippingAddress().City(),
			District:      order.ShippingAddress().District(),
			DetailAddress: order.ShippingAddress().DetailAddress(),
			PostalCode:    order.ShippingAddress().PostalCode(),
			ContactName:   order.ShippingAddress().ContactName(),
			ContactPhone:  order.ShippingAddress().ContactPhone(),
		},
		BillingAddress: &dto.AddressResponse{
			Country:       order.BillingAddress().Country(),
			Province:      order.BillingAddress().Province(),
			City:          order.BillingAddress().City(),
			District:      order.BillingAddress().District(),
			DetailAddress: order.BillingAddress().DetailAddress(),
			PostalCode:    order.BillingAddress().PostalCode(),
			ContactName:   order.BillingAddress().ContactName(),
			ContactPhone:  order.BillingAddress().ContactPhone(),
		},
		PaymentMethod:  order.PaymentMethod().Value(),
		ShippingMethod: order.ShippingMethod().Value(),
		Subtotal:       order.Subtotal().Amount(),
		ShippingFee:    order.ShippingFee().Amount(),
		Tax:            order.Tax().Amount(),
		DiscountAmount: order.DiscountAmount().Amount(),
		TotalAmount:    order.TotalAmount().Amount(),
		Currency:       order.TotalAmount().Currency(),
		Note:           order.Note(),
		TrackingNumber: order.TrackingNumber(),
		PaidAt:         order.PaidAt(),
		ShippedAt:      order.ShippedAt(),
		DeliveredAt:    order.DeliveredAt(),
		CompletedAt:    order.CompletedAt(),
		CancelledAt:    order.CancelledAt(),
		RefundedAt:     order.RefundedAt(),
		CreatedAt:      order.CreatedAt(),
		UpdatedAt:      order.UpdatedAt(),
	}
} 