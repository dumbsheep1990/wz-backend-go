package valueobject

import (
	"errors"
	"fmt"
)

// OrderID 订单ID值对象
type OrderID string

// NewOrderID 创建订单ID值对象
func NewOrderID(id string) OrderID {
	return OrderID(id)
}

// MustNewOrderID 创建订单ID值对象，如果无效则panic
func MustNewOrderID(id string) OrderID {
	if id == "" {
		panic("订单ID不能为空")
	}
	return OrderID(id)
}

// ValidateOrderID 验证订单ID
func ValidateOrderID(id string) error {
	if id == "" {
		return errors.New("订单ID不能为空")
	}
	return nil
}

// Value 获取ID值
func (id OrderID) Value() string {
	return string(id)
}

// String 字符串表示
func (id OrderID) String() string {
	return fmt.Sprintf("%s", id)
}
