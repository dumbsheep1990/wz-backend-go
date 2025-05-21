package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler 处理健康检查请求
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "community-service",
		"version": "1.0.0",
	})
}
