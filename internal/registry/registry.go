package registry

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	// ErrInvalidServiceConfig 表示无效的服务配置错误
	ErrInvalidServiceConfig = errors.New("无效的服务配置")

	// ErrNacosClientNotInitialized 表示Nacos客户端未初始化错误
	ErrNacosClientNotInitialized = errors.New("Nacos客户端未初始化")
)

// NacosConfig 表示Nacos客户端的配置
type NacosConfig struct {
	ServerAddr      string `yaml:"server_addr"`      // Nacos服务器地址
	ServerPort      uint64 `yaml:"server_port"`      // Nacos服务器端口
	Namespace       string `yaml:"namespace"`        // 命名空间ID
	Group           string `yaml:"group"`            // 分组
	DataID          string `yaml:"data_id"`          // 配置ID
	LogDir          string `yaml:"log_dir"`          // 日志目录
	CacheDir        string `yaml:"cache_dir"`        // 缓存目录
	LogLevel        string `yaml:"log_level"`        // 日志级别
	Username        string `yaml:"username"`         // 用户名(可选)
	Password        string `yaml:"password"`         // 密码(可选)
	HealthCheckPort int    `yaml:"health_check_port"` // 健康检查端口
}

// ServiceRegistry 定义服务注册接口
type ServiceRegistry interface {
	// Register 注册服务
	Register(serviceName, ip string, port int, meta map[string]string) error

	// Deregister 注销服务
	Deregister(serviceName, ip string, port int) error

	// GetService 获取服务实例
	GetService(serviceName string) ([]ServiceInstance, error)

	// Subscribe 订阅服务变更
	Subscribe(serviceName string, callback func(instances []ServiceInstance)) error

	// Unsubscribe 取消订阅服务变更
	Unsubscribe(serviceName string) error

	// Close 关闭注册中心客户端
	Close() error
}

// ServiceInstance 定义服务实例
type ServiceInstance struct {
	ServiceName string            // 服务名称
	IP          string            // 实例IP地址
	Port        int               // 实例端口
	Metadata    map[string]string // 元数据
	Healthy     bool              // 是否健康
}

// NacosRegistry 基于Nacos的服务注册实现
type NacosRegistry struct {
	client            naming_client.INamingClient
	config            *NacosConfig
	serviceInstances  map[string][]ServiceInstance
	subscriptions     map[string]func(instances []ServiceInstance)
	instancesLock     sync.RWMutex
	subscriptionsLock sync.RWMutex
}

// NewNacosRegistry 创建新的Nacos服务注册中心
func NewNacosRegistry(config *NacosConfig) (*NacosRegistry, error) {
	if config == nil || config.ServerAddr == "" || config.ServerPort == 0 {
		return nil, ErrInvalidServiceConfig
	}

	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			config.ServerAddr,
			config.ServerPort,
			constant.WithContextPath("/nacos"),
		),
	}

	// 创建ClientConfig
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(config.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(config.LogDir),
		constant.WithCacheDir(config.CacheDir),
		constant.WithLogLevel(config.LogLevel),
	)

	// 如果提供了用户名和密码，则设置它们
	if config.Username != "" && config.Password != "" {
		clientConfig.Username = config.Username
		clientConfig.Password = config.Password
	}

	// 创建服务发现客户端
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("创建Nacos客户端失败: %w", err)
	}

	return &NacosRegistry{
		client:           client,
		config:           config,
		serviceInstances: make(map[string][]ServiceInstance),
		subscriptions:    make(map[string]func(instances []ServiceInstance)),
	}, nil
}

// Register 向Nacos注册服务
func (r *NacosRegistry) Register(serviceName, ip string, port int, meta map[string]string) error {
	if r.client == nil {
		return ErrNacosClientNotInitialized
	}

	// 默认元数据
	if meta == nil {
		meta = make(map[string]string)
	}

	// 用于健康检查的端口
	healthCheckPort := port
	if r.config.HealthCheckPort > 0 {
		healthCheckPort = r.config.HealthCheckPort
	}

	// 注册服务实例
	_, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: serviceName,
		Weight:      10.0,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    meta,
		GroupName:   r.config.Group,
		ClusterName: "DEFAULT",
	})

	return err
}

