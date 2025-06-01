package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"wz-backend-go/internal/domain/event"
	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// GatewayDomainService handles domain-level gateway operations
type GatewayDomainService struct {
	routeRepository   repository.RouteRepository
	serviceRepository repository.ServiceRepository
	eventBus          event.EventBus
}

// NewGatewayDomainService creates a new GatewayDomainService
func NewGatewayDomainService(
	routeRepository repository.RouteRepository,
	serviceRepository repository.ServiceRepository,
	eventBus event.EventBus,
) *GatewayDomainService {
	return &GatewayDomainService{
		routeRepository:   routeRepository,
		serviceRepository: serviceRepository,
		eventBus:          eventBus,
	}
}

// RegisterService registers a new service with the gateway
func (s *GatewayDomainService) RegisterService(
	ctx context.Context,
	name string,
	baseURL string,
	healthURL string,
	defaultAuth valueobject.AuthType,
	rateLimitType valueobject.RateLimitType,
	rateLimitRequests int,
	rateLimitDuration int,
) (*entity.Service, error) {
	// Create service name value object
	serviceName, err := valueobject.NewServiceName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Check if service already exists
	existingService, err := s.serviceRepository.FindByName(ctx, serviceName)
	if err == nil && existingService != nil {
		return nil, fmt.Errorf("service with name '%s' already exists", name)
	}
	
	// Create service builder
	builder := entity.NewServiceBuilder().
		WithName(name).
		WithBaseURL(baseURL).
		WithDefaultAuth(defaultAuth)
	
	// Add health URL if provided
	if healthURL != "" {
		builder.WithHealthURL(healthURL)
	}
	
	// Add rate limit if provided
	if rateLimitType != valueobject.NoRateLimit {
		builder.WithRateLimit(rateLimitType, rateLimitRequests, rateLimitDuration)
	}
	
	// Build service entity
	service, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	
	// Save service
	if err := s.serviceRepository.Save(ctx, service); err != nil {
		return nil, fmt.Errorf("failed to save service: %w", err)
	}
	
	// Publish event
	if s.eventBus != nil {
		if err := s.eventBus.Publish(ctx, service.ServiceCreatedEvent()); err != nil {
			// Log but don't fail
			fmt.Printf("failed to publish service created event: %v\n", err)
		}
	}
	
	return service, nil
}

// UpdateServiceHealth updates the health status of a service
func (s *GatewayDomainService) UpdateServiceHealth(
	ctx context.Context,
	serviceName string,
	isHealthy bool,
	errorMessage string,
) error {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Check if health status changed
	healthChanged := service.IsHealthy() != isHealthy
	
	// Update health status
	if isHealthy {
		service.SetHealthy()
	} else {
		service.SetUnhealthy(errorMessage)
	}
	
	// Save service
	if err := s.serviceRepository.Save(ctx, service); err != nil {
		return fmt.Errorf("failed to save service health status: %w", err)
	}
	
	// Publish event if health status changed
	if healthChanged && s.eventBus != nil {
		if err := s.eventBus.Publish(ctx, service.ServiceHealthStatusChangedEvent()); err != nil {
			// Log but don't fail
			fmt.Printf("failed to publish service health status changed event: %v\n", err)
		}
	}
	
	return nil
}

// CheckServiceHealth performs a health check for a service
func (s *GatewayDomainService) CheckServiceHealth(
	ctx context.Context,
	serviceName string,
) error {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Check if health URL is configured
	healthURL := service.HealthURL()
	if healthURL.URL() == nil {
		return fmt.Errorf("service '%s' does not have a health URL configured", serviceName)
	}
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	// Perform health check
	resp, err := client.Get(healthURL.String())
	if err != nil {
		// Update service as unhealthy
		s.UpdateServiceHealth(ctx, serviceName, false, err.Error())
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("health check returned non-2xx status code: %d", resp.StatusCode)
		// Update service as unhealthy
		s.UpdateServiceHealth(ctx, serviceName, false, errMsg)
		return errors.New(errMsg)
	}
	
	// Update service as healthy
	s.UpdateServiceHealth(ctx, serviceName, true, "")
	return nil
}

// RegisterRoute registers a new route with the gateway
func (s *GatewayDomainService) RegisterRoute(
	ctx context.Context,
	routeID string,
	path string,
	serviceName string,
	targetURL string,
	authType valueobject.AuthType,
	rateLimitType valueobject.RateLimitType,
	rateLimitRequests int,
	rateLimitDuration int,
) (*entity.Route, error) {
	// Create value objects
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Check if service exists
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Create route ID value object
	routeIDVO, err := valueobject.NewRouteID(routeID)
	if err != nil {
		return nil, fmt.Errorf("invalid route ID: %w", err)
	}
	
	// Check if route already exists
	existingRoute, err := s.routeRepository.FindByID(ctx, routeIDVO)
	if err == nil && existingRoute != nil {
		return nil, fmt.Errorf("route with ID '%s' already exists", routeID)
	}
	
	// Create route builder
	builder := entity.NewRouteBuilder().
		WithID(routeID).
		WithPath(path).
		WithServiceName(serviceName).
		WithTargetURL(targetURL).
		WithAuthType(authType)
	
	// Add rate limit if provided
	if rateLimitType != valueobject.NoRateLimit {
		builder.WithRateLimit(rateLimitType, rateLimitRequests, rateLimitDuration)
	}
	
	// Build route entity
	route, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}
	
	// Save route
	if err := s.routeRepository.Save(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to save route: %w", err)
	}
	
	// Publish event
	if s.eventBus != nil {
		if err := s.eventBus.Publish(ctx, route.RouteCreatedEvent()); err != nil {
			// Log but don't fail
			fmt.Printf("failed to publish route created event: %v\n", err)
		}
	}
	
	return route, nil
}

// FindMatchingRoutes finds all routes that match a request path
func (s *GatewayDomainService) FindMatchingRoutes(
	ctx context.Context,
	requestPath string,
) ([]*entity.Route, error) {
	// Get all active routes
	routes, err := s.routeRepository.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find active routes: %w", err)
	}
	
	// Filter routes that match the request path
	var matchingRoutes []*entity.Route
	for _, route := range routes {
		if route.Matches(requestPath) {
			matchingRoutes = append(matchingRoutes, route)
		}
	}
	
	return matchingRoutes, nil
}

// CheckServiceAvailability checks if a service is available
func (s *GatewayDomainService) CheckServiceAvailability(
	ctx context.Context,
	serviceName string,
) (bool, error) {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return false, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return false, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return false, fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Check if service is active and healthy
	return service.IsActive() && service.IsHealthy(), nil
}
