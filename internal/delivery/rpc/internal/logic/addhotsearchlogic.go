package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHotSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddHotSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHotSearchLogic {
	return &AddHotSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理热搜词
func (l *AddHotSearchLogic) AddHotSearch(in *search.AddHotSearchRequest) (*search.AddHotSearchResponse, error) {
	// todo: add your logic here and delete this line

	return &search.AddHotSearchResponse{}, nil
}
