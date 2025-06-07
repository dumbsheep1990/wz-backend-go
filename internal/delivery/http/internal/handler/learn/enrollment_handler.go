package learn

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// EnrollmentHandler 处理课程报名相关的HTTP请求
type EnrollmentHandler struct {
	enrollmentAppService *learn.EnrollmentAppService
}

// NewEnrollmentHandler 创建课程报名处理器
func NewEnrollmentHandler(enrollmentAppService *learn.EnrollmentAppService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentAppService: enrollmentAppService,
	}
}

// RegisterRoutes 注册路由
func (h *EnrollmentHandler) RegisterRoutes(router *gin.RouterGroup) {
	enrollmentRouter := router.Group("/enrollments")
	{
		// 学员使用的接口
		enrollmentRouter.POST("", middleware.Auth(), h.enrollCourse)
		enrollmentRouter.PUT("/:id/progress", middleware.Auth(), h.updateProgress)
		enrollmentRouter.POST("/:id/complete", middleware.Auth(), h.completeCourse)
		enrollmentRouter.POST("/:id/rating", middleware.Auth(), h.rateEnrollment)
		enrollmentRouter.GET("/user/:userId", middleware.Auth(), h.listUserEnrollments)
		
		// 管理员使用的接口
		enrollmentRouter.POST("/:id/refund", middleware.AdminAuth(), h.refundEnrollment)
		enrollmentRouter.GET("", middleware.AdminAuth(), h.listEnrollments)
		enrollmentRouter.POST("/process-expired", middleware.AdminAuth(), h.processExpiredEnrollments)
		enrollmentRouter.GET("/stats", middleware.AdminAuth(), h.getEnrollmentStats)
	}
}

// 报名课程请求
type enrollCourseRequest struct {
	UserID   string `json:"userId" binding:"required"`
	CourseID string `json:"courseId" binding:"required"`
}

// 报名课程
func (h *EnrollmentHandler) enrollCourse(c *gin.Context) {
	var req enrollCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证当前登录用户与请求中的用户ID是否一致
	currentUserID := c.GetString("userId")
	if currentUserID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权为其他用户报名课程"})
		return
	}
	
	enrollment, err := h.enrollmentAppService.EnrollCourse(c.Request.Context(), req.UserID, req.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, enrollment)
}

// 更新学习进度请求
type updateProgressRequest struct {
	LessonID string `json:"lessonId" binding:"required"`
	Progress int    `json:"progress" binding:"required,min=0,max=100"`
}

// 更新学习进度
func (h *EnrollmentHandler) updateProgress(c *gin.Context) {
	id := c.Param("id")
	var req updateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证用户是否有权更新此报名记录
	enrollment, err := h.enrollmentAppService.GetEnrollment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	currentUserID := c.GetString("userId")
	if currentUserID != enrollment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权更新其他用户的学习进度"})
		return
	}
	
	updatedEnrollment, err := h.enrollmentAppService.UpdateProgress(
		c.Request.Context(), id, req.LessonID, req.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedEnrollment)
}

// 完成课程
func (h *EnrollmentHandler) completeCourse(c *gin.Context) {
	id := c.Param("id")
	
	// 验证用户是否有权完成此课程
	enrollment, err := h.enrollmentAppService.GetEnrollment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	currentUserID := c.GetString("userId")
	if currentUserID != enrollment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权完成其他用户的课程"})
		return
	}
	
	updatedEnrollment, err := h.enrollmentAppService.CompleteCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedEnrollment)
}

// 评价课程请求
type rateEnrollmentRequest struct {
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	ReviewTitle string `json:"reviewTitle"`
	ReviewBody  string `json:"reviewBody"`
}

// 评价课程
func (h *EnrollmentHandler) rateEnrollment(c *gin.Context) {
	id := c.Param("id")
	var req rateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证用户是否有权评价此课程
	enrollment, err := h.enrollmentAppService.GetEnrollment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	currentUserID := c.GetString("userId")
	if currentUserID != enrollment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权评价其他用户的课程"})
		return
	}
	
	updatedEnrollment, err := h.enrollmentAppService.RateEnrollment(
		c.Request.Context(), id, req.Rating, req.ReviewTitle, req.ReviewBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedEnrollment)
}

// 退款报名请求
type refundEnrollmentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// 退款报名
func (h *EnrollmentHandler) refundEnrollment(c *gin.Context) {
	id := c.Param("id")
	var req refundEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	enrollment, err := h.enrollmentAppService.RefundEnrollment(c.Request.Context(), id, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, enrollment)
}

// 获取用户的所有报名
func (h *EnrollmentHandler) listUserEnrollments(c *gin.Context) {
	userID := c.Param("userId")
	
	// 验证当前登录用户与请求的用户ID是否一致
	currentUserID := c.GetString("userId")
	if currentUserID != userID && !c.GetBool("isAdmin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看其他用户的报名记录"})
		return
	}
	
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 解析状态过滤参数
	status := c.Query("status") // 可能的值: active, completed, expired, refunded
	
	enrollments, total, err := h.enrollmentAppService.ListUserEnrollments(
		c.Request.Context(), userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"items": enrollments,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取所有报名记录（管理员）
func (h *EnrollmentHandler) listEnrollments(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 解析过滤参数
	params := repository.EnrollmentQueryParams{
		Status:   c.Query("status"),
		UserID:   c.Query("userId"),
		CourseID: c.Query("courseId"),
		FromDate: c.Query("fromDate"),
		ToDate:   c.Query("toDate"),
	}
	
	enrollments, total, err := h.enrollmentAppService.ListEnrollments(
		c.Request.Context(), params, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"items": enrollments,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 处理过期的报名记录
func (h *EnrollmentHandler) processExpiredEnrollments(c *gin.Context) {
	count, err := h.enrollmentAppService.ProcessExpiredEnrollments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "成功处理过期报名记录",
		"count":   count,
	})
}

// 获取报名统计信息
func (h *EnrollmentHandler) getEnrollmentStats(c *gin.Context) {
	stats, err := h.enrollmentAppService.GetEnrollmentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}
