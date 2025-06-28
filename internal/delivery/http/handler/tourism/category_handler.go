package tourism

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/wanzhitouzi/wz-backend-go/internal/application/tourism"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/handler"
)

// CategoryHandler handles HTTP requests for tourism categories
type CategoryHandler struct {
	categoryAppService *tourism.CategoryAppService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(categoryAppService *tourism.CategoryAppService) *CategoryHandler {
	return &CategoryHandler{
		categoryAppService: categoryAppService,
	}
}

// Create handles the creation of a new category
// @Summary Create a new category
// @Description Create a new tourism category
// @Tags Tourism - Categories
// @Accept json
// @Produce json
// @Param request body dto.CategoryCreateRequest true "Category creation data"
// @Success 201 {object} handler.Response{data=dto.CategoryResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/categories [post]
// @Security ApiKeyAuth
func (h *CategoryHandler) Create(c *gin.Context) {
	// Parse request body
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create category
	resp, err := h.categoryAppService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to create category: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusCreated, "Category created successfully", resp)
}

// Get handles the retrieval of a category by ID
// @Summary Get a category
// @Description Get detailed information about a specific category
// @Tags Tourism - Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} handler.Response{data=dto.CategoryResponse}
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	// Get category ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Get category
	resp, err := h.categoryAppService.GetCategory(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Failed to get category: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Category retrieved successfully", resp)
}

// Update handles the update of a category
// @Summary Update a category
// @Description Update a category with the provided information
// @Tags Tourism - Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.CategoryUpdateRequest true "Category update data"
// @Success 200 {object} handler.Response{data=dto.CategoryResponse}
// @Failure 400 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/categories/{id} [put]
// @Security ApiKeyAuth
func (h *CategoryHandler) Update(c *gin.Context) {
	// Get category ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Parse request body
	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Update category
	resp, err := h.categoryAppService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to update category: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Category updated successfully", resp)
}

// Delete handles the deletion of a category
// @Summary Delete a category
// @Description Delete a specific category
// @Tags Tourism - Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/categories/{id} [delete]
// @Security ApiKeyAuth
func (h *CategoryHandler) Delete(c *gin.Context) {
	// Get category ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Delete category
	err := h.categoryAppService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to delete category: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Category deleted successfully", nil)
}

// List handles the listing of all categories
// @Summary List categories
// @Description Get a paginated list of categories
// @Tags Tourism - Categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.CategoryListResponse}
// @Failure 500 {object} handler.Response
// @Router /tourism/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List categories
	resp, err := h.categoryAppService.ListCategories(c.Request.Context(), offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list categories: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Categories retrieved successfully", resp)
}

// GetTree handles the retrieval of the category tree
// @Summary Get category tree
// @Description Get the hierarchical category tree structure
// @Tags Tourism - Categories
// @Produce json
// @Success 200 {object} handler.Response{data=[]dto.CategoryTreeNode}
// @Failure 500 {object} handler.Response
// @Router /tourism/categories/tree [get]
func (h *CategoryHandler) GetTree(c *gin.Context) {
	// Get category tree
	tree, err := h.categoryAppService.GetCategoryTree(c.Request.Context())
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to get category tree: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Category tree retrieved successfully", tree)
}
