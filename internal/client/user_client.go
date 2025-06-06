package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/user"
)

// UserClient 是用户服务的客户端封装
type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient 创建一个新的用户服务客户端
func NewUserClient(serviceAddr string) (*UserClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50051" // 默认地址
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("无法连接到用户服务: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUserStats 获取用户统计数据
func (c *UserClient) GetUserStats(ctx context.Context) (*pb.UserStatsResponse, error) {
	return c.client.GetUserStats(ctx, &pb.GetUserStatsRequest{})
}

// ListUsers 获取用户列表
func (c *UserClient) ListUsers(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, &pb.ListUsersRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetUser 获取用户详情
func (c *UserClient) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	return c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: id,
	})
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	return c.client.CreateUser(ctx, req)
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	return c.client.UpdateUser(ctx, req)
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(ctx context.Context, id int64) (*pb.DeleteUserResponse, error) {
	return c.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id,
	})
}

// GetUserByUsername 通过用户名获取用户
func (c *UserClient) GetUserByUsername(ctx context.Context, username string) (*pb.User, error) {
	return c.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})
}

// ListUserRoles 获取用户角色列表
func (c *UserClient) ListUserRoles(ctx context.Context, userId int64) (*pb.ListUserRolesResponse, error) {
	return c.client.ListUserRoles(ctx, &pb.ListUserRolesRequest{
		UserId: userId,
	})
}

// AssignRoleToUser 为用户分配角色
func (c *UserClient) AssignRoleToUser(ctx context.Context, userId int64, roleId string) (*pb.AssignRoleToUserResponse, error) {
	return c.client.AssignRoleToUser(ctx, &pb.AssignRoleToUserRequest{
		UserId: userId,
		RoleId: roleId,
	})
}

// RemoveRoleFromUser 从用户移除角色
func (c *UserClient) RemoveRoleFromUser(ctx context.Context, userId int64, roleId string) (*pb.RemoveRoleFromUserResponse, error) {
	return c.client.RemoveRoleFromUser(ctx, &pb.RemoveRoleFromUserRequest{
		UserId: userId,
		RoleId: roleId,
	})
}

// VerifyUserPassword 验证用户密码
func (c *UserClient) VerifyUserPassword(ctx context.Context, username, password string) (*pb.VerifyUserPasswordResponse, error) {
	return c.client.VerifyUserPassword(ctx, &pb.VerifyUserPasswordRequest{
		Username: username,
		Password: password,
	})
}

// ChangeUserPassword 修改用户密码
func (c *UserClient) ChangeUserPassword(ctx context.Context, userId int64, oldPassword, newPassword string) (*pb.ChangeUserPasswordResponse, error) {
	return c.client.ChangeUserPassword(ctx, &pb.ChangeUserPasswordRequest{
		UserId:      userId,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
}

// ResetUserPassword 重置用户密码
func (c *UserClient) ResetUserPassword(ctx context.Context, userId int64, newPassword string) (*pb.ResetUserPasswordResponse, error) {
	return c.client.ResetUserPassword(ctx, &pb.ResetUserPasswordRequest{
		UserId:      userId,
		NewPassword: newPassword,
	})
}
