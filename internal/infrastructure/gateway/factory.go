package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/service"
	gatewayRepo "wz-backend-go/internal/infrastructure/gateway/repository"
	"wz-backend-go/internal/infrastructure/gateway/middleware"
)

// Factory 网关基础设施工厂，负责创建所有网关基础设施层组件
type Factory struct {
	db          *sqlx.DB
	redisClient *redis.Client
	jwtSecret   string
	jwtIssuer   string
	version     string
}

// NewFactory 创建新的网关基础设施工厂
func NewFactory(db *sqlx.DB, redisClient *redis.Client, jwtSecret string, jwtIssuer string, version string) *Factory {
	return &Factory{
		db:          db,
		redisClient: redisClient,
		jwtSecret:   jwtSecret,
		jwtIssuer:   jwtIssuer,
		version:     version,
	}
}

// CreateServiceRepository 创建服务仓储
func (f *Factory) CreateServiceRepository() repository.ServiceRepository {
	return gatewayRepo.NewServiceRepository(f.db)
}

// CreateRouteRepository 创建路由仓储
func (f *Factory) CreateRouteRepository() repository.RouteRepository {
	return gatewayRepo.NewRouteRepository(f.db)
}

// CreateRateLimiterRepository 创建限流器仓储
func (f *Factory) CreateRateLimiterRepository() repository.RateLimiterRepository {
	return gatewayRepo.NewRateLimiterRepository(f.db)
}

// CreateRouter 创建路由管理器
func (f *Factory) CreateRouter(routeRepo repository.RouteRepository) *RouterImpl {
	return NewRouter(routeRepo)
}

// CreateServiceRegistry 创建服务注册表
func (f *Factory) CreateServiceRegistry(serviceRepo repository.ServiceRepository) *ServiceRegistryImpl {
	return NewServiceRegistry(serviceRepo)
}

// CreateRateLimiter 创建限流器
func (f *Factory) CreateRateLimiter() *RateLimiterImpl {
	return NewRateLimiter(f.redisClient)
}

// CreateAuthHandler 创建认证处理器
func (f *Factory) CreateAuthHandler(tokenTTL int) *AuthHandlerImpl {
	// 将tokenTTL转换为持续时间（秒）
	tokenTTLDuration := tokenTTL * 1000000000 // 转换为纳秒
	return NewAuthHandler(f.redisClient, f.jwtSecret, f.jwtIssuer, tokenTTLDuration)
}

// CreateMiddlewareManager 创建中间件管理器
func (f *Factory) CreateMiddlewareManager(domainService service.GatewayDomainService) *middleware.MiddlewareManager {
	return middleware.NewMiddlewareManager(domainService, f.redisClient, f.version)
}

// ApplyGlobalMiddlewares 应用全局中间件
func (f *Factory) ApplyGlobalMiddlewares(engine *gin.Engine, middlewareManager *middleware.MiddlewareManager) {
	middlewareManager.ApplyGlobalMiddlewares(engine)
}

// ApplyRouteMiddlewares 应用路由中间件
func (f *Factory) ApplyRouteMiddlewares(routerGroup *gin.RouterGroup, middlewareManager *middleware.MiddlewareManager) {
	middlewareManager.ApplyRouteMiddlewares(routerGroup)
}
