package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSuggestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSuggestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSuggestionsLogic {
	return &GetSuggestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索推荐
func (l *GetSuggestionsLogic) GetSuggestions(in *search.GetSuggestionsRequest) (*search.GetSuggestionsResponse, error) {
	// todo: add your logic here and delete this line

	return &search.GetSuggestionsResponse{}, nil
}
