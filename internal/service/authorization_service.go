package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
)

// AuthorizationService 授权服务接口
type AuthorizationService interface {
	// 检查用户对资源的访问权限
	CheckPermission(ctx context.Context, sub string, obj string, act string) (bool, error)
	// 为用户添加角色
	AddRoleForUser(ctx context.Context, user string, role string) (bool, error)
	// 为角色添加权限
	AddPermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error)
	// 删除用户角色
	DeleteRoleForUser(ctx context.Context, user string, role string) (bool, error)
	// 删除角色权限
	DeletePermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error)
	// 获取用户角色
	GetRolesForUser(ctx context.Context, user string) ([]string, error)
	// 获取角色权限
	GetPermissionsForRole(ctx context.Context, role string) ([][]string, error)
	// 强制刷新策略
	RefreshPolicy(ctx context.Context) error
}

type authorizationService struct {
	enforcer *casbin.Enforcer
	redis    *redis.Client
}

// NewAuthorizationService 创建授权服务
func NewAuthorizationService(enforcer *casbin.Enforcer, redis *redis.Client) AuthorizationService {
	return &authorizationService{
		enforcer: enforcer,
		redis:    redis,
	}
}

// CheckPermission 检查用户对资源的访问权限
func (s *authorizationService) CheckPermission(ctx context.Context, sub string, obj string, act string) (bool, error) {
	if s.enforcer == nil {
		return false, errors.New("enforcer not initialized")
	}

	// 使用缓存优化权限检查
	cacheKey := fmt.Sprintf("perm:%s:%s:%s", sub, obj, act)
	
	// 尝试从缓存获取结果
	result, err := s.redis.Get(ctx, cacheKey).Int()
	if err == nil {
		// 缓存命中
		return result == 1, nil
	}

	// 缓存未命中，调用Casbin检查权限
	allowed, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, err
	}

	// 缓存结果（5分钟过期）
	value := 0
	if allowed {
		value = 1
	}
	s.redis.Set(ctx, cacheKey, value, 5*60*1000*1000*1000) // 5分钟

	return allowed, nil
}

// AddRoleForUser 为用户添加角色
func (s *authorizationService) AddRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	if s.enforcer == nil {
		return false, errors.New("enforcer not initialized")
	}

	// 添加角色
	success, err := s.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return false, err
	}

	// 保存更改
	err = s.enforcer.SavePolicy()
	if err != nil {
		return false, err
	}

	// 清理受影响的缓存
	s.clearPermissionCache(ctx, user)

	return success, nil
}

// AddPermissionForRole 为角色添加权限
func (s *authorizationService) AddPermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error) {
	if s.enforcer == nil {
		return false, errors.New("enforcer not initialized")
	}

	// 添加权限
	success, err := s.enforcer.AddPolicy(role, obj, act)
	if err != nil {
		return false, err
	}

	// 保存更改
	err = s.enforcer.SavePolicy()
	if err != nil {
		return false, err
	}

	// 清理受影响的缓存
	s.clearPermissionCacheForRole(ctx, role)

	return success, nil
}

// DeleteRoleForUser 删除用户角色
func (s *authorizationService) DeleteRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	if s.enforcer == nil {
		return false, errors.New("enforcer not initialized")
	}

	// 删除角色
	success, err := s.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return false, err
	}

	// 保存更改
	err = s.enforcer.SavePolicy()
	if err != nil {
		return false, err
	}

	// 清理受影响的缓存
	s.clearPermissionCache(ctx, user)

	return success, nil
}

// DeletePermissionForRole 删除角色权限
func (s *authorizationService) DeletePermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error) {
	if s.enforcer == nil {
		return false, errors.New("enforcer not initialized")
	}

	// 删除权限
	success, err := s.enforcer.RemovePolicy(role, obj, act)
	if err != nil {
		return false, err
	}

	// 保存更改
	err = s.enforcer.SavePolicy()
	if err != nil {
		return false, err
	}

	// 清理受影响的缓存
	s.clearPermissionCacheForRole(ctx, role)

	return success, nil
}

// GetRolesForUser 获取用户角色
func (s *authorizationService) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	if s.enforcer == nil {
		return nil, errors.New("enforcer not initialized")
	}

	// 缓存键
	cacheKey := fmt.Sprintf("roles:%s", user)
	
	// 尝试从缓存获取
	rolesStr, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var roles []string
		err = json.Unmarshal([]byte(rolesStr), &roles)
		if err == nil {
			return roles, nil
		}
	}

	// 缓存未命中或解析失败，从Casbin获取
	roles, err := s.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, err
	}

	// 缓存结果（5分钟过期）
	rolesBytes, _ := json.Marshal(roles)
	s.redis.Set(ctx, cacheKey, string(rolesBytes), 5*60*1000*1000*1000) // 5分钟

	return roles, nil
}

// GetPermissionsForRole 获取角色权限
func (s *authorizationService) GetPermissionsForRole(ctx context.Context, role string) ([][]string, error) {
	if s.enforcer == nil {
		return nil, errors.New("enforcer not initialized")
	}

	// 缓存键
	cacheKey := fmt.Sprintf("perms:%s", role)
	
	// 尝试从缓存获取
	permsStr, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var perms [][]string
		err = json.Unmarshal([]byte(permsStr), &perms)
		if err == nil {
			return perms, nil
		}
	}

	// 缓存未命中或解析失败，从Casbin获取
	perms, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}

	// 缓存结果（5分钟过期）
	permsBytes, _ := json.Marshal(perms)
	s.redis.Set(ctx, cacheKey, string(permsBytes), 5*60*1000*1000*1000) // 5分钟

	return perms, nil
}

// RefreshPolicy 强制刷新策略
func (s *authorizationService) RefreshPolicy(ctx context.Context) error {
	if s.enforcer == nil {
		return errors.New("enforcer not initialized")
	}

	// 重新加载策略
	err := s.enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	// 清理所有权限相关缓存
	// 使用通配符删除所有权限相关缓存
	keys, _, _ := s.redis.Scan(ctx, 0, "perm:*", 1000).Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}
	
	// 清理角色缓存
	keys, _, _ = s.redis.Scan(ctx, 0, "roles:*", 1000).Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}
	
	// 清理权限缓存
	keys, _, _ = s.redis.Scan(ctx, 0, "perms:*", 1000).Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}

	return nil
}

// 清理用户的权限缓存
func (s *authorizationService) clearPermissionCache(ctx context.Context, user string) {
	// 删除用户角色缓存
	s.redis.Del(ctx, fmt.Sprintf("roles:%s", user))
	
	// 删除用户权限缓存
	keys, _, _ := s.redis.Scan(ctx, 0, fmt.Sprintf("perm:%s:*", user), 1000).Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}
}

// 清理角色的权限缓存
func (s *authorizationService) clearPermissionCacheForRole(ctx context.Context, role string) {
	// 删除角色权限缓存
	s.redis.Del(ctx, fmt.Sprintf("perms:%s", role))
	
	// 获取所有具有该角色的用户
	users, err := s.enforcer.GetUsersForRole(role)
	if err != nil {
		// 记录错误但继续执行，因为这只是缓存清理操作
		log.Printf("获取角色用户列表失败: %v", err)
		return
	}
	
	// 删除这些用户的权限缓存
	for _, user := range users {
		s.clearPermissionCache(ctx, user)
	}
}
