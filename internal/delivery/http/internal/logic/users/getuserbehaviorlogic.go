package users

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBehaviorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBehaviorLogic {
	return &GetUserBehaviorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBehaviorLogic) GetUserBehavior(req *types.UserBehaviorReq) (resp *types.UserBehaviorResp, err error) {
	// todo: add your logic here and delete this line

	return
}
