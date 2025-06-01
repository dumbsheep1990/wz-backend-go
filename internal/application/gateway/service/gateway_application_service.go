package service

import (
	"context"
	"fmt"
	"time"

	"wz-backend-go/internal/application/gateway/dto"
	domainService "wz-backend-go/internal/domain/gateway/service"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// GatewayApplicationService handles application-level gateway operations
type GatewayApplicationService struct {
	gatewayDomainService *domainService.GatewayDomainService
	routeRepository      repository.RouteRepository
	serviceRepository    repository.ServiceRepository
}

// NewGatewayApplicationService creates a new GatewayApplicationService
func NewGatewayApplicationService(
	gatewayDomainService *domainService.GatewayDomainService,
	routeRepository repository.RouteRepository,
	serviceRepository repository.ServiceRepository,
) *GatewayApplicationService {
	return &GatewayApplicationService{
		gatewayDomainService: gatewayDomainService,
		routeRepository:      routeRepository,
		serviceRepository:    serviceRepository,
	}
}

// RegisterService registers a new service with the gateway
func (s *GatewayApplicationService) RegisterService(
	ctx context.Context,
	input dto.RegisterServiceInput,
) (*dto.ServiceDTO, error) {
	// Convert string types to value objects
	authType := dto.StringToAuthType(input.DefaultAuth)
	rateLimitType := dto.StringToRateLimitType(input.RateLimitType)
	
	// Register service via domain service
	service, err := s.gatewayDomainService.RegisterService(
		ctx,
		input.Name,
		input.BaseURL,
		input.HealthURL,
		authType,
		rateLimitType,
		input.RateLimitReqs,
		input.RateLimitPeriod,
	)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return dto.ServiceDTOFromEntity(service), nil
}

// UpdateService updates an existing service
func (s *GatewayApplicationService) UpdateService(
	ctx context.Context,
	serviceName string,
	input dto.UpdateServiceInput,
) (*dto.ServiceDTO, error) {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Update base URL if provided
	if input.BaseURL != "" {
		baseURL, err := valueobject.NewTargetURL(input.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		service.UpdateBaseURL(baseURL)
	}
	
	// Update health URL if provided
	if input.HealthURL != "" {
		healthURL, err := valueobject.NewTargetURL(input.HealthURL)
		if err != nil {
			return nil, fmt.Errorf("invalid health URL: %w", err)
		}
		service.UpdateHealthURL(healthURL)
	}
	
	// Update default auth if provided
	if input.DefaultAuth != "" {
		authType := dto.StringToAuthType(input.DefaultAuth)
		service.UpdateDefaultAuth(authType)
	}
	
	// Update rate limit if provided
	if input.RateLimitType != "" {
		rateLimitType := dto.StringToRateLimitType(input.RateLimitType)
		rateLimit, err := valueobject.NewRateLimit(
			rateLimitType,
			input.RateLimitReqs,
			input.RateLimitPeriod,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid rate limit: %w", err)
		}
		service.UpdateRateLimit(rateLimit)
	}
	
	// Update active status if provided
	if input.IsActive != nil {
		if *input.IsActive {
			service.Activate()
		} else {
			service.Deactivate()
		}
	}
	
	// Save service
	if err := s.serviceRepository.Save(ctx, service); err != nil {
		return nil, fmt.Errorf("failed to save service: %w", err)
	}
	
	// Convert to DTO
	return dto.ServiceDTOFromEntity(service), nil
}

// GetService gets a service by name
func (s *GatewayApplicationService) GetService(
	ctx context.Context,
	serviceName string,
) (*dto.ServiceDTO, error) {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Convert to DTO
	return dto.ServiceDTOFromEntity(service), nil
}

// ListServices lists all services with pagination
func (s *GatewayApplicationService) ListServices(
	ctx context.Context,
	page, pageSize int,
) (*dto.ServiceListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Find services
	services, totalCount, err := s.serviceRepository.FindAll(ctx, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find services: %w", err)
	}
	
	// Convert to DTOs
	serviceDTOs := make([]*dto.ServiceDTO, len(services))
	for i, service := range services {
		serviceDTOs[i] = dto.ServiceDTOFromEntity(service)
	}
	
	// Create response
	return &dto.ServiceListResponse{
		Services:    serviceDTOs,
		TotalCount:  totalCount,
		PageSize:    pageSize,
		CurrentPage: page,
	}, nil
}

// DeleteService deletes a service
func (s *GatewayApplicationService) DeleteService(
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
	
	// Check if routes exist for this service
	routes, err := s.routeRepository.FindByServiceName(ctx, serviceNameVO)
	if err != nil {
		return fmt.Errorf("failed to check for routes: %w", err)
	}
	if len(routes) > 0 {
		return fmt.Errorf("cannot delete service '%s' because it has %d routes, delete the routes first", serviceName, len(routes))
	}
	
	// Delete service
	if err := s.serviceRepository.Delete(ctx, serviceNameVO); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	
	return nil
}

// RegisterRoute registers a new route
func (s *GatewayApplicationService) RegisterRoute(
	ctx context.Context,
	input dto.RegisterRouteInput,
) (*dto.RouteDTO, error) {
	// Get service to build target URL
	serviceNameVO, err := valueobject.NewServiceName(input.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found", input.ServiceName)
	}
	
	// Build target URL
	baseURL := service.BaseURL().String()
	targetURL := baseURL
	if input.TargetPath != "" {
		// Remove trailing slash from base URL
		if baseURL[len(baseURL)-1] == '/' {
			baseURL = baseURL[:len(baseURL)-1]
		}
		
		// Add leading slash to target path if needed
		targetPath := input.TargetPath
		if targetPath[0] != '/' {
			targetPath = "/" + targetPath
		}
		
		targetURL = baseURL + targetPath
	}
	
	// Convert string types to value objects
	authType := dto.StringToAuthType(input.AuthType)
	rateLimitType := dto.StringToRateLimitType(input.RateLimitType)
	
	// Register route via domain service
	route, err := s.gatewayDomainService.RegisterRoute(
		ctx,
		input.ID,
		input.Path,
		input.ServiceName,
		targetURL,
		authType,
		rateLimitType,
		input.RateLimitReqs,
		input.RateLimitPeriod,
	)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return dto.RouteDTOFromEntity(route), nil
}

// UpdateRoute updates an existing route
func (s *GatewayApplicationService) UpdateRoute(
	ctx context.Context,
	routeID string,
	input dto.UpdateRouteInput,
) (*dto.RouteDTO, error) {
	// Create route ID value object
	routeIDVO, err := valueobject.NewRouteID(routeID)
	if err != nil {
		return nil, fmt.Errorf("invalid route ID: %w", err)
	}
	
	// Find route
	route, err := s.routeRepository.FindByID(ctx, routeIDVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find route: %w", err)
	}
	if route == nil {
		return nil, fmt.Errorf("route '%s' not found", routeID)
	}
	
	// Update path if provided
	if input.Path != "" {
		path, err := valueobject.NewPath(input.Path)
		if err != nil {
			return nil, fmt.Errorf("invalid path: %w", err)
		}
		route.UpdatePath(path)
	}
	
	// Update target URL if provided
	if input.TargetPath != "" {
		// Get service to build target URL
		service, err := s.serviceRepository.FindByName(ctx, route.ServiceName())
		if err != nil {
			return nil, fmt.Errorf("failed to find service: %w", err)
		}
		if service == nil {
			return nil, fmt.Errorf("service '%s' not found", route.ServiceName().String())
		}
		
		// Build target URL
		baseURL := service.BaseURL().String()
		if baseURL[len(baseURL)-1] == '/' {
			baseURL = baseURL[:len(baseURL)-1]
		}
		
		targetPath := input.TargetPath
		if targetPath[0] != '/' {
			targetPath = "/" + targetPath
		}
		
		targetURL := baseURL + targetPath
		
		// Update target URL
		targetURLVO, err := valueobject.NewTargetURL(targetURL)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL: %w", err)
		}
		route.UpdateTargetURL(targetURLVO)
	}
	
	// Update auth type if provided
	if input.AuthType != "" {
		authType := dto.StringToAuthType(input.AuthType)
		route.UpdateAuthType(authType)
	}
	
	// Update rate limit if provided
	if input.RateLimitType != "" {
		rateLimitType := dto.StringToRateLimitType(input.RateLimitType)
		rateLimit, err := valueobject.NewRateLimit(
			rateLimitType,
			input.RateLimitReqs,
			input.RateLimitPeriod,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid rate limit: %w", err)
		}
		route.UpdateRateLimit(rateLimit)
	}
	
	// Update active status if provided
	if input.IsActive != nil {
		if *input.IsActive {
			route.Activate()
		} else {
			route.Deactivate()
		}
	}
	
	// Save route
	if err := s.routeRepository.Save(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to save route: %w", err)
	}
	
	// Convert to DTO
	return dto.RouteDTOFromEntity(route), nil
}

// GetRoute gets a route by ID
func (s *GatewayApplicationService) GetRoute(
	ctx context.Context,
	routeID string,
) (*dto.RouteDTO, error) {
	// Create route ID value object
	routeIDVO, err := valueobject.NewRouteID(routeID)
	if err != nil {
		return nil, fmt.Errorf("invalid route ID: %w", err)
	}
	
	// Find route
	route, err := s.routeRepository.FindByID(ctx, routeIDVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find route: %w", err)
	}
	if route == nil {
		return nil, fmt.Errorf("route '%s' not found", routeID)
	}
	
	// Convert to DTO
	return dto.RouteDTOFromEntity(route), nil
}

// ListRoutes lists all routes with pagination
func (s *GatewayApplicationService) ListRoutes(
	ctx context.Context,
	page, pageSize int,
) (*dto.RouteListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Find routes
	routes, totalCount, err := s.routeRepository.FindAll(ctx, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find routes: %w", err)
	}
	
	// Convert to DTOs
	routeDTOs := make([]*dto.RouteDTO, len(routes))
	for i, route := range routes {
		routeDTOs[i] = dto.RouteDTOFromEntity(route)
	}
	
	// Create response
	return &dto.RouteListResponse{
		Routes:      routeDTOs,
		TotalCount:  totalCount,
		PageSize:    pageSize,
		CurrentPage: page,
	}, nil
}

// ListServiceRoutes lists all routes for a service
func (s *GatewayApplicationService) ListServiceRoutes(
	ctx context.Context,
	serviceName string,
) ([]*dto.RouteDTO, error) {
	// Create service name value object
	serviceNameVO, err := valueobject.NewServiceName(serviceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// Find service
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found", serviceName)
	}
	
	// Find routes
	routes, err := s.routeRepository.FindByServiceName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find routes: %w", err)
	}
	
	// Convert to DTOs
	routeDTOs := make([]*dto.RouteDTO, len(routes))
	for i, route := range routes {
		routeDTOs[i] = dto.RouteDTOFromEntity(route)
	}
	
	return routeDTOs, nil
}

