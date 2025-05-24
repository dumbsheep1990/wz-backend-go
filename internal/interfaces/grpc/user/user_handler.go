package user

import (
	"context"

	pb "wz-backend-go/services/user-service/api/proto"
	appuser "wz-backend-go/internal/application/user"
)

// UserHandler 实现用户gRPC服务接口
type UserHandler struct {
	pb.UnimplementedUserServer
	userAppService *appuser.UserApplicationService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userAppService *appuser.UserApplicationService) *UserHandler {
	return &UserHandler{
		userAppService: userAppService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 转换请求
	appReq := appuser.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 调用应用服务
	resp, err := h.userAppService.Register(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.RegisterResponse{
		UserId:   resp.UserID,
		Username: resp.Username,
		Token:    resp.Token,
	}, nil
}

// Login 用户登录
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 转换请求
	appReq := appuser.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// 调用应用服务
	resp, err := h.userAppService.Login(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.LoginResponse{
		UserId:   resp.UserID,
		Username: resp.Username,
		Token:    resp.Token,
	}, nil
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 转换请求
	appReq := appuser.GetUserRequest{
		UserID: req.UserId,
	}

	// 调用应用服务
	resp, err := h.userAppService.GetUser(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.GetUserResponse{
		UserId:            resp.UserID,
		Username:          resp.Username,
		Email:             resp.Email,
		Phone:             resp.Phone,
		Status:            resp.Status,
		IsVerified:        resp.IsVerified,
		IsCompanyVerified: resp.IsCompanyVerified,
		DefaultTenantId:   resp.DefaultTenantID,
		Role:              resp.Role,
		CreatedAt:         resp.CreatedAt.Unix(),
		UpdatedAt:         resp.UpdatedAt.Unix(),
	}, nil
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 转换请求
	appReq := appuser.UpdateUserRequest{
		UserID:   req.UserId,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 调用应用服务
	resp, err := h.userAppService.UpdateUser(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.UpdateUserResponse{
		UserId:            resp.UserID,
		Username:          resp.Username,
		Email:             resp.Email,
		Phone:             resp.Phone,
		Status:            resp.Status,
		IsVerified:        resp.IsVerified,
		IsCompanyVerified: resp.IsCompanyVerified,
		DefaultTenantId:   resp.DefaultTenantID,
		Role:              resp.Role,
		UpdatedAt:         resp.UpdatedAt.Unix(),
	}, nil
}

// VerifyUser 验证用户
func (h *UserHandler) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	// 转换请求
	appReq := appuser.VerifyUserRequest{
		UserID: req.UserId,
	}

	// 调用应用服务
	resp, err := h.userAppService.VerifyUser(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.VerifyUserResponse{
		UserId:     resp.UserID,
		Username:   resp.Username,
		IsVerified: resp.IsVerified,
		UpdatedAt:  resp.UpdatedAt.Unix(),
	}, nil
}

// VerifyCompany 验证企业
func (h *UserHandler) VerifyCompany(ctx context.Context, req *pb.VerifyCompanyRequest) (*pb.VerifyCompanyResponse, error) {
	// 转换请求
	appReq := appuser.VerifyCompanyRequest{
		UserID: req.UserId,
	}

	// 调用应用服务
	resp, err := h.userAppService.VerifyCompany(ctx, appReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &pb.VerifyCompanyResponse{
		UserId:            resp.UserID,
		Username:          resp.Username,
		IsCompanyVerified: resp.IsCompanyVerified,
		UpdatedAt:         resp.UpdatedAt.Unix(),
	}, nil
}

// GetUserBehavior 获取用户行为
func (h *UserHandler) GetUserBehavior(ctx context.Context, req *pb.GetUserBehaviorRequest) (*pb.GetUserBehaviorResponse, error) {
	// TODO: 实现用户行为相关功能
	// 这部分功能可能涉及到其他领域，例如交互/行为领域，可能需要单独的DDD设计
	return &pb.GetUserBehaviorResponse{
		UserId: req.UserId,
		// 其他字段暂时留空
	}, nil
}
