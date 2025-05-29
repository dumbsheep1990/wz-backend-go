package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/user/dto"
	"wz-backend-go/internal/application/user/service"
	"wz-backend-go/internal/interfaces/http/middleware"
)

// UserFavoriteController 用户收藏控制器
type UserFavoriteController struct {
	favoriteService *service.UserFavoriteApplicationService
}

// NewUserFavoriteController 创建用户收藏控制器
func NewUserFavoriteController(favoriteService *service.UserFavoriteApplicationService) *UserFavoriteController {
	return &UserFavoriteController{
		favoriteService: favoriteService,
	}
}

// RegisterRoutes 注册路由
func (c *UserFavoriteController) RegisterRoutes(router *gin.RouterGroup) {
	favoriteGroup := router.Group("/favorites")
	{
		// 公共接口
		favoriteGroup.GET("/check", c.CheckFavorite)      // 检查是否已收藏
		favoriteGroup.GET("/hot", c.GetHotContent)        // 获取热门收藏内容
		favoriteGroup.GET("/trend", c.GetFavoritesTrend)  // 获取收藏趋势数据
		favoriteGroup.GET("/statistics", c.GetStatistics) // 获取收藏统计数据

		// 需要用户认证的接口
		userAuth := favoriteGroup.Group("", middleware.UserAuth())
		{
			userAuth.GET("/users/:userId", c.GetUserFavorites) // 获取用户收藏列表
			userAuth.POST("", c.CreateFavorite)                // 创建收藏记录
			userAuth.GET("/:id", c.GetFavoriteByID)            // 获取收藏记录详情
			userAuth.DELETE("/:id", c.DeleteFavorite)          // 删除收藏记录
		}

		// 需要管理员认证的接口
		adminAuth := favoriteGroup.Group("", middleware.AdminAuth())
		{
			adminAuth.POST("/batch-delete", c.BatchDeleteFavorites) // 批量删除收藏记录
			adminAuth.GET("/export", c.ExportFavoritesData)         // 导出收藏数据
		}
	}
}

// CreateFavorite 创建收藏记录
func (c *UserFavoriteController) CreateFavorite(ctx *gin.Context) {
	var req dto.CreateFavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, _ := middleware.GetCurrentUserID(ctx)
	req.UserID = userID

	// 调用应用服务创建收藏
	favorite, err := c.favoriteService.CreateFavorite(context.Background(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建收藏失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建收藏成功",
		"data":    favorite,
	})
}

// GetFavoriteByID 获取收藏详情
func (c *UserFavoriteController) GetFavoriteByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID参数",
		})
		return
	}

	// 调用应用服务获取收藏详情
	favorite, err := c.favoriteService.GetFavoriteByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取收藏详情失败: " + err.Error(),
		})
		return
	}

	// 验证权限：只有管理员或记录创建者可以查看详情
	userID, _ := middleware.GetCurrentUserID(ctx)
	isAdmin := middleware.IsAdmin(ctx)
	if !isAdmin && favorite.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限查看此收藏详情",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取收藏详情成功",
		"data":    favorite,
	})
}

// GetUserFavorites 获取用户收藏列表
func (c *UserFavoriteController) GetUserFavorites(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID参数",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.DefaultQuery("pageSize", "10"), 10, 64)
	itemType := ctx.DefaultQuery("itemType", "")

	// 验证权限：只有管理员或自己可以查看自己的收藏列表
	currentUserID, _ := middleware.GetCurrentUserID(ctx)
	isAdmin := middleware.IsAdmin(ctx)
	if !isAdmin && currentUserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限查看此用户的收藏列表",
		})
		return
	}

	// 构造请求参数
	req := &dto.ListFavoritesRequest{
		Page:     page,
		PageSize: pageSize,
		UserID:   userID,
		ItemType: itemType,
	}

	// 调用应用服务获取收藏列表
	favorites, err := c.favoriteService.ListFavorites(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取收藏列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取收藏列表成功",
		"data":    favorites,
	})
}

