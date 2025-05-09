package handler

import (
	"net/http"

	"wz-backend-go/internal/admin/logic"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetRecommendRulesHandler 获取推荐规则列表处理器
func GetRecommendRulesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecommendRuleListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetRecommendRulesLogic(r.Context(), svcCtx)
		resp, err := l.GetRecommendRules(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// SaveRecommendRuleHandler 保存推荐规则处理器
func SaveRecommendRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveRecommendRuleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSaveRecommendRuleLogic(r.Context(), svcCtx)
		resp, err := l.SaveRecommendRule(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// SetContentWeightHandler 设置内容权重处理器
func SetContentWeightHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetContentWeightReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSetContentWeightLogic(r.Context(), svcCtx)
		resp, err := l.SetContentWeight(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetHotContentHandler 获取热门内容处理器
func GetHotContentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotContentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetHotContentLogic(r.Context(), svcCtx)
		resp, err := l.GetHotContent(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetHotKeywordsHandler 获取热门关键词处理器
func GetHotKeywordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotKeywordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetHotKeywordsLogic(r.Context(), svcCtx)
		resp, err := l.GetHotKeywords(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetHotCategoriesHandler 获取热门分类处理器
func GetHotCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotCategoriesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetHotCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.GetHotCategories(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
