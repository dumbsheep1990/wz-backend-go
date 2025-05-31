package valueobject

import (
	"errors"
	"strconv"
)

// ProductID 产品ID值对象
type ProductID struct {
	value int64
}

// NewProductID 创建产品ID值对象
func NewProductID(id int64) ProductID {
	return ProductID{value: id}
}

// NewProductIDFromString 从字符串创建产品ID值对象
func NewProductIDFromString(idStr string) (ProductID, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ProductID{}, errors.New("无效的产品ID格式")
	}
	if id <= 0 {
		return ProductID{}, errors.New("产品ID必须大于0")
	}
	return ProductID{value: id}, nil
}

// MustNewProductID 创建产品ID值对象，如果无效则panic
func MustNewProductID(id int64) ProductID {
	if id <= 0 {
		panic("产品ID必须大于0")
	}
	return ProductID{value: id}
}

// Value 获取产品ID值
func (p ProductID) Value() int64 {
	return p.value
}

// String 获取产品ID的字符串表示
func (p ProductID) String() string {
	return strconv.FormatInt(p.value, 10)
}

// IsValid 检查产品ID是否有效
func (p ProductID) IsValid() bool {
	return p.value > 0
}

// Equals 比较两个产品ID是否相等
func (p ProductID) Equals(other ProductID) bool {
	return p.value == other.value
}

// ValidateProductID 验证产品ID
func ValidateProductID(id int64) error {
	if id <= 0 {
		return errors.New("产品ID必须大于0")
	}
	return nil
}
