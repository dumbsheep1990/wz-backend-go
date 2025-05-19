package handlers

import (
	"net/http"
	"wz-backend-go/models"
	"wz-backend-go/services/component-service/service"

	"github.com/gin-gonic/gin"
)

// ListComponentCategories 获取组件分类列表
// @Summary 获取组件分类列表
// @Description 获取所有可用的组件分类列表
// @Tags Component
// @Accept json
// @Produce json
// @Success 200 {array} models.ComponentCategory
// @Router /component-categories [get]
func ListComponentCategories(c *gin.Context) {
	categories, err := service.ListComponentCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetComponentDefinition 获取组件定义
// @Summary 获取组件定义
// @Description 获取指定类型的组件定义
// @Tags Component
// @Accept json
// @Produce json
// @Param type path string true "组件类型"
// @Success 200 {object} models.ComponentDefinition
// @Router /components/{type} [get]
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
// @Summary 添加组件到区块
// @Description 将组件添加到指定的区块
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param pageId path string true "页面ID"
// @Param sectionId path string true "区块ID"
// @Param data body models.Component true "组件信息"
// @Success 201 {object} models.Component
// @Router /sites/{siteId}/pages/{pageId}/sections/{sectionId}/components [post]
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
// @Summary 更新组件
// @Description 更新指定组件的信息
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param pageId path string true "页面ID"
// @Param sectionId path string true "区块ID"
// @Param id path string true "组件ID"
// @Param data body models.Component true "组件信息"
// @Success 200 {object} models.Component
// @Router /sites/{siteId}/pages/{pageId}/sections/{sectionId}/components/{id} [put]
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
// @Summary 删除组件
// @Description 删除指定的组件
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param pageId path string true "页面ID"
// @Param sectionId path string true "区块ID"
// @Param id path string true "组件ID"
// @Success 200 {object} gin.H
// @Router /sites/{siteId}/pages/{pageId}/sections/{sectionId}/components/{id} [delete]
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
// @Summary 重新排序组件
// @Description 重新排序指定区块中的组件
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param pageId path string true "页面ID"
// @Param sectionId path string true "区块ID"
// @Param data body []string true "组件顺序"
// @Success 200 {object} gin.H
// @Router /sites/{siteId}/pages/{pageId}/sections/{sectionId}/components/reorder [post]
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

// CreateComponent 创建组件
// @Summary 创建组件
// @Description 创建一个新的组件
// @Tags Component
// @Accept json
// @Produce json
// @Param data body models.Component true "组件信息"
// @Success 201 {object} models.Component
// @Router /components [post]
func CreateComponent(c *gin.Context) {
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

	// 创建组件
	createdComponent, err := service.CreateComponent(siteID, pageID, sectionID, component)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdComponent)
}

// ListComponents 获取组件列表
// @Summary 获取组件列表
// @Description 获取指定站点的所有组件
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Success 200 {array} models.Component
// @Router /sites/{siteId}/components [get]
func ListComponents(c *gin.Context) {
	siteID := c.Param("siteId")
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

	components, err := service.ListComponents(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, components)
}

// GetComponent 获取组件
// @Summary 获取组件
// @Description 获取指定站点的指定组件
// @Tags Component
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param id path string true "组件ID"
// @Success 200 {object} models.Component
// @Router /sites/{siteId}/components/{id} [get]
func GetComponent(c *gin.Context) {
	siteID := c.Param("siteId")
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

	component, err := service.GetComponent(siteID, componentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, component)
}

// CreateCategory 创建组件分类
// @Summary 创建组件分类
// @Description 创建一个新的组件分类
// @Tags Component
// @Accept json
// @Produce json
// @Param name path string true "分类名称"
// @Success 201 {object} models.ComponentCategory
// @Router /component-categories/{name} [post]
func CreateCategory(c *gin.Context) {
	categoryName := c.Param("name")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 创建组件分类
	createdCategory, err := service.CreateCategory(categoryName, tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// ListCategories 获取组件分类列表
// @Summary 获取组件分类列表
// @Description 获取所有可用的组件分类列表
// @Tags Component
// @Accept json
// @Produce json
// @Success 200 {array} models.ComponentCategory
// @Router /component-categories [get]
func ListCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	categories, err := service.ListCategories(tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
