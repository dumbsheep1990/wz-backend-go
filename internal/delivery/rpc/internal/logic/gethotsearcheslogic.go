package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotSearchesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotSearchesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotSearchesLogic {
	return &GetHotSearchesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取热搜词
func (l *GetHotSearchesLogic) GetHotSearches(in *search.GetHotSearchesRequest) (*search.GetHotSearchesResponse, error) {
	// todo: add your logic here and delete this line

	return &search.GetHotSearchesResponse{}, nil
}
