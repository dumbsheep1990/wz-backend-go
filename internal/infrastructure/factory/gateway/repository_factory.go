package gateway

import (
	"database/sql"

	"github.com/go-redis/redis/v8"

	"wz-backend-go/internal/domain/gateway/repository"
	mysqlRepo "wz-backend-go/internal/infrastructure/persistence/mysql/gateway"
	redisRepo "wz-backend-go/internal/infrastructure/persistence/redis/gateway"
)

// RepositoryFactory creates and provides access to different Gateway Service repositories
type RepositoryFactory struct {
	db          *sql.DB
	redisClient *redis.Client
	
	// MySQL repositories
	mysqlRouteRepo    *mysqlRepo.MySQLRouteRepository
	mysqlServiceRepo  *mysqlRepo.MySQLServiceRepository
	
	// Redis repositories
	redisRouteCacheRepo    *redisRepo.RedisRouteCacheRepository
	redisServiceCacheRepo  *redisRepo.RedisServiceCacheRepository
	redisRateLimiterRepo   *redisRepo.RedisRateLimiterRepository
	
	// Repository facades
	routeRepoFacade    *RouteRepositoryFacade
	serviceRepoFacade  *ServiceRepositoryFacade
}

// NewRepositoryFactory creates a new RepositoryFactory
func NewRepositoryFactory(db *sql.DB, redisClient *redis.Client) *RepositoryFactory {
	return &RepositoryFactory{
		db:          db,
		redisClient: redisClient,
	}
}

// GetRouteRepository returns a route repository that combines MySQL and Redis caching
func (f *RepositoryFactory) GetRouteRepository() repository.RouteRepository {
	if f.routeRepoFacade == nil {
		f.initRouteRepositories()
		f.routeRepoFacade = NewRouteRepositoryFacade(
			f.mysqlRouteRepo,
			f.redisRouteCacheRepo,
		)
	}
	return f.routeRepoFacade
}

// GetServiceRepository returns a service repository that combines MySQL and Redis caching
func (f *RepositoryFactory) GetServiceRepository() repository.ServiceRepository {
	if f.serviceRepoFacade == nil {
		f.initServiceRepositories()
		f.serviceRepoFacade = NewServiceRepositoryFacade(
			f.mysqlServiceRepo,
			f.redisServiceCacheRepo,
		)
	}
	return f.serviceRepoFacade
}

// GetRateLimiterRepository returns a rate limiter repository implementation
func (f *RepositoryFactory) GetRateLimiterRepository() *redisRepo.RedisRateLimiterRepository {
	if f.redisRateLimiterRepo == nil {
		f.redisRateLimiterRepo = redisRepo.NewRedisRateLimiterRepository(*f.redisClient)
	}
	return f.redisRateLimiterRepo
}

// GetMySQLRouteRepository returns the MySQL route repository directly
func (f *RepositoryFactory) GetMySQLRouteRepository() *mysqlRepo.MySQLRouteRepository {
	if f.mysqlRouteRepo == nil {
		f.mysqlRouteRepo = mysqlRepo.NewMySQLRouteRepository(f.db)
	}
	return f.mysqlRouteRepo
}

// GetMySQLServiceRepository returns the MySQL service repository directly
func (f *RepositoryFactory) GetMySQLServiceRepository() *mysqlRepo.MySQLServiceRepository {
	if f.mysqlServiceRepo == nil {
		f.mysqlServiceRepo = mysqlRepo.NewMySQLServiceRepository(f.db)
	}
	return f.mysqlServiceRepo
}

// GetRedisRouteCacheRepository returns the Redis route cache repository directly
func (f *RepositoryFactory) GetRedisRouteCacheRepository() *redisRepo.RedisRouteCacheRepository {
	if f.redisRouteCacheRepo == nil {
		f.redisRouteCacheRepo = redisRepo.NewRedisRouteCacheRepository(*f.redisClient)
	}
	return f.redisRouteCacheRepo
}

// GetRedisServiceCacheRepository returns the Redis service cache repository directly
func (f *RepositoryFactory) GetRedisServiceCacheRepository() *redisRepo.RedisServiceCacheRepository {
	if f.redisServiceCacheRepo == nil {
		f.redisServiceCacheRepo = redisRepo.NewRedisServiceCacheRepository(*f.redisClient)
	}
	return f.redisServiceCacheRepo
}

// initRouteRepositories initializes the route repositories if they don't exist
func (f *RepositoryFactory) initRouteRepositories() {
	if f.mysqlRouteRepo == nil {
		f.mysqlRouteRepo = mysqlRepo.NewMySQLRouteRepository(f.db)
	}
	if f.redisRouteCacheRepo == nil {
		f.redisRouteCacheRepo = redisRepo.NewRedisRouteCacheRepository(*f.redisClient)
	}
}

// initServiceRepositories initializes the service repositories if they don't exist
func (f *RepositoryFactory) initServiceRepositories() {
	if f.mysqlServiceRepo == nil {
		f.mysqlServiceRepo = mysqlRepo.NewMySQLServiceRepository(f.db)
	}
	if f.redisServiceCacheRepo == nil {
		f.redisServiceCacheRepo = redisRepo.NewRedisServiceCacheRepository(*f.redisClient)
	}
}
