package statistics

import (
	"context"

	"google.golang.org/grpc"
)

// 统计服务接口定义
type StatisticsService interface {
	// 占位方法
	Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error)
}

// 统计服务RPC客户端
type statisticsServiceClient struct {
	conn *grpc.ClientConn
}

// 创建统计服务客户端
func NewStatisticsService(conn *grpc.ClientConn) StatisticsService {
	return &statisticsServiceClient{conn: conn}
}

// 占位请求
type PlaceholderReq struct{}

// 占位响应
type PlaceholderResp struct {
	Success bool `json:"success"`
}

// 实现StatisticsService接口的方法
func (c *statisticsServiceClient) Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error) {
	return &PlaceholderResp{Success: true}, nil
}
