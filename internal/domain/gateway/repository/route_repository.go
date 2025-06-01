package repository

import (
	"context"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// RouteRepository defines the interface for route persistence
type RouteRepository interface {
	// FindByID finds a route by ID
	FindByID(ctx context.Context, id valueobject.RouteID) (*entity.Route, error)
	
	// FindByPath finds routes that match a path
	FindByPath(ctx context.Context, path string) ([]*entity.Route, error)
	
	// FindByServiceName finds routes for a service
	FindByServiceName(ctx context.Context, serviceName valueobject.ServiceName) ([]*entity.Route, error)
	
	// FindAll finds all routes with pagination
	FindAll(ctx context.Context, offset, limit int) ([]*entity.Route, int, error)
	
	// FindActive finds all active routes
	FindActive(ctx context.Context) ([]*entity.Route, error)
	
	// Save persists a route
	Save(ctx context.Context, route *entity.Route) error
	
	// Delete removes a route
	Delete(ctx context.Context, id valueobject.RouteID) error
}
