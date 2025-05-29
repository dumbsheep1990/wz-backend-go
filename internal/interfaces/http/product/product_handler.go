package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/product/dto"
	"wz-backend-go/internal/application/product/service"
)

// ProductHandler 商品HTTP处理器
type ProductHandler struct {
	productAppService *service.ProductApplicationService
}

// NewProductHandler 创建商品HTTP处理器
func NewProductHandler(productAppService *service.ProductApplicationService) *ProductHandler {
	return &ProductHandler{
		productAppService: productAppService,
	}
}

// CreateProduct 创建商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.ProductCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文中获取创建者ID
	creatorID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 将creatorID转换为int64
	creatorIDInt64, ok := creatorID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	response, err := h.productAppService.CreateProduct(&req, creatorIDInt64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProduct 更新商品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	var req dto.ProductUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.productAppService.UpdateProduct(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductByID 根据ID获取商品
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	response, err := h.productAppService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductBySKU 根据SKU获取商品
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Query("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SKU不能为空",
		})
		return
	}

	response, err := h.productAppService.GetProductBySKU(sku)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PublishProduct 发布商品
func (h *ProductHandler) PublishProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	response, err := h.productAppService.PublishProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UnpublishProduct 下架商品
func (h *ProductHandler) UnpublishProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	response, err := h.productAppService.UnpublishProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProduct 删除商品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	response, err := h.productAppService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateStock 更新库存
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	var req dto.StockUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.productAppService.UpdateStock(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DecrementStock 减少库存
func (h *ProductHandler) DecrementStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	var req dto.StockDecrementRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.productAppService.DecrementStock(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// IncrementStock 增加库存
func (h *ProductHandler) IncrementStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商品ID",
		})
		return
	}

	var req dto.StockIncrementRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.productAppService.IncrementStock(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProducts 分页获取商品列表
func (h *ProductHandler) GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := h.productAppService.GetProducts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByCategory 获取分类下的商品
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的分类ID",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := h.productAppService.GetProductsByCategory(categoryID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByCreator 获取创建者的商品
func (h *ProductHandler) GetProductsByCreator(c *gin.Context) {
	creatorIDStr := c.Param("creatorId")
	creatorID, err := strconv.ParseInt(creatorIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的创建者ID",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := h.productAppService.GetProductsByCreator(creatorID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchProducts 搜索商品
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "关键词不能为空",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := h.productAppService.SearchProducts(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetActiveProducts 获取上架商品
func (h *ProductHandler) GetActiveProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	response, err := h.productAppService.GetActiveProducts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes 注册路由
func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/products")
	{
		group.POST("", h.CreateProduct)
		group.PUT("/:id", h.UpdateProduct)
		group.GET("", h.GetProducts)
		group.GET("/:id", h.GetProductByID)
		group.GET("/sku", h.GetProductBySKU)
		group.PUT("/:id/publish", h.PublishProduct)
		group.PUT("/:id/unpublish", h.UnpublishProduct)
		group.DELETE("/:id", h.DeleteProduct)
		group.PUT("/:id/stock", h.UpdateStock)
		group.PUT("/:id/stock/decrement", h.DecrementStock)
		group.PUT("/:id/stock/increment", h.IncrementStock)
		group.GET("/category/:categoryId", h.GetProductsByCategory)
		group.GET("/creator/:creatorId", h.GetProductsByCreator)
		group.GET("/search", h.SearchProducts)
		group.GET("/active", h.GetActiveProducts)
	}
}
