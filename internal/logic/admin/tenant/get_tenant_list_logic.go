package tenant

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewGetTenantListLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *GetTenantListLogic {
	return &GetTenantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTenantListLogic) GetTenantList(req *http.TenantListReq) (*http.TenantListResp, error) {
	// 构建查询过滤条件
	filters := make(map[string]interface{})
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.Subdomain != "" {
		filters["subdomain"] = req.Subdomain
	}
	if req.Type > 0 {
		filters["type"] = req.Type
	}
	if req.Status > 0 {
		filters["status"] = req.Status
	}
	if req.StartTime != "" {
		filters["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		filters["endTime"] = req.EndTime
	}

	// 调用仓库层查询数据
	tenants, total, err := l.svcCtx.TenantRepo.GetTenantList(l.ctx, req.Page, req.PageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("查询租户列表失败: %v", err)
	}

	// 转换数据格式
	var list []http.TenantDetail
	for _, tenant := range tenants {
		var expireAt string
		if !tenant.ExpireAt.IsZero() {
			expireAt = tenant.ExpireAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, http.TenantDetail{
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
		})
	}

	return &http.TenantListResp{
		Total: total,
		List:  list,
	}, nil
} 