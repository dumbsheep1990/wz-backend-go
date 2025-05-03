package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/gateway/auth"
	"wz-backend-go/internal/gateway/config"
)

// EnhancedAuthentication 增强的认证中间件
func EnhancedAuthentication(authManager *auth.AuthManager, securityConfig config.SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于OPTIONS请求直接放行
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 检测认证方式
		authType := detectAuthType(c, securityConfig)
		if authType == auth.AuthTypeNone {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "需要认证",
			})
			return
		}

		// 获取认证提供者
		provider, err := authManager.GetProvider(authType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "不支持的认证方式",
			})
			return
		}

		// 验证用户
		user, err := provider.Authenticate(c)
		if err != nil {
			var statusCode int
			var message string

			switch err {
			case auth.ErrCredentialsNotFound:
				statusCode = http.StatusUnauthorized
				message = "未提供认证凭据"
			case auth.ErrInvalidCredentials:
				statusCode = http.StatusUnauthorized
				message = "无效的认证凭据"
			case auth.ErrTokenExpired:
				statusCode = http.StatusUnauthorized
				message = "认证凭据已过期"
			default:
				statusCode = http.StatusInternalServerError
				message = "认证过程中发生错误"
			}

			c.AbortWithStatusJSON(statusCode, gin.H{
				"code":    statusCode,
				"message": message,
				"error":   err.Error(),
			})
			return
		}

		// 认证成功，将用户信息存储到上下文
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)
		c.Set("tenantID", user.TenantID)
		c.Set("authUser", user)

		c.Next()
	}
}

// RBACAuthorization RBAC授权中间件
func RBACAuthorization(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无法获取用户角色信息",
			})
			return
		}

		userRole, ok := role.(string)
		if !ok || userRole == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无效的用户角色",
			})
			return
		}

		// 检查角色是否拥有权限
		if len(requiredRoles) > 0 {
			authorized := false
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					authorized = true
					break
				}
			}

			if !authorized {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "权限不足，需要以下角色之一: " + strings.Join(requiredRoles, ", "),
				})
				return
			}
		}

		c.Next()
	}
}

// ResourceAuthorization 资源级别授权中间件
func ResourceAuthorization(authManager *auth.AuthManager, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无法获取用户信息",
			})
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无效的用户ID",
			})
			return
		}

		// 检查用户是否有权限访问资源
		if !authManager.HasPermission(userIDStr, resource, action) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权访问该资源",
				"details": gin.H{
					"resource": resource,
					"action":   action,
				},
			})
			return
		}

		c.Next()
	}
}

// detectAuthType 检测认证类型
func detectAuthType(c *gin.Context, config config.SecurityConfig) string {
	// 检查请求头中的 Authorization
	authHeader := c.GetHeader("Authorization")

	// 检测 Bearer 令牌 (JWT/OAuth2)
	if strings.HasPrefix(authHeader, "Bearer ") {
		// 如果 OAuth2 启用，需要更复杂的逻辑来区分 JWT 和 OAuth2
		// 简单实现：默认使用 JWT
		for _, method := range config.AuthMethods {
			if method == auth.AuthTypeOAuth2 && config.OAuth2.Enabled {
				return auth.AuthTypeOAuth2
			}
			if method == auth.AuthTypeJWT {
				return auth.AuthTypeJWT
			}
		}
		return auth.AuthTypeJWT
	}

	// 检测 Basic 认证
	if strings.HasPrefix(authHeader, "Basic ") {
		for _, method := range config.AuthMethods {
			if method == auth.AuthTypeBasic {
				return auth.AuthTypeBasic
			}
		}
	}

	// 检测 API Key
	apiKeyHeader := config.APIKey.HeaderName
	if apiKeyHeader == "" {
		apiKeyHeader = "X-API-Key"
	}
	if c.GetHeader(apiKeyHeader) != "" {
		for _, method := range config.AuthMethods {
			if method == auth.AuthTypeAPIKey && config.APIKey.Enabled {
				return auth.AuthTypeAPIKey
			}
		}
	}

	// 检测查询参数中的 API Key
	apiKeyParam := config.APIKey.QueryParamName
	if apiKeyParam == "" {
		apiKeyParam = "api_key"
	}
	if c.Query(apiKeyParam) != "" {
		for _, method := range config.AuthMethods {
			if method == auth.AuthTypeAPIKey && config.APIKey.Enabled {
				return auth.AuthTypeAPIKey
			}
		}
	}

	// 如果没有检测到有效的认证凭据，使用默认认证类型
	if config.DefaultAuth != "" {
		return config.DefaultAuth
	}

	// 没有有效的认证类型
	return auth.AuthTypeNone
}
