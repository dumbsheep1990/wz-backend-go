package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// SqlUserRepository SQL用户仓库实现
type SqlUserRepository struct {
	conn sqlx.SqlConn
}

// NewSqlUserRepository 创建SQL用户仓库实现
func NewSqlUserRepository(conn sqlx.SqlConn) *SqlUserRepository {
	return &SqlUserRepository{
		conn: conn,
	}
}

// GetUserById 根据ID获取用户
func (r *SqlUserRepository) GetUserById(ctx context.Context, id int64) (*User, error) {
	var user User
	query := `SELECT 
		id, username, password, email, phone, 
		role, status, is_verified, is_company_verified, 
		default_tenant_id, created_at, updated_at 
	FROM users 
	WHERE id = ?`

	err := r.conn.QueryRowCtx(ctx, &user, query, id)
	if err != nil {
		logx.WithContext(ctx).Errorf("根据ID查询用户失败: %v, id: %d", err, id)
		return nil, err
	}

	return &user, nil
}

// GetUserList 获取用户列表
func (r *SqlUserRepository) GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*User, int64, error) {
	// 构建查询条件
	whereClause := "1=1"
	args := []interface{}{}

	// 动态添加查询条件
	if filters != nil {
		for key, value := range filters {
			switch key {
			case "username":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND username LIKE ?"
					args = append(args, fmt.Sprintf("%%%s%%", strVal))
				}
			case "email":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND email LIKE ?"
					args = append(args, fmt.Sprintf("%%%s%%", strVal))
				}
			case "phone":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND phone LIKE ?"
					args = append(args, fmt.Sprintf("%%%s%%", strVal))
				}
			case "status":
				whereClause += " AND status = ?"
				args = append(args, value)
			case "role":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND role = ?"
					args = append(args, strVal)
				}
			case "startTime":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND created_at >= ?"
					args = append(args, strVal)
				}
			case "endTime":
				if strVal, ok := value.(string); ok && strVal != "" {
					whereClause += " AND created_at <= ?"
					args = append(args, strVal)
				}
			}
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereClause)
	var count int64
	err := r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询用户总数失败: %v", err)
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT 
			id, username, password, email, phone, 
			role, status, is_verified, is_company_verified, 
			default_tenant_id, created_at, updated_at 
		FROM users 
		WHERE %s 
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)
	var users []*User
	err = r.conn.QueryRowsCtx(ctx, &users, query, args...)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询用户列表失败: %v", err)
		return nil, 0, err
	}

	return users, count, nil
}

// CreateUser 创建用户
func (r *SqlUserRepository) CreateUser(ctx context.Context, user *User) (int64, error) {
	query := `
		INSERT INTO users (
			username, password, email, phone, role, status, 
			is_verified, is_company_verified, default_tenant_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.conn.ExecCtx(ctx, query,
		user.Username, user.Password, user.Email, user.Phone, user.Role,
		user.Status, user.IsVerified, user.IsCompanyVerified, user.DefaultTenantId,
	)
	if err != nil {
		logx.WithContext(ctx).Errorf("创建用户失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.WithContext(ctx).Errorf("获取新创建用户ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// UpdateUser 更新用户
func (r *SqlUserRepository) UpdateUser(ctx context.Context, user *User) error {
	// 构建动态更新SQL
	var setValues []string
	var args []interface{}

	// 只更新非空字段
	if user.Username != "" {
		setValues = append(setValues, "username = ?")
		args = append(args, user.Username)
	}

	if user.Password != "" {
		setValues = append(setValues, "password = ?")
		args = append(args, user.Password)
	}

	if user.Email != "" {
		setValues = append(setValues, "email = ?")
		args = append(args, user.Email)
	}

	if user.Phone != "" {
		setValues = append(setValues, "phone = ?")
		args = append(args, user.Phone)
	}

	if user.Role != "" {
		setValues = append(setValues, "role = ?")
		args = append(args, user.Role)
	}

	// 状态字段特殊处理，可以为0
	setValues = append(setValues, "status = ?")
	args = append(args, user.Status)

	setValues = append(setValues, "is_verified = ?")
	args = append(args, user.IsVerified)

	setValues = append(setValues, "is_company_verified = ?")
	args = append(args, user.IsCompanyVerified)

	if user.DefaultTenantId > 0 {
		setValues = append(setValues, "default_tenant_id = ?")
		args = append(args, user.DefaultTenantId)
	}

	// 添加更新时间
	setValues = append(setValues, "updated_at = NOW()")

	// 构建最终SQL
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setValues, ", "))
	args = append(args, user.Id)

	_, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		logx.WithContext(ctx).Errorf("更新用户失败: %v, id: %d", err, user.Id)
		return err
	}

	return nil
}

// DeleteUser 删除用户
func (r *SqlUserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := r.conn.ExecCtx(ctx, query, id)
	if err != nil {
		logx.WithContext(ctx).Errorf("删除用户失败: %v, id: %d", err, id)
		return err
	}

	return nil
}
