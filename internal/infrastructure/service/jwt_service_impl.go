package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"

	appService "wz-backend-go/internal/application/admin/service"
)

// JWTConfig contains configuration for JWT token generation and validation
type JWTConfig struct {
	SecretKey     string
	ExpirationMin int
	Issuer        string
}

// JWTClaims contains the standard JWT claims plus custom claims
type JWTClaims struct {
	AdminID  int64  `json:"adminId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTServiceImpl implements JWTService interface
type JWTServiceImpl struct {
	config JWTConfig
	redis  *redis.Client
}

// NewJWTService creates a new JWTServiceImpl
func NewJWTService(config JWTConfig, redis *redis.Client) *JWTServiceImpl {
	return &JWTServiceImpl{
		config: config,
		redis:  redis,
	}
}

// Ensure JWTServiceImpl implements the JWTService interface
var _ appService.JWTService = (*JWTServiceImpl)(nil)

// GenerateToken creates a new JWT token for an admin user
func (s *JWTServiceImpl) GenerateToken(adminID int64, username string) (string, int64, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Minute * time.Duration(s.config.ExpirationMin))
	expiresAt := expirationTime.Unix()

	// Create claims
	claims := &JWTClaims{
		AdminID:  adminID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    s.config.Issuer,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTServiceImpl) ValidateToken(tokenString string) (int64, string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, "", errors.New("could not parse claims")
	}

	// Check if token is revoked
	isRevoked, err := s.IsTokenRevoked(context.Background(), tokenString)
	if err != nil {
		return 0, "", err
	}
	if isRevoked {
		return 0, "", errors.New("token has been revoked")
	}

	return claims.AdminID, claims.Username, nil
}

// RevokeToken blacklists a token to prevent its further use
func (s *JWTServiceImpl) RevokeToken(ctx context.Context, tokenString string) error {
	// Parse the token without validation to get the expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return err
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return errors.New("could not parse claims")
	}

	// Calculate TTL for the revoked token in Redis
	// The TTL should be the time until the token expires
	now := time.Now().Unix()
	ttl := time.Duration(claims.ExpiresAt - now) * time.Second
	if ttl <= 0 {
		// Token already expired, no need to revoke
		return nil
	}

	// Store token in Redis blacklist with TTL
	key := fmt.Sprintf("token:blacklist:%s", tokenString)
	err = s.redis.Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (s *JWTServiceImpl) IsTokenRevoked(ctx context.Context, tokenString string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", tokenString)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
