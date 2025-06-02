package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRefundsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRefundsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRefundsLogic {
	return &ListRefundsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRefundsLogic) ListRefunds(in *trade.ListRefundsRequest) (*trade.ListRefundsResponse, error) {
	// 构建退款列表查询
	query := dto.ListRefundsQuery{
		CustomerID:  in.CustomerId,
		OrderID:     in.OrderId,
		OrderNumber: in.OrderNumber,
		Status:      in.Status,
		Page:        int(in.Page),
		PageSize:    int(in.PageSize),
		SortBy:      in.SortBy,
		SortOrder:   in.SortOrder,
	}

	// 如果没有指定页面大小，使用默认值
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 如果没有指定页码，从第一页开始
	if query.Page <= 0 {
		query.Page = 1
	}

	// 调用应用服务查询退款列表
	result, err := l.svcCtx.OrderApplicationService.ListRefunds(l.ctx, query)
	if err != nil {
		l.Error("Failed to list refunds", logx.Field("error", err),
			logx.Field("customerId", in.CustomerId),
			logx.Field("orderId", in.OrderId))
		return &trade.ListRefundsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 构建响应对象
	response := &trade.ListRefundsResponse{
		Success:     true,
		Total:       int64(result.Total),
		Page:        int64(result.Page),
		PageSize:    int64(result.PageSize),
		TotalPages:  int64(result.TotalPages),
		HasMore:     result.HasMore,
		Refunds:     make([]*trade.Refund, 0, len(result.Refunds)),
	}

	// 转换退款数据
	for _, refund := range result.Refunds {
		// 构建退款项目
		refundItems := make([]*trade.RefundItem, 0, len(refund.Items))
		for _, item := range refund.Items {
			refundItems = append(refundItems, &trade.RefundItem{
				OrderItemId: item.OrderItemID,
				Quantity:    int32(item.Quantity),
				Reason:      item.Reason,
			})
		}

		// 构建退款对象
		refundObj := &trade.Refund{
			Id:           refund.ID,
			OrderId:      refund.OrderID,
			OrderNumber:  refund.OrderNumber,
			CustomerId:   refund.CustomerID,
			Status:       refund.Status,
			Reason:       refund.Reason,
			Amount: &trade.Money{
				Amount:   refund.Amount.Amount,
				Currency: refund.Amount.Currency,
			},
			Items:        refundItems,
			CustomerNote: refund.CustomerNote,
			AdminNote:    refund.AdminNote,
			RefundMethod: refund.RefundMethod,
			Transaction:  refund.Transaction,
			CreatedAt:    refund.CreatedAt.Unix(),
		}

		// 如果有处理时间
		if !refund.ProcessedAt.IsZero() {
			refundObj.ProcessedAt = refund.ProcessedAt.Unix()
		}

		response.Refunds = append(response.Refunds, refundObj)
	}

	return response, nil
}
