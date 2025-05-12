package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/service"
	"wz-backend-go/internal/types"
	"wz-backend-go/internal/utils/response"
)

// AdminFavoritesHandler 管理员收藏管理API处理器
type AdminFavoritesHandler struct {
	favoriteService *service.UserFavoriteService
}

// NewAdminFavoritesHandler 创建管理员收藏处理器
func NewAdminFavoritesHandler(favoriteService *service.UserFavoriteService) *AdminFavoritesHandler {
	return &AdminFavoritesHandler{
		favoriteService: favoriteService,
	}
}

// ListFavorites 获取收藏记录列表
func (h *AdminFavoritesHandler) ListFavorites(c *gin.Context) {
	var req types.ListFavoritesRequest

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req.Page = page
	req.PageSize = pageSize

	// 获取筛选参数
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			req.UserID = userID
		}
	}

	req.Username = c.Query("username")
	req.Title = c.Query("title")
	req.ItemType = c.Query("item_type")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务层获取数据
	favorites, total, err := h.favoriteService.ListFavoritesWithTotal(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 转换为API响应格式
	var list []*types.UserFavoriteResponse
	for _, f := range favorites {
		list = append(list, &types.UserFavoriteResponse{
			ID:        f.ID,
			UserID:    f.UserID,
			Username:  f.Username,
			ItemID:    f.ItemID,
			ItemType:  f.ItemType,
			Title:     f.Title,
			Cover:     f.Cover,
			Summary:   f.Summary,
			URL:       f.URL,
			Remark:    f.Remark,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		})
	}

	response.Success(c, &types.PagedUserFavoriteResponse{
		List:  list,
		Total: total,
	})
}

// GetFavoriteDetail 获取收藏详情
func (h *AdminFavoritesHandler) GetFavoriteDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID参数")
		return
	}

	favorite, err := h.favoriteService.GetFavoriteByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, "收藏记录不存在")
		} else {
			response.Fail(c, err.Error())
		}
		return
	}

	response.Success(c, &types.UserFavoriteResponse{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		Username:  favorite.Username,
		ItemID:    favorite.ItemID,
		ItemType:  favorite.ItemType,
		Title:     favorite.Title,
		Cover:     favorite.Cover,
		Summary:   favorite.Summary,
		URL:       favorite.URL,
		Remark:    favorite.Remark,
		CreatedAt: favorite.CreatedAt,
		UpdatedAt: favorite.UpdatedAt,
	})
}

// DeleteFavorite 删除收藏记录
func (h *AdminFavoritesHandler) DeleteFavorite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID参数")
		return
	}

	// 获取要删除的记录
	_, err = h.favoriteService.GetFavoriteByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, "收藏记录不存在")
		} else {
			response.Fail(c, err.Error())
		}
		return
	}

	// 调用服务层删除记录
	err = h.favoriteService.DeleteFavorite(id)
	if err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDeleteFavorites 批量删除收藏记录
func (h *AdminFavoritesHandler) BatchDeleteFavorites(c *gin.Context) {
	var req types.BatchDeleteFavoritesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := h.favoriteService.BatchDeleteFavorites(req.IDs)
	if err != nil {
		response.Fail(c, "批量删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserFavorites 获取指定用户的收藏列表
func (h *AdminFavoritesHandler) GetUserFavorites(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID参数")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := types.ListFavoritesRequest{
		Page:     page,
		PageSize: pageSize,
		UserID:   userID,
	}

	// 获取筛选参数
	req.Title = c.Query("title")
	req.ItemType = c.Query("item_type")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务层获取数据
	favorites, total, err := h.favoriteService.ListFavoritesWithTotal(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 转换为API响应格式
	var list []*types.UserFavoriteResponse
	for _, f := range favorites {
		list = append(list, &types.UserFavoriteResponse{
			ID:        f.ID,
			UserID:    f.UserID,
			Username:  f.Username,
			ItemID:    f.ItemID,
			ItemType:  f.ItemType,
			Title:     f.Title,
			Cover:     f.Cover,
			Summary:   f.Summary,
			URL:       f.URL,
			Remark:    f.Remark,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		})
	}

	response.Success(c, &types.PagedUserFavoriteResponse{
		List:  list,
		Total: total,
	})
}

// ExportFavoritesData 导出收藏数据
func (h *AdminFavoritesHandler) ExportFavoritesData(c *gin.Context) {
	var req types.ListFavoritesRequest

	// 获取筛选参数
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			req.UserID = userID
		}
	}

	req.Username = c.Query("username")
	req.Title = c.Query("title")
	req.ItemType = c.Query("item_type")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务导出数据
	data, err := h.favoriteService.ExportFavoritesData(&req)
	if err != nil {
		response.Fail(c, "导出失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename=favorites_export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// GetFavoritesStatistics 获取收藏统计数据
func (h *AdminFavoritesHandler) GetFavoritesStatistics(c *gin.Context) {
	stats, err := h.favoriteService.GetFavoritesStatistics()
	if err != nil {
		response.Fail(c, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetHotContent 获取热门收藏内容
func (h *AdminFavoritesHandler) GetHotContent(c *gin.Context) {
	hotContent, err := h.favoriteService.GetHotContent()
	if err != nil {
		response.Fail(c, "获取热门内容失败: "+err.Error())
		return
	}

	response.Success(c, hotContent)
}

// GetFavoritesTrend 获取收藏趋势数据
func (h *AdminFavoritesHandler) GetFavoritesTrend(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	// 验证周期参数
	if period != "week" && period != "month" && period != "year" {
		period = "month"
	}

	trendData, err := h.favoriteService.GetFavoritesTrend(period)
	if err != nil {
		response.Fail(c, "获取趋势数据失败: "+err.Error())
		return
	}

	response.Success(c, trendData)
}
