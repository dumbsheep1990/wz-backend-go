package gateway

import (
	"context"
	"log"
	"time"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
	mysqlRepo "wz-backend-go/internal/infrastructure/persistence/mysql/gateway"
	redisRepo "wz-backend-go/internal/infrastructure/persistence/redis/gateway"
)

// ServiceRepositoryFacade implements repository.ServiceRepository using a combination
// of MySQL persistence and Redis caching
type ServiceRepositoryFacade struct {
	mysqlRepo      *mysqlRepo.MySQLServiceRepository
	redisCacheRepo *redisRepo.RedisServiceCacheRepository
}

// NewServiceRepositoryFacade creates a new ServiceRepositoryFacade
func NewServiceRepositoryFacade(
	mysqlRepo *mysqlRepo.MySQLServiceRepository,
	redisCacheRepo *redisRepo.RedisServiceCacheRepository,
) *ServiceRepositoryFacade {
	return &ServiceRepositoryFacade{
		mysqlRepo:      mysqlRepo,
		redisCacheRepo: redisCacheRepo,
	}
}

// Save saves a service to both MySQL and Redis cache
func (f *ServiceRepositoryFacade) Save(ctx context.Context, service *entity.Service) error {
	// First save to MySQL for persistence
	err := f.mysqlRepo.Save(ctx, service)
	if err != nil {
		return err
	}

	// Then update the cache
	err = f.redisCacheRepo.SaveService(ctx, service)
	if err != nil {
		// Log but don't fail if cache update fails
		log.Printf("Failed to update service cache: %v", err)
	}

	return nil
}

// FindByName finds a service by name, trying cache first then falling back to MySQL
func (f *ServiceRepositoryFacade) FindByName(ctx context.Context, name valueobject.ServiceName) (*entity.Service, error) {
	// Try to get from cache first
	service, err := f.redisCacheRepo.FindServiceByName(ctx, name)
	if err != nil {
		log.Printf("Cache error when finding service by name: %v", err)
	}

	// If found in cache, return it
	if service != nil {
		return service, nil
	}

	// Otherwise, get from MySQL
	service, err = f.mysqlRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache
	if service != nil {
		err = f.redisCacheRepo.SaveService(ctx, service)
		if err != nil {
			// Log but don't fail if cache update fails
			log.Printf("Failed to update service cache after MySQL lookup: %v", err)
		}
	}

	return service, nil
}

// FindActive finds all active services, trying cache first then falling back to MySQL
func (f *ServiceRepositoryFacade) FindActive(ctx context.Context) ([]*entity.Service, error) {
	// Try to get from cache first
	services, err := f.redisCacheRepo.FindActiveServices(ctx)
	if err != nil {
		log.Printf("Cache error when finding active services: %v", err)
	}

	// If found in cache, return them
	if services != nil && len(services) > 0 {
		return services, nil
	}

	// Otherwise, get from MySQL
	services, err = f.mysqlRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache for each service
	if services != nil && len(services) > 0 {
		for _, service := range services {
			err = f.redisCacheRepo.SaveService(ctx, service)
			if err != nil {
				// Log but don't fail if cache update fails
				log.Printf("Failed to update service cache after MySQL lookup: %v", err)
			}
		}
	}

	return services, nil
}

// FindHealthy finds all healthy services, trying cache first then falling back to MySQL
func (f *ServiceRepositoryFacade) FindHealthy(ctx context.Context) ([]*entity.Service, error) {
	// Try to get from cache first
	services, err := f.redisCacheRepo.FindHealthyServices(ctx)
	if err != nil {
		log.Printf("Cache error when finding healthy services: %v", err)
	}

	// If found in cache, return them
	if services != nil && len(services) > 0 {
		return services, nil
	}

	// Otherwise, get from MySQL
	services, err = f.mysqlRepo.FindHealthy(ctx)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache for each service
	if services != nil && len(services) > 0 {
		for _, service := range services {
			err = f.redisCacheRepo.SaveService(ctx, service)
			if err != nil {
				// Log but don't fail if cache update fails
				log.Printf("Failed to update service cache after MySQL lookup: %v", err)
			}
		}
	}

	return services, nil
}

// FindAll finds all services with pagination, only using MySQL
// This is not cached because it's typically used for admin interfaces and full listing
func (f *ServiceRepositoryFacade) FindAll(ctx context.Context, offset, limit int) ([]*entity.Service, int, error) {
	return f.mysqlRepo.FindAll(ctx, offset, limit)
}

// Delete deletes a service from both MySQL and Redis cache
func (f *ServiceRepositoryFacade) Delete(ctx context.Context, name valueobject.ServiceName) error {
	// First delete from MySQL
	err := f.mysqlRepo.Delete(ctx, name)
	if err != nil {
		return err
	}

	// Then delete from cache
	err = f.redisCacheRepo.DeleteService(ctx, name)
	if err != nil {
		// Log but don't fail if cache update fails
		log.Printf("Failed to delete service from cache: %v", err)
	}

	return nil
}

// UpdateHealth updates the health status of a service in both MySQL and Redis cache
func (f *ServiceRepositoryFacade) UpdateHealth(ctx context.Context, name valueobject.ServiceName, isHealthy bool, lastHealthy *time.Time, errorMessage string) error {
	// First get the existing service
	service, err := f.FindByName(ctx, name)
	if err != nil {
		return err
	}
	
	if service == nil {
		return repository.ErrServiceNotFound
	}
	
	// Update the service health properties
	updatedService := entity.NewServiceBuilder().
		FromService(service).
		WithHealthy(isHealthy).
		WithErrorMessage(errorMessage).
		WithUpdatedAt(time.Now()).
		Build()
	
	if lastHealthy != nil {
		updatedService = entity.NewServiceBuilder().
			FromService(updatedService).
			WithLastHealthy(lastHealthy).
			Build()
	}
	
	// Save the updated service
	return f.Save(ctx, updatedService)
}

// EnsureTable ensures the MySQL table exists
func (f *ServiceRepositoryFacade) EnsureTable(ctx context.Context) error {
	return f.mysqlRepo.EnsureTable(ctx)
}
