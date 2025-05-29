package service

import (
	"errors"

	"wz-backend-go/internal/application/user/dto"
	domainService "wz-backend-go/internal/domain/user/service"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserApplicationService 用户应用服务
type UserApplicationService struct {
	userDomainService *domainService.UserDomainService
}

// NewUserApplicationService 创建用户应用服务
func NewUserApplicationService(userDomainService *domainService.UserDomainService) *UserApplicationService {
	return &UserApplicationService{
		userDomainService: userDomainService,
	}
}

// Register 用户注册
func (s *UserApplicationService) Register(req *dto.UserRegisterRequest) (*dto.UserResponse, error) {
	// 转换请求参数为值对象
	username, err := valueobject.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	phone, err := valueobject.NewPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	// 调用领域服务进行注册
	user, err := s.userDomainService.Register(username, req.Password, email, phone)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "注册成功",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// Login 用户登录
func (s *UserApplicationService) Login(req *dto.UserLoginRequest, ip string, userAgent string) (*dto.UserResponse, error) {
	// 调用领域服务进行登录
	user, err := s.userDomainService.Login(req.UsernameOrEmail, req.Password, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "登录成功",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserApplicationService) GetUserByID(id int64) (*dto.UserResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateUserID(id); err != nil {
		return nil, err
	}

	userID := valueobject.NewUserID(id)

	// 获取用户
	user, err := s.userDomainService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "success",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// VerifyUser 验证用户
func (s *UserApplicationService) VerifyUser(id int64) (*dto.UserResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateUserID(id); err != nil {
		return nil, err
	}

	userID := valueobject.NewUserID(id)

	// 验证用户
	if err := s.userDomainService.VerifyUser(userID); err != nil {
		return nil, err
	}

	// 获取更新后的用户
	user, err := s.userDomainService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "验证成功",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// VerifyCompany 完成企业认证
func (s *UserApplicationService) VerifyCompany(id int64) (*dto.UserResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateUserID(id); err != nil {
		return nil, err
	}

	userID := valueobject.NewUserID(id)

	// 完成企业认证
	if err := s.userDomainService.VerifyCompany(userID); err != nil {
		return nil, err
	}

	// 获取更新后的用户
	user, err := s.userDomainService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "企业认证成功",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// ChangePassword 修改密码
func (s *UserApplicationService) ChangePassword(id int64, req *dto.ChangePasswordRequest) (*dto.UserResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateUserID(id); err != nil {
		return nil, err
	}

	userID := valueobject.NewUserID(id)

	// 修改密码
	if err := s.userDomainService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}

	// 获取更新后的用户
	user, err := s.userDomainService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.UserResponse{
		Code:    0,
		Message: "密码修改成功",
		Data:    dto.ToUserDTO(user),
	}, nil
}

// GetUsers 分页获取用户列表
func (s *UserApplicationService) GetUsers(page, pageSize int) (*dto.UsersResponse, error) {
	// 调用领域服务获取用户列表
	users, total, err := s.userDomainService.GetUsers(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToUsersResponse(users, total), nil
}

// RecordUserBehavior 记录用户行为
func (s *UserApplicationService) RecordUserBehavior(userID int64, req *dto.UserBehaviorRequest) error {
	// 转换为值对象
	if err := valueobject.ValidateUserID(userID); err != nil {
		return err
	}

	uid := valueobject.NewUserID(userID)

	if err := valueobject.ValidateUserID(req.ResourceID); err != nil {
		return err
	}

	resourceID := valueobject.NewUserID(req.ResourceID)

	// 记录用户行为
	return s.userDomainService.RecordUserBehavior(uid, req.Action, req.ResourceType, resourceID)
}

// GetUserBehaviors 获取用户行为列表
func (s *UserApplicationService) GetUserBehaviors(userID int64, page, pageSize int) (*dto.UserBehaviorsResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateUserID(userID); err != nil {
		return nil, err
	}

	uid := valueobject.NewUserID(userID)

	// 获取用户行为列表
	behaviors, total, err := s.userDomainService.GetUserBehaviors(uid, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToUserBehaviorsResponse(behaviors, total), nil
}
