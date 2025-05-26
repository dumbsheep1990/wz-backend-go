package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"
	"github.com/wz-project/wz-backend-go/internal/domain/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEnterpriseRegistrationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEnterpriseRegistrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEnterpriseRegistrationLogic {
	return &UpdateEnterpriseRegistrationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEnterpriseRegistrationLogic) UpdateEnterpriseRegistration(req *types.UpdateEnterpriseRegistrationReq) (resp *types.UpdateEnterpriseRegistrationResp, err error) {
	// Get current registration
	registration, err := l.svcCtx.EnterpriseRegistrationService.GetEnterpriseRegistration(l.ctx, req.UserID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	if registration == nil {
		return nil, types.NewInternalError("企业入驻信息不存在")
	}

	// Update registration fields
	registration.CompanyName = req.CompanyName
	registration.CompanyType = model.CompanyType(req.CompanyType)
	registration.ContactPerson = req.ContactPerson
	registration.Region = req.Region
	registration.DetailedAddress = req.DetailedAddress
	registration.LocationLatitude = req.LocationLatitude
	registration.LocationLongitude = req.LocationLongitude

	// Call domain service
	err = l.svcCtx.EnterpriseRegistrationService.UpdateEnterpriseRegistration(l.ctx, registration)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	return &types.UpdateEnterpriseRegistrationResp{
		Success: true,
	}, nil
}
