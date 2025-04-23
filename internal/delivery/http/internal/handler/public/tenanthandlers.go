package public

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/delivery/http/internal/logic/public"
	"wz-backend-go/internal/delivery/http/internal/middleware"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// GetTenantNavigationHandler 处理获取租户导航分类请求
func GetTenantNavigationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取语言参数，默认为zh-CN
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "zh-CN"
		}

		// 从上下文中获取租户ID
		tenantID, ok := middleware.GetTenantIDFromContext(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("未指定租户上下文"))
			return
		}

		l := public.NewGetTenantNavigationLogic(r.Context(), svcCtx)
		resp, err := l.GetTenantNavigation(tenantID, lang)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "导航分类获取成功"))
		}
	}
}

// SearchTenantHandler 处理租户内搜索请求
func SearchTenantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取参数
		q := r.URL.Query().Get("q")
		if q == "" {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("搜索关键词不能为空"))
			return
		}

		// 获取分页参数
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		
		page := 1
		limit := 10
		
		if pageStr != "" {
			pageNum, err := strconv.Atoi(pageStr)
			if err == nil && pageNum > 0 {
				page = pageNum
			}
		}
		
		if limitStr != "" {
			limitNum, err := strconv.Atoi(limitStr)
			if err == nil && limitNum > 0 {
				limit = limitNum
			}
		}

		// 从上下文中获取租户ID
		tenantID, ok := middleware.GetTenantIDFromContext(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("未指定租户上下文"))
			return
		}

		l := public.NewSearchTenantLogic(r.Context(), svcCtx)
		resp, err := l.SearchTenant(tenantID, q, page, limit)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "搜索成功"))
		}
	}
}

// GetRecommendationsHandler 处理获取推荐内容请求
func GetRecommendationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取参数
		limitStr := r.URL.Query().Get("limit")
		typeParam := r.URL.Query().Get("type")
		
		limit := 5
		
		if limitStr != "" {
			limitNum, err := strconv.Atoi(limitStr)
			if err == nil && limitNum > 0 {
				limit = limitNum
			}
		}

		// 从上下文中获取租户ID
		tenantID, ok := middleware.GetTenantIDFromContext(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("未指定租户上下文"))
			return
		}

		l := public.NewGetRecommendationsLogic(r.Context(), svcCtx)
		resp, err := l.GetRecommendations(tenantID, limit, typeParam)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "推荐内容获取成功"))
		}
	}
}

// GetCategoryDetailHandler 处理获取分类详情请求
func GetCategoryDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取路径参数
		vars := mux.Vars(r)
		categoryIDStr := vars["id"]
		
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("无效的分类ID"))
			return
		}

		// 从上下文中获取租户ID
		tenantID, ok := middleware.GetTenantIDFromContext(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("未指定租户上下文"))
			return
		}

		l := public.NewGetCategoryDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetCategoryDetail(tenantID, categoryID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "分类详情获取成功"))
		}
	}
}

// GetStaticPageHandler 处理获取静态页面请求
func GetStaticPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取路径参数
		vars := mux.Vars(r)
		pageName := vars["page"]
		
		if pageName == "" {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("页面标识符不能为空"))
			return
		}

		// 从上下文中获取租户ID
		tenantID, ok := middleware.GetTenantIDFromContext(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewErrorResponse("未指定租户上下文"))
			return
		}

		l := public.NewGetStaticPageLogic(r.Context(), svcCtx)
		resp, err := l.GetStaticPage(tenantID, pageName)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.NewSuccessResponse(resp, "静态页面获取成功"))
		}
	}
}
