package content

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/repository"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	contentRepo repository.ContentRepository
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(contentRepo repository.ContentRepository) *CategoryHandler {
	return &CategoryHandler{
		contentRepo: contentRepo,
	}
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentId    int64  `json:"parent_id"`
		SortOrder   int32  `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &repository.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentId:    req.ParentId,
		SortOrder:   req.SortOrder,
		Status:      1, // 默认启用
	}

	id, err := h.contentRepo.CreateCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	category.ID = id
	c.JSON(http.StatusCreated, gin.H{"category": category})
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ParentId    int64  `json:"parent_id"`
		SortOrder   int32  `json:"sort_order"`
		Status      int32  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ParentId != 0 {
		updates["parent_id"] = req.ParentId
	}
	if req.SortOrder != 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	err = h.contentRepo.UpdateCategory(c.Request.Context(), categoryId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的分类信息
	category, err := h.contentRepo.GetCategoryById(c.Request.Context(), categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = h.contentRepo.DeleteCategory(c.Request.Context(), categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetCategory 获取分类详情
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.contentRepo.GetCategoryById(c.Request.Context(), categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListCategories 获取分类列表
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	filters := make(map[string]interface{})
	if parentIdStr := c.Query("parent_id"); parentIdStr != "" {
		if parentId, err := strconv.ParseInt(parentIdStr, 10, 64); err == nil {
			filters["parent_id"] = parentId
		}
	}

	categories, total, err := h.contentRepo.GetCategoryList(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
	})
} 