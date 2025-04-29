package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionsLogic {
	return &GetTransactionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTransactionsLogic) GetTransactions(in *trade.GetTransactionsRequest) (*trade.GetTransactionsResponse, error) {
	// todo: add your logic here and delete this line

	return &trade.GetTransactionsResponse{}, nil
}
