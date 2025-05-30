package valueobject

import (
	"errors"
	"strings"
)

// PageID 页面ID值对象
type PageID struct {
	value string
}

// NewPageID 创建页面ID
func NewPageID(value string) (PageID, error) {
	if err := validatePageID(value); err != nil {
		return PageID{}, err
	}
	return PageID{value: value}, nil
}

// NewPageIDFromString 从字符串创建页面ID（用于数据库映射）
func NewPageIDFromString(value string) PageID {
	return PageID{value: value}
}

// Value 获取ID值
func (id PageID) Value() string {
	return id.value
}

// IsEmpty 判断ID是否为空
func (id PageID) IsEmpty() bool {
	return id.value == ""
}

// Equals 判断两个ID是否相等
func (id PageID) Equals(other PageID) bool {
	return id.value == other.value
}

// String 字符串表示
func (id PageID) String() string {
	return id.value
}

// validatePageID 验证页面ID
func validatePageID(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("页面ID不能为空")
	}
	if len(value) > 50 {
		return errors.New("页面ID长度不能超过50个字符")
	}
	return nil
} 