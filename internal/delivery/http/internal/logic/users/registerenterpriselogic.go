package users

import (
	"context"
	"errors"
	"fmt"

	"wz-backend-go/internal/delivery/http/internal/middleware"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/delivery/rpc/user"
	"wz-backend-go/internal/domain/model"

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

	// 多租户相关验证
	if req.Subdomain == "" {
		return errors.New("子域名为必填项")
	}
	if req.TenantName == "" {
		return errors.New("租户名称为必填项")
	}
	
	return nil
}

func (l *RegisterEnterpriseLogic) RegisterEnterprise(req *types.EnterpriseRegistrationReq) (resp *types.EnterpriseRegistrationResp, err error) {
	// 验证企业入驻信息
	if err = l.validateEnterpriseRegistration(req); err != nil {
		return nil, err
	}
	
	// 从上下文中获取用户ID
	userID, ok := middleware.GetUserIDFromContext(l.ctx)
	if !ok {
		return nil, errors.New("未登录或登录信息无效")
	}
	
	// 首先验证子域名是否已存在
	existingTenant, err := l.svcCtx.TenantService.GetTenantBySubdomain(l.ctx, req.Subdomain)
	if err == nil && existingTenant != nil {
		return nil, errors.New("子域名已被占用，请选择其他子域名")
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
		return nil, fmt.Errorf("企业注册失败: %v", err)
	}
	
	// 创建租户
	tenantReq := &model.CreateTenantRequest{
		Name:        req.TenantName,
		Subdomain:   req.Subdomain,
		TenantType:  model.TenantTypeEnterprise, // 根据公司类型设置对应的租户类型
		Description: req.TenantDesc,
		Logo:        "", // 可以在后续更新
	}
	
	tenant, err := l.svcCtx.TenantService.CreateTenant(l.ctx, tenantReq, userID)
	if err != nil {
		return nil, fmt.Errorf("租户创建失败: %v", err)
	}
	
	// 将用户与租户关联，角色为租户管理员
	err = l.svcCtx.TenantService.AddUserToTenant(l.ctx, tenant.ID, userID, string(model.RoleTenantAdmin))
	if err != nil {
		return nil, fmt.Errorf("关联用户与租户失败: %v", err)
	}
	
	// 生成租户管理员的JWT令牌
	tokenPair, err := middleware.GenerateTokenPair(
		userID,
		req.ContactPerson, // 使用联系人作为用户名
		tenant.ID,
		model.RoleTenantAdmin,
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}
	
	resp = &types.EnterpriseRegistrationResp{
		Success:      result.Success,
		TenantID:     tenant.ID,
		Subdomain:    tenant.Subdomain,
		TenantName:   tenant.Name,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
	}
	
	return
}
