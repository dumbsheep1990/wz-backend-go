package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	// 在实际生产环境中，这个密钥应该从配置或环境变量中获取
	jwtSecret = []byte("万知相同社区平台JWT密钥")
)

// JWTAuthMiddleware 创建JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			c.Abort()
			return
		}

		// 解析和验证令牌
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 验证令牌声明
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户ID和用户名存储在上下文中，供后续处理器使用
			c.Set("user_id", claims["user_id"])
			c.Set("user_name", claims["user_name"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌声明",
			})
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID, userName string) (string, error) {
	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
		"iat":       time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString(jwtSecret)
}

// 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}
