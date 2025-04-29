package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetHotContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetHotContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetHotContentLogic {
	return &SetHotContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetHotContentLogic) SetHotContent(in *content.SetHotContentRequest) (*content.SetHotContentResponse, error) {
	// todo: add your logic here and delete this line

	return &content.SetHotContentResponse{}, nil
}
