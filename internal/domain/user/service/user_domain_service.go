package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/event"
	"wz-backend-go/internal/domain/user/repository"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserDomainService 用户领域服务
type UserDomainService struct {
	userRepository repository.UserRepository
	eventPublisher repository.EventPublisher
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(
	userRepository repository.UserRepository,
	eventPublisher repository.EventPublisher,
) *UserDomainService {
	return &UserDomainService{
		userRepository: userRepository,
		eventPublisher: eventPublisher,
	}
}

// RegisterUser 注册用户
func (s *UserDomainService) RegisterUser(
	ctx context.Context,
	username valueobject.Username,
	password string,
	email valueobject.Email,
	phone valueobject.Phone,
) (*entity.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 如果手机号不为空，检查是否已存在
	if !phone.IsEmpty() {
		exists, err = s.userRepository.ExistsByPhone(ctx, phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 创建用户
	user, err := entity.NewUser(username, password, email, phone)
	if err != nil {
		return nil, err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	// 发布用户注册事件
	s.eventPublisher.Publish(event.NewUserRegisteredEvent(user))

	return user, nil
}

// AuthenticateUser 用户认证
func (s *UserDomainService) AuthenticateUser(
	ctx context.Context,
	usernameOrEmail string,
	password string,
	clientIP, userAgent string,
) (*entity.User, error) {
	// 根据用户名或邮箱查找用户
	var user *entity.User
	var err error

	// 尝试用邮箱查找
	if email, emailErr := valueobject.NewEmail(usernameOrEmail); emailErr == nil {
		user, err = s.userRepository.FindByEmail(ctx, email)
	}

	// 如果邮箱查找失败，尝试用用户名查找
	if user == nil && err == nil {
		if username, usernameErr := valueobject.NewUsername(usernameOrEmail); usernameErr == nil {
			user, err = s.userRepository.FindByUsername(ctx, username)
		}
	}

	if err != nil || user == nil {
		// 发布登录失败事件（没有找到用户）
		if user != nil {
			s.eventPublisher.Publish(event.NewUserLoginFailedEvent(user, "用户不存在", clientIP, userAgent))
		}
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否可以登录
	if !user.CanLogin() {
		reason := "账户状态不允许登录"
		if user.Status().IsLocked() {
			reason = "账户已被锁定"
		} else if user.Status().IsBanned() {
			reason = "账户已被封禁"
		} else if user.Status().IsSuspended() {
			reason = "账户已被暂停"
		} else if user.LoginFailCount() >= 5 {
			reason = "登录失败次数过多，账户已锁定"
		}

		s.eventPublisher.Publish(event.NewUserLoginFailedEvent(user, reason, clientIP, userAgent))
		return nil, errors.New(reason)
	}

	// 验证密码
	if !user.CheckPassword(password) {
		// 记录登录失败
		user.RecordLoginFail()
		s.userRepository.Save(user)

		s.eventPublisher.Publish(event.NewUserLoginFailedEvent(user, "密码错误", clientIP, userAgent))
		return nil, errors.New("用户名或密码错误")
	}

	// 登录成功，记录登录信息
	user.RecordLogin()
	s.userRepository.Save(user)

	// 发布登录成功事件
	s.eventPublisher.Publish(event.NewUserLoginEvent(user, clientIP, userAgent))

	return user, nil
}

// ActivateUser 激活用户
func (s *UserDomainService) ActivateUser(ctx context.Context, userID valueobject.UserID, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Activate()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserActivatedEvent(user))
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, "用户激活", operator))

	return nil
}

// VerifyUser 验证用户邮箱
func (s *UserDomainService) VerifyUser(ctx context.Context, userID valueobject.UserID, verificationCode string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 这里应该验证验证码的有效性
	// 简化实现，实际项目中需要验证码服务
	if verificationCode == "" {
		return errors.New("验证码不能为空")
	}

	// 验证用户
	user.Verify()

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布验证事件
	s.eventPublisher.Publish(event.NewUserVerifiedEvent(user, "email"))

	return nil
}

// VerifyCompany 企业验证
func (s *UserDomainService) VerifyCompany(ctx context.Context, userID valueobject.UserID, companyInfo map[string]interface{}) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 检查是否已完成邮箱验证
	if !user.IsVerified() {
		return errors.New("请先完成邮箱验证")
	}

	// 检查是否为企业邮箱
	if !user.Email().IsBusinessEmail() {
		return errors.New("请使用企业邮箱进行企业认证")
	}

	// 这里应该验证企业信息的真实性
	// 简化实现，实际项目中需要企业验证服务

	// 企业验证
	user.VerifyCompany()

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布验证事件
	s.eventPublisher.Publish(event.NewUserVerifiedEvent(user, "company"))

	return nil
}

// SuspendUser 暂停用户
func (s *UserDomainService) SuspendUser(ctx context.Context, userID valueobject.UserID, reason, operator string, duration *time.Duration) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Suspend()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserSuspendedEvent(user, reason, operator, duration))
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, reason, operator))

	return nil
}

// UnsuspendUser 恢复用户
func (s *UserDomainService) UnsuspendUser(ctx context.Context, userID valueobject.UserID, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Unsuspend()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, "恢复用户", operator))

	return nil
}

