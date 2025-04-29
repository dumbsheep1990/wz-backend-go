package logic

import (
	"context"

	"wz-backend-go/api/rpc/search"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogSearchQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogSearchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogSearchQueryLogic {
	return &LogSearchQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索日志管理
func (l *LogSearchQueryLogic) LogSearchQuery(in *search.LogSearchQueryRequest) (*search.LogSearchQueryResponse, error) {
	// todo: add your logic here and delete this line

	return &search.LogSearchQueryResponse{}, nil
}
