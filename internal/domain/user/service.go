package user

import (
	"context"
	"errors"
	"time"

	"wz-project/wz-backend-go/internal/domain/event"
)

// Service 用户领域服务
// 处理跨实体的复杂业务逻辑，不属于单个实体的职责
type Service interface {
	// Register 用户注册
	Register(ctx context.Context, username, password, email, phone string) (*User, error)
	
	// Login 用户登录
	Login(ctx context.Context, usernameOrEmail, password string, ip, userAgent string) (*User, *UserLoginLog, error)
	
	// VerifyUser 验证用户
	VerifyUser(ctx context.Context, userID UserID, verificationCode string) error
	
	// ResetPassword 重置密码
	ResetPassword(ctx context.Context, userID UserID, oldPassword, newPassword string) error
	
	// UpdateUserProfile 更新用户资料
	UpdateUserProfile(ctx context.Context, userID UserID, profile *UserProfile) error
	
	// ChangeUserRole 修改用户角色
	ChangeUserRole(ctx context.Context, userID UserID, role UserRole, operatorID UserID) error
	
	// DisableUser 禁用用户
	DisableUser(ctx context.Context, userID UserID, reason string, operatorID UserID) error
	
	// EnableUser 启用用户
	EnableUser(ctx context.Context, userID UserID, operatorID UserID) error
	
	// VerifyCompany 企业认证
	VerifyCompany(ctx context.Context, userID UserID, companyVerification *CompanyVerification) error
}

// ServiceImpl 用户领域服务实现
type ServiceImpl struct {
	userRepo     Repository
	profileRepo  ProfileRepository
	loginRepo    LoginRepository
	behaviorRepo BehaviorRepository
	eventBus     event.Bus
}

// NewService 创建用户领域服务
func NewService(
	userRepo Repository,
	profileRepo ProfileRepository,
	loginRepo LoginRepository,
	behaviorRepo BehaviorRepository,
	eventBus event.Bus,
) Service {
	return &ServiceImpl{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		loginRepo:    loginRepo,
		behaviorRepo: behaviorRepo,
		eventBus:     eventBus,
	}
}

// Register 用户注册
func (s *ServiceImpl) Register(ctx context.Context, username, password, email, phone string) (*User, error) {
	// 验证用户名是否已存在
	usernameObj, err := NewUsername(username)
	if err != nil {
		return nil, err
	}
	
	exists, err := s.userRepo.ExistsByUsername(ctx, usernameObj)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}
	
	// 验证邮箱是否已存在
	emailObj, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	
	if !emailObj.IsEmpty() {
		exists, err = s.userRepo.ExistsByEmail(ctx, emailObj)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("邮箱已被注册")
		}
	}
	
	// 验证手机号是否已存在
	phoneObj, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}
	
	if !phoneObj.IsEmpty() {
		exists, err = s.userRepo.ExistsByPhone(ctx, phoneObj)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("手机号已被注册")
		}
	}
	
	// 创建用户
	user, err := NewUser(usernameObj, emailObj, phoneObj, password)
	if err != nil {
		return nil, err
	}
	
	// 保存用户
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}
	
	// 发布用户注册事件
	s.eventBus.Publish(NewUserRegisteredEvent(user.ID(), user.Username().Value()))
	
	return user, nil
}

// Login 用户登录
func (s *ServiceImpl) Login(ctx context.Context, usernameOrEmail, password string, ip, userAgent string) (*User, *UserLoginLog, error) {
	var user *User
	var err error
	
	// 判断是用户名还是邮箱
	if emailRegex.MatchString(usernameOrEmail) {
		emailObj, err := NewEmail(usernameOrEmail)
		if err != nil {
			return nil, nil, err
		}
		user, err = s.userRepo.FindByEmail(ctx, emailObj)
	} else {
		usernameObj, err := NewUsername(usernameOrEmail)
		if err != nil {
			return nil, nil, err
		}
		user, err = s.userRepo.FindByUsername(ctx, usernameObj)
	}
	
	if err != nil {
		return nil, nil, errors.New("用户名或密码错误")
	}
	
	// 验证用户状态
	if !user.CanLogin() {
		return nil, nil, errors.New("账号已被禁用")
	}
	
	// 验证密码
	if !user.VerifyPassword(password) {
		// 记录失败登录日志
		loginLog := &UserLoginLog{
			userID:    user.ID(),
			loginType: "password",
			ip:        ip,
			userAgent: userAgent,
			status:    0, // 失败
			errorMsg:  "密码错误",
			loginAt:   time.Now(),
		}
		s.loginRepo.SaveLoginLog(ctx, loginLog)
		
		return nil, nil, errors.New("用户名或密码错误")
	}
	
	// 记录成功登录日志
	loginLog := &UserLoginLog{
		userID:    user.ID(),
		loginType: "password",
		ip:        ip,
		userAgent: userAgent,
		status:    1, // 成功
		loginAt:   time.Now(),
	}
	s.loginRepo.SaveLoginLog(ctx, loginLog)
	
	// 发布用户登录事件
	s.eventBus.Publish(NewUserLoggedInEvent(user.ID(), loginLog.loginAt))
	
	return user, loginLog, nil
}

