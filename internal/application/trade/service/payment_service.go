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

// PaymentService 支付应用服务
type PaymentService struct {
	paymentRepository repository.PaymentRepository
	orderRepository   repository.OrderRepository
	eventBus          event.EventBus
	validator         *validator.Validate
	unitOfWork        database.UnitOfWork
}

// NewPaymentService 创建支付应用服务
func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	orderRepository repository.OrderRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *PaymentService {
	return &PaymentService{
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		eventBus:          eventBus,
		validator:         validator.New(),
		unitOfWork:        unitOfWork,
	}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证创建支付请求失败: %w", err)
	}

	// 转换订单ID
	orderID, err := valueobject.NewOrderID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 转换金额
	amount, err := valueobject.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("无效的支付金额: %w", err)
	}

	// 转换支付方式
	method, err := valueobject.NewPaymentMethod(req.Method)
	if err != nil {
		return nil, fmt.Errorf("无效的支付方式: %w", err)
	}

	// 检查订单是否存在
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	if order == nil {
		return nil, errors.New("未找到订单")
	}

	// 检查订单是否可以支付
	if order.Status() != valueobject.OrderStatusPending {
		return nil, errors.New("只有待支付状态的订单可以创建支付")
	}

	// 检查用户ID是否匹配
	if order.UserID().String() != userID.String() {
		return nil, errors.New("订单用户ID不匹配")
	}

	// 检查支付金额是否匹配
	if order.TotalAmount().Amount() != amount.Amount() || order.TotalAmount().Currency() != amount.Currency() {
		return nil, errors.New("支付金额与订单金额不匹配")
	}

	// 检查是否已存在支付
	existingPayment, err := s.paymentRepository.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	}

	if existingPayment != nil {
		// 如果存在支付但是失败了，可以重新创建
		if existingPayment.Status() == valueobject.PaymentStatusFailed {
			// 删除失败的支付记录
			err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
				return s.paymentRepository.Delete(ctx, existingPayment.ID())
			})
			if err != nil {
				return nil, fmt.Errorf("删除失败的支付记录失败: %w", err)
			}
		} else {
			return nil, errors.New("订单已存在支付记录")
		}
	}

	// 创建支付实体
	payment, err := entity.NewPayment(orderID, userID, amount, method)
	if err != nil {
		return nil, fmt.Errorf("创建支付失败: %w", err)
	}

	// 使用工作单元保存支付并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepository.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}

		// 发布领域事件
		for _, event := range payment.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布支付事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		payment.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// CompletePayment 完成支付
func (s *PaymentService) CompletePayment(ctx context.Context, req dto.CompletePaymentRequest) (*dto.PaymentDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证完成支付请求失败: %w", err)
	}

	// 转换支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 完成支付
	if err := payment.Pay(req.TransactionID); err != nil {
		return nil, fmt.Errorf("完成支付失败: %w", err)
	}

	// 使用工作单元保存支付、更新订单并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		// 保存支付
		if err := s.paymentRepository.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}

		// 查询订单
		order, err := s.orderRepository.FindByID(ctx, payment.OrderID())
		if err != nil {
			return fmt.Errorf("查询订单失败: %w", err)
		}

		if order == nil {
			return errors.New("未找到订单")
		}

		// 更新订单状态为已支付
		if err := order.MarkAsPaid(paymentID); err != nil {
			return fmt.Errorf("标记订单为已支付失败: %w", err)
		}

		// 保存订单
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 发布支付事件
		for _, event := range payment.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布支付事件失败: %w", err)
			}
		}

		// 发布订单事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		payment.ClearDomainEvents()
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// FailPayment 标记支付失败
func (s *PaymentService) FailPayment(ctx context.Context, req dto.FailPaymentRequest) (*dto.PaymentDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证标记支付失败请求失败: %w", err)
	}

	// 转换支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 标记支付失败
	if err := payment.Fail(req.Reason); err != nil {
		return nil, fmt.Errorf("标记支付失败失败: %w", err)
	}

	// 使用工作单元保存支付并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepository.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}

		// 发布领域事件
		for _, event := range payment.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布支付事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		payment.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// RequestRefund 申请退款
func (s *PaymentService) RequestRefund(ctx context.Context, req dto.RequestRefundRequest) (*dto.PaymentDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证申请退款请求失败: %w", err)
	}

	// 转换支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 检查用户ID是否匹配
	if payment.UserID().String() != userID.String() {
		return nil, errors.New("支付用户ID不匹配")
	}

	// 申请退款
	if err := payment.RequestRefund(); err != nil {
		return nil, fmt.Errorf("申请退款失败: %w", err)
	}

	// 使用工作单元保存支付、更新订单并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		// 保存支付
		if err := s.paymentRepository.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}

		// 查询订单
		order, err := s.orderRepository.FindByID(ctx, payment.OrderID())
		if err != nil {
			return fmt.Errorf("查询订单失败: %w", err)
		}

		if order == nil {
			return errors.New("未找到订单")
		}

		// 更新订单状态为退款中
		if err := order.RequestRefund(); err != nil {
			return fmt.Errorf("标记订单为退款中失败: %w", err)
		}

		// 保存订单
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 发布支付事件
		for _, event := range payment.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布支付事件失败: %w", err)
			}
		}

		// 发布订单事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		payment.ClearDomainEvents()
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// CompleteRefund 完成退款
func (s *PaymentService) CompleteRefund(ctx context.Context, req dto.CompleteRefundRequest) (*dto.PaymentDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证完成退款请求失败: %w", err)
	}

	// 转换支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 完成退款
	if err := payment.CompleteRefund(); err != nil {
		return nil, fmt.Errorf("完成退款失败: %w", err)
	}

	// 使用工作单元保存支付、更新订单并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		// 保存支付
		if err := s.paymentRepository.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}

		// 查询订单
		order, err := s.orderRepository.FindByID(ctx, payment.OrderID())
		if err != nil {
			return fmt.Errorf("查询订单失败: %w", err)
		}

		if order == nil {
			return errors.New("未找到订单")
		}

		// 更新订单状态为已退款
		if err := order.MarkAsRefunded(); err != nil {
			return fmt.Errorf("标记订单为已退款失败: %w", err)
		}

		// 保存订单
		if err := s.orderRepository.Save(ctx, order); err != nil {
			return fmt.Errorf("保存订单失败: %w", err)
		}

		// 发布支付事件
		for _, event := range payment.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布支付事件失败: %w", err)
			}
		}

		// 发布订单事件
		for _, event := range order.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布订单事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		payment.ClearDomainEvents()
		order.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// GetPayment 获取支付
