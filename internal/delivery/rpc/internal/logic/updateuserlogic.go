package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/internal/svc"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	// todo: add your logic here and delete this line

	return &user.UpdateUserResponse{}, nil
}
