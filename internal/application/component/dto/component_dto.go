package dto

import (
	"time"
)

// CreateComponentRequest 创建组件请求
type CreateComponentRequest struct {
	Name          string   `json:"name" validate:"required,min=2,max=100"`
	Description   string   `json:"description" validate:"max=500"`
	ComponentType string   `json:"componentType" validate:"required,oneof=header navbar hero feature content footer sidebar card gallery contact custom"`
	Template      string   `json:"template" validate:"required"`
	Config        string   `json:"config"`
	Preview       string   `json:"preview"`
	Category      string   `json:"category" validate:"max=50"`
	Tags          []string `json:"tags"`
	IsPublic      bool     `json:"isPublic"`
}

// UpdateComponentRequest 更新组件请求
type UpdateComponentRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Template    *string  `json:"template,omitempty"`
	Config      *string  `json:"config,omitempty"`
	Preview     *string  `json:"preview,omitempty"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,max=50"`
	Tags        []string `json:"tags,omitempty"`
}

// ComponentDTO 组件数据传输对象
type ComponentDTO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ComponentType string    `json:"componentType"`
	Template      string    `json:"template"`
	Config        string    `json:"config"`
	Preview       string    `json:"preview"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	IsPublic      bool      `json:"isPublic"`
	TenantID      string    `json:"tenantId"`
	Version       string    `json:"version"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ComponentListRequest 组件列表查询请求
type ComponentListRequest struct {
	ComponentType string   `json:"componentType" form:"componentType"`
	Category      string   `json:"category" form:"category"`
	IsPublic      *bool    `json:"isPublic" form:"isPublic"`
	Search        string   `json:"search" form:"search"`
	Tags          []string `json:"tags" form:"tags"`
	Page          int      `json:"page" form:"page" validate:"min=1"`
	Size          int      `json:"size" form:"size" validate:"min=1,max=100"`
	SortBy        string   `json:"sortBy" form:"sortBy" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder     string   `json:"sortOrder" form:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

// ComponentListResponse 组件列表响应
type ComponentListResponse struct {
	Components []ComponentDTO `json:"components"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
}

// PublicComponentListRequest 公开组件列表查询请求
type PublicComponentListRequest struct {
	ComponentType string   `json:"componentType" form:"componentType"`
	Category      string   `json:"category" form:"category"`
	Search        string   `json:"search" form:"search"`
	Tags          []string `json:"tags" form:"tags"`
	Page          int      `json:"page" form:"page" validate:"min=1"`
	Size          int      `json:"size" form:"size" validate:"min=1,max=100"`
	SortBy        string   `json:"sortBy" form:"sortBy" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder     string   `json:"sortOrder" form:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

// PublicComponentListResponse 公开组件列表响应
type PublicComponentListResponse struct {
	Components []ComponentDTO `json:"components"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
}

// AvailableComponentListRequest 可用组件列表查询请求（自有+公开）
type AvailableComponentListRequest struct {
	ComponentType string   `json:"componentType" form:"componentType" validate:"required"`
	Category      string   `json:"category" form:"category"`
	Search        string   `json:"search" form:"search"`
	Tags          []string `json:"tags" form:"tags"`
	Page          int      `json:"page" form:"page" validate:"min=1"`
	Size          int      `json:"size" form:"size" validate:"min=1,max=100"`
}

// AvailableComponentListResponse 可用组件列表响应
type AvailableComponentListResponse struct {
	OwnComponents    []ComponentDTO `json:"ownComponents"`
	PublicComponents []ComponentDTO `json:"publicComponents"`
	Total            int64          `json:"total"`
	Page             int            `json:"page"`
	Size             int            `json:"size"`
}

// MakePublicRequest 设为公开请求
type MakePublicRequest struct {
	ComponentID string `json:"componentId" validate:"required"`
}

// MakePublicResponse 设为公开响应
type MakePublicResponse struct {
	Component ComponentDTO `json:"component"`
	Message   string       `json:"message"`
}

// ComponentPreviewRequest 组件预览请求
type ComponentPreviewRequest struct {
	Template string `json:"template" validate:"required"`
	Config   string `json:"config"`
}

// ComponentPreviewResponse 组件预览响应
type ComponentPreviewResponse struct {
	HTML    string `json:"html"`
	CSS     string `json:"css"`
	Preview string `json:"preview"`
} 