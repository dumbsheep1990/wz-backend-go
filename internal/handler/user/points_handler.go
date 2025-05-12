package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/service"
	"wz-backend-go/internal/types"
	"wz-backend-go/internal/utils/response"
)

// PointsHandler 用户积分API处理器
type PointsHandler struct {
	pointsService *service.UserPointsService
}

// NewPointsHandler 创建用户积分处理器
func NewPointsHandler(pointsService *service.UserPointsService) *PointsHandler {
	return &PointsHandler{
		pointsService: pointsService,
	}
}

// CreatePoints 创建积分记录
func (h *PointsHandler) CreatePoints(c *gin.Context) {
	var req types.CreateUserPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID, exists := c.Get("UserID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	req.UserID = userID.(int64)

	id, err := h.pointsService.CreateUserPointsRequest(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, "创建积分记录失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// GetPointsDetail 获取积分记录详情
func (h *PointsHandler) GetPointsDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, exists := c.Get("UserID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	pointsDetail, err := h.pointsService.GetUserPointsById(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, "获取积分记录失败: "+err.Error())
		return
	}

	// 检查记录是否属于当前用户
	if pointsDetail.UserID != userID.(int64) {
		response.Forbidden(c, "无权访问")
		return
	}

	response.Success(c, pointsDetail)
}

// ListPoints 获取用户积分记录列表
func (h *PointsHandler) ListPoints(c *gin.Context) {
	// 从上下文获取当前用户ID
	userID, exists := c.Get("UserID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 解析类型过滤参数
	typeFilter, err := strconv.Atoi(c.DefaultQuery("type", "0"))
	if err != nil {
		typeFilter = 0 // 0表示不过滤
	}

	req := &types.ListUserPointsRequest{
		UserID:   userID.(int64),
		Page:     page,
		PageSize: pageSize,
	}

	// 如果有类型过滤
	if typeFilter > 0 {
		req.Type = typeFilter
	}

	result, err := h.pointsService.ListUserPointsByUser(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, "获取积分记录列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetTotalPoints 获取用户总积分
func (h *PointsHandler) GetTotalPoints(c *gin.Context) {
	// 从上下文获取当前用户ID
	userID, exists := c.Get("UserID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	totalPoints, err := h.pointsService.GetUserTotalPoints(userID.(int64))
	if err != nil {
		response.Fail(c, "获取用户总积分失败: "+err.Error())
		return
	}

	response.Success(c, totalPoints)
}

// RegisterRoutes 注册路由
func (h *PointsHandler) RegisterRoutes(router *gin.RouterGroup) {
	pointsGroup := router.Group("/points")
	{
		pointsGroup.POST("", h.CreatePoints)
		pointsGroup.GET("", h.ListPoints)
		pointsGroup.GET("/:id", h.GetPointsDetail)
		pointsGroup.GET("/total", h.GetTotalPoints)
	}
}
