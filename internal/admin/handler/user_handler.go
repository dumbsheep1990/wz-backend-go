package handler

import (
	"net/http"
	"strconv"

	"wz-backend-go/internal/admin/logic"
	"wz-backend-go/internal/admin/svc"
	"wz-backend-go/internal/admin/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetUserListHandler 获取用户列表处理器
func GetUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetUserListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserList(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// GetUserDetailHandler 获取用户详情处理器
func GetUserDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取用户ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.UserDetailReq{
			ID: id,
		}

		l := logic.NewGetUserDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetUserDetail(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// CreateUserHandler 创建用户处理器
func CreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminCreateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateUserLogic(r.Context(), svcCtx)
		resp, err := l.CreateUser(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// UpdateUserHandler 更新用户处理器
func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取用户ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		var req types.AdminUpdateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		req.ID = id

		l := logic.NewUpdateUserLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUser(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

// DeleteUserHandler 删除用户处理器
func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中获取用户ID
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		req := &types.DeleteUserReq{
			ID: id,
		}

		l := logic.NewDeleteUserLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUser(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
