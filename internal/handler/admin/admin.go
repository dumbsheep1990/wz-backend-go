package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"wz-backend-go/internal/handler/admin/dashboard"
	"wz-backend-go/internal/handler/admin/tenant"
	"wz-backend-go/internal/handler/admin/user"
	"wz-backend-go/internal/svc"
)

// RegisterHandlers 注册所有后台管理API的处理函数
func RegisterHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	// 注册用户管理API
	registerUserHandlers(server, serverCtx)

	// 注册租户管理API
	registerTenantHandlers(server, serverCtx)

	// 注册仪表盘API
	registerDashboardHandlers(server, serverCtx)
}

// registerUserHandlers 注册用户管理API
func registerUserHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/users",
					Handler: user.GetUserListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/users/:id",
					Handler: user.GetUserDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/admin/users",
					Handler: user.CreateUserHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/users/:id",
					Handler: user.UpdateUserHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/users/:id",
					Handler: user.DeleteUserHandler(serverCtx),
				},
			},
		),
	)
}

// registerTenantHandlers 注册租户管理API
func registerTenantHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/tenants",
					Handler: tenant.GetTenantListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/tenants/:id",
					Handler: tenant.GetTenantDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/admin/tenants",
					Handler: tenant.CreateTenantHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/tenants/:id",
					Handler: tenant.UpdateTenantHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/tenants/:id",
					Handler: tenant.DeleteTenantHandler(serverCtx),
				},
			},
		),
	)
}

// registerDashboardHandlers 注册仪表盘API
func registerDashboardHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/dashboard/overview",
					Handler: dashboard.GetOverviewHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/settings",
					Handler: dashboard.GetSettingsHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/settings",
					Handler: dashboard.UpdateSettingHandler(serverCtx),
				},
			},
		),
	)
}
