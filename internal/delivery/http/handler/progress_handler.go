package handler

import (
	"net/http"
	"strconv"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
)

// ProgressHandler 学习进度HTTP处理器
type ProgressHandler struct {
	progressAppService *learn.ProgressAppService
}

// NewProgressHandler 创建学习进度处理器
func NewProgressHandler(progressAppService *learn.ProgressAppService) *ProgressHandler {
	return &ProgressHandler{
		progressAppService: progressAppService,
	}
}

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	WatchedDuration int `json:"watchedDuration" binding:"required,min=0"`
}

// UpdateLessonProgress 更新课时学习进度
func (h *ProgressHandler) UpdateLessonProgress(c *gin.Context) {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		response.Error(c, http.StatusBadRequest, "课时ID不能为空", "")
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	progress, err := h.progressAppService.UpdateLessonProgress(c.Request.Context(), userID.(string), lessonID, req.WatchedDuration)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "更新学习进度失败", err.Error())
		return
	}

	response.Success(c, progress)
}

// CompleteLesson 完成课时
func (h *ProgressHandler) CompleteLesson(c *gin.Context) {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		response.Error(c, http.StatusBadRequest, "课时ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	progress, err := h.progressAppService.CompleteLessonProgress(c.Request.Context(), userID.(string), lessonID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "完成课时失败", err.Error())
		return
	}

	response.Success(c, progress)
}

// ResetLesson 重置课时进度
func (h *ProgressHandler) ResetLesson(c *gin.Context) {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		response.Error(c, http.StatusBadRequest, "课时ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	progress, err := h.progressAppService.ResetLessonProgress(c.Request.Context(), userID.(string), lessonID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "重置课时进度失败", err.Error())
		return
	}

	response.Success(c, progress)
}

// GetMyProgress 获取我的学习进度
func (h *ProgressHandler) GetMyProgress(c *gin.Context) {
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

	progresses, total, err := h.progressAppService.GetUserProgress(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取学习进度失败", err.Error())
		return
	}

	response.SuccessWithPagination(c, progresses, total, page, pageSize)
}

// GetCourseProgress 获取课程学习进度
func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "课程ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	progresses, err := h.progressAppService.GetCourseProgress(c.Request.Context(), userID.(string), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取课程学习进度失败", err.Error())
		return
	}

	response.Success(c, progresses)
}

// GetRecentProgress 获取最近学习进度
func (h *ProgressHandler) GetRecentProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	progresses, err := h.progressAppService.GetRecentProgress(c.Request.Context(), userID.(string), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取最近学习进度失败", err.Error())
		return
	}

	response.Success(c, progresses)
}

// GetCourseProgressStats 获取课程进度统计
func (h *ProgressHandler) GetCourseProgressStats(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "课程ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	stats, err := h.progressAppService.GetCourseProgressStats(c.Request.Context(), userID.(string), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取课程进度统计失败", err.Error())
		return
	}

	response.Success(c, stats)
}

// GetUserProgressStats 获取用户整体进度统计
func (h *ProgressHandler) GetUserProgressStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	stats, err := h.progressAppService.GetUserOverallStats(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户进度统计失败", err.Error())
		return
	}

	response.Success(c, stats)
}

// InitializeCourseProgress 初始化课程进度
func (h *ProgressHandler) InitializeCourseProgress(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "课程ID不能为空", "")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未认证", "")
		return
	}

	err := h.progressAppService.InitializeCourseProgress(c.Request.Context(), userID.(string), courseID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "初始化课程进度失败", err.Error())
		return
	}

	response.Success(c, gin.H{"message": "初始化成功"})
}
