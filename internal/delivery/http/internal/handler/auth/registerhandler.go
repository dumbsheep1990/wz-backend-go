package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/delivery/http/internal/logic/auth"
	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
