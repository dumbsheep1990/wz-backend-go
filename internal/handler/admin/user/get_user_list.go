package user

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/logic/admin/user"
	"wz-backend-go/internal/svc"
)

// GetUserListHandler 处理获取用户列表请求
func GetUserListHandler(svcCtx *svc.AdminServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析查询参数
		page, _ := strconv.Atoi(r.FormValue("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
		if pageSize <= 0 {
			pageSize = 10
		}

		req := &user.UserListReq{
			Page:      page,
			PageSize:  pageSize,
			Username:  r.FormValue("username"),
			Email:     r.FormValue("email"),
			Phone:     r.FormValue("phone"),
			StartTime: r.FormValue("startTime"),
			EndTime:   r.FormValue("endTime"),
		}

		// 将status字符串转为整数
		if statusStr := r.FormValue("status"); statusStr != "" {
			if status, err := strconv.Atoi(statusStr); err == nil {
				req.Status = status
			}
		}

		// 设置role
		req.Role = r.FormValue("role")

		// 创建逻辑层对象并调用
		l := user.NewGetUserListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserList(req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
