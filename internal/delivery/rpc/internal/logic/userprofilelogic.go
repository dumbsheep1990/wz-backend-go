package logic

import (
	"context"

	"wz-backend-go/api/rpc/statistics"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户画像
func (l *UserProfileLogic) UserProfile(in *statistics.UserProfileRequest) (*statistics.UserProfileResponse, error) {
	// todo: add your logic here and delete this line

	return &statistics.UserProfileResponse{}, nil
}
