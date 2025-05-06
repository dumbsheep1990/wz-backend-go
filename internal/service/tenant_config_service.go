package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wxnacy/wz-backend-go/internal/pkg/tenantctx"
)

// TenantConfig 表示租户的配置信息
type TenantConfig struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Status          int                     `json:"status"`  // 1: 正常, 2: 禁用
	ExpirationDate  time.Time               `json:"expiration_date"`
	Features        map[string]bool         `json:"features"`
	PlatformConfigs map[string]PlatformConfig `json:"platform_configs"`
}

// PlatformConfig 包含平台特定的配置信息
type PlatformConfig struct {
	Enabled         bool              `json:"enabled"`
	Theme           string            `json:"theme"`
	AppVersion      string            `json:"app_version"`
	MinVersion      string            `json:"min_version"`
	MaxVersion      string            `json:"max_version"`
	CustomSettings  map[string]string `json:"custom_settings"`
	APIEndpoints    map[string]string `json:"api_endpoints"`
}

// TenantConfigService 处理租户配置操作
type TenantConfigService struct {
	tenantctx.BaseTenantService
	redisClient *redis.Client
	configCache map[string]*TenantConfig
	cacheMutex  sync.RWMutex
	cacheExp    time.Duration
}

// NewTenantConfigService 创建一个新的租户配置服务
func NewTenantConfigService(ctx context.Context, redisClient *redis.Client) *TenantConfigService {
	return &TenantConfigService{
		BaseTenantService: *tenantctx.NewBaseTenantService(ctx),
		redisClient:       redisClient,
		configCache:       make(map[string]*TenantConfig),
		cacheExp:          time.Minute * 5, // 缓存租户配置5分钟
	}
}

// WithContext 实现TenantAwareService.WithContext接口
func (s *TenantConfigService) WithContext(ctx context.Context) *TenantConfigService {
	return &TenantConfigService{
		BaseTenantService: *s.BaseTenantService.WithContext(ctx),
		redisClient:       s.redisClient,
		configCache:       s.configCache,
		cacheMutex:        s.cacheMutex,
		cacheExp:          s.cacheExp,
	}
}

// GetTenantConfig 获取当前租户的配置信息
func (s *TenantConfigService) GetTenantConfig() (*TenantConfig, error) {
	tenantID, err := s.GetCurrentTenant()
	if err != nil {
		return nil, err
	}
	
	return s.GetTenantConfigByID(tenantID)
}

// GetTenantConfigByID 获取指定租户的配置信息
func (s *TenantConfigService) GetTenantConfigByID(tenantID string) (*TenantConfig, error) {
	// 首先检查缓存
	s.cacheMutex.RLock()
	if config, exists := s.configCache[tenantID]; exists {
		s.cacheMutex.RUnlock()
		return config, nil
	}
	s.cacheMutex.RUnlock()
	
	// 从Redis获取
	key := fmt.Sprintf("tenant:config:%s", tenantID)
	configJSON, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis错误: %w", err)
	}
	
	var config *TenantConfig
	
	if err == redis.Nil || configJSON == "" {
		// 需要从数据库加载
		config, err = s.loadTenantConfigFromDB(tenantID)
		if err != nil {
			return nil, err
		}
		
		// 保存到Redis
		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("序列化错误: %w", err)
		}
		
		// 设置过期时间
		err = s.redisClient.Set(context.Background(), key, configJSON, time.Hour).Err()
		if err != nil {
			// 记录错误但继续执行，这不是关键错误
			fmt.Printf("缓存租户配置到redis失败: %v", err)
		}
	} else {
		// 从Redis解析
		config = &TenantConfig{}
		if err := json.Unmarshal([]byte(configJSON), config); err != nil {
			return nil, fmt.Errorf("反序列化错误: %w", err)
		}
	}
	
	// 更新缓存
	s.cacheMutex.Lock()
	s.configCache[tenantID] = config
	s.cacheMutex.Unlock()
	
	return config, nil
}

