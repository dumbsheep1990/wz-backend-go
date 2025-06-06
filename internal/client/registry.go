package client

import (
	"fmt"
	"sync"
)

// ClientRegistry 管理客户端实例和生命周期
type ClientRegistry struct {
	mutex     sync.RWMutex
	factory   *ClientFactory
	clients   map[ServiceType]interface{}
	isClosing bool
}

// NewClientRegistry 创建一个新的客户端注册表
func NewClientRegistry(factory *ClientFactory) *ClientRegistry {
	if factory == nil {
		factory = NewClientFactory()
	}
	return &ClientRegistry{
		factory: factory,
		clients: make(map[ServiceType]interface{}),
	}
}

// GetUserClient 获取或创建用户服务客户端
func (r *ClientRegistry) GetUserClient() (*UserClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[UserService]; ok {
		if userClient, ok := client.(*UserClient); ok {
			return userClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewUserClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[UserService] = client
	return client, nil
}

// GetContentClient 获取或创建内容服务客户端
func (r *ClientRegistry) GetContentClient() (*ContentClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[ContentService]; ok {
		if contentClient, ok := client.(*ContentClient); ok {
			return contentClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewContentClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[ContentService] = client
	return client, nil
}

// GetCommunityClient 获取或创建社区服务客户端
func (r *ClientRegistry) GetCommunityClient() (*CommunityClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[CommunityService]; ok {
		if communityClient, ok := client.(*CommunityClient); ok {
			return communityClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewCommunityClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[CommunityService] = client
	return client, nil
}

// GetTradeClient 获取或创建交易服务客户端
func (r *ClientRegistry) GetTradeClient() (*TradeClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[TradeService]; ok {
		if tradeClient, ok := client.(*TradeClient); ok {
			return tradeClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewTradeClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[TradeService] = client
	return client, nil
}

// GetSiteClient 获取或创建站点服务客户端
func (r *ClientRegistry) GetSiteClient() (*SiteClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[SiteService]; ok {
		if siteClient, ok := client.(*SiteClient); ok {
			return siteClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewSiteClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[SiteService] = client
	return client, nil
}

// GetComponentClient 获取或创建组件服务客户端
func (r *ClientRegistry) GetComponentClient() (*ComponentClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[ComponentService]; ok {
		if componentClient, ok := client.(*ComponentClient); ok {
			return componentClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewComponentClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[ComponentService] = client
	return client, nil
}

// GetRenderClient 获取或创建渲染服务客户端
func (r *ClientRegistry) GetRenderClient() (*RenderClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isClosing {
		return nil, fmt.Errorf("客户端注册表正在关闭")
	}

	// 检查是否已有实例
	if client, ok := r.clients[RenderService]; ok {
		if renderClient, ok := client.(*RenderClient); ok {
			return renderClient, nil
		}
	}

	// 创建新实例
	client, err := r.factory.NewRenderClient()
	if err != nil {
		return nil, err
	}

	// 存入注册表
	r.clients[RenderService] = client
	return client, nil
}

// Close 关闭所有客户端连接
func (r *ClientRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.isClosing = true
	var lastErr error

	// 关闭所有客户端
	for serviceType, client := range r.clients {
		var err error

		switch c := client.(type) {
		case *UserClient:
			err = c.Close()
		case *ContentClient:
			err = c.Close()
		case *CommunityClient:
			err = c.Close()
		case *TradeClient:
			err = c.Close()
		case *SiteClient:
			err = c.Close()
		case *ComponentClient:
			err = c.Close()
		case *RenderClient:
			err = c.Close()
		}

		if err != nil {
			lastErr = fmt.Errorf("关闭 %s 服务客户端失败: %v", serviceType, err)
		}

		delete(r.clients, serviceType)
	}

	return lastErr
}
