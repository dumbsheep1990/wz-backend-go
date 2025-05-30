package service

import (
	"context"
	"fmt"
	"time"
	"github.com/go-playground/validator/v10"
	"wz-backend-go/internal/application/trade/dto"
	"wz-backend-go/internal/domain/trade/entity"
	"wz-backend-go/internal/domain/trade/repository"
	tradeService "wz-backend-go/internal/domain/trade/service"
	"wz-backend-go/internal/domain/trade/valueobject"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/database"
)

// PaymentApplicationService 支付应用服务
type PaymentApplicationService struct {
	paymentRepo   repository.PaymentRepository
	domainService *tradeService.PaymentDomainService
	eventBus      event.EventBus
	validator     *validator.Validate
	unitOfWork    database.UnitOfWork
}

// NewPaymentApplicationService 创建支付应用服务
func NewPaymentApplicationService(
	paymentRepo repository.PaymentRepository,
	domainService *tradeService.PaymentDomainService,
	eventBus event.EventBus,
	validator *validator.Validate,
	unitOfWork database.UnitOfWork,
) *PaymentApplicationService {
	return &PaymentApplicationService{
		paymentRepo:   paymentRepo,
		domainService: domainService,
		eventBus:      eventBus,
		validator:     validator,
		unitOfWork:    unitOfWork,
	}
}

// CreatePayment 创建支付
func (s *PaymentApplicationService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest, userID string) (*dto.CreatePaymentResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 验证订单支付唯一性
	if err := s.domainService.ValidateOrderPaymentUniqueness(ctx, req.OrderID); err != nil {
		return nil, err
	}
	
	// 创建值对象
	paymentID, err := valueobject.NewPaymentID(generatePaymentID())
	if err != nil {
		return nil, fmt.Errorf("创建支付ID失败: %w", err)
	}
	
	amount, err := valueobject.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("创建金额失败: %w", err)
	}
	
	method, err := valueobject.NewPaymentMethod(req.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("创建支付方式失败: %w", err)
	}
	
	// 验证金额和支付方式
	if err := s.domainService.ValidatePaymentAmount(amount); err != nil {
		return nil, err
	}
	
	if err := s.domainService.ValidatePaymentMethod(method, amount); err != nil {
		return nil, err
	}
	
	// 创建支付聚合
	payment, err := entity.NewPayment(paymentID, req.OrderID, userID, amount, method, req.ClientIP)
	if err != nil {
		return nil, fmt.Errorf("创建支付失败: %w", err)
	}
	
	// 设置回调URL
	if req.ReturnURL != "" || req.NotifyURL != "" {
		payment.SetCallbackURLs(req.ReturnURL, req.NotifyURL)
	}
	
	// 设置扩展信息
	if req.Metadata != "" {
		payment.SetMetadata(req.Metadata)
	}
	
	var paymentURL, qrCode string
	
	// 在事务中保存支付
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}
		
		// 如果是第三方支付，生成支付信息
		if method.RequiresCallback() {
			paymentURL, qrCode, err = s.generatePaymentInfo(payment, method)
			if err != nil {
				return fmt.Errorf("生成支付信息失败: %w", err)
			}
			
			// 设置支付信息
			if err := payment.SetPaymentInfo(paymentURL, qrCode, ""); err != nil {
				return fmt.Errorf("设置支付信息失败: %w", err)
			}
			
			// 再次保存
			if err := s.paymentRepo.Save(ctx, payment); err != nil {
				return fmt.Errorf("保存支付信息失败: %w", err)
			}
		}
		
		// 发布领域事件
		events := payment.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		payment.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.CreatePaymentResponse{
		Payment:    *s.paymentToDTO(payment),
		PaymentURL: paymentURL,
		QRCode:     qrCode,
		Message:    "支付创建成功",
	}, nil
}

