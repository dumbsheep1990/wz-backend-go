package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetContentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentStatusLogic {
	return &GetContentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetContentStatusLogic) GetContentStatus(in *content.GetContentStatusRequest) (*content.GetContentStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &content.GetContentStatusResponse{}, nil
}
