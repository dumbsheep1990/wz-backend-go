package handler

import (
	"net/http"
	"strconv"

	"wz-backend-go/internal/admin/logic"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetAdSpaceListHandler 获取广告位列表处理器
func GetAdSpaceListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdSpaceListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetAdSpaceListLogic(r.Context(), svcCtx)
		resp, err := l.GetAdSpaceList(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetAdSpaceDetailHandler 获取广告位详情处理器
func GetAdSpaceDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取广告位ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.AdSpaceDetailReq{
			ID: id,
		}

		l := logic.NewGetAdSpaceDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetAdSpaceDetail(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// CreateAdSpaceHandler 创建广告位处理器
func CreateAdSpaceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateAdSpaceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateAdSpaceLogic(r.Context(), svcCtx)
		resp, err := l.CreateAdSpace(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// UpdateAdSpaceHandler 更新广告位处理器
func UpdateAdSpaceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取广告位ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		var req types.UpdateAdSpaceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		req.ID = id

		l := logic.NewUpdateAdSpaceLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAdSpace(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// DeleteAdSpaceHandler 删除广告位处理器
func DeleteAdSpaceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取广告位ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.DeleteAdSpaceReq{
			ID: id,
		}

		l := logic.NewDeleteAdSpaceLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAdSpace(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
