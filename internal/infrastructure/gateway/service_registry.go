package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
)

// ServiceRegistryImpl 服务注册表实现
type ServiceRegistryImpl struct {
	serviceRepo    repository.ServiceRepository
	services       map[string]*entity.Service
	healthStatuses map[string]bool
	httpClient     *http.Client
	mutex          sync.RWMutex
}

// NewServiceRegistry 创建新的服务注册表
func NewServiceRegistry(serviceRepo repository.ServiceRepository) *ServiceRegistryImpl {
	registry := &ServiceRegistryImpl{
		serviceRepo:    serviceRepo,
		services:       make(map[string]*entity.Service),
		healthStatuses: make(map[string]bool),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	
	// 初始加载所有服务
	ctx := context.Background()
	if err := registry.LoadServices(ctx); err != nil {
		log.Printf("初始服务加载失败: %v", err)
	}
	
	return registry
}

// LoadServices 从存储中加载所有服务
func (r *ServiceRegistryImpl) LoadServices(ctx context.Context) error {
	services, err := r.serviceRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("从存储加载服务失败: %w", err)
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// 清空当前服务表并重新加载
	r.services = make(map[string]*entity.Service)
	for _, service := range services {
		r.services[service.ID] = service
		// 初始化健康状态为未知（false）
		r.healthStatuses[service.ID] = false
	}
	
	log.Printf("已加载 %d 个服务", len(r.services))
	return nil
}

// GetService 获取服务信息
func (r *ServiceRegistryImpl) GetService(serviceID string) (*entity.Service, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	service, exists := r.services[serviceID]
	return service, exists
}

// GetServiceByName 通过名称获取服务
func (r *ServiceRegistryImpl) GetServiceByName(name string) (*entity.Service, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, service := range r.services {
		if service.Name == name {
			return service, true
		}
	}
	
	return nil, false
}

// AddService 添加服务
func (r *ServiceRegistryImpl) AddService(service *entity.Service) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.services[service.ID] = service
	r.healthStatuses[service.ID] = false
}

// RemoveService 移除服务
func (r *ServiceRegistryImpl) RemoveService(serviceID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.services, serviceID)
	delete(r.healthStatuses, serviceID)
}

// GetAllServices 获取所有服务
func (r *ServiceRegistryImpl) GetAllServices() []*entity.Service {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	services := make([]*entity.Service, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	return services
}

// CheckHealth 检查服务健康状态
func (r *ServiceRegistryImpl) CheckHealth(ctx context.Context, serviceID string) (bool, string, error) {
	r.mutex.RLock()
	service, exists := r.services[serviceID]
	r.mutex.RUnlock()
	
	if !exists {
		return false, "服务不存在", fmt.Errorf("服务ID %s 不存在", serviceID)
	}
	
	// 构建健康检查URL
	healthURL := fmt.Sprintf("%s/health", service.BaseURL)
	if service.HealthCheckPath != "" {
		healthURL = fmt.Sprintf("%s%s", service.BaseURL, service.HealthCheckPath)
	}
	
	// 创建带上下文的请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	if err != nil {
		return false, "创建健康检查请求失败", err
	}
	
	// 执行健康检查请求
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.updateHealthStatus(serviceID, false)
		return false, fmt.Sprintf("健康检查请求失败: %v", err), err
	}
	defer resp.Body.Close()
	
	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		r.updateHealthStatus(serviceID, false)
		return false, fmt.Sprintf("健康检查返回非成功状态码: %d", resp.StatusCode), nil
	}
	
	// 更新健康状态
	r.updateHealthStatus(serviceID, true)
	return true, "服务健康", nil
}

// IsHealthy 检查服务是否健康
func (r *ServiceRegistryImpl) IsHealthy(serviceID string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	healthy, exists := r.healthStatuses[serviceID]
	return exists && healthy
}

// 更新服务健康状态
func (r *ServiceRegistryImpl) updateHealthStatus(serviceID string, healthy bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.healthStatuses[serviceID] = healthy
}
