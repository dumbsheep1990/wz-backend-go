package service

import (
	"context"
	"errors"
	"log"
	"time"

	pb "wz-backend-go/services/user-service/api/proto"
	"wz-backend-go/services/user-service/internal/model"
	"wz-backend-go/services/user-service/internal/repository"
)

// UserService 实现User服务接口
type UserService struct {
	pb.UnimplementedUserServer
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("检查用户名失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("检查邮箱失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	if exists {
		return nil, errors.New("邮箱已注册")
	}

	// 创建用户
	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    1, // 活跃状态
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置密码（实际中应该使用安全的密码哈希）
	if err := user.SetPassword(req.Password); err != nil {
		log.Printf("设置密码失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 保存用户
	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("创建用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 生成token（实际中应使用JWT或其他令牌机制）
	token, err := generateToken(userID, req.Username)
	if err != nil {
		log.Printf("生成令牌失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &pb.RegisterResponse{
		UserId:   userID,
		Username: req.Username,
		Token:    token,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成令牌失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &pb.LoginResponse{
		UserId:   user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 根据ID查找用户
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	return &pb.GetUserResponse{
		UserId:            user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		Status:            user.Status,
		IsVerified:        user.IsVerified,
		IsCompanyVerified: user.IsCompanyVerified,
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 根据ID查找用户
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			log.Printf("设置密码失败: %v", err)
			return nil, errors.New("服务器内部错误")
		}
	}

	user.UpdatedAt = time.Now()

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("更新用户失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &pb.UpdateUserResponse{
		Success: true,
	}, nil
}

// VerifyUser 验证用户
func (s *UserService) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	// 根据ID查找用户
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// TODO: 实现验证码校验逻辑

	// 更新用户验证状态
	user.IsVerified = true
	user.UpdatedAt = time.Now()

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("更新用户验证状态失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &pb.VerifyUserResponse{
		Success: true,
	}, nil
}

// VerifyCompany 验证企业
func (s *UserService) VerifyCompany(ctx context.Context, req *pb.VerifyCompanyRequest) (*pb.VerifyCompanyResponse, error) {
	// 根据ID查找用户
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		log.Printf("查找用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// TODO: 实现企业认证逻辑

	// 更新用户企业验证状态
	user.IsCompanyVerified = true
	user.UpdatedAt = time.Now()

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("更新用户企业验证状态失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	return &pb.VerifyCompanyResponse{
		Success: true,
	}, nil
}

// GetUserBehavior 获取用户行为
func (s *UserService) GetUserBehavior(ctx context.Context, req *pb.GetUserBehaviorRequest) (*pb.GetUserBehaviorResponse, error) {
	// 解析时间范围
	var startTime, endTime time.Time
	var err error

	if req.StartTime != "" {
		startTime, err = time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
		}
	}

	if req.EndTime != "" {
		endTime, err = time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return nil, errors.New("结束时间格式错误")
		}
	}

	// 查询用户行为
	behaviors, err := s.userRepo.GetUserBehaviors(ctx, req.UserId, startTime, endTime)
	if err != nil {
		log.Printf("获取用户行为失败: %v", err)
		return nil, errors.New("服务器内部错误")
	}

	// 转换为响应格式
	var respBehaviors []*pb.UserBehavior
	for _, b := range behaviors {
		respBehaviors = append(respBehaviors, &pb.UserBehavior{
			BehaviorId:   b.ID,
			UserId:       b.UserID,
			Action:       b.Action,
			ResourceType: b.ResourceType,
			ResourceId:   b.ResourceID,
			CreatedAt:    b.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetUserBehaviorResponse{
		Behaviors: respBehaviors,
	}, nil
}

// generateToken 生成认证令牌
func generateToken(userID int64, username string) (string, error) {
	// TODO: 实现JWT令牌生成
	// 此处为简化实现
	return "sample_token_" + username, nil
}
