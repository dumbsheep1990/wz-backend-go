package server

import (
	"context"

	pb "wz-backend-go/services/user-service/api/proto"
	"wz-backend-go/services/user-service/internal/service"
)

// RegisterUserServer 注册用户服务
func RegisterUserServer(server interface{}, srv interface{}) {
	pb.RegisterUserServer(server.(pb.UserServer), srv.(*service.UserService))
}

// UserServer 用户服务实现
type UserServer struct {
	pb.UnimplementedUserServer
	userService *service.UserService
}

// NewUserServer 创建用户服务服务器
func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

// Register 用户注册
func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.userService.Register(ctx, req)
}

// Login 用户登录
func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.userService.Login(ctx, req)
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.userService.GetUser(ctx, req)
}

// UpdateUser 更新用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.userService.UpdateUser(ctx, req)
}

// VerifyUser 验证用户
func (s *UserServer) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	return s.userService.VerifyUser(ctx, req)
}

// VerifyCompany 企业认证
func (s *UserServer) VerifyCompany(ctx context.Context, req *pb.VerifyCompanyRequest) (*pb.VerifyCompanyResponse, error) {
	return s.userService.VerifyCompany(ctx, req)
}

// GetUserBehavior 获取用户行为
func (s *UserServer) GetUserBehavior(ctx context.Context, req *pb.GetUserBehaviorRequest) (*pb.GetUserBehaviorResponse, error) {
	return s.userService.GetUserBehavior(ctx, req)
}
