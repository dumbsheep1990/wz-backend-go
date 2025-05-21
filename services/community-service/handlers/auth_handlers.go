package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wz-backend-go/services/community-service/middleware"
)

// 简单的内存用户存储，实际应用中应使用数据库
var users = map[string]User{
	"admin": {
		ID:       "user-admin",
		Username: "admin",
		Password: "admin123",
		Name:     "管理员",
	},
	"test": {
		ID:       "user-test",
		Username: "test",
		Password: "test123",
		Name:     "测试用户",
	},
}

// User 用户模型
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不输出密码
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", LoginHandler)
		auth.POST("/register", RegisterHandler)
		
		// 需要认证的路由
		authorized := auth.Group("")
		authorized.Use(middleware.JWTAuthMiddleware())
		{
			authorized.GET("/me", CurrentUserHandler)
		}
	}
}

// LoginHandler 处理用户登录
func LoginHandler(c *gin.Context) {
	var req middleware.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 验证用户
	user, exists := users[req.Username]
	if !exists || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成令牌
	token, err := middleware.GenerateToken(user.ID, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败: " + err.Error(),
		})
		return
	}

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, middleware.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		UserName: user.Name,
	})
}

// RegisterHandler 处理用户注册
func RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 检查用户名是否已存在
	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 创建新用户
	userID := "user-" + req.Username
	users[req.Username] = User{
		ID:        userID,
		Username:  req.Username,
		Password:  req.Password,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	// 生成令牌
	token, err := middleware.GenerateToken(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败: " + err.Error(),
		})
		return
	}

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, middleware.LoginResponse{
		Token:    token,
		UserID:   userID,
		UserName: req.Name,
	})
}

// CurrentUserHandler 获取当前用户信息
func CurrentUserHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"user_name": userName,
	})
}
