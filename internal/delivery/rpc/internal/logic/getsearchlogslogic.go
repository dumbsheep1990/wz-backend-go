package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSearchLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSearchLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSearchLogsLogic {
	return &GetSearchLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSearchLogsLogic) GetSearchLogs(in *search.GetSearchLogsRequest) (*search.GetSearchLogsResponse, error) {
	// todo: add your logic here and delete this line

	return &search.GetSearchLogsResponse{}, nil
}
