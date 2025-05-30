package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// PageSlug 页面路径值对象
type PageSlug struct {
	value string
}

// NewPageSlug 创建页面路径
func NewPageSlug(value string) (PageSlug, error) {
	if err := validatePageSlug(value); err != nil {
		return PageSlug{}, err
	}
	return PageSlug{value: strings.ToLower(strings.TrimSpace(value))}, nil
}

// Value 获取路径值
func (s PageSlug) Value() string {
	return s.value
}

// IsEmpty 判断路径是否为空
func (s PageSlug) IsEmpty() bool {
	return s.value == ""
}

// Equals 判断两个路径是否相等
func (s PageSlug) Equals(other PageSlug) bool {
	return s.value == other.value
}

// String 字符串表示
func (s PageSlug) String() string {
	return s.value
}

// GetFullPath 获取完整路径
func (s PageSlug) GetFullPath() string {
	if s.value == "" || s.value == "index" || s.value == "home" {
		return "/"
	}
	return "/" + s.value
}

// validatePageSlug 验证页面路径
func validatePageSlug(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("页面路径不能为空")
	}
	
	if len(trimmed) > 100 {
		return errors.New("页面路径长度不能超过100个字符")
	}
	
	// 页面路径只能包含字母、数字、连字符和下划线
	slugRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !slugRegex.MatchString(trimmed) {
		return errors.New("页面路径只能包含字母、数字、连字符和下划线")
	}
	
	// 不能以连字符开始或结束
	if strings.HasPrefix(trimmed, "-") || strings.HasSuffix(trimmed, "-") {
		return errors.New("页面路径不能以连字符开始或结束")
	}
	
	return nil
} 