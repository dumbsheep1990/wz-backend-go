package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/wz-project/wz-backend-go/services/community-service/models"
)

// RegisterCategoryRoutes 注册分类相关的路由
func RegisterCategoryRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	// 获取所有同X分类
	v1.GET("/similar-categories", func(c *gin.Context) {
		categories := models.GetAllSimilarCategories()
		
		// 转换为JSON友好的格式
		result := make([]map[string]string, 0, len(categories))
		for _, cat := range categories {
			result = append(result, map[string]string{
				"code": string(cat),
				"name": string(cat),
			})
		}
		
		c.JSON(200, gin.H{
			"data": result,
		})
	})
}
