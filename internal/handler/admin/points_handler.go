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

// AdminPointsHandler 管理员积分管理API处理器
type AdminPointsHandler struct {
	pointsService *service.UserPointsService
}

// NewAdminPointsHandler 创建管理员积分处理器
func NewAdminPointsHandler(pointsService *service.UserPointsService) *AdminPointsHandler {
	return &AdminPointsHandler{
		pointsService: pointsService,
	}
}

// ListPoints 获取积分记录列表
func (h *AdminPointsHandler) ListPoints(c *gin.Context) {
	var req types.ListPointsRequest

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

	if typeStr := c.Query("type"); typeStr != "" {
		t, err := strconv.Atoi(typeStr)
		if err == nil {
			req.Type = t
		}
	}

	req.Source = c.Query("source")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务层获取数据
	points, total, err := h.pointsService.ListPointsWithTotal(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 转换为API响应格式
	var list []*types.UserPointsResponse
	for _, p := range points {
		list = append(list, &types.UserPointsResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Points:      p.Points,
			TotalPoints: p.TotalPoints,
			Type:        p.Type,
			Source:      p.Source,
			Description: p.Description,
			RelatedID:   p.RelatedID,
			RelatedType: p.RelatedType,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	response.Success(c, &types.PagedUserPointsResponse{
		List:  list,
		Total: total,
	})
}

// GetPointByID 获取积分记录详情
func (h *AdminPointsHandler) GetPointByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID参数")
		return
	}

	point, err := h.pointsService.GetPointByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, "积分记录不存在")
		} else {
			response.Fail(c, err.Error())
		}
		return
	}

	response.Success(c, &types.UserPointsResponse{
		ID:          point.ID,
		UserID:      point.UserID,
		Points:      point.Points,
		TotalPoints: point.TotalPoints,
		Type:        point.Type,
		Source:      point.Source,
		Description: point.Description,
		RelatedID:   point.RelatedID,
		RelatedType: point.RelatedType,
		CreatedAt:   point.CreatedAt,
		UpdatedAt:   point.UpdatedAt,
	})
}

// GetUserPoints 获取用户积分明细
func (h *AdminPointsHandler) GetUserPoints(c *gin.Context) {
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

	req := types.ListPointsRequest{
		Page:     page,
		PageSize: pageSize,
		UserID:   userID,
	}

	// 获取筛选参数
	if typeStr := c.Query("type"); typeStr != "" {
		t, err := strconv.Atoi(typeStr)
		if err == nil {
			req.Type = t
		}
	}

	req.Source = c.Query("source")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务层获取数据
	points, total, err := h.pointsService.ListPointsWithTotal(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 转换为API响应格式
	var list []*types.UserPointsResponse
	for _, p := range points {
		list = append(list, &types.UserPointsResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Points:      p.Points,
			TotalPoints: p.TotalPoints,
			Type:        p.Type,
			Source:      p.Source,
			Description: p.Description,
			RelatedID:   p.RelatedID,
			RelatedType: p.RelatedType,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	response.Success(c, &types.PagedUserPointsResponse{
		List:  list,
		Total: total,
	})
}

// AddPoints 管理员添加/调整用户积分
func (h *AdminPointsHandler) AddPoints(c *gin.Context) {
	var req types.AdminAddPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取管理员ID
	adminID, exists := c.Get("AdminID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 创建积分记录
	point := &domain.UserPoints{
		UserID:      req.UserID,
		Points:      req.Points,
		Type:        req.Type,
		Source:      "admin", // 固定来源为管理员
		Description: req.Description,
		OperatorID:  adminID.(int64), // 记录操作管理员ID
	}

	id, err := h.pointsService.CreatePoints(point)
	if err != nil {
		response.Fail(c, "添加积分失败: "+err.Error())
		return
	}

	// 获取创建后的记录
	createdPoint, err := h.pointsService.GetPointByID(id)
	if err != nil {
		response.Success(c, map[string]int64{"id": id})
		return
	}

	response.Success(c, &types.UserPointsResponse{
		ID:          createdPoint.ID,
		UserID:      createdPoint.UserID,
		Points:      createdPoint.Points,
		TotalPoints: createdPoint.TotalPoints,
		Type:        createdPoint.Type,
		Source:      createdPoint.Source,
		Description: createdPoint.Description,
		RelatedID:   createdPoint.RelatedID,
		RelatedType: createdPoint.RelatedType,
		CreatedAt:   createdPoint.CreatedAt,
		UpdatedAt:   createdPoint.UpdatedAt,
	})
}

// DeletePoint 删除/撤销积分调整
func (h *AdminPointsHandler) DeletePoint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID参数")
		return
	}

	// 获取要删除的记录
	point, err := h.pointsService.GetPointByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, "积分记录不存在")
		} else {
			response.Fail(c, err.Error())
		}
		return
	}

	// 只允许删除管理员创建的积分记录
	if point.Source != "admin" {
		response.Fail(c, "只能撤销管理员调整的积分")
		return
	}

	// 调用服务层删除记录
	err = h.pointsService.DeletePoint(id)
	if err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ExportPointsData 导出积分数据
func (h *AdminPointsHandler) ExportPointsData(c *gin.Context) {
	var req types.ListPointsRequest

	// 获取筛选参数
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			req.UserID = userID
		}
	}

	req.Username = c.Query("username")

	if typeStr := c.Query("type"); typeStr != "" {
		t, err := strconv.Atoi(typeStr)
		if err == nil {
			req.Type = t
		}
	}

	req.Source = c.Query("source")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 调用服务导出数据
	data, err := h.pointsService.ExportPointsData(&req)
	if err != nil {
		response.Fail(c, "导出失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename=points_export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// GetPointsStatistics 获取积分统计数据
func (h *AdminPointsHandler) GetPointsStatistics(c *gin.Context) {
	stats, err := h.pointsService.GetPointsStatistics()
	if err != nil {
		response.Fail(c, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPointsRules 获取积分规则
func (h *AdminPointsHandler) GetPointsRules(c *gin.Context) {
	rules, err := h.pointsService.GetPointsRules()
	if err != nil {
		response.Fail(c, "获取积分规则失败: "+err.Error())
		return
	}

	response.Success(c, rules)
}

// UpdatePointsRules 更新积分规则
func (h *AdminPointsHandler) UpdatePointsRules(c *gin.Context) {
	var req types.PointsRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := h.pointsService.UpdatePointsRules(&req)
	if err != nil {
		response.Fail(c, "更新积分规则失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
