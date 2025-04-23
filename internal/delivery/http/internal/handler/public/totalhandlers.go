package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/delivery/http/internal/logic/public"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// GetTenantsHandler 处理获取租户列表请求
func GetTenantsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewGetTenantsLogic(r.Context(), svcCtx)
		resp, err := l.GetTenants()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 使用统一的响应格式
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "租户列表获取成功"))
		}
	}
}

// GetTotalNavigationHandler 处理获取全局导航请求
func GetTotalNavigationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取语言参数，默认为zh-CN
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "zh-CN"
		}

		l := public.NewGetTotalNavigationLogic(r.Context(), svcCtx)
		resp, err := l.GetTotalNavigation(lang)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "全局导航获取成功"))
		}
	}
}
