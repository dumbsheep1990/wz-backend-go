package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessPaymentLogic {
	return &ProcessPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付管理
func (l *ProcessPaymentLogic) ProcessPayment(in *trade.ProcessPaymentRequest) (*trade.ProcessPaymentResponse, error) {
	// todo: add your logic here and delete this line

	return &trade.ProcessPaymentResponse{}, nil
}
