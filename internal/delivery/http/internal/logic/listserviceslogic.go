package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListServicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListServicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListServicesLogic {
	return &ListServicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListServices 获取所有已注册的服务列表
func (l *ListServicesLogic) ListServices(req *types.ListServicesReq) (resp *types.ListServicesResp, err error) {
	// 这里我们需要通过Nacos SDK获取所有服务列表
	// 由于Nacos SDK没有直接提供获取所有服务的方法，我们可以通过服务订阅列表来获取
	// 或者从配置中读取已知服务列表
	
	// 这里是示例实现，实际项目中可能需要从配置或者数据库中获取
	services := []string{
		"user-service",
		"content-service",
		"search-service",
		"transaction-service",
		"public-api-service",
		"interaction-service",
		"ai-service",
		"notification-service",
		"file-service",
		"statistics-service",
	}
	
	// 实际项目中，应该通过Nacos SDK获取服务列表
	// 例如：services := l.svcCtx.Registry.ListServices()
	
	return &types.ListServicesResp{
		Services: services,
	}, nil
}
