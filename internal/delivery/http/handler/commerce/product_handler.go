package commerce

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	"wz-backend-go/internal/application/commerce"
	"wz-backend-go/internal/domain/commerce/dto"
)

// ProductHandler handles HTTP requests for commerce products
type ProductHandler struct {
	productAppService *commerce.ProductAppService
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(productAppService *commerce.ProductAppService) *ProductHandler {
	return &ProductHandler{
		productAppService: productAppService,
	}
}

// CreateProduct handles requests to create a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productAppService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct handles requests to update an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the one in the request body
	if c.Param("id") != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID mismatch"})
		return
	}

	product, err := h.productAppService.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductByID handles requests to retrieve a product by its ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productAppService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles requests to retrieve products with filtering
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter dto.ProductFilterRequest
	
	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	
	products, err := h.productAppService.GetProducts(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, products)
}

// GetFeaturedProducts handles requests to retrieve featured products
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit := 10 // Default limit
	
	// Get limit from query parameter if provided
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	products, err := h.productAppService.GetFeaturedProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetNewProducts handles requests to retrieve new products
func (h *ProductHandler) GetNewProducts(c *gin.Context) {
	limit := 10 // Default limit
	
	// Get limit from query parameter if provided
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	products, err := h.productAppService.GetNewProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetPopularProducts handles requests to retrieve popular products
func (h *ProductHandler) GetPopularProducts(c *gin.Context) {
	limit := 10 // Default limit
	
	// Get limit from query parameter if provided
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	products, err := h.productAppService.GetPopularProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// TrackProductView handles requests to track a product view
func (h *ProductHandler) TrackProductView(c *gin.Context) {
	id := c.Param("id")
	
	err := h.productAppService.TrackProductView(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// SearchProducts handles requests to search for products
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}
	
	var filter dto.ProductFilterRequest
	
	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	
	products, err := h.productAppService.SearchProducts(c.Request.Context(), query, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, products)
}
