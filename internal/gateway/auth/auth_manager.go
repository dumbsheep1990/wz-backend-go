package auth

import (
	"errors"
	"sync"

	"wz-backend-go/internal/gateway/config"
)

// 认证类型常量
const (
	AuthTypeJWT     = "jwt"
	AuthTypeAPIKey  = "apikey"
	AuthTypeOAuth2  = "oauth2"
	AuthTypeBasic   = "basic"
	AuthTypeNone    = "none"
)

// 错误定义
var (
	ErrInvalidAuthType     = errors.New("无效的认证类型")
	ErrAuthProviderNotFound = errors.New("未找到认证提供者")
	ErrPermissionDenied    = errors.New("权限不足")
)

// AuthManager 认证管理器
type AuthManager struct {
	providers map[string]AuthProvider
	config    config.SecurityConfig
	mu        sync.RWMutex
}

// NewAuthManager 创建新的认证管理器
func NewAuthManager(conf config.SecurityConfig) *AuthManager {
	manager := &AuthManager{
		providers: make(map[string]AuthProvider),
		config:    conf,
	}

	// 注册默认的认证提供者
	manager.RegisterProvider(AuthTypeJWT, NewJWTProvider(conf))
	manager.RegisterProvider(AuthTypeAPIKey, NewAPIKeyProvider(conf))
	
	// 如果启用了OAuth2，注册OAuth2提供者
	if conf.OAuth2.Enabled {
		manager.RegisterProvider(AuthTypeOAuth2, NewOAuth2Provider(conf.OAuth2))
	}
	
	// 注册基本认证提供者
	manager.RegisterProvider(AuthTypeBasic, NewBasicAuthProvider())

	return manager
}

// RegisterProvider 注册认证提供者
func (am *AuthManager) RegisterProvider(authType string, provider AuthProvider) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	am.providers[authType] = provider
}

// GetProvider 获取认证提供者
func (am *AuthManager) GetProvider(authType string) (AuthProvider, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	provider, exists := am.providers[authType]
	if !exists {
		return nil, ErrAuthProviderNotFound
	}
	
	return provider, nil
}

// IsAuthorized 检查是否有权限
func (am *AuthManager) IsAuthorized(userRole string, requiredRoles []string) bool {
	// 如果未设置必要角色，则允许访问
	if len(requiredRoles) == 0 {
		return true
	}
	
	// 检查用户角色是否在必要角色列表中
	for _, role := range requiredRoles {
		if role == userRole {
			return true
		}
	}
	
	return false
}

// HasPermission 检查是否拥有权限
func (am *AuthManager) HasPermission(userID, resource, action string) bool {
	// 这里应当根据权限系统检查用户是否有对资源执行操作的权限
	// 简单实现，后续可以扩展为更完整的RBAC实现
	return true
}
