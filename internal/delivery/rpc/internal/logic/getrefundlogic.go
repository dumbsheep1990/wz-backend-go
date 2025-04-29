package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundLogic {
	return &GetRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRefundLogic) GetRefund(in *trade.GetRefundRequest) (*trade.GetRefundResponse, error) {
	// todo: add your logic here and delete this line

	return &trade.GetRefundResponse{}, nil
}
