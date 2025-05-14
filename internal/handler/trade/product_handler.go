package tradehandler // 包名可以根据项目规范调整，例如 handler 或 http_handler

import (
	"net/http"
	"strconv"

	"wz-project/wz-backend-go/internal/service" // 假设的项目路径
	// "wz-project/wz-backend-go/pkg/apperror"   // 假设的自定义错误包
	// "wz-project/wz-backend-go/pkg/response" // 假设的统一响应包

	"github.com/gin-gonic/gin" // 假设使用 Gin 框架，根据实际情况调整
)

// ProductHandler 封装了产品相关的 HTTP Handler 方法
type ProductHandler struct {
	productService service.ProductService
	// 其他依赖，例如 logger
}

// NewProductHandler 创建一个新的 ProductHandler 实例
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// RegisterRoutes 在 Gin Engine 上注册产品相关的路由
// 根据项目结构，这个方法可能在更高层次的 router.go 文件中被调用
func (h *ProductHandler) RegisterRoutes(router *gin.Engine) { // 或者 *gin.RouterGroup
	productGroup := router.Group("/api/v1/products") // 或者根据 trade.api 定义调整
	{
		productGroup.GET("/:product_id", h.GetProductDetail)
		// 其他产品相关路由...
	}
}

// GetProductDetail 处理获取产品详情的请求
// GET /api/v1/products/:product_id
func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil || productID == 0 {
		// response.BadRequest(c, apperror.ErrInvalidParam.WithMessage("Invalid product ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		// // 根据错误类型返回不同的 HTTP 状态码
		// if errors.Is(err, model.ErrProductNotFound) { // 需要导入 model 和 errors
		// 	response.NotFound(c, apperror.ErrNotFound.Wrap(err, "product not found"))
		// 	return
		// }
		// // 其他特定业务错误处理
		// response.InternalError(c, apperror.ErrInternal.Wrap(err, "failed to get product detail"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product detail"}) // 简化错误处理
		return
	}

	// response.Success(c, product)
	c.JSON(http.StatusOK, product)
}

// 其他 Product Handler 方法 (ListProducts, GetRelatedProducts) 可以后续添加
