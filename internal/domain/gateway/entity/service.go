package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/gateway/valueobject"
)

// Service represents a backend service registered in the gateway
type Service struct {
	name         valueobject.ServiceName
	baseURL      valueobject.TargetURL
	healthURL    valueobject.TargetURL
	defaultAuth  valueobject.AuthType
	rateLimit    valueobject.RateLimit
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
	lastHealthy  *time.Time
	isHealthy    bool
	errorMessage string
}

// ServiceBuilder is a builder for Service
type ServiceBuilder struct {
	name         valueobject.ServiceName
	baseURL      valueobject.TargetURL
	healthURL    valueobject.TargetURL
	defaultAuth  valueobject.AuthType
	rateLimit    valueobject.RateLimit
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
	lastHealthy  *time.Time
	isHealthy    bool
	errorMessage string
	err          error
}

// NewServiceBuilder creates a new ServiceBuilder
func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{
		defaultAuth: valueobject.NoAuth,
		isActive:    true,
		isHealthy:   true,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

// WithName sets the service name
func (b *ServiceBuilder) WithName(name string) *ServiceBuilder {
	if b.err != nil {
		return b
	}
	b.name, b.err = valueobject.NewServiceName(name)
	return b
}

// WithBaseURL sets the base URL
func (b *ServiceBuilder) WithBaseURL(url string) *ServiceBuilder {
	if b.err != nil {
		return b
	}
	b.baseURL, b.err = valueobject.NewTargetURL(url)
	return b
}

// WithHealthURL sets the health check URL
func (b *ServiceBuilder) WithHealthURL(url string) *ServiceBuilder {
	if b.err != nil {
		return b
	}
	b.healthURL, b.err = valueobject.NewTargetURL(url)
	return b
}

// WithDefaultAuth sets the default authentication type
func (b *ServiceBuilder) WithDefaultAuth(authType valueobject.AuthType) *ServiceBuilder {
	if b.err != nil {
		return b
	}
	if !authType.IsValid() {
		b.err = errors.New("invalid auth type")
		return b
	}
	b.defaultAuth = authType
	return b
}

// WithRateLimit sets the rate limit
func (b *ServiceBuilder) WithRateLimit(limitType valueobject.RateLimitType, requests, durationSeconds int) *ServiceBuilder {
	if b.err != nil {
		return b
	}
	b.rateLimit, b.err = valueobject.NewRateLimit(limitType, requests, durationSeconds)
	return b
}

// WithIsActive sets whether the service is active
func (b *ServiceBuilder) WithIsActive(isActive bool) *ServiceBuilder {
	b.isActive = isActive
	return b
}

// Build creates a new Service
func (b *ServiceBuilder) Build() (*Service, error) {
	if b.err != nil {
		return nil, b.err
	}
	
	// Validate required fields
	if b.name.String() == "" {
		return nil, errors.New("service name is required")
	}
	if b.baseURL.String() == "" {
		return nil, errors.New("base URL is required")
	}
	
	return &Service{
		name:         b.name,
		baseURL:      b.baseURL,
		healthURL:    b.healthURL,
		defaultAuth:  b.defaultAuth,
		rateLimit:    b.rateLimit,
		isActive:     b.isActive,
		createdAt:    b.createdAt,
		updatedAt:    b.updatedAt,
		isHealthy:    b.isHealthy,
		lastHealthy:  b.lastHealthy,
		errorMessage: b.errorMessage,
	}, nil
}

// Name returns the service name
func (s *Service) Name() valueobject.ServiceName {
	return s.name
}

// BaseURL returns the base URL
func (s *Service) BaseURL() valueobject.TargetURL {
	return s.baseURL
}

// HealthURL returns the health check URL
func (s *Service) HealthURL() valueobject.TargetURL {
	return s.healthURL
}

// DefaultAuth returns the default authentication type
func (s *Service) DefaultAuth() valueobject.AuthType {
	return s.defaultAuth
}

// RateLimit returns the rate limit
func (s *Service) RateLimit() valueobject.RateLimit {
	return s.rateLimit
}

// IsActive returns whether the service is active
func (s *Service) IsActive() bool {
	return s.isActive
}

// CreatedAt returns the creation time
func (s *Service) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns the last update time
func (s *Service) UpdatedAt() time.Time {
	return s.updatedAt
}

// IsHealthy returns whether the service is healthy
func (s *Service) IsHealthy() bool {
	return s.isHealthy
}

// LastHealthy returns the last time the service was healthy
func (s *Service) LastHealthy() *time.Time {
	return s.lastHealthy
}

// ErrorMessage returns the error message from the last health check
func (s *Service) ErrorMessage() string {
	return s.errorMessage
}

// Activate activates the service
func (s *Service) Activate() {
	s.isActive = true
	s.updatedAt = time.Now()
}

// Deactivate deactivates the service
func (s *Service) Deactivate() {
	s.isActive = false
	s.updatedAt = time.Now()
}

// UpdateBaseURL updates the base URL
func (s *Service) UpdateBaseURL(baseURL valueobject.TargetURL) {
	s.baseURL = baseURL
	s.updatedAt = time.Now()
}

// UpdateHealthURL updates the health check URL
func (s *Service) UpdateHealthURL(healthURL valueobject.TargetURL) {
	s.healthURL = healthURL
	s.updatedAt = time.Now()
}

// UpdateDefaultAuth updates the default authentication type
func (s *Service) UpdateDefaultAuth(defaultAuth valueobject.AuthType) {
	s.defaultAuth = defaultAuth
	s.updatedAt = time.Now()
}

// UpdateRateLimit updates the rate limit
func (s *Service) UpdateRateLimit(rateLimit valueobject.RateLimit) {
	s.rateLimit = rateLimit
	s.updatedAt = time.Now()
}

// SetHealthy marks the service as healthy
func (s *Service) SetHealthy() {
	now := time.Now()
	s.isHealthy = true
	s.lastHealthy = &now
	s.errorMessage = ""
	s.updatedAt = now
}

// SetUnhealthy marks the service as unhealthy
func (s *Service) SetUnhealthy(errorMessage string) {
	s.isHealthy = false
	s.errorMessage = errorMessage
	s.updatedAt = time.Now()
}
