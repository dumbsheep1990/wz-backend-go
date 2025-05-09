package handler

import (
	"net/http"
	"wz-backend-go/internal/admin/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 注册所有路由处理器
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册通用中间件
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// CORS设置
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// 继续处理请求
			next(w, r)
		}
	})

	// 注册路由
	registerSystemHandlers(server, serverCtx)
	registerUserHandlers(server, serverCtx)
	registerContentHandlers(server, serverCtx)
	registerTenantHandlers(server, serverCtx)
	registerTradeHandlers(server, serverCtx)
	registerInteractionHandlers(server, serverCtx)
	registerAIHandlers(server, serverCtx)
	registerAdHandlers(server, serverCtx)        // 注册广告管理路由
	registerRecommendHandlers(server, serverCtx) // 注册推荐管理路由
}

// 注册系统管理相关路由
func registerSystemHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现系统管理路由
}

// 注册用户管理相关路由
func registerUserHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 用户管理路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/users",
		Handler: GetUserListHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/users/:id",
		Handler: GetUserDetailHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/admin/users",
		Handler: CreateUserHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/v1/admin/users/:id",
		Handler: UpdateUserHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/v1/admin/users/:id",
		Handler: DeleteUserHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))
}

// 注册内容管理相关路由
func registerContentHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现内容管理路由
}

// 注册租户管理相关路由
func registerTenantHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现租户管理路由
}

// 注册交易管理相关路由
func registerTradeHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现交易管理路由
}

// 注册交互管理相关路由
func registerInteractionHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现交互管理路由
}

// 注册AI管理相关路由
func registerAIHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// TODO: 实现AI管理路由
}

// 注册广告管理相关路由
func registerAdHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 广告位管理路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/ad/spaces",
		Handler: GetAdSpaceListHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/ad/spaces/:id",
		Handler: GetAdSpaceDetailHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/admin/ad/spaces",
		Handler: CreateAdSpaceHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/v1/admin/ad/spaces/:id",
		Handler: UpdateAdSpaceHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/v1/admin/ad/spaces/:id",
		Handler: DeleteAdSpaceHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))
}

// 注册推荐管理相关路由
func registerRecommendHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 推荐规则管理路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/recommend/rules",
		Handler: GetRecommendRulesHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/admin/recommend/rules",
		Handler: SaveRecommendRuleHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/admin/recommend/content/weight",
		Handler: SetContentWeightHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	// 热门数据路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/recommend/hot/content",
		Handler: GetHotContentHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/recommend/hot/keywords",
		Handler: GetHotKeywordsHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/admin/recommend/hot/categories",
		Handler: GetHotCategoriesHandler(serverCtx),
	}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret))
}
