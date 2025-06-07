package learn

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// TeacherHandler 处理讲师相关的HTTP请求
type TeacherHandler struct {
	teacherAppService *learn.TeacherAppService
}

// NewTeacherHandler 创建讲师处理器
func NewTeacherHandler(teacherAppService *learn.TeacherAppService) *TeacherHandler {
	return &TeacherHandler{
		teacherAppService: teacherAppService,
	}
}

// RegisterRoutes 注册路由
func (h *TeacherHandler) RegisterRoutes(router *gin.RouterGroup) {
	teacherRouter := router.Group("/teachers")
	{
		teacherRouter.POST("", middleware.AdminAuth(), h.createTeacher)
		teacherRouter.PUT("/:id/profile", middleware.AdminAuth(), h.updateTeacherProfile)
		teacherRouter.PUT("/:id/contact", middleware.AdminAuth(), h.updateTeacherContact)
		teacherRouter.PUT("/:id/specialties", middleware.AdminAuth(), h.setTeacherSpecialties)
		teacherRouter.PUT("/:id/social-profiles", middleware.AdminAuth(), h.setTeacherSocialProfiles)
		teacherRouter.PUT("/:id/activate", middleware.AdminAuth(), h.activateTeacher)
		teacherRouter.PUT("/:id/deactivate", middleware.AdminAuth(), h.deactivateTeacher)
		teacherRouter.DELETE("/:id", middleware.AdminAuth(), h.deleteTeacher)
		teacherRouter.GET("/:id", h.getTeacherDetail)
		teacherRouter.GET("/by-user/:userId", h.getTeacherByUserID)
		teacherRouter.GET("", h.listTeachers)
		teacherRouter.GET("/:id/courses", h.getTeacherWithCourses)
		teacherRouter.GET("/stats", middleware.AdminAuth(), h.getTeacherStats)
	}
}

// 创建讲师请求
type createTeacherRequest struct {
	UserID       string `json:"userId" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
}

// 创建讲师
func (h *TeacherHandler) createTeacher(c *gin.Context) {
	var req createTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherAppService.CreateTeacher(c.Request.Context(), req.UserID, req.Name, req.Title, req.Introduction, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

// 更新讲师基本资料请求
type updateTeacherProfileRequest struct {
	Name         string `json:"name" binding:"required"`
	Avatar       string `json:"avatar"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

// 更新讲师基本资料
func (h *TeacherHandler) updateTeacherProfile(c *gin.Context) {
	id := c.Param("id")
	var req updateTeacherProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherAppService.UpdateTeacherProfile(c.Request.Context(), id, req.Name, req.Avatar, req.Title, req.Introduction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 更新讲师联系信息请求
type updateTeacherContactRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// 更新讲师联系信息
func (h *TeacherHandler) updateTeacherContact(c *gin.Context) {
	id := c.Param("id")
	var req updateTeacherContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherAppService.UpdateTeacherContact(c.Request.Context(), id, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 设置讲师专长请求
type setTeacherSpecialtiesRequest struct {
	Specialties []string `json:"specialties" binding:"required"`
}

// 设置讲师专长
func (h *TeacherHandler) setTeacherSpecialties(c *gin.Context) {
	id := c.Param("id")
	var req setTeacherSpecialtiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherAppService.SetTeacherSpecialties(c.Request.Context(), id, req.Specialties)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 设置讲师社交档案请求
type setTeacherSocialProfilesRequest struct {
	SocialProfiles []string `json:"socialProfiles" binding:"required"`
}

// 设置讲师社交档案
func (h *TeacherHandler) setTeacherSocialProfiles(c *gin.Context) {
	id := c.Param("id")
	var req setTeacherSocialProfilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherAppService.SetTeacherSocialProfiles(c.Request.Context(), id, req.SocialProfiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 激活讲师
func (h *TeacherHandler) activateTeacher(c *gin.Context) {
	id := c.Param("id")
	teacher, err := h.teacherAppService.ActivateTeacher(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 停用讲师
func (h *TeacherHandler) deactivateTeacher(c *gin.Context) {
	id := c.Param("id")
	teacher, err := h.teacherAppService.DeactivateTeacher(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 删除讲师
func (h *TeacherHandler) deleteTeacher(c *gin.Context) {
	id := c.Param("id")
	err := h.teacherAppService.DeleteTeacher(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "讲师已删除"})
}

// 获取讲师详情
func (h *TeacherHandler) getTeacherDetail(c *gin.Context) {
	id := c.Param("id")
	teacher, err := h.teacherAppService.GetTeacherDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 通过用户ID获取讲师
func (h *TeacherHandler) getTeacherByUserID(c *gin.Context) {
	userID := c.Param("userId")
	teacher, err := h.teacherAppService.GetTeacherByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// 获取讲师列表
func (h *TeacherHandler) listTeachers(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤参数
	params := repository.TeacherQueryParams{
		Status:  c.Query("status"),
		Keyword: c.Query("keyword"),
	}

	teachers, total, err := h.teacherAppService.ListTeachers(c.Request.Context(), params, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": teachers,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取讲师及其课程
func (h *TeacherHandler) getTeacherWithCourses(c *gin.Context) {
	id := c.Param("id")
	teacher, courses, err := h.teacherAppService.GetTeacherWithCourses(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teacher": teacher,
		"courses": courses,
	})
}

// 获取讲师统计信息
func (h *TeacherHandler) getTeacherStats(c *gin.Context) {
	stats, err := h.teacherAppService.GetTeacherStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
