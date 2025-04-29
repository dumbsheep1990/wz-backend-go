package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchContentLogic {
	return &SearchContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 全文搜索
func (l *SearchContentLogic) SearchContent(in *search.SearchContentRequest) (*search.SearchContentResponse, error) {
	// todo: add your logic here and delete this line

	return &search.SearchContentResponse{}, nil
}