// ProcessPayment 处理支付
func (s *PaymentApplicationService) ProcessPayment(ctx context.Context, req dto.ProcessPaymentRequest, userID string) (*dto.ProcessPaymentResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}
	
	// 验证支付所有权
	payment, err := s.domainService.ValidatePaymentOwnership(ctx, paymentID, userID)
	if err != nil {
		return nil, err
	}
	
	// 验证支付是否可以处理
	if _, err := s.domainService.ValidatePaymentForProcessing(ctx, paymentID); err != nil {
		return nil, err
	}
	
	// 开始处理支付
	if err := payment.StartProcessing(); err != nil {
		return nil, fmt.Errorf("开始处理支付失败: %w", err)
	}
	
	// 在事务中保存
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}
		
		// 发布领域事件
		events := payment.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		payment.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.ProcessPaymentResponse{
		Payment: *s.paymentToDTO(payment),
		Message: "支付处理中",
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (s *PaymentApplicationService) HandlePaymentCallback(ctx context.Context, req dto.PaymentCallbackRequest) (*dto.PaymentCallbackResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}
	
	// 查找支付
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查找支付失败: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("支付不存在")
	}
	
	// 在事务中处理回调
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		switch req.Status {
		case "success":
			if err := payment.MarkAsSuccess(req.TransactionID); err != nil {
				return fmt.Errorf("标记支付成功失败: %w", err)
			}
		case "failed":
			reason := "支付失败"
			if data, ok := req.CallbackData["reason"].(string); ok {
				reason = data
			}
			if err := payment.MarkAsFailed(reason); err != nil {
				return fmt.Errorf("标记支付失败失败: %w", err)
			}
		case "cancelled":
			reason := "支付取消"
			if data, ok := req.CallbackData["reason"].(string); ok {
				reason = data
			}
			if err := payment.Cancel(reason); err != nil {
				return fmt.Errorf("取消支付失败: %w", err)
			}
		default:
			return fmt.Errorf("未知的回调状态: %s", req.Status)
		}
		
		// 保存支付
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}
		
		// 发布领域事件
		events := payment.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		payment.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return &dto.PaymentCallbackResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &dto.PaymentCallbackResponse{
		Success: true,
		Message: "回调处理成功",
	}, nil
}

// CancelPayment 取消支付
func (s *PaymentApplicationService) CancelPayment(ctx context.Context, req dto.CancelPaymentRequest, userID string) (*dto.CancelPaymentResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}
	
	// 验证支付所有权
	payment, err := s.domainService.ValidatePaymentOwnership(ctx, paymentID, userID)
	if err != nil {
		return nil, err
	}
	
	// 检查是否可以取消
	if err := s.domainService.CanCancelPayment(payment); err != nil {
		return nil, err
	}
	
	// 取消支付
	reason := req.Reason
	if reason == "" {
		reason = "用户取消"
	}
	
	if err := payment.Cancel(reason); err != nil {
		return nil, fmt.Errorf("取消支付失败: %w", err)
	}
	
	// 在事务中保存
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}
		
		// 发布领域事件
		events := payment.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		payment.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.CancelPaymentResponse{
		Payment: *s.paymentToDTO(payment),
		Message: "支付已取消",
	}, nil
}

// RefundPayment 退款
func (s *PaymentApplicationService) RefundPayment(ctx context.Context, req dto.RefundPaymentRequest, userID string) (*dto.RefundPaymentResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建支付ID
	paymentID, err := valueobject.NewPaymentID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}
	
	// 验证退款资格
	payment, err := s.domainService.ValidateRefundEligibility(ctx, paymentID, userID)
	if err != nil {
		return nil, err
	}
	
	// 执行退款
	if err := payment.MarkAsRefunded(); err != nil {
		return nil, fmt.Errorf("退款失败: %w", err)
	}
	
	// 在事务中保存
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			return fmt.Errorf("保存支付失败: %w", err)
		}
		
		// 发布领域事件
		events := payment.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		payment.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.RefundPaymentResponse{
		Payment: *s.paymentToDTO(payment),
		Message: "退款成功",
	}, nil
}

// GetPayment 获取支付详情
func (s *PaymentApplicationService) GetPayment(ctx context.Context, paymentID, userID string) (*dto.PaymentDTO, error) {
	id, err := valueobject.NewPaymentID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("无效的支付ID: %w", err)
	}
	
	payment, err := s.domainService.ValidatePaymentOwnership(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	
	return s.paymentToDTO(payment), nil
}

