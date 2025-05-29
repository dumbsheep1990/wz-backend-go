package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取认证头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// TODO: 实际项目中应该验证令牌并解析用户ID和角色
		// 这里简化处理，假设已经验证通过并设置用户ID
		c.Set("user_id", int64(1))
		c.Set("role", "admin")

		c.Next()
	}
}

// UserAuth 普通用户认证中间件
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取认证头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// TODO: 实际项目中应该验证令牌并解析用户ID和角色
		// 这里简化处理，假设已经验证通过并设置用户ID
		c.Set("user_id", int64(1))
		c.Set("role", "user")

		c.Next()
	}
}
