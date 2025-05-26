package users

import (
	"context"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"
	"github.com/wz-project/wz-backend-go/internal/domain/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnterpriseRegistrationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnterpriseRegistrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnterpriseRegistrationLogic {
	return &EnterpriseRegistrationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnterpriseRegistrationLogic) EnterpriseRegistration(req *types.EnterpriseRegistrationReq) (resp *types.EnterpriseRegistrationResp, err error) {
	// Convert HTTP request to domain request
	domainReq := &model.EnterpriseRegistrationRequest{
		UserID:             req.UserID,
		CompanyName:        req.CompanyName,
		CompanyType:        model.CompanyType(req.CompanyType),
		ContactPerson:      req.ContactPerson,
		JobPosition:        req.JobPosition,
		Region:             req.Region,
		VerificationMethod: req.VerificationMethod,
		DetailedAddress:    req.DetailedAddress,
		LocationLatitude:   req.LocationLatitude,
		LocationLongitude:  req.LocationLongitude,
		Subdomain:          req.Subdomain,
		TenantType:         model.TenantType(1), // Default to organization
		TenantName:         req.TenantName,
		TenantDesc:         req.TenantDesc,
	}

	// Call domain service
	registration, err := l.svcCtx.EnterpriseRegistrationService.CreateEnterpriseRegistration(l.ctx, domainReq)
	if err != nil {
		return nil, types.NewInternalError(err.Error())
	}

	// Create tenant and tokens for response
	// In a real implementation, you would handle tenant creation and token generation here
	// For simplicity, we'll just return a success response

	resp = &types.EnterpriseRegistrationResp{
		Success:      true,
		TenantID:     0, // Would be set in real implementation
		Subdomain:    req.Subdomain,
		TenantName:   req.TenantName,
		AccessToken:  "", // Would be set in real implementation
		RefreshToken: "", // Would be set in real implementation
		ExpiresAt:    0,  // Would be set in real implementation
		TokenType:    "Bearer",
	}

	return resp, nil
}