// ListPayments 获取支付列表
func (s *PaymentApplicationService) ListPayments(ctx context.Context, req dto.PaymentListRequest, userID string) (*dto.PaymentListResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	
	// 构建过滤器
	filters := repository.PaymentFilters{
		UserID:    userID,
		OrderID:   req.OrderID,
		Currency:  req.Currency,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Limit:     req.Size,
		Offset:    (req.Page - 1) * req.Size,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}
	
	// 处理状态过滤
	if req.Status != "" {
		status, err := valueobject.NewPaymentStatus(req.Status)
		if err != nil {
			return nil, fmt.Errorf("无效的支付状态: %w", err)
		}
		filters.Status = &status
	}
	
	// 处理支付方式过滤
	if req.PaymentMethod != "" {
		method, err := valueobject.NewPaymentMethod(req.PaymentMethod)
		if err != nil {
			return nil, fmt.Errorf("无效的支付方式: %w", err)
		}
		filters.Method = &method
	}
	
	// 处理金额过滤
	if req.MinAmount != nil {
		minAmount, err := valueobject.NewMoney(*req.MinAmount, filters.Currency)
		if err != nil {
			return nil, fmt.Errorf("无效的最小金额: %w", err)
		}
		filters.MinAmount = &minAmount
	}
	
	if req.MaxAmount != nil {
		maxAmount, err := valueobject.NewMoney(*req.MaxAmount, filters.Currency)
		if err != nil {
			return nil, fmt.Errorf("无效的最大金额: %w", err)
		}
		filters.MaxAmount = &maxAmount
	}
	
	// 查询支付列表
	payments, err := s.paymentRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("查询支付列表失败: %w", err)
	}
	
	// 统计总数
	total, err := s.paymentRepo.CountByUser(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("统计支付数量失败: %w", err)
	}
	
	// 转换为DTO
	paymentDTOs := make([]dto.PaymentDTO, 0, len(payments))
	for _, payment := range payments {
		paymentDTOs = append(paymentDTOs, *s.paymentToDTO(payment))
	}
	
	return &dto.PaymentListResponse{
		Payments: paymentDTOs,
		Total:    total,
		Page:     req.Page,
		Size:     req.Size,
	}, nil
}

// paymentToDTO 将支付实体转换为DTO
func (s *PaymentApplicationService) paymentToDTO(payment *entity.Payment) *dto.PaymentDTO {
	return &dto.PaymentDTO{
		ID:            payment.ID().Value(),
		OrderID:       payment.OrderID(),
		UserID:        payment.UserID(),
		Amount:        payment.Amount().Amount(),
		Currency:      payment.Amount().Currency(),
		DisplayAmount: payment.GetDisplayAmount(),
		PaymentMethod: payment.Method().Value(),
		Status:        payment.Status().Value(),
		TransactionID: payment.TransactionID(),
		PaymentURL:    payment.PaymentURL(),
		QRCode:        payment.QRCode(),
		ClientIP:      payment.ClientIP(),
		ReturnURL:     payment.ReturnURL(),
		NotifyURL:     payment.NotifyURL(),
		FailureReason: payment.FailureReason(),
		Metadata:      payment.Metadata(),
		PaymentTime:   payment.PaymentTime(),
		ExpiredAt:     payment.ExpiredAt(),
		CreatedAt:     payment.CreatedAt(),
		UpdatedAt:     payment.UpdatedAt(),
	}
}

// generatePaymentID 生成支付ID (简单实现，实际应使用UUID)
func generatePaymentID() string {
	return fmt.Sprintf("pay_%d", time.Now().UnixNano())
}

// generatePaymentInfo 生成支付信息（模拟第三方支付接口）
func (s *PaymentApplicationService) generatePaymentInfo(payment *entity.Payment, method valueobject.PaymentMethod) (paymentURL, qrCode string, err error) {
	// 这里应该调用真实的第三方支付接口
	switch {
	case method.IsAliPay():
		paymentURL = fmt.Sprintf("https://alipay.com/pay?id=%s", payment.ID().Value())
		qrCode = fmt.Sprintf("alipay_qr_%s", payment.ID().Value())
	case method.IsWeChatPay():
		paymentURL = fmt.Sprintf("https://pay.weixin.qq.com/pay?id=%s", payment.ID().Value())
		qrCode = fmt.Sprintf("wechat_qr_%s", payment.ID().Value())
	case method.IsPayPal():
		paymentURL = fmt.Sprintf("https://paypal.com/pay?id=%s", payment.ID().Value())
	case method.IsStripe():
		paymentURL = fmt.Sprintf("https://stripe.com/pay?id=%s", payment.ID().Value())
	}
	
	return paymentURL, qrCode, nil
} 