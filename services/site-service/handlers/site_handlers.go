package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"wz-backend-go/models"
	"wz-backend-go/services/site-service/service"
)

// ListSites 获取站点列表
func ListSites(c *gin.Context) {
	// 从认证中间件获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	// 获取查询参数
	status := c.Query("status")
	search := c.Query("search")
	
	// 调用服务层获取站点列表
	sites, err := service.ListSites(tenantID.(string), status, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, sites)
}

// GetSite 获取单个站点
func GetSite(c *gin.Context) {
	siteID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	site, err := service.GetSite(siteID, tenantID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "站点不存在"})
		return
	}
	
	c.JSON(http.StatusOK, site)
}

// CreateSite 创建新站点
func CreateSite(c *gin.Context) {
	var site models.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	// 设置默认值
	site.TenantID = tenantID.(string)
	site.Status = "draft"
	site.CreatedAt = time.Now()
	site.UpdatedAt = time.Now()
	
	// 调用服务层创建站点
	createdSite, err := service.CreateSite(site)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, createdSite)
}

// UpdateSite 更新站点
func UpdateSite(c *gin.Context) {
	siteID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	var site models.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 设置站点ID和更新时间
	site.ID = siteID
	site.UpdatedAt = time.Now()
	
	// 检查站点所有权
	exists, err := service.CheckSiteOwnership(siteID, tenantID.(string))
	if err != nil || !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权更新此站点"})
		return
	}
	
	// 调用服务层更新站点
	updatedSite, err := service.UpdateSite(site)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedSite)
}

// DeleteSite 删除站点
func DeleteSite(c *gin.Context) {
	siteID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	// 检查站点所有权
	exists, err := service.CheckSiteOwnership(siteID, tenantID.(string))
	if err != nil || !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此站点"})
		return
	}
	
	// 调用服务层删除站点
	if err := service.DeleteSite(siteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "站点已删除"})
}

// PublishSite 发布站点
func PublishSite(c *gin.Context) {
	siteID := c.Param("id")
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
		return
	}
	
	// 检查站点所有权
	exists, err := service.CheckSiteOwnership(siteID, tenantID.(string))
	if err != nil || !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权发布此站点"})
		return
	}
	
	// 调用服务层发布站点
	publishedSite, err := service.PublishSite(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, publishedSite)
} 