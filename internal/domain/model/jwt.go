package model

import (
	"time"
)

// JWTClaims JWT令牌声明
type JWTClaims struct {
	UserID       int64    `json:"user_id"`
	Username     string   `json:"username"`
	TenantID     int64    `json:"tenant_id,omitempty"` // 租户ID，可选，如果不指定则是平台级用户
	Role         UserRole `json:"role"`                // 用户角色
	ExpiresAt    int64    `json:"exp"`                 // 过期时间
	IssuedAt     int64    `json:"iat"`                 // 签发时间
	TokenType    string   `json:"token_type"`          // 令牌类型：access/refresh
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
