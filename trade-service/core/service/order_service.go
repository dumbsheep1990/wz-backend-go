package service

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// OrderService 订单服务接口
type OrderService interface {
	// 创建订单
	CreateOrder(ctx context.Context, order *model.Order) error

	// 获取订单详情
	GetOrderDetail(ctx context.Context, id int64) (*model.Order, error)

	// 根据订单号获取订单
	GetOrderByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)

	// 获取用户订单列表
	GetUserOrders(ctx context.Context, userID int64, page, pageSize int) ([]*model.Order, int64, error)

	// 取消订单
	CancelOrder(ctx context.Context, id int64) error

	// 订单支付成功处理
	PaymentSuccess(ctx context.Context, orderNo string, payType int, tradeNo string) error

	// 订单发货
	ShipOrder(ctx context.Context, id int64, logisticsCompany, logisticsNo string) error

	// 确认收货
	ConfirmReceipt(ctx context.Context, id int64) error

	// 删除订单
	DeleteOrder(ctx context.Context, id int64) error

	// 退款
	RefundOrder(ctx context.Context, id int64, amount float64, reason string) error

	// 获取订单列表
	ListOrders(ctx context.Context, query model.OrderQuery) ([]*model.Order, int64, error)

	// 导出订单数据
	ExportOrders(ctx context.Context, query model.OrderQuery) ([]byte, error)

	// 获取订单统计信息
	GetOrderStatistics(ctx context.Context) (*model.OrderStatistics, error)

	// 获取订单状态统计
	GetStatusStatistics(ctx context.Context) ([]*model.StatusCount, error)

	// 获取支付方式统计
	GetPaymentTypeStatistics(ctx context.Context) ([]*model.PaymentTypeCount, error)

	// 获取订单趋势
	GetOrderTrend(ctx context.Context, period string) ([]*model.TrendData, error)
}