// DeleteFavorite 删除收藏记录
func (c *UserFavoriteController) DeleteFavorite(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID参数",
		})
		return
	}

	// 获取当前用户ID和管理员状态
	userID, _ := middleware.GetCurrentUserID(ctx)
	isAdmin := middleware.IsAdmin(ctx)

	// 调用应用服务删除收藏
	err = c.favoriteService.DeleteFavorite(context.Background(), id, userID, isAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除收藏失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除收藏成功",
	})
}

// BatchDeleteFavorites 批量删除收藏记录
func (c *UserFavoriteController) BatchDeleteFavorites(ctx *gin.Context) {
	var req dto.BatchDeleteFavoritesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID和管理员状态
	userID, _ := middleware.GetCurrentUserID(ctx)
	isAdmin := middleware.IsAdmin(ctx)

	// 调用应用服务批量删除收藏
	err := c.favoriteService.BatchDeleteFavorites(context.Background(), &req, userID, isAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量删除收藏失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量删除收藏成功",
	})
}

// CheckFavorite 检查是否已收藏
func (c *UserFavoriteController) CheckFavorite(ctx *gin.Context) {
	userID, _ := strconv.ParseInt(ctx.Query("userId"), 10, 64)
	itemID, _ := strconv.ParseInt(ctx.Query("itemId"), 10, 64)
	itemType := ctx.Query("itemType")

	if userID <= 0 || itemID <= 0 || itemType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要参数",
		})
		return
	}

	// 构造请求参数
	req := &dto.CheckFavoriteRequest{
		UserID:   userID,
		ItemID:   itemID,
		ItemType: itemType,
	}

	// 调用应用服务检查是否已收藏
	result, err := c.favoriteService.CheckFavorite(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查收藏状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "检查收藏状态成功",
		"data":    result,
	})
}

// GetStatistics 获取收藏统计数据
func (c *UserFavoriteController) GetStatistics(ctx *gin.Context) {
	// 调用应用服务获取统计数据
	statistics, err := c.favoriteService.GetFavoritesStatistics(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取收藏统计数据失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取收藏统计数据成功",
		"data":    statistics,
	})
}

// GetHotContent 获取热门收藏内容
func (c *UserFavoriteController) GetHotContent(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	// 调用应用服务获取热门内容
	hotContent, err := c.favoriteService.GetHotContent(context.Background(), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取热门收藏内容失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取热门收藏内容成功",
		"data":    hotContent,
	})
}

// GetFavoritesTrend 获取收藏趋势数据
func (c *UserFavoriteController) GetFavoritesTrend(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	// 构造请求参数
	req := &dto.GetFavoritesTrendRequest{
		Period: period,
	}

	// 调用应用服务获取趋势数据
	trendData, err := c.favoriteService.GetFavoritesTrend(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取收藏趋势数据失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取收藏趋势数据成功",
		"data":    trendData,
	})
}

// ExportFavoritesData 导出收藏数据
func (c *UserFavoriteController) ExportFavoritesData(ctx *gin.Context) {
	// 获取查询参数
	userID, _ := strconv.ParseInt(ctx.Query("userId"), 10, 64)
	username := ctx.Query("username")
	title := ctx.Query("title")
	itemType := ctx.Query("itemType")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	// 构造请求参数
	req := &dto.ExportFavoritesDataRequest{
		UserID:    userID,
		Username:  username,
		Title:     title,
		ItemType:  itemType,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 调用应用服务导出数据
	csvData, err := c.favoriteService.ExportFavoritesData(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导出收藏数据失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	fileName := fmt.Sprintf("favorites_export_%s.csv", time.Now().Format("20060102150405"))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", strconv.Itoa(len(csvData)))

	// 返回CSV数据
	ctx.Data(http.StatusOK, "text/csv", csvData)
}
