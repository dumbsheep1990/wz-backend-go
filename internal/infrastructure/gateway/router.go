package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// RouterImpl 路由管理器实现
type RouterImpl struct {
	routeRepo repository.RouteRepository
	routes    map[string]*entity.Route
	mutex     sync.RWMutex
}

// NewRouter 创建新的路由管理器
func NewRouter(routeRepo repository.RouteRepository) *RouterImpl {
	router := &RouterImpl{
		routeRepo: routeRepo,
		routes:    make(map[string]*entity.Route),
	}
	
	// 初始加载所有路由
	ctx := context.Background()
	if err := router.LoadRoutes(ctx); err != nil {
		log.Printf("初始路由加载失败: %v", err)
	}
	
	return router
}

// LoadRoutes 从存储中加载所有路由
func (r *RouterImpl) LoadRoutes(ctx context.Context) error {
	routes, err := r.routeRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("从存储加载路由失败: %w", err)
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// 清空当前路由表并重新加载
	r.routes = make(map[string]*entity.Route)
	for _, route := range routes {
		r.routes[route.ID] = route
	}
	
	log.Printf("已加载 %d 个路由", len(r.routes))
	return nil
}

// ResolveRoute 根据HTTP请求解析匹配的路由
func (r *RouterImpl) ResolveRoute(req *http.Request) (*entity.Route, error) {
	path := req.URL.Path
	method := req.Method
	
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// 按照路径精确度排序（从最精确到最不精确）匹配路由
	var matchedRoute *entity.Route
	var maxPathSegments int
	
	for _, route := range r.routes {
		// 检查路由是否启用
		if !route.IsActive {
			continue
		}
		
		// 检查HTTP方法是否匹配
		methodMatch := false
		for _, m := range route.Methods {
			if m == method || m == "*" {
				methodMatch = true
				break
			}
		}
		
		if !methodMatch {
			continue
		}
		
		// 路径匹配
		pathMatch := false
		
		// 精确匹配
		if route.PathType == valueobject.ExactPath && route.Path == path {
			pathMatch = true
		}
		
		// 前缀匹配
		if route.PathType == valueobject.PrefixPath && strings.HasPrefix(path, route.Path) {
			pathMatch = true
		}
		
		// 正则表达式匹配
		if route.PathType == valueobject.RegexPath && route.CompiledRegex != nil {
			if route.CompiledRegex.MatchString(path) {
				pathMatch = true
			}
		}
		
		if pathMatch {
			// 计算路径段数，用于确定最精确的匹配
			segments := len(strings.Split(route.Path, "/"))
			if matchedRoute == nil || segments > maxPathSegments {
				matchedRoute = route
				maxPathSegments = segments
			}
		}
	}
	
	if matchedRoute == nil {
		return nil, errors.New("没有找到匹配的路由")
	}
	
	return matchedRoute, nil
}

// AddRoute 添加新路由
func (r *RouterImpl) AddRoute(route *entity.Route) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.routes[route.ID] = route
}

// RemoveRoute 移除路由
func (r *RouterImpl) RemoveRoute(routeID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.routes, routeID)
}

// GetAllRoutes 获取所有路由
func (r *RouterImpl) GetAllRoutes() []*entity.Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	routes := make([]*entity.Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	return routes
}
