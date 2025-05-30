package valueobject

import (
	"errors"
	"strings"
)

// PaymentID 支付ID值对象
type PaymentID struct {
	value string
}

// NewPaymentID 创建支付ID
func NewPaymentID(value string) (PaymentID, error) {
	if err := validatePaymentID(value); err != nil {
		return PaymentID{}, err
	}
	return PaymentID{value: value}, nil
}

// NewPaymentIDFromString 从字符串创建支付ID（用于数据库映射）
func NewPaymentIDFromString(value string) PaymentID {
	return PaymentID{value: value}
}

// Value 获取ID值
func (id PaymentID) Value() string {
	return id.value
}

// IsEmpty 判断ID是否为空
func (id PaymentID) IsEmpty() bool {
	return id.value == ""
}

// Equals 判断两个ID是否相等
func (id PaymentID) Equals(other PaymentID) bool {
	return id.value == other.value
}

// String 字符串表示
func (id PaymentID) String() string {
	return id.value
}

// validatePaymentID 验证支付ID
func validatePaymentID(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("支付ID不能为空")
	}
	if len(value) > 50 {
		return errors.New("支付ID长度不能超过50个字符")
	}
	return nil
} 