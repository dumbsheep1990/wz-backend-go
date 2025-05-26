package users

import (
	"context"
	"time"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEnterpriseRegistrationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEnterpriseRegistrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEnterpriseRegistrationLogic {
	return &GetEnterpriseRegistrationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEnterpriseRegistrationLogic) GetEnterpriseRegistration(userID int64) (resp *types.GetEnterpriseRegistrationResp, err error) {
	// Call domain service
	registration, err := l.svcCtx.EnterpriseRegistrationService.GetEnterpriseRegistration(l.ctx, userID)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	if registration == nil {
		return &types.GetEnterpriseRegistrationResp{}, nil
	}

	// Convert domain model to response
	resp = &types.GetEnterpriseRegistrationResp{
		ID:               registration.ID,
		UserID:           registration.UserID,
		CompanyName:      registration.CompanyName,
		CompanyType:      types.CompanyType(registration.CompanyType),
		ContactPerson:    registration.ContactPerson,
		Region:           registration.Region,
		DetailedAddress:  registration.DetailedAddress,
		LocationLatitude: registration.LocationLatitude,
		LocationLongitude: registration.LocationLongitude,
		Status:           registration.Status,
		CreatedAt:        registration.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        registration.UpdatedAt.Format(time.RFC3339),
	}

	return resp, nil
}
