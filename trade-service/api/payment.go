package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/trade-service/core/model"
	"wz-backend-go/trade-service/core/service"
)

// PaymentHandler 支付API处理器
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler 创建支付API处理器
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// RegisterRoutes 注册路由
func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	paymentGroup := router.Group("/payments")
	{
		paymentGroup.POST("", h.CreatePayment)
		paymentGroup.GET("/status/:order_no", h.QueryPaymentStatus)
		paymentGroup.POST("/refund", h.RefundPayment)
		paymentGroup.GET("/order/:order_no", h.GetPaymentByOrderNo)
		paymentGroup.GET("/statistics", h.GetPaymentStatistics)

		// 支付回调
		paymentGroup.POST("/notify/alipay", h.HandleAlipayNotify)
		paymentGroup.POST("/notify/wechat", h.HandleWechatNotify)
	}

	// 支付配置管理
	configGroup := router.Group("/payment-configs")
	{
		configGroup.GET("", h.GetPaymentConfigs)
		configGroup.GET("/:type", h.GetPaymentConfigByType)
		configGroup.POST("", h.SavePaymentConfig)
		configGroup.PUT("/:id/status", h.UpdatePaymentConfigStatus)
		configGroup.DELETE("/:id", h.DeletePaymentConfig)
	}
}

// CreatePayment 创建支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req model.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 设置客户端IP
	req.ClientIP = c.ClientIP()

	resp, err := h.paymentService.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建支付失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// QueryPaymentStatus 查询支付状态
func (h *PaymentHandler) QueryPaymentStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}

	status, err := h.paymentService.QueryPaymentStatus(c.Request.Context(), orderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询支付状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": status})
}

// RefundPayment 申请退款
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var req model.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	resp, err := h.paymentService.RefundPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "申请退款失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// GetPaymentByOrderNo 根据订单号获取支付记录
func (h *PaymentHandler) GetPaymentByOrderNo(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取支付记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payment})
}

// GetPaymentStatistics 获取支付统计
func (h *PaymentHandler) GetPaymentStatistics(c *gin.Context) {
	stats, err := h.paymentService.GetPaymentStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取支付统计失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": stats})
}

// HandleAlipayNotify 处理支付宝回调
func (h *PaymentHandler) HandleAlipayNotify(c *gin.Context) {
	// 获取所有回调参数
	c.Request.ParseForm()
	notifyData := make(map[string]string)
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			notifyData[k] = v[0]
		}
	}

	if err := h.paymentService.HandlePaymentNotify(c.Request.Context(), 1, notifyData); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// HandleWechatNotify 处理微信回调
func (h *PaymentHandler) HandleWechatNotify(c *gin.Context) {
	// 读取请求体内容
	data, err := c.GetRawData()
	if err != nil {
		c.XML(http.StatusOK, gin.H{"return_code": "FAIL", "return_msg": "读取请求失败"})
		return
	}

	// 解析成map
	notifyData := make(map[string]string)
	notifyData["raw_data"] = string(data)

	if err := h.paymentService.HandlePaymentNotify(c.Request.Context(), 2, notifyData); err != nil {
		c.XML(http.StatusOK, gin.H{"return_code": "FAIL", "return_msg": err.Error()})
		return
	}

	c.XML(http.StatusOK, gin.H{"return_code": "SUCCESS", "return_msg": "OK"})
}

// GetPaymentConfigs 获取支付配置列表
func (h *PaymentHandler) GetPaymentConfigs(c *gin.Context) {
	configs, err := h.paymentService.GetPaymentConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取支付配置失败: " + err.Error()})
		return
	}

	// 安全起见，返回前清除敏感信息
	for i := range configs {
		configs[i].PrivateKey = "******"
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": configs})
}

// GetPaymentConfigByType 根据类型获取支付配置
func (h *PaymentHandler) GetPaymentConfigByType(c *gin.Context) {
	typeStr := c.Param("type")
	payType, err := strconv.Atoi(typeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的支付类型"})
		return
	}

	config, err := h.paymentService.GetPaymentConfigByType(c.Request.Context(), payType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取支付配置失败: " + err.Error()})
		return
	}

	if config == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "支付配置不存在"})
		return
	}

	// 安全起见，返回前清除敏感信息
	config.PrivateKey = "******"

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// SavePaymentConfig 保存支付配置
func (h *PaymentHandler) SavePaymentConfig(c *gin.Context) {
	var config model.PaymentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.paymentService.SavePaymentConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存支付配置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UpdatePaymentConfigStatus 更新支付配置状态
func (h *PaymentHandler) UpdatePaymentConfigStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的配置ID"})
		return
	}

	var req struct {
		Status bool `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.paymentService.UpdatePaymentConfigStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新支付配置状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DeletePaymentConfig 删除支付配置
func (h *PaymentHandler) DeletePaymentConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的配置ID"})
		return
	}

	if err := h.paymentService.DeletePaymentConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除支付配置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
