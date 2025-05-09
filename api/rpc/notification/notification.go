package notification

import (
	"context"

	"google.golang.org/grpc"
)

// 通知服务接口定义
type NotificationService interface {
	// 占位方法
	Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error)
}

// 通知服务RPC客户端
type notificationServiceClient struct {
	conn *grpc.ClientConn
}

// 创建通知服务客户端
func NewNotificationService(conn *grpc.ClientConn) NotificationService {
	return &notificationServiceClient{conn: conn}
}

// 占位请求
type PlaceholderReq struct{}

// 占位响应
type PlaceholderResp struct {
	Success bool `json:"success"`
}

// 实现NotificationService接口的方法
func (c *notificationServiceClient) Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error) {
	return &PlaceholderResp{Success: true}, nil
}
