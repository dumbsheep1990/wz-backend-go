package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaymentCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentCallbackLogic {
	return &PaymentCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaymentCallbackLogic) PaymentCallback(in *trade.PaymentCallbackRequest) (*trade.PaymentCallbackResponse, error) {
	// 检查回调请求参数
	if in.OrderId == "" {
		l.Error("Invalid payment callback request: missing order ID")
		return &trade.PaymentCallbackResponse{
			Success: false,
			Error:   "Missing order ID",
		}, nil
	}

	if in.PaymentProvider == "" {
		l.Error("Invalid payment callback request: missing payment provider")
		return &trade.PaymentCallbackResponse{
			Success: false,
			Error:   "Missing payment provider",
		}, nil
	}

	// 构建支付订单命令
	cmd := dto.PayOrderCommand{
		OrderID:       in.OrderId,
		PaymentMethod: in.PaymentProvider,
		TransactionID: in.TransactionId,
		PaidAmount: dto.Money{
			Amount:   in.Amount.Amount,
			Currency: in.Amount.Currency,
		},
		PaymentData: in.PaymentData,
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.PayOrder(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to process payment callback", logx.Field("error", err),
			logx.Field("orderId", in.OrderId),
			logx.Field("provider", in.PaymentProvider))
		return &trade.PaymentCallbackResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 返回响应
	return &trade.PaymentCallbackResponse{
		Success:       true,
		OrderId:       orderDTO.ID,
		OrderNumber:   orderDTO.OrderNumber,
		OrderStatus:   orderDTO.Status,
		OrderStatusCode: orderDTO.StatusCode,
		TransactionId: orderDTO.Transaction,
		Amount: &trade.Money{
			Amount:   orderDTO.TotalAmount.Amount,
			Currency: orderDTO.TotalAmount.Currency,
		},
	}, nil
}
