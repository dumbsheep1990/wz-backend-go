package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/go-redis/redis/v8"
	"github.com/casbin/casbin/v2"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务接口
type AuthService interface {
	// 生成JWT令牌
	GenerateToken(ctx context.Context, userID int64, role string, tenantID *int64) (string, error)
	// 验证JWT令牌
	VerifyToken(ctx context.Context, tokenString string) (*Claims, error)
	// 刷新JWT令牌
	RefreshToken(ctx context.Context, tokenString string) (string, error)
	// 撤销JWT令牌
	RevokeToken(ctx context.Context, tokenString string) error
	// 验证密码
	VerifyPassword(password, hashedPassword string) bool
	// 哈希密码
	HashPassword(password string) (string, error)
	// 检查用户权限
	CheckPermission(ctx context.Context, userID int64, role string, tenantID *int64, obj string, act string) (bool, error)
	// 检查多因素认证
	CheckMFA(ctx context.Context, userID int64, code string) (bool, error)
	// 生成多因素认证密钥
	GenerateMFASecret(ctx context.Context, userID int64) (string, string, error)
}

// Claims JWT令牌的声明
type Claims struct {
	UserID   int64   `json:"uid"`
	Role     string  `json:"role"`
	TenantID *int64  `json:"tid,omitempty"`
	jwt.RegisteredClaims
}

type authService struct {
	jwtSecret     []byte
	jwtExpiration time.Duration
	redis         *redis.Client
	enforcer      *casbin.Enforcer
	mfaService    MFAService
}

// NewAuthService 创建认证服务
func NewAuthService(
	jwtSecret string,
	jwtExpiration time.Duration,
	redis *redis.Client,
	enforcer *casbin.Enforcer,
	mfaService MFAService,
) AuthService {
	return &authService{
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
		redis:         redis,
		enforcer:      enforcer,
		mfaService:    mfaService,
	}
}

// GenerateToken 生成JWT令牌
func (s *authService) GenerateToken(ctx context.Context, userID int64, role string, tenantID *int64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiration)

	claims := Claims{
		UserID:   userID,
		Role:     role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "wz-backend-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	// 存储令牌到Redis以支持撤销功能
	tokenKey := fmt.Sprintf("token:%d", userID)
	err = s.redis.Set(ctx, tokenKey, tokenString, s.jwtExpiration).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken 验证JWT令牌
func (s *authService) VerifyToken(ctx context.Context, tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌并获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 检查令牌是否被撤销
	tokenKey := fmt.Sprintf("token:%d", claims.UserID)
	storedToken, err := s.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("token revoked")
		}
		return nil, err
	}

	if storedToken != tokenString {
		return nil, errors.New("token revoked")
	}

	return claims, nil
}

// RefreshToken 刷新JWT令牌
func (s *authService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// 验证旧令牌
	claims, err := s.VerifyToken(ctx, tokenString)
	if err != nil {
		return "", err
	}

	// 生成新令牌
	return s.GenerateToken(ctx, claims.UserID, claims.Role, claims.TenantID)
}

// RevokeToken 撤销JWT令牌
func (s *authService) RevokeToken(ctx context.Context, tokenString string) error {
	// 验证令牌
	claims, err := s.VerifyToken(ctx, tokenString)
	if err != nil {
		return err
	}

	// 从Redis中删除令牌
	tokenKey := fmt.Sprintf("token:%d", claims.UserID)
	return s.redis.Del(ctx, tokenKey).Err()
}

// VerifyPassword 验证密码
func (s *authService) VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPassword 哈希密码
func (s *authService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPermission 检查用户权限
func (s *authService) CheckPermission(ctx context.Context, userID int64, role string, tenantID *int64, obj string, act string) (bool, error) {
	var sub string
	if tenantID != nil {
		// 租户内权限检查
		sub = fmt.Sprintf("%d:%s:%d", userID, role, *tenantID)
	} else {
		// 平台级权限检查
		sub = fmt.Sprintf("%d:%s", userID, role)
	}

	return s.enforcer.Enforce(sub, obj, act)
}

// CheckMFA 检查多因素认证
func (s *authService) CheckMFA(ctx context.Context, userID int64, code string) (bool, error) {
	return s.mfaService.ValidateCode(ctx, userID, code)
}

// GenerateMFASecret 生成多因素认证密钥
func (s *authService) GenerateMFASecret(ctx context.Context, userID int64) (string, string, error) {
	return s.mfaService.GenerateSecret(ctx, userID)
}
