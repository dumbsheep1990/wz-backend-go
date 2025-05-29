package valueobject

import (
	"errors"
	"fmt"
)

// OrderItemID 订单项ID值对象
type OrderItemID string

// NewOrderItemID 创建订单项ID值对象
func NewOrderItemID(id string) OrderItemID {
	return OrderItemID(id)
}

// MustNewOrderItemID 创建订单项ID值对象，如果无效则panic
func MustNewOrderItemID(id string) OrderItemID {
	if id == "" {
		panic("订单项ID不能为空")
	}
	return OrderItemID(id)
}

// ValidateOrderItemID 验证订单项ID
func ValidateOrderItemID(id string) error {
	if id == "" {
		return errors.New("订单项ID不能为空")
	}
	return nil
}

// Value 获取ID值
func (id OrderItemID) Value() string {
	return string(id)
}

// String 字符串表示
func (id OrderItemID) String() string {
	return fmt.Sprintf("%s", id)
}
