package learn

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// CourseHandler 处理课程相关的HTTP请求
type CourseHandler struct {
	courseAppService *learn.CourseAppService
}

// NewCourseHandler 创建课程处理器
func NewCourseHandler(courseAppService *learn.CourseAppService) *CourseHandler {
	return &CourseHandler{
		courseAppService: courseAppService,
	}
}

// RegisterRoutes 注册路由
func (h *CourseHandler) RegisterRoutes(router *gin.RouterGroup) {
	courseRouter := router.Group("/courses")
	{
		courseRouter.POST("", middleware.AdminAuth(), h.createCourse)
		courseRouter.PUT("/:id", middleware.AdminAuth(), h.updateCourse)
		courseRouter.PUT("/:id/media", middleware.AdminAuth(), h.updateCourseMedia)
		courseRouter.PUT("/:id/categories", middleware.AdminAuth(), h.setCourseCategories)
		courseRouter.PUT("/:id/tags", middleware.AdminAuth(), h.setCourseTags)
		courseRouter.PUT("/:id/publish", middleware.AdminAuth(), h.publishCourse)
		courseRouter.PUT("/:id/archive", middleware.AdminAuth(), h.archiveCourse)
		courseRouter.GET("/:id", h.getCourseDetail)
		courseRouter.GET("", h.listCourses)
		courseRouter.GET("/:id/full", h.getCourseFull)
		courseRouter.GET("/stats", middleware.AdminAuth(), h.getCourseStats)
	}
}

// 创建课程请求
type createCourseRequest struct {
	TeacherID   string             `json:"teacherId" binding:"required"`
	Title       string             `json:"title" binding:"required"`
	Subtitle    string             `json:"subtitle"`
	Description string             `json:"description"`
	Level       entity.CourseLevel `json:"level" binding:"required"`
}

// 创建课程
func (h *CourseHandler) createCourse(c *gin.Context) {
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseAppService.CreateCourse(c.Request.Context(), req.TeacherID, req.Title, req.Subtitle, req.Description, req.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// 更新课程请求
type updateCourseRequest struct {
	Title       string             `json:"title" binding:"required"`
	Subtitle    string             `json:"subtitle"`
	Description string             `json:"description"`
	Level       entity.CourseLevel `json:"level" binding:"required"`
}

// 更新课程
func (h *CourseHandler) updateCourse(c *gin.Context) {
	id := c.Param("id")
	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseAppService.UpdateCourse(c.Request.Context(), id, req.Title, req.Subtitle, req.Description, req.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 更新课程媒体请求
type updateCourseMediaRequest struct {
	CoverImage   string `json:"coverImage"`
	PreviewVideo string `json:"previewVideo"`
}

// 更新课程媒体
func (h *CourseHandler) updateCourseMedia(c *gin.Context) {
	id := c.Param("id")
	var req updateCourseMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseAppService.UpdateCourseMedia(c.Request.Context(), id, req.CoverImage, req.PreviewVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 设置课程分类请求
type setCourseCategoriesRequest struct {
	CategoryIDs []string `json:"categoryIds" binding:"required"`
}

// 设置课程分类
func (h *CourseHandler) setCourseCategories(c *gin.Context) {
	id := c.Param("id")
	var req setCourseCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseAppService.SetCourseCategories(c.Request.Context(), id, req.CategoryIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 设置课程标签请求
type setCourseTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

// 设置课程标签
func (h *CourseHandler) setCourseTags(c *gin.Context) {
	id := c.Param("id")
	var req setCourseTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseAppService.SetCourseTags(c.Request.Context(), id, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 发布课程
func (h *CourseHandler) publishCourse(c *gin.Context) {
	id := c.Param("id")
	course, err := h.courseAppService.PublishCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 归档课程
func (h *CourseHandler) archiveCourse(c *gin.Context) {
	id := c.Param("id")
	course, err := h.courseAppService.ArchiveCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 获取课程详情
func (h *CourseHandler) getCourseDetail(c *gin.Context) {
	id := c.Param("id")
	course, err := h.courseAppService.GetCourseDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// 获取课程列表
func (h *CourseHandler) listCourses(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤参数
	params := repository.CourseQueryParams{
		Status:     c.Query("status"),
		Level:      c.Query("level"),
		TeacherID:  c.Query("teacherId"),
		CategoryID: c.Query("categoryId"),
		Keyword:    c.Query("keyword"),
	}

	courses, total, err := h.courseAppService.ListCourses(c.Request.Context(), params, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": courses,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取完整课程信息
func (h *CourseHandler) getCourseFull(c *gin.Context) {
	id := c.Param("id")
	courseFull, err := h.courseAppService.GetCourseFull(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courseFull)
}

// 获取课程统计信息
func (h *CourseHandler) getCourseStats(c *gin.Context) {
	stats, err := h.courseAppService.GetCourseStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
