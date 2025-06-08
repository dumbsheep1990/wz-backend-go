package commerce

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	"wz-backend-go/internal/application/commerce"
	"wz-backend-go/internal/domain/commerce/dto"
)

// CategoryHandler handles HTTP requests for commerce categories
type CategoryHandler struct {
	categoryAppService *commerce.CategoryAppService
}

// NewCategoryHandler creates a new instance of CategoryHandler
func NewCategoryHandler(categoryAppService *commerce.CategoryAppService) *CategoryHandler {
	return &CategoryHandler{
		categoryAppService: categoryAppService,
	}
}

// CreateCategory handles requests to create a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryAppService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory handles requests to update an existing category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the one in the request body
	if c.Param("id") != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID mismatch"})
		return
	}

	category, err := h.categoryAppService.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryByID handles requests to retrieve a category by its ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	
	includeChildren := false
	if includeChildrenStr := c.Query("includeChildren"); includeChildrenStr != "" {
		includeChildren, _ = strconv.ParseBool(includeChildrenStr)
	}

	category, err := h.categoryAppService.GetCategoryByID(c.Request.Context(), id, includeChildren)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetRootCategories handles requests to retrieve root categories
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	activeOnly := false
	if activeOnlyStr := c.Query("activeOnly"); activeOnlyStr != "" {
		activeOnly, _ = strconv.ParseBool(activeOnlyStr)
	}
	
	categories, err := h.categoryAppService.GetRootCategories(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetCategories handles requests to retrieve categories with filtering
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var filter dto.CategoryFilterRequest
	
	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	categories, err := h.categoryAppService.GetCategories(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// ReorderCategories handles requests to update category sort order
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	var req dto.ReorderCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryAppService.ReorderCategories(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetCategoryHierarchy handles requests to retrieve the complete category hierarchy
func (h *CategoryHandler) GetCategoryHierarchy(c *gin.Context) {
	activeOnly := false
	if activeOnlyStr := c.Query("activeOnly"); activeOnlyStr != "" {
		activeOnly, _ = strconv.ParseBool(activeOnlyStr)
	}
	
	hierarchy, err := h.categoryAppService.GetCategoryHierarchy(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"categories": hierarchy})
}
