package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
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
	// todo: add your logic here and delete this line

	return &trade.PaymentCallbackResponse{}, nil
}
