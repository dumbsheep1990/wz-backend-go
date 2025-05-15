package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/token"
)

// AdminAuthMiddleware 后台管理认证中间件
type AdminAuthMiddleware struct {
	enforcer *casbin.Enforcer // 权限管理器
	redis    *redis.Redis     // Redis缓存
}

// NewAdminAuthMiddleware 创建新的管理后台认证中间件
func NewAdminAuthMiddleware(enforcer *casbin.Enforcer, redis *redis.Redis) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		enforcer: enforcer,
		redis:    redis,
	}
}

// Handle 权限验证处理函数
func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求头中的Authorization
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "未提供认证信息", http.StatusUnauthorized)
			return
		}

		// 提取JWT令牌
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "认证格式错误", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		// 解析JWT令牌
		claims, err := token.ParseToken(tokenStr, []byte("your-jwt-secret-here")) // 实际中应从配置获取密钥
		if err != nil {
			logx.Errorf("解析JWT失败: %v", err)
			http.Error(w, "认证失败", http.StatusUnauthorized)
			return
		}

		// 提取用户ID和角色
		var userClaims map[string]interface{}
		jsonClaims, _ := json.Marshal(claims)
		json.Unmarshal(jsonClaims, &userClaims)

		userId, ok := userClaims["id"].(float64)
		if !ok {
			http.Error(w, "无效的用户ID", http.StatusForbidden)
			return
		}

		role, ok := userClaims["role"].(string)
		if !ok || role == "" {
			http.Error(w, "无效的用户角色", http.StatusForbidden)
			return
		}

		// 检查用户状态（可以从Redis缓存中获取）
		userKey := "admin:user:" + string(int64(userId))
		status, err := m.redis.Hget(userKey, "status")
		if err != nil || status != "1" {
			http.Error(w, "用户已被禁用或不存在", http.StatusForbidden)
			return
		}

		// 检查资源权限
		path := r.URL.Path
		method := r.Method
		allowed, err := m.enforcer.Enforce(role, path, method)
		if err != nil {
			logx.Errorf("权限检查错误: %v", err)
			http.Error(w, "权限检查失败", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "无权访问该资源", http.StatusForbidden)
			return
		}

		// 将用户信息添加到请求上下文中
		r.Header.Set("X-User-Id", string(int64(userId)))
		r.Header.Set("X-User-Role", role)

		// 记录API访问日志（可以在这里记录，或者在服务层面统一处理）

		// 调用下一个处理器
		next(w, r)
	}
}
