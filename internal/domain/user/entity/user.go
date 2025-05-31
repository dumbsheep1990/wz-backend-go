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
	lastLoginAt       *time.Time
	loginFailCount    int32
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
	if username.IsEmpty() {
		return nil, errors.New("用户名不能为空")
	}
	if password == "" {
		return nil, errors.New("密码不能为空")
	}
	if email.Value() == "" {
		return nil, errors.New("邮箱不能为空")
	}

	// 验证密码强度
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// 默认状态为未激活
	status, _ := valueobject.NewUserStatus(int32(valueobject.UserStatusInactive))

	now := time.Now()

	user := &User{
		username:          username,
		email:             email,
		phone:             phone,
		status:            status,
		isVerified:        false,
		isCompanyVerified: false,
		loginFailCount:    0,
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
	lastLoginAt *time.Time,
	loginFailCount int32,
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
		lastLoginAt:       lastLoginAt,
		loginFailCount:    loginFailCount,
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

// UpdateUsername 更新用户名
func (u *User) UpdateUsername(username valueobject.Username) error {
	if username.IsEmpty() {
		return errors.New("用户名不能为空")
	}
	
	u.username = username
	u.updatedAt = time.Now()
	return nil
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

	// 验证密码强度
	if err := validatePassword(password); err != nil {
		return err
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

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(email valueobject.Email) error {
	if email.Value() == "" {
		return errors.New("邮箱不能为空")
	}
	
	u.email = email
	u.isVerified = false // 更新邮箱后需要重新验证
	u.updatedAt = time.Now()
	return nil
}

// Phone 获取手机号
func (u *User) Phone() valueobject.Phone {
	return u.phone
}

// UpdatePhone 更新手机号
func (u *User) UpdatePhone(phone valueobject.Phone) error {
	u.phone = phone
	u.updatedAt = time.Now()
	return nil
}

// Status 获取状态
func (u *User) Status() valueobject.UserStatus {
	return u.status
}

// ChangeStatus 更改用户状态
func (u *User) ChangeStatus(newStatus valueobject.UserStatus) error {
	if !u.status.CanTransitionTo(newStatus) {
		return errors.New("不允许的状态转换: " + u.status.String() + " -> " + newStatus.String())
	}
	
	u.status = newStatus
	u.updatedAt = time.Now()
	return nil
}

// Activate 激活用户
func (u *User) Activate() error {
	return u.ChangeStatus(valueobject.UserStatusActive)
}

// Suspend 暂停用户
func (u *User) Suspend() error {
	return u.ChangeStatus(valueobject.UserStatusSuspended)
}

// Unsuspend 恢复用户
func (u *User) Unsuspend() error {
	return u.ChangeStatus(valueobject.UserStatusActive)
}

// Ban 封禁用户
func (u *User) Ban() error {
	return u.ChangeStatus(valueobject.UserStatusBanned)
}

// Unban 解封用户
func (u *User) Unban() error {
	return u.ChangeStatus(valueobject.UserStatusActive)
}

// Lock 锁定用户
func (u *User) Lock() error {
	return u.ChangeStatus(valueobject.UserStatusLocked)
}

// Unlock 解锁用户
func (u *User) Unlock() error {
	u.loginFailCount = 0 // 重置登录失败次数
	return u.ChangeStatus(valueobject.UserStatusActive)
}

// Delete 删除用户（软删除）
func (u *User) Delete() error {
	return u.ChangeStatus(valueobject.UserStatusDeleted)
}

// IsVerified 获取验证状态
func (u *User) IsVerified() bool {
	return u.isVerified
}

// Verify 验证用户
func (u *User) Verify() {
	u.isVerified = true
	u.updatedAt = time.Now()
}

// IsCompanyVerified 获取企业验证状态
func (u *User) IsCompanyVerified() bool {
	return u.isCompanyVerified
}

// VerifyCompany 企业验证
func (u *User) VerifyCompany() {
	u.isCompanyVerified = true
	u.updatedAt = time.Now()
}

// LastLoginAt 获取最后登录时间
func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// RecordLogin 记录登录
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.loginFailCount = 0 // 重置登录失败次数
	u.updatedAt = now
}

// LoginFailCount 获取登录失败次数
func (u *User) LoginFailCount() int32 {
	return u.loginFailCount
}

// RecordLoginFail 记录登录失败
func (u *User) RecordLoginFail() error {
	u.loginFailCount++
	u.updatedAt = time.Now()
	
	// 连续失败5次自动锁定账户
	if u.loginFailCount >= 5 && u.status.CanBeLocked() {
		return u.Lock()
	}
	
	return nil
}

// ResetLoginFailCount 重置登录失败次数
func (u *User) ResetLoginFailCount() {
	u.loginFailCount = 0
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

// CanLogin 判断用户是否可以登录
func (u *User) CanLogin() bool {
	return u.status.CanLogin() && u.loginFailCount < 5
}

// CanPerformAction 判断用户是否可以执行某个操作
func (u *User) CanPerformAction(action string) bool {
	// 只有活跃状态的用户才能执行操作
	if !u.status.IsActive() {
		return false
	}
	
	// 根据具体操作判断权限
	switch action {
	case "create_content", "edit_profile", "comment":
		return u.isVerified
	case "create_company_content", "company_action":
		return u.isVerified && u.isCompanyVerified
	default:
		return true
	}
}

// GetDisplayName 获取显示名称
func (u *User) GetDisplayName() string {
	return u.username.Value()
}

// GetContactInfo 获取联系信息（脱敏）
func (u *User) GetContactInfo() map[string]string {
	return map[string]string{
		"email": u.email.Mask(),
		"phone": u.phone.Mask(),
	}
}

// IsBusinessUser 判断是否为企业用户
func (u *User) IsBusinessUser() bool {
	return u.isCompanyVerified && u.email.IsBusinessEmail()
}

// GetAccountAge 获取账户年龄（天数）
func (u *User) GetAccountAge() int {
	return int(time.Since(u.createdAt).Hours() / 24)
}

// IsNewUser 判断是否为新用户（注册7天内）
func (u *User) IsNewUser() bool {
	return u.GetAccountAge() <= 7
}

// NeedsPasswordReset 判断是否需要重置密码
func (u *User) NeedsPasswordReset() bool {
	// 如果账户被锁定且登录失败次数过多
	return u.status.IsLocked() && u.loginFailCount >= 5
}

// validatePassword 验证密码强度
func validatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	if len(password) > 30 {
		return errors.New("密码长度不能超过30位")
	}
	
	// 这里可以添加更复杂的密码强度验证
	// 比如：必须包含字母、数字、特殊字符等
	
	return nil
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
