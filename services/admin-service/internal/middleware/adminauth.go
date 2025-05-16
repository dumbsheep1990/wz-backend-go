package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/services/admin-service/internal/model"
)

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenNotProvided = errors.New("token is not provided")
)

// AdminJWT 管理员JWT认证中间件
func AdminJWT(accessSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求头中获取token
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				httpx.WriteJson(w, http.StatusUnauthorized, model.CommonResponse{
					Code:    401,
					Message: "请先登录",
				})
				return
			}

			// 去掉Bearer前缀
			tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

			// 解析token
			token, err := jwt.ParseWithClaims(tokenStr, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(accessSecret), nil
			})

			if err != nil {
				var message string
				if errors.Is(err, ErrTokenExpired) {
					message = "登录已过期"
				} else {
					message = "无效的token"
				}
				httpx.WriteJson(w, http.StatusUnauthorized, model.CommonResponse{
					Code:    401,
					Message: message,
				})
				return
			}

			claims, ok := token.Claims.(*model.JwtClaims)
			if !ok || !token.Valid {
				httpx.WriteJson(w, http.StatusUnauthorized, model.CommonResponse{
					Code:    401,
					Message: "登录状态校验失败",
				})
				return
			}

			// 检查token是否过期
			if time.Now().Unix() > claims.ExpiresAt {
				httpx.WriteJson(w, http.StatusUnauthorized, model.CommonResponse{
					Code:    401,
					Message: "登录已过期",
				})
				return
			}

			// 将用户信息存储到请求上下文中
			r = r.WithContext(WithJWTContext(r.Context(), claims))
			
			// 调用下一个处理器
			next.ServeHTTP(w, r)
		})
	}
}

// JWT上下文key
type jwtContextKey struct{}

// WithJWTContext 添加JWT上下文
func WithJWTContext(ctx context.Context, claims *model.JwtClaims) context.Context {
	return context.WithValue(ctx, jwtContextKey{}, claims)
}

// GetJWTClaims 从上下文中获取JWT信息
func GetJWTClaims(ctx context.Context) (*model.JwtClaims, bool) {
	claims, ok := ctx.Value(jwtContextKey{}).(*model.JwtClaims)
	return claims, ok
}
