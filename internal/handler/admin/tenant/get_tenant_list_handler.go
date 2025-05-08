package tenant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	httphttp "wz-backend-go/api/http"
	"wz-backend-go/internal/logic/admin/tenant"
	"wz-backend-go/internal/svc"
)

func GetTenantListHandler(svcCtx *svc.AdminServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req httphttp.TenantListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := tenant.NewGetTenantListLogic(r.Context(), svcCtx)
		resp, err := l.GetTenantList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
} 