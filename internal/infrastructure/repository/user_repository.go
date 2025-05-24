package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"wz-backend-go/internal/domain/user"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	if u.ID().Value() == 0 {
		// 新用户插入
		query := `
			INSERT INTO users (
				username, password, email, phone, status, is_verified, is_company_verified, 
				default_tenant_id, role, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		
		result, err := r.db.ExecContext(ctx, query,
			u.Username().Value(),
			string(u.Password()),
			u.Email().Value(),
			u.Phone().Value(),
			int32(u.Status()),
			u.IsVerified(),
			u.IsCompanyVerified(),
			u.DefaultTenantID().Value(),
			string(u.Role()),
			u.CreatedAt(),
			u.UpdatedAt(),
		)
		
		if err != nil {
			return err
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		
		// 更新ID（这需要User实体暴露设置ID的方法，或者通过重构方式创建新实例）
		u.SetID(user.NewUserID(id))
		
		return nil
	} else {
		// 更新现有用户
		query := `
			UPDATE users SET 
				username = ?,
				password = ?,
				email = ?,
				phone = ?,
				status = ?,
				is_verified = ?,
				is_company_verified = ?,
				default_tenant_id = ?,
				role = ?,
				updated_at = ?
			WHERE id = ?
		`
		
		_, err := r.db.ExecContext(ctx, query,
			u.Username().Value(),
			string(u.Password()),
			u.Email().Value(),
			u.Phone().Value(),
			int32(u.Status()),
			u.IsVerified(),
			u.IsCompanyVerified(),
			u.DefaultTenantID().Value(),
			string(u.Role()),
			time.Now(),
			u.ID().Value(),
		)
		
		return err
	}
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users WHERE id = ?
	`
	
	var (
		idVal             int64
		username          string
		password          string
		email             string
		phone             string
		status            int32
		isVerified        bool
		isCompanyVerified bool
		defaultTenantID   int64
		role              string
		createdAt         time.Time
		updatedAt         time.Time
	)
	
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idVal, &username, &password, &email, &phone, &status, &isVerified, &isCompanyVerified,
		&defaultTenantID, &role, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 重构用户实体
	return user.ReconstructUser(
		idVal, username, email, phone, []byte(password), status, isVerified, isCompanyVerified,
		defaultTenantID, role, createdAt, updatedAt,
	)
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username user.Username) (*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users WHERE username = ?
	`
	
	var (
		id                int64
		usernameVal       string
		password          string
		email             string
		phone             string
		status            int32
		isVerified        bool
		isCompanyVerified bool
		defaultTenantID   int64
		role              string
		createdAt         time.Time
		updatedAt         time.Time
	)
	
	err := r.db.QueryRowContext(ctx, query, username.Value()).Scan(
		&id, &usernameVal, &password, &email, &phone, &status, &isVerified, &isCompanyVerified,
		&defaultTenantID, &role, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 重构用户实体
	return user.ReconstructUser(
		id, usernameVal, email, phone, []byte(password), status, isVerified, isCompanyVerified,
		defaultTenantID, role, createdAt, updatedAt,
	)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users WHERE email = ?
	`
	
	var (
		id                int64
		username          string
		password          string
		emailVal          string
		phone             string
		status            int32
		isVerified        bool
		isCompanyVerified bool
		defaultTenantID   int64
		role              string
		createdAt         time.Time
		updatedAt         time.Time
	)
	
	err := r.db.QueryRowContext(ctx, query, email.Value()).Scan(
		&id, &username, &password, &emailVal, &phone, &status, &isVerified, &isCompanyVerified,
		&defaultTenantID, &role, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 重构用户实体
	return user.ReconstructUser(
		id, username, emailVal, phone, []byte(password), status, isVerified, isCompanyVerified,
		defaultTenantID, role, createdAt, updatedAt,
	)
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(ctx context.Context, phone user.Phone) (*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users WHERE phone = ?
	`
	
	var (
		id                int64
		username          string
		password          string
		email             string
		phoneVal          string
		status            int32
		isVerified        bool
		isCompanyVerified bool
		defaultTenantID   int64
		role              string
		createdAt         time.Time
		updatedAt         time.Time
	)
	
	err := r.db.QueryRowContext(ctx, query, phone.Value()).Scan(
		&id, &username, &password, &email, &phoneVal, &status, &isVerified, &isCompanyVerified,
		&defaultTenantID, &role, &createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 重构用户实体
	return user.ReconstructUser(
		id, username, email, phoneVal, []byte(password), status, isVerified, isCompanyVerified,
		defaultTenantID, role, createdAt, updatedAt,
	)
}

// UpdateUser 更新用户
func (r *UserRepository) UpdateUser(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users SET 
			username = ?,
			password = ?,
			email = ?,
			phone = ?,
			status = ?,
			is_verified = ?,
			is_company_verified = ?,
			default_tenant_id = ?,
			role = ?,
			updated_at = ?
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query,
		u.Username().Value(),
		string(u.Password()),
		u.Email().Value(),
		u.Phone().Value(),
		int32(u.Status()),
		u.IsVerified(),
		u.IsCompanyVerified(),
		u.DefaultTenantID().Value(),
		string(u.Role()),
		time.Now(),
		u.ID().Value(),
	)
	
	return err
}

// DeleteUser 删除用户
func (r *UserRepository) DeleteUser(ctx context.Context, id user.UserID) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id.Value())
	return err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username user.Username) (bool, error) {
	query := "SELECT 1 FROM users WHERE username = ? LIMIT 1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, username.Value()).Scan(&exists)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	query := "SELECT 1 FROM users WHERE email = ? LIMIT 1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, email.Value()).Scan(&exists)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// ExistsByPhone 检查手机号是否存在
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone user.Phone) (bool, error) {
	query := "SELECT 1 FROM users WHERE phone = ? LIMIT 1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, phone.Value()).Scan(&exists)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// FindUsersByTenant 根据租户ID查找用户
func (r *UserRepository) FindUsersByTenant(ctx context.Context, tenantID user.TenantID, offset, limit int) ([]*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users 
		WHERE default_tenant_id = ?
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID.Value(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*user.User
	
	for rows.Next() {
		var (
			id                int64
			username          string
			password          string
			email             string
			phone             string
			status            int32
			isVerified        bool
			isCompanyVerified bool
			defaultTenantID   int64
			role              string
			createdAt         time.Time
			updatedAt         time.Time
		)
		
		if err := rows.Scan(&id, &username, &password, &email, &phone, &status, &isVerified, &isCompanyVerified,
			&defaultTenantID, &role, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		
		u, err := user.ReconstructUser(
			id, username, email, phone, []byte(password), status, isVerified, isCompanyVerified,
			defaultTenantID, role, createdAt, updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		users = append(users, u)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// CountUsersByTenant 统计租户用户数量
func (r *UserRepository) CountUsersByTenant(ctx context.Context, tenantID user.TenantID) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE default_tenant_id = ?"
	
	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID.Value()).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// FindUsersByRole 根据角色查找用户
func (r *UserRepository) FindUsersByRole(ctx context.Context, role user.UserRole, offset, limit int) ([]*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users 
		WHERE role = ?
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, string(role), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*user.User
	
	for rows.Next() {
		var (
			id                int64
			username          string
			password          string
			email             string
			phone             string
			status            int32
			isVerified        bool
			isCompanyVerified bool
			defaultTenantID   int64
			roleVal           string
			createdAt         time.Time
			updatedAt         time.Time
		)
		
		if err := rows.Scan(&id, &username, &password, &email, &phone, &status, &isVerified, &isCompanyVerified,
			&defaultTenantID, &roleVal, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		
		u, err := user.ReconstructUser(
			id, username, email, phone, []byte(password), status, isVerified, isCompanyVerified,
			defaultTenantID, roleVal, createdAt, updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		users = append(users, u)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// CountUsersByRole 统计角色用户数量
func (r *UserRepository) CountUsersByRole(ctx context.Context, role user.UserRole) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE role = ?"
	
	var count int
	err := r.db.QueryRowContext(ctx, query, string(role)).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// FindActiveUsers 查找活跃用户
func (r *UserRepository) FindActiveUsers(ctx context.Context, offset, limit int) ([]*user.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, 
		       default_tenant_id, role, created_at, updated_at
		FROM users 
		WHERE status = ?
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, int32(user.UserStatusActive), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*user.User
	
	for rows.Next() {
		var (
			id                int64
			username          string
			password          string
			email             string
			phone             string
			status            int32
			isVerified        bool
			isCompanyVerified bool
			defaultTenantID   int64
			role              string
			createdAt         time.Time
			updatedAt         time.Time
		)
		
		if err := rows.Scan(&id, &username, &password, &email, &phone, &status, &isVerified, &isCompanyVerified,
			&defaultTenantID, &role, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		
		u, err := user.ReconstructUser(
			id, username, email, phone, []byte(password), status, isVerified, isCompanyVerified,
			defaultTenantID, role, createdAt, updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		users = append(users, u)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// CountActiveUsers 统计活跃用户数量
func (r *UserRepository) CountActiveUsers(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE status = ?"
	
	var count int
	err := r.db.QueryRowContext(ctx, query, int32(user.UserStatusActive)).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
