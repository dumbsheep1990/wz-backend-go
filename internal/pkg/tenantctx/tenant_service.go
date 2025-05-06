package tenantctx

import (
	"context"
)

// TenantAwareService 定义需要租户隔离的服务接口
type TenantAwareService interface {
	// WithContext 将包含租户信息的上下文应用到服务操作中
	WithContext(ctx context.Context) TenantAwareService
	
	// GetCurrentTenant 返回服务当前使用的租户ID
	GetCurrentTenant() (string, error)
	
	// IsTenantValid 验证租户ID是否存在且活跃
	IsTenantValid(tenantID string) (bool, error)
}

// BaseTenantService 提供了TenantAwareService的基本实现
type BaseTenantService struct {
	ctx context.Context
}

// NewBaseTenantService 使用给定的上下文创建新的BaseTenantService
func NewBaseTenantService(ctx context.Context) *BaseTenantService {
	return &BaseTenantService{ctx: ctx}
}

// WithContext 实现TenantAwareService.WithContext接口
func (s *BaseTenantService) WithContext(ctx context.Context) *BaseTenantService {
	return &BaseTenantService{ctx: ctx}
}

// GetCurrentTenant 实现TenantAwareService.GetCurrentTenant接口
func (s *BaseTenantService) GetCurrentTenant() (string, error) {
	tenantID, ok := GetTenantID(s.ctx)
	if !ok {
		return "", ErrMissingTenantID
	}
	return tenantID, nil
}

// GetCurrentPlatform 从上下文中获取当前平台
func (s *BaseTenantService) GetCurrentPlatform() (AppPlatform, error) {
	platform, ok := GetAppPlatform(s.ctx)
	if !ok {
		// 如果未指定，默认为Web平台
		return PlatformWeb, nil
	}
	return platform, nil
}

// IsSystemContext 检查当前上下文是否为系统上下文
func (s *BaseTenantService) IsSystemContext() bool {
	return IsSystemInternal(s.ctx)
}
