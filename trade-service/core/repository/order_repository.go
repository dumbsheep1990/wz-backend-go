package repository

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// 创建订单
	Create(ctx context.Context, order *model.Order) error

	// 获取订单详情
	GetByID(ctx context.Context, id int64) (*model.Order, error)

	// 根据订单号获取订单
	GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)

	// 根据用户ID获取订单列表
	GetByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*model.Order, int64, error)

	// 更新订单状态
	UpdateStatus(ctx context.Context, id int64, status int) error

	// 更新订单支付信息
	UpdatePayment(ctx context.Context, id int64, payType int, payTime string) error

	// 更新订单物流信息
	UpdateShipping(ctx context.Context, id int64, logisticsCompany, logisticsNo string) error

	// 删除订单
	Delete(ctx context.Context, id int64) error

	// 获取订单列表
	ListOrders(ctx context.Context, query model.OrderQuery) ([]*model.Order, int64, error)

	// 创建订单项
	CreateOrderItems(ctx context.Context, items []*model.OrderItem) error

	// 获取订单项
	GetOrderItems(ctx context.Context, orderID int64) ([]*model.OrderItem, error)

	// 获取订单统计信息
	GetOrderStatistics(ctx context.Context) (*model.OrderStatistics, error)

	// 获取状态统计
	GetStatusCounts(ctx context.Context) ([]*model.StatusCount, error)

	// 获取支付方式统计
	GetPaymentTypeCounts(ctx context.Context) ([]*model.PaymentTypeCount, error)

	// 获取订单趋势
	GetOrderTrend(ctx context.Context, period string) ([]*model.TrendData, error)
}
