package repository

import (
	"context"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// ServiceRepository defines the interface for service persistence
type ServiceRepository interface {
	// FindByName finds a service by name
	FindByName(ctx context.Context, name valueobject.ServiceName) (*entity.Service, error)
	
	// FindAll finds all services with pagination
	FindAll(ctx context.Context, offset, limit int) ([]*entity.Service, int, error)
	
	// FindActive finds all active services
	FindActive(ctx context.Context) ([]*entity.Service, error)
	
	// FindHealthy finds all healthy services
	FindHealthy(ctx context.Context) ([]*entity.Service, error)
	
	// FindUnhealthy finds all unhealthy services
	FindUnhealthy(ctx context.Context) ([]*entity.Service, error)
	
	// Save persists a service
	Save(ctx context.Context, service *entity.Service) error
	
	// Delete removes a service
	Delete(ctx context.Context, name valueobject.ServiceName) error
}
