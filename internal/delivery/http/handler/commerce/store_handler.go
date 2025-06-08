package commerce

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"wz-backend-go/internal/application/commerce"
	"wz-backend-go/internal/domain/commerce/dto"
)

// StoreHandler handles HTTP requests for commerce stores
type StoreHandler struct {
	storeAppService *commerce.StoreAppService
}

// NewStoreHandler creates a new instance of StoreHandler
func NewStoreHandler(storeAppService *commerce.StoreAppService) *StoreHandler {
	return &StoreHandler{
		storeAppService: storeAppService,
	}
}

// CreateStore handles requests to create a new store
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req dto.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, err := h.storeAppService.CreateStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// UpdateStore handles requests to update an existing store
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var req dto.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the one in the request body
	if c.Param("id") != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "store ID mismatch"})
		return
	}

	store, err := h.storeAppService.UpdateStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}

// GetStoreByID handles requests to retrieve a store by its ID
func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	id := c.Param("id")

	store, err := h.storeAppService.GetStoreByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if store == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}

// GetStoresByOwner handles requests to retrieve stores by owner ID
func (h *StoreHandler) GetStoresByOwner(c *gin.Context) {
	ownerID := c.Param("ownerID")
	
	stores, err := h.storeAppService.GetStoresByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"stores": stores})
}

// GetStores handles requests to retrieve stores with filtering
func (h *StoreHandler) GetStores(c *gin.Context) {
	var filter dto.StoreFilterRequest
	
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
	
	stores, err := h.storeAppService.GetStores(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stores)
}

// SearchStores handles requests to search for stores
func (h *StoreHandler) SearchStores(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}
	
	stores, err := h.storeAppService.SearchStores(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"stores": stores})
}
