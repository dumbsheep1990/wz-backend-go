package user

import (
	"context"

	"google.golang.org/grpc"
)

// 用户服务接口定义
type UserService interface {
	// 获取用户列表
	GetUserList(ctx context.Context, in *GetUserListReq) (*GetUserListResp, error)
	// 获取用户详情
	GetUserDetail(ctx context.Context, in *GetUserDetailReq) (*UserDetailResp, error)
	// 创建用户
	CreateUser(ctx context.Context, in *CreateUserReq) (*CreateUserResp, error)
	// 更新用户
	UpdateUser(ctx context.Context, in *UpdateUserReq) (*UpdateUserResp, error)
	// 删除用户
	DeleteUser(ctx context.Context, in *DeleteUserReq) (*DeleteUserResp, error)
}

// 用户服务RPC客户端
type userServiceClient struct {
	conn *grpc.ClientConn
}

// 创建用户服务客户端
func NewUserService(conn *grpc.ClientConn) UserService {
	return &userServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// 获取用户列表请求
type GetUserListReq struct {
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Status    int32  `json:"status,omitempty"`
	Role      string `json:"role,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// 获取用户列表响应
type GetUserListResp struct {
	Total int64             `json:"total"`
	List  []*UserDetailResp `json:"list"`
}

// 获取用户详情请求
type GetUserDetailReq struct {
	Id int64 `json:"id"`
}

// 用户详情响应
type UserDetailResp struct {
	Id                int64  `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Role              string `json:"role"`
	Status            int32  `json:"status"`
	IsVerified        bool   `json:"is_verified"`
	IsCompanyVerified bool   `json:"is_company_verified"`
	DefaultTenantID   int64  `json:"default_tenant_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// 创建用户请求
type CreateUserReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Role            string `json:"role"`
	Status          int32  `json:"status"`
	DefaultTenantID int64  `json:"default_tenant_id,omitempty"`
}

// 创建用户响应
type CreateUserResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

// 更新用户请求
type UpdateUserReq struct {
	Id              int64  `json:"id"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Role            string `json:"role,omitempty"`
	Status          int32  `json:"status,omitempty"`
	DefaultTenantID int64  `json:"default_tenant_id,omitempty"`
}

// 更新用户响应
type UpdateUserResp struct {
	Success bool `json:"success"`
}

// 删除用户请求
type DeleteUserReq struct {
	Id int64 `json:"id"`
}

// 删除用户响应
type DeleteUserResp struct {
	Success bool `json:"success"`
}

// 实现UserService接口的方法
func (c *userServiceClient) GetUserList(ctx context.Context, in *GetUserListReq) (*GetUserListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	// 此处为了演示，返回模拟数据
	return &GetUserListResp{
		Total: 0,
		List:  []*UserDetailResp{},
	}, nil
}

func (c *userServiceClient) GetUserDetail(ctx context.Context, in *GetUserDetailReq) (*UserDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UserDetailResp{}, nil
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserReq) (*CreateUserResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &CreateUserResp{}, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserReq) (*UpdateUserResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateUserResp{Success: true}, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserReq) (*DeleteUserResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteUserResp{Success: true}, nil
}
