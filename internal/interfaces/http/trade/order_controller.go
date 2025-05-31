package trade

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/trade/dto"
	"wz-backend-go/internal/application/trade/service"
)

// OrderController 订单控制器
type OrderController struct {
	orderAppSvc *service.OrderApplicationService
}

// NewOrderController 创建订单控制器
func NewOrderController(orderAppSvc *service.OrderApplicationService) *OrderController {
	return &OrderController{
		orderAppSvc: orderAppSvc,
	}
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 创建新的订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "订单信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.CreateOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 根据订单ID获取订单详情
// @Tags 订单
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param user_id query int64 true "用户ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/{id} [get]
func (c *OrderController) GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	userIDStr := ctx.Query("user_id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	req := &dto.GetOrderRequest{
		OrderID: orderID,
		UserID:  userID,
	}

	order, err := c.orderAppSvc.GetOrder(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// GetOrderByNumber 根据订单号获取订单
// @Summary 根据订单号获取订单
// @Description 根据订单号获取订单详情
// @Tags 订单
// @Accept json
// @Produce json
// @Param number path string true "订单号"
// @Param user_id query int64 true "用户ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/number/{number} [get]
func (c *OrderController) GetOrderByNumber(ctx *gin.Context) {
	orderNumber := ctx.Param("number")
	userIDStr := ctx.Query("user_id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	req := &dto.GetOrderByNumberRequest{
		OrderNumber: orderNumber,
		UserID:      userID,
	}

	order, err := c.orderAppSvc.GetOrderByNumber(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// ListUserOrders 获取用户订单列表
// @Summary 获取用户订单列表
// @Description 分页获取用户的订单列表
// @Tags 订单
// @Accept json
// @Produce json
// @Param customer_id query int64 true "客户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.ListOrdersResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders [get]
func (c *OrderController) ListUserOrders(ctx *gin.Context) {
	customerIDStr := ctx.Query("customer_id")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := &dto.ListUserOrdersRequest{
		CustomerID: customerID,
		Page:       page,
		PageSize:   pageSize,
	}

	orders, err := c.orderAppSvc.ListUserOrders(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// PayOrder 支付订单
// @Summary 支付订单
// @Description 支付指定的订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.PayOrderRequest true "支付信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/pay [post]
func (c *OrderController) PayOrder(ctx *gin.Context) {
	var req dto.PayOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.PayOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// ShipOrder 发货
// @Summary 订单发货
// @Description 为订单安排发货
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.ShipOrderRequest true "发货信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/ship [post]
func (c *OrderController) ShipOrder(ctx *gin.Context) {
	var req dto.ShipOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.ShipOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// DeliverOrder 确认送达
// @Summary 确认订单送达
// @Description 确认订单已送达
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.DeliverOrderRequest true "送达信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/deliver [post]
func (c *OrderController) DeliverOrder(ctx *gin.Context) {
	var req dto.DeliverOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.DeliverOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// CompleteOrder 完成订单
// @Summary 完成订单
// @Description 标记订单为已完成
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.CompleteOrderRequest true "完成信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/complete [post]
func (c *OrderController) CompleteOrder(ctx *gin.Context) {
	var req dto.CompleteOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.CompleteOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 取消指定的订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.CancelOrderRequest true "取消信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/cancel [post]
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	var req dto.CancelOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.CancelOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// RefundOrder 退款
// @Summary 订单退款
// @Description 申请订单退款
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body dto.RefundOrderRequest true "退款信息"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/refund [post]
func (c *OrderController) RefundOrder(ctx *gin.Context) {
	var req dto.RefundOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderAppSvc.RefundOrder(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// SearchOrders 搜索订单
// @Summary 搜索订单
// @Description 根据关键词搜索订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.ListOrdersResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/search [get]
func (c *OrderController) SearchOrders(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := &dto.SearchOrdersRequest{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	orders, err := c.orderAppSvc.SearchOrders(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
} 