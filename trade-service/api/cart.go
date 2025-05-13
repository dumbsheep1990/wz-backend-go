package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/trade-service/core/model"
	"wz-backend-go/trade-service/core/service"
)

// CartHandler 购物车API处理器
type CartHandler struct {
	cartService service.CartService
}

// NewCartHandler 创建购物车API处理器
func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// RegisterRoutes 注册路由
func (h *CartHandler) RegisterRoutes(router *gin.RouterGroup) {
	cartGroup := router.Group("/carts")
	{
		cartGroup.POST("", h.AddToCart)
		cartGroup.GET("/user/:user_id", h.GetUserCart)
		cartGroup.PUT("/:id", h.UpdateQuantity)
		cartGroup.DELETE("/:id", h.RemoveFromCart)
		cartGroup.DELETE("/user/:user_id", h.ClearCart)
		cartGroup.PUT("/user/:user_id/product/:product_id/selected", h.UpdateSelected)
		cartGroup.PUT("/user/:user_id/selected", h.SelectAll)
		cartGroup.GET("/user/:user_id/selected", h.GetSelectedItems)
		cartGroup.GET("/user/:user_id/count", h.GetCartItemCount)
		cartGroup.GET("/user/:user_id/total", h.GetCartTotal)
		cartGroup.POST("/checkout", h.Checkout)
	}
}

// AddToCart 添加商品到购物车
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req struct {
		UserID int64                 `json:"user_id" binding:"required"`
		Item   model.CartItemRequest `json:"item" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.cartService.AddToCart(c.Request.Context(), req.UserID, req.Item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加购物车失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetUserCart 获取用户购物车
func (h *CartHandler) GetUserCart(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	cart, err := h.cartService.GetUserCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取购物车失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": cart})
}

// UpdateQuantity 更新购物车商品数量
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	var req struct {
		UserID    int64 `json:"user_id" binding:"required"`
		ProductID int64 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.cartService.UpdateQuantity(c.Request.Context(), req.UserID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新购物车数量失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// RemoveFromCart 删除购物车商品
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的商品ID"})
		return
	}

	if err := h.cartService.RemoveFromCart(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除购物车商品失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ClearCart 清空购物车
func (h *CartHandler) ClearCart(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清空购物车失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UpdateSelected 更新选中状态
func (h *CartHandler) UpdateSelected(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的商品ID"})
		return
	}

	var req struct {
		Selected bool `json:"selected" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.cartService.UpdateSelected(c.Request.Context(), userID, productID, req.Selected); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新选中状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// SelectAll 全选/全不选
func (h *CartHandler) SelectAll(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	var req struct {
		Selected bool `json:"selected" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.cartService.SelectAll(c.Request.Context(), userID, req.Selected); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "全选/全不选失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetSelectedItems 获取选中的购物车项
func (h *CartHandler) GetSelectedItems(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	items, err := h.cartService.GetSelectedItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取选中商品失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": items})
}

// GetCartItemCount 获取购物车商品数量
func (h *CartHandler) GetCartItemCount(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	count, err := h.cartService.GetCartItemCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取购物车商品数量失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": count})
}

// GetCartTotal 获取购物车总金额
func (h *CartHandler) GetCartTotal(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	total, err := h.cartService.GetCartTotal(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取购物车总金额失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": total})
}

// Checkout 购物车结算
func (h *CartHandler) Checkout(c *gin.Context) {
	var req struct {
		UserID    int64  `json:"user_id" binding:"required"`
		Address   string `json:"address" binding:"required"`
		Consignee string `json:"consignee" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Remark    string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	order, err := h.cartService.Checkout(c.Request.Context(), req.UserID, req.Address, req.Consignee, req.Phone, req.Remark)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "结算失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}
