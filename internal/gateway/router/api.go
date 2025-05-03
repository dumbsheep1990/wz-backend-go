package router

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/gateway/config"
)

// RouteAPI 处理路由管理API
type RouteAPI struct {
	router *DynamicRouter
}

// NewRouteAPI 创建新的路由API处理器
func NewRouteAPI(router *DynamicRouter) *RouteAPI {
	return &RouteAPI{
		router: router,
	}
}

// RegisterRouteAPI 注册路由管理API
func (api *RouteAPI) RegisterRouteAPI(r *gin.Engine) {
	routeGroup := r.Group("/api/v1/gateway/routes")
	{
		// 获取所有服务路由
		routeGroup.GET("", api.GetAllRoutes)
		// 获取指定服务的路由
		routeGroup.GET("/:serviceName", api.GetServiceRoutes)
		// 注册新服务
		routeGroup.POST("", api.RegisterService)
		// 更新服务
		routeGroup.PUT("/:serviceName", api.UpdateService)
		// 删除服务
		routeGroup.DELETE("/:serviceName", api.DeleteService)
		// 添加路由
		routeGroup.POST("/:serviceName/routes", api.AddRoute)
		// 删除路由
		routeGroup.DELETE("/:serviceName/routes/:routeID", api.DeleteRoute)
	}
}

// GetAllRoutes 获取所有路由信息
func (api *RouteAPI) GetAllRoutes(c *gin.Context) {
	services := api.router.GetServices()
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功获取所有路由信息",
		"data":    services,
	})
}

// GetServiceRoutes 获取指定服务的路由信息
func (api *RouteAPI) GetServiceRoutes(c *gin.Context) {
	serviceName := c.Param("serviceName")
	service, exists := api.router.GetService(serviceName)
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "服务不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功获取服务路由信息",
		"data":    service,
	})
}

// RegisterService 注册新服务
func (api *RouteAPI) RegisterService(c *gin.Context) {
	var service config.ServiceConfig
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的服务配置",
			"error":   err.Error(),
		})
		return
	}
	
	// 验证服务配置
	if service.Name == "" || service.Prefix == "" || service.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少必要的服务配置字段",
		})
		return
	}
	
	// 注册服务
	api.router.RegisterService(service, func(group *gin.RouterGroup, svc config.ServiceConfig) {
		// 这里应根据服务类型调用相应的路由注册函数
		// 简单实现，仅支持HTTP路由
		for _, route := range svc.Routes {
			api.registerRouteHandler(group, svc.Prefix, route)
		}
	})
	
	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "服务注册成功",
		"data":    service,
	})
}

// UpdateService 更新服务配置
func (api *RouteAPI) UpdateService(c *gin.Context) {
	serviceName := c.Param("serviceName")
	_, exists := api.router.GetService(serviceName)
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "服务不存在",
		})
		return
	}
	
	var service config.ServiceConfig
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的服务配置",
			"error":   err.Error(),
		})
		return
	}
	
	// 确保服务名称一致
	service.Name = serviceName
	
	// 更新服务
	api.router.RegisterService(service, func(group *gin.RouterGroup, svc config.ServiceConfig) {
		// 这里应根据服务类型调用相应的路由注册函数
		for _, route := range svc.Routes {
			api.registerRouteHandler(group, svc.Prefix, route)
		}
	})
	
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "服务更新成功",
		"data":    service,
	})
}

// DeleteService 删除服务
func (api *RouteAPI) DeleteService(c *gin.Context) {
	serviceName := c.Param("serviceName")
	_, exists := api.router.GetService(serviceName)
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "服务不存在",
		})
		return
	}
	
	// 删除服务
	api.router.DeregisterService(serviceName)
	
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "服务删除成功",
	})
}

// AddRoute 添加路由
func (api *RouteAPI) AddRoute(c *gin.Context) {
	serviceName := c.Param("serviceName")
	service, exists := api.router.GetService(serviceName)
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "服务不存在",
		})
		return
	}
	
	var route config.RouteConfig
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的路由配置",
			"error":   err.Error(),
		})
		return
	}
	
	// 验证路由配置
	if route.Path == "" || route.Method == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少必要的路由配置字段",
		})
		return
	}
	
	// 添加路由到服务
	service.Routes = append(service.Routes, route)
	
	// 重新注册服务
	api.router.RegisterService(service, func(group *gin.RouterGroup, svc config.ServiceConfig) {
		for _, r := range svc.Routes {
			api.registerRouteHandler(group, svc.Prefix, r)
		}
	})
	
	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "路由添加成功",
		"data":    route,
	})
}

// DeleteRoute 删除路由
func (api *RouteAPI) DeleteRoute(c *gin.Context) {
	serviceName := c.Param("serviceName")
	routeID := c.Param("routeID")
	
	service, exists := api.router.GetService(serviceName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "服务不存在",
		})
		return
	}
	
	// 查找并删除路由
	var updatedRoutes []config.RouteConfig
	var found bool
	for _, route := range service.Routes {
		if generateRouteID(route) != routeID {
			updatedRoutes = append(updatedRoutes, route)
		} else {
			found = true
		}
	}
	
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "路由不存在",
		})
		return
	}
	
	// 更新服务路由
	service.Routes = updatedRoutes
	
	// 重新注册服务
	api.router.RegisterService(service, func(group *gin.RouterGroup, svc config.ServiceConfig) {
		for _, route := range svc.Routes {
			api.registerRouteHandler(group, svc.Prefix, route)
		}
	})
	
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "路由删除成功",
	})
}

// registerRouteHandler 注册路由处理器
func (api *RouteAPI) registerRouteHandler(group *gin.RouterGroup, prefix string, route config.RouteConfig) {
	// 创建路由处理函数
	handler := createHTTPProxy(route, prefix)
	
	// 注册到路由管理器
	api.router.RegisterRouteHandler(prefix, route.Path, route.Method, handler)
}

// generateRouteID 生成路由标识符
func generateRouteID(route config.RouteConfig) string {
	// 简单实现：使用路径和方法组合作为ID
	// 实际应用中可能需要更复杂的唯一标识生成方法
	routeMap := map[string]interface{}{
		"path":   route.Path,
		"method": route.Method,
	}
	
	jsonBytes, _ := json.Marshal(routeMap)
	return string(jsonBytes)
}

// createHTTPProxy 创建HTTP代理处理器
func createHTTPProxy(route config.RouteConfig, prefix string) gin.HandlerFunc {
	// 这里应该复用http_routes.go中的代理逻辑
	// 简单实现，返回一个空的处理器
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "路由已注册，但处理器逻辑尚未实现",
			"route":   route,
		})
	}
}
