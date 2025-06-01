package dto

import (
	"time"
)

// RenderRequestDTO 渲染请求DTO
type RenderRequestDTO struct {
	TemplateID  string                 `json:"templateId"`
	SiteID      string                 `json:"siteId"`
	PageID      string                 `json:"pageId,omitempty"`
	Format      string                 `json:"format"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
	CacheConfig CacheConfigDTO         `json:"cacheConfig"`
}

// CacheConfigDTO 缓存配置DTO
type CacheConfigDTO struct {
	Enabled bool   `json:"enabled"`
	TTL     int64  `json:"ttl"` // 单位：秒
	Level   string `json:"level"`
}

// RenderResponseDTO 渲染响应DTO
type RenderResponseDTO struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	ContentType string    `json:"contentType"`
	Format      string    `json:"format"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	CacheHit    bool      `json:"cacheHit"`
}

// TemplateDTO 模板DTO
type TemplateDTO struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Content     string            `json:"content"`
	Type        string            `json:"type"`
	Version     string            `json:"version"`
	SiteID      string            `json:"siteId"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// CreateTemplateDTO 创建模板DTO
type CreateTemplateDTO struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Content     string            `json:"content" binding:"required"`
	Type        string            `json:"type" binding:"required"`
	Version     string            `json:"version"`
	SiteID      string            `json:"siteId" binding:"required"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateTemplateDTO 更新模板DTO
type UpdateTemplateDTO struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Content     string            `json:"content"`
	Type        string            `json:"type"`
	Version     string            `json:"version"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TemplateListDTO 模板列表DTO
type TemplateListDTO struct {
	Templates []TemplateDTO `json:"templates"`
	Total     int           `json:"total"`
}

// PreviewRequestDTO 预览请求DTO
type PreviewRequestDTO struct {
	SiteID   string                 `json:"siteId" binding:"required"`
	PageID   string                 `json:"pageId,omitempty"`
	DeviceType string               `json:"deviceType"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// PageRenderRequestDTO 页面渲染请求DTO
type PageRenderRequestDTO struct {
	SiteID  string `json:"siteId"`
	Slug    string `json:"slug,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Format  string `json:"format,omitempty"`
	Version string `json:"version,omitempty"`
}

// RenderErrorDTO 渲染错误DTO
type RenderErrorDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
