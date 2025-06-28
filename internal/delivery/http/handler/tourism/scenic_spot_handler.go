package tourism

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	"github.com/wanzhitouzi/wz-backend-go/internal/application/tourism"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/handler"
)

// ScenicSpotHandler handles HTTP requests for scenic spots
type ScenicSpotHandler struct {
	scenicSpotAppService *tourism.ScenicSpotAppService
}

// NewScenicSpotHandler creates a new ScenicSpotHandler instance
func NewScenicSpotHandler(scenicSpotAppService *tourism.ScenicSpotAppService) *ScenicSpotHandler {
	return &ScenicSpotHandler{
		scenicSpotAppService: scenicSpotAppService,
	}
}

// Create handles the creation of a new scenic spot
// @Summary Create a new scenic spot
// @Description Create a new scenic spot with the provided information
// @Tags Tourism - Scenic Spots
// @Accept json
// @Produce json
// @Param request body dto.ScenicSpotCreateRequest true "Scenic spot creation data"
// @Success 201 {object} handler.Response{data=dto.ScenicSpotResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots [post]
// @Security ApiKeyAuth
func (h *ScenicSpotHandler) Create(c *gin.Context) {
	// Parse request body
	var req dto.ScenicSpotCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create scenic spot
	resp, err := h.scenicSpotAppService.CreateScenicSpot(c.Request.Context(), &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to create scenic spot: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusCreated, "Scenic spot created successfully", resp)
}

// Get handles the retrieval of a scenic spot by ID
// @Summary Get a scenic spot
// @Description Get detailed information about a specific scenic spot
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param id path string true "Scenic Spot ID"
// @Success 200 {object} handler.Response{data=dto.ScenicSpotResponse}
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots/{id} [get]
func (h *ScenicSpotHandler) Get(c *gin.Context) {
	// Get scenic spot ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Scenic spot ID is required")
		return
	}

	// Get scenic spot
	resp, err := h.scenicSpotAppService.GetScenicSpot(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Failed to get scenic spot: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spot retrieved successfully", resp)
}

// Update handles the update of a scenic spot
// @Summary Update a scenic spot
// @Description Update a scenic spot with the provided information
// @Tags Tourism - Scenic Spots
// @Accept json
// @Produce json
// @Param id path string true "Scenic Spot ID"
// @Param request body dto.ScenicSpotUpdateRequest true "Scenic spot update data"
// @Success 200 {object} handler.Response{data=dto.ScenicSpotResponse}
// @Failure 400 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots/{id} [put]
// @Security ApiKeyAuth
func (h *ScenicSpotHandler) Update(c *gin.Context) {
	// Get scenic spot ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Scenic spot ID is required")
		return
	}

	// Parse request body
	var req dto.ScenicSpotUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Update scenic spot
	resp, err := h.scenicSpotAppService.UpdateScenicSpot(c.Request.Context(), id, &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to update scenic spot: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spot updated successfully", resp)
}

// Delete handles the deletion of a scenic spot
// @Summary Delete a scenic spot
// @Description Delete a specific scenic spot
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param id path string true "Scenic Spot ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots/{id} [delete]
// @Security ApiKeyAuth
func (h *ScenicSpotHandler) Delete(c *gin.Context) {
	// Get scenic spot ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Scenic spot ID is required")
		return
	}

	// Delete scenic spot
	err := h.scenicSpotAppService.DeleteScenicSpot(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to delete scenic spot: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spot deleted successfully", nil)
}

// List handles the listing of all scenic spots
// @Summary List scenic spots
// @Description Get a paginated list of scenic spots
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ScenicSpotListResponse}
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots [get]
func (h *ScenicSpotHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List scenic spots
	resp, err := h.scenicSpotAppService.ListScenicSpots(c.Request.Context(), offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list scenic spots: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spots retrieved successfully", resp)
}

// ListByCategory handles the listing of scenic spots by category
// @Summary List scenic spots by category
// @Description Get a paginated list of scenic spots filtered by category
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param category_id path string true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ScenicSpotListResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/categories/{category_id}/scenic-spots [get]
func (h *ScenicSpotHandler) ListByCategory(c *gin.Context) {
	// Get category ID from path
	categoryID := c.Param("category_id")
	if categoryID == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List scenic spots by category
	resp, err := h.scenicSpotAppService.ListScenicSpotsByCategory(c.Request.Context(), categoryID, offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list scenic spots: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spots retrieved successfully", resp)
}

// ListByArea handles the listing of scenic spots by area
// @Summary List scenic spots by area
// @Description Get a paginated list of scenic spots filtered by area
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param area path string true "Location Area"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ScenicSpotListResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/areas/{area}/scenic-spots [get]
func (h *ScenicSpotHandler) ListByArea(c *gin.Context) {
	// Get area from path
	area := c.Param("area")
	if area == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Area is required")
		return
	}

	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List scenic spots by area
	resp, err := h.scenicSpotAppService.ListScenicSpotsByArea(c.Request.Context(), area, offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list scenic spots: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spots retrieved successfully", resp)
}

// Search handles the search for scenic spots
// @Summary Search scenic spots
// @Description Search for scenic spots by keyword
// @Tags Tourism - Scenic Spots
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ScenicSpotListResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots/search [get]
func (h *ScenicSpotHandler) Search(c *gin.Context) {
	// Get keyword from query
	keyword := c.Query("keyword")
	if keyword == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Search keyword is required")
		return
	}

	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// Search scenic spots
	resp, err := h.scenicSpotAppService.SearchScenicSpots(c.Request.Context(), keyword, offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to search scenic spots: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Scenic spots retrieved successfully", resp)
}

// Helper methods

// getPaginationParams extracts and validates pagination parameters from the request
func getPaginationParams(c *gin.Context) (int, int) {
	// Parse page parameter
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}
