package users

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCompanyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCompanyLogic {
	return &VerifyCompanyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCompanyLogic) VerifyCompany(req *types.VerifyCompanyReq) (resp *types.VerifyCompanyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
