package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"wz-backend-go/internal/handler/admin/content"
	"wz-backend-go/internal/handler/admin/dashboard"
	"wz-backend-go/internal/handler/admin/tenant"
	"wz-backend-go/internal/handler/admin/trade"
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

	// 注册内容管理API
	registerContentHandlers(server, serverCtx)

	// 注册交易管理API
	registerTradeHandlers(server, serverCtx)

	// 注册积分管理API
	registerPointsHandlers(server, serverCtx)

	// 注册收藏管理API
	registerFavoritesHandlers(server, serverCtx)
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

// registerContentHandlers 注册内容管理API
func registerContentHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/contents",
					Handler: content.GetContentListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/contents/:id",
					Handler: content.GetContentDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/contents/:id/status",
					Handler: content.UpdateContentStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/contents/:id",
					Handler: content.DeleteContentHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/admin/contents/recommend",
					Handler: content.RecommendContentHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/contents/recommend/:id",
					Handler: content.CancelRecommendationHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/categories",
					Handler: content.GetCategoryListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/categories/:id",
					Handler: content.GetCategoryDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/admin/categories",
					Handler: content.CreateCategoryHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/categories/:id",
					Handler: content.UpdateCategoryHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/categories/:id",
					Handler: content.DeleteCategoryHandler(serverCtx),
				},
			},
		),
	)
}

// registerTradeHandlers 注册交易管理API
func registerTradeHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/orders",
					Handler: trade.GetOrderListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/orders/:id",
					Handler: trade.GetOrderDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/orders/:id/status",
					Handler: trade.UpdateOrderStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/refunds",
					Handler: trade.GetRefundListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/refunds/:id",
					Handler: trade.GetRefundDetailHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/refunds/:id/process",
					Handler: trade.ProcessRefundHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/transactions",
					Handler: trade.GetTransactionListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/reports/financial",
					Handler: trade.GetFinancialReportHandler(serverCtx),
				},
			},
		),
	)
}

// registerPointsHandlers 注册积分管理API
func registerPointsHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	pointsHandler := NewAdminPointsHandler(serverCtx.UserPointsService)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/points",
					Handler: serverCtx.WrapHandler(pointsHandler.ListPoints),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/points/:id",
					Handler: serverCtx.WrapHandler(pointsHandler.GetPointByID),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/users/:userId/points",
					Handler: serverCtx.WrapHandler(pointsHandler.GetUserPoints),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/admin/points",
					Handler: serverCtx.WrapHandler(pointsHandler.AddPoints),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/points/:id",
					Handler: serverCtx.WrapHandler(pointsHandler.DeletePoint),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/points/export",
					Handler: serverCtx.WrapHandler(pointsHandler.ExportPointsData),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/points/statistics",
					Handler: serverCtx.WrapHandler(pointsHandler.GetPointsStatistics),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/points/rules",
					Handler: serverCtx.WrapHandler(pointsHandler.GetPointsRules),
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/admin/points/rules",
					Handler: serverCtx.WrapHandler(pointsHandler.UpdatePointsRules),
				},
			},
		),
	)
}

// registerFavoritesHandlers 注册收藏管理API
func registerFavoritesHandlers(server *rest.Server, serverCtx *svc.AdminServiceContext) {
	favoritesHandler := NewAdminFavoritesHandler(serverCtx.UserFavoriteService)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites",
					Handler: serverCtx.WrapHandler(favoritesHandler.ListFavorites),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites/:id",
					Handler: serverCtx.WrapHandler(favoritesHandler.GetFavoriteByID),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/users/:userId/favorites",
					Handler: serverCtx.WrapHandler(favoritesHandler.GetUserFavorites),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/favorites/:id",
					Handler: serverCtx.WrapHandler(favoritesHandler.DeleteFavorite),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/admin/favorites/batch",
					Handler: serverCtx.WrapHandler(favoritesHandler.BatchDeleteFavorites),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites/export",
					Handler: serverCtx.WrapHandler(favoritesHandler.ExportFavoritesData),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites/statistics",
					Handler: serverCtx.WrapHandler(favoritesHandler.GetFavoritesStatistics),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites/hot",
					Handler: serverCtx.WrapHandler(favoritesHandler.GetHotContent),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/admin/favorites/trend",
					Handler: serverCtx.WrapHandler(favoritesHandler.GetFavoritesTrend),
				},
			},
		),
	)
}
