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

type UpdateTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewUpdateTenantLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *UpdateTenantLogic {
	return &UpdateTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTenantLogic) UpdateTenant(id string, req *http.AdminUpdateTenantReq) (*http.SuccessResp, error) {
	// 转换ID为int64
	tenantId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("租户ID格式错误: %v", err)
	}

	// 检查租户是否存在
	existingTenant, err := l.svcCtx.TenantRepo.GetTenantById(l.ctx, tenantId)
	if err != nil {
		return nil, fmt.Errorf("查询租户失败: %v", err)
	}
	if existingTenant == nil {
		return nil, fmt.Errorf("租户不存在")
	}

	// 准备更新数据
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.ContactEmail != "" {
		updates["contactEmail"] = req.ContactEmail
	}
	if req.ContactPhone != "" {
		updates["contactPhone"] = req.ContactPhone
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.ExpireAt != "" {
		// 解析过期时间
		expireTime, err := time.Parse("2006-01-02 15:04:05", req.ExpireAt)
		if err != nil {
			expireTime, err = time.Parse("2006-01-02", req.ExpireAt)
			if err != nil {
				return nil, fmt.Errorf("过期时间格式错误: %v", err)
			}
		}
		updates["expireAt"] = expireTime
	}

	// 没有要更新的字段
	if len(updates) == 0 {
		return &http.SuccessResp{
			Success: true,
			Message: "没有更新内容",
		}, nil
	}

	// 调用仓库层更新租户
	err = l.svcCtx.TenantRepo.UpdateTenant(l.ctx, tenantId, updates)
	if err != nil {
		return nil, fmt.Errorf("更新租户失败: %v", err)
	}

	// 返回成功结果
	return &http.SuccessResp{
		Success: true,
		Message: "更新成功",
	}, nil
} 