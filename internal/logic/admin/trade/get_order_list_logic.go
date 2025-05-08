package trade

import (
	"context"
	"fmt"
	"strconv"

	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *http.OrderListReq) (*http.OrderListResp, error) {
	// 构建查询过滤条件
	filters := make(map[string]interface{})
	if req.OrderId != "" {
		filters["orderId"] = req.OrderId
	}
	if req.UserId > 0 {
		filters["userId"] = req.UserId
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.ProductType != "" {
		filters["productType"] = req.ProductType
	}
	if req.StartTime != "" {
		filters["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		filters["endTime"] = req.EndTime
	}
	if req.TenantId > 0 {
		filters["tenantId"] = req.TenantId
	}

	// 调用仓库层查询数据
	orders, total, err := l.svcCtx.TradeRepo.GetOrderList(l.ctx, req.Page, req.PageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("查询订单列表失败: %v", err)
	}

	// 转换数据格式
	var list []http.OrderDetail
	for _, order := range orders {
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

		list = append(list, http.OrderDetail{
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
		})
	}

	return &http.OrderListResp{
		Total: total,
		List:  list,
	}, nil
} 