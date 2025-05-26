package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyEnterpriseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyEnterpriseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEnterpriseLogic {
	return &VerifyEnterpriseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyEnterpriseLogic) VerifyEnterprise(req *types.VerifyEnterpriseReq) (resp *types.VerifyEnterpriseResp, err error) {
	// Call domain service
	err = l.svcCtx.EnterpriseRegistrationService.VerifyEnterprise(l.ctx, req.UserID, req.VerificationCode)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	return &types.VerifyEnterpriseResp{
		Success: true,
	}, nil
}
