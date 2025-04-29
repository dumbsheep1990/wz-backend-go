package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotContentLogic {
	return &GetHotContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 热门内容管理
func (l *GetHotContentLogic) GetHotContent(in *content.GetHotContentRequest) (*content.GetHotContentResponse, error) {
	// todo: add your logic here and delete this line

	return &content.GetHotContentResponse{}, nil
}
