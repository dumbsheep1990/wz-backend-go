package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/trade-service/core/model"
	"wz-backend-go/trade-service/core/service"
)

// OrderHandler 订单API处理器
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler 创建订单API处理器
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// RegisterRoutes 注册路由
func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("", h.CreateOrder)
		orderGroup.GET("", h.ListOrders)
		orderGroup.GET("/:id", h.GetOrderDetail)
		orderGroup.GET("/user/:user_id", h.GetUserOrders)
		orderGroup.PUT("/:id/cancel", h.CancelOrder)
		orderGroup.PUT("/:id/ship", h.ShipOrder)
		orderGroup.PUT("/:id/confirm", h.ConfirmReceipt)
		orderGroup.POST("/:id/refund", h.RefundOrder)
		orderGroup.DELETE("/:id", h.DeleteOrder)
		orderGroup.POST("/export", h.ExportOrders)
		orderGroup.GET("/statistics", h.GetOrderStatistics)
		orderGroup.GET("/status-statistics", h.GetStatusStatistics)
		orderGroup.GET("/payment-statistics", h.GetPaymentTypeStatistics)
		orderGroup.GET("/trend", h.GetOrderTrend)
	}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.orderService.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// GetOrderDetail 获取订单详情
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	order, err := h.orderService.GetOrderDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取订单详情失败: " + err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// ListOrders 获取订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var query model.OrderQuery

	// 解析查询参数
	query.OrderNo = c.Query("order_no")
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			query.UserID = userID
		}
	}

	statusStr := c.Query("status")
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			statusInt := int(status)
			query.Status = &statusInt
		}
	}

	payTypeStr := c.Query("pay_type")
	if payTypeStr != "" {
		payType, err := strconv.Atoi(payTypeStr)
		if err == nil {
			payTypeInt := int(payType)
			query.PayType = &payTypeInt
		}
	}

	query.StartDate = c.Query("start_date")
	query.EndDate = c.Query("end_date")

	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	query.Page = page

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	query.PageSize = pageSize

	// 查询订单列表
	orders, total, err := h.orderService.ListOrders(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取订单列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":     orders,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetUserOrders 获取用户订单列表
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询用户订单列表
	orders, total, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户订单列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":     orders,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	if err := h.orderService.CancelOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ShipOrder 订单发货
func (h *OrderHandler) ShipOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	var req struct {
		LogisticsCompany string `json:"logistics_company" binding:"required"`
		LogisticsNo      string `json:"logistics_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.orderService.ShipOrder(c.Request.Context(), id, req.LogisticsCompany, req.LogisticsNo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "订单发货失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ConfirmReceipt 确认收货
func (h *OrderHandler) ConfirmReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	if err := h.orderService.ConfirmReceipt(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "确认收货失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// RefundOrder 订单退款
func (h *OrderHandler) RefundOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Reason string  `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.orderService.RefundOrder(c.Request.Context(), id, req.Amount, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "订单退款失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DeleteOrder 删除订单
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订单ID"})
		return
	}

	if err := h.orderService.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ExportOrders 导出订单数据
func (h *OrderHandler) ExportOrders(c *gin.Context) {
	var query model.OrderQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	data, err := h.orderService.ExportOrders(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出订单数据失败: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename=orders.xlsx")
	c.Data(http.StatusOK, "application/vnd.ms-excel", data)
}

// GetOrderStatistics 获取订单统计信息
func (h *OrderHandler) GetOrderStatistics(c *gin.Context) {
	stats, err := h.orderService.GetOrderStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取订单统计信息失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": stats})
}

// GetStatusStatistics 获取订单状态统计
func (h *OrderHandler) GetStatusStatistics(c *gin.Context) {
	stats, err := h.orderService.GetStatusStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取订单状态统计失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": stats})
}

// GetPaymentTypeStatistics 获取支付方式统计
func (h *OrderHandler) GetPaymentTypeStatistics(c *gin.Context) {
	stats, err := h.orderService.GetPaymentTypeStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取支付方式统计失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": stats})
}

// GetOrderTrend 获取订单趋势
func (h *OrderHandler) GetOrderTrend(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	if period != "week" && period != "month" && period != "year" {
		period = "month"
	}

	trend, err := h.orderService.GetOrderTrend(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取订单趋势失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": trend})
}
