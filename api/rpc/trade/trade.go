package trade

import (
	"context"

	"google.golang.org/grpc"
)

// 交易服务接口定义
type TradeService interface {
	// 占位方法
	Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error)
}

// 交易服务RPC客户端
type tradeServiceClient struct {
	conn *grpc.ClientConn
}

// 创建交易服务客户端
func NewTradeService(conn *grpc.ClientConn) TradeService {
	return &tradeServiceClient{conn: conn}
}

// 占位请求
type PlaceholderReq struct{}

// 占位响应
type PlaceholderResp struct {
	Success bool `json:"success"`
}

// 实现TradeService接口的方法
func (c *tradeServiceClient) Placeholder(ctx context.Context, in *PlaceholderReq) (*PlaceholderResp, error) {
	return &PlaceholderResp{Success: true}, nil
}
