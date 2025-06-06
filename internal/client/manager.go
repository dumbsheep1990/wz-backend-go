package client

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Options 客户端管理器配置选项
type Options struct {
	// Etcd 服务发现配置
	EtcdEndpoints []string
	// 是否启用服务发现
	EnableDiscovery bool
	// 服务地址手动配置
	ServiceAddrs map[ServiceType]string
	// 重试配置
	MaxRetries          int
	RetryInterval       time.Duration
	RequestTimeout      time.Duration
	EnableCircuitBreaker bool
}

// DefaultOptions 返回默认选项
func DefaultOptions() *Options {
	return &Options{
		EtcdEndpoints:       []string{"localhost:2379"},
		EnableDiscovery:     true,
		ServiceAddrs:        make(map[ServiceType]string),
		MaxRetries:          3,
		RetryInterval:       time.Second,
		RequestTimeout:      time.Second * 10,
		EnableCircuitBreaker: true,
	}
}

// ClientManager 客户端管理器
type ClientManager struct {
	factory    *ClientFactory
	registry   *ClientRegistry
	discovery  *ServiceDiscovery
	options    *Options
	ctx        context.Context
	cancelFunc context.CancelFunc
	isRunning  bool
}

// NewClientManager 创建客户端管理器
func NewClientManager(options *Options) *ClientManager {
	if options == nil {
		options = DefaultOptions()
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	return &ClientManager{
		factory:    NewClientFactory(),
		options:    options,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Init 初始化客户端管理器
func (m *ClientManager) Init() error {
	// 注册手动配置的服务地址
	for serviceType, addr := range m.options.ServiceAddrs {
		m.factory.RegisterServiceAddr(serviceType, addr)
	}

	// 创建客户端注册表
	m.registry = NewClientRegistry(m.factory)

	// 初始化服务发现
	if m.options.EnableDiscovery {
		discovery, err := NewServiceDiscovery(m.options.EtcdEndpoints, m.factory)
		if err != nil {
			return fmt.Errorf("初始化服务发现失败: %v", err)
		}
		
		m.discovery = discovery
		if err := m.discovery.Start(); err != nil {
			return fmt.Errorf("启动服务发现失败: %v", err)
		}
	}

	m.isRunning = true
	return nil
}

// GetUserClient 获取用户服务客户端
func (m *ClientManager) GetUserClient() (*UserClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetUserClient()
}

// GetContentClient 获取内容服务客户端
func (m *ClientManager) GetContentClient() (*ContentClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetContentClient()
}

// GetCommunityClient 获取社区服务客户端
func (m *ClientManager) GetCommunityClient() (*CommunityClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetCommunityClient()
}

// GetTradeClient 获取交易服务客户端
func (m *ClientManager) GetTradeClient() (*TradeClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetTradeClient()
}

// GetSiteClient 获取站点服务客户端
func (m *ClientManager) GetSiteClient() (*SiteClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetSiteClient()
}

// GetComponentClient 获取组件服务客户端
func (m *ClientManager) GetComponentClient() (*ComponentClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetComponentClient()
}

// GetRenderClient 获取渲染服务客户端
func (m *ClientManager) GetRenderClient() (*RenderClient, error) {
	if !m.isRunning {
		return nil, fmt.Errorf("客户端管理器未初始化")
	}
	return m.registry.GetRenderClient()
}

// Execute 执行带有重试逻辑的请求
func (m *ClientManager) Execute(operation func() error) error {
	var err error
	for i := 0; i <= m.options.MaxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		// 如果不是最后一次重试，则等待一段时间
		if i < m.options.MaxRetries {
			log.Printf("操作失败，将在 %v 后重试 (%d/%d): %v\n", 
				m.options.RetryInterval, i+1, m.options.MaxRetries, err)
			time.Sleep(m.options.RetryInterval)
		}
	}
	return fmt.Errorf("操作在 %d 次重试后仍失败: %v", m.options.MaxRetries, err)
}

// ExecuteWithTimeout 执行带有超时和重试逻辑的请求
func (m *ClientManager) ExecuteWithTimeout(operation func(ctx context.Context) error) error {
	return m.Execute(func() error {
		ctx, cancel := context.WithTimeout(m.ctx, m.options.RequestTimeout)
		defer cancel()
		return operation(ctx)
	})
}

// Close 关闭客户端管理器
func (m *ClientManager) Close() error {
	m.cancelFunc()
	m.isRunning = false

	var errs []error
	
	// 关闭客户端注册表
	if m.registry != nil {
		if err := m.registry.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭客户端注册表失败: %v", err))
		}
	}

	// 关闭服务发现
	if m.discovery != nil {
		if err := m.discovery.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭服务发现失败: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("关闭客户端管理器时发生错误: %v", errs)
	}
	
	return nil
}
