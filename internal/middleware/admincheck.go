package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminCheckMiddleware 管理员权限检查中间件
type AdminCheckMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewAdminCheckMiddleware 创建新的管理员检查中间件
func NewAdminCheckMiddleware(enforcer *casbin.Enforcer) *AdminCheckMiddleware {
	return &AdminCheckMiddleware{
		enforcer: enforcer,
	}
}

// Handle 处理请求
func (m *AdminCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文中获取用户ID和角色
		userId := r.Context().Value("userId")
		if userId == nil {
			httpx.Error(w, ErrUnauthorized)
			return
		}

		// 获取用户角色
		role := r.Context().Value("role")
		if role == nil {
			httpx.Error(w, ErrForbidden)
			return
		}

		// 获取请求路径和方法
		path := r.URL.Path
		method := r.Method

		// 简化路径（去除版本号和API前缀）
		path = simplifyPath(path)

		// 检查权限
		allowed, err := m.enforcer.Enforce(role, path, method)
		if err != nil {
			logx.Errorf("权限检查错误: %v", err)
			httpx.Error(w, ErrInternalServer)
			return
		}

		if !allowed {
			httpx.Error(w, ErrForbidden)
			return
		}

		// 继续处理请求
		next(w, r)
	}
}

// 简化路径，处理形如 /api/v1/admin/users 的路径
func simplifyPath(path string) string {
	// 去除API版本
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		// 例如 /api/v1/admin/users -> /admin/users
		result := "/" + strings.Join(parts[3:], "/")
		return result
	}
	return path
}

// 错误定义
var (
	ErrUnauthorized   = NewCustomError(http.StatusUnauthorized, "未授权")
	ErrForbidden      = NewCustomError(http.StatusForbidden, "权限不足")
	ErrInternalServer = NewCustomError(http.StatusInternalServerError, "服务器内部错误")
)

// CustomError 自定义错误
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (e CustomError) Error() string {
	return e.Message
}

// NewCustomError 创建新的自定义错误
func NewCustomError(code int, message string) CustomError {
	return CustomError{
		Code:    code,
		Message: message,
	}
}
