package svc

import (
	"wz-backend-go/internal/delivery/rpc/internal/config"
	"wz-backend-go/internal/application/order/service"
	"wz-backend-go/internal/domain/order/repository"
	"wz-backend-go/internal/infrastructure/persistence"
)

type ServiceContext struct {
	Config config.Config
	OrderApplicationService service.OrderApplicationService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db := persistence.NewDBConnection(c.DB)
	
	// 初始化订单仓储
	orderRepo := repository.NewOrderRepository(db)
	
	// 初始化订单应用服务
	orderAppService := service.NewOrderApplicationService(orderRepo)
	
	return &ServiceContext{
		Config: c,
		OrderApplicationService: orderAppService,
	}
}
