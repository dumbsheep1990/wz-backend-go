package handlers

import (
	"net/http"
	"time"
	"wz-backend-go/models"
	"wz-backend-go/services/page-service/service"

	"github.com/gin-gonic/gin"
)

// ListPages 获取站点下的所有页面
// @Summary 获取站点下的所有页面
// @Description 获取站点下的所有页面
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Success 200 {array} models.Page
// @Router /pages [get]
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
// @Summary 获取单个页面
// @Description 获取单个页面
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param id path string true "页面ID"
// @Success 200 {object} models.Page
// @Router /pages/{siteId}/{id} [get]
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
// @Summary 创建页面
// @Description 创建一个新的页面
// @Tags Page
// @Accept json
// @Produce json
// @Param data body models.Page true "页面信息"
// @Success 201 {object} models.Page
// @Router /pages [post]
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
// @Summary 更新页面
// @Description 更新页面的信息
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param id path string true "页面ID"
// @Param data body models.Page true "页面信息"
// @Success 200 {object} models.Page
// @Router /pages/{siteId}/{id} [put]
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
// @Summary 删除页面
// @Description 删除指定的页面
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param id path string true "页面ID"
// @Success 200 {object} gin.H
// @Router /pages/{siteId}/{id} [delete]
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
// @Summary 设置页面为首页
// @Description 将指定的页面设置为首页
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param id path string true "页面ID"
// @Success 200 {object} models.Page
// @Router /pages/{siteId}/{id}/set-homepage [put]
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
// @Summary 重新排序页面
// @Description 重新排序页面
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Param data body []string true "页面顺序"
// @Success 200 {object} gin.H
// @Router /pages/{siteId}/reorder [put]
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

// CreateCategory 创建新分类
// @Summary 创建新分类
// @Description 创建一个新的分类
// @Tags Page
// @Accept json
// @Produce json
// @Param data body models.PageCategory true "分类信息"
// @Success 201 {object} models.PageCategory
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
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

	var category models.PageCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置分类属性
	category.SiteID = siteID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	createdCategory, err := service.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// ListCategories 获取站点下的所有分类
// @Summary 获取站点下的所有分类
// @Description 获取站点下的所有分类
// @Tags Page
// @Accept json
// @Produce json
// @Param siteId path string true "站点ID"
// @Success 200 {array} models.PageCategory
// @Router /categories [get]
func ListCategories(c *gin.Context) {
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

	categories, err := service.ListCategories(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
