package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/handler"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminCheck 检查用户是否有管理员权限
func AdminCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从JWT中获取用户角色信息，通常在go-zero中JWT信息已经通过JwtAuthMiddleware解析
		// 并存储在请求上下文中
		role := r.Context().Value("role")
		if role == nil {
			logx.Errorf("未找到用户角色信息")
			httpx.Error(w, handler.NewHttpError(http.StatusForbidden, "没有管理员权限"))
			return
		}

		// role是interface{}类型，需要转为string
		roleStr, ok := role.(string)
		if !ok {
			logx.Errorf("角色信息类型错误: %v", role)
			httpx.Error(w, handler.NewHttpError(http.StatusForbidden, "没有管理员权限"))
			return
		}

		// 检查角色是否是管理员
		if roleStr != "platform_admin" && !strings.HasPrefix(roleStr, "admin") {
			logx.Errorf("用户没有管理员权限: %s", roleStr)
			httpx.Error(w, handler.NewHttpError(http.StatusForbidden, "没有管理员权限"))
			return
		}

		next(w, r)
	}
}