// BanUser 封禁用户
func (s *UserDomainService) BanUser(ctx context.Context, userID valueobject.UserID, reason, operator string, duration *time.Duration) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Ban()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserBannedEvent(user, reason, operator, duration))
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, reason, operator))

	return nil
}

// UnbanUser 解封用户
func (s *UserDomainService) UnbanUser(ctx context.Context, userID valueobject.UserID, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Unban()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, "解封用户", operator))

	return nil
}

// LockUser 锁定用户
func (s *UserDomainService) LockUser(ctx context.Context, userID valueobject.UserID, reason, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Lock()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, reason, operator))

	return nil
}

// UnlockUser 解锁用户
func (s *UserDomainService) UnlockUser(ctx context.Context, userID valueobject.UserID, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	oldStatus := user.Status()
	err = user.Unlock()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布事件
	s.eventPublisher.Publish(event.NewUserStatusChangedEvent(user, oldStatus, "解锁用户", operator))

	return nil
}

// ChangePassword 修改密码
func (s *UserDomainService) ChangePassword(
	ctx context.Context,
	userID valueobject.UserID,
	oldPassword, newPassword string,
	clientIP, userAgent string,
) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !user.CheckPassword(oldPassword) {
		return errors.New("原密码错误")
	}

	// 设置新密码
	err = user.SetPassword(newPassword)
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布密码修改事件
	s.eventPublisher.Publish(event.NewUserPasswordChangedEvent(user, clientIP, userAgent))

	return nil
}

// UpdateProfile 更新用户资料
func (s *UserDomainService) UpdateProfile(
	ctx context.Context,
	userID valueobject.UserID,
	updates map[string]interface{},
) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	updatedFields := make(map[string]interface{})

	// 更新用户名
	if username, ok := updates["username"].(string); ok && username != "" {
		newUsername, err := valueobject.NewUsername(username)
		if err != nil {
			return err
		}

		// 检查用户名是否已存在
		exists, err := s.userRepository.ExistsByUsername(ctx, newUsername)
		if err != nil {
			return err
		}
		if exists && !user.Username().IsEquals(newUsername) {
			return errors.New("用户名已存在")
		}

		err = user.UpdateUsername(newUsername)
		if err != nil {
			return err
		}
		updatedFields["username"] = username
	}

	// 更新邮箱
	if email, ok := updates["email"].(string); ok && email != "" {
		newEmail, err := valueobject.NewEmail(email)
		if err != nil {
			return err
		}

		// 检查邮箱是否已存在
		exists, err := s.userRepository.ExistsByEmail(ctx, newEmail)
		if err != nil {
			return err
		}
		if exists && !user.Email().IsEquals(newEmail) {
			return errors.New("邮箱已被注册")
		}

		err = user.UpdateEmail(newEmail)
		if err != nil {
			return err
		}
		updatedFields["email"] = email
	}

	// 更新手机号
	if phone, ok := updates["phone"].(string); ok {
		var newPhone valueobject.Phone
		var err error

		if phone != "" {
			newPhone, err = valueobject.NewPhone(phone)
			if err != nil {
				return err
			}

			// 检查手机号是否已存在
			exists, err := s.userRepository.ExistsByPhone(ctx, newPhone)
			if err != nil {
				return err
			}
			if exists && !user.Phone().IsEquals(newPhone) {
				return errors.New("手机号已被注册")
			}
		} else {
			// 清空手机号
			newPhone = valueobject.Phone{}
		}

		err = user.UpdatePhone(newPhone)
		if err != nil {
			return err
		}
		updatedFields["phone"] = phone
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布资料更新事件
	if len(updatedFields) > 0 {
		s.eventPublisher.Publish(event.NewUserProfileUpdatedEvent(user, updatedFields))
	}

	return nil
}

// RecordUserBehavior 记录用户行为
func (s *UserDomainService) RecordUserBehavior(
	ctx context.Context,
	userID valueobject.UserID,
	action, resourceType string,
	resourceID int64,
	clientIP, userAgent string,
) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 发布用户行为记录事件
	s.eventPublisher.Publish(event.NewUserBehaviorRecordedEvent(
		user, action, resourceType, resourceID, clientIP, userAgent,
	))

	return nil
}

// CanUserPerformAction 检查用户是否可以执行某个操作
func (s *UserDomainService) CanUserPerformAction(
	ctx context.Context,
	userID valueobject.UserID,
	action string,
) (bool, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return false, errors.New("用户不存在")
	}

	return user.CanPerformAction(action), nil
}

// DeleteUser 删除用户
func (s *UserDomainService) DeleteUser(ctx context.Context, userID valueobject.UserID, reason, operator string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	err = user.Delete()
	if err != nil {
		return err
	}

	// 保存用户
	err = s.userRepository.Save(user)
	if err != nil {
		return err
	}

	// 发布删除事件
	s.eventPublisher.Publish(event.NewUserDeletedEvent(user, reason, operator))

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserDomainService) GetUserByID(ctx context.Context, userID valueobject.UserID) (*entity.User, error) {
	return s.userRepository.FindByID(ctx, userID)
}

// GetUserByUsername 根据用户名获取用户
func (s *UserDomainService) GetUserByUsername(ctx context.Context, username valueobject.Username) (*entity.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserDomainService) GetUserByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	return s.userRepository.FindByEmail(ctx, email)
} 