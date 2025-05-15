package repository

import (
	"context"
	"wz-backend-go/services/admin-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.User, int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, id int64, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}

// TenantRepository 租户仓库接口
type TenantRepository interface {
	GetTenantList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Tenant, int64, error)
	GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error)
	CreateTenant(ctx context.Context, tenant *model.Tenant) (int64, error)
	UpdateTenant(ctx context.Context, id int64, tenant *model.Tenant) error
	DeleteTenant(ctx context.Context, id int64) error
}

// ContentRepository 内容仓库接口
type ContentRepository interface {
	GetContentList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Content, int64, error)
	GetContentByID(ctx context.Context, id int64) (*model.Content, error)
	UpdateContentStatus(ctx context.Context, id int64, status int, reason string) error
	DeleteContent(ctx context.Context, id int64) error
	RecommendContent(ctx context.Context, contentId int64, priority int) error
	CancelRecommendation(ctx context.Context, id int64) error
}

// TradeRepository 交易仓库接口
type TradeRepository interface {
	GetOrderList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Order, int64, error)
	GetOrderByID(ctx context.Context, id int64) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status int, remark string) error
	DeleteOrder(ctx context.Context, id int64) error
}

// SettingsRepository 系统设置仓库接口
type SettingsRepository interface {
	GetSettings(ctx context.Context) (map[string]string, error)
	UpdateSetting(ctx context.Context, key, value string) error
}

// AdminRepository 管理员仓库接口
type AdminRepository interface {
	GetAdminList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Admin, int64, error)
	GetAdminByID(ctx context.Context, id int64) (*model.Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*model.Admin, error)
	CreateAdmin(ctx context.Context, admin *model.Admin) (int64, error)
	UpdateAdmin(ctx context.Context, id int64, admin *model.Admin) error
	DeleteAdmin(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

// RoleRepository 角色仓库接口
type RoleRepository interface {
	GetRoleList(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error)
	GetRoleByID(ctx context.Context, id int64) (*model.Role, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) (int64, error)
	UpdateRole(ctx context.Context, id int64, role *model.Role) error
	DeleteRole(ctx context.Context, id int64) error
}

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	GetOperationLogList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error)
	GetOperationLogByID(ctx context.Context, id int64) (*model.OperationLog, error)
	CreateOperationLog(ctx context.Context, log *model.OperationLog) error
}

// 以下是仓库实现的工厂函数

// NewUserRepository 创建用户仓库
func NewUserRepository(conn sqlx.SqlConn) UserRepository {
	// TODO: 实现用户仓库
	return nil
}

// NewTenantRepository 创建租户仓库
func NewTenantRepository(conn sqlx.SqlConn) TenantRepository {
	// TODO: 实现租户仓库
	return nil
}

// NewContentRepository 创建内容仓库
func NewContentRepository(conn sqlx.SqlConn) ContentRepository {
	// TODO: 实现内容仓库
	return nil
}

// NewTradeRepository 创建交易仓库
func NewTradeRepository(conn sqlx.SqlConn) TradeRepository {
	// TODO: 实现交易仓库
	return nil
}

// NewSettingsRepository 创建系统设置仓库
func NewSettingsRepository(conn sqlx.SqlConn) SettingsRepository {
	// TODO: 实现系统设置仓库
	return nil
}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository(conn sqlx.SqlConn) AdminRepository {
	// TODO: 实现管理员仓库
	return nil
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository(conn sqlx.SqlConn) RoleRepository {
	// TODO: 实现角色仓库
	return nil
}

// NewOperationLogRepository 创建操作日志仓库
func NewOperationLogRepository(conn sqlx.SqlConn) OperationLogRepository {
	// TODO: 实现操作日志仓库
	return nil
}
