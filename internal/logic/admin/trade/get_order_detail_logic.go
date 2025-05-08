package trade

import (
	"context"
	"fmt"

	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(id string) (*http.OrderDetail, error) {
	// 调用仓库层查询数据
	order, err := l.svcCtx.TradeRepo.GetOrderById(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询订单详情失败: %v", err)
	}

	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 转换数据格式
	var paymentTime, createdAt, updatedAt, expireTime string
	
	if !order.PaymentTime.IsZero() {
		paymentTime = order.PaymentTime.Format("2006-01-02 15:04:05")
	}
	if !order.CreatedAt.IsZero() {
		createdAt = order.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !order.UpdatedAt.IsZero() {
		updatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	if !order.ExpireTime.IsZero() {
		expireTime = order.ExpireTime.Format("2006-01-02 15:04:05")
	}

	result := &http.OrderDetail{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		TenantId:    order.TenantId,
		ProductId:   order.ProductId,
		ProductType: order.ProductType,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		Amount:      order.Amount,
		Currency:    order.Currency,
		Status:      order.Status,
		PaymentId:   order.PaymentId,
		PaymentType: order.PaymentType,
		PaymentTime: paymentTime,
		Description: order.Description,
		Metadata:    order.Metadata,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		ExpireTime:  expireTime,
	}

	return result, nil
} 