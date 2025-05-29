package service

import (
	"errors"

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

// Register 用户注册
func (s *UserDomainService) Register(
	username valueobject.Username,
	password string,
	email valueobject.Email,
	phone valueobject.Phone,
) (*entity.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepository.FindByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepository.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 如果提供了手机号，检查手机号是否已存在
	if phone.Value() != "" {
		existingUser, err = s.userRepository.FindByPhone(phone)
		if err == nil && existingUser != nil {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 创建新用户
	user, err := entity.NewUser(username, password, email, phone)
	if err != nil {
		return nil, err
	}

	// 保存用户
	if err := s.userRepository.Save(user); err != nil {
		return nil, err
	}

	// 发布用户创建事件
	createdEvent := event.NewUserCreatedEvent(user.ID(), user.Username())
	if err := s.eventPublisher.Publish(createdEvent); err != nil {
		// 记录错误但不影响用户创建
		// 可以考虑使用日志记录这个错误
	}

	return user, nil
}

// Login 用户登录
func (s *UserDomainService) Login(
	usernameOrEmail string,
	password string,
	ip string,
	userAgent string,
) (*entity.User, error) {
	var user *entity.User
	var err error

	// 尝试使用用户名查找
	username, usernameErr := valueobject.NewUsername(usernameOrEmail)
	if usernameErr == nil {
		user, err = s.userRepository.FindByUsername(username)
	}

	// 如果找不到用户，尝试使用邮箱查找
	if user == nil {
		email, emailErr := valueobject.NewEmail(usernameOrEmail)
		if emailErr == nil {
			user, err = s.userRepository.FindByEmail(email)
		}
	}

	// 找不到用户
	if user == nil || err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status() != valueobject.UserStatusActive {
		return nil, errors.New("用户账号已被锁定或未激活")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 发布用户登录事件
	loginEvent := event.NewUserLoggedInEvent(user.ID(), ip, userAgent)
	if err := s.eventPublisher.Publish(loginEvent); err != nil {
		// 记录错误但不影响登录
		// 可以考虑使用日志记录这个错误
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserDomainService) GetUserByID(userID valueobject.UserID) (*entity.User, error) {
	return s.userRepository.FindByID(userID)
}

// GetUsers 分页获取用户列表
func (s *UserDomainService) GetUsers(page, pageSize int) ([]*entity.User, int64, error) {
	return s.userRepository.FindAll(page, pageSize)
}

// VerifyUser 验证用户
func (s *UserDomainService) VerifyUser(userID valueobject.UserID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	// 如果用户已经验证过了，返回成功
	if user.IsVerified() {
		return nil
	}

	// 验证用户
	user.Verify()

	// 保存用户
	if err := s.userRepository.Save(user); err != nil {
		return err
	}

	// 发布用户验证事件
	verifiedEvent := event.NewUserVerifiedEvent(user.ID())
	if err := s.eventPublisher.Publish(verifiedEvent); err != nil {
		// 记录错误但不影响验证
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// VerifyCompany 完成企业认证
func (s *UserDomainService) VerifyCompany(userID valueobject.UserID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	// 如果用户已经完成企业认证，返回成功
	if user.IsCompanyVerified() {
		return nil
	}

	// 用户必须先完成个人认证
	if !user.IsVerified() {
		return errors.New("用户必须先完成个人认证")
	}

	// 完成企业认证
	user.VerifyCompany()

	// 保存用户
	if err := s.userRepository.Save(user); err != nil {
		return err
	}

	// 发布企业认证事件
	verifiedEvent := event.NewUserCompanyVerifiedEvent(user.ID())
	if err := s.eventPublisher.Publish(verifiedEvent); err != nil {
		// 记录错误但不影响企业认证
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserDomainService) ChangePassword(
	userID valueobject.UserID,
	oldPassword string,
	newPassword string,
) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !user.CheckPassword(oldPassword) {
		return errors.New("旧密码不正确")
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// 保存用户
	if err := s.userRepository.Save(user); err != nil {
		return err
	}

	// 发布密码修改事件
	passwordChangedEvent := event.NewUserPasswordChangedEvent(user.ID())
	if err := s.eventPublisher.Publish(passwordChangedEvent); err != nil {
		// 记录错误但不影响密码修改
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// RecordUserBehavior 记录用户行为
func (s *UserDomainService) RecordUserBehavior(
	userID valueobject.UserID,
	action string,
	resourceType string,
	resourceID valueobject.UserID,
) error {
	// 创建用户行为
	behavior, err := entity.NewUserBehavior(userID, action, resourceType, resourceID)
	if err != nil {
		return err
	}

	// 保存用户行为
	if err := s.userRepository.SaveBehavior(behavior); err != nil {
		return err
	}

	// 发布用户行为记录事件
	behaviorEvent := event.NewUserBehaviorRecordedEvent(userID, action, resourceType, resourceID)
	if err := s.eventPublisher.Publish(behaviorEvent); err != nil {
		// 记录错误但不影响行为记录
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// GetUserBehaviors 获取用户行为列表
func (s *UserDomainService) GetUserBehaviors(userID valueobject.UserID, page, pageSize int) ([]*entity.UserBehavior, int64, error) {
	return s.userRepository.FindBehaviorsByUserID(userID, page, pageSize)
}
