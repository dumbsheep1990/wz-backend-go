package valueobject

import (
	"errors"
	"strings"
)

// ComponentID 组件ID值对象
type ComponentID struct {
	value string
}

// NewComponentID 创建组件ID
func NewComponentID(value string) (ComponentID, error) {
	if err := validateComponentID(value); err != nil {
		return ComponentID{}, err
	}
	return ComponentID{value: value}, nil
}

// NewComponentIDFromString 从字符串创建组件ID（用于数据库映射）
func NewComponentIDFromString(value string) ComponentID {
	return ComponentID{value: value}
}

// Value 获取ID值
func (id ComponentID) Value() string {
	return id.value
}

// IsEmpty 判断ID是否为空
func (id ComponentID) IsEmpty() bool {
	return id.value == ""
}

// Equals 判断两个ID是否相等
func (id ComponentID) Equals(other ComponentID) bool {
	return id.value == other.value
}

// String 字符串表示
func (id ComponentID) String() string {
	return id.value
}

// validateComponentID 验证组件ID
func validateComponentID(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("组件ID不能为空")
	}
	if len(value) > 50 {
		return errors.New("组件ID长度不能超过50个字符")
	}
	return nil
} 