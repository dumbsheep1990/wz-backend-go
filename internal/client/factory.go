package client

import (
	"fmt"
	"sync"
)

// ServiceType 表示微服务类型
type ServiceType string

const (
	UserService      ServiceType = "user"
	ContentService   ServiceType = "content"
	CommunityService ServiceType = "community"
	TradeService     ServiceType = "trade"
	SiteService      ServiceType = "site"
	ComponentService ServiceType = "component"
	RenderService    ServiceType = "render"
)

// ClientFactory 负责创建和管理各种微服务客户端
type ClientFactory struct {
	mutex           sync.RWMutex
	serviceAddrs    map[ServiceType]string
	defaultAddrs    map[ServiceType]string
	registeredAddrs map[ServiceType]string // 从注册中心获取的地址
}

// NewClientFactory 创建一个新的客户端工厂
func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		serviceAddrs: make(map[ServiceType]string),
		defaultAddrs: map[ServiceType]string{
			UserService:      "localhost:50051",
			ContentService:   "localhost:50052",
			TradeService:     "localhost:50053",
			CommunityService: "localhost:50054",
			SiteService:      "localhost:50055",
			ComponentService: "localhost:50056",
			RenderService:    "localhost:50057",
		},
		registeredAddrs: make(map[ServiceType]string),
	}
}

// RegisterServiceAddr 注册服务地址
func (f *ClientFactory) RegisterServiceAddr(serviceType ServiceType, addr string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.serviceAddrs[serviceType] = addr
}

// UpdateServiceAddrs 从服务发现系统批量更新服务地址
func (f *ClientFactory) UpdateServiceAddrs(addrs map[ServiceType]string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for serviceType, addr := range addrs {
		f.registeredAddrs[serviceType] = addr
	}
}

// getServiceAddr 获取服务地址，优先级：手动注册 > 服务发现 > 默认地址
func (f *ClientFactory) getServiceAddr(serviceType ServiceType) string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	// 优先使用手动注册的地址
	if addr, ok := f.serviceAddrs[serviceType]; ok && addr != "" {
		return addr
	}
	
	// 其次使用从服务发现获取的地址
	if addr, ok := f.registeredAddrs[serviceType]; ok && addr != "" {
		return addr
	}
	
	// 最后使用默认地址
	if addr, ok := f.defaultAddrs[serviceType]; ok {
		return addr
	}
	
	return ""
}

// NewUserClient 创建用户服务客户端
func (f *ClientFactory) NewUserClient() (*UserClient, error) {
	addr := f.getServiceAddr(UserService)
	if addr == "" {
		return nil, fmt.Errorf("未配置用户服务地址")
	}
	return NewUserClient(addr)
}

// NewContentClient 创建内容服务客户端
func (f *ClientFactory) NewContentClient() (*ContentClient, error) {
	addr := f.getServiceAddr(ContentService)
	if addr == "" {
		return nil, fmt.Errorf("未配置内容服务地址")
	}
	return NewContentClient(addr)
}

// NewCommunityClient 创建社区服务客户端
func (f *ClientFactory) NewCommunityClient() (*CommunityClient, error) {
	addr := f.getServiceAddr(CommunityService)
	if addr == "" {
		return nil, fmt.Errorf("未配置社区服务地址")
	}
	return NewCommunityClient(addr)
}

// NewTradeClient 创建交易服务客户端
func (f *ClientFactory) NewTradeClient() (*TradeClient, error) {
	addr := f.getServiceAddr(TradeService)
	if addr == "" {
		return nil, fmt.Errorf("未配置交易服务地址")
	}
	return NewTradeClient(addr)
}

// NewSiteClient 创建站点服务客户端
func (f *ClientFactory) NewSiteClient() (*SiteClient, error) {
	addr := f.getServiceAddr(SiteService)
	if addr == "" {
		return nil, fmt.Errorf("未配置站点服务地址")
	}
	return NewSiteClient(addr)
}

// NewComponentClient 创建组件服务客户端
func (f *ClientFactory) NewComponentClient() (*ComponentClient, error) {
	addr := f.getServiceAddr(ComponentService)
	if addr == "" {
		return nil, fmt.Errorf("未配置组件服务地址")
	}
	return NewComponentClient(addr)
}

// NewRenderClient 创建渲染服务客户端
func (f *ClientFactory) NewRenderClient() (*RenderClient, error) {
	addr := f.getServiceAddr(RenderService)
	if addr == "" {
		return nil, fmt.Errorf("未配置渲染服务地址")
	}
	return NewRenderClient(addr)
}
