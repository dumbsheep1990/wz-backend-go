package repository

import (
	"context"
	
	"wz-backend-go/internal/domain/model"
)

// TenantRepository 租户仓库接口
type TenantRepository interface {
	// 创建新租户
	CreateTenant(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error)
	// 获取租户详情
	GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error)
	// 通过子域名获取租户
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error)
	// 更新租户信息
	UpdateTenant(ctx context.Context, tenant *model.Tenant) error
	// 列出所有活跃租户
	ListActiveTenants(ctx context.Context) ([]*model.Tenant, error)
	// 添加用户到租户
	AddUserToTenant(ctx context.Context, tenantUser *model.TenantUser) error
	// 获取租户下的用户列表
	GetTenantUsers(ctx context.Context, tenantID int64) ([]*model.User, error)
	// 从租户中移除用户
	RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error
	// 检查用户是否属于租户，返回是否存在、角色和错误
	CheckUserInTenant(ctx context.Context, tenantID, userID int64) (bool, string, error)
}
