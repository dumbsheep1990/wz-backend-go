package model

import (
	"time"
)

// TenantType 租户类型
type TenantType int32

const (
	TenantTypeEnterprise  TenantType = 1 // 企业类型，主要功能为商品交易
	TenantTypePersonal    TenantType = 2 // 个人类型，主要功能为博客作品集
	TenantTypeEducational TenantType = 3 // 教育机构类型，主要功能为课程，学术资料库等
)

// Tenant 租户信息
type Tenant struct {
	ID              int64      `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`               // 租户名称
	Subdomain       string     `db:"subdomain" json:"subdomain"`     // 子域名
	TenantType      TenantType `db:"tenant_type" json:"tenant_type"` // 租户类型
	Description     string     `db:"description" json:"description"` // 租户描述
	Logo            string     `db:"logo" json:"logo"`               // 租户Logo
	CreatorUserID   int64      `db:"creator_user_id" json:"creator_user_id"` // 创建者用户ID
	Status          int32      `db:"status" json:"status"`           // 状态 1-正常 2-禁用
	ExpirationDate  time.Time  `db:"expiration_date" json:"expiration_date"` // 过期时间
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// TenantUser 租户-用户关联
type TenantUser struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   int64     `db:"tenant_id" json:"tenant_id"`     // 租户ID
	UserID     int64     `db:"user_id" json:"user_id"`         // 用户ID
	Role       string    `db:"role" json:"role"`               // 角色: admin(租户管理员), user(普通用户)
	Status     int32     `db:"status" json:"status"`           // 状态 1-正常 2-禁用
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name        string     `json:"name" validate:"required"`
	Subdomain   string     `json:"subdomain" validate:"required,alphanum"` // 子域名只允许字母和数字
	TenantType  TenantType `json:"tenant_type" validate:"required"`
	Description string     `json:"description"`
	Logo        string     `json:"logo"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	ID          int64      `json:"-"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Logo        string     `json:"logo,omitempty"`
	Status      int32      `json:"status,omitempty"`
}

// TenantResponse 租户信息响应
type TenantResponse struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Subdomain   string     `json:"subdomain"`
	TenantType  TenantType `json:"tenant_type"`
	Description string     `json:"description"`
	Logo        string     `json:"logo"`
	Status      int32      `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}
