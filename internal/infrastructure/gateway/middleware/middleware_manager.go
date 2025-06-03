package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"wz-backend-go/internal/domain/gateway/service"
)

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	domainService service.GatewayDomainService
	redisClient   *redis.Client
	version       string
	startTime     time.Time
}

// NewMiddlewareManager 创建新的中间件管理器
func NewMiddlewareManager(
	domainService service.GatewayDomainService,
	redisClient *redis.Client,
	version string,
) *MiddlewareManager {
	return &MiddlewareManager{
		domainService: domainService,
		redisClient:   redisClient,
		version:       version,
		startTime:     time.Now(),
	}
}

// ApplyGlobalMiddlewares 应用全局中间件
func (m *MiddlewareManager) ApplyGlobalMiddlewares(engine *gin.Engine) {
	// 设置全局中间件，按照从外到内的顺序应用
	
	// 1. 错误处理 - 最外层，捕获所有异常
	engine.Use(ErrorMiddleware())
	
	// 2. 日志记录 - 记录所有请求和响应
	engine.Use(LoggingMiddleware())
	
	// 3. CORS - 处理跨域请求
	engine.Use(CORSMiddleware())
	
	// 4. 健康检查 - 监控服务状态
	engine.Use(HealthCheckMiddleware(m.version, m.startTime))
}

// ApplyRouteMiddlewares 应用路由中间件
func (m *MiddlewareManager) ApplyRouteMiddlewares(routerGroup *gin.RouterGroup) {
	// 设置路由中间件，按照从外到内的顺序应用
	
	// 1. 路由解析 - 解析和匹配请求路由
	routerGroup.Use(NewRoutingMiddleware(m.domainService).Handle())
	
	// 2. 认证 - 处理身份验证
	routerGroup.Use(NewAuthMiddleware(m.domainService).Handle())
	
	// 3. 限流 - 防止API过度使用
	routerGroup.Use(NewRateLimitMiddleware(m.redisClient).Handle())
	
	// 4. 代理 - 转发请求到后端服务
	routerGroup.Use(ProxyMiddleware())
}

// GetRoutingMiddleware 获取路由中间件
func (m *MiddlewareManager) GetRoutingMiddleware() gin.HandlerFunc {
	return NewRoutingMiddleware(m.domainService).Handle()
}

// GetAuthMiddleware 获取认证中间件
func (m *MiddlewareManager) GetAuthMiddleware() gin.HandlerFunc {
	return NewAuthMiddleware(m.domainService).Handle()
}

// GetRateLimitMiddleware 获取限流中间件
func (m *MiddlewareManager) GetRateLimitMiddleware() gin.HandlerFunc {
	return NewRateLimitMiddleware(m.redisClient).Handle()
}

// GetProxyMiddleware 获取代理中间件
func (m *MiddlewareManager) GetProxyMiddleware() gin.HandlerFunc {
	return ProxyMiddleware()
}
