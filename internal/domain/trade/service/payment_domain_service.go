package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wz-backend-go/internal/domain/trade/entity"
	"wz-backend-go/internal/domain/trade/repository"
	"wz-backend-go/internal/domain/trade/valueobject"
)

// PaymentDomainService 支付领域服务
type PaymentDomainService struct {
	paymentRepo repository.PaymentRepository
}

// NewPaymentDomainService 创建支付领域服务
func NewPaymentDomainService(paymentRepo repository.PaymentRepository) *PaymentDomainService {
	return &PaymentDomainService{
		paymentRepo: paymentRepo,
	}
}

// ValidatePaymentOwnership 验证支付所有权
func (s *PaymentDomainService) ValidatePaymentOwnership(ctx context.Context, paymentID valueobject.PaymentID, userID string) (*entity.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查找支付失败: %w", err)
	}
	if payment == nil {
		return nil, errors.New("支付不存在")
	}
	
	if !payment.IsOwnedBy(userID) {
		return nil, errors.New("无权限操作此支付")
	}
	
	return payment, nil
}

// ValidatePaymentForProcessing 验证支付是否可以处理
func (s *PaymentDomainService) ValidatePaymentForProcessing(ctx context.Context, paymentID valueobject.PaymentID) (*entity.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("查找支付失败: %w", err)
	}
	if payment == nil {
		return nil, errors.New("支付不存在")
	}
	
	// 检查支付状态
	if !payment.Status().IsPending() {
		return nil, errors.New("支付状态不允许处理")
	}
	
	// 检查是否过期
	if payment.IsExpired() {
		return nil, errors.New("支付已过期")
	}
	
	return payment, nil
}

// ValidateOrderPaymentUniqueness 验证订单支付唯一性
func (s *PaymentDomainService) ValidateOrderPaymentUniqueness(ctx context.Context, orderID string) error {
	// 检查订单是否已经有成功的支付
	payments, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("查询订单支付失败: %w", err)
	}
	
	for _, payment := range payments {
		if payment.Status().IsSuccess() {
			return errors.New("订单已经有成功的支付")
		}
	}
	
	return nil
}

// ValidateRefundEligibility 验证退款资格
func (s *PaymentDomainService) ValidateRefundEligibility(ctx context.Context, paymentID valueobject.PaymentID, userID string) (*entity.Payment, error) {
	payment, err := s.ValidatePaymentOwnership(ctx, paymentID, userID)
	if err != nil {
		return nil, err
	}
	
	if !payment.Status().IsSuccess() {
		return nil, errors.New("只有成功的支付才能退款")
	}
	
	// 检查退款时限（比如30天内）
	if payment.PaymentTime() != nil {
		refundDeadline := payment.PaymentTime().Add(30 * 24 * time.Hour) // 30天
		if time.Now().After(refundDeadline) {
			return nil, errors.New("已超过退款时限")
		}
	}
	
	return payment, nil
}

// CalculatePaymentStatistics 计算支付统计
func (s *PaymentDomainService) CalculatePaymentStatistics(ctx context.Context, userID string, filters repository.PaymentFilters) (*repository.PaymentStatistics, error) {
	// 设置用户过滤器
	filters.UserID = userID
	
	// 获取总数和总金额
	totalCount, err := s.paymentRepo.CountByUser(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("统计总支付数失败: %w", err)
	}
	
	// 按状态统计
	successStatus := valueobject.NewSuccessStatus()
	successCount, err := s.paymentRepo.CountByStatus(ctx, successStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("统计成功支付数失败: %w", err)
	}
	
	failedStatus := valueobject.NewFailedStatus()
	failedCount, err := s.paymentRepo.CountByStatus(ctx, failedStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("统计失败支付数失败: %w", err)
	}
	
	pendingStatus := valueobject.NewPendingStatus()
	pendingCount, err := s.paymentRepo.CountByStatus(ctx, pendingStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("统计待支付数失败: %w", err)
	}
	
	refundedStatus := valueobject.NewRefundedStatus()
	refundedCount, err := s.paymentRepo.CountByStatus(ctx, refundedStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("统计退款数失败: %w", err)
	}
	
	// 这里简化处理，实际应该查询具体的金额统计
	totalAmount, _ := valueobject.NewZeroMoney("CNY")
	successAmount, _ := valueobject.NewZeroMoney("CNY")
	refundedAmount, _ := valueobject.NewZeroMoney("CNY")
	
	return &repository.PaymentStatistics{
		TotalCount:     totalCount,
		TotalAmount:    totalAmount,
		SuccessCount:   successCount,
		SuccessAmount:  successAmount,
		FailedCount:    failedCount,
		PendingCount:   pendingCount,
		RefundedCount:  refundedCount,
		RefundedAmount: refundedAmount,
		MethodStats:    make(map[string]repository.PaymentMethodStatistics),
		DailyStats:     make([]repository.DailyPaymentStatistics, 0),
	}, nil
}

