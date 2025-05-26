package users

import (
	"net/http"

	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/logic/users"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/svc"
	"github.com/wz-project/wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// EnterpriseRegistrationHandler handles enterprise registration requests
func EnterpriseRegistrationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnterpriseRegistrationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Extract user ID from context (set by authentication middleware)
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewInternalError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewEnterpriseRegistrationLogic(r.Context(), svcCtx)
		resp, err := l.EnterpriseRegistration(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// GetEnterpriseRegistrationHandler handles get enterprise registration requests
func GetEnterpriseRegistrationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context (set by authentication middleware)
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewInternalError("未授权访问"))
			return
		}

		l := users.NewGetEnterpriseRegistrationLogic(r.Context(), svcCtx)
		resp, err := l.GetEnterpriseRegistration(userID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// UpdateEnterpriseRegistrationHandler handles update enterprise registration requests
func UpdateEnterpriseRegistrationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateEnterpriseRegistrationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Extract user ID from context (set by authentication middleware)
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewInternalError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewUpdateEnterpriseRegistrationLogic(r.Context(), svcCtx)
		resp, err := l.UpdateEnterpriseRegistration(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// VerifyEnterpriseHandler handles verify enterprise requests
func VerifyEnterpriseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyEnterpriseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Extract user ID from context (set by authentication middleware)
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			httpx.ErrorCtx(r.Context(), w, types.NewInternalError("未授权访问"))
			return
		}
		req.UserID = userID

		l := users.NewVerifyEnterpriseLogic(r.Context(), svcCtx)
		resp, err := l.VerifyEnterprise(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
