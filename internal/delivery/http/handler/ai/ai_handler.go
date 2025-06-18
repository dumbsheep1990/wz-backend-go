package ai

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/delivery/rpc/aiclient"
	"wz-backend-go/api/rpc/ai"
)

// AIHandler AI处理器
type AIHandler struct {
	aiClient aiclient.AI
}

// NewAIHandler 创建AI处理器
func NewAIHandler(aiClient aiclient.AI) *AIHandler {
	return &AIHandler{
		aiClient: aiClient,
	}
}

// Recommend 获取推荐内容
func (h *AIHandler) Recommend(c *gin.Context) {
	var req struct {
		UserId    int64  `json:"userId"`
		Scene     string `json:"scene" binding:"required"`
		PageSize  int32  `json:"pageSize"`
		PageToken string `json:"pageToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID（如果需要）
	if req.UserId == 0 {
		if userID, exists := c.Get("userId"); exists {
			req.UserId = userID.(int64)
		}
	}

	// 默认值
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 调用RPC服务
	rpcReq := &ai.RecommendRequest{
		UserId:    req.UserId,
		Scene:     req.Scene,
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}

	resp, err := h.aiClient.Recommend(c.Request.Context(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": gin.H{
			"items":     resp.Items,
			"pageToken": resp.PageToken,
			"hasMore":   resp.HasMore,
		},
	})
}

// ContentReview 内容审核
func (h *AIHandler) ContentReview(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用RPC服务
	rpcReq := &ai.ContentReviewRequest{
		Content: req.Content,
		Type:    req.Type,
	}

	resp, err := h.aiClient.ContentReview(c.Request.Context(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": gin.H{
			"passed":     resp.Passed,
			"reason":     resp.Reason,
			"labels":     resp.Labels,
			"confidence": resp.Confidence,
		},
	})
}

// Chat 客服对话
func (h *AIHandler) Chat(c *gin.Context) {
	var req struct {
		UserId  int64  `json:"userId"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID（如果需要）
	if req.UserId == 0 {
		if userID, exists := c.Get("userId"); exists {
			req.UserId = userID.(int64)
		}
	}

	// 调用RPC服务
	rpcReq := &ai.ChatRequest{
		UserId:  req.UserId,
		Message: req.Message,
	}

	resp, err := h.aiClient.Chat(c.Request.Context(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": gin.H{
			"reply": resp.Reply,
		},
	})
}

// GetRecommendByScene 根据场景获取推荐（GET方式）
func (h *AIHandler) GetRecommendByScene(c *gin.Context) {
	scene := c.Query("scene")
	if scene == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scene is required"})
		return
	}

	pageSize := int32(10)
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.ParseInt(pageSizeStr, 10, 32); err == nil && ps > 0 {
			pageSize = int32(ps)
		}
	}

	pageToken := c.Query("pageToken")

	// 从JWT中获取用户ID
	var userId int64
	if userID, exists := c.Get("userId"); exists {
		userId = userID.(int64)
	}

	// 调用RPC服务
	rpcReq := &ai.RecommendRequest{
		UserId:    userId,
		Scene:     scene,
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	resp, err := h.aiClient.Recommend(c.Request.Context(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": gin.H{
			"items":     resp.Items,
			"pageToken": resp.PageToken,
			"hasMore":   resp.HasMore,
		},
	})
} 