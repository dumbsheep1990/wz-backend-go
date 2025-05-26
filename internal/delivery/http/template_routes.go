package http

import (
	"net/http"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/handler/users"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/middleware"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterTemplateRoutes 注册模板管理相关路由
func RegisterTemplateRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 模板管理路由，需要认证
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/users/templates",
				Handler: users.GetTemplatesHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/users/templates/:id",
				Handler: users.GetTemplateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/users/templates",
				Handler: users.CreateTemplateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/users/templates/:id",
				Handler: users.UpdateTemplateHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/users/templates/:id",
				Handler: users.DeleteTemplateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/users/templates/:id/status",
				Handler: users.UpdateTemplateStatusHandler(serverCtx),
			},
		},
		rest.WithMiddlewares(
			[]rest.Middleware{
				middleware.NewAuthMiddleware(serverCtx).Handle,
			},
		),
	)
}
