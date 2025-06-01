package dto

import (
	"time"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// -----------------------------------------------
// Input DTOs for request validation and mapping
// -----------------------------------------------

// RegisterServiceInput represents input for registering a service
type RegisterServiceInput struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	BaseURL         string `json:"base_url" binding:"required,url"`
	HealthURL       string `json:"health_url" binding:"omitempty,url"`
	DefaultAuth     string `json:"default_auth" binding:"omitempty,oneof=none jwt apikey oauth"`
	RateLimitType   string `json:"rate_limit_type" binding:"omitempty,oneof=none ip token service"`
	RateLimitReqs   int    `json:"rate_limit_requests" binding:"omitempty,min=1"`
	RateLimitPeriod int    `json:"rate_limit_period_seconds" binding:"omitempty,min=1"`
}

// UpdateServiceInput represents input for updating a service
type UpdateServiceInput struct {
	BaseURL         string `json:"base_url" binding:"omitempty,url"`
	HealthURL       string `json:"health_url" binding:"omitempty,url"`
	DefaultAuth     string `json:"default_auth" binding:"omitempty,oneof=none jwt apikey oauth"`
	RateLimitType   string `json:"rate_limit_type" binding:"omitempty,oneof=none ip token service"`
	RateLimitReqs   int    `json:"rate_limit_requests" binding:"omitempty,min=1"`
	RateLimitPeriod int    `json:"rate_limit_period_seconds" binding:"omitempty,min=1"`
	IsActive        *bool  `json:"is_active" binding:"omitempty"`
}

// RegisterRouteInput represents input for registering a route
type RegisterRouteInput struct {
	ID              string `json:"id" binding:"required,min=1,max=100"`
	Path            string `json:"path" binding:"required,min=1,max=500"`
	ServiceName     string `json:"service_name" binding:"required,min=1,max=100"`
	TargetPath      string `json:"target_path" binding:"omitempty"`
	AuthType        string `json:"auth_type" binding:"omitempty,oneof=none jwt apikey oauth"`
	RateLimitType   string `json:"rate_limit_type" binding:"omitempty,oneof=none ip token service"`
	RateLimitReqs   int    `json:"rate_limit_requests" binding:"omitempty,min=1"`
	RateLimitPeriod int    `json:"rate_limit_period_seconds" binding:"omitempty,min=1"`
}

// UpdateRouteInput represents input for updating a route
type UpdateRouteInput struct {
	Path            string `json:"path" binding:"omitempty,min=1,max=500"`
	TargetPath      string `json:"target_path" binding:"omitempty"`
	AuthType        string `json:"auth_type" binding:"omitempty,oneof=none jwt apikey oauth"`
	RateLimitType   string `json:"rate_limit_type" binding:"omitempty,oneof=none ip token service"`
	RateLimitReqs   int    `json:"rate_limit_requests" binding:"omitempty,min=1"`
	RateLimitPeriod int    `json:"rate_limit_period_seconds" binding:"omitempty,min=1"`
	IsActive        *bool  `json:"is_active" binding:"omitempty"`
}

// -----------------------------------------------
// Output DTOs for responses
// -----------------------------------------------

