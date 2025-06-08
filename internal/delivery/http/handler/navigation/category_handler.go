package navigation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/application/navigation"
	"wz-backend-go/internal/domain/navigation/dto"
)

// CategoryHandler handles HTTP requests for navigation categories
type CategoryHandler struct {
	navigationAppService *navigation.NavigationAppService
}

// NewCategoryHandler creates a new instance of CategoryHandler
func NewCategoryHandler(navigationAppService *navigation.NavigationAppService) *CategoryHandler {
	return &CategoryHandler{
		navigationAppService: navigationAppService,
	}
}

// CreateCategory handles requests to create a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.navigationAppService.CreateCategory(c.Request.Context(), &req)
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

	category, err := h.navigationAppService.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategories handles requests to retrieve all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	categories, err := h.navigationAppService.GetCategories(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID handles requests to retrieve a category by its ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.navigationAppService.GetCategoryByID(c.Request.Context(), id)
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

// ReorderCategories handles requests to reorder categories
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	var req dto.ReorderCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.navigationAppService.ReorderCategories(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteCategory handles requests to delete a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	// In a real implementation, you would have a delete method in the app service
	// For now, we'll return a not implemented error
	c.JSON(http.StatusNotImplemented, gin.H{"error": "delete operation not implemented"})
}
