package valueobject

import (
	"errors"
	"strings"
)

// SiteName 站点名称值对象
type SiteName struct {
	value string
}

// NewSiteName 创建站点名称
func NewSiteName(value string) (SiteName, error) {
	if err := validateSiteName(value); err != nil {
		return SiteName{}, err
	}
	return SiteName{value: strings.TrimSpace(value)}, nil
}

// Value 获取名称值
func (n SiteName) Value() string {
	return n.value
}

// IsEmpty 判断名称是否为空
func (n SiteName) IsEmpty() bool {
	return n.value == ""
}

// Equals 判断两个名称是否相等
func (n SiteName) Equals(other SiteName) bool {
	return n.value == other.value
}

// String 字符串表示
func (n SiteName) String() string {
	return n.value
}

// validateSiteName 验证站点名称
func validateSiteName(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("站点名称不能为空")
	}
	if len(trimmed) < 2 {
		return errors.New("站点名称至少需要2个字符")
	}
	if len(trimmed) > 100 {
		return errors.New("站点名称长度不能超过100个字符")
	}
	return nil
} 