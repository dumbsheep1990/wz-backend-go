package users

import (
	"context"
	"errors"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterEnterpriseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterEnterpriseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterEnterpriseLogic {
	return &RegisterEnterpriseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// validateEnterpriseRegistration 验证企业入驻信息
func (l *RegisterEnterpriseLogic) validateEnterpriseRegistration(req *types.EnterpriseRegistrationReq) error {
	// 验证必填字段
	if req.CompanyName == "" {
		return errors.New("公司名称为必填项")
	}
	if req.ContactPerson == "" {
		return errors.New("联系人为必填项")
	}
	if req.JobPosition == "" {
		return errors.New("职位为必填项")
	}
	if req.Region == "" {
		return errors.New("所在地区为必填项")
	}
	if req.VerificationMethod == "" {
		return errors.New("认证方式为必填项")
	}
	if req.DetailedAddress == "" {
		return errors.New("详细地址为必填项")
	}

	// 验证公司类型
	switch req.CompanyType {
	case types.CompanyTypeEnterprise, 
		 types.CompanyTypeGroup, 
		 types.CompanyTypeGovernment, 
		 types.CompanyTypeResearchInst:
		// 有效的公司类型
	default:
		return errors.New("公司类型不正确")
	}

	return nil
}

func (l *RegisterEnterpriseLogic) RegisterEnterprise(req *types.EnterpriseRegistrationReq) (resp *types.EnterpriseRegistrationResp, err error) {
	// 验证企业入驻信息
	if err = l.validateEnterpriseRegistration(req); err != nil {
		return nil, err
	}
	
	// 从上下文中获取用户ID
	userID, ok := l.ctx.Value("user_id").(int64)
	if !ok {
		return nil, errors.New("未登录或登录信息无效")
	}
	
	// 调用RPC服务进行企业入驻请求处理
	rpcReq := &user.EnterpriseRegistrationRequest{
		UserID:             userID,
		CompanyName:        req.CompanyName,
		CompanyType:        int32(req.CompanyType),
		ContactPerson:      req.ContactPerson,
		JobPosition:        req.JobPosition,
		Region:             req.Region,
		VerificationMethod: req.VerificationMethod,
		DetailedAddress:    req.DetailedAddress,
		LocationLatitude:   float32(req.LocationLatitude),
		LocationLongitude:  float32(req.LocationLongitude),
	}
	
	result, err := l.svcCtx.UserRpc.RegisterEnterprise(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	
	resp = &types.EnterpriseRegistrationResp{
		Success: result.Success,
	}
	
	return
}
