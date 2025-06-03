package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/go-redis/redis/v8"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// AuthHandlerImpl 认证处理器实现
type AuthHandlerImpl struct {
	redisClient *redis.Client
	jwtSecret   []byte
	jwtIssuer   string
	tokenTTL    time.Duration
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(redisClient *redis.Client, jwtSecret string, jwtIssuer string, tokenTTL time.Duration) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		redisClient: redisClient,
		jwtSecret:   []byte(jwtSecret),
		jwtIssuer:   jwtIssuer,
		tokenTTL:    tokenTTL,
	}
}

// VerifyToken 验证访问令牌
func (a *AuthHandlerImpl) VerifyToken(ctx context.Context, tokenString string) (*entity.UserSession, error) {
	// 解析JWT令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}
	
	// 验证令牌声明
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 检查令牌是否已被撤销
		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			return nil, errors.New("令牌缺少JTI字段")
		}
		
		// 检查Redis中的撤销记录
		revokedKey := fmt.Sprintf("token:revoked:%s", jti)
		revoked, err := a.redisClient.Exists(ctx, revokedKey).Result()
		if err != nil {
			log.Printf("检查令牌撤销状态失败: %v", err)
		}
		
		if revoked > 0 {
			return nil, errors.New("令牌已被撤销")
		}
		
		// 获取用户ID
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			return nil, errors.New("令牌缺少用户ID")
		}
		
		// 获取用户名
		username, _ := claims["username"].(string)
		
		// 获取角色
		roleStr, _ := claims["role"].(string)
		role := valueobject.UserRole(roleStr)
		
		// 获取权限
		var permissions []string
		if permInterface, ok := claims["permissions"].([]interface{}); ok {
			for _, p := range permInterface {
				if pStr, ok := p.(string); ok {
					permissions = append(permissions, pStr)
				}
			}
		}
		
		// 创建用户会话
		session := &entity.UserSession{
			UserID:      userID,
			Username:    username,
			Role:        role,
			Permissions: permissions,
			TokenID:     jti,
			ExpiresAt:   time.Unix(int64(claims["exp"].(float64)), 0),
		}
		
		return session, nil
	}
	
	return nil, errors.New("无效的令牌")
}

// ExtractToken 从HTTP请求中提取令牌
func (a *AuthHandlerImpl) ExtractToken(r *http.Request) (string, error) {
	// 从Authorization头中提取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Bearer {token}
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
	}
	
	// 从查询参数中提取
	token := r.URL.Query().Get("access_token")
	if token != "" {
		return token, nil
	}
	
	// 从Cookie中提取
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	
	return "", errors.New("找不到访问令牌")
}

// RevokeToken 撤销令牌
func (a *AuthHandlerImpl) RevokeToken(ctx context.Context, tokenID string) error {
	// 计算令牌过期时间
	revokedKey := fmt.Sprintf("token:revoked:%s", tokenID)
	
	// 将令牌ID添加到撤销列表，使用原令牌相同的过期时间
	_, err := a.redisClient.Set(ctx, revokedKey, "1", a.tokenTTL).Result()
	if err != nil {
		return fmt.Errorf("撤销令牌失败: %w", err)
	}
	
	return nil
}

// GenerateToken 生成新的访问令牌
func (a *AuthHandlerImpl) GenerateToken(ctx context.Context, userID string, username string, role valueobject.UserRole, permissions []string) (string, error) {
	// 生成唯一的令牌ID
	tokenID := fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())
	
	// 创建JWT声明
	now := time.Now()
	expiresAt := now.Add(a.tokenTTL)
	
	claims := jwt.MapClaims{
		"iss":         a.jwtIssuer,
		"sub":         userID,
		"jti":         tokenID,
		"iat":         now.Unix(),
		"exp":         expiresAt.Unix(),
		"username":    username,
		"role":        string(role),
		"permissions": permissions,
	}
	
	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名令牌
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("签名令牌失败: %w", err)
	}
	
	// 存储令牌信息到Redis（可选，用于令牌管理）
	sessionKey := fmt.Sprintf("token:session:%s", tokenID)
	sessionData := map[string]interface{}{
		"user_id":     userID,
		"username":    username,
		"role":        string(role),
		"permissions": strings.Join(permissions, ","),
		"expires_at":  expiresAt.Unix(),
	}
	
	_, err = a.redisClient.HMSet(ctx, sessionKey, sessionData).Result()
	if err != nil {
		log.Printf("存储令牌会话数据失败: %v", err)
	}
	
	// 设置过期时间
	a.redisClient.Expire(ctx, sessionKey, a.tokenTTL)
	
	return tokenString, nil
}

// CheckPermission 检查用户是否有特定权限
func (a *AuthHandlerImpl) CheckPermission(session *entity.UserSession, requiredPermission string) bool {
	// 超级管理员拥有所有权限
	if session.Role == valueobject.RoleAdmin {
		return true
	}
	
	// 检查用户具体权限
	for _, permission := range session.Permissions {
		// 精确匹配
		if permission == requiredPermission {
			return true
		}
		
		// 通配符匹配（如 "users.*" 匹配 "users.read"）
		if strings.HasSuffix(permission, ".*") {
			prefix := strings.TrimSuffix(permission, ".*")
			if strings.HasPrefix(requiredPermission, prefix) {
				return true
			}
		}
	}
	
	return false
}

// CheckRole 检查用户是否有特定角色
func (a *AuthHandlerImpl) CheckRole(session *entity.UserSession, requiredRole valueobject.UserRole) bool {
	// 超级管理员可以执行任何角色的操作
	if session.Role == valueobject.RoleAdmin {
		return true
	}
	
	// 检查用户角色
	return session.Role == requiredRole
}
