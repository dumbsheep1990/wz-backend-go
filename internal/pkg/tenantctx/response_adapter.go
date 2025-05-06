package tenantctx

import (
	"context"
)

// ResponseAdapter 定义基于平台适配的API响应方法
type ResponseAdapter interface {
	// AdaptResponse 根据平台需求修改响应
	// 如果响应被修改返回true，否则返回false
	AdaptResponse(ctx context.Context, response interface{}) (interface{}, bool)
	
	// SupportsPlatform 检查适配器是否支持给定平台
	SupportsPlatform(platform AppPlatform) bool
}

// BaseResponseAdapter 为特定平台的适配器提供基本实现
type BaseResponseAdapter struct {
	supportedPlatforms []AppPlatform
}

// NewBaseResponseAdapter 创建一个新的基本响应适配器
func NewBaseResponseAdapter(platforms ...AppPlatform) *BaseResponseAdapter {
	return &BaseResponseAdapter{
		supportedPlatforms: platforms,
	}
}

// SupportsPlatform 实现ResponseAdapter.SupportsPlatform接口
func (a *BaseResponseAdapter) SupportsPlatform(platform AppPlatform) bool {
	for _, p := range a.supportedPlatforms {
		if p == platform {
			return true
		}
	}
	return false
}

// ResponseAdapterRegistry 是响应适配器的注册表
type ResponseAdapterRegistry struct {
	adapters []ResponseAdapter
}

// NewResponseAdapterRegistry 创建一个新的响应适配器注册表
func NewResponseAdapterRegistry() *ResponseAdapterRegistry {
	return &ResponseAdapterRegistry{
		adapters: make([]ResponseAdapter, 0),
	}
}

// RegisterAdapter 注册一个响应适配器
func (r *ResponseAdapterRegistry) RegisterAdapter(adapter ResponseAdapter) {
	r.adapters = append(r.adapters, adapter)
}

// AdaptResponse 根据上下文中的平台调整响应
func (r *ResponseAdapterRegistry) AdaptResponse(ctx context.Context, response interface{}) interface{} {
	// 从上下文中获取平台
	platform, ok := GetAppPlatform(ctx)
	if !ok {
		// 如果未指定，默认为Web平台
		platform = PlatformWeb
	}
	
	// 查找适用的适配器
	for _, adapter := range r.adapters {
		if adapter.SupportsPlatform(platform) {
			// 应用适配器
			if adaptedResponse, modified := adapter.AdaptResponse(ctx, response); modified {
				response = adaptedResponse
			}
		}
	}
	
	return response
}
