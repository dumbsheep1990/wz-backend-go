package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/interaction/interaction"
	"wz-backend-go/internal/delivery/rpc/interaction/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 关注
func (l *FollowLogic) Follow(in *interaction.FollowRequest) (*interaction.FollowResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.FollowResponse{}, nil
}
