package tenant

import (
	"context"
	"fmt"
	"strconv"

	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewGetTenantDetailLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *GetTenantDetailLogic {
	return &GetTenantDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTenantDetailLogic) GetTenantDetail(id string) (*http.TenantDetail, error) {
	// 转换ID为int64
	tenantId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("租户ID格式错误: %v", err)
	}

	// 调用仓库层查询数据
	tenant, err := l.svcCtx.TenantRepo.GetTenantById(l.ctx, tenantId)
	if err != nil {
		return nil, fmt.Errorf("查询租户详情失败: %v", err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("租户不存在")
	}

	// 转换数据格式
	var expireAt string
	if !tenant.ExpireAt.IsZero() {
		expireAt = tenant.ExpireAt.Format("2006-01-02 15:04:05")
	}

	result := &http.TenantDetail{
		ID:           tenant.ID,
		Name:         tenant.Name,
		Description:  tenant.Description,
		Subdomain:    tenant.Subdomain,
		Type:         int(tenant.Type),
		Status:       int(tenant.Status),
		Logo:         tenant.Logo,
		ContactEmail: tenant.ContactEmail,
		ContactPhone: tenant.ContactPhone,
		AdminUserID:  tenant.AdminUserID,
		CreatedAt:    tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		ExpireAt:     expireAt,
	}

	return result, nil
} 