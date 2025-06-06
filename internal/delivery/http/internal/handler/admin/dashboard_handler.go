package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"wz-backend-go/internal/application/admin"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/pkg/response"
)

// DashboardHandler 处理仪表盘相关的HTTP请求
type DashboardHandler struct {
	dashboardService *admin.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler(dashboardService *admin.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取仪表盘统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats, c)
}

// GetUserStats 获取用户统计数据
func (h *DashboardHandler) GetUserStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取用户统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.UserStats, c)
}

// GetContentStats 获取内容统计数据
func (h *DashboardHandler) GetContentStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取内容统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.ContentStats, c)
}

// GetTradeStats 获取交易统计数据
func (h *DashboardHandler) GetTradeStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取交易统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.TradeStats, c)
}

// GetSiteStats 获取站点统计数据
func (h *DashboardHandler) GetSiteStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取站点统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.SiteStats, c)
}

// GetCommunityStats 获取社区统计数据
func (h *DashboardHandler) GetCommunityStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取社区统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.CommunityStats, c)
}

// GetComponentStats 获取组件统计数据
func (h *DashboardHandler) GetComponentStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取组件统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.ComponentStats, c)
}

// GetRenderStats 获取渲染统计数据
func (h *DashboardHandler) GetRenderStats(c *gin.Context) {
	// 从上下文获取请求ID，用于日志跟踪
	requestId, _ := c.Get("X-Request-Id")

	// 获取仪表盘统计数据
	stats, err := h.dashboardService.GetDashboardStats(c)
	if err != nil {
		response.FailWithMessage(errors.Wrap(err, "获取渲染统计数据失败").Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithData(stats.RenderStats, c)
}
