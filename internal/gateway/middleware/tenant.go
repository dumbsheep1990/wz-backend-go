package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TenantIdentifier 返回一个租户识别中间件
func TenantIdentifier() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID string

		// 尝试从请求头获取租户ID
		tenantID = c.GetHeader("X-Tenant-ID")

		// 如果请求头中没有租户ID，则尝试从子域名解析
		if tenantID == "" {
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
					// TODO: 查询数据库，将子域名映射到租户ID
					// 这里为了演示，直接使用子域名作为租户ID
					tenantID = tenantSubdomain
				}
			}
		}

		// 如果我们确定了租户ID，则将其保存在上下文中
		if tenantID != "" {
			c.Set("tenantID", tenantID)
			// 强制设置请求头，以便传递给后端服务
			c.Request.Header.Set("X-Tenant-ID", tenantID)
		} else {
			// 检查是否是公共路由，如果不是，可能需要返回错误
			path := c.Request.URL.Path
			if !isPublicRoute(path) {
				// 注意：实际项目中，可能会重定向到错误页面或总站，这里简单返回错误
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "无法确定租户信息",
				})
				return
			}
		}

		c.Next()
	}
}

// isIPAddress 检查主机名是否为IP地址
func isIPAddress(host string) bool {
	// 简单检查：IP地址第一部分通常是数字
	parts := strings.Split(host, ".")
	if len(parts) != 4 {
		return false
	}
	
	// 检查所有部分是否为数字
	for _, part := range parts {
		if !isNumeric(part) {
			return false
		}
	}
	
	return true
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isPublicRoute 检查是否为公共路由
func isPublicRoute(path string) bool {
	// 定义公共路由前缀列表
	publicPrefixes := []string{
		"/api/public",
		"/health",
		"/swagger",
		"/metrics",
		"/auth/login",
		"/auth/register",
	}
	
	// 检查路径是否以公共前缀开头
	for _, prefix := range publicPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	
	return false
}
