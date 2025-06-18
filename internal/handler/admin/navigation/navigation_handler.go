package navigation

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/zrpc"
	
	"wz-backend-go/internal/client"
)

// NavigationAdminHandler admin导航管理处理器
type NavigationAdminHandler struct {
	navigationClient client.NavigationClient
}

// NewNavigationAdminHandler 创建导航管理处理器
func NewNavigationAdminHandler(navigationClient client.NavigationClient) *NavigationAdminHandler {
	return &NavigationAdminHandler{
		navigationClient: navigationClient,
	}
}

// GetMainNavigation 获取主导航
func (h *NavigationAdminHandler) GetMainNavigation(c *gin.Context) {
	ctx := c.Request.Context()
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 调用navigation微服务
	resp, err := h.navigationClient.GetCategories(ctx, &client.GetCategoriesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取主导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": resp,
	})
}

// SaveMainNavigation 保存主导航
func (h *NavigationAdminHandler) SaveMainNavigation(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		SortOrder   int32  `json:"sortOrder"`
		IsActive    bool   `json:"isActive"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"details": err.Error(),
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务创建分类
	resp, err := h.navigationClient.CreateCategory(ctx, &client.CreateCategoryRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IconUrl:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存主导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "保存成功",
		"data": resp,
	})
}

// DeleteMainNavigationItem 删除主导航项
func (h *NavigationAdminHandler) DeleteMainNavigationItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID参数不能为空",
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务删除分类
	_, err := h.navigationClient.DeleteCategory(ctx, &client.DeleteCategoryRequest{
		Id: id,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除主导航项失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetFooterNavigation 获取底部导航
func (h *NavigationAdminHandler) GetFooterNavigation(c *gin.Context) {
	ctx := c.Request.Context()
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 调用navigation微服务获取网站列表（底部导航通常是网站链接）
	resp, err := h.navigationClient.GetWebsites(ctx, &client.GetWebsitesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Type:     "footer", // 标记为底部导航类型
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取底部导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": resp,
	})
}

// SaveFooterNavigation 保存底部导航
func (h *NavigationAdminHandler) SaveFooterNavigation(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		URL         string `json:"url" binding:"required"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		SortOrder   int32  `json:"sortOrder"`
		IsActive    bool   `json:"isActive"`
		CategoryID  string `json:"categoryId"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"details": err.Error(),
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务创建网站
	resp, err := h.navigationClient.CreateWebsite(ctx, &client.CreateWebsiteRequest{
		Name:        req.Name,
		Url:         req.URL,
		Description: req.Description,
		IconUrl:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		CategoryId:  req.CategoryID,
		Type:        "footer",
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存底部导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "保存成功",
		"data": resp,
	})
}

// DeleteFooterNavigationItem 删除底部导航项
func (h *NavigationAdminHandler) DeleteFooterNavigationItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID参数不能为空",
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务删除网站
	_, err := h.navigationClient.DeleteWebsite(ctx, &client.DeleteWebsiteRequest{
		Id: id,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除底部导航项失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetSideNavigation 获取侧边导航
func (h *NavigationAdminHandler) GetSideNavigation(c *gin.Context) {
	ctx := c.Request.Context()
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 调用navigation微服务获取网站列表（侧边导航）
	resp, err := h.navigationClient.GetWebsites(ctx, &client.GetWebsitesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Type:     "sidebar", // 标记为侧边导航类型
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取侧边导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": resp,
	})
}

// SaveSideNavigation 保存侧边导航
func (h *NavigationAdminHandler) SaveSideNavigation(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		URL         string `json:"url" binding:"required"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		SortOrder   int32  `json:"sortOrder"`
		IsActive    bool   `json:"isActive"`
		CategoryID  string `json:"categoryId"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"details": err.Error(),
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务创建网站
	resp, err := h.navigationClient.CreateWebsite(ctx, &client.CreateWebsiteRequest{
		Name:        req.Name,
		Url:         req.URL,
		Description: req.Description,
		IconUrl:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		CategoryId:  req.CategoryID,
		Type:        "sidebar",
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存侧边导航失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "保存成功",
		"data": resp,
	})
}

// DeleteSideNavigationItem 删除侧边导航项
func (h *NavigationAdminHandler) DeleteSideNavigationItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID参数不能为空",
		})
		return
	}
	
	ctx := c.Request.Context()
	
	// 调用navigation微服务删除网站
	_, err := h.navigationClient.DeleteWebsite(ctx, &client.DeleteWebsiteRequest{
		Id: id,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除侧边导航项失败",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
} 