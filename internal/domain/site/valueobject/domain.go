package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Domain 域名值对象
type Domain struct {
	value string
}

// NewDomain 创建域名
func NewDomain(value string) (Domain, error) {
	if err := validateDomain(value); err != nil {
		return Domain{}, err
	}
	return Domain{value: strings.ToLower(strings.TrimSpace(value))}, nil
}

// Value 获取域名值
func (d Domain) Value() string {
	return d.value
}

// IsEmpty 判断域名是否为空
func (d Domain) IsEmpty() bool {
	return d.value == ""
}

// Equals 判断两个域名是否相等
func (d Domain) Equals(other Domain) bool {
	return d.value == other.value
}

// String 字符串表示
func (d Domain) String() string {
	return d.value
}

// GetSubdomain 获取子域名部分
func (d Domain) GetSubdomain() string {
	parts := strings.Split(d.value, ".")
	if len(parts) > 2 {
		return parts[0]
	}
	return ""
}

// GetMainDomain 获取主域名部分
func (d Domain) GetMainDomain() string {
	parts := strings.Split(d.value, ".")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return d.value
}

// validateDomain 验证域名
func validateDomain(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("域名不能为空")
	}
	
	if len(trimmed) > 253 {
		return errors.New("域名长度不能超过253个字符")
	}
	
	// 基本的域名格式验证
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !domainRegex.MatchString(trimmed) {
		return errors.New("域名格式无效")
	}
	
	return nil
} 