package file

import (
	"context"

	"google.golang.org/grpc"
)

// 文件服务接口定义
type FileService interface {
	// 占位方法
	Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error)
}

// 文件服务RPC客户端
type fileServiceClient struct {
	conn *grpc.ClientConn
}

// 创建文件服务客户端
func NewFileService(conn *grpc.ClientConn) FileService {
	return &fileServiceClient{conn: conn}
}

// 占位请求
type PlaceholderReq struct{}

// 占位响应
type PlaceholderResp struct {
	Success bool `json:"success"`
}

// 实现FileService接口的方法
func (c *fileServiceClient) Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error) {
	return &PlaceholderResp{Success: true}, nil
}