// ServiceDTO represents a service in responses
type ServiceDTO struct {
	Name           string    `json:"name"`
	BaseURL        string    `json:"base_url"`
	HealthURL      string    `json:"health_url,omitempty"`
	DefaultAuth    string    `json:"default_auth"`
	RateLimitType  string    `json:"rate_limit_type,omitempty"`
	RateLimitReqs  int       `json:"rate_limit_requests,omitempty"`
	RateLimitSecs  int       `json:"rate_limit_period_seconds,omitempty"`
	IsActive       bool      `json:"is_active"`
	IsHealthy      bool      `json:"is_healthy"`
	LastHealthy    time.Time `json:"last_healthy,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// FromService converts a service entity to DTO
func ServiceDTOFromEntity(service *entity.Service) *ServiceDTO {
	dto := &ServiceDTO{
		Name:        service.Name().String(),
		BaseURL:     service.BaseURL().String(),
		DefaultAuth: service.DefaultAuth().String(),
		IsActive:    service.IsActive(),
		IsHealthy:   service.IsHealthy(),
		CreatedAt:   service.CreatedAt(),
		UpdatedAt:   service.UpdatedAt(),
	}
	
	// Add health URL if available
	if service.HealthURL().URL() != nil {
		dto.HealthURL = service.HealthURL().String()
	}
	
	// Add rate limit if available
	if service.RateLimit().IsEnabled() {
		dto.RateLimitType = service.RateLimit().Type().String()
		dto.RateLimitReqs = service.RateLimit().Requests()
		dto.RateLimitSecs = int(service.RateLimit().Duration().Seconds())
	}
	
	// Add last healthy time if available
	if service.LastHealthy() != nil {
		dto.LastHealthy = *service.LastHealthy()
	}
	
	// Add error message if available
	if service.ErrorMessage() != "" {
		dto.ErrorMessage = service.ErrorMessage()
	}
	
	return dto
}

// RouteDTO represents a route in responses
type RouteDTO struct {
	ID             string    `json:"id"`
	Path           string    `json:"path"`
	ServiceName    string    `json:"service_name"`
	TargetURL      string    `json:"target_url"`
	AuthType       string    `json:"auth_type"`
	RateLimitType  string    `json:"rate_limit_type,omitempty"`
	RateLimitReqs  int       `json:"rate_limit_requests,omitempty"`
	RateLimitSecs  int       `json:"rate_limit_period_seconds,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// FromRoute converts a route entity to DTO
func RouteDTOFromEntity(route *entity.Route) *RouteDTO {
	dto := &RouteDTO{
		ID:          route.ID().String(),
		Path:        route.Path().String(),
		ServiceName: route.ServiceName().String(),
		TargetURL:   route.TargetURL().String(),
		AuthType:    route.AuthType().String(),
		IsActive:    route.IsActive(),
		CreatedAt:   route.CreatedAt(),
		UpdatedAt:   route.UpdatedAt(),
	}
	
	// Add rate limit if available
	if route.RateLimit().IsEnabled() {
		dto.RateLimitType = route.RateLimit().Type().String()
		dto.RateLimitReqs = route.RateLimit().Requests()
		dto.RateLimitSecs = int(route.RateLimit().Duration().Seconds())
	}
	
	return dto
}

// ServiceListResponse represents a paginated list of services
type ServiceListResponse struct {
	Services    []*ServiceDTO `json:"services"`
	TotalCount  int           `json:"total_count"`
	PageSize    int           `json:"page_size"`
	CurrentPage int           `json:"current_page"`
}

// RouteListResponse represents a paginated list of routes
type RouteListResponse struct {
	Routes      []*RouteDTO `json:"routes"`
	TotalCount  int         `json:"total_count"`
	PageSize    int         `json:"page_size"`
	CurrentPage int         `json:"current_page"`
}

// GatewayStatsDTO represents gateway statistics
type GatewayStatsDTO struct {
	TotalServices        int `json:"total_services"`
	ActiveServices       int `json:"active_services"`
	HealthyServices      int `json:"healthy_services"`
	TotalRoutes          int `json:"total_routes"`
	ActiveRoutes         int `json:"active_routes"`
	RequestsPerMinute    int `json:"requests_per_minute"`
	SuccessfulRequests   int `json:"successful_requests"`
	FailedRequests       int `json:"failed_requests"`
	AverageResponseTime  int `json:"average_response_time_ms"`
	P95ResponseTime      int `json:"p95_response_time_ms"`
	RateLimitedRequests  int `json:"rate_limited_requests"`
	AuthFailedRequests   int `json:"auth_failed_requests"`
	LastUpdated          time.Time `json:"last_updated"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Helper functions to convert between value objects and strings

// StringToAuthType converts a string to an AuthType
func StringToAuthType(authType string) valueobject.AuthType {
	switch authType {
	case "jwt":
		return valueobject.JWTAuth
	case "apikey":
		return valueobject.APIKeyAuth
	case "oauth":
		return valueobject.OAuthAuth
	default:
		return valueobject.NoAuth
	}
}

// StringToRateLimitType converts a string to a RateLimitType
func StringToRateLimitType(limitType string) valueobject.RateLimitType {
	switch limitType {
	case "ip":
		return valueobject.IPRateLimit
	case "token":
		return valueobject.TokenRateLimit
	case "service":
		return valueobject.ServiceRateLimit
	default:
		return valueobject.NoRateLimit
	}
}
