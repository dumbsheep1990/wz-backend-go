package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Tenant 租户实体
type Tenant struct {
	ID           int64
	Name         string
	Description  string
	Subdomain    string
	Type         int32
	Status       int32
	Logo         string
	ContactEmail string
	ContactPhone string
	AdminUserID  int64
	ExpireAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TenantRepository 租户数据仓库接口
type TenantRepository interface {
	GetTenantById(ctx context.Context, id int64) (*Tenant, error)
	GetTenantList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Tenant, int64, error)
	CreateTenant(ctx context.Context, tenant *Tenant) (int64, error)
	UpdateTenant(ctx context.Context, id int64, updates map[string]interface{}) error
	DeleteTenant(ctx context.Context, id int64) error
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*Tenant, error)
}

// SqlTenantRepository SQL租户仓库实现
type SqlTenantRepository struct {
	conn sqlx.SqlConn
}

// NewSqlTenantRepository 创建租户仓库实例
func NewSqlTenantRepository(conn sqlx.SqlConn) TenantRepository {
	return &SqlTenantRepository{
		conn: conn,
	}
}

// GetTenantById 根据ID获取租户
func (r *SqlTenantRepository) GetTenantById(ctx context.Context, id int64) (*Tenant, error) {
	query := `SELECT id, name, description, subdomain, type, status, logo, contact_email, 
	          contact_phone, admin_user_id, expire_at, created_at, updated_at 
	          FROM tenants WHERE id = ?`
	var tenant Tenant
	err := r.conn.QueryRowCtx(ctx, &tenant, query, id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// GetTenantBySubdomain 根据子域名获取租户
func (r *SqlTenantRepository) GetTenantBySubdomain(ctx context.Context, subdomain string) (*Tenant, error) {
	query := `SELECT id, name, description, subdomain, type, status, logo, contact_email, 
	          contact_phone, admin_user_id, expire_at, created_at, updated_at 
	          FROM tenants WHERE subdomain = ?`
	var tenant Tenant
	err := r.conn.QueryRowCtx(ctx, &tenant, query, subdomain)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// GetTenantList 获取租户列表
func (r *SqlTenantRepository) GetTenantList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Tenant, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if name, ok := filters["name"].(string); ok && name != "" {
			whereClause += " AND name LIKE ?"
			args = append(args, "%"+name+"%")
		}
		if subdomain, ok := filters["subdomain"].(string); ok && subdomain != "" {
			whereClause += " AND subdomain = ?"
			args = append(args, subdomain)
		}
		if tenantType, ok := filters["type"].(int); ok && tenantType > 0 {
			whereClause += " AND type = ?"
			args = append(args, tenantType)
		}
		if status, ok := filters["status"].(int); ok && status > 0 {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if startTime, ok := filters["startTime"].(string); ok && startTime != "" {
			whereClause += " AND created_at >= ?"
			args = append(args, startTime)
		}
		if endTime, ok := filters["endTime"].(string); ok && endTime != "" {
			whereClause += " AND created_at <= ?"
			args = append(args, endTime)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT id, name, description, subdomain, type, status, logo, 
	                      contact_email, contact_phone, admin_user_id, expire_at, 
	                      created_at, updated_at FROM tenants %s ORDER BY id DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM tenants %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var tenants []*Tenant
	err := r.conn.QueryRowsCtx(ctx, &tenants, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return tenants, count, nil
}

// CreateTenant 创建租户
func (r *SqlTenantRepository) CreateTenant(ctx context.Context, tenant *Tenant) (int64, error) {
	// 使用事务保证一致性
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	// 如果未设置过期时间，默认一年后过期
	if tenant.ExpireAt.IsZero() {
		tenant.ExpireAt = now.AddDate(1, 0, 0)
	}

	query := `INSERT INTO tenants (name, description, subdomain, type, status, logo, 
	         contact_email, contact_phone, admin_user_id, expire_at, created_at, updated_at) 
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := tx.ExecCtx(ctx, query,
		tenant.Name,
		tenant.Description,
		tenant.Subdomain,
		tenant.Type,
		tenant.Status,
		tenant.Logo,
		tenant.ContactEmail,
		tenant.ContactPhone,
		tenant.AdminUserID,
		tenant.ExpireAt,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateTenant 更新租户信息
func (r *SqlTenantRepository) UpdateTenant(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	sets := []string{}
	args := []interface{}{}

	// 添加更新字段
	for k, v := range updates {
		// 转换键名为下划线形式
		fieldName := camelToSnake(k)
		sets = append(sets, fmt.Sprintf("%s = ?", fieldName))
		args = append(args, v)
	}

	// 添加更新时间
	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now())

	// 添加ID条件
	args = append(args, id)

	query := fmt.Sprintf("UPDATE tenants SET %s WHERE id = ?", strings.Join(sets, ", "))
	_, err := r.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteTenant 删除租户
func (r *SqlTenantRepository) DeleteTenant(ctx context.Context, id int64) error {
	// 使用事务确保一致性
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 先查询租户是否存在
	tenant, err := r.GetTenantById(ctx, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if tenant == nil {
		tx.Rollback()
		return sql.ErrNoRows
	}

	// 执行删除操作
	query := "DELETE FROM tenants WHERE id = ?"
	_, err = tx.ExecCtx(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// 辅助函数：驼峰命名转换为下划线命名
func camelToSnake(s string) string {
	var result strings.Builder
	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(v)
	}
	return strings.ToLower(result.String())
} 