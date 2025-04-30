package svc

import (
	"wz-backend-go/internal/delivery/http/internal/config"
	"wz-backend-go/internal/registry"
	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/repository/mysql"
	"wz-backend-go/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 服务
	TenantService service.TenantService
	// 服务注册与发现
	Registry         registry.ServiceRegistry
	InstanceManager  registry.InstanceManager
	HealthChecker    registry.HealthChecker
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.DB.DataSource)

	// 初始化租户仓库
	tenantRepo := mysql.NewTenantRepository(conn)
	
	// 初始化租户服务
	tenantService := service.NewTenantService(tenantRepo)

	// 初始化服务注册与发现
	nacosConfig := &registry.NacosConfig{
		ServerAddr: c.Registry.ServerAddr,
		ServerPort: c.Registry.ServerPort,
		Namespace:  c.Registry.Namespace,
		Group:      c.Registry.Group,
		LogDir:     c.Registry.LogDir,
		CacheDir:   c.Registry.CacheDir,
		LogLevel:   c.Registry.LogLevel,
		Username:   c.Registry.Username,
		Password:   c.Registry.Password,
	}
	
	// 创建注册中心客户端
	nacosRegistry, err := registry.NewNacosRegistry(nacosConfig)
	if err != nil {
		panic(err)
	}
	
	// 创建实例管理器
	instanceManager := registry.NewNacosInstanceManager(nacosRegistry)
	
	// 创建健康检查服务
	healthChecker := registry.NewHealthCheckServer(c.Registry.HealthCheckPort, c.Registry.CheckInterval)
	
	return &ServiceContext{
		Config:           c,
		TenantService:    tenantService,
		Registry:         nacosRegistry,
		InstanceManager:  instanceManager,
		HealthChecker:    healthChecker,
	}
}
