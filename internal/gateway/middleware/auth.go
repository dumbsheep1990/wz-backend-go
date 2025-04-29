package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"wz-backend-go/internal/gateway/config"
	"github.com/gin-gonic/gin"
)

// 定义JWT声明结构
type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Authentication 返回一个身份验证中间件
func Authentication(securityConfig config.SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			return
		}

		// 提取令牌
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式无效",
			})
			return
		}
		tokenString := parts[1]

		// 解析并验证令牌
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无效的签名方法")
			}
			return []byte(securityConfig.JwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
				"error":   err.Error(),
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌已过期或无效",
			})
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("tenantID", claims.TenantID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// GenerateToken 生成新的JWT令牌
func GenerateToken(userID, tenantID, username, role string, expiration int, secret string) (string, error) {
	// 设置令牌的过期时间
	expirationTime := time.Now().Add(time.Duration(expiration) * time.Minute)
	
	// 创建JWT声明
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	
	// 使用指定的签名方法创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 使用密钥对令牌进行签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}
