package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/domain/model"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// CtxUserIDKey 上下文中的用户ID键
	CtxUserIDKey ContextKey = "user_id"
	// CtxTenantIDKey 上下文中的租户ID键
	CtxTenantIDKey ContextKey = "tenant_id"
	// CtxUserRoleKey 上下文中的用户角色键
	CtxUserRoleKey ContextKey = "user_role"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(accessSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从Authorization头中获取token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.Error(w, errors.New("未提供认证令牌"))
				return
			}

			// 检查认证方式是否为Bearer
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				httpx.Error(w, errors.New("认证格式无效"))
				return
			}

			// 解析JWT令牌
			token := parts[1]
			claims, err := parseToken(token, accessSecret)
			if err != nil {
				httpx.Error(w, errors.New("无效的令牌"))
				return
			}

			// 检查令牌是否过期
			if claims.ExpiresAt < time.Now().Unix() {
				httpx.Error(w, errors.New("令牌已过期"))
				return
			}

			// 将用户ID和租户ID加入请求上下文
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUserIDKey, claims.UserID)
			if claims.TenantID > 0 {
				ctx = context.WithValue(ctx, CtxTenantIDKey, claims.TenantID)
			}
			ctx = context.WithValue(ctx, CtxUserRoleKey, claims.Role)

			// 使用新的上下文继续请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalJWTAuthMiddleware 可选的JWT认证中间件，不强制要求认证
func OptionalJWTAuthMiddleware(accessSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从Authorization头中获取token
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// 检查认证方式是否为Bearer
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					// 解析JWT令牌
					token := parts[1]
					claims, err := parseToken(token, accessSecret)
					if err == nil && claims.ExpiresAt >= time.Now().Unix() {
						// 将用户ID和租户ID加入请求上下文
						ctx := r.Context()
						ctx = context.WithValue(ctx, CtxUserIDKey, claims.UserID)
						if claims.TenantID > 0 {
							ctx = context.WithValue(ctx, CtxTenantIDKey, claims.TenantID)
						}
						ctx = context.WithValue(ctx, CtxUserRoleKey, claims.Role)
						
						// 使用新的上下文继续请求
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					// 令牌无效或已过期，忽略并继续
					logx.Infof("无效或过期的可选令牌: %v", err)
				}
			}
			
			// 无令牌或令牌无效，继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// parseToken 解析JWT令牌
func parseToken(tokenString, secretKey string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateTokenPair 生成令牌对（访问令牌和刷新令牌）
func GenerateTokenPair(userID int64, username string, tenantID int64, role model.UserRole, accessSecret string, accessExpire int64) (*model.TokenPair, error) {
	now := time.Now().Unix()
	accessExpireTime := now + accessExpire
	
	// 创建访问令牌
	accessClaims := model.JWTClaims{
		UserID:     userID,
		Username:   username,
		TenantID:   tenantID,
		Role:       role,
		IssuedAt:   now,
		ExpiresAt:  accessExpireTime,
		TokenType:  "access",
	}
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}
	
	// 创建刷新令牌（有效期更长）
	refreshClaims := model.JWTClaims{
		UserID:     userID,
		Username:   username,
		TenantID:   tenantID,
		Role:       role,
		IssuedAt:   now,
		ExpiresAt:  now + (accessExpire * 7), // 刷新令牌有效期为访问令牌的7倍
		TokenType:  "refresh",
	}
	
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}
	
	return &model.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Unix(accessExpireTime, 0),
		TokenType:    "Bearer",
	}, nil
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(CtxUserIDKey).(int64)
	return userID, ok
}

// GetTenantIDFromContext 从上下文中获取租户ID
func GetTenantIDFromContext(ctx context.Context) (int64, bool) {
	tenantID, ok := ctx.Value(CtxTenantIDKey).(int64)
	return tenantID, ok
}

// GetUserRoleFromContext 从上下文中获取用户角色
func GetUserRoleFromContext(ctx context.Context) (model.UserRole, bool) {
	role, ok := ctx.Value(CtxUserRoleKey).(model.UserRole)
	return role, ok
}
