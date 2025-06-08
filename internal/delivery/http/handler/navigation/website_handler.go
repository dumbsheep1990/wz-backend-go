package navigation

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/application/navigation"
	"wz-backend-go/internal/domain/navigation/dto"
)

// WebsiteHandler handles HTTP requests for websites in the navigation system
type WebsiteHandler struct {
	navigationAppService *navigation.NavigationAppService
}

// NewWebsiteHandler creates a new instance of WebsiteHandler
func NewWebsiteHandler(navigationAppService *navigation.NavigationAppService) *WebsiteHandler {
	return &WebsiteHandler{
		navigationAppService: navigationAppService,
	}
}

// CreateWebsite handles requests to create a new website
func (h *WebsiteHandler) CreateWebsite(c *gin.Context) {
	var req dto.CreateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	website, err := h.navigationAppService.CreateWebsite(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, website)
}

// UpdateWebsite handles requests to update an existing website
func (h *WebsiteHandler) UpdateWebsite(c *gin.Context) {
	var req dto.UpdateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the one in the request body
	if c.Param("id") != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "website ID mismatch"})
		return
	}

	website, err := h.navigationAppService.UpdateWebsite(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, website)
}

// GetWebsites handles requests to retrieve websites with optional filtering
func (h *WebsiteHandler) GetWebsites(c *gin.Context) {
	categoryID := c.Query("category_id")
	featuredOnly := c.Query("featured") == "true"

	websites, err := h.navigationAppService.GetWebsites(c.Request.Context(), categoryID, featuredOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, websites)
}

// GetPopularWebsites handles requests to retrieve popular websites
func (h *WebsiteHandler) GetPopularWebsites(c *gin.Context) {
	limit := 10 // Default limit
	// Parse limit from query parameter if provided
	if c.Query("limit") != "" {
		var err error
		limitParam, err := c.GetQuery("limit")
		if err == nil {
			// Convert limit to int
			var parsedLimit int
			_, err = fmt.Sscan(limitParam, &parsedLimit)
			if err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}
	}

	websites, err := h.navigationAppService.GetPopularWebsites(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, websites)
}

// TrackWebsiteView handles requests to track a website view
func (h *WebsiteHandler) TrackWebsiteView(c *gin.Context) {
	id := c.Param("id")

	err := h.navigationAppService.TrackWebsiteView(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ReorderWebsites handles requests to reorder websites
func (h *WebsiteHandler) ReorderWebsites(c *gin.Context) {
	var req dto.ReorderWebsitesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.navigationAppService.ReorderWebsites(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteWebsite handles requests to delete a website
func (h *WebsiteHandler) DeleteWebsite(c *gin.Context) {
	id := c.Param("id")
	
	// In a real implementation, you would have a delete method in the app service
	// For now, we'll return a not implemented error
	c.JSON(http.StatusNotImplemented, gin.H{"error": "delete operation not implemented"})
}
