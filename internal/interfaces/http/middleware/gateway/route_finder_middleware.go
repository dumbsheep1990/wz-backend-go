package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/gateway/dto"
	"wz-backend-go/internal/domain/gateway/service"
)

// RouteFinder finds matching routes for requests
type RouteFinder struct {
	gatewayDomainService *service.GatewayDomainService
}

// NewRouteFinder creates a new RouteFinder
func NewRouteFinder(gatewayDomainService *service.GatewayDomainService) *RouteFinder {
	return &RouteFinder{
		gatewayDomainService: gatewayDomainService,
	}
}

// FindRoute returns a handler function to find a matching route
func (f *RouteFinder) FindRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request path
		path := c.Request.URL.Path

		// Find matching route
		route, err := f.gatewayDomainService.FindRouteByPath(c, path)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "ROUTE_FINDER_ERROR",
				Message: "Failed to find route",
				Details: err.Error(),
			})
			return
		}

		// If no route found, continue to next handler (likely 404)
		if route == nil {
			c.Next()
			return
		}

		// Check if route is active
		if !route.IsActive() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, dto.ErrorResponse{
				Code:    "ROUTE_INACTIVE",
				Message: "This route is currently inactive",
			})
			return
		}

		// Check if service is available
		serviceName := route.ServiceName().String()
		isAvailable, err := f.gatewayDomainService.IsServiceAvailable(c, serviceName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "SERVICE_CHECK_ERROR",
				Message: "Failed to check service availability",
				Details: err.Error(),
			})
			return
		}

		if !isAvailable {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, dto.ErrorResponse{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "The service for this route is currently unavailable",
			})
			return
		}

		// Store route in context for other middlewares
		c.Set("route", route)
		c.Set("target_url", route.TargetURL().String())
		c.Next()
	}
}
