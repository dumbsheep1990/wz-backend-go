package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"wz-backend-go/internal/domain/user/valueobject"
)

// User 用户实体
type User struct {
	id                valueobject.UserID
	username          valueobject.Username
	password          string
	email             valueobject.Email
	phone             valueobject.Phone
	status            valueobject.UserStatus
	isVerified        bool
	isCompanyVerified bool
	createdAt         time.Time
	updatedAt         time.Time
}

// NewUser 创建新用户
func NewUser(
	username valueobject.Username,
	password string,
	email valueobject.Email,
	phone valueobject.Phone,
) (*User, error) {
	// 验证必填参数
	if username.Value() == "" {
		return nil, errors.New("用户名不能为空")
	}
	if password == "" {
		return nil, errors.New("密码不能为空")
	}
	if email.Value() == "" {
		return nil, errors.New("邮箱不能为空")
	}

	// 默认状态为活跃
	status, _ := valueobject.NewUserStatus(1)

	now := time.Now()

	user := &User{
		username:          username,
		email:             email,
		phone:             phone,
		status:            status,
		isVerified:        false,
		isCompanyVerified: false,
		createdAt:         now,
		updatedAt:         now,
	}

	// 设置密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

// ReconstructUser 从存储中重建用户实体
func ReconstructUser(
	id valueobject.UserID,
	username valueobject.Username,
	password string,
	email valueobject.Email,
	phone valueobject.Phone,
	status valueobject.UserStatus,
	isVerified bool,
	isCompanyVerified bool,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:                id,
		username:          username,
		password:          password,
		email:             email,
		phone:             phone,
		status:            status,
		isVerified:        isVerified,
		isCompanyVerified: isCompanyVerified,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// ID 获取用户ID
func (u *User) ID() valueobject.UserID {
	return u.id
}

// SetID 设置用户ID
func (u *User) SetID(id valueobject.UserID) {
	u.id = id
}

// Username 获取用户名
func (u *User) Username() valueobject.Username {
	return u.username
}

// Password 获取密码（哈希值）
func (u *User) Password() string {
	return u.password
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	// 使用bcrypt加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.password = string(hashedPassword)
	u.updatedAt = time.Now()
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

// Email 获取邮箱
func (u *User) Email() valueobject.Email {
	return u.email
}

// SetEmail 更新邮箱
func (u *User) SetEmail(email valueobject.Email) {
	u.email = email
	u.updatedAt = time.Now()
}

// Phone 获取手机号
func (u *User) Phone() valueobject.Phone {
	return u.phone
}

// SetPhone 更新手机号
func (u *User) SetPhone(phone valueobject.Phone) {
	u.phone = phone
	u.updatedAt = time.Now()
}

// Status 获取状态
func (u *User) Status() valueobject.UserStatus {
	return u.status
}

// SetStatus 更新状态
func (u *User) SetStatus(status valueobject.UserStatus) {
	u.status = status
	u.updatedAt = time.Now()
}

// IsVerified 获取是否已验证
func (u *User) IsVerified() bool {
	return u.isVerified
}

// Verify 验证用户
func (u *User) Verify() {
	u.isVerified = true
	u.updatedAt = time.Now()
}

// IsCompanyVerified 获取是否已完成企业认证
func (u *User) IsCompanyVerified() bool {
	return u.isCompanyVerified
}

// VerifyCompany 完成企业认证
func (u *User) VerifyCompany() {
	u.isCompanyVerified = true
	u.updatedAt = time.Now()
}

// CreatedAt 获取创建时间
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt 获取更新时间
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// UserBehavior 用户行为实体
type UserBehavior struct {
	id           valueobject.UserID
	userID       valueobject.UserID
	action       string
	resourceType string
	resourceID   valueobject.UserID
	createdAt    time.Time
}

// NewUserBehavior 创建用户行为
func NewUserBehavior(
	userID valueobject.UserID,
	action string,
	resourceType string,
	resourceID valueobject.UserID,
) (*UserBehavior, error) {
	if userID.Value() <= 0 {
		return nil, errors.New("用户ID无效")
	}
	if action == "" {
		return nil, errors.New("行为不能为空")
	}
	if resourceType == "" {
		return nil, errors.New("资源类型不能为空")
	}

	return &UserBehavior{
		userID:       userID,
		action:       action,
		resourceType: resourceType,
		resourceID:   resourceID,
		createdAt:    time.Now(),
	}, nil
}

// ReconstructUserBehavior 从存储中重建用户行为实体
func ReconstructUserBehavior(
	id valueobject.UserID,
	userID valueobject.UserID,
	action string,
	resourceType string,
	resourceID valueobject.UserID,
	createdAt time.Time,
) *UserBehavior {
	return &UserBehavior{
		id:           id,
		userID:       userID,
		action:       action,
		resourceType: resourceType,
		resourceID:   resourceID,
		createdAt:    createdAt,
	}
}

// ID 获取行为ID
func (b *UserBehavior) ID() valueobject.UserID {
	return b.id
}

// SetID 设置行为ID
func (b *UserBehavior) SetID(id valueobject.UserID) {
	b.id = id
}

// UserID 获取用户ID
func (b *UserBehavior) UserID() valueobject.UserID {
	return b.userID
}

// Action 获取行为
func (b *UserBehavior) Action() string {
	return b.action
}

// ResourceType 获取资源类型
func (b *UserBehavior) ResourceType() string {
	return b.resourceType
}

// ResourceID 获取资源ID
func (b *UserBehavior) ResourceID() valueobject.UserID {
	return b.resourceID
}

// CreatedAt 获取创建时间
func (b *UserBehavior) CreatedAt() time.Time {
	return b.createdAt
}
