package dto

import (
	"time"
)

// CreateSiteRequest 创建站点请求
type CreateSiteRequest struct {
	Name        string      `json:"name" validate:"required,min=2,max=100"`
	Description string      `json:"description" validate:"max=500"`
	Domain      string      `json:"domain" validate:"required"`
	Logo        string      `json:"logo"`
	Favicon     string      `json:"favicon"`
	Theme       ThemeDTO    `json:"theme"`
	Thumbnail   string      `json:"thumbnail"`
}

// UpdateSiteRequest 更新站点请求
type UpdateSiteRequest struct {
	Name        *string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=500"`
	Domain      *string     `json:"domain,omitempty"`
	Logo        *string     `json:"logo,omitempty"`
	Favicon     *string     `json:"favicon,omitempty"`
	Theme       *ThemeDTO   `json:"theme,omitempty"`
	Thumbnail   *string     `json:"thumbnail,omitempty"`
}

// SiteDTO 站点数据传输对象
type SiteDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Domain      string     `json:"domain"`
	Logo        string     `json:"logo"`
	Favicon     string     `json:"favicon"`
	TenantID    string     `json:"tenantId"`
	Theme       ThemeDTO   `json:"theme"`
	Status      string     `json:"status"`
	Thumbnail   string     `json:"thumbnail"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	PublishedAt *time.Time `json:"publishedAt"`
}

// ThemeDTO 主题配置数据传输对象
type ThemeDTO struct {
	PrimaryColor    string `json:"primaryColor" validate:"required"`
	SecondaryColor  string `json:"secondaryColor" validate:"required"`
	AccentColor     string `json:"accentColor" validate:"required"`
	TextColor       string `json:"textColor" validate:"required"`
	BackgroundColor string `json:"backgroundColor" validate:"required"`
	FontFamily      string `json:"fontFamily" validate:"required"`
	HeaderStyle     string `json:"headerStyle" validate:"required,oneof=standard centered minimal"`
	BorderRadius    string `json:"borderRadius" validate:"required,oneof=none small medium large"`
	CustomCSS       string `json:"customCSS" validate:"max=10000"`
}

// SiteListRequest 站点列表查询请求
type SiteListRequest struct {
	Status string `json:"status" form:"status"`
	Search string `json:"search" form:"search"`
	Page   int    `json:"page" form:"page" validate:"min=1"`
	Size   int    `json:"size" form:"size" validate:"min=1,max=100"`
}

// SiteListResponse 站点列表响应
type SiteListResponse struct {
	Sites []SiteDTO `json:"sites"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

// PublishSiteRequest 发布站点请求
type PublishSiteRequest struct {
	SiteID string `json:"siteId" validate:"required"`
}

// PublishSiteResponse 发布站点响应
type PublishSiteResponse struct {
	Site        SiteDTO   `json:"site"`
	PublishedAt time.Time `json:"publishedAt"`
} 