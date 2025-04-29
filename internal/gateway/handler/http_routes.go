package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"wz-backend-go/internal/gateway/config"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPRoutes 注册HTTP服务路由
func RegisterHTTPRoutes(group *gin.RouterGroup, service config.ServiceConfig) {
	// 处理配置中的每个路由
	for _, route := range service.Routes {
		// 构建完整路径
		path := route.Path
		
		// 创建对应的HTTP方法处理器
		switch strings.ToUpper(route.Method) {
		case "GET":
			group.GET(path, createHTTPProxy(route, service))
		case "POST":
			group.POST(path, createHTTPProxy(route, service))
		case "PUT":
			group.PUT(path, createHTTPProxy(route, service))
		case "DELETE":
			group.DELETE(path, createHTTPProxy(route, service))
		case "PATCH":
			group.PATCH(path, createHTTPProxy(route, service))
		case "HEAD":
			group.HEAD(path, createHTTPProxy(route, service))
		case "OPTIONS":
			group.OPTIONS(path, createHTTPProxy(route, service))
		case "ANY", "*":
			group.Any(path, createHTTPProxy(route, service))
		}
	}
	
	// 如果没有定义具体路由，则设置一个通配符路由
	if len(service.Routes) == 0 {
		group.Any("/*path", createServiceProxy(service))
	}
}

// createHTTPProxy 创建HTTP代理处理器
func createHTTPProxy(route config.RouteConfig, service config.ServiceConfig) gin.HandlerFunc {
	// 构建目标URL
	target := route.Target
	if target == "" {
		target = service.Target
	}
	
	// 解析目标URL
	targetURL, err := url.Parse(target)
	if err != nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无效的目标URL配置",
				"error":   err.Error(),
			})
		}
	}
	
	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	// 自定义导演函数
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// 调用原始导演函数
		originalDirector(req)
		
		// 修改请求路径
		if route.StripPath {
			path := req.URL.Path
			for _, prefix := range []string{service.Prefix, route.Path} {
				if prefix != "" && strings.HasPrefix(path, prefix) {
					req.URL.Path = strings.TrimPrefix(path, prefix)
					if req.URL.Path == "" {
						req.URL.Path = "/"
					}
					break
				}
			}
		}
		
		// 添加自定义请求头
		for key, value := range route.Headers {
			req.Header.Set(key, value)
		}
		
		// 添加自定义查询参数
		if len(route.QueryParams) > 0 {
			query := req.URL.Query()
			for key, value := range route.QueryParams {
				query.Set(key, value)
			}
			req.URL.RawQuery = query.Encode()
		}
		
		// 保留原始主机头
		req.Host = targetURL.Host
	}
	
	// 自定义错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// 直接使用http.ResponseWriter，不尝试创建gin.Context
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		
		// 创建错误响应
		response := map[string]interface{}{
			"code":    502,
			"message": "代理请求错误",
			"error":   err.Error(),
		}
		
		// 序列化为JSON并写入响应
		responseBytes, _ := json.Marshal(response)
		w.Write(responseBytes)
	}
	
	// 返回代理处理器
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// createServiceProxy 为整个服务创建代理处理器
func createServiceProxy(service config.ServiceConfig) gin.HandlerFunc {
	// 解析目标URL
	targetURL, err := url.Parse(service.Target)
	if err != nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无效的服务目标URL配置",
				"error":   err.Error(),
			})
		}
	}
	
	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	// 自定义导演函数
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// 调用原始导演函数
		originalDirector(req)
		
		// 修改请求路径
		path := req.URL.Path
		if strings.HasPrefix(path, service.Prefix) {
			req.URL.Path = strings.TrimPrefix(path, service.Prefix)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
		
		// 保留原始主机头
		req.Host = targetURL.Host
	}
	
	// 返回代理处理器
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
