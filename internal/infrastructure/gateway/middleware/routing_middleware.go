package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/domain/gateway/service"
)

// RoutingMiddleware 路由中间件，负责解析和匹配请求路由
type RoutingMiddleware struct {
	domainService service.GatewayDomainService
}

// NewRoutingMiddleware 创建新的路由中间件
func NewRoutingMiddleware(domainService service.GatewayDomainService) *RoutingMiddleware {
	return &RoutingMiddleware{
		domainService: domainService,
	}
}

// Handle 处理请求路由
func (m *RoutingMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中解析匹配的路由
		route, err := m.domainService.ResolveRoute(c.Request)
		if err != nil {
			log.Printf("路由解析失败: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "没有找到匹配的路由",
			})
			c.Abort()
			return
		}

		// 获取目标服务
		service, err := m.domainService.GetService(route.ServiceID)
		if err != nil {
			log.Printf("获取服务失败: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    502,
				"message": "无法访问目标服务",
			})
			c.Abort()
			return
		}

		// 检查服务是否激活
		if !service.IsActive {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    503,
				"message": "目标服务当前不可用",
			})
			c.Abort()
			return
		}

		// 将路由和服务存储在上下文中，供后续中间件使用
		c.Set("route", route)
		c.Set("service", service)

		// 继续处理请求
		c.Next()
	}
}
