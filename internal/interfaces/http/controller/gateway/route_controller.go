package gateway

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/gateway/dto"
	"wz-backend-go/internal/application/gateway/service"
)

// RouteController handles HTTP requests for routes
type RouteController struct {
	gatewayAppService *service.GatewayApplicationService
}

// NewRouteController creates a new RouteController
func NewRouteController(gatewayAppService *service.GatewayApplicationService) *RouteController {
	return &RouteController{
		gatewayAppService: gatewayAppService,
	}
}

// RegisterRoutes registers routes for the controller
func (c *RouteController) RegisterRoutes(router *gin.RouterGroup) {
	routes := router.Group("/routes")
	{
		routes.GET("", c.ListRoutes)
		routes.POST("", c.RegisterRoute)
		routes.GET("/:id", c.GetRoute)
		routes.PUT("/:id", c.UpdateRoute)
		routes.DELETE("/:id", c.DeleteRoute)
	}
}

// RegisterRoute godoc
// @Summary Register a new route
// @Description Register a new route with the gateway
// @Tags gateway,routes
// @Accept json
// @Produce json
// @Param route body dto.RegisterRouteInput true "Route registration details"
// @Success 201 {object} dto.RouteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/routes [post]
func (c *RouteController) RegisterRoute(ctx *gin.Context) {
	var input dto.RegisterRouteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid input format",
			Details: err.Error(),
		})
		return
	}
	
	route, err := c.gatewayAppService.RegisterRoute(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "ROUTE_REGISTRATION_FAILED",
			Message: "Failed to register route",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusCreated, route)
}

// UpdateRoute godoc
// @Summary Update a route
// @Description Update an existing route
// @Tags gateway,routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID"
// @Param route body dto.UpdateRouteInput true "Route update details"
// @Success 200 {object} dto.RouteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/routes/{id} [put]
func (c *RouteController) UpdateRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Route ID is required",
		})
		return
	}
	
	var input dto.UpdateRouteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid input format",
			Details: err.Error(),
		})
		return
	}
	
	route, err := c.gatewayAppService.UpdateRoute(ctx, id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "ROUTE_UPDATE_FAILED",
			Message: "Failed to update route",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, route)
}

// GetRoute godoc
// @Summary Get a route
// @Description Get a route by ID
// @Tags gateway,routes
// @Produce json
// @Param id path string true "Route ID"
// @Success 200 {object} dto.RouteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/routes/{id} [get]
func (c *RouteController) GetRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Route ID is required",
		})
		return
	}
	
	route, err := c.gatewayAppService.GetRoute(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    "ROUTE_NOT_FOUND",
			Message: "Route not found",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, route)
}

// ListRoutes godoc
// @Summary List routes
// @Description List all routes with pagination
// @Tags gateway,routes
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.RouteListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/routes [get]
func (c *RouteController) ListRoutes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	routes, err := c.gatewayAppService.ListRoutes(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "ROUTE_LIST_FAILED",
			Message: "Failed to list routes",
			Details: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, routes)
}

// DeleteRoute godoc
// @Summary Delete a route
// @Description Delete a route by ID
// @Tags gateway,routes
// @Produce json
// @Param id path string true "Route ID"
// @Success 204 "No content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/gateway/routes/{id} [delete]
func (c *RouteController) DeleteRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Route ID is required",
		})
		return
	}
	
	if err := c.gatewayAppService.DeleteRoute(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "ROUTE_DELETE_FAILED",
			Message: "Failed to delete route",
			Details: err.Error(),
		})
		return
	}
	
	ctx.Status(http.StatusNoContent)
}
