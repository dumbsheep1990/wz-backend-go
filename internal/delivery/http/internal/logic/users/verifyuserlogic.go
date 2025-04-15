package users

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyUserLogic {
	return &VerifyUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyUserLogic) VerifyUser(req *types.VerifyUserReq) (resp *types.VerifyUserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
