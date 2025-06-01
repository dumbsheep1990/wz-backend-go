package gateway

import (
	"context"
	"log"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
	mysqlRepo "wz-backend-go/internal/infrastructure/persistence/mysql/gateway"
	redisRepo "wz-backend-go/internal/infrastructure/persistence/redis/gateway"
)

// RouteRepositoryFacade implements repository.RouteRepository using a combination
// of MySQL persistence and Redis caching
type RouteRepositoryFacade struct {
	mysqlRepo     *mysqlRepo.MySQLRouteRepository
	redisCacheRepo *redisRepo.RedisRouteCacheRepository
}

// NewRouteRepositoryFacade creates a new RouteRepositoryFacade
func NewRouteRepositoryFacade(
	mysqlRepo *mysqlRepo.MySQLRouteRepository,
	redisCacheRepo *redisRepo.RedisRouteCacheRepository,
) *RouteRepositoryFacade {
	return &RouteRepositoryFacade{
		mysqlRepo:      mysqlRepo,
		redisCacheRepo: redisCacheRepo,
	}
}

// Save saves a route to both MySQL and Redis cache
func (f *RouteRepositoryFacade) Save(ctx context.Context, route *entity.Route) error {
	// First save to MySQL for persistence
	err := f.mysqlRepo.Save(ctx, route)
	if err != nil {
		return err
	}

	// Then update the cache
	err = f.redisCacheRepo.SaveRoute(ctx, route)
	if err != nil {
		// Log but don't fail if cache update fails
		log.Printf("Failed to update route cache: %v", err)
	}

	return nil
}

// FindByID finds a route by ID, trying cache first then falling back to MySQL
func (f *RouteRepositoryFacade) FindByID(ctx context.Context, id valueobject.RouteID) (*entity.Route, error) {
	// Try to get from cache first
	route, err := f.redisCacheRepo.FindRouteByID(ctx, id)
	if err != nil {
		log.Printf("Cache error when finding route by ID: %v", err)
	}

	// If found in cache, return it
	if route != nil {
		return route, nil
	}

	// Otherwise, get from MySQL
	route, err = f.mysqlRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache
	if route != nil {
		err = f.redisCacheRepo.SaveRoute(ctx, route)
		if err != nil {
			// Log but don't fail if cache update fails
			log.Printf("Failed to update route cache after MySQL lookup: %v", err)
		}
	}

	return route, nil
}

// FindByPath finds routes by path, trying cache first then falling back to MySQL
func (f *RouteRepositoryFacade) FindByPath(ctx context.Context, path valueobject.Path) ([]*entity.Route, error) {
	// Try to get from cache first
	routes, err := f.redisCacheRepo.FindRoutesByPath(ctx, path)
	if err != nil {
		log.Printf("Cache error when finding routes by path: %v", err)
	}

	// If found in cache, return them
	if routes != nil && len(routes) > 0 {
		return routes, nil
	}

	// Otherwise, get from MySQL
	routes, err = f.mysqlRepo.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache for each route
	if routes != nil && len(routes) > 0 {
		for _, route := range routes {
			err = f.redisCacheRepo.SaveRoute(ctx, route)
			if err != nil {
				// Log but don't fail if cache update fails
				log.Printf("Failed to update route cache after MySQL lookup: %v", err)
			}
		}
	}

	return routes, nil
}

// FindByServiceName finds routes by service name, trying cache first then falling back to MySQL
func (f *RouteRepositoryFacade) FindByServiceName(ctx context.Context, serviceName valueobject.ServiceName) ([]*entity.Route, error) {
	// Try to get from cache first
	routes, err := f.redisCacheRepo.FindRoutesByServiceName(ctx, serviceName)
	if err != nil {
		log.Printf("Cache error when finding routes by service name: %v", err)
	}

	// If found in cache, return them
	if routes != nil && len(routes) > 0 {
		return routes, nil
	}

	// Otherwise, get from MySQL
	routes, err = f.mysqlRepo.FindByServiceName(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache for each route
	if routes != nil && len(routes) > 0 {
		for _, route := range routes {
			err = f.redisCacheRepo.SaveRoute(ctx, route)
			if err != nil {
				// Log but don't fail if cache update fails
				log.Printf("Failed to update route cache after MySQL lookup: %v", err)
			}
		}
	}

	return routes, nil
}

// FindActive finds all active routes, trying cache first then falling back to MySQL
func (f *RouteRepositoryFacade) FindActive(ctx context.Context) ([]*entity.Route, error) {
	// Try to get from cache first
	routes, err := f.redisCacheRepo.FindActiveRoutes(ctx)
	if err != nil {
		log.Printf("Cache error when finding active routes: %v", err)
	}

	// If found in cache, return them
	if routes != nil && len(routes) > 0 {
		return routes, nil
	}

	// Otherwise, get from MySQL
	routes, err = f.mysqlRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	// If found in MySQL but not in cache, update the cache for each route
	if routes != nil && len(routes) > 0 {
		for _, route := range routes {
			err = f.redisCacheRepo.SaveRoute(ctx, route)
			if err != nil {
				// Log but don't fail if cache update fails
				log.Printf("Failed to update route cache after MySQL lookup: %v", err)
			}
		}
	}

	return routes, nil
}

// FindAll finds all routes with pagination, only using MySQL
// This is not cached because it's typically used for admin interfaces and full listing
func (f *RouteRepositoryFacade) FindAll(ctx context.Context, offset, limit int) ([]*entity.Route, int, error) {
	return f.mysqlRepo.FindAll(ctx, offset, limit)
}

// Delete deletes a route from both MySQL and Redis cache
func (f *RouteRepositoryFacade) Delete(ctx context.Context, id valueobject.RouteID) error {
	// First delete from MySQL
	err := f.mysqlRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Then delete from cache
	err = f.redisCacheRepo.DeleteRoute(ctx, id)
	if err != nil {
		// Log but don't fail if cache update fails
		log.Printf("Failed to delete route from cache: %v", err)
	}

	return nil
}

// EnsureTable ensures the MySQL table exists
func (f *RouteRepositoryFacade) EnsureTable(ctx context.Context) error {
	return f.mysqlRepo.EnsureTable(ctx)
}
