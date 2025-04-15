package logic

import (
	"context"

	"wz-backend-go/internal/delivery/rpc/internal/svc"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCompanyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCompanyLogic {
	return &VerifyCompanyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 企业认证
func (l *VerifyCompanyLogic) VerifyCompany(in *user.VerifyCompanyRequest) (*user.VerifyCompanyResponse, error) {
	// todo: add your logic here and delete this line

	return &user.VerifyCompanyResponse{}, nil
}