// DeleteRoute deletes a route
func (s *GatewayApplicationService) DeleteRoute(
	ctx context.Context,
	routeID string,
) error {
	// Create route ID value object
	routeIDVO, err := valueobject.NewRouteID(routeID)
	if err != nil {
		return fmt.Errorf("invalid route ID: %w", err)
	}
	
	// Find route
	route, err := s.routeRepository.FindByID(ctx, routeIDVO)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}
	if route == nil {
		return fmt.Errorf("route '%s' not found", routeID)
	}
	
	// Delete route
	if err := s.routeRepository.Delete(ctx, routeIDVO); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	
	return nil
}

// CheckHealth checks the health of all services
func (s *GatewayApplicationService) CheckHealth(
	ctx context.Context,
) (map[string]bool, error) {
	// Get all services
	services, _, err := s.serviceRepository.FindAll(ctx, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to find services: %w", err)
	}
	
	// Check health of each service
	results := make(map[string]bool)
	for _, service := range services {
		// Skip services without health URL
		if service.HealthURL().URL() == nil {
			continue
		}
		
		err := s.gatewayDomainService.CheckServiceHealth(ctx, service.Name().String())
		results[service.Name().String()] = err == nil
	}
	
	return results, nil
}

// GetGatewayStats gets gateway statistics
func (s *GatewayApplicationService) GetGatewayStats(
	ctx context.Context,
) (*dto.GatewayStatsDTO, error) {
	// Get all services
	allServices, totalServices, err := s.serviceRepository.FindAll(ctx, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to find services: %w", err)
	}
	
	// Get active services
	activeServices, err := s.serviceRepository.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find active services: %w", err)
	}
	
	// Get healthy services
	healthyServices, err := s.serviceRepository.FindHealthy(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find healthy services: %w", err)
	}
	
	// Get all routes
	_, totalRoutes, err := s.routeRepository.FindAll(ctx, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to find routes: %w", err)
	}
	
	// Get active routes
	activeRoutes, err := s.routeRepository.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find active routes: %w", err)
	}
	
	// Create stats DTO
	stats := &dto.GatewayStatsDTO{
		TotalServices:   totalServices,
		ActiveServices:  len(activeServices),
		HealthyServices: len(healthyServices),
		TotalRoutes:     totalRoutes,
		ActiveRoutes:    len(activeRoutes),
		LastUpdated:     time.Now(),
		
		// TODO: These would come from metrics collection in a real implementation
		RequestsPerMinute:   0,
		SuccessfulRequests:  0,
		FailedRequests:      0,
		AverageResponseTime: 0,
		P95ResponseTime:     0,
		RateLimitedRequests: 0,
		AuthFailedRequests:  0,
	}
	
	return stats, nil
}

