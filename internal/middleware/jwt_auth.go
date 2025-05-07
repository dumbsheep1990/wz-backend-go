package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/handler"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// JwtAuthMiddleware JWT认证中间件
type JwtAuthMiddleware struct {
	accessSecret string
}

// NewJwtAuthMiddleware 创建JWT认证中间件
func NewJwtAuthMiddleware(accessSecret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		accessSecret: accessSecret,
	}
}

// CustomClaims 自定义JWT载荷
type CustomClaims struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Handle JWT认证中间件处理函数
func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.Error(w, handler.NewHttpError(http.StatusUnauthorized, "未授权，请先登录"))
			return
		}

		// 解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			httpx.Error(w, handler.NewHttpError(http.StatusUnauthorized, "认证格式有误，请使用Bearer <token>"))
			return
		}

		// 验证token
		token := parts[1]
		claims := &CustomClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.accessSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			logx.Errorf("解析token失败: %v", err)
			httpx.Error(w, handler.NewHttpError(http.StatusUnauthorized, "无效的token"))
			return
		}

		// 将解析出的用户ID和角色存入上下文
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "role", claims.Role)

		// 使用更新后的上下文继续处理请求
		next(w, r.WithContext(ctx))
	}
}
