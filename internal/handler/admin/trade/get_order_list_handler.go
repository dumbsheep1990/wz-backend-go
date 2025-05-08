package trade

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	httphttp "wz-backend-go/api/http"
	"wz-backend-go/internal/logic/admin/trade"
	"wz-backend-go/internal/svc"
)

func GetOrderListHandler(svcCtx *svc.AdminServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req httphttp.OrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := trade.NewGetOrderListLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
} 