package handler

import (
	"net/http"
	"strconv"

	"wz-backend-go/internal/admin/logic"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetSectionListHandler 获取板块列表处理器
func GetSectionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SectionListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetSectionListLogic(r.Context(), svcCtx)
		resp, err := l.GetSectionList(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetSectionDetailHandler 获取板块详情处理器
func GetSectionDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取板块ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.SectionDetailReq{
			ID: id,
		}

		l := logic.NewGetSectionDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetSectionDetail(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// CreateSectionHandler 创建板块处理器
func CreateSectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSectionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateSectionLogic(r.Context(), svcCtx)
		resp, err := l.CreateSection(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// UpdateSectionHandler 更新板块处理器
func UpdateSectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取板块ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		var req types.UpdateSectionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		req.ID = id

		l := logic.NewUpdateSectionLogic(r.Context(), svcCtx)
		resp, err := l.UpdateSection(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// DeleteSectionHandler 删除板块处理器
func DeleteSectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取板块ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.DeleteSectionReq{
			ID: id,
		}

		l := logic.NewDeleteSectionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSection(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
