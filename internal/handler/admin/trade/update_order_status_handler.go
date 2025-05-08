package trade

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	httphttp "wz-backend-go/api/http"
	"wz-backend-go/internal/logic/admin/trade"
	"wz-backend-go/internal/svc"
)

func UpdateOrderStatusHandler(svcCtx *svc.AdminServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中提取ID参数
		id := r.PathValue("id")
		if id == "" {
			httpx.Error(w, http.ErrMissingFile)
			return
		}

		var req httphttp.UpdateOrderStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := trade.NewUpdateOrderStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdateOrderStatus(id, &req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
} 