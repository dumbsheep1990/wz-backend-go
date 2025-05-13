package service

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// PaymentService 支付服务接口
type PaymentService interface {
	// 创建支付
	CreatePayment(ctx context.Context, req model.PaymentRequest) (*model.PaymentResponse, error)

	// 处理支付回调
	HandlePaymentNotify(ctx context.Context, payType int, notifyData map[string]string) error

	// 查询支付状态
	QueryPaymentStatus(ctx context.Context, orderNo string) (int, error)

	// 申请退款
	RefundPayment(ctx context.Context, req model.RefundRequest) (*model.RefundResponse, error)

	// 获取支付配置列表
	GetPaymentConfigs(ctx context.Context) ([]*model.PaymentConfig, error)

	// 保存支付配置
	SavePaymentConfig(ctx context.Context, config *model.PaymentConfig) error

	// 更新支付配置状态
	UpdatePaymentConfigStatus(ctx context.Context, id int64, status bool) error

	// 删除支付配置
	DeletePaymentConfig(ctx context.Context, id int64) error

	// 获取支付记录
	GetPaymentByOrderNo(ctx context.Context, orderNo string) (*model.Payment, error)

	// 获取支付统计
	GetPaymentStatistics(ctx context.Context) (map[string]interface{}, error)

	// 获取指定支付方式配置
	GetPaymentConfigByType(ctx context.Context, payType int) (*model.PaymentConfig, error)
}
