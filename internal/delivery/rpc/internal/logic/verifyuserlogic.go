package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/internal/svc"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyUserLogic {
	return &VerifyUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证用户
func (l *VerifyUserLogic) VerifyUser(in *user.VerifyUserRequest) (*user.VerifyUserResponse, error) {
	// todo: add your logic here and delete this line

	return &user.VerifyUserResponse{}, nil
}
