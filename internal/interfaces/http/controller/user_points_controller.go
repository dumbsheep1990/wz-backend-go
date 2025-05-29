package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/user/dto"
	"wz-backend-go/internal/application/user/service"
	"wz-backend-go/internal/interfaces/http/middleware"
)

// UserPointsController 用户积分控制器
type UserPointsController struct {
	userPointsService *service.UserPointsApplicationService
}

// NewUserPointsController 创建用户积分控制器
func NewUserPointsController(userPointsService *service.UserPointsApplicationService) *UserPointsController {
	return &UserPointsController{
		userPointsService: userPointsService,
	}
}

// Register 注册路由
func (c *UserPointsController) Register(router *gin.RouterGroup) {
	points := router.Group("/points")
	{
		// 获取用户积分记录
		points.GET("/users/:userId", c.GetUserPoints)

		// 创建积分记录
		points.POST("", middleware.AdminAuth(), c.CreatePoints)

		// 获取积分记录详情
		points.GET("/:id", c.GetPointDetails)

		// 撤销积分记录
		points.PUT("/:id/revoke", middleware.AdminAuth(), c.RevokePoint)

		// 获取积分统计
		points.GET("/statistics", middleware.AdminAuth(), c.GetStatistics)

		// 积分规则相关
		rules := points.Group("/rules")
		{
			// 获取积分规则
			rules.GET("", c.GetPointsRules)

			// 更新积分规则
			rules.PUT("", middleware.AdminAuth(), c.UpdatePointsRules)
		}
	}
}

// GetUserPoints 获取用户积分记录
func (c *UserPointsController) GetUserPoints(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.DefaultQuery("page_size", "10"), 10, 64)

	result, err := c.userPointsService.ListPointsByUserID(ctx, userID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// CreatePoints 创建积分记录
func (c *UserPointsController) CreatePoints(ctx *gin.Context) {
	var req dto.CreatePointsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置操作者ID（从认证中获取）
	operatorID, _ := ctx.Get("user_id")
	req.OperatorID = operatorID.(int64)

	result, err := c.userPointsService.CreatePoints(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// GetPointDetails 获取积分记录详情
func (c *UserPointsController) GetPointDetails(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.userPointsService.GetPointByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// RevokePoint 撤销积分记录
func (c *UserPointsController) RevokePoint(ctx *gin.Context) {
	id := ctx.Param("id")

	// 获取操作者ID（从认证中获取）
	operatorID, _ := ctx.Get("user_id")

	err := c.userPointsService.RevokePoint(ctx, id, operatorID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "积分记录已撤销"})
}

// GetStatistics 获取积分统计
func (c *UserPointsController) GetStatistics(ctx *gin.Context) {
	result, err := c.userPointsService.GetPointsStatistics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetPointsRules 获取积分规则
func (c *UserPointsController) GetPointsRules(ctx *gin.Context) {
	tenantID, _ := strconv.ParseInt(ctx.DefaultQuery("tenant_id", "0"), 10, 64)

	result, err := c.userPointsService.GetPointsRules(ctx, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdatePointsRules 更新积分规则
func (c *UserPointsController) UpdatePointsRules(ctx *gin.Context) {
	var req dto.PointsRulesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.userPointsService.UpdatePointsRules(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
