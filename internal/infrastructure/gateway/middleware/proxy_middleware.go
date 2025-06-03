package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/infrastructure/gateway/proxy"
)

// ProxyMiddleware 创建代理中间件处理函数
func ProxyMiddleware() gin.HandlerFunc {
	// 创建代理处理器
	proxyHandler := proxy.NewProxyHandler()

	return func(c *gin.Context) {
		// 从上下文中获取路由和服务信息
		routeInterface, exists := c.Get("route")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误：路由信息缺失",
			})
			c.Abort()
			return
		}

		serviceInterface, exists := c.Get("service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误：服务信息缺失",
			})
			c.Abort()
			return
		}

		route, ok := routeInterface.(*entity.Route)
		if !ok {
			log.Println("代理中间件：路由信息类型错误")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误",
			})
			c.Abort()
			return
		}

		service, ok := serviceInterface.(*entity.Service)
		if !ok {
			log.Println("代理中间件：服务信息类型错误")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "内部服务器错误",
			})
			c.Abort()
			return
		}

		// 执行代理请求
		err := proxyHandler.ProxyRequest(c.Request.Context(), c.Writer, c.Request, route, service)
		if err != nil {
			log.Printf("代理请求失败: %v", err)
			// 如果代理请求失败且尚未写入响应，则返回错误
			if !c.Writer.Written() {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":    502,
					"message": "代理请求失败",
					"error":   err.Error(),
				})
			}
			c.Abort()
			return
		}

		// 代理请求已处理完成，无需继续处理
		c.Abort()
	}
}
