package logic

import (
	"context"

	"wz-backend-go/api/rpc/interaction"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeLogic {
	return &UnlikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消点赞
func (l *UnlikeLogic) Unlike(in *interaction.LikeRequest) (*interaction.LikeResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.LikeResponse{}, nil
}
