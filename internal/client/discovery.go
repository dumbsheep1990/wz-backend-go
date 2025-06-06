package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceDiscovery 负责服务发现和地址解析
type ServiceDiscovery struct {
	client     *clientv3.Client
	ctx        context.Context
	cancel     context.CancelFunc
	watchCh    clientv3.WatchChan
	serviceMap sync.Map
	factory    *ClientFactory
	prefix     string
	mutex      sync.RWMutex
}

// NewServiceDiscovery 创建服务发现客户端
func NewServiceDiscovery(endpoints []string, factory *ClientFactory) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("连接服务发现系统失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ServiceDiscovery{
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		factory: factory,
		prefix:  "/services/",
	}, nil
}

// Start 启动服务发现
func (d *ServiceDiscovery) Start() error {
	// 初始化获取所有服务
	err := d.loadServices()
	if err != nil {
		return err
	}

	// 监听服务变化
	go d.watch()
	return nil
}

// loadServices 加载所有服务
func (d *ServiceDiscovery) loadServices() error {
	resp, err := d.client.Get(d.ctx, d.prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("获取服务列表失败: %v", err)
	}

	for _, kv := range resp.Kvs {
		d.updateServiceMap(string(kv.Key), string(kv.Value))
	}

	// 更新工厂的服务地址
	d.updateClientFactory()
	return nil
}

// updateServiceMap 更新服务映射
func (d *ServiceDiscovery) updateServiceMap(key, value string) {
	key = strings.TrimPrefix(key, d.prefix)
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return
	}

	serviceType := parts[0]
	d.serviceMap.Store(serviceType, value)
	log.Printf("服务 [%s] 地址更新为: %s\n", serviceType, value)
}

// updateClientFactory 更新客户端工厂
func (d *ServiceDiscovery) updateClientFactory() {
	addrMap := make(map[ServiceType]string)
	
	d.serviceMap.Range(func(key, value interface{}) bool {
		serviceTypeStr, ok1 := key.(string)
		addr, ok2 := value.(string)
		
		if ok1 && ok2 {
			serviceType := ServiceType(serviceTypeStr)
			addrMap[serviceType] = addr
		}
		return true
	})

	// 更新工厂地址
	d.factory.UpdateServiceAddrs(addrMap)
}

// watch 监听服务变化
func (d *ServiceDiscovery) watch() {
	watchCh := d.client.Watch(d.ctx, d.prefix, clientv3.WithPrefix())
	for {
		select {
		case <-d.ctx.Done():
			return
		case resp := <-watchCh:
			for _, ev := range resp.Events {
				switch ev.Type {
				case clientv3.EventTypePut: // 新增或更新
					d.updateServiceMap(string(ev.Kv.Key), string(ev.Kv.Value))
				case clientv3.EventTypeDelete: // 删除
					key := strings.TrimPrefix(string(ev.Kv.Key), d.prefix)
					parts := strings.Split(key, "/")
					if len(parts) < 2 {
						continue
					}
					serviceType := parts[0]
					d.serviceMap.Delete(serviceType)
					log.Printf("服务 [%s] 已下线\n", serviceType)
				}
			}
			
			// 更新工厂的服务地址
			d.updateClientFactory()
		}
	}
}

// GetServiceAddr 获取服务地址
func (d *ServiceDiscovery) GetServiceAddr(serviceType ServiceType) (string, bool) {
	value, ok := d.serviceMap.Load(string(serviceType))
	if !ok {
		return "", false
	}
	
	addr, ok := value.(string)
	return addr, ok
}

// Close 关闭服务发现客户端
func (d *ServiceDiscovery) Close() error {
	d.cancel()
	return d.client.Close()
}
