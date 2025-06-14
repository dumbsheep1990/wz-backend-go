package handler

import (
	"net/http"
	"strconv"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 评价HTTP处理器
type ReviewHandler struct {
	reviewAppService *learn.ReviewAppService
}

// NewReviewHandler 创建评价处理器
func NewReviewHandler(reviewAppService *learn.ReviewAppService) *ReviewHandler {
	return &ReviewHandler{
		reviewAppService: reviewAppService,
	}
}

// CreateReviewRequest 创建评价请求
type CreateReviewRequest struct {
	CourseID string `json:"courseId" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Content  string `json:"content" binding:"required,max=1000"`
}

// UpdateReviewRequest 更新评价请求
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"required,max=1000"`
}

// CreateReview 创建评价
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID（假设已通过认证中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	review, err := h.reviewAppService.CreateReview(c.Request.Context(), userID.(string), req.CourseID, req.Rating, req.Content)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "创建评价失败", err.Error())
		return
	}

	response.Success(c, review)
}

// UpdateReview 更新评价
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		response.Error(c, http.StatusBadRequest, "评价ID不能为空", "")
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	review, err := h.reviewAppService.UpdateReview(c.Request.Context(), userID.(string), reviewID, req.Rating, req.Content)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "更新评价失败", err.Error())
		return
	}

	response.Success(c, review)
}

// DeleteReview 删除评价
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		response.Error(c, http.StatusBadRequest, "评价ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	err := h.reviewAppService.DeleteReview(c.Request.Context(), userID.(string), reviewID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "删除评价失败", err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetMyReviews 获取我的评价
func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	reviews, total, err := h.reviewAppService.GetReviewsByUser(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评价列表失败", err.Error())
		return
	}

	response.SuccessWithPagination(c, reviews, total, page, pageSize)
}

// GetCourseReviews 获取课程评价
func (h *ReviewHandler) GetCourseReviews(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "课程ID不能为空", "")
		return
	}

	statusStr := c.DefaultQuery("status", string(entity.ReviewStatusApproved))
	status := entity.ReviewStatus(statusStr)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	reviews, total, err := h.reviewAppService.GetReviewsByCourse(c.Request.Context(), courseID, status, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取课程评价失败", err.Error())
		return
	}

	response.SuccessWithPagination(c, reviews, total, page, pageSize)
}

// GetRatingStats 获取课程评分统计
func (h *ReviewHandler) GetRatingStats(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "课程ID不能为空", "")
		return
	}

	stats, err := h.reviewAppService.GetCourseRatingStats(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评分统计失败", err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPendingReviews 获取待审核评价（管理员）
func (h *ReviewHandler) GetPendingReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	reviews, total, err := h.reviewAppService.GetPendingReviews(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待审核评价失败", err.Error())
		return
	}

	response.SuccessWithPagination(c, reviews, total, page, pageSize)
}

// ApproveReview 审核通过评价（管理员）
func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		response.Error(c, http.StatusBadRequest, "评价ID不能为空", "")
		return
	}

	review, err := h.reviewAppService.ApproveReview(c.Request.Context(), reviewID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "审核通过失败", err.Error())
		return
	}

	response.Success(c, review)
}

// RejectReview 审核拒绝评价（管理员）
func (h *ReviewHandler) RejectReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		response.Error(c, http.StatusBadRequest, "评价ID不能为空", "")
		return
	}

	review, err := h.reviewAppService.RejectReview(c.Request.Context(), reviewID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "审核拒绝失败", err.Error())
		return
	}

	response.Success(c, review)
}
