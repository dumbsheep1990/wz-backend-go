package models

import "time"

// Page 页面模型
type Page struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	SiteID      string    `json:"siteId" gorm:"index"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug" gorm:"index"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    []string  `json:"keywords" gorm:"type:json"`
	IsHomepage  bool      `json:"isHomepage"`
	Layout      string    `json:"layout"` // default, full-width, sidebar
	Sections    []Section `json:"sections" gorm:"-"` // 不存储在同一表
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	SortOrder   int       `json:"sortOrder" gorm:"index"`
}

// Section 页面区块模型
type Section struct {
	ID         string      `json:"id" gorm:"primaryKey"`
	PageID     string      `json:"pageId" gorm:"index"`
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	Settings   interface{} `json:"settings" gorm:"type:json"`
	Components []Component `json:"components" gorm:"-"` // 不存储在同一表
	Style      interface{} `json:"style" gorm:"type:json"`
	SortOrder  int         `json:"sortOrder" gorm:"index"`
}

// Component 组件模型
type Component struct {
	ID        string      `json:"id" gorm:"primaryKey"`
	SectionID string      `json:"sectionId" gorm:"index"`
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Settings  interface{} `json:"settings" gorm:"type:json"`
	Content   interface{} `json:"content" gorm:"type:json"`
	Style     interface{} `json:"style" gorm:"type:json"`
	SortOrder int         `json:"sortOrder" gorm:"index"`
}

// Navigation 导航配置
type Navigation struct {
	Type  string          `json:"type"` // horizontal, vertical, mega-menu
	Items []NavigationItem `json:"items"`
	Style interface{}     `json:"style"`
}

// NavigationItem 导航项
type NavigationItem struct {
	ID            string           `json:"id"`
	Label         string           `json:"label"`
	Link          string           `json:"link"`
	Icon          string           `json:"icon,omitempty"`
	Children      []NavigationItem `json:"children,omitempty"`
	IsExternalLink bool            `json:"isExternalLink"`
}

// ComponentCategory 组件分类
type ComponentCategory struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Components  []ComponentDefinition `json:"components"`
}

// ComponentDefinition 组件定义
type ComponentDefinition struct {
	Type            string      `json:"type"`
	Name            string      `json:"name"`
	Icon            string      `json:"icon"`
	Description     string      `json:"description"`
	DefaultSettings interface{} `json:"defaultSettings"`
} 