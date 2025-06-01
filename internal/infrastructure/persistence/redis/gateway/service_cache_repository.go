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
	serviceKeyPrefix    = "gateway:service:"
	serviceActivesKey   = "gateway:active-services"
	serviceHealthyKey   = "gateway:healthy-services"
	serviceCacheTTL     = 15 * time.Minute
	healthStatusTTL     = 1 * time.Minute // Shorter TTL for health status
)

// ServiceDTO is a data transfer object for serializing services to Redis
type ServiceDTO struct {
	Name                   string    `json:"name"`
	BaseURL                string    `json:"base_url"`
	HealthURL              string    `json:"health_url,omitempty"`
	DefaultAuth            string    `json:"default_auth"`
	RateLimitType          string    `json:"rate_limit_type"`
	RateLimitRequests      int       `json:"rate_limit_requests"`
	RateLimitDurationSecs  int       `json:"rate_limit_duration_secs"`
	IsActive               bool      `json:"is_active"`
	IsHealthy              bool      `json:"is_healthy"`
	LastHealthy            *int64    `json:"last_healthy,omitempty"` // Unix timestamp
	ErrorMessage           string    `json:"error_message,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// RedisServiceCacheRepository implements a caching layer for services using Redis
type RedisServiceCacheRepository struct {
	client redis.Client
}

// NewRedisServiceCacheRepository creates a new RedisServiceCacheRepository
func NewRedisServiceCacheRepository(client redis.Client) *RedisServiceCacheRepository {
	return &RedisServiceCacheRepository{
		client: client,
	}
}

// SaveService caches a service in Redis
func (r *RedisServiceCacheRepository) SaveService(ctx context.Context, service *entity.Service) error {
	// Convert service to DTO for serialization
	dto := r.toDTO(service)
	
	// Serialize DTO to JSON
	serviceJSON, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}
	
	// Create pipeline for multiple operations
	pipe := r.client.Pipeline()
	
	// Store service by name
	serviceKey := fmt.Sprintf("%s%s", serviceKeyPrefix, service.Name().String())
	pipe.Set(ctx, serviceKey, serviceJSON, serviceCacheTTL)
	
	// Update active services set
	if service.IsActive() {
		pipe.SAdd(ctx, serviceActivesKey, service.Name().String())
		pipe.Expire(ctx, serviceActivesKey, serviceCacheTTL)
	} else {
		pipe.SRem(ctx, serviceActivesKey, service.Name().String())
	}
	
	// Update healthy services set
	if service.IsHealthy() {
		pipe.SAdd(ctx, serviceHealthyKey, service.Name().String())
		pipe.Expire(ctx, serviceHealthyKey, healthStatusTTL)
	} else {
		pipe.SRem(ctx, serviceHealthyKey, service.Name().String())
	}
	
	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save service to cache: %w", err)
	}
	
	return nil
}

// FindServiceByName finds a service by name from cache
func (r *RedisServiceCacheRepository) FindServiceByName(ctx context.Context, name valueobject.ServiceName) (*entity.Service, error) {
	serviceKey := fmt.Sprintf("%s%s", serviceKeyPrefix, name.String())
	serviceJSON, err := r.client.Get(ctx, serviceKey).Result()
	
	if err == redis.Nil {
		return nil, nil // Not found in cache
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get service from cache: %w", err)
	}
	
	return r.fromJSON(serviceJSON)
}

// FindActiveServices finds all active services from cache
func (r *RedisServiceCacheRepository) FindActiveServices(ctx context.Context) ([]*entity.Service, error) {
	serviceNames, err := r.client.SMembers(ctx, serviceActivesKey).Result()
	
	if err == redis.Nil || len(serviceNames) == 0 {
		return nil, nil // No active services found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active services from cache: %w", err)
	}
	
	return r.getServicesByNames(ctx, serviceNames)
}

// FindHealthyServices finds all healthy services from cache
func (r *RedisServiceCacheRepository) FindHealthyServices(ctx context.Context) ([]*entity.Service, error) {
	serviceNames, err := r.client.SMembers(ctx, serviceHealthyKey).Result()
	
	if err == redis.Nil || len(serviceNames) == 0 {
		return nil, nil // No healthy services found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get healthy services from cache: %w", err)
	}
	
	return r.getServicesByNames(ctx, serviceNames)
}

// DeleteService removes a service from cache
func (r *RedisServiceCacheRepository) DeleteService(ctx context.Context, name valueobject.ServiceName) error {
	// Create pipeline for multiple operations
	pipe := r.client.Pipeline()
	
	// Remove service by name
	serviceKey := fmt.Sprintf("%s%s", serviceKeyPrefix, name.String())
	pipe.Del(ctx, serviceKey)
	
	// Remove from active services set
	pipe.SRem(ctx, serviceActivesKey, name.String())
	
	// Remove from healthy services set
	pipe.SRem(ctx, serviceHealthyKey, name.String())
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete service from cache: %w", err)
	}
	
	return nil
}

// UpdateServiceHealth updates only the health status of a service in cache
func (r *RedisServiceCacheRepository) UpdateServiceHealth(
	ctx context.Context, 
	name valueobject.ServiceName, 
	isHealthy bool, 
	lastHealthy *time.Time, 
	errorMessage string,
) error {
	// First get the existing service
	service, err := r.FindServiceByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get service for health update: %w", err)
	}
	
	if service == nil {
		return fmt.Errorf("service not found in cache for health update")
	}
	
	// Update the service health properties
	updatedService := entity.NewServiceBuilder().
		WithName(service.Name()).
		WithBaseURL(service.BaseURL()).
		WithHealthURL(service.HealthURL()).
		WithDefaultAuth(service.DefaultAuth()).
		WithRateLimit(service.RateLimit()).
		WithActive(service.IsActive()).
		WithHealthy(isHealthy).
		WithErrorMessage(errorMessage).
		WithCreatedAt(service.CreatedAt()).
		WithUpdatedAt(time.Now()).
		Build()
	
	if lastHealthy != nil {
		updatedService = entity.NewServiceBuilder().
			FromService(updatedService).
			WithLastHealthy(lastHealthy).
			Build()
	}
	
	// Save the updated service
	return r.SaveService(ctx, updatedService)
}

// getServicesByNames gets multiple services by their names
func (r *RedisServiceCacheRepository) getServicesByNames(ctx context.Context, serviceNames []string) ([]*entity.Service, error) {
	if len(serviceNames) == 0 {
		return nil, nil
	}
	
	// Create pipeline for batch operations
	pipe := r.client.Pipeline()
	
	// Queue up GET commands for each service
	gets := make([]*redis.StringCmd, len(serviceNames))
	for i, serviceName := range serviceNames {
		serviceKey := fmt.Sprintf("%s%s", serviceKeyPrefix, serviceName)
		gets[i] = pipe.Get(ctx, serviceKey)
	}
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get services by names: %w", err)
	}
	
	// Process results
	var services []*entity.Service
	for _, cmd := range gets {
		serviceJSON, err := cmd.Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			// Skip services that can't be retrieved
			continue
		}
		
		service, err := r.fromJSON(serviceJSON)
		if err != nil {
			// Skip services that can't be parsed
			continue
		}
		
		services = append(services, service)
	}
	
	return services, nil
}

// toDTO converts a service entity to a DTO for serialization
func (r *RedisServiceCacheRepository) toDTO(service *entity.Service) ServiceDTO {
	dto := ServiceDTO{
		Name:                  service.Name().String(),
		BaseURL:               service.BaseURL().String(),
		DefaultAuth:           service.DefaultAuth().String(),
		RateLimitType:         service.RateLimit().Type().String(),
		RateLimitRequests:     service.RateLimit().Requests(),
		RateLimitDurationSecs: int(service.RateLimit().Duration().Seconds()),
		IsActive:              service.IsActive(),
		IsHealthy:             service.IsHealthy(),
		ErrorMessage:          service.ErrorMessage(),
		CreatedAt:             service.CreatedAt(),
		UpdatedAt:             service.UpdatedAt(),
	}
	
	// Add health URL if not empty
	if url := service.HealthURL().URL(); url != nil {
		dto.HealthURL = service.HealthURL().String()
	}
	
	// Add last healthy time if available
	if lastHealthy := service.LastHealthy(); lastHealthy != nil {
		unixTime := lastHealthy.Unix()
		dto.LastHealthy = &unixTime
	}
	
	return dto
}

// fromJSON deserializes a service from JSON
func (r *RedisServiceCacheRepository) fromJSON(serviceJSON string) (*entity.Service, error) {
	var dto ServiceDTO
	err := json.Unmarshal([]byte(serviceJSON), &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service: %w", err)
	}
	
	// Create value objects
	serviceName, err := valueobject.NewServiceName(dto.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	
	baseURL, err := valueobject.NewTargetURL(dto.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	
	var healthURL valueobject.TargetURL
	if dto.HealthURL != "" {
		healthURL, err = valueobject.NewTargetURL(dto.HealthURL)
		if err != nil {
			return nil, fmt.Errorf("invalid health URL: %w", err)
		}
	} else {
		healthURL, _ = valueobject.NewEmptyTargetURL()
	}
	
	defaultAuth := valueobject.AuthTypeFromString(dto.DefaultAuth)
	
	rateLimitType := valueobject.RateLimitTypeFromString(dto.RateLimitType)
	rateLimit, err := valueobject.NewRateLimit(
		rateLimitType,
		dto.RateLimitRequests,
		dto.RateLimitDurationSecs,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limit: %w", err)
	}
	
	// Create service entity
	serviceBuilder := entity.NewServiceBuilder()
	serviceBuilder = serviceBuilder.
		WithName(serviceName).
		WithBaseURL(baseURL).
		WithHealthURL(healthURL).
		WithDefaultAuth(defaultAuth).
		WithRateLimit(rateLimit).
		WithActive(dto.IsActive).
		WithHealthy(dto.IsHealthy).
		WithErrorMessage(dto.ErrorMessage).
		WithCreatedAt(dto.CreatedAt).
		WithUpdatedAt(dto.UpdatedAt)
	
	// Add last healthy time if available
	if dto.LastHealthy != nil {
		lastHealthyTime := time.Unix(*dto.LastHealthy, 0)
		serviceBuilder = serviceBuilder.WithLastHealthy(&lastHealthyTime)
	}
	
	service := serviceBuilder.Build()
	return service, nil
}
