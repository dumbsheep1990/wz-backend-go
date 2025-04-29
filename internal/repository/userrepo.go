package repository

import (
	"context"
	
	"wz-backend-go/internal/domain/model"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// 创建新用户
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	
	// 根据ID获取用户
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	
	// 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	
	// 根据邮箱获取用户
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	
	// 根据手机号获取用户
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	
	// 更新用户信息
	UpdateUser(ctx context.Context, user *model.User) error
	
	// 删除用户（软删除）
	DeleteUser(ctx context.Context, id int64) error
	
	// 验证用户密码
	VerifyPassword(ctx context.Context, username, password string) (bool, *model.User, error)
	
	// 更新用户密码
	UpdatePassword(ctx context.Context, id int64, newPassword string) error
	
	// 更新用户角色
	UpdateUserRole(ctx context.Context, id int64, role string) error
	
	// 列出所有用户（分页）
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	
	// 根据角色查询用户
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*model.User, int64, error)
	
	// 检查用户名是否存在
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	
	// 检查邮箱是否存在
	IsEmailExists(ctx context.Context, email string) (bool, error)
	
	// 检查手机号是否存在
	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	
	// 获取用户的所有租户
	GetUserTenants(ctx context.Context, userID int64) ([]*model.Tenant, error)
	
	// 保存用户登录记录
	SaveLoginRecord(ctx context.Context, userID int64, ip, device, location string) error
	
	// 获取用户登录记录
	GetLoginRecords(ctx context.Context, userID int64, limit int) ([]*model.UserLoginLog, error)
}
