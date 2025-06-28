package tourism

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/wanzhitouzi/wz-backend-go/internal/application/tourism"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/handler"
	"github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/middleware"
)

// ReviewHandler handles HTTP requests for scenic spot reviews
type ReviewHandler struct {
	reviewAppService *tourism.ReviewAppService
}

// NewReviewHandler creates a new ReviewHandler instance
func NewReviewHandler(reviewAppService *tourism.ReviewAppService) *ReviewHandler {
	return &ReviewHandler{
		reviewAppService: reviewAppService,
	}
}

// Create handles the creation of a new review for a scenic spot
// @Summary Create a new review
// @Description Create a new review for a scenic spot
// @Tags Tourism - Reviews
// @Accept json
// @Produce json
// @Param request body dto.ReviewCreateRequest true "Review creation data"
// @Success 201 {object} handler.Response{data=dto.ReviewResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews [post]
// @Security ApiKeyAuth
func (h *ReviewHandler) Create(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, userName := middleware.GetUserIDAndName(c)
	if userID == "" {
		handler.ResponseError(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	// Parse request body
	var req dto.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create review
	resp, err := h.reviewAppService.CreateReview(c.Request.Context(), userID, userName, &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to create review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusCreated, "Review created successfully", resp)
}

// Get handles the retrieval of a review by ID
// @Summary Get a review
// @Description Get detailed information about a specific review
// @Tags Tourism - Reviews
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} handler.Response{data=dto.ReviewResponse}
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews/{id} [get]
func (h *ReviewHandler) Get(c *gin.Context) {
	// Get review ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Review ID is required")
		return
	}

	// Get review
	resp, err := h.reviewAppService.GetReview(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Failed to get review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Review retrieved successfully", resp)
}

// Update handles the update of a review
// @Summary Update a review
// @Description Update a review with the provided information
// @Tags Tourism - Reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Param request body dto.ReviewUpdateRequest true "Review update data"
// @Success 200 {object} handler.Response{data=dto.ReviewResponse}
// @Failure 400 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews/{id} [put]
// @Security ApiKeyAuth
func (h *ReviewHandler) Update(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, _ := middleware.GetUserIDAndName(c)
	if userID == "" {
		handler.ResponseError(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	// Get review ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Review ID is required")
		return
	}

	// Parse request body
	var req dto.ReviewUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Update review
	resp, err := h.reviewAppService.UpdateReview(c.Request.Context(), id, userID, &req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to update review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Review updated successfully", resp)
}

// Delete handles the deletion of a review
// @Summary Delete a review
// @Description Delete a specific review
// @Tags Tourism - Reviews
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews/{id} [delete]
// @Security ApiKeyAuth
func (h *ReviewHandler) Delete(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, _ := middleware.GetUserIDAndName(c)
	if userID == "" {
		handler.ResponseError(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	// Get review ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Review ID is required")
		return
	}

	// Delete review
	err := h.reviewAppService.DeleteReview(c.Request.Context(), id, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to delete review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Review deleted successfully", nil)
}

// Like handles the action of liking a review
// @Summary Like a review
// @Description Like a specific review
// @Tags Tourism - Reviews
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews/{id}/like [post]
// @Security ApiKeyAuth
func (h *ReviewHandler) Like(c *gin.Context) {
	// Check if user is authenticated
	userID, _ := middleware.GetUserIDAndName(c)
	if userID == "" {
		handler.ResponseError(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	// Get review ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Review ID is required")
		return
	}

	// Like review
	err := h.reviewAppService.LikeReview(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to like review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Review liked successfully", nil)
}

// Unlike handles the action of unliking a review
// @Summary Unlike a review
// @Description Unlike a specific review
// @Tags Tourism - Reviews
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/reviews/{id}/unlike [post]
// @Security ApiKeyAuth
func (h *ReviewHandler) Unlike(c *gin.Context) {
	// Check if user is authenticated
	userID, _ := middleware.GetUserIDAndName(c)
	if userID == "" {
		handler.ResponseError(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	// Get review ID from path
	id := c.Param("id")
	if id == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Review ID is required")
		return
	}

	// Unlike review
	err := h.reviewAppService.UnlikeReview(c.Request.Context(), id)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to unlike review: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Review unliked successfully", nil)
}

// ListByScenicSpot handles the listing of reviews for a specific scenic spot
// @Summary List reviews by scenic spot
// @Description Get a paginated list of reviews for a specific scenic spot
// @Tags Tourism - Reviews
// @Produce json
// @Param scenic_spot_id path string true "Scenic Spot ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ReviewListResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/scenic-spots/{scenic_spot_id}/reviews [get]
func (h *ReviewHandler) ListByScenicSpot(c *gin.Context) {
	// Get scenic spot ID from path
	scenicSpotID := c.Param("scenic_spot_id")
	if scenicSpotID == "" {
		handler.ResponseError(c, http.StatusBadRequest, "Scenic spot ID is required")
		return
	}

	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List reviews by scenic spot
	resp, err := h.reviewAppService.ListReviewsByScenicSpot(c.Request.Context(), scenicSpotID, offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list reviews: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Reviews retrieved successfully", resp)
}

// ListByUser handles the listing of reviews by a specific user
// @Summary List reviews by user
// @Description Get a paginated list of reviews from a specific user
// @Tags Tourism - Reviews
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=dto.ReviewListResponse}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /tourism/users/{user_id}/reviews [get]
func (h *ReviewHandler) ListByUser(c *gin.Context) {
	// Get user ID from path
	userID := c.Param("user_id")
	if userID == "" {
		handler.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Parse pagination parameters
	page, limit := getPaginationParams(c)
	offset := (page - 1) * limit

	// List reviews by user
	resp, err := h.reviewAppService.ListReviewsByUser(c.Request.Context(), userID, offset, limit)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Failed to list reviews: "+err.Error())
		return
	}

	// Return success response
	handler.ResponseSuccess(c, http.StatusOK, "Reviews retrieved successfully", resp)
}
