package repository

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	// 创建支付记录
	Create(ctx context.Context, payment *model.Payment) error

	// 获取支付记录
	GetByID(ctx context.Context, id int64) (*model.Payment, error)

	// 根据订单号获取支付记录
	GetByOrderNo(ctx context.Context, orderNo string) (*model.Payment, error)

	// 更新支付状态
	UpdateStatus(ctx context.Context, id int64, status int, tradeNo string) error

	// 获取支付配置列表
	GetPaymentConfigs(ctx context.Context) ([]*model.PaymentConfig, error)

	// 根据支付类型获取支付配置
	GetPaymentConfigByType(ctx context.Context, payType int) (*model.PaymentConfig, error)

	// 保存支付配置
	SavePaymentConfig(ctx context.Context, config *model.PaymentConfig) error

	// 更新支付配置状态
	UpdatePaymentConfigStatus(ctx context.Context, id int64, status bool) error

	// 删除支付配置
	DeletePaymentConfig(ctx context.Context, id int64) error

	// 获取用户的支付统计数据
	GetUserPaymentStats(ctx context.Context, userID int64) (map[string]interface{}, error)
}
