package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRefundLogic {
	return &CreateRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退款管理
func (l *CreateRefundLogic) CreateRefund(in *trade.CreateRefundRequest) (*trade.CreateRefundResponse, error) {
	// 构建退款请求命令
	cmd := dto.RequestRefundCommand{
		OrderID:      in.OrderId,
		Reason:       in.Reason,
		RefundAmount: dto.Money{Amount: in.Amount.Amount, Currency: in.Amount.Currency},
		CustomerNote: in.CustomerNote,
		RefundItems:  make([]dto.RefundItemDTO, 0, len(in.Items)),
	}

	// 处理退款项目
	for _, item := range in.Items {
		cmd.RefundItems = append(cmd.RefundItems, dto.RefundItemDTO{
			OrderItemID: item.OrderItemId,
			Quantity:    int(item.Quantity),
			Reason:      item.Reason,
		})
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.RequestRefund(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to create refund", logx.Field("error", err),
			logx.Field("orderId", in.OrderId))
		return &trade.CreateRefundResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 构建返回的退款对象
	refundItems := make([]*trade.RefundItem, 0, len(cmd.RefundItems))
	for _, item := range cmd.RefundItems {
		refundItems = append(refundItems, &trade.RefundItem{
			OrderItemId: item.OrderItemID,
			Quantity:    int32(item.Quantity),
			Reason:      item.Reason,
		})
	}

	// 构建响应
	return &trade.CreateRefundResponse{
		Success: true,
		Refund: &trade.Refund{
			Id:           orderDTO.ID + "-refund", // 临时生成退款编号
			OrderId:      orderDTO.ID,
			OrderNumber:  orderDTO.OrderNumber,
			CustomerId:   orderDTO.CustomerID,
			Status:       "requested",
			Reason:       cmd.Reason,
			Amount: &trade.Money{
				Amount:   cmd.RefundAmount.Amount,
				Currency: cmd.RefundAmount.Currency,
			},
			Items:        refundItems,
			CustomerNote: cmd.CustomerNote,
			CreatedAt:    orderDTO.RefundRequestedAt.Unix(),
		},
	}, nil
}
