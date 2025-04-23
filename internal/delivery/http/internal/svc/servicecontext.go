package svc

import (
	"wz-backend-go/internal/delivery/http/internal/config"
	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/repository/mysql"
	"wz-backend-go/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 服务
	TenantService service.TenantService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.DB.DataSource)

	// 初始化租户仓库
	tenantRepo := mysql.NewTenantRepository(conn)
	
	// 初始化租户服务
	tenantService := service.NewTenantService(tenantRepo)

	return &ServiceContext{
		Config:        c,
		TenantService: tenantService,
	}
}
