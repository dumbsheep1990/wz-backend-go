package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
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
	// todo: add your logic here and delete this line

	return &trade.ListRefundsResponse{}, nil
}
