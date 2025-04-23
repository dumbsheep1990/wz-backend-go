package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/repository"
)

type tenantRepository struct {
	conn sqlx.SqlConn
}

// NewTenantRepository 创建租户仓库实例
func NewTenantRepository(conn sqlx.SqlConn) repository.TenantRepository {
	return &tenantRepository{
		conn: conn,
	}
}

// CreateTenant 创建新租户
func (r *tenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error) {
	query := `
		INSERT INTO tenants (
			name, subdomain, tenant_type, description, logo, 
			creator_user_id, status, expiration_date, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	expirationDate := now.AddDate(1, 0, 0) // 默认一年有效期
	
	result, err := r.conn.ExecCtx(ctx, query,
		tenant.Name, tenant.Subdomain, tenant.TenantType, tenant.Description, tenant.Logo,
		tenant.CreatorUserID, tenant.Status, expirationDate, now, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	tenant.ID = id
	tenant.CreatedAt = now
	tenant.UpdatedAt = now
	tenant.ExpirationDate = expirationDate

	return tenant, nil
}

// GetTenantByID 通过ID获取租户
func (r *tenantRepository) GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error) {
	query := `
		SELECT id, name, subdomain, tenant_type, description, logo, 
			   creator_user_id, status, expiration_date, created_at, updated_at
		FROM tenants
		WHERE id = ?
	`

	var tenant model.Tenant
	err := r.conn.QueryRowCtx(ctx, &tenant, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("租户不存在: %d", id)
		}
		return nil, err
	}

	return &tenant, nil
}

// GetTenantBySubdomain 通过子域名获取租户
func (r *tenantRepository) GetTenantBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error) {
	query := `
		SELECT id, name, subdomain, tenant_type, description, logo, 
			   creator_user_id, status, expiration_date, created_at, updated_at
		FROM tenants
		WHERE subdomain = ? AND status = 1
	`

	var tenant model.Tenant
	err := r.conn.QueryRowCtx(ctx, &tenant, query, subdomain)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("租户不存在: %s", subdomain)
		}
		return nil, err
	}

	return &tenant, nil
}

// UpdateTenant 更新租户信息
func (r *tenantRepository) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	query := `
		UPDATE tenants
		SET name = ?, description = ?, logo = ?, 
		    status = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.conn.ExecCtx(ctx, query,
		tenant.Name, tenant.Description, tenant.Logo,
		tenant.Status, now, tenant.ID,
	)
	if err != nil {
		return err
	}

	tenant.UpdatedAt = now
	return nil
}

// ListActiveTenants 列出所有活跃租户
func (r *tenantRepository) ListActiveTenants(ctx context.Context) ([]*model.Tenant, error) {
	query := `
		SELECT id, name, subdomain, tenant_type, description, logo, 
			   creator_user_id, status, expiration_date, created_at, updated_at
		FROM tenants
		WHERE status = 1 AND expiration_date > NOW()
		ORDER BY id DESC
	`

	var tenants []*model.Tenant
	err := r.conn.QueryRowsCtx(ctx, &tenants, query)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

// AddUserToTenant 添加用户到租户
func (r *tenantRepository) AddUserToTenant(ctx context.Context, tenantUser *model.TenantUser) error {
	query := `
		INSERT INTO tenant_users (
			tenant_id, user_id, role, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE role = ?, status = ?, updated_at = ?
	`

	now := time.Now()
	_, err := r.conn.ExecCtx(ctx, query,
		tenantUser.TenantID, tenantUser.UserID, tenantUser.Role, tenantUser.Status, now, now,
		tenantUser.Role, tenantUser.Status, now,
	)
	
	return err
}

// GetTenantUsers 获取租户下的用户列表
func (r *tenantRepository) GetTenantUsers(ctx context.Context, tenantID int64) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.phone, u.status, 
		       u.is_verified, u.is_company_verified, u.default_tenant_id, u.role,
			   u.created_at, u.updated_at
		FROM users u
		JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.tenant_id = ? AND tu.status = 1
		ORDER BY tu.id DESC
	`

	var users []*model.User
	err := r.conn.QueryRowsCtx(ctx, &users, query, tenantID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// RemoveUserFromTenant 从租户中移除用户
func (r *tenantRepository) RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error {
	query := `DELETE FROM tenant_users WHERE tenant_id = ? AND user_id = ?`
	_, err := r.conn.ExecCtx(ctx, query, tenantID, userID)
	return err
}

// CheckUserInTenant 检查用户是否属于租户，返回是否存在、角色和错误
func (r *tenantRepository) CheckUserInTenant(ctx context.Context, tenantID, userID int64) (bool, string, error) {
	query := `
		SELECT role
		FROM tenant_users
		WHERE tenant_id = ? AND user_id = ? AND status = 1
	`

	var role string
	err := r.conn.QueryRowCtx(ctx, &role, query, tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	return true, role, nil
}