// Deregister 从Nacos注销服务
func (r *NacosRegistry) Deregister(serviceName, ip string, port int) error {
	if r.client == nil {
		return ErrNacosClientNotInitialized
	}

	_, err := r.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: serviceName,
		GroupName:   r.config.Group,
		Ephemeral:   true,
	})

	return err
}

// GetService 获取服务实例列表
func (r *NacosRegistry) GetService(serviceName string) ([]ServiceInstance, error) {
	if r.client == nil {
		return nil, ErrNacosClientNotInitialized
	}

	// 先检查缓存
	r.instancesLock.RLock()
	instances, ok := r.serviceInstances[serviceName]
	r.instancesLock.RUnlock()

	if ok {
		return instances, nil
	}

	// 从Nacos获取服务实例
	resp, err := r.client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   r.config.Group,
	})
	if err != nil {
		return nil, err
	}

	// 转换为内部ServiceInstance结构
	var result []ServiceInstance
	for _, host := range resp.Hosts {
		result = append(result, ServiceInstance{
			ServiceName: serviceName,
			IP:          host.Ip,
			Port:        int(host.Port),
			Metadata:    host.Metadata,
			Healthy:     host.Healthy,
		})
	}

	// 更新缓存
	r.instancesLock.Lock()
	r.serviceInstances[serviceName] = result
	r.instancesLock.Unlock()

	return result, nil
}

// Subscribe 订阅服务变更
func (r *NacosRegistry) Subscribe(serviceName string, callback func(instances []ServiceInstance)) error {
	if r.client == nil {
		return ErrNacosClientNotInitialized
	}

	// 保存回调函数
	r.subscriptionsLock.Lock()
	r.subscriptions[serviceName] = callback
	r.subscriptionsLock.Unlock()

	// 订阅服务变更
	err := r.client.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   r.config.Group,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				fmt.Printf("服务变更通知失败: %v\n", err)
				return
			}

			// 转换为内部ServiceInstance结构
			var instances []ServiceInstance
			for _, service := range services {
				instances = append(instances, ServiceInstance{
					ServiceName: serviceName,
					IP:          service.Ip,
					Port:        int(service.Port),
					Metadata:    service.Metadata,
					Healthy:     service.Healthy,
				})
			}

			// 更新缓存
			r.instancesLock.Lock()
			r.serviceInstances[serviceName] = instances
			r.instancesLock.Unlock()

			// 执行回调
			r.subscriptionsLock.RLock()
			cb, exists := r.subscriptions[serviceName]
			r.subscriptionsLock.RUnlock()

			if exists && cb != nil {
				cb(instances)
			}
		},
	})

	return err
}

// Unsubscribe 取消订阅服务变更
func (r *NacosRegistry) Unsubscribe(serviceName string) error {
	if r.client == nil {
		return ErrNacosClientNotInitialized
	}

	// 从Nacos取消订阅
	err := r.client.Unsubscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   r.config.Group,
	})

	if err != nil {
		return err
	}

	// 移除回调函数
	r.subscriptionsLock.Lock()
	delete(r.subscriptions, serviceName)
	r.subscriptionsLock.Unlock()

	return nil
}

// Close 关闭注册中心客户端
func (r *NacosRegistry) Close() error {
	// Nacos SDK没有提供明确的关闭方法，所以我们清理内部状态
	r.instancesLock.Lock()
	r.serviceInstances = make(map[string][]ServiceInstance)
	r.instancesLock.Unlock()

	r.subscriptionsLock.Lock()
	r.subscriptions = make(map[string]func(instances []ServiceInstance))
	r.subscriptionsLock.Unlock()

	return nil
}