// WZ 万知网站专属区域路由模板 - 支持"入同"分类
func (s *GatewayApplicationService) CreateWZCategoryRoutes(
	ctx context.Context,
	servicePrefix string,
) ([]*dto.RouteDTO, error) {
	// 21种"入同"分类
	categories := []string{
		"同用", "同好", "同购", "同年", "同游", "同在", "同市", 
		"同企", "同亲", "同班", "同师", "同业", "同网", "同工", 
		"同务", "同艺", "同玩", "同闲", "同拍", "同乡", "同学",
	}
	
	// 创建服务名称
	serviceNameVO, err := valueobject.NewServiceName(servicePrefix + "-category-service")
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	// 检查服务是否存在
	service, err := s.serviceRepository.FindByName(ctx, serviceNameVO)
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("service '%s' not found, please create it first", serviceNameVO.String())
	}
	
	// 为每个分类创建路由
	var routeDTOs []*dto.RouteDTO
	for _, category := range categories {
		// 为分类API创建路由
		apiRoute, err := s.RegisterRoute(ctx, dto.RegisterRouteInput{
			ID:          fmt.Sprintf("%s-api-%s", servicePrefix, category),
			Path:        fmt.Sprintf("/api/v1/categories/%s", category),
			ServiceName: serviceNameVO.String(),
			TargetPath:  fmt.Sprintf("/api/categories/%s", category),
			AuthType:    "jwt",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create API route for category %s: %w", category, err)
		}
		routeDTOs = append(routeDTOs, apiRoute)
		
		// 为分类页面创建路由
		pageRoute, err := s.RegisterRoute(ctx, dto.RegisterRouteInput{
			ID:          fmt.Sprintf("%s-page-%s", servicePrefix, category),
			Path:        fmt.Sprintf("/%s", category),
			ServiceName: serviceNameVO.String(),
			TargetPath:  fmt.Sprintf("/categories/%s", category),
			AuthType:    "none",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create page route for category %s: %w", category, err)
		}
		routeDTOs = append(routeDTOs, pageRoute)
	}
	
	return routeDTOs, nil
}
