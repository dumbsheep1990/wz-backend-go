package repository

import (
	"context"
)

// 用户相关接口

// UserRepository 用户数据仓库接口
type UserRepository interface {
	// 用户管理
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*User, int64, error)
	CreateUser(ctx context.Context, user *User) (int64, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}

// TenantRepository 租户数据仓库接口
type TenantRepository interface {
	// 租户管理
	GetTenantById(ctx context.Context, id int64) (*Tenant, error)
	GetTenantList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Tenant, int64, error)
	CreateTenant(ctx context.Context, tenant *Tenant) (int64, error)
	UpdateTenant(ctx context.Context, tenant *Tenant) error
	DeleteTenant(ctx context.Context, id int64) error
}

// ContentRepository 内容数据仓库接口
type ContentRepository interface {
	// 内容管理相关方法
}

// TradeRepository 交易数据仓库接口
type TradeRepository interface {
	// 交易管理相关方法
}

// SettingsRepository 系统设置数据仓库接口
type SettingsRepository interface {
	// 系统设置相关方法
	GetSetting(ctx context.Context, key string) (*Setting, error)
	GetAllSettings(ctx context.Context) ([]*Setting, error)
	UpdateSetting(ctx context.Context, key, value string) error
}

// 基础模型定义

// User 用户模型
type User struct {
	Id                int64  `db:"id"`
	Username          string `db:"username"`
	Password          string `db:"password"`
	Email             string `db:"email"`
	Phone             string `db:"phone"`
	Role              string `db:"role"`
	Status            int32  `db:"status"`
	IsVerified        bool   `db:"is_verified"`
	IsCompanyVerified bool   `db:"is_company_verified"`
	DefaultTenantId   int64  `db:"default_tenant_id"`
	CreatedAt         string `db:"created_at"`
	UpdatedAt         string `db:"updated_at"`
}

// Tenant 租户模型
type Tenant struct {
	Id           int64  `db:"id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Subdomain    string `db:"subdomain"`
	Type         int32  `db:"type"`
	Status       int32  `db:"status"`
	Logo         string `db:"logo"`
	ContactEmail string `db:"contact_email"`
	ContactPhone string `db:"contact_phone"`
	AdminUserId  int64  `db:"admin_user_id"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
	ExpireAt     string `db:"expire_at"`
}

// Setting 系统设置模型
type Setting struct {
	Id        int64  `db:"id"`
	Key       string `db:"key"`
	Value     string `db:"value"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
