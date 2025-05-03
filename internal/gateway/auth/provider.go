package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"wz-backend-go/internal/gateway/config"
)

// AuthUser 表示认证后的用户信息
type AuthUser struct {
	ID       string            // 用户ID
	Username string            // 用户名
	Role     string            // 角色
	TenantID string            // 租户ID
	Claims   map[string]string // 其他声明
}

// AuthProvider 认证提供者接口
type AuthProvider interface {
	// Authenticate 从请求中提取凭据并验证
	Authenticate(c *gin.Context) (*AuthUser, error)
	
	// GenerateCredentials 生成新的凭据
	GenerateCredentials(user *AuthUser) (map[string]interface{}, error)
	
	// Name 获取提供者名称
	Name() string
}

// 通用错误
var (
	ErrAuthFailed              = errors.New("认证失败")
	ErrCredentialsNotFound     = errors.New("未找到认证凭据")
	ErrInvalidCredentials      = errors.New("无效的认证凭据")
	ErrTokenExpired            = errors.New("令牌已过期")
	ErrGenerateCredentialsFailed = errors.New("生成凭据失败")
)

// ===== JWT认证提供者 =====

// JWTProvider JWT认证提供者
type JWTProvider struct {
	config config.SecurityConfig
}

// NewJWTProvider 创建新的JWT认证提供者
func NewJWTProvider(conf config.SecurityConfig) *JWTProvider {
	return &JWTProvider{
		config: conf,
	}
}

// Name 获取提供者名称
func (p *JWTProvider) Name() string {
	return AuthTypeJWT
}

// Authenticate 从请求中提取JWT令牌并验证
func (p *JWTProvider) Authenticate(c *gin.Context) (*AuthUser, error) {
	// 从请求头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, ErrCredentialsNotFound
	}

	// 提取令牌
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, ErrInvalidCredentials
	}
	tokenString := parts[1]

	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidCredentials
		}
		return []byte(p.config.JwtSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrAuthFailed
	}

	// 验证令牌
	if !token.Valid {
		return nil, ErrInvalidCredentials
	}

	// 提取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	// 创建用户信息
	user := &AuthUser{
		Claims: make(map[string]string),
	}

	// 提取关键字段
	if id, ok := claims["user_id"].(string); ok {
		user.ID = id
	}
	if username, ok := claims["username"].(string); ok {
		user.Username = username
	}
	if role, ok := claims["role"].(string); ok {
		user.Role = role
	}
	if tenantID, ok := claims["tenant_id"].(string); ok {
		user.TenantID = tenantID
	}

	// 提取其他声明
	for key, val := range claims {
		if strVal, ok := val.(string); ok {
			if key != "user_id" && key != "username" && key != "role" && key != "tenant_id" {
				user.Claims[key] = strVal
			}
		}
	}

	return user, nil
}

// GenerateCredentials 生成JWT令牌
func (p *JWTProvider) GenerateCredentials(user *AuthUser) (map[string]interface{}, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(p.config.JwtExpiration) * time.Minute)

	// 创建声明
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"role":      user.Role,
		"tenant_id": user.TenantID,
		"exp":       expirationTime.Unix(),
		"iat":       time.Now().Unix(),
	}

	// 添加其他声明
	for key, val := range user.Claims {
		claims[key] = val
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(p.config.JwtSecret))
	if err != nil {
		return nil, ErrGenerateCredentialsFailed
	}

	// 返回凭据
	return map[string]interface{}{
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": p.config.JwtExpiration * 60, // 转换为秒
	}, nil
}

// ===== API Key认证提供者 =====

// APIKeyProvider API Key认证提供者
type APIKeyProvider struct {
	config config.SecurityConfig
	// 在真实实现中，应当有API Key的存储和管理机制
}

// NewAPIKeyProvider 创建新的API Key认证提供者
func NewAPIKeyProvider(conf config.SecurityConfig) *APIKeyProvider {
	return &APIKeyProvider{
		config: conf,
	}
}

// Name 获取提供者名称
func (p *APIKeyProvider) Name() string {
	return AuthTypeAPIKey
}

