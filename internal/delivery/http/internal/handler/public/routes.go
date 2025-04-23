package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"wz-backend-go/internal/delivery/http/internal/middleware"
	"wz-backend-go/internal/delivery/http/internal/svc"
)

// RegisterPublicHandlers 注册公共API处理器
func RegisterPublicHandlers(server *rest.Server, serverCtx *svc.ServiceContext, mainDomain string) {
	// 总站公共接口 - 无需租户上下文
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/total/tenants",
				Handler: GetTenantsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/total/navigation",
				Handler: GetTotalNavigationHandler(serverCtx),
			},
		},
	)

	// 租户专属公共接口 - 需要租户上下文
	// 这些路由需要租户中间件来解析租户ID
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/navigation",
				Handler: GetTenantNavigationHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/search",
				Handler: SearchTenantHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/recommendations",
				Handler: GetRecommendationsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/category/{id}",
				Handler: GetCategoryDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/static/{page}",
				Handler: GetStaticPageHandler(serverCtx),
			},
		},
		rest.WithMiddlewares(
			// 可选的JWT认证中间件，即使未认证也可以访问，但会解析令牌
			middleware.OptionalJWTAuthMiddleware(serverCtx.Config.Auth.AccessSecret),
			// 租户中间件，解析租户ID
			middleware.TenantMiddleware(mainDomain, serverCtx.TenantService),
		),
	)
}
