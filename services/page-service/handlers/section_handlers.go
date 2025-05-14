package handlers

import (
	"net/http"
	"wz-backend-go/models"
	"wz-backend-go/services/page-service/service"

	"github.com/gin-gonic/gin"
)

// ListSections 获取页面下的所有区块
func ListSections(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
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

	sections, err := service.ListSections(siteID, pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// AddSection 添加新区块
func AddSection(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
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

	var section models.Section
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置区块属性
	section.PageID = pageID

	// 添加区块
	addedSection, err := service.AddSection(siteID, pageID, section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, addedSection)
}

// UpdateSection 更新区块
func UpdateSection(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("id")
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

	var section models.Section
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置区块属性
	section.ID = sectionID
	section.PageID = pageID

	updatedSection, err := service.UpdateSection(siteID, pageID, section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 同时更新页面的更新时间
	service.UpdatePageTimestamp(siteID, pageID)

	c.JSON(http.StatusOK, updatedSection)
}

// DeleteSection 删除区块
func DeleteSection(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	sectionID := c.Param("id")
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

	// 删除区块
	if err := service.DeleteSection(siteID, pageID, sectionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 同时更新页面的更新时间
	service.UpdatePageTimestamp(siteID, pageID)

	c.JSON(http.StatusOK, gin.H{"message": "区块已删除"})
}

// ReorderSections 重新排序区块
func ReorderSections(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
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

	var sectionOrder []string
	if err := c.ShouldBindJSON(&sectionOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新排序区块
	if err := service.ReorderSections(siteID, pageID, sectionOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 同时更新页面的更新时间
	service.UpdatePageTimestamp(siteID, pageID)

	c.JSON(http.StatusOK, gin.H{"message": "区块顺序已更新"})
}
