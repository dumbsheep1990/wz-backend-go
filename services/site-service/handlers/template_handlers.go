package handlers

import (
	"net/http"
	"wz-backend-go/services/site-service/service"

	"github.com/gin-gonic/gin"
)

// ListTemplates 获取站点模板列表
func ListTemplates(c *gin.Context) {
	// 可选的分类过滤
	category := c.Query("category")

	templates, err := service.ListTemplates(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetTemplate 获取站点模板详情
func GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	template, err := service.GetTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	c.JSON(http.StatusOK, template)
}
