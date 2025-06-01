package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/gateway/valueobject"
)

// Route represents a routing rule in the gateway
type Route struct {
	id          valueobject.RouteID
	path        valueobject.Path
	serviceName valueobject.ServiceName
	targetURL   valueobject.TargetURL
	authType    valueobject.AuthType
	rateLimit   valueobject.RateLimit
	createdAt   time.Time
	updatedAt   time.Time
	isActive    bool
}

// RouteBuilder is a builder for Route
type RouteBuilder struct {
	id          valueobject.RouteID
	path        valueobject.Path
	serviceName valueobject.ServiceName
	targetURL   valueobject.TargetURL
	authType    valueobject.AuthType
	rateLimit   valueobject.RateLimit
	createdAt   time.Time
	updatedAt   time.Time
	isActive    bool
	err         error
}

// NewRouteBuilder creates a new RouteBuilder
func NewRouteBuilder() *RouteBuilder {
	return &RouteBuilder{
		authType:  valueobject.NoAuth,
		isActive:  true,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

// WithID sets the route ID
func (b *RouteBuilder) WithID(id string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.id, b.err = valueobject.NewRouteID(id)
	return b
}

// WithPath sets the route path
func (b *RouteBuilder) WithPath(path string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.path, b.err = valueobject.NewPath(path)
	return b
}

// WithServiceName sets the service name
func (b *RouteBuilder) WithServiceName(name string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.serviceName, b.err = valueobject.NewServiceName(name)
	return b
}

// WithTargetURL sets the target URL
func (b *RouteBuilder) WithTargetURL(url string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.targetURL, b.err = valueobject.NewTargetURL(url)
	return b
}

// WithAuthType sets the authentication type
func (b *RouteBuilder) WithAuthType(authType valueobject.AuthType) *RouteBuilder {
	if b.err != nil {
		return b
	}
	if !authType.IsValid() {
		b.err = errors.New("invalid auth type")
		return b
	}
	b.authType = authType
	return b
}

// WithRateLimit sets the rate limit
func (b *RouteBuilder) WithRateLimit(limitType valueobject.RateLimitType, requests, durationSeconds int) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.rateLimit, b.err = valueobject.NewRateLimit(limitType, requests, durationSeconds)
	return b
}

// WithIsActive sets whether the route is active
func (b *RouteBuilder) WithIsActive(isActive bool) *RouteBuilder {
	b.isActive = isActive
	return b
}

// Build creates a new Route
func (b *RouteBuilder) Build() (*Route, error) {
	if b.err != nil {
		return nil, b.err
	}
	
	// Validate required fields
	if b.id.String() == "" {
		return nil, errors.New("route ID is required")
	}
	if b.path.String() == "" {
		return nil, errors.New("path is required")
	}
	if b.serviceName.String() == "" {
		return nil, errors.New("service name is required")
	}
	if b.targetURL.String() == "" {
		return nil, errors.New("target URL is required")
	}
	
	return &Route{
		id:          b.id,
		path:        b.path,
		serviceName: b.serviceName,
		targetURL:   b.targetURL,
		authType:    b.authType,
		rateLimit:   b.rateLimit,
		createdAt:   b.createdAt,
		updatedAt:   b.updatedAt,
		isActive:    b.isActive,
	}, nil
}

// ID returns the route ID
func (r *Route) ID() valueobject.RouteID {
	return r.id
}

// Path returns the route path
func (r *Route) Path() valueobject.Path {
	return r.path
}

// ServiceName returns the service name
func (r *Route) ServiceName() valueobject.ServiceName {
	return r.serviceName
}

// TargetURL returns the target URL
func (r *Route) TargetURL() valueobject.TargetURL {
	return r.targetURL
}

// AuthType returns the authentication type
func (r *Route) AuthType() valueobject.AuthType {
	return r.authType
}

// RateLimit returns the rate limit
func (r *Route) RateLimit() valueobject.RateLimit {
	return r.rateLimit
}

// CreatedAt returns the creation time
func (r *Route) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the last update time
func (r *Route) UpdatedAt() time.Time {
	return r.updatedAt
}

// IsActive returns whether the route is active
func (r *Route) IsActive() bool {
	return r.isActive
}

// Activate activates the route
func (r *Route) Activate() {
	r.isActive = true
	r.updatedAt = time.Now()
}

// Deactivate deactivates the route
func (r *Route) Deactivate() {
	r.isActive = false
	r.updatedAt = time.Now()
}

// UpdatePath updates the route path
func (r *Route) UpdatePath(path valueobject.Path) {
	r.path = path
	r.updatedAt = time.Now()
}

// UpdateTargetURL updates the target URL
func (r *Route) UpdateTargetURL(targetURL valueobject.TargetURL) {
	r.targetURL = targetURL
	r.updatedAt = time.Now()
}

// UpdateAuthType updates the authentication type
func (r *Route) UpdateAuthType(authType valueobject.AuthType) {
	r.authType = authType
	r.updatedAt = time.Now()
}

// UpdateRateLimit updates the rate limit
func (r *Route) UpdateRateLimit(rateLimit valueobject.RateLimit) {
	r.rateLimit = rateLimit
	r.updatedAt = time.Now()
}

// Matches checks if this route matches the given request path
func (r *Route) Matches(requestPath string) bool {
	return r.path.Matches(requestPath)
}