// loadTenantConfigFromDB 从数据库加载租户配置
// 在实际实现中，这将查询数据库
func (s *TenantConfigService) loadTenantConfigFromDB(tenantID string) (*TenantConfig, error) {
	// 这是一个占位符。在实际实现中，你需要查询数据库。
	// 为了演示目的，我们创建一个虚拟配置。
	config := &TenantConfig{
		ID:     tenantID,
		Name:   "租户 " + tenantID,
		Status: 1, // 正常
		ExpirationDate: time.Now().AddDate(1, 0, 0), // 从现在起一年后
		Features: map[string]bool{
			"content_management": true,
			"e_commerce":        true,
			"analytics":         true,
		},
		PlatformConfigs: map[string]PlatformConfig{
			string(tenantctx.PlatformWeb): {
				Enabled:    true,
				Theme:      "default",
				AppVersion: "1.0.0",
				CustomSettings: map[string]string{
					"max_upload_size": "10MB",
					"enable_chat":     "true",
				},
				APIEndpoints: map[string]string{
					"content": "/api/v1/content",
					"users":   "/api/v1/users",
				},
			},
			string(tenantctx.PlatformMobile): {
				Enabled:    true,
				Theme:      "mobile_dark",
				AppVersion: "1.0.0",
				MinVersion: "1.0.0",
				MaxVersion: "2.0.0",
				CustomSettings: map[string]string{
					"max_upload_size": "5MB",
					"enable_push":     "true",
				},
				APIEndpoints: map[string]string{
					"content": "/api/v1/content",
					"users":   "/api/v1/users",
					"offline": "/api/v1/sync",
				},
			},
			string(tenantctx.PlatformUniApp): {
				Enabled:    true,
				Theme:      "uniapp_light",
				AppVersion: "1.0.0",
				MinVersion: "1.0.0",
				CustomSettings: map[string]string{
					"max_upload_size": "8MB",
				},
				APIEndpoints: map[string]string{
					"content": "/api/v1/content",
					"users":   "/api/v1/users",
				},
			},
		},
	}
	
	return config, nil
}

// GetPlatformConfig 获取当前租户的平台特定配置
func (s *TenantConfigService) GetPlatformConfig() (*PlatformConfig, error) {
	// 获取租户配置
	config, err := s.GetTenantConfig()
	if err != nil {
		return nil, err
	}
	
	// 获取平台
	platform, err := s.GetCurrentPlatform()
	if err != nil {
		return nil, err
	}
	
	// 获取平台配置
	platformConfig, exists := config.PlatformConfigs[string(platform)]
	if !exists {
		// 降级到Web配置
		platformConfig, exists = config.PlatformConfigs[string(tenantctx.PlatformWeb)]
		if !exists {
			return nil, fmt.Errorf("未找到平台 %s 的配置", platform)
		}
	}
	
	return &platformConfig, nil
}

// IsTenantEnabled 检查租户是否启用且未过期
func (s *TenantConfigService) IsTenantEnabled(tenantID string) (bool, error) {
	config, err := s.GetTenantConfigByID(tenantID)
	if err != nil {
		return false, err
	}
	
	if config.Status != 1 {
		return false, nil
	}
	
	if time.Now().After(config.ExpirationDate) {
		return false, nil
	}
	
	return true, nil
}

// IsPlatformEnabled 检查当前平台是否为当前租户启用
func (s *TenantConfigService) IsPlatformEnabled() (bool, error) {
	platformConfig, err := s.GetPlatformConfig()
	if err != nil {
		return false, err
	}
	
	return platformConfig.Enabled, nil
}

// FeatureEnabled 检查特定功能是否为当前租户启用
func (s *TenantConfigService) FeatureEnabled(featureName string) (bool, error) {
	config, err := s.GetTenantConfig()
	if err != nil {
		return false, err
	}
	
	enabled, exists := config.Features[featureName]
	if !exists {
		return false, nil
	}
	
	return enabled, nil
}

// RefreshTenantConfig 强制从数据库刷新租户配置
func (s *TenantConfigService) RefreshTenantConfig(tenantID string) error {
	// 从缓存中移除
	s.cacheMutex.Lock()
	delete(s.configCache, tenantID)
	s.cacheMutex.Unlock()
	
	// 从Redis中移除
	key := fmt.Sprintf("tenant:config:%s", tenantID)
	err := s.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("redis错误: %w", err)
	}
	
	return nil
}