// VerifyUser 验证用户
func (s *ServiceImpl) VerifyUser(ctx context.Context, userID UserID, verificationCode string) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// TODO: 验证码校验逻辑
	// 这里简化处理，实际应该有验证码匹配逻辑
	
	// 验证用户
	user.VerifyUser()
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 发布用户验证事件
	s.eventBus.Publish(NewUserVerifiedEvent(userID))
	
	return nil
}

// ResetPassword 重置密码
func (s *ServiceImpl) ResetPassword(ctx context.Context, userID UserID, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// 修改密码
	if err := user.ChangePassword(oldPassword, newPassword); err != nil {
		return err
	}
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 发布密码修改事件
	s.eventBus.Publish(NewPasswordChangedEvent(userID))
	
	return nil
}

// UpdateUserProfile 更新用户资料
func (s *ServiceImpl) UpdateUserProfile(ctx context.Context, userID UserID, profile *UserProfile) error {
	// 验证用户ID
	if profile.userID != userID {
		return errors.New("用户ID不匹配")
	}
	
	// 保存或更新用户资料
	existingProfile, err := s.profileRepo.FindProfileByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		// 更新现有资料
		if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
			return err
		}
	} else {
		// 创建新资料
		if err := s.profileRepo.SaveProfile(ctx, profile); err != nil {
			return err
		}
	}
	
	// 发布资料更新事件
	s.eventBus.Publish(NewProfileUpdatedEvent(userID))
	
	return nil
}

// ChangeUserRole 修改用户角色
func (s *ServiceImpl) ChangeUserRole(ctx context.Context, userID UserID, role UserRole, operatorID UserID) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// 获取操作员
	operator, err := s.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return err
	}
	
	// 权限检查：只有平台管理员或租户管理员可以修改角色
	if operator.Role() != RolePlatformAdmin && operator.Role() != RoleTenantAdmin {
		return errors.New("无权限执行此操作")
	}
	
	// 修改角色
	user.ChangeRole(role)
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 记录操作日志
	behaviorLog := &UserBehaviorLog{
		userID:       operatorID,
		action:       "change_role",
		resourceType: "user",
		resourceID:   userID.Value(),
		createdAt:    time.Now(),
	}
	s.behaviorRepo.SaveBehaviorLog(ctx, behaviorLog)
	
	// 发布角色变更事件
	s.eventBus.Publish(NewRoleChangedEvent(userID, role))
	
	return nil
}

// DisableUser 禁用用户
func (s *ServiceImpl) DisableUser(ctx context.Context, userID UserID, reason string, operatorID UserID) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// 获取操作员
	operator, err := s.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return err
	}
	
	// 权限检查：只有平台管理员或租户管理员可以禁用用户
	if operator.Role() != RolePlatformAdmin && operator.Role() != RoleTenantAdmin {
		return errors.New("无权限执行此操作")
	}
	
	// 禁止禁用平台管理员
	if user.Role() == RolePlatformAdmin && operator.Role() != RolePlatformAdmin {
		return errors.New("无权限禁用平台管理员")
	}
	
	// 修改状态
	if err := user.ChangeStatus(UserStatusDisabled); err != nil {
		return err
	}
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 记录操作日志
	behaviorLog := &UserBehaviorLog{
		userID:       operatorID,
		action:       "disable_user",
		resourceType: "user",
		resourceID:   userID.Value(),
		createdAt:    time.Now(),
	}
	s.behaviorRepo.SaveBehaviorLog(ctx, behaviorLog)
	
	// 发布用户禁用事件
	s.eventBus.Publish(NewUserDisabledEvent(userID, reason, operatorID))
	
	return nil
}

// EnableUser 启用用户
func (s *ServiceImpl) EnableUser(ctx context.Context, userID UserID, operatorID UserID) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// 获取操作员
	operator, err := s.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return err
	}
	
	// 权限检查：只有平台管理员或租户管理员可以启用用户
	if operator.Role() != RolePlatformAdmin && operator.Role() != RoleTenantAdmin {
		return errors.New("无权限执行此操作")
	}
	
	// 修改状态
	if err := user.ChangeStatus(UserStatusActive); err != nil {
		return err
	}
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 记录操作日志
	behaviorLog := &UserBehaviorLog{
		userID:       operatorID,
		action:       "enable_user",
		resourceType: "user",
		resourceID:   userID.Value(),
		createdAt:    time.Now(),
	}
	s.behaviorRepo.SaveBehaviorLog(ctx, behaviorLog)
	
	// 发布用户启用事件
	s.eventBus.Publish(NewUserEnabledEvent(userID, operatorID))
	
	return nil
}

// VerifyCompany 企业认证
func (s *ServiceImpl) VerifyCompany(ctx context.Context, userID UserID, companyVerification *CompanyVerification) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	
	// TODO: 企业认证逻辑
	// 这里简化处理，实际应该有认证资料校验逻辑
	
	// 设置企业认证标记
	user.VerifyCompany()
	
	// 保存更新
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}
	
	// 发布企业认证事件
	s.eventBus.Publish(NewCompanyVerifiedEvent(userID, companyVerification.CompanyName))
	
	return nil
}

// CompanyVerification 企业认证信息
type CompanyVerification struct {
	UserID             UserID
	CompanyType        int
	CompanyName        string
	BusinessLicense    string
	OrgCodeCert        string
	UnifiedSocialCredit string
	ContactPerson      string
	ContactPhone       string
}
