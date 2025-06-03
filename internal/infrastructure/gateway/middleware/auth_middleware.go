package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/service"
)

// AuthMiddleware 认证中间件，负责处理请求身份验证
type AuthMiddleware struct {
	domainService service.GatewayDomainService
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(domainService service.GatewayDomainService) *AuthMiddleware {
	return &AuthMiddleware{
		domainService: domainService,
	}
}

// Handle 处理请求身份验证
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取路由信息
		routeInterface, exists := c.Get("route")
		if !exists {
			log.Println("认证中间件：找不到路由信息")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误",
			})
			c.Abort()
			return
		}

		route, ok := routeInterface.(*entity.Route)
		if !ok {
			log.Println("认证中间件：路由信息类型错误")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误",
			})
			c.Abort()
			return
		}

		// 检查路由是否需要身份验证
		if !route.AuthRequired {
			// 不需要身份验证，继续处理请求
			c.Next()
			return
		}

		// 从请求中提取访问令牌
		tokenString, err := m.domainService.ExtractToken(c.Request)
		if err != nil {
			log.Printf("提取访问令牌失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "需要身份验证",
			})
			c.Abort()
			return
		}

		// 验证访问令牌
		session, err := m.domainService.VerifyToken(c, tokenString)
		if err != nil {
			log.Printf("令牌验证失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的访问令牌",
			})
			c.Abort()
			return
		}

		// 检查是否有权限访问
		if len(route.Permissions) > 0 {
			hasPermission := false
			for _, permission := range route.Permissions {
				if m.domainService.CheckPermission(session, permission) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "没有访问权限",
				})
				c.Abort()
				return
			}
		}

		// 将用户会话存储在上下文中，供后续处理使用
		c.Set("user_session", session)

		// 继续处理请求
		c.Next()
	}
}
