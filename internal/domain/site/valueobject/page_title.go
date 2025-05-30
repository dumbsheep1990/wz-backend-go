package valueobject

import (
	"errors"
	"strings"
)

// PageTitle 页面标题值对象
type PageTitle struct {
	value string
}

// NewPageTitle 创建页面标题
func NewPageTitle(value string) (PageTitle, error) {
	if err := validatePageTitle(value); err != nil {
		return PageTitle{}, err
	}
	return PageTitle{value: strings.TrimSpace(value)}, nil
}

// Value 获取标题值
func (t PageTitle) Value() string {
	return t.value
}

// IsEmpty 判断标题是否为空
func (t PageTitle) IsEmpty() bool {
	return t.value == ""
}

// Equals 判断两个标题是否相等
func (t PageTitle) Equals(other PageTitle) bool {
	return t.value == other.value
}

// String 字符串表示
func (t PageTitle) String() string {
	return t.value
}

// validatePageTitle 验证页面标题
func validatePageTitle(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("页面标题不能为空")
	}
	if len(trimmed) < 1 {
		return errors.New("页面标题至少需要1个字符")
	}
	if len(trimmed) > 200 {
		return errors.New("页面标题长度不能超过200个字符")
	}
	return nil
} 