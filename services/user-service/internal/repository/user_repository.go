package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"wz-backend-go/services/user-service/internal/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(ctx context.Context, user model.User) (int64, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *model.User) error
	GetUserBehaviors(ctx context.Context, userID int64, startTime, endTime time.Time) ([]model.UserBehavior, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository() UserRepository {
	// TODO: 从配置中获取数据库连接信息
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=wz_user sslmode=disable")
	if err != nil {
		// 实际项目中应该处理这个错误，而不是panic
		panic(err)
	}

	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user model.User) (int64, error) {
	query := `
		INSERT INTO users (
			username, password, email, phone, status, is_verified, is_company_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
		user.Status,
		user.IsVerified,
		user.IsCompanyVerified,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, password, email, phone, status, is_verified, is_company_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, username)

	return exists, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)

	return exists, err
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET 
			username = $1,
			password = $2,
			email = $3,
			phone = $4,
			status = $5,
			is_verified = $6,
			is_company_verified = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
		user.Status,
		user.IsVerified,
		user.IsCompanyVerified,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// GetUserBehaviors 获取用户行为
func (r *userRepository) GetUserBehaviors(ctx context.Context, userID int64, startTime, endTime time.Time) ([]model.UserBehavior, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, created_at
		FROM user_behaviors
		WHERE user_id = $1
	`

	args := []interface{}{userID}

	// 添加时间过滤条件
	if !startTime.IsZero() {
		query += " AND created_at >= $2"
		args = append(args, startTime)
	}

	if !endTime.IsZero() {
		if !startTime.IsZero() {
			query += " AND created_at <= $3"
		} else {
			query += " AND created_at <= $2"
		}
		args = append(args, endTime)
	}

	query += " ORDER BY created_at DESC"

	var behaviors []model.UserBehavior
	err := r.db.SelectContext(ctx, &behaviors, query, args...)

	return behaviors, err
}
