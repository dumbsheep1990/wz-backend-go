package models

import "time"

// Site 站点模型
type Site struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Domain      string       `json:"domain"`
	Logo        string       `json:"logo"`
	Favicon     string       `json:"favicon"`
	TenantID    string       `json:"tenantId"` // 企业/组织ID
	Theme       ThemeConfig  `json:"theme" gorm:"embedded"`
	Pages       []Page       `json:"pages" gorm:"-"` // 不存储在同一表
	Navigation  Navigation   `json:"navigation" gorm:"type:json"`
	Footer      interface{}  `json:"footer" gorm:"type:json"`
	Thumbnail   string       `json:"thumbnail"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	PublishedAt *time.Time   `json:"publishedAt"`
	Status      string       `json:"status"` // draft, published, archived
}

// ThemeConfig 主题配置
type ThemeConfig struct {
	PrimaryColor    string `json:"primaryColor"`
	SecondaryColor  string `json:"secondaryColor"`
	AccentColor     string `json:"accentColor"`
	TextColor       string `json:"textColor"`
	BackgroundColor string `json:"backgroundColor"`
	FontFamily      string `json:"fontFamily"`
	HeaderStyle     string `json:"headerStyle"` // standard, centered, minimal
	BorderRadius    string `json:"borderRadius"` // none, small, medium, large
	CustomCSS       string `json:"customCSS"`
}

// SiteTemplate 站点模板
type SiteTemplate struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Config      string `json:"config" gorm:"type:json"` // 模板配置，JSON格式
} 