func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*dto.PaymentDTO, error) {
	// 验证支付ID
	id, err := valueobject.NewPaymentID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// GetPaymentByOrderID 根据订单ID获取支付
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*dto.PaymentDTO, error) {
	// 验证订单ID
	id, err := valueobject.NewOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}

	// 查询支付
	payment, err := s.paymentRepository.FindByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	if payment == nil {
		return nil, errors.New("未找到支付")
	}

	// 转换为DTO
	return s.toPaymentDTO(payment), nil
}

// QueryPayments 查询支付列表
func (s *PaymentService) QueryPayments(ctx context.Context, req dto.PaymentQueryRequest) (*dto.PaymentQueryResponse, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证支付查询请求失败: %w", err)
	}

	var userID valueobject.UserID
	var orderID valueobject.OrderID
	var status valueobject.PaymentStatus
	var err error

	// 转换参数
	if req.UserID != "" {
		userID, err = valueobject.NewUserID(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("无效的用户ID: %w", err)
		}
	}

	if req.OrderID != "" {
		orderID, err = valueobject.NewOrderID(req.OrderID)
		if err != nil {
			return nil, fmt.Errorf("无效的订单ID: %w", err)
		}
	}

	if req.Status != "" {
		status, err = valueobject.NewPaymentStatus(req.Status)
		if err != nil {
			return nil, fmt.Errorf("无效的支付状态: %w", err)
		}
	}

	var payments []*entity.Payment
	var totalCount int64

	// 根据条件查询
	if req.OrderID != "" {
		// 根据订单ID查询单个支付
		payment, err := s.paymentRepository.FindByOrderID(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("查询支付失败: %w", err)
		}

		if payment != nil {
			payments = []*entity.Payment{payment}
			totalCount = 1
		} else {
			payments = []*entity.Payment{}
			totalCount = 0
		}
	} else if req.UserID != "" && req.Status != "" {
		// 查询指定用户和状态的支付
		// 这里假设仓储接口有根据用户ID和状态查询的方法，如果没有，需要添加
		// 此处简化为只根据用户ID查询
		payments, totalCount, err = s.paymentRepository.FindByUserID(ctx, userID, req.Page, req.PerPage)
		// 过滤状态
		var filteredPayments []*entity.Payment
		for _, p := range payments {
			if p.Status() == status {
				filteredPayments = append(filteredPayments, p)
			}
		}
		payments = filteredPayments
		totalCount = int64(len(filteredPayments))
	} else if req.UserID != "" {
		// 查询指定用户的所有支付
		payments, totalCount, err = s.paymentRepository.FindByUserID(ctx, userID, req.Page, req.PerPage)
	} else if req.Status != "" {
		// 查询指定状态的所有支付
		payments, totalCount, err = s.paymentRepository.FindByStatus(ctx, status, req.Page, req.PerPage)
	} else {
		// 不应该出现的情况，至少需要指定用户ID、订单ID或状态
		return nil, errors.New("查询参数不足，至少需要指定用户ID、订单ID或支付状态")
	}

	if err != nil {
		return nil, fmt.Errorf("查询支付失败: %w", err)
	}

	// 转换为DTO
	paymentDTOs := make([]dto.PaymentDTO, 0, len(payments))
	for _, payment := range payments {
		paymentDTOs = append(paymentDTOs, *s.toPaymentDTO(payment))
	}

	return &dto.PaymentQueryResponse{
		Payments:   paymentDTOs,
		TotalCount: totalCount,
		Page:       req.Page,
		PerPage:    req.PerPage,
	}, nil
}

// 将支付实体转换为DTO
func (s *PaymentService) toPaymentDTO(payment *entity.Payment) *dto.PaymentDTO {
	return &dto.PaymentDTO{
		ID:            payment.ID().String(),
		OrderID:       payment.OrderID().String(),
		UserID:        payment.UserID().String(),
		Amount:        payment.Amount().Amount(),
		Currency:      payment.Amount().Currency(),
		Method:        string(payment.Method()),
		Status:        string(payment.Status()),
		TransactionID: payment.TransactionID(),
		PaymentTime:   payment.PaymentTime(),
		RefundTime:    payment.RefundTime(),
		CreatedAt:     payment.CreatedAt(),
		UpdatedAt:     payment.UpdatedAt(),
	}
}
