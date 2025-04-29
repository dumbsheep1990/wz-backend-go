package logic

import (
	"context"

	"wz-backend-go/api/rpc/statistics"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HotContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHotContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HotContentLogic {
	return &HotContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取热门内容
func (l *HotContentLogic) HotContent(in *statistics.HotContentRequest) (*statistics.HotContentResponse, error) {
	// todo: add your logic here and delete this line

	return &statistics.HotContentResponse{}, nil
}
