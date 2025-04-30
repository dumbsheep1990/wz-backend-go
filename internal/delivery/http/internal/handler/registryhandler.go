package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/delivery/http/internal/logic"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// 注册所有服务注册与发现相关的处理器
func RegisterRegistryHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 获取所有服务列表
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/registry/services",
				Handler: ListServicesHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/registry/services/:serviceName/instances",
				Handler: GetServiceInstancesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/registry/services/:serviceName/instances",
				Handler: RegisterInstanceHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/v1/registry/services/:serviceName/instances/:instanceId",
				Handler: DeregisterInstanceHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/api/v1/registry/services/:serviceName/instances/:instanceId/status",
				Handler: UpdateInstanceStatusHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/registry/health/status",
				Handler: GetHealthStatusHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/registry/health/check",
				Handler: TriggerHealthCheckHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/registry/services/:serviceName/dependencies",
				Handler: GetServiceDependenciesHandler(serverCtx),
			},
		},
	)
}

// ListServicesHandler 获取所有服务列表处理器
func ListServicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListServicesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewListServicesLogic(r.Context(), svcCtx)
		resp, err := l.ListServices(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// GetServiceInstancesHandler 获取服务实例列表处理器
func GetServiceInstancesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetServiceInstancesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetServiceInstancesLogic(r.Context(), svcCtx)
		resp, err := l.GetServiceInstances(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// RegisterInstanceHandler 注册服务实例处理器
func RegisterInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterInstanceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRegisterInstanceLogic(r.Context(), svcCtx)
		resp, err := l.RegisterInstance(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// DeregisterInstanceHandler 注销服务实例处理器
func DeregisterInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeregisterInstanceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDeregisterInstanceLogic(r.Context(), svcCtx)
		resp, err := l.DeregisterInstance(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// UpdateInstanceStatusHandler 更新实例状态处理器
func UpdateInstanceStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateInstanceStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUpdateInstanceStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdateInstanceStatus(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// GetHealthStatusHandler 获取健康状态处理器
func GetHealthStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetHealthStatusLogic(r.Context(), svcCtx)
		resp, err := l.GetHealthStatus()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// TriggerHealthCheckHandler 触发健康检查处理器
func TriggerHealthCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TriggerHealthCheckReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewTriggerHealthCheckLogic(r.Context(), svcCtx)
		resp, err := l.TriggerHealthCheck(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// GetServiceDependenciesHandler 获取服务依赖关系处理器
func GetServiceDependenciesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetServiceDependenciesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetServiceDependenciesLogic(r.Context(), svcCtx)
		resp, err := l.GetServiceDependencies(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
