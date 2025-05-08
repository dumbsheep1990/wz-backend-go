package trade

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wz-backend-go/internal/logic/admin/trade"
	"wz-backend-go/internal/svc"
)

func GetOrderDetailHandler(svcCtx *svc.AdminServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径中提取ID参数
		id := r.PathValue("id")
		if id == "" {
			httpx.Error(w, http.ErrMissingFile)
			return
		}

		l := trade.NewGetOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderDetail(id)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
} 