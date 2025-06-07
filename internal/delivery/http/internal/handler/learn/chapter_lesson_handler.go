package learn

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/delivery/http/internal/middleware"
)

// ChapterLessonHandler 处理章节和课时相关的HTTP请求
type ChapterLessonHandler struct {
	chapterLessonAppService *learn.ChapterLessonAppService
}

// NewChapterLessonHandler 创建章节和课时处理器
func NewChapterLessonHandler(chapterLessonAppService *learn.ChapterLessonAppService) *ChapterLessonHandler {
	return &ChapterLessonHandler{
		chapterLessonAppService: chapterLessonAppService,
	}
}

// RegisterRoutes 注册路由
func (h *ChapterLessonHandler) RegisterRoutes(router *gin.RouterGroup) {
	// 章节路由
	chapterRouter := router.Group("/chapters")
	{
		chapterRouter.POST("", middleware.AdminAuth(), h.createChapter)
		chapterRouter.PUT("/:id", middleware.AdminAuth(), h.updateChapter)
		chapterRouter.DELETE("/:id", middleware.AdminAuth(), h.deleteChapter)
		chapterRouter.GET("/:id", h.getChapterByID)
		chapterRouter.GET("/course/:courseId", h.getChaptersByCourseID)
	}
	
	// 课时路由
	lessonRouter := router.Group("/lessons")
	{
		lessonRouter.POST("", middleware.AdminAuth(), h.createLesson)
		lessonRouter.PUT("/:id", middleware.AdminAuth(), h.updateLesson)
		lessonRouter.PUT("/:id/publish", middleware.AdminAuth(), h.publishLesson)
		lessonRouter.DELETE("/:id", middleware.AdminAuth(), h.deleteLesson)
		lessonRouter.GET("/:id", h.getLessonByID)
		lessonRouter.PUT("/:id/video", middleware.AdminAuth(), h.setVideoContent)
		lessonRouter.PUT("/:id/article", middleware.AdminAuth(), h.setArticleContent)
		lessonRouter.PUT("/:id/audio", middleware.AdminAuth(), h.setAudioContent)
	}
}

// 创建章节请求
type createChapterRequest struct {
	CourseID    string `json:"courseId" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// 创建章节
func (h *ChapterLessonHandler) createChapter(c *gin.Context) {
	var req createChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	chapter, err := h.chapterLessonAppService.CreateChapter(
		c.Request.Context(), req.CourseID, req.Title, req.Description, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, chapter)
}

// 更新章节请求
type updateChapterRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// 更新章节
func (h *ChapterLessonHandler) updateChapter(c *gin.Context) {
	id := c.Param("id")
	var req updateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	chapter, err := h.chapterLessonAppService.UpdateChapter(
		c.Request.Context(), id, req.Title, req.Description, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, chapter)
}

// 删除章节
func (h *ChapterLessonHandler) deleteChapter(c *gin.Context) {
	id := c.Param("id")
	err := h.chapterLessonAppService.DeleteChapter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "章节已删除"})
}

// 获取章节详情
func (h *ChapterLessonHandler) getChapterByID(c *gin.Context) {
	id := c.Param("id")
	chapter, err := h.chapterLessonAppService.GetChapterByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, chapter)
}

// 获取课程的所有章节
func (h *ChapterLessonHandler) getChaptersByCourseID(c *gin.Context) {
	courseID := c.Param("courseId")
	chapters, err := h.chapterLessonAppService.GetChaptersByCourseID(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, chapters)
}

// 创建课时请求
type createLessonRequest struct {
	CourseID    string             `json:"courseId" binding:"required"`
	ChapterID   string             `json:"chapterId" binding:"required"`
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description"`
	Order       int                `json:"order"`
	Type        entity.LessonType  `json:"type" binding:"required"`
	IsFree      bool               `json:"isFree"`
}

// 创建课时
func (h *ChapterLessonHandler) createLesson(c *gin.Context) {
	var req createLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lesson, err := h.chapterLessonAppService.CreateLesson(
		c.Request.Context(), req.CourseID, req.ChapterID, req.Title,
		req.Description, req.Order, req.Type, req.IsFree)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, lesson)
}

// 更新课时请求
type updateLessonRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	IsFree      bool   `json:"isFree"`
}

// 更新课时
func (h *ChapterLessonHandler) updateLesson(c *gin.Context) {
	id := c.Param("id")
	var req updateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lesson, err := h.chapterLessonAppService.UpdateLesson(
		c.Request.Context(), id, req.Title, req.Description, req.Order, req.IsFree)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}

// 发布课时
func (h *ChapterLessonHandler) publishLesson(c *gin.Context) {
	id := c.Param("id")
	lesson, err := h.chapterLessonAppService.PublishLesson(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}

// 删除课时
func (h *ChapterLessonHandler) deleteLesson(c *gin.Context) {
	id := c.Param("id")
	err := h.chapterLessonAppService.DeleteLesson(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "课时已删除"})
}

// 获取课时详情
func (h *ChapterLessonHandler) getLessonByID(c *gin.Context) {
	id := c.Param("id")
	lesson, err := h.chapterLessonAppService.GetLessonByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}

// 设置视频内容请求
type setVideoContentRequest struct {
	VideoURL string `json:"videoUrl" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}

// 设置视频内容
func (h *ChapterLessonHandler) setVideoContent(c *gin.Context) {
	id := c.Param("id")
	var req setVideoContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lesson, err := h.chapterLessonAppService.SetVideoContent(
		c.Request.Context(), id, req.VideoURL, req.Duration, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}

// 设置文章内容请求
type setArticleContentRequest struct {
	Content string `json:"content" binding:"required"`
}

// 设置文章内容
func (h *ChapterLessonHandler) setArticleContent(c *gin.Context) {
	id := c.Param("id")
	var req setArticleContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lesson, err := h.chapterLessonAppService.SetArticleContent(
		c.Request.Context(), id, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}

// 设置音频内容请求
type setAudioContentRequest struct {
	AudioURL string `json:"audioUrl" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}

// 设置音频内容
func (h *ChapterLessonHandler) setAudioContent(c *gin.Context) {
	id := c.Param("id")
	var req setAudioContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lesson, err := h.chapterLessonAppService.SetAudioContent(
		c.Request.Context(), id, req.AudioURL, req.Duration, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lesson)
}
