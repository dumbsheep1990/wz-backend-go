package user

import (
	"context"
	"errors"
	"log"
	"time"

	"wz-backend-go/internal/domain/user"
)

// UserApplicationService 用户应用服务
type UserApplicationService struct {
	userRepo    user.Repository
	userService user.Service
	eventBus    EventBus
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event interface{}) error
}

// NewUserApplicationService 创建用户应用服务
func NewUserApplicationService(repo user.Repository, service user.Service, eventBus EventBus) *UserApplicationService {
	return &UserApplicationService{
		userRepo:    repo,
		userService: service,
		eventBus:    eventBus,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string
	Password string
	Email    string
	Phone    string
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   int64
	Username string
	Token    string
}

// Register 用户注册
func (s *UserApplicationService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// 创建值对象
	username, err := user.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	phone, err := user.NewPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		log.Printf("检查用户名失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		log.Printf("检查邮箱失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	if exists {
		return nil, errors.New("邮箱已注册")
	}

	// 如果手机号不为空，检查手机号是否已存在
	if req.Phone != "" {
		exists, err = s.userRepo.ExistsByPhone(ctx, phone)
		if err != nil {
			log.Printf("检查手机号失败: %v", err)
			return nil, errors.New("服务器内部错误")
		}

		if exists {
			return nil, errors.New("手机号已注册")
		}
	}

	// 通过领域服务创建用户
	newUser, err := s.userService.Register(ctx, username, email, phone, req.Password)
	if err != nil {
		log.Printf("创建用户失败: %v", err)
		return nil, errors.New("用户注册失败")
	}

	// 保存用户
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		log.Printf("保存用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 发布用户注册事件
	event := user.NewUserRegisteredEvent(newUser.ID(), newUser.Username(), newUser.Email(), time.Now())
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("发布用户注册事件失败: %v", err)
		// 不影响主流程，继续执行
	}

	// 生成令牌（实际中应使用JWT或其他令牌机制）
	token, err := generateToken(newUser.ID().Value(), newUser.Username().Value())
	if err != nil {
		log.Printf("生成令牌失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &RegisterResponse{
		UserID:   newUser.ID().Value(),
		Username: newUser.Username().Value(),
		Token:    token,
	}, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID   int64
	Username string
	Token    string
}

// Login 用户登录
func (s *UserApplicationService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 创建用户名值对象
	username, err := user.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// 查找用户
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户名或密码错误")
	}

	// 通过领域服务验证密码
	if err := s.userService.Login(ctx, existingUser, req.Password); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 发布用户登录事件
	event := user.NewUserLoggedInEvent(existingUser.ID(), existingUser.Username(), time.Now())
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("发布用户登录事件失败: %v", err)
		// 不影响主流程，继续执行
	}

	// 生成令牌
	token, err := generateToken(existingUser.ID().Value(), existingUser.Username().Value())
	if err != nil {
		log.Printf("生成令牌失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &LoginResponse{
		UserID:   existingUser.ID().Value(),
		Username: existingUser.Username().Value(),
		Token:    token,
	}, nil
}

// GetUserRequest 获取用户请求
type GetUserRequest struct {
	UserID int64
}

// GetUserResponse 获取用户响应
type GetUserResponse struct {
	UserID            int64
	Username          string
	Email             string
	Phone             string
	Status            int32
	IsVerified        bool
	IsCompanyVerified bool
	DefaultTenantID   int64
	Role              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// GetUser 获取用户信息
func (s *UserApplicationService) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	// 创建用户ID值对象
	userID := user.NewUserID(req.UserID)

	// 查找用户
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	return &GetUserResponse{
		UserID:            existingUser.ID().Value(),
		Username:          existingUser.Username().Value(),
		Email:             existingUser.Email().Value(),
		Phone:             existingUser.Phone().Value(),
		Status:            int32(existingUser.Status()),
		IsVerified:        existingUser.IsVerified(),
		IsCompanyVerified: existingUser.IsCompanyVerified(),
		DefaultTenantID:   existingUser.DefaultTenantID().Value(),
		Role:              string(existingUser.Role()),
		CreatedAt:         existingUser.CreatedAt(),
		UpdatedAt:         existingUser.UpdatedAt(),
	}, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	UserID   int64
	Username string
	Email    string
	Phone    string
}

// UpdateUserResponse 更新用户响应
type UpdateUserResponse struct {
	UserID            int64
	Username          string
	Email             string
	Phone             string
	Status            int32
	IsVerified        bool
	IsCompanyVerified bool
	DefaultTenantID   int64
	Role              string
	UpdatedAt         time.Time
}

// UpdateUser 更新用户信息
func (s *UserApplicationService) UpdateUser(ctx context.Context, req UpdateUserRequest) (*UpdateUserResponse, error) {
	// 创建用户ID值对象
	userID := user.NewUserID(req.UserID)

	// 查找用户
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// 如果更新用户名
	if req.Username != "" && req.Username != existingUser.Username().Value() {
		username, err := user.NewUsername(req.Username)
		if err != nil {
			return nil, err
		}

		// 检查用户名是否已存在
		exists, err := s.userRepo.ExistsByUsername(ctx, username)
		if err != nil {
			log.Printf("检查用户名失败: %v", err)
			return nil, errors.New("服务器内部错误")
		}

		if exists {
			return nil, errors.New("用户名已存在")
		}

		// 更新用户名
		existingUser.UpdateUsername(username)
	}

	// 如果更新邮箱
	if req.Email != "" && req.Email != existingUser.Email().Value() {
		email, err := user.NewEmail(req.Email)
		if err != nil {
			return nil, err
		}

		// 检查邮箱是否已存在
		exists, err := s.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			log.Printf("检查邮箱失败: %v", err)
			return nil, errors.New("服务器内部错误")
		}

		if exists {
			return nil, errors.New("邮箱已注册")
		}

		// 更新邮箱
		existingUser.UpdateEmail(email)
	}

	// 如果更新手机号
	if req.Phone != "" && req.Phone != existingUser.Phone().Value() {
		phone, err := user.NewPhone(req.Phone)
		if err != nil {
			return nil, err
		}

		// 检查手机号是否已存在
		exists, err := s.userRepo.ExistsByPhone(ctx, phone)
		if err != nil {
			log.Printf("检查手机号失败: %v", err)
			return nil, errors.New("服务器内部错误")
		}

		if exists {
			return nil, errors.New("手机号已注册")
		}

		// 更新手机号
		existingUser.UpdatePhone(phone)
	}

	// 保存更新后的用户
	if err := s.userRepo.Save(ctx, existingUser); err != nil {
		log.Printf("保存用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 发布用户更新事件
	event := user.NewUserUpdatedEvent(existingUser.ID(), existingUser.Username(), time.Now())
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("发布用户更新事件失败: %v", err)
		// 不影响主流程，继续执行
	}

	return &UpdateUserResponse{
		UserID:            existingUser.ID().Value(),
		Username:          existingUser.Username().Value(),
		Email:             existingUser.Email().Value(),
		Phone:             existingUser.Phone().Value(),
		Status:            int32(existingUser.Status()),
		IsVerified:        existingUser.IsVerified(),
		IsCompanyVerified: existingUser.IsCompanyVerified(),
		DefaultTenantID:   existingUser.DefaultTenantID().Value(),
		Role:              string(existingUser.Role()),
		UpdatedAt:         existingUser.UpdatedAt(),
	}, nil
}

// VerifyUserRequest 用户验证请求
type VerifyUserRequest struct {
	UserID int64
}

// VerifyUserResponse 用户验证响应
type VerifyUserResponse struct {
	UserID     int64
	Username   string
	IsVerified bool
	UpdatedAt  time.Time
}

// VerifyUser 验证用户
func (s *UserApplicationService) VerifyUser(ctx context.Context, req VerifyUserRequest) (*VerifyUserResponse, error) {
	// 创建用户ID值对象
	userID := user.NewUserID(req.UserID)

	// 查找用户
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// 验证用户
	existingUser.Verify()

	// 保存更新后的用户
	if err := s.userRepo.Save(ctx, existingUser); err != nil {
		log.Printf("保存用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 发布用户验证事件
	event := user.NewUserVerifiedEvent(existingUser.ID(), existingUser.Username(), time.Now())
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("发布用户验证事件失败: %v", err)
		// 不影响主流程，继续执行
	}

	return &VerifyUserResponse{
		UserID:     existingUser.ID().Value(),
		Username:   existingUser.Username().Value(),
		IsVerified: existingUser.IsVerified(),
		UpdatedAt:  existingUser.UpdatedAt(),
	}, nil
}

// VerifyCompanyRequest 企业验证请求
type VerifyCompanyRequest struct {
	UserID int64
}

// VerifyCompanyResponse 企业验证响应
type VerifyCompanyResponse struct {
	UserID            int64
	Username          string
	IsCompanyVerified bool
	UpdatedAt         time.Time
}

// VerifyCompany 验证企业
func (s *UserApplicationService) VerifyCompany(ctx context.Context, req VerifyCompanyRequest) (*VerifyCompanyResponse, error) {
	// 创建用户ID值对象
	userID := user.NewUserID(req.UserID)

	// 查找用户
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// 验证企业
	existingUser.VerifyCompany()

	// 保存更新后的用户
	if err := s.userRepo.Save(ctx, existingUser); err != nil {
		log.Printf("保存用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 发布企业验证事件
	event := user.NewCompanyVerifiedEvent(existingUser.ID(), existingUser.Username(), time.Now())
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("发布企业验证事件失败: %v", err)
		// 不影响主流程，继续执行
	}

	return &VerifyCompanyResponse{
		UserID:            existingUser.ID().Value(),
		Username:          existingUser.Username().Value(),
		IsCompanyVerified: existingUser.IsCompanyVerified(),
		UpdatedAt:         existingUser.UpdatedAt(),
	}, nil
}

// 生成token（实际实现中应使用JWT或其他安全的令牌机制）
func generateToken(userID int64, username string) (string, error) {
	// TODO: 实现安全的令牌生成逻辑
	return "sample-token-" + username, nil
}
