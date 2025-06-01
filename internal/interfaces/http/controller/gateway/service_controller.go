package gateway

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/gateway/dto"
	"wz-backend-go/internal/application/gateway/service"
)

// ServiceController handles HTTP requests for services
type ServiceController struct {
	gatewayAppService *service.GatewayApplicationService
}

// NewServiceController creates a new ServiceController
func NewServiceController(gatewayAppService *service.GatewayApplicationService) *ServiceController {
	return &ServiceController{
		gatewayAppService: gatewayAppService,
	}
}

// RegisterRoutes registers routes for the controller
func (c *ServiceController) RegisterRoutes(router *gin.RouterGroup) {
	services := router.Group("/services")
	{
		services.GET("", c.ListServices)
		services.POST("", c.RegisterService)
		services.GET("/:name", c.GetService)
		services.PUT("/:name", c.UpdateService)
		services.DELETE("/:name", c.DeleteService)
		services.GET("/:name/routes", c.ListServiceRoutes)
		services.POST("/:name/check-health", c.CheckServiceHealth)
		services.POST("/check-health", c.CheckAllServicesHealth)
		services.POST("/wz-categories", c.CreateWZCategoryRoutes)
	}
}

// RegisterService godoc
// @Summary Register a new service
// @Description Register a new service with the gateway
// @Tags gateway,services
// @Accept json
// @Produce json
// @Param service body dto.RegisterServiceInput true "Service registration details"
// @Success 201 {object} dto.ServiceDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services [post]
func (c *ServiceController) RegisterService(ctx *gin.Context) {
	var input dto.RegisterServiceInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid input format",
			Details: err.Error(),
		})
		return
	}
	
	service, err := c.gatewayAppService.RegisterService(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_REGISTRATION_FAILED",
			Message: "Failed to register service",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusCreated, service)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service
// @Tags gateway,services
// @Accept json
// @Produce json
// @Param name path string true "Service name"
// @Param service body dto.UpdateServiceInput true "Service update details"
// @Success 200 {object} dto.ServiceDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/{name} [put]
func (c *ServiceController) UpdateService(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service name is required",
		})
		return
	}
	
	var input dto.UpdateServiceInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid input format",
			Details: err.Error(),
		})
		return
	}
	
	service, err := c.gatewayAppService.UpdateService(ctx, name, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_UPDATE_FAILED",
			Message: "Failed to update service",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, service)
}

// GetService godoc
// @Summary Get a service
// @Description Get a service by name
// @Tags gateway,services
// @Produce json
// @Param name path string true "Service name"
// @Success 200 {object} dto.ServiceDTO
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/{name} [get]
func (c *ServiceController) GetService(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service name is required",
		})
		return
	}
	
	service, err := c.gatewayAppService.GetService(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    "SERVICE_NOT_FOUND",
			Message: "Service not found",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, service)
}

// ListServices godoc
// @Summary List services
// @Description List all services with pagination
// @Tags gateway,services
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.ServiceListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services [get]
func (c *ServiceController) ListServices(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	services, err := c.gatewayAppService.ListServices(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_LIST_FAILED",
			Message: "Failed to list services",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, services)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service by name
// @Tags gateway,services
// @Produce json
// @Param name path string true "Service name"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/{name} [delete]
func (c *ServiceController) DeleteService(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service name is required",
		})
		return
	}
	
	if err := c.gatewayAppService.DeleteService(ctx, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_DELETE_FAILED",
			Message: "Failed to delete service",
			Details: err.Error(),
		})
		return
	}
	
	ctx.Status(http.StatusNoContent)
}

// ListServiceRoutes godoc
// @Summary List service routes
// @Description List all routes for a service
// @Tags gateway,services,routes
// @Produce json
// @Param name path string true "Service name"
// @Success 200 {array} dto.RouteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/{name}/routes [get]
func (c *ServiceController) ListServiceRoutes(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service name is required",
		})
		return
	}
	
	routes, err := c.gatewayAppService.ListServiceRoutes(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_ROUTES_LIST_FAILED",
			Message: "Failed to list service routes",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, routes)
}

// CheckServiceHealth godoc
// @Summary Check service health
// @Description Check the health of a service
// @Tags gateway,services
// @Produce json
// @Param name path string true "Service name"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/{name}/check-health [post]
func (c *ServiceController) CheckServiceHealth(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service name is required",
		})
		return
	}
	
	if err := c.gatewayAppService.CheckServiceHealth(ctx, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_HEALTH_CHECK_FAILED",
			Message: "Service health check failed",
			Details: err.Error(),
		})
		return
	}
	
	// Get updated service
	service, err := c.gatewayAppService.GetService(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICE_GET_FAILED",
			Message: "Failed to get service after health check",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"service_name": name,
		"is_healthy":   service.IsHealthy,
		"error":        service.ErrorMessage,
	})
}

// CheckAllServicesHealth godoc
// @Summary Check all services health
// @Description Check the health of all services
// @Tags gateway,services
// @Produce json
// @Success 200 {object} map[string]bool
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/check-health [post]
func (c *ServiceController) CheckAllServicesHealth(ctx *gin.Context) {
	results, err := c.gatewayAppService.CheckHealth(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "SERVICES_HEALTH_CHECK_FAILED",
			Message: "Services health check failed",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, results)
}

// CreateWZCategoryRoutes godoc
// @Summary Create 万知 category routes
// @Description Create routes for all 21 万知 "入同" categories
// @Tags gateway,services,routes,万知
// @Produce json
// @Param prefix query string true "Service prefix for route IDs"
// @Success 200 {array} dto.RouteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/services/wz-categories [post]
func (c *ServiceController) CreateWZCategoryRoutes(ctx *gin.Context) {
	prefix := ctx.Query("prefix")
	if prefix == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Service prefix is required",
		})
		return
	}
	
	routes, err := c.gatewayAppService.CreateWZCategoryRoutes(ctx, prefix)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "CATEGORY_ROUTES_CREATION_FAILED",
			Message: "Failed to create category routes",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, routes)
}
