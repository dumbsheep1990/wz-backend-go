package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/gateway/config"
)

// DynamicRouter 动态路由管理器
type DynamicRouter struct {
	engine      *gin.Engine
	routeGroups map[string]*gin.RouterGroup
	services    map[string]config.ServiceConfig
	handlers    map[string]gin.HandlerFunc
	mu          sync.RWMutex
}

// NewDynamicRouter 创建新的动态路由管理器
func NewDynamicRouter(engine *gin.Engine) *DynamicRouter {
	return &DynamicRouter{
		engine:      engine,
		routeGroups: make(map[string]*gin.RouterGroup),
		services:    make(map[string]config.ServiceConfig),
		handlers:    make(map[string]gin.HandlerFunc),
	}
}

// RegisterService 注册或更新服务
func (r *DynamicRouter) RegisterService(service config.ServiceConfig, registerFunc func(*gin.RouterGroup, config.ServiceConfig)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查服务是否已存在
	if _, exists := r.services[service.Name]; exists {
		// 如果存在，则先删除旧的路由
		r.removeServiceRoutes(service.Name)
	}

	// 保存服务配置
	r.services[service.Name] = service

	// 获取或创建路由组
	var group *gin.RouterGroup
	if existingGroup, exists := r.routeGroups[service.Prefix]; exists {
		group = existingGroup
	} else {
		group = r.engine.Group(service.Prefix)
		r.routeGroups[service.Prefix] = group
	}

	// 应用路由注册函数
	registerFunc(group, service)
}

// DeregisterService 注销服务
func (r *DynamicRouter) DeregisterService(serviceName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 删除服务路由
	r.removeServiceRoutes(serviceName)

	// 从映射中删除服务
	delete(r.services, serviceName)
}

// removeServiceRoutes 删除服务的所有路由
// 注意：这是内部方法，调用前必须已获取锁
func (r *DynamicRouter) removeServiceRoutes(serviceName string) {
	// 在Gin中，我们不能直接删除路由
	// 但我们可以通过记录哪些路由已被删除，在路由处理器中跳过这些路由
	service, exists := r.services[serviceName]
	if !exists {
		return
	}
	
	// 移除所有与此服务相关联的路由处理器
	for _, route := range service.Routes {
		routeKey := service.Prefix + route.Path + ":" + route.Method
		delete(r.handlers, routeKey)
	}
	
	// 如果服务类型是gRPC，还需要处理gRPC方法
	if service.Type == "grpc" {
		for _, method := range service.GrpcOptions.Methods {
			routeKey := service.Prefix + method.Path + ":" + method.HTTPMethod
			delete(r.handlers, routeKey)
		}
	}
}

// GetServices 获取所有注册的服务
func (r *DynamicRouter) GetServices() []config.ServiceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	services := make([]config.ServiceConfig, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	
	return services
}

// GetService 获取指定服务的配置
func (r *DynamicRouter) GetService(serviceName string) (config.ServiceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	service, exists := r.services[serviceName]
	return service, exists
}

// RegisterRouteHandler 注册路由处理器
func (r *DynamicRouter) RegisterRouteHandler(prefix, path, method string, handler gin.HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	routeKey := prefix + path + ":" + method
	r.handlers[routeKey] = handler
	
	// 获取或创建路由组
	var group *gin.RouterGroup
	if existingGroup, exists := r.routeGroups[prefix]; exists {
		group = existingGroup
	} else {
		group = r.engine.Group(prefix)
		r.routeGroups[prefix] = group
	}
	
	// 根据HTTP方法注册路由
	switch method {
	case "GET":
		group.GET(path, handler)
	case "POST":
		group.POST(path, handler)
	case "PUT":
		group.PUT(path, handler)
	case "DELETE":
		group.DELETE(path, handler)
	case "PATCH":
		group.PATCH(path, handler)
	case "HEAD":
		group.HEAD(path, handler)
	case "OPTIONS":
		group.OPTIONS(path, handler)
	default:
		group.Any(path, handler)
	}
}
