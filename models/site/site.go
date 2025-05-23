package site

import (
	"time"

	"wz-project/wz-backend-go/models/common"
	"wz-project/wz-backend-go/models/page"
)

// Site 站点模型
type Site struct {
	ID          string       `json:"id" gorm:"primaryKey" db:"id"`
	Name        string       `json:"name" gorm:"not null" db:"name"`
	Description string       `json:"description" db:"description"`
	Domain      string       `json:"domain" gorm:"index" db:"domain"`
	Logo        string       `json:"logo" db:"logo"`
	Favicon     string       `json:"favicon" db:"favicon"`
	TenantID    string       `json:"tenantId" gorm:"index" db:"tenant_id"` // 企业/组织ID
	Theme       ThemeConfig  `json:"theme" gorm:"embedded" db:"theme"`
	Pages       []page.Page  `json:"pages" gorm:"-"` // 不存储在同一表
	Navigation  Navigation   `json:"navigation" gorm:"type:json" db:"navigation"`
	Footer      interface{}  `json:"footer" gorm:"type:json" db:"footer"`
	Thumbnail   string       `json:"thumbnail" db:"thumbnail"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"autoUpdateTime" db:"updated_at"`
	PublishedAt *time.Time   `json:"publishedAt" db:"published_at"`
	Status      string       `json:"status" gorm:"default:draft" db:"status"` // draft, published, archived
}

// ThemeConfig 主题配置
type ThemeConfig struct {
	PrimaryColor    string `json:"primaryColor" db:"primary_color"`
	SecondaryColor  string `json:"secondaryColor" db:"secondary_color"`
	AccentColor     string `json:"accentColor" db:"accent_color"`
	TextColor       string `json:"textColor" db:"text_color"`
	BackgroundColor string `json:"backgroundColor" db:"background_color"`
	FontFamily      string `json:"fontFamily" db:"font_family"`
	HeaderStyle     string `json:"headerStyle" db:"header_style"` // standard, centered, minimal
	BorderRadius    string `json:"borderRadius" db:"border_radius"` // none, small, medium, large
	CustomCSS       string `json:"customCSS" db:"custom_css"`
}

// SiteTemplate 站点模板
type SiteTemplate struct {
	ID          string `json:"id" gorm:"primaryKey" db:"id"`
	Name        string `json:"name" gorm:"not null" db:"name"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail"`
	Description string `json:"description" db:"description"`
	Config      string `json:"config" gorm:"type:json" db:"config"` // 模板配置，JSON格式
}

// Navigation 导航配置
type Navigation struct {
	Type  string           `json:"type" db:"type"` // horizontal, vertical, mega-menu
	Items []NavigationItem `json:"items" db:"items"`
	Style interface{}      `json:"style" db:"style"`
}

// NavigationItem 导航项
type NavigationItem struct {
	ID            string           `json:"id" db:"id"`
	Label         string           `json:"label" db:"label"`
	Link          string           `json:"link" db:"link"`
	Icon          string           `json:"icon,omitempty" db:"icon"`
	Children      []NavigationItem `json:"children,omitempty" db:"children"`
	IsExternalLink bool            `json:"isExternalLink" db:"is_external_link"`
}

// SiteAnalytics 站点分析数据
type SiteAnalytics struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID       string `json:"siteId" gorm:"index" db:"site_id"`
	TenantID     string `json:"tenantId" gorm:"index" db:"tenant_id"`
	PageViews    int64  `json:"pageViews" db:"page_views"`
	UniqueVisitors int64  `json:"uniqueVisitors" db:"unique_visitors"`
	BounceRate   float64 `json:"bounceRate" db:"bounce_rate"`
	AvgSessionDuration int64 `json:"avgSessionDuration" db:"avg_session_duration"`
	TopPages     string `json:"topPages" gorm:"type:json" db:"top_pages"`
	TopSources   string `json:"topSources" gorm:"type:json" db:"top_sources"`
	DeviceStats  string `json:"deviceStats" gorm:"type:json" db:"device_stats"`
	DatePeriod   string `json:"datePeriod" db:"date_period"` // daily, weekly, monthly
	Date         time.Time `json:"date" db:"date"`
}

// SiteSettings 站点设置
type SiteSettings struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string `json:"siteId" gorm:"uniqueIndex:idx_site_setting_key" db:"site_id"`
	SettingKey  string `json:"settingKey" gorm:"uniqueIndex:idx_site_setting_key" db:"setting_key"`
	SettingValue string `json:"settingValue" db:"setting_value"`
	SettingType string `json:"settingType" db:"setting_type"` // string, number, boolean, json
	Group       string `json:"group" db:"group"`
	IsSystem    bool   `json:"isSystem" db:"is_system"`
	Description string `json:"description" db:"description"`
}
