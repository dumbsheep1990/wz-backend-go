package valueobject

import (
	"errors"
	"strings"
)

// SiteID 站点ID值对象
type SiteID struct {
	value string
}

// NewSiteID 创建站点ID
func NewSiteID(value string) (SiteID, error) {
	if err := validateSiteID(value); err != nil {
		return SiteID{}, err
	}
	return SiteID{value: value}, nil
}

// NewSiteIDFromString 从字符串创建站点ID（用于数据库映射）
func NewSiteIDFromString(value string) SiteID {
	return SiteID{value: value}
}

// Value 获取ID值
func (id SiteID) Value() string {
	return id.value
}

// IsEmpty 判断ID是否为空
func (id SiteID) IsEmpty() bool {
	return id.value == ""
}

// Equals 判断两个ID是否相等
func (id SiteID) Equals(other SiteID) bool {
	return id.value == other.value
}

// String 字符串表示
func (id SiteID) String() string {
	return id.value
}

// validateSiteID 验证站点ID
func validateSiteID(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("站点ID不能为空")
	}
	if len(value) > 50 {
		return errors.New("站点ID长度不能超过50个字符")
	}
	return nil
} 