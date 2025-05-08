package tenant

import (
	"context"
	"fmt"
	"time"

	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTenantLogic) CreateTenant(req *http.AdminCreateTenantReq) (*http.AdminCreateTenantResp, error) {
	// 检查子域名是否已存在
	existingTenant, err := l.svcCtx.TenantRepo.GetTenantBySubdomain(l.ctx, req.Subdomain)
	if err != nil {
		return nil, fmt.Errorf("验证子域名失败: %v", err)
	}
	if existingTenant != nil {
		return nil, fmt.Errorf("子域名已被使用")
	}

	// 准备创建租户数据
	tenant := &repository.Tenant{
		Name:         req.Name,
		Description:  req.Description,
		Subdomain:    req.Subdomain,
		Type:         int32(req.Type),
		Status:       int32(req.Status),
		Logo:         req.Logo,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		AdminUserID:  req.AdminUserID,
	}

	// 设置过期时间，如果提供了
	if req.ExpireAt != "" {
		expireTime, err := time.Parse("2006-01-02 15:04:05", req.ExpireAt)
		if err != nil {
			expireTime, err = time.Parse("2006-01-02", req.ExpireAt)
			if err != nil {
				return nil, fmt.Errorf("过期时间格式错误: %v", err)
			}
		}
		tenant.ExpireAt = expireTime
	} else {
		// 默认一年后过期
		tenant.ExpireAt = time.Now().AddDate(1, 0, 0)
	}

	// 如果没有提供管理员用户ID，需要创建一个默认管理员
	if tenant.AdminUserID == 0 {
		// 实际实现中，这里应该创建一个默认的管理员用户
		// 暂时为保持简单，设置一个占位符
		tenant.AdminUserID = 1 // 假设ID为1的用户是系统默认管理员
	}

	// 调用仓库层创建租户
	id, err := l.svcCtx.TenantRepo.CreateTenant(l.ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("创建租户失败: %v", err)
	}

	// 返回创建结果
	return &http.AdminCreateTenantResp{
		ID:        id,
		Name:      req.Name,
		Subdomain: req.Subdomain,
	}, nil
} 