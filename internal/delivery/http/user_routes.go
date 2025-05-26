package http

import (
	"net/http"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/handler/users"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/middleware"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterUserRoutes registers user related routes
func RegisterUserRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
	// Public routes
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/users/register",
				Handler: users.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/users/login",
				Handler: users.LoginHandler(serverCtx),
			},
		},
	)

	// Protected routes requiring authentication
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/users/profile",
				Handler: users.GetProfileHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/users/profile",
				Handler: users.UpdateProfileHandler(serverCtx),
			},
			// Enterprise registration routes
			{
				Method:  http.MethodPost,
				Path:    "/api/users/enterprise-registration",
				Handler: users.EnterpriseRegistrationHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/users/enterprise-registration",
				Handler: users.GetEnterpriseRegistrationHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/users/enterprise-registration",
				Handler: users.UpdateEnterpriseRegistrationHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/users/enterprise-registration/verify",
				Handler: users.VerifyEnterpriseHandler(serverCtx),
			},
		},
		rest.WithMiddlewares(
			[]rest.Middleware{
				middleware.NewAuthMiddleware(serverCtx).Handle,
			},
		),
	)
}
