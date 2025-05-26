package model

import (
	"time"
)

// TemplateType 模板类型
type TemplateType string

const (
	TemplateBanner  TemplateType = "banner"  // Banner模板
	TemplateProduct TemplateType = "product" // 产品模板
	TemplateArticle TemplateType = "article" // 文章模板
)

// Template 模板信息
type Template struct {
	ID          int64        `db:"id" json:"id"`
	UserID      int64        `db:"user_id" json:"user_id"`           // 所属用户ID
	Name        string       `db:"name" json:"name"`                 // 模板名称
	Type        TemplateType `db:"type" json:"type"`                 // 模板类型
	Preview     string       `db:"preview" json:"preview"`           // 预览图路径
	Content     string       `db:"content" json:"content"`           // 模板内容（JSON格式）
	Enabled     bool         `db:"enabled" json:"enabled"`           // 是否启用
	IsNew       bool         `db:"is_new" json:"is_new"`             // 是否为新模板
	PublicShare bool         `db:"public_share" json:"public_share"` // 是否公开分享
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	UserID      int64        `json:"-"`                              // 所属用户ID
	Name        string       `json:"name" validate:"required"`       // 模板名称
	Type        TemplateType `json:"type" validate:"required"`       // 模板类型
	Preview     string       `json:"preview"`                        // 预览图路径
	Content     string       `json:"content"`                        // 模板内容（JSON格式）
	PublicShare bool         `json:"public_share"`                   // 是否公开分享
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	TemplateID  int64        `json:"-"`                              // 模板ID
	UserID      int64        `json:"-"`                              // 所属用户ID
	Name        string       `json:"name" validate:"required"`       // 模板名称
	Type        TemplateType `json:"type" validate:"required"`       // 模板类型
	Preview     string       `json:"preview"`                        // 预览图路径
	Content     string       `json:"content"`                        // 模板内容（JSON格式）
	PublicShare bool         `json:"public_share"`                   // 是否公开分享
}

// UpdateTemplateStatusRequest 更新模板状态请求
type UpdateTemplateStatusRequest struct {
	TemplateID int64 `json:"-"`                                      // 模板ID
	UserID     int64 `json:"-"`                                      // 所属用户ID
	Enabled    bool  `json:"enabled"`                                // 是否启用
}

// TemplateListItem 模板列表项
type TemplateListItem struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Type        TemplateType `json:"type"`
	Preview     string       `json:"preview"`
	Enabled     bool         `json:"enabled"`
	IsNew       bool         `json:"is_new"`
	PublicShare bool         `json:"public_share"`
	CreatedAt   time.Time    `json:"created_at"`
}
