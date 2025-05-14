package handlers

import (
	"net/http"
	"wz-backend-go/services/render-service/service"

	"github.com/gin-gonic/gin"
)

// PreviewSite 预览整个站点
func PreviewSite(c *gin.Context) {
	siteID := c.Param("siteId")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 模拟设备类型
	device := c.DefaultQuery("device", "desktop") // desktop, tablet, mobile

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	// 获取站点数据
	site, err := service.GetSiteWithAllPages(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "站点不存在"})
		return
	}

	// 生成站点预览HTML
	htmlContent, err := service.GenerateSitePreview(site, device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回HTML内容
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// PreviewPage 预览单个页面
func PreviewPage(c *gin.Context) {
	siteID := c.Param("siteId")
	pageID := c.Param("pageId")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}

	// 模拟设备类型
	device := c.DefaultQuery("device", "desktop") // desktop, tablet, mobile

	// 校验站点所有权
	if !service.CheckSiteAccess(siteID, tenantID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该站点"})
		return
	}

	// 获取站点和页面数据
	site, page, err := service.GetSiteAndPage(siteID, pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 生成页面预览HTML
	htmlContent, err := service.GeneratePagePreview(site, page, device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回HTML内容
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}
