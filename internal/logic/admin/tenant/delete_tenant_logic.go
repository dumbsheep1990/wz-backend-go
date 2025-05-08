package tenant

import (
	"context"
	"fmt"
	"strconv"

	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewDeleteTenantLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *DeleteTenantLogic {
	return &DeleteTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTenantLogic) DeleteTenant(id string) (*http.SuccessResp, error) {
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

	// 调用仓库层删除租户
	err = l.svcCtx.TenantRepo.DeleteTenant(l.ctx, tenantId)
	if err != nil {
		return nil, fmt.Errorf("删除租户失败: %v", err)
	}

	// 返回成功结果
	return &http.SuccessResp{
		Success: true,
		Message: "删除成功",
	}, nil
} 