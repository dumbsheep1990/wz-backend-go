package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/service"
	"wz-backend-go/internal/types"
	"wz-backend-go/internal/utils/response"
)

// FavoritesHandler 用户收藏API处理器
type FavoritesHandler struct {
	favoriteService *service.UserFavoriteService
}

// NewFavoritesHandler 创建用户收藏处理器
func NewFavoritesHandler(favoriteService *service.UserFavoriteService) *FavoritesHandler {
	return &FavoritesHandler{
		favoriteService: favoriteService,
	}
}

// CreateFavorite 创建收藏
func (h *FavoritesHandler) CreateFavorite(c *gin.Context) {
	var req types.CreateUserFavoriteRequest
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

	id, err := h.favoriteService.CreateUserFavorite(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, "创建收藏失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// GetFavoriteDetail 获取收藏详情
func (h *FavoritesHandler) GetFavoriteDetail(c *gin.Context) {
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

	favoriteDetail, err := h.favoriteService.GetUserFavoriteById(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, "获取收藏详情失败: "+err.Error())
		return
	}

	// 检查记录是否属于当前用户
	if favoriteDetail.UserID != userID.(int64) {
		response.Forbidden(c, "无权访问")
		return
	}

	response.Success(c, favoriteDetail)
}

// ListFavorites 获取用户收藏列表
func (h *FavoritesHandler) ListFavorites(c *gin.Context) {
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
	itemType := c.DefaultQuery("item_type", "")

	req := &types.ListUserFavoritesRequest{
		UserID:   userID.(int64),
		Page:     page,
		PageSize: pageSize,
		ItemType: itemType,
	}

	result, err := h.favoriteService.ListUserFavoritesByUser(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, "获取收藏列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteFavorite 删除收藏
func (h *FavoritesHandler) DeleteFavorite(c *gin.Context) {
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

	err = h.favoriteService.DeleteUserFavorite(c.Request.Context(), id, userID.(int64))
	if err != nil {
		response.Fail(c, "删除收藏失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// CheckFavorite 检查是否已收藏
func (h *FavoritesHandler) CheckFavorite(c *gin.Context) {
	// 从上下文获取当前用户ID
	userID, exists := c.Get("UserID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取查询参数
	itemIDStr := c.Query("item_id")
	itemType := c.Query("item_type")

	if itemIDStr == "" || itemType == "" {
		response.BadRequest(c, "缺少必要参数")
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的item_id")
		return
	}

	exists, err = h.favoriteService.CheckUserFavorite(c.Request.Context(), userID.(int64), itemID, itemType)
	if err != nil {
		response.Fail(c, "检查收藏状态失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"exists": exists})
}

// RegisterRoutes 注册路由
func (h *FavoritesHandler) RegisterRoutes(router *gin.RouterGroup) {
	favoritesGroup := router.Group("/favorites")
	{
		favoritesGroup.POST("", h.CreateFavorite)
		favoritesGroup.GET("", h.ListFavorites)
		favoritesGroup.GET("/check", h.CheckFavorite)
		favoritesGroup.GET("/:id", h.GetFavoriteDetail)
		favoritesGroup.DELETE("/:id", h.DeleteFavorite)
	}
}
