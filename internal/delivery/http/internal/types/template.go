package types

import "time"

// TemplateType 模板类型
type TemplateType string

const (
	TemplateBanner  TemplateType = "banner"  // Banner模板
	TemplateProduct TemplateType = "product" // 产品模板
	TemplateArticle TemplateType = "article" // 文章模板
)

// CreateTemplateReq 创建模板请求
type CreateTemplateReq struct {
	UserID      int64       `json:"-"`                              // 所属用户ID
	Name        string      `json:"name" validate:"required"`       // 模板名称
	Type        TemplateType `json:"type" validate:"required"`      // 模板类型
	Preview     string      `json:"preview"`                        // 预览图路径
	Content     string      `json:"content"`                        // 模板内容（JSON格式）
	PublicShare bool        `json:"public_share"`                   // 是否公开分享
}

// CreateTemplateResp 创建模板响应
type CreateTemplateResp struct {
	ID        int64       `json:"id"`                             // 模板ID
	Name      string      `json:"name"`                           // 模板名称
	Type      TemplateType `json:"type"`                          // 模板类型
	Preview   string      `json:"preview"`                        // 预览图路径
	Enabled   bool        `json:"enabled"`                        // 是否启用
	IsNew     bool        `json:"is_new"`                         // 是否为新模板
	CreatedAt time.Time   `json:"created_at"`                     // 创建时间
}

// UpdateTemplateReq 更新模板请求
type UpdateTemplateReq struct {
	TemplateID  int64       `json:"-"`                             // 模板ID
	UserID      int64       `json:"-"`                             // 所属用户ID
	Name        string      `json:"name" validate:"required"`      // 模板名称
	Type        TemplateType `json:"type" validate:"required"`     // 模板类型
	Preview     string      `json:"preview"`                       // 预览图路径
	Content     string      `json:"content"`                       // 模板内容（JSON格式）
	PublicShare bool        `json:"public_share"`                  // 是否公开分享
}

// UpdateTemplateResp 更新模板响应
type UpdateTemplateResp struct {
	Success bool `json:"success"` // 是否成功
}

// UpdateTemplateStatusReq 更新模板状态请求
type UpdateTemplateStatusReq struct {
	TemplateID int64 `json:"-"`          // 模板ID
	UserID     int64 `json:"-"`          // 所属用户ID
	Enabled    bool  `json:"enabled"`    // 是否启用
}

// UpdateTemplateStatusResp 更新模板状态响应
type UpdateTemplateStatusResp struct {
	Success bool `json:"success"` // 是否成功
}

// DeleteTemplateResp 删除模板响应
type DeleteTemplateResp struct {
	Success bool `json:"success"` // 是否成功
}

// TemplateItem 模板列表项
type TemplateItem struct {
	ID          int64       `json:"id"`          // 模板ID
	Name        string      `json:"name"`        // 模板名称
	Type        TemplateType `json:"type"`       // 模板类型
	Preview     string      `json:"preview"`     // 预览图路径
	Enabled     bool        `json:"enabled"`     // 是否启用
	IsNew       bool        `json:"is_new"`      // 是否为新模板
	PublicShare bool        `json:"public_share"`// 是否公开分享
	CreatedAt   time.Time   `json:"created_at"`  // 创建时间
}

// GetTemplatesResp 获取模板列表响应
type GetTemplatesResp struct {
	Total   int           `json:"total"`    // 总数
	Page    int           `json:"page"`     // 当前页码
	Size    int           `json:"size"`     // 每页大小
	Items   []TemplateItem `json:"items"`   // 模板列表
}

// GetTemplateResp 获取单个模板响应
type GetTemplateResp struct {
	ID          int64       `json:"id"`          // 模板ID
	Name        string      `json:"name"`        // 模板名称
	Type        TemplateType `json:"type"`       // 模板类型
	Preview     string      `json:"preview"`     // 预览图路径
	Content     string      `json:"content"`     // 模板内容
	Enabled     bool        `json:"enabled"`     // 是否启用
	IsNew       bool        `json:"is_new"`      // 是否为新模板
	PublicShare bool        `json:"public_share"`// 是否公开分享
	CreatedAt   time.Time   `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time   `json:"updated_at"`  // 更新时间
}
