package page

import (
	"time"
)

// Page 页面模型
type Page struct {
	ID          string    `json:"id" gorm:"primaryKey" db:"id"`
	SiteID      string    `json:"siteId" gorm:"index" db:"site_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" gorm:"index" db:"slug"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Keywords    []string  `json:"keywords" gorm:"type:json" db:"keywords"`
	IsHomepage  bool      `json:"isHomepage" db:"is_homepage"`
	Layout      string    `json:"layout" db:"layout"` // default, full-width, sidebar
	Sections    []Section `json:"sections" gorm:"-" db:"-"` // 不存储在同一表
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime" db:"updated_at"`
	SortOrder   int       `json:"sortOrder" gorm:"index" db:"sort_order"`
}

// Section 页面区块模型
type Section struct {
	ID         string      `json:"id" gorm:"primaryKey" db:"id"`
	PageID     string      `json:"pageId" gorm:"index" db:"page_id"`
	Type       string      `json:"type" db:"type"`
	Title      string      `json:"title" db:"title"`
	Settings   interface{} `json:"settings" gorm:"type:json" db:"settings"`
	Components []Component `json:"components" gorm:"-" db:"-"` // 不存储在同一表
	Style      interface{} `json:"style" gorm:"type:json" db:"style"`
	SortOrder  int         `json:"sortOrder" gorm:"index" db:"sort_order"`
}

// Component 组件模型
type Component struct {
	ID        string      `json:"id" gorm:"primaryKey" db:"id"`
	SectionID string      `json:"sectionId" gorm:"index" db:"section_id"`
	Type      string      `json:"type" db:"type"`
	Name      string      `json:"name" db:"name"`
	Settings  interface{} `json:"settings" gorm:"type:json" db:"settings"`
	Content   interface{} `json:"content" gorm:"type:json" db:"content"`
	Style     interface{} `json:"style" gorm:"type:json" db:"style"`
	SortOrder int         `json:"sortOrder" gorm:"index" db:"sort_order"`
}

// ComponentCategory 组件分类
type ComponentCategory struct {
	ID          string                `json:"id" db:"id"`
	Name        string                `json:"name" db:"name"`
	Components  []ComponentDefinition `json:"components" db:"components"`
}

// ComponentDefinition 组件定义
type ComponentDefinition struct {
	Type            string      `json:"type" db:"type"`
	Name            string      `json:"name" db:"name"`
	Icon            string      `json:"icon" db:"icon"`
	Description     string      `json:"description" db:"description"`
	DefaultSettings interface{} `json:"defaultSettings" db:"default_settings"`
}

// PageVersion 页面版本
type PageVersion struct {
	ID          string    `json:"id" gorm:"primaryKey" db:"id"`
	PageID      string    `json:"pageId" gorm:"index" db:"page_id"`
	Version     int       `json:"version" db:"version"`
	Content     string    `json:"content" gorm:"type:json" db:"content"`
	Description string    `json:"description" db:"description"`
	CreatedBy   string    `json:"createdBy" db:"created_by"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime" db:"created_at"`
	IsPublished bool      `json:"isPublished" db:"is_published"`
}

// PageTemplate 页面模板
type PageTemplate struct {
	ID          string    `json:"id" gorm:"primaryKey" db:"id"`
	Name        string    `json:"name" db:"name"`
	Thumbnail   string    `json:"thumbnail" db:"thumbnail"`
	Description string    `json:"description" db:"description"`
	Content     string    `json:"content" gorm:"type:json" db:"content"`
	Type        string    `json:"type" db:"type"` // general, landing, product, blog
	CreatedBy   string    `json:"createdBy" db:"created_by"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime" db:"updated_at"`
}

// PageAnalytics 页面分析数据
type PageAnalytics struct {
	ID                string    `json:"id" gorm:"primaryKey" db:"id"`
	PageID            string    `json:"pageId" gorm:"index" db:"page_id"`
	SiteID            string    `json:"siteId" gorm:"index" db:"site_id"`
	Views             int       `json:"views" db:"views"`
	UniqueVisitors    int       `json:"uniqueVisitors" db:"unique_visitors"`
	AverageTimeOnPage int       `json:"averageTimeOnPage" db:"average_time_on_page"` // 秒
	BounceRate        float64   `json:"bounceRate" db:"bounce_rate"`
	ExitRate          float64   `json:"exitRate" db:"exit_rate"`
	Date              time.Time `json:"date" db:"date"`
}
