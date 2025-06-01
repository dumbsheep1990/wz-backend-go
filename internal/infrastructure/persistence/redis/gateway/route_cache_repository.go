package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	
	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

const (
	routeKeyPrefix     = "gateway:route:"
	routeServicePrefix = "gateway:service-routes:"
	routePathPrefix    = "gateway:path-routes:"
	routeActivesKey    = "gateway:active-routes"
	routeCacheTTL      = 30 * time.Minute
)

// RouteDTO is a data transfer object for serializing routes to Redis
type RouteDTO struct {
	ID                     string    `json:"id"`
	Path                   string    `json:"path"`
	ServiceName            string    `json:"service_name"`
	TargetURL              string    `json:"target_url"`
	AuthType               string    `json:"auth_type"`
	RateLimitType          string    `json:"rate_limit_type"`
	RateLimitRequests      int       `json:"rate_limit_requests"`
	RateLimitDurationSecs  int       `json:"rate_limit_duration_secs"`
	IsActive               bool      `json:"is_active"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// RedisRouteCacheRepository implements a caching layer for routes using Redis
type RedisRouteCacheRepository struct {
	client redis.Client
}

// NewRedisRouteCacheRepository creates a new RedisRouteCacheRepository
func NewRedisRouteCacheRepository(client redis.Client) *RedisRouteCacheRepository {
	return &RedisRouteCacheRepository{
		client: client,
	}
}

// SaveRoute caches a route in Redis
func (r *RedisRouteCacheRepository) SaveRoute(ctx context.Context, route *entity.Route) error {
	// Convert route to DTO for serialization
	dto := r.toDTO(route)
	
	// Serialize DTO to JSON
	routeJSON, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal route: %w", err)
	}
	
	// Create pipeline for multiple operations
	pipe := r.client.Pipeline()
	
	// Store route by ID
	routeKey := fmt.Sprintf("%s%s", routeKeyPrefix, route.ID().String())
	pipe.Set(ctx, routeKey, routeJSON, routeCacheTTL)
	
	// Store route ID in service routes set
	serviceKey := fmt.Sprintf("%s%s", routeServicePrefix, route.ServiceName().String())
	pipe.SAdd(ctx, serviceKey, route.ID().String())
	pipe.Expire(ctx, serviceKey, routeCacheTTL)
	
	// Store route ID in path routes set if active
	if route.IsActive() {
		pathKey := fmt.Sprintf("%s%s", routePathPrefix, route.Path().String())
		pipe.SAdd(ctx, pathKey, route.ID().String())
		pipe.Expire(ctx, pathKey, routeCacheTTL)
		
		// Add to active routes set
		pipe.SAdd(ctx, routeActivesKey, route.ID().String())
		pipe.Expire(ctx, routeActivesKey, routeCacheTTL)
	} else {
		// Remove from active routes set if not active
		pipe.SRem(ctx, routeActivesKey, route.ID().String())
		
		// Remove from path routes set
		pathKey := fmt.Sprintf("%s%s", routePathPrefix, route.Path().String())
		pipe.SRem(ctx, pathKey, route.ID().String())
	}
	
	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save route to cache: %w", err)
	}
	
	return nil
}

// FindRouteByID finds a route by ID from cache
func (r *RedisRouteCacheRepository) FindRouteByID(ctx context.Context, id valueobject.RouteID) (*entity.Route, error) {
	routeKey := fmt.Sprintf("%s%s", routeKeyPrefix, id.String())
	routeJSON, err := r.client.Get(ctx, routeKey).Result()
	
	if err == redis.Nil {
		return nil, nil // Not found in cache
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get route from cache: %w", err)
	}
	
	return r.fromJSON(routeJSON)
}

// FindRoutesByPath finds routes by path from cache
func (r *RedisRouteCacheRepository) FindRoutesByPath(ctx context.Context, path valueobject.Path) ([]*entity.Route, error) {
	pathKey := fmt.Sprintf("%s%s", routePathPrefix, path.String())
	routeIDs, err := r.client.SMembers(ctx, pathKey).Result()
	
	if err == redis.Nil || len(routeIDs) == 0 {
		return nil, nil // No routes found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get routes by path from cache: %w", err)
	}
	
	return r.getRoutesByIDs(ctx, routeIDs)
}

// FindRoutesByServiceName finds routes by service name from cache
func (r *RedisRouteCacheRepository) FindRoutesByServiceName(ctx context.Context, serviceName valueobject.ServiceName) ([]*entity.Route, error) {
	serviceKey := fmt.Sprintf("%s%s", routeServicePrefix, serviceName.String())
	routeIDs, err := r.client.SMembers(ctx, serviceKey).Result()
	
	if err == redis.Nil || len(routeIDs) == 0 {
		return nil, nil // No routes found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get routes by service name from cache: %w", err)
	}
	
	return r.getRoutesByIDs(ctx, routeIDs)
}

// FindActiveRoutes finds all active routes from cache
func (r *RedisRouteCacheRepository) FindActiveRoutes(ctx context.Context) ([]*entity.Route, error) {
	routeIDs, err := r.client.SMembers(ctx, routeActivesKey).Result()
	
	if err == redis.Nil || len(routeIDs) == 0 {
		return nil, nil // No active routes found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active routes from cache: %w", err)
	}
	
	return r.getRoutesByIDs(ctx, routeIDs)
}

// DeleteRoute removes a route from cache
func (r *RedisRouteCacheRepository) DeleteRoute(ctx context.Context, id valueobject.RouteID) error {
	// First get the route to know its service name and path
	route, err := r.FindRouteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get route for deletion: %w", err)
	}
	
	if route == nil {
		return nil // Already not in cache
	}
	
	// Create pipeline for multiple operations
	pipe := r.client.Pipeline()
	
	// Remove route by ID
	routeKey := fmt.Sprintf("%s%s", routeKeyPrefix, id.String())
	pipe.Del(ctx, routeKey)
	
	// Remove from service routes set
	serviceKey := fmt.Sprintf("%s%s", routeServicePrefix, route.ServiceName().String())
	pipe.SRem(ctx, serviceKey, id.String())
	
	// Remove from path routes set
	pathKey := fmt.Sprintf("%s%s", routePathPrefix, route.Path().String())
	pipe.SRem(ctx, pathKey, id.String())
	
	// Remove from active routes set
	pipe.SRem(ctx, routeActivesKey, id.String())
	
	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete route from cache: %w", err)
	}
	
	return nil
}

// InvalidateServiceRoutes removes all routes for a service from cache
func (r *RedisRouteCacheRepository) InvalidateServiceRoutes(ctx context.Context, serviceName valueobject.ServiceName) error {
	// Get route IDs for the service
	serviceKey := fmt.Sprintf("%s%s", routeServicePrefix, serviceName.String())
	routeIDs, err := r.client.SMembers(ctx, serviceKey).Result()
	
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get routes for service: %w", err)
	}
	
	if len(routeIDs) == 0 {
		return nil
	}
	
	// Create pipeline for multiple operations
	pipe := r.client.Pipeline()
	
	// Delete service key
	pipe.Del(ctx, serviceKey)
	
	// For each route, remove it from cache
	for _, routeIDStr := range routeIDs {
		routeID, err := valueobject.NewRouteID(routeIDStr)
		if err != nil {
			continue
		}
		
		// Get the route to know its path
		route, err := r.FindRouteByID(ctx, routeID)
		if err != nil || route == nil {
			continue
		}
		
		// Remove route by ID
		routeKey := fmt.Sprintf("%s%s", routeKeyPrefix, routeIDStr)
		pipe.Del(ctx, routeKey)
		
		// Remove from path routes set
		pathKey := fmt.Sprintf("%s%s", routePathPrefix, route.Path().String())
		pipe.SRem(ctx, pathKey, routeIDStr)
		
		// Remove from active routes set
		pipe.SRem(ctx, routeActivesKey, routeIDStr)
	}
	
	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to invalidate service routes from cache: %w", err)
	}
	
	return nil
}

// getRoutesByIDs gets multiple routes by their IDs
func (r *RedisRouteCacheRepository) getRoutesByIDs(ctx context.Context, routeIDs []string) ([]*entity.Route, error) {
	if len(routeIDs) == 0 {
		return nil, nil
	}
	
	// Create pipeline for batch operations
	pipe := r.client.Pipeline()
	
	// Queue up GET commands for each route
	gets := make([]*redis.StringCmd, len(routeIDs))
	for i, routeID := range routeIDs {
		routeKey := fmt.Sprintf("%s%s", routeKeyPrefix, routeID)
		gets[i] = pipe.Get(ctx, routeKey)
	}
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get routes by IDs: %w", err)
	}
	
	// Process results
	var routes []*entity.Route
	for _, cmd := range gets {
		routeJSON, err := cmd.Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			// Skip routes that can't be retrieved
			continue
		}
		
		route, err := r.fromJSON(routeJSON)
		if err != nil {
			// Skip routes that can't be parsed
			continue
		}
		
		routes = append(routes, route)
	}
	
	return routes, nil
}

// toDTO converts a route entity to a DTO for serialization
func (r *RedisRouteCacheRepository) toDTO(route *entity.Route) RouteDTO {
	return RouteDTO{
		ID:                    route.ID().String(),
		Path:                  route.Path().String(),
		ServiceName:           route.ServiceName().String(),
		TargetURL:             route.TargetURL().String(),
		AuthType:              route.AuthType().String(),
		RateLimitType:         route.RateLimit().Type().String(),
		RateLimitRequests:     route.RateLimit().Requests(),
		RateLimitDurationSecs: int(route.RateLimit().Duration().Seconds()),
		IsActive:              route.IsActive(),
		CreatedAt:             route.CreatedAt(),
		UpdatedAt:             route.UpdatedAt(),
	}
}

// fromJSON deserializes a route from JSON
func (r *RedisRouteCacheRepository) fromJSON(routeJSON string) (*entity.Route, error) {
	var dto RouteDTO
	err := json.Unmarshal([]byte(routeJSON), &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal route: %w", err)
	}
	
	// Create value objects
	routeID, err := valueobject.NewRouteID(dto.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid route ID: %w", err)
	}
	
	path, err := valueobject.NewPath(dto.Path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	
	serviceName, err := valueobject.NewServiceName(dto.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	targetURL, err := valueobject.NewTargetURL(dto.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}
	
	authType := valueobject.AuthTypeFromString(dto.AuthType)
	
	rateLimitType := valueobject.RateLimitTypeFromString(dto.RateLimitType)
	rateLimit, err := valueobject.NewRateLimit(
		rateLimitType,
		dto.RateLimitRequests,
		dto.RateLimitDurationSecs,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limit: %w", err)
	}
	
	// Create route entity
	routeBuilder := entity.NewRouteBuilder()
	route := routeBuilder.
		WithID(routeID).
		WithPath(path).
		WithServiceName(serviceName).
		WithTargetURL(targetURL).
		WithAuthType(authType).
		WithRateLimit(rateLimit).
		WithActive(dto.IsActive).
		WithCreatedAt(dto.CreatedAt).
		WithUpdatedAt(dto.UpdatedAt).
		Build()
	
	return route, nil
}
