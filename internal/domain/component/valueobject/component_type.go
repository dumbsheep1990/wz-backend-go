package valueobject

import (
	"errors"
)

// ComponentType 组件类型值对象
type ComponentType struct {
	value string
}

// 组件类型常量
const (
	TypeHeader    = "header"    // 页头组件
	TypeNavbar    = "navbar"    // 导航栏组件
	TypeHero      = "hero"      // 英雄区组件
	TypeFeature   = "feature"   // 特性展示组件
	TypeContent   = "content"   // 内容组件
	TypeFooter    = "footer"    // 页脚组件
	TypeSidebar   = "sidebar"   // 侧边栏组件
	TypeCard      = "card"      // 卡片组件
	TypeGallery   = "gallery"   // 图片画廊组件
	TypeContact   = "contact"   // 联系表单组件
	TypeCustom    = "custom"    // 自定义组件
)

// NewComponentType 创建组件类型
func NewComponentType(value string) (ComponentType, error) {
	if err := validateComponentType(value); err != nil {
		return ComponentType{}, err
	}
	return ComponentType{value: value}, nil
}

// NewHeaderType 创建页头类型
func NewHeaderType() ComponentType {
	return ComponentType{value: TypeHeader}
}

// NewNavbarType 创建导航栏类型
func NewNavbarType() ComponentType {
	return ComponentType{value: TypeNavbar}
}

// NewHeroType 创建英雄区类型
func NewHeroType() ComponentType {
	return ComponentType{value: TypeHero}
}

// NewContentType 创建内容类型
func NewContentType() ComponentType {
	return ComponentType{value: TypeContent}
}

// NewFooterType 创建页脚类型
func NewFooterType() ComponentType {
	return ComponentType{value: TypeFooter}
}

// Value 获取类型值
func (t ComponentType) Value() string {
	return t.value
}

// IsHeader 是否为页头组件
func (t ComponentType) IsHeader() bool {
	return t.value == TypeHeader
}

// IsNavbar 是否为导航栏组件
func (t ComponentType) IsNavbar() bool {
	return t.value == TypeNavbar
}

// IsHero 是否为英雄区组件
func (t ComponentType) IsHero() bool {
	return t.value == TypeHero
}

// IsContent 是否为内容组件
func (t ComponentType) IsContent() bool {
	return t.value == TypeContent
}

// IsFooter 是否为页脚组件
func (t ComponentType) IsFooter() bool {
	return t.value == TypeFooter
}

// IsCustom 是否为自定义组件
func (t ComponentType) IsCustom() bool {
	return t.value == TypeCustom
}

// Equals 判断两个类型是否相等
func (t ComponentType) Equals(other ComponentType) bool {
	return t.value == other.value
}

// String 字符串表示
func (t ComponentType) String() string {
	return t.value
}

// validateComponentType 验证组件类型
func validateComponentType(value string) error {
	validTypes := []string{
		TypeHeader, TypeNavbar, TypeHero, TypeFeature, TypeContent,
		TypeFooter, TypeSidebar, TypeCard, TypeGallery, TypeContact, TypeCustom,
	}
	
	for _, validType := range validTypes {
		if value == validType {
			return nil
		}
	}
	
	return errors.New("无效的组件类型")
} 