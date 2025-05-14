package handlers

import (
	"net/http"
	"time"
	"wz-backend-go/models"
	"wz-backend-go/services/page-service/service"

	"github.com/gin-gonic/gin"
)

// ListPages 获取站点下的所有页面
func ListPages(c *gin.Context) {
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

	pages, err := service.ListPages(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetPage 获取单个页面
func GetPage(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("id")
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

	page, err := service.GetPage(siteID, pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "页面不存在"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage 创建新页面
func CreatePage(c *gin.Context) {
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

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置页面属性
	page.SiteID = siteID
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	// 处理首页设置
	if page.IsHomepage {
		if err := service.UnsetOtherHomepages(siteID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	createdPage, err := service.CreatePage(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPage)
}

// UpdatePage 更新页面
func UpdatePage(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("id")
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

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置页面属性
	page.ID = pageID
	page.SiteID = siteID
	page.UpdatedAt = time.Now()

	// 处理首页设置
	if page.IsHomepage {
		if err := service.UnsetOtherHomepages(siteID, pageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	updatedPage, err := service.UpdatePage(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPage)
}

// DeletePage 删除页面
func DeletePage(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("id")
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

	// 检查是否为首页
	page, err := service.GetPage(siteID, pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "页面不存在"})
		return
	}

	if page.IsHomepage {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除首页，请先设置其他页面为首页"})
		return
	}

	if err := service.DeletePage(siteID, pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "页面已删除"})
}

// SetHomepage 设置页面为首页
func SetHomepage(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("id")
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

	// 将其他页面设为非首页
	if err := service.UnsetOtherHomepages(siteID, pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置当前页面为首页
	updatedPage, err := service.SetHomepage(siteID, pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPage)
}

// ReorderPages 重新排序页面
func ReorderPages(c *gin.Context) {
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

	var pageOrder []string
	if err := c.ShouldBindJSON(&pageOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.ReorderPages(siteID, pageOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "页面顺序已更新"})
}
