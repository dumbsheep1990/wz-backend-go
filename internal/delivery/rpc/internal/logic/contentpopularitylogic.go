package logic

import (
	"context"

	"wz-backend-go/api/rpc/statistics"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentPopularityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewContentPopularityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentPopularityLogic {
	return &ContentPopularityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取内容流行度统计
func (l *ContentPopularityLogic) ContentPopularity(in *statistics.ContentPopularityRequest) (*statistics.ContentPopularityResponse, error) {
	// todo: add your logic here and delete this line

	return &statistics.ContentPopularityResponse{}, nil
}