// Authenticate 从请求中提取API Key并验证
func (p *APIKeyProvider) Authenticate(c *gin.Context) (*AuthUser, error) {
	// 从请求头获取API Key
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		// 尝试从查询参数获取
		apiKey = c.Query("api_key")
		if apiKey == "" {
			return nil, ErrCredentialsNotFound
		}
	}

	// 在实际实现中，这里应该查询数据库验证API Key
	// 并获取关联的用户信息
	// 简单示例，假设有一个有效的API Key
	if apiKey == "test-api-key" {
		return &AuthUser{
			ID:       "api-user-1",
			Username: "api-client",
			Role:     "api",
			TenantID: "default",
			Claims: map[string]string{
				"scope": "read",
			},
		}, nil
	}

	return nil, ErrInvalidCredentials
}

// GenerateCredentials 生成API Key
func (p *APIKeyProvider) GenerateCredentials(user *AuthUser) (map[string]interface{}, error) {
	// 在实际实现中，这里应该生成新的API Key并保存到数据库
	// 简单示例，返回固定的API Key
	return map[string]interface{}{
		"api_key": "test-api-key",
	}, nil
}

// ===== OAuth2认证提供者 =====

// OAuth2Provider OAuth2认证提供者
type OAuth2Provider struct {
	config config.OAuth2Config
}

// NewOAuth2Provider 创建新的OAuth2认证提供者
func NewOAuth2Provider(conf config.OAuth2Config) *OAuth2Provider {
	return &OAuth2Provider{
		config: conf,
	}
}

// Name 获取提供者名称
func (p *OAuth2Provider) Name() string {
	return AuthTypeOAuth2
}

// Authenticate 从请求中提取OAuth2令牌并验证
func (p *OAuth2Provider) Authenticate(c *gin.Context) (*AuthUser, error) {
	// 从请求头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, ErrCredentialsNotFound
	}

	// 提取令牌
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, ErrInvalidCredentials
	}
	tokenString := parts[1]

	// 在实际实现中，这里应该验证OAuth2令牌
	// 并从令牌或用户信息端点获取用户信息
	// 简单示例，假设令牌有效
	if tokenString != "" {
		return &AuthUser{
			ID:       "oauth-user-1",
			Username: "oauth-client",
			Role:     "user",
			TenantID: "default",
			Claims: map[string]string{
				"scope": "read write",
			},
		}, nil
	}

	return nil, ErrInvalidCredentials
}

// GenerateCredentials 获取OAuth2授权URL
func (p *OAuth2Provider) GenerateCredentials(user *AuthUser) (map[string]interface{}, error) {
	// 在实际实现中，这里应该生成授权URL或处理授权码流程
	// 简单示例，返回授权URL
	return map[string]interface{}{
		"authorization_url": p.config.AuthorizationURL,
		"client_id":         p.config.ClientID,
		"scope":             p.config.Scope,
		"redirect_uri":      p.config.RedirectURI,
	}, nil
}

// ===== 基本认证提供者 =====

// BasicAuthProvider 基本认证提供者
type BasicAuthProvider struct {
	// 在真实实现中，应当有用户凭据的存储和管理机制
}

// NewBasicAuthProvider 创建新的基本认证提供者
func NewBasicAuthProvider() *BasicAuthProvider {
	return &BasicAuthProvider{}
}

// Name 获取提供者名称
func (p *BasicAuthProvider) Name() string {
	return AuthTypeBasic
}

// Authenticate 从请求中提取基本认证凭据并验证
func (p *BasicAuthProvider) Authenticate(c *gin.Context) (*AuthUser, error) {
	// 获取基本认证信息
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return nil, ErrCredentialsNotFound
	}

	// 在实际实现中，这里应该查询数据库验证用户名和密码
	// 简单示例，假设有一个有效的用户
	if username == "admin" && password == "password" {
		return &AuthUser{
			ID:       "basic-user-1",
			Username: username,
			Role:     "admin",
			TenantID: "default",
			Claims:   map[string]string{},
		}, nil
	}

	return nil, ErrInvalidCredentials
}

// GenerateCredentials 基本认证不需要生成凭据
func (p *BasicAuthProvider) GenerateCredentials(user *AuthUser) (map[string]interface{}, error) {
	// 基本认证不需要生成凭据，返回用户信息
	return map[string]interface{}{
		"username": user.Username,
		"role":     user.Role,
	}, nil
}
