package logic

import (
	"context"

	"wz-backend-go/api/rpc/content"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *content.GetPostRequest) (*content.GetPostResponse, error) {
	// todo: add your logic here and delete this line

	return &content.GetPostResponse{}, nil
}
