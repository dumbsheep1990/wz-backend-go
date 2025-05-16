package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/services/admin-service/internal/model"
)

// AdminCasbin Casbin权限校验中间件
func AdminCasbin(enforcer *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从JWT上下文中获取用户信息
			claims, ok := GetJWTClaims(r.Context())
			if !ok {
				httpx.WriteJson(w, http.StatusUnauthorized, model.CommonResponse{
					Code:    401,
					Message: "未找到登录信息",
				})
				return
			}

			// 获取请求的路径和方法
			path := r.URL.Path
			method := r.Method

			// 检查权限
			sub := claims.AuthorityId // 主体(角色)
			obj := path               // 客体(资源路径)
			act := method             // 操作(HTTP方法)
			
			// 检查白名单路径
			if isWhitelist(path) {
				next.ServeHTTP(w, r)
				return
			}

			// Casbin权限检查
			allowed, err := enforcer.Enforce(sub, obj, act)
			if err != nil {
				httpx.WriteJson(w, http.StatusInternalServerError, model.CommonResponse{
					Code:    500,
					Message: "权限检查失败: " + err.Error(),
				})
				return
			}

			if !allowed {
				httpx.WriteJson(w, http.StatusForbidden, model.CommonResponse{
					Code:    403,
					Message: "权限不足",
				})
				return
			}

			// 记录操作日志
			r = r.WithContext(WithOperationContext(r.Context(), model.OperationLog{
				UserId: claims.ID,
				Method: method,
				Path:   path,
			}))

			// 通过权限检查，继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// 白名单路径，不需要权限校验
func isWhitelist(path string) bool {
	whitelist := []string{
		"/api/v1/admin/login",
		"/api/v1/admin/refresh_token",
		"/api/v1/admin/captcha",
		"/api/v1/admin/user/getUserInfo", // 获取当前用户信息不需要权限校验
	}

	for _, p := range whitelist {
		if p == path {
			return true
		}
	}
	return false
}

// Operation上下文key
type operationContextKey struct{}

// WithOperationContext 添加操作日志上下文
func WithOperationContext(ctx context.Context, log model.OperationLog) context.Context {
	return context.WithValue(ctx, operationContextKey{}, log)
}

// GetOperationLog 从上下文中获取操作日志
func GetOperationLog(ctx context.Context) (model.OperationLog, bool) {
	log, ok := ctx.Value(operationContextKey{}).(model.OperationLog)
	return log, ok
} 