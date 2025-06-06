package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/trade"
)

// TradeClient 是交易服务的客户端封装
type TradeClient struct {
	client pb.TradeServiceClient
	conn   *grpc.ClientConn
}

// NewTradeClient 创建一个新的交易服务客户端
func NewTradeClient(serviceAddr string) (*TradeClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50053" // 默认地址
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
		return nil, fmt.Errorf("无法连接到交易服务: %v", err)
	}

	client := pb.NewTradeServiceClient(conn)
	return &TradeClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *TradeClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetTradeStats 获取交易统计数据
func (c *TradeClient) GetTradeStats(ctx context.Context) (*pb.TradeStatsResponse, error) {
	return c.client.GetTradeStats(ctx, &pb.GetTradeStatsRequest{})
}

// ListOrders 获取订单列表
func (c *TradeClient) ListOrders(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListOrdersResponse, error) {
	return c.client.ListOrders(ctx, &pb.ListOrdersRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetOrder 获取订单详情
func (c *TradeClient) GetOrder(ctx context.Context, id string) (*pb.Order, error) {
	return c.client.GetOrder(ctx, &pb.GetOrderRequest{
		Id: id,
	})
}

// CreateOrder 创建订单
func (c *TradeClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return c.client.CreateOrder(ctx, req)
}

// UpdateOrder 更新订单
func (c *TradeClient) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.Order, error) {
	return c.client.UpdateOrder(ctx, req)
}

// CancelOrder 取消订单
func (c *TradeClient) CancelOrder(ctx context.Context, orderId string, reason string) (*pb.CancelOrderResponse, error) {
	return c.client.CancelOrder(ctx, &pb.CancelOrderRequest{
		OrderId: orderId,
		Reason:  reason,
	})
}

// PayOrder 支付订单
func (c *TradeClient) PayOrder(ctx context.Context, orderId string, paymentMethod string, paymentId string) (*pb.PayOrderResponse, error) {
	return c.client.PayOrder(ctx, &pb.PayOrderRequest{
		OrderId:       orderId,
		PaymentMethod: paymentMethod,
		PaymentId:     paymentId,
	})
}

// RefundOrder 退款订单
func (c *TradeClient) RefundOrder(ctx context.Context, orderId string, reason string) (*pb.RefundOrderResponse, error) {
	return c.client.RefundOrder(ctx, &pb.RefundOrderRequest{
		OrderId: orderId,
		Reason:  reason,
	})
}

// ListPayments 获取支付记录列表
func (c *TradeClient) ListPayments(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListPaymentsResponse, error) {
	return c.client.ListPayments(ctx, &pb.ListPaymentsRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetPayment 获取支付记录详情
func (c *TradeClient) GetPayment(ctx context.Context, id string) (*pb.Payment, error) {
	return c.client.GetPayment(ctx, &pb.GetPaymentRequest{
		Id: id,
	})
}

// ListRefunds 获取退款记录列表
func (c *TradeClient) ListRefunds(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListRefundsResponse, error) {
	return c.client.ListRefunds(ctx, &pb.ListRefundsRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetRefund 获取退款记录详情
func (c *TradeClient) GetRefund(ctx context.Context, id string) (*pb.Refund, error) {
	return c.client.GetRefund(ctx, &pb.GetRefundRequest{
		Id: id,
	})
}
