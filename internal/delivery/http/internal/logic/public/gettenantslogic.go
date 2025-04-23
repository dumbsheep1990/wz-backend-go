package public

import (
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

type GetTenantsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTenantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantsLogic {
	return &GetTenantsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTenants 获取所有可用租户列表
func (l *GetTenantsLogic) GetTenants() (*model.TenantListResponse, error) {
	// 调用租户服务获取所有可用租户
	tenants, err := l.svcCtx.TenantService.ListActiveTenants(l.ctx)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var tenantInfos []model.TenantInfo
	for _, tenant := range tenants {
		tenantInfos = append(tenantInfos, model.TenantInfo{
			ID:          strconv.FormatInt(tenant.ID, 10),
			Name:        tenant.Name,
			Description: tenant.Description,
			Subdomain:   tenant.Subdomain,
			TenantType:  tenantTypeToString(tenant.TenantType),
			Logo:        tenant.Logo,
		})
	}

	return &model.TenantListResponse{
		Tenants: tenantInfos,
	}, nil
}

// tenantTypeToString 将租户类型转换为字符串描述
func tenantTypeToString(tenantType model.TenantType) string {
	switch tenantType {
	case model.TenantTypeEnterprise:
		return "企业"
	case model.TenantTypePersonal:
		return "个人"
	case model.TenantTypeEducational:
		return "教育机构"
	default:
		return "未知"
	}
}
