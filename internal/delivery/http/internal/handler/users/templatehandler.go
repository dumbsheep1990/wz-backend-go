package users

import (
	"net/http"
	"strconv"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/logic/users"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetTemplatesHandler 获取模板列表处理器
func GetTemplatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取分页参数
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}

		l := users.NewGetTemplatesLogic(r.Context(), svcCtx)
		resp, err := l.GetTemplates(userID, page, pageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// GetTemplateHandler 获取单个模板处理器
func GetTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取模板ID
		templateID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil || templateID <= 0 {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("无效的模板ID"))
			return
		}

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}

		l := users.NewGetTemplateLogic(r.Context(), svcCtx)
		resp, err := l.GetTemplate(userID, templateID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// CreateTemplateHandler 创建模板处理器
func CreateTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTemplateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewCreateTemplateLogic(r.Context(), svcCtx)
		resp, err := l.CreateTemplate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// UpdateTemplateHandler 更新模板处理器
func UpdateTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTemplateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取模板ID
		templateID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil || templateID <= 0 {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("无效的模板ID"))
			return
		}
		req.TemplateID = templateID

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewUpdateTemplateLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTemplate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// DeleteTemplateHandler 删除模板处理器
func DeleteTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取模板ID
		templateID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil || templateID <= 0 {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("无效的模板ID"))
			return
		}

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}

		l := users.NewDeleteTemplateLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTemplate(userID, templateID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// UpdateTemplateStatusHandler 更新模板状态处理器
func UpdateTemplateStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTemplateStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取模板ID
		templateID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil || templateID <= 0 {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("无效的模板ID"))
			return
		}
		req.TemplateID = templateID

		// 获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewUnauthorizedError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewUpdateTemplateStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTemplateStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
