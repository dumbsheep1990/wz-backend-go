package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessRefundLogic {
	return &ProcessRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessRefundLogic) ProcessRefund(in *trade.ProcessRefundRequest) (*trade.ProcessRefundResponse, error) {
	// todo: add your logic here and delete this line

	return &trade.ProcessRefundResponse{}, nil
}
