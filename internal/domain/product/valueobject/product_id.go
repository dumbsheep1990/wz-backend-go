package valueobject

import (
	"errors"
	"fmt"
)

// ProductID 商品ID值对象
type ProductID int64

// NewProductID 创建商品ID值对象
func NewProductID(id int64) ProductID {
	return ProductID(id)
}

// MustNewProductID 创建商品ID值对象，如果无效则panic
func MustNewProductID(id int64) ProductID {
	if id <= 0 {
		panic("商品ID必须大于0")
	}
	return ProductID(id)
}

// ValidateProductID 验证商品ID
func ValidateProductID(id int64) error {
	if id <= 0 {
		return errors.New("商品ID必须大于0")
	}
	return nil
}

// Value 获取ID值
func (id ProductID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id ProductID) String() string {
	return fmt.Sprintf("%d", id)
}
