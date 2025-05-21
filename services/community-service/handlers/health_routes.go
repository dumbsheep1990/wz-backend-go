package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterHealthRoutes 注册健康检查相关路由
func RegisterHealthRoutes(router *gin.Engine) {
	// 健康检查路由
	api := router.Group("/api/v1")
	{
		api.GET("/health", HealthHandler)
	}
}
