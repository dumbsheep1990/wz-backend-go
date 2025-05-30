package valueobject

import (
	"errors"
)

// PaymentStatus 支付状态值对象
type PaymentStatus struct {
	value string
}

// 支付状态常量
const (
	StatusPending   = "pending"   // 待支付
	StatusProcessing = "processing" // 处理中
	StatusSuccess   = "success"   // 支付成功
	StatusFailed    = "failed"    // 支付失败
	StatusCancelled = "cancelled" // 已取消
	StatusRefunded  = "refunded"  // 已退款
	StatusExpired   = "expired"   // 已过期
)

// NewPaymentStatus 创建支付状态
func NewPaymentStatus(value string) (PaymentStatus, error) {
	if err := validatePaymentStatus(value); err != nil {
		return PaymentStatus{}, err
	}
	return PaymentStatus{value: value}, nil
}

// NewPendingStatus 创建待支付状态
func NewPendingStatus() PaymentStatus {
	return PaymentStatus{value: StatusPending}
}

// NewProcessingStatus 创建处理中状态
func NewProcessingStatus() PaymentStatus {
	return PaymentStatus{value: StatusProcessing}
}

// NewSuccessStatus 创建支付成功状态
func NewSuccessStatus() PaymentStatus {
	return PaymentStatus{value: StatusSuccess}
}

// NewFailedStatus 创建支付失败状态
func NewFailedStatus() PaymentStatus {
	return PaymentStatus{value: StatusFailed}
}

// NewCancelledStatus 创建已取消状态
func NewCancelledStatus() PaymentStatus {
	return PaymentStatus{value: StatusCancelled}
}

// NewRefundedStatus 创建已退款状态
func NewRefundedStatus() PaymentStatus {
	return PaymentStatus{value: StatusRefunded}
}

// NewExpiredStatus 创建已过期状态
func NewExpiredStatus() PaymentStatus {
	return PaymentStatus{value: StatusExpired}
}

// Value 获取状态值
func (s PaymentStatus) Value() string {
	return s.value
}

// IsPending 是否为待支付状态
func (s PaymentStatus) IsPending() bool {
	return s.value == StatusPending
}

// IsProcessing 是否为处理中状态
func (s PaymentStatus) IsProcessing() bool {
	return s.value == StatusProcessing
}

// IsSuccess 是否为支付成功状态
func (s PaymentStatus) IsSuccess() bool {
	return s.value == StatusSuccess
}

// IsFailed 是否为支付失败状态
func (s PaymentStatus) IsFailed() bool {
	return s.value == StatusFailed
}

// IsCancelled 是否为已取消状态
func (s PaymentStatus) IsCancelled() bool {
	return s.value == StatusCancelled
}

// IsRefunded 是否为已退款状态
func (s PaymentStatus) IsRefunded() bool {
	return s.value == StatusRefunded
}

// IsExpired 是否为已过期状态
func (s PaymentStatus) IsExpired() bool {
	return s.value == StatusExpired
}

// IsFinal 是否为最终状态（不可再变更）
func (s PaymentStatus) IsFinal() bool {
	return s.value == StatusSuccess || s.value == StatusFailed || 
		   s.value == StatusCancelled || s.value == StatusRefunded || 
		   s.value == StatusExpired
}

// CanTransitionTo 是否可以转换到指定状态
func (s PaymentStatus) CanTransitionTo(target PaymentStatus) bool {
	switch s.value {
	case StatusPending:
		return target.value == StatusProcessing || target.value == StatusSuccess || 
			   target.value == StatusFailed || target.value == StatusCancelled || 
			   target.value == StatusExpired
	case StatusProcessing:
		return target.value == StatusSuccess || target.value == StatusFailed || 
			   target.value == StatusCancelled
	case StatusSuccess:
		return target.value == StatusRefunded
	case StatusFailed, StatusCancelled, StatusRefunded, StatusExpired:
		return false // 最终状态不能再变更
	default:
		return false
	}
}

// Equals 判断两个状态是否相等
func (s PaymentStatus) Equals(other PaymentStatus) bool {
	return s.value == other.value
}

// String 字符串表示
func (s PaymentStatus) String() string {
	return s.value
}

// validatePaymentStatus 验证支付状态
func validatePaymentStatus(value string) error {
	validStatuses := []string{
		StatusPending, StatusProcessing, StatusSuccess, 
		StatusFailed, StatusCancelled, StatusRefunded, StatusExpired,
	}
	
	for _, validStatus := range validStatuses {
		if value == validStatus {
			return nil
		}
	}
	
	return errors.New("无效的支付状态")
} 