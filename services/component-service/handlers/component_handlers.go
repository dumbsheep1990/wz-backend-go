package handlers

import (
	"net/http"
	"wz-backend-go/models"
	"wz-backend-go/services/component-service/service"

	"github.com/gin-gonic/gin"
)

// ListComponentCategories 获取组件分类列表
func ListComponentCategories(c *gin.Context) {
	categories, err := service.ListComponentCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetComponentDefinition 获取组件定义
func GetComponentDefinition(c *gin.Context) {
	componentType := c.Param("type")

	definition, err := service.GetComponentDefinition(componentType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "组件类型不存在"})
		return
	}

	c.JSON(http.StatusOK, definition)
}

// AddComponent 添加组件到区块
func AddComponent(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("sectionId")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	var component models.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置组件属性
	component.SectionID = sectionID

	// 添加组件
	addedComponent, err := service.AddComponent(siteID, pageID, sectionID, component)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, addedComponent)
}

// UpdateComponent 更新组件
func UpdateComponent(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("sectionId")
	componentID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	var component models.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置组件属性
	component.ID = componentID
	component.SectionID = sectionID

	// 更新组件
	updatedComponent, err := service.UpdateComponent(siteID, pageID, sectionID, component)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedComponent)
}

// DeleteComponent 删除组件
func DeleteComponent(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("sectionId")
	componentID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	// 删除组件
	if err := service.DeleteComponent(siteID, pageID, sectionID, componentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组件已删除"})
}

// ReorderComponents 重新排序组件
func ReorderComponents(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("sectionId")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	var componentOrder []string
	if err := c.ShouldBindJSON(&componentOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新排序组件
	if err := service.ReorderComponents(siteID, pageID, sectionID, componentOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组件顺序已更新"})
}
