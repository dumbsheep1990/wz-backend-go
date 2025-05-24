package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户实体（聚合根）
type User struct {
	id                UserID
	username          Username
	email             Email
	phone             Phone
	password          []byte
	status            UserStatus
	isVerified        bool
	isCompanyVerified bool
	defaultTenantID   TenantID
	role              UserRole
	createdAt         time.Time
	updatedAt         time.Time
}

// UserRole 用户角色枚举
type UserRole string

const (
	// RolePlatformAdmin 平台管理员
	RolePlatformAdmin UserRole = "platform_admin"
	// RoleTenantAdmin 租户管理员
	RoleTenantAdmin UserRole = "tenant_admin"
	// RoleTenantUser 租户普通用户
	RoleTenantUser UserRole = "tenant_user"
	// RolePersonalUser 个人用户
	RolePersonalUser UserRole = "personal_user"
)

// NewUser 创建新用户
func NewUser(username Username, email Email, phone Phone, password string) (*User, error) {
	if username.Value() == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 密码强度验证
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	now := time.Now()
	return &User{
		username:          username,
		email:             email,
		phone:             phone,
		password:          hashedPassword,
		status:            UserStatusUnverified, // 默认未验证状态
		isVerified:        false,
		isCompanyVerified: false,
		role:              RolePersonalUser, // 默认个人用户角色
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

// ReconstructUser 从持久化数据重建用户实体
func ReconstructUser(
	id int64,
	username string,
	email string,
	phone string,
	hashedPassword []byte,
	status int32,
	isVerified bool,
	isCompanyVerified bool,
	defaultTenantID int64,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	usernameObj, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	emailObj, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	phoneObj, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}

	return &User{
		id:                NewUserID(id),
		username:          usernameObj,
		email:             emailObj,
		phone:             phoneObj,
		password:          hashedPassword,
		status:            UserStatus(status),
		isVerified:        isVerified,
		isCompanyVerified: isCompanyVerified,
		defaultTenantID:   NewTenantID(defaultTenantID),
		role:              UserRole(role),
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}, nil
}

// ID 获取用户ID
func (u *User) ID() UserID {
	return u.id
}

// Username 获取用户名
func (u *User) Username() Username {
	return u.username
}

// Email 获取用户邮箱
func (u *User) Email() Email {
	return u.email
}

// Phone 获取用户手机号
func (u *User) Phone() Phone {
	return u.phone
}

// Status 获取用户状态
func (u *User) Status() UserStatus {
	return u.status
}

// IsVerified 获取用户是否已验证
func (u *User) IsVerified() bool {
	return u.isVerified
}

// IsCompanyVerified 获取用户是否已企业认证
func (u *User) IsCompanyVerified() bool {
	return u.isCompanyVerified
}

// DefaultTenantID 获取默认租户ID
func (u *User) DefaultTenantID() TenantID {
	return u.defaultTenantID
}

// Role 获取用户角色
func (u *User) Role() UserRole {
	return u.role
}

// CreatedAt 获取创建时间
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt 获取更新时间
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword(u.password, []byte(oldPassword)); err != nil {
		return errors.New("旧密码不正确")
	}

	// 密码强度验证
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	u.password = hashedPassword
	u.updatedAt = time.Now()
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.password, []byte(password))
	return err == nil
}

// UpdateProfile 更新用户基本信息
func (u *User) UpdateProfile(username Username, email Email, phone Phone) error {
	changed := false

	if username.Value() != "" && !username.Equals(u.username) {
		u.username = username
		changed = true
	}

	if !email.IsEmpty() && !email.Equals(u.email) {
		u.email = email
		// 邮箱更改后需要重新验证
		u.isVerified = false
		changed = true
	}

	if !phone.IsEmpty() && !phone.Equals(u.phone) {
		u.phone = phone
		changed = true
	}

	if changed {
		u.updatedAt = time.Now()
	}

	return nil
}

// VerifyUser 验证用户
func (u *User) VerifyUser() {
	u.isVerified = true
	u.status = UserStatusActive
	u.updatedAt = time.Now()
}

// VerifyCompany 企业认证
func (u *User) VerifyCompany() {
	u.isCompanyVerified = true
	u.updatedAt = time.Now()
}

// ChangeStatus 改变用户状态
func (u *User) ChangeStatus(status UserStatus) error {
	if !status.IsValid() {
		return errors.New("无效的用户状态")
	}
	u.status = status
	u.updatedAt = time.Now()
	return nil
}

// ChangeRole 修改用户角色
func (u *User) ChangeRole(role UserRole) {
	u.role = role
	u.updatedAt = time.Now()
}

// SetDefaultTenant 设置默认租户
func (u *User) SetDefaultTenant(tenantID TenantID) {
	u.defaultTenantID = tenantID
	u.updatedAt = time.Now()
}

// IsActive 判断用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.status == UserStatusActive
}

// CanLogin 判断用户是否可登录
func (u *User) CanLogin() bool {
	return u.status != UserStatusDisabled
}

// validatePassword 验证密码强度
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度必须至少为8个字符")
	}

	hasLetter := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z'):
			hasLetter = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasLetter || !hasDigit {
		return errors.New("密码必须包含字母和数字")
	}

	return nil
}
