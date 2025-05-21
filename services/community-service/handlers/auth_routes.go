package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/wz-backend-go/services/community-service/middleware"
)

// RegisterAuthRoutes 注册所有认证相关的路由
func RegisterAuthRoutes(router *gin.Engine) {
	// 认证相关路由
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", LoginHandler)
		auth.POST("/register", RegisterHandler)
	}
}
