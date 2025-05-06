package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/wxnacy/wz-backend-go/internal/pkg/tenantctx"
)

// EnhancedTenantMiddleware 增强版租户中间件，支持多种平台类型识别
// 同时支持子域名和请求头方式传递租户信息
func EnhancedTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 识别请求来源的平台类型
		platform := detectPlatform(c)
		
		// 2. 从多种来源获取租户ID
		tenantID := extractTenantID(c)
		
		// 3. 将信息添加到上下文
		if tenantID != "" {
			// 创建带有租户ID的上下文
			ctx := tenantctx.WithTenantID(c.Request.Context(), tenantID)
			
			// 添加平台信息
			ctx = tenantctx.WithAppPlatform(ctx, platform)
			
			// 更新请求上下文
			c.Request = c.Request.WithContext(ctx)
			
			// 同时保存在Gin上下文中，方便其他中间件使用
			c.Set("tenantID", tenantID)
			c.Set("appPlatform", string(platform))
			
			// 设置请求头，传递给后端服务
			c.Request.Header.Set("X-Tenant-ID", tenantID)
			c.Request.Header.Set("X-App-Platform", string(platform))
			
			logx.Infof("租户请求: TenantID=%s, Platform=%s, Path=%s", 
				tenantID, platform, c.Request.URL.Path)
		} else {
			// 检查是否是公共路由
			path := c.Request.URL.Path
			if !isPublicRoute(path) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "无法确定租户信息",
				})
				return
			}
			
			// 即使是公共路由，也需要携带平台信息
			ctx := tenantctx.WithAppPlatform(c.Request.Context(), platform)
			c.Request = c.Request.WithContext(ctx)
			c.Set("appPlatform", string(platform))
			c.Request.Header.Set("X-App-Platform", string(platform))
		}
		
		c.Next()
	}
}

// detectPlatform 根据请求信息检测平台类型
func detectPlatform(c *gin.Context) tenantctx.AppPlatform {
	// 优先从请求头获取平台信息
	if platform := c.GetHeader("X-App-Platform"); platform != "" {
		switch platform {
		case string(tenantctx.PlatformMobile):
			return tenantctx.PlatformMobile
		case string(tenantctx.PlatformUniApp):
			return tenantctx.PlatformUniApp
		case string(tenantctx.PlatformWeb):
			return tenantctx.PlatformWeb
		}
	}
	
	// 从User-Agent判断
	userAgent := c.GetHeader("User-Agent")
	userAgent = strings.ToLower(userAgent)
	
	// 移动设备UA通常包含 mobile, android, iphone 等关键词
	if strings.Contains(userAgent, "mobile") || 
	   strings.Contains(userAgent, "android") || 
	   strings.Contains(userAgent, "iphone") {
		return tenantctx.PlatformMobile
	}
	
	// UniApp通常会在UA中包含特定标识
	if strings.Contains(userAgent, "uni-app") || 
	   strings.Contains(userAgent, "uniapp") {
		return tenantctx.PlatformUniApp
	}
	
	// 默认为Web平台
	return tenantctx.PlatformWeb
}

// extractTenantID 从多种来源提取租户ID
func extractTenantID(c *gin.Context) string {
	// 1. 优先从请求头获取
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID != "" {
		return tenantID
	}
	
	// 2. 尝试从子域名解析
	host := c.Request.Host
	// 移除端口号部分（如果有）
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}
	
	// 如果不是IP地址，解析子域名
	if !isIPAddress(host) {
		parts := strings.Split(host, ".")
		// 简单判断：如果有足够的部分，并且不是 www，则第一部分可能是租户子域名
		if len(parts) >= 3 && parts[0] != "www" {
			// 提取第一部分作为租户标识符
			tenantSubdomain := parts[0]
			// 实际应用中需要查询数据库，将子域名映射到租户ID
			// 这里为了示例，直接返回子域名
			return tenantSubdomain
		}
	}
	
	// 3. 从URL查询参数中获取（适用于移动端和UniApp）
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		return tenantID
	}
	
	return ""
}
