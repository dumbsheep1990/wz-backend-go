package logic

import (
	"context"

	"wz-backend-go/api/rpc/statistics"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBehaviorLogic {
	return &UserBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户行为统计
func (l *UserBehaviorLogic) UserBehavior(in *statistics.UserBehaviorRequest) (*statistics.UserBehaviorResponse, error) {
	// todo: add your logic here and delete this line

	return &statistics.UserBehaviorResponse{}, nil
}
