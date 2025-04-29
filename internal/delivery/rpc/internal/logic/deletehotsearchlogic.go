package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHotSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteHotSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHotSearchLogic {
	return &DeleteHotSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteHotSearchLogic) DeleteHotSearch(in *search.DeleteHotSearchRequest) (*search.DeleteHotSearchResponse, error) {
	// todo: add your logic here and delete this line

	return &search.DeleteHotSearchResponse{}, nil
}
