package content

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/repository"
)

// PostHandler 帖子处理器
type PostHandler struct {
	contentRepo repository.ContentRepository
}

// NewPostHandler 创建帖子处理器
func NewPostHandler(contentRepo repository.ContentRepository) *PostHandler {
	return &PostHandler{
		contentRepo: contentRepo,
	}
}

// CreatePost 创建帖子
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryId int64  `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	content := &repository.Content{
		Type:       "post",
		Title:      req.Title,
		Content:    req.Content,
		UserId:     userID.(int64),
		CategoryId: req.CategoryId,
		Status:     1, // 默认发布状态
	}

	id, err := h.contentRepo.CreateContent(c.Request.Context(), content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	content.ID = id
	c.JSON(http.StatusCreated, gin.H{"post": content})
}

// UpdatePost 更新帖子
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryId int64  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查帖子是否存在
	post, err := h.contentRepo.GetContentById(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查权限（只能编辑自己的帖子）
	userID, exists := c.Get("userId")
	if !exists || post.UserId != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryId != 0 {
		updates["category_id"] = req.CategoryId
	}

	err = h.contentRepo.UpdateContent(c.Request.Context(), postId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的帖子信息
	updatedPost, err := h.contentRepo.GetContentById(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": updatedPost})
}

// DeletePost 删除帖子
func (h *PostHandler) DeletePost(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 检查帖子是否存在
	post, err := h.contentRepo.GetContentById(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查权限（只能删除自己的帖子）
	userID, exists := c.Get("userId")
	if !exists || post.UserId != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = h.contentRepo.DeleteContent(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPost 获取帖子详情
func (h *PostHandler) GetPost(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.contentRepo.GetContentById(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 增加浏览计数
	go func() {
		h.contentRepo.IncrementViewCount(c.Request.Context(), postId)
	}()

	c.JSON(http.StatusOK, post)
}

// ListPosts 获取帖子列表
func (h *PostHandler) ListPosts(c *gin.Context) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	filters := make(map[string]interface{})
	filters["type"] = "post"

	if categoryIdStr := c.Query("category_id"); categoryIdStr != "" {
		if categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64); err == nil {
			filters["category_id"] = categoryId
		}
	}

	if userIdStr := c.Query("user_id"); userIdStr != "" {
		if userId, err := strconv.ParseInt(userIdStr, 10, 64); err == nil {
			filters["user_id"] = userId
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.ParseInt(statusStr, 10, 32); err == nil {
			filters["status"] = int32(status)
		}
	}

	posts, total, err := h.contentRepo.GetContentList(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
} 