// ProcessExpiredPayments 处理过期支付
func (s *PaymentDomainService) ProcessExpiredPayments(ctx context.Context) error {
	// 查找过期的支付
	now := time.Now()
	expiredPayments, err := s.paymentRepo.FindExpiredPayments(ctx, now)
	if err != nil {
		return fmt.Errorf("查找过期支付失败: %w", err)
	}
	
	for _, payment := range expiredPayments {
		// 检查并设置为过期状态
		if err := payment.CheckExpiration(); err != nil {
			continue // 记录日志但继续处理其他支付
		}
		
		// 保存更新
		if err := s.paymentRepo.Save(ctx, payment); err != nil {
			// 记录日志但继续处理其他支付
			continue
		}
	}
	
	return nil
}

// ValidatePaymentAmount 验证支付金额
func (s *PaymentDomainService) ValidatePaymentAmount(amount valueobject.Money) error {
	if !amount.IsPositive() {
		return errors.New("支付金额必须大于零")
	}
	
	// 检查金额范围（比如最小1分，最大100万）
	minAmount, _ := valueobject.NewMoneyFromCents(1, amount.Currency())
	if isLess, _ := amount.LessThan(minAmount); isLess {
		return errors.New("支付金额过小")
	}
	
	maxAmount, _ := valueobject.NewMoney(1000000, amount.Currency()) // 100万
	if isGreater, _ := amount.GreaterThan(maxAmount); isGreater {
		return errors.New("支付金额过大")
	}
	
	return nil
}

// ValidatePaymentMethod 验证支付方式
func (s *PaymentDomainService) ValidatePaymentMethod(method valueobject.PaymentMethod, amount valueobject.Money) error {
	// 根据支付方式验证金额限制
	switch {
	case method.IsBalance():
		// 余额支付可能有额度限制
		maxBalanceAmount, _ := valueobject.NewMoney(50000, amount.Currency()) // 5万限额
		if isGreater, _ := amount.GreaterThan(maxBalanceAmount); isGreater {
			return errors.New("余额支付金额超过限额")
		}
	case method.IsAliPay() || method.IsWeChatPay():
		// 移动支付限额
		maxMobileAmount, _ := valueobject.NewMoney(500000, amount.Currency()) // 50万限额
		if isGreater, _ := amount.GreaterThan(maxMobileAmount); isGreater {
			return errors.New("移动支付金额超过限额")
		}
	}
	
	return nil
}

// CanCancelPayment 检查是否可以取消支付
func (s *PaymentDomainService) CanCancelPayment(payment *entity.Payment) error {
	if payment.Status().IsFinal() {
		return errors.New("支付已完成，无法取消")
	}
	
	return nil
}

// GetRetryablePayments 获取可重试的支付
func (s *PaymentDomainService) GetRetryablePayments(ctx context.Context, userID string) ([]*entity.Payment, error) {
	// 查找失败和过期的支付
	failedStatus := valueobject.NewFailedStatus()
	filters := repository.PaymentFilters{
		UserID: userID,
		Status: &failedStatus,
		Limit:  50, // 限制数量
	}
	
	failedPayments, err := s.paymentRepo.FindByStatus(ctx, failedStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("查找失败支付失败: %w", err)
	}
	
	expiredStatus := valueobject.NewExpiredStatus()
	filters.Status = &expiredStatus
	expiredPayments, err := s.paymentRepo.FindByStatus(ctx, expiredStatus, filters)
	if err != nil {
		return nil, fmt.Errorf("查找过期支付失败: %w", err)
	}
	
	// 合并结果
	retryablePayments := make([]*entity.Payment, 0, len(failedPayments)+len(expiredPayments))
	retryablePayments = append(retryablePayments, failedPayments...)
	retryablePayments = append(retryablePayments, expiredPayments...)
	
	return retryablePayments, nil
} 