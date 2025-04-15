package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/delivery/http/internal/logic/users"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
)

func GetUserBehaviorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserBehaviorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewGetUserBehaviorLogic(r.Context(), svcCtx)
		resp, err := l.GetUserBehavior(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
