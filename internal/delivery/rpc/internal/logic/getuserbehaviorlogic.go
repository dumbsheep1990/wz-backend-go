package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/internal/svc"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBehaviorLogic {
	return &GetUserBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户行为分析
func (l *GetUserBehaviorLogic) GetUserBehavior(in *user.GetUserBehaviorRequest) (*user.GetUserBehaviorResponse, error) {
	// todo: add your logic here and delete this line

	return &user.GetUserBehaviorResponse{}, nil
}
