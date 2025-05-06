package tenantctx

import (
	"context"
	"errors"
)

type contextKey string

const (
	// TenantIDKey 用于在上下文中存储租户ID的键
	TenantIDKey contextKey = "tenantID"
	
	// TenantTypeKey 用于在上下文中存储租户类型的键
	TenantTypeKey contextKey = "tenantType"
	
	// AppPlatformKey 用于在上下文中存储应用平台信息的键
	AppPlatformKey contextKey = "appPlatform"
	
	// SystemInternalKey 用于标识系统内部操作的键
	SystemInternalKey contextKey = "systemInternal"
)

// AppPlatform 表示请求来源的平台
type AppPlatform string

const (
	// PlatformWeb 表示来自网页浏览器的请求
	PlatformWeb AppPlatform = "web"
	
	// PlatformMobile 表示来自原生移动应用的请求
	PlatformMobile AppPlatform = "mobile"
	
	// PlatformUniApp 表示来自uniapp应用的请求
	PlatformUniApp AppPlatform = "uniapp"
)

// 常见错误
var (
	ErrMissingTenantID   = errors.New("缺少租户ID")
	ErrInvalidTenantID   = errors.New("无效的租户ID")
	ErrMissingTenantType = errors.New("缺少租户类型")
)

// WithTenantID 向上下文添加租户ID
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID 从上下文中获取租户ID
func GetTenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok && tenantID != ""
}

// WithTenantType 向上下文添加租户类型
func WithTenantType(ctx context.Context, tenantType int) context.Context {
	return context.WithValue(ctx, TenantTypeKey, tenantType)
}

// GetTenantType 从上下文中获取租户类型
func GetTenantType(ctx context.Context) (int, bool) {
	tenantType, ok := ctx.Value(TenantTypeKey).(int)
	return tenantType, ok && tenantType > 0
}

// WithAppPlatform 向上下文添加应用平台信息
func WithAppPlatform(ctx context.Context, platform AppPlatform) context.Context {
	return context.WithValue(ctx, AppPlatformKey, platform)
}

// GetAppPlatform 从上下文中获取应用平台信息
func GetAppPlatform(ctx context.Context) (AppPlatform, bool) {
	platform, ok := ctx.Value(AppPlatformKey).(AppPlatform)
	return platform, ok && platform != ""
}

// SystemContext 创建用于系统内部操作的上下文
// 该上下文会绕过租户隔离
func SystemContext() context.Context {
	return context.WithValue(context.Background(), SystemInternalKey, true)
}

// IsSystemInternal 检查上下文是否用于系统内部操作
func IsSystemInternal(ctx context.Context) bool {
	isInternal, ok := ctx.Value(SystemInternalKey).(bool)
	return ok && isInternal
}
