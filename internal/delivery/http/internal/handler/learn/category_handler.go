package learn

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// CategoryHandler 处理分类相关的HTTP请求
type CategoryHandler struct {
	categoryAppService *learn.CategoryAppService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryAppService *learn.CategoryAppService) *CategoryHandler {
	return &CategoryHandler{
		categoryAppService: categoryAppService,
	}
}

// RegisterRoutes 注册路由
func (h *CategoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	categoryRouter := router.Group("/categories")
	{
		// 管理员接口
		categoryRouter.POST("", middleware.AdminAuth(), h.createCategory)
		categoryRouter.PUT("/:id", middleware.AdminAuth(), h.updateCategory)
		categoryRouter.PUT("/:id/activate", middleware.AdminAuth(), h.activateCategory)
		categoryRouter.PUT("/:id/deactivate", middleware.AdminAuth(), h.deactivateCategory)
		categoryRouter.DELETE("/:id", middleware.AdminAuth(), h.deleteCategory)
		categoryRouter.GET("/stats", middleware.AdminAuth(), h.getCategoryStats)
		
		// 公共接口
		categoryRouter.GET("/:id", h.getCategoryByID)
		categoryRouter.GET("/tree", h.getCategoryTree)
		categoryRouter.GET("/active", h.getActiveCategories)
		categoryRouter.GET("/with-courses", h.getCategoriesWithCourseCount)
	}
}

// 创建分类请求
type createCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parentId"`
	Order       int     `json:"order"`
}

// 创建分类
func (h *CategoryHandler) createCategory(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	category, err := h.categoryAppService.CreateCategory(
		c.Request.Context(), req.Name, req.Description, req.Icon, req.ParentID, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, category)
}

// 更新分类请求
type updateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

// 更新分类
func (h *CategoryHandler) updateCategory(c *gin.Context) {
	id := c.Param("id")
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	category, err := h.categoryAppService.UpdateCategory(
		c.Request.Context(), id, req.Name, req.Description, req.Icon, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// 激活分类
func (h *CategoryHandler) activateCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryAppService.ActivateCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// 停用分类
func (h *CategoryHandler) deactivateCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryAppService.DeactivateCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// 删除分类
func (h *CategoryHandler) deleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := h.categoryAppService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}

// 获取分类详情
func (h *CategoryHandler) getCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryAppService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// 获取分类树
func (h *CategoryHandler) getCategoryTree(c *gin.Context) {
	tree, err := h.categoryAppService.GetCategoryTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tree)
}

// 获取所有激活状态的分类
func (h *CategoryHandler) getActiveCategories(c *gin.Context) {
	categories, err := h.categoryAppService.GetActiveCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, categories)
}

// 获取带课程数的分类列表
func (h *CategoryHandler) getCategoriesWithCourseCount(c *gin.Context) {
	categories, err := h.categoryAppService.GetCategoriesWithCourseCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, categories)
}

// 获取分类统计信息
func (h *CategoryHandler) getCategoryStats(c *gin.Context) {
	stats, err := h.categoryAppService.GetCategoryStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}
