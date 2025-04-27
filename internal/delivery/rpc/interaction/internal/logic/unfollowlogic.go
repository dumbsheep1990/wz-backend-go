package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/interaction/interaction"
	"wz-backend-go/internal/delivery/rpc/interaction/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消关注
func (l *UnfollowLogic) Unfollow(in *interaction.FollowRequest) (*interaction.FollowResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.FollowResponse{}, nil
}
