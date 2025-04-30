package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServiceDependenciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServiceDependenciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceDependenciesLogic {
	return &GetServiceDependenciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetServiceDependencies 获取服务依赖关系
func (l *GetServiceDependenciesLogic) GetServiceDependencies(req *types.GetServiceDependenciesReq) (resp *types.GetServiceDependenciesResp, err error) {
	// 参数验证
	if req.ServiceName == "" {
		return nil, errorServiceNameRequired
	}
	
	// 这里我们需要从服务依赖管理系统或配置中获取依赖关系
	// 由于项目中可能还没有实现服务依赖管理，这里使用模拟数据

	// 模拟服务依赖关系
	dependencies := map[string][]types.ServiceDependency{
		"user-service": {
			{Source: "user-service", Target: "notification-service", Type: "RPC", Description: "发送用户通知"},
			{Source: "user-service", Target: "file-service", Type: "REST", Description: "用户头像存储"},
		},
		"content-service": {
			{Source: "content-service", Target: "user-service", Type: "RPC", Description: "获取用户信息"},
			{Source: "content-service", Target: "search-service", Type: "RPC", Description: "内容索引"},
			{Source: "content-service", Target: "file-service", Type: "REST", Description: "内容文件存储"},
		},
		"search-service": {
			{Source: "search-service", Target: "content-service", Type: "RPC", Description: "获取内容数据"},
		},
		"transaction-service": {
			{Source: "transaction-service", Target: "user-service", Type: "RPC", Description: "用户账户操作"},
			{Source: "transaction-service", Target: "notification-service", Type: "RPC", Description: "交易通知"},
		},
		"public-api-service": {
			{Source: "public-api-service", Target: "user-service", Type: "RPC", Description: "用户认证和授权"},
			{Source: "public-api-service", Target: "content-service", Type: "RPC", Description: "内容管理"},
			{Source: "public-api-service", Target: "search-service", Type: "RPC", Description: "内容搜索"},
			{Source: "public-api-service", Target: "transaction-service", Type: "RPC", Description: "交易处理"},
		},
		"interaction-service": {
			{Source: "interaction-service", Target: "user-service", Type: "RPC", Description: "用户信息"},
			{Source: "interaction-service", Target: "content-service", Type: "RPC", Description: "内容数据"},
			{Source: "interaction-service", Target: "notification-service", Type: "RPC", Description: "互动通知"},
		},
		"ai-service": {
			{Source: "ai-service", Target: "content-service", Type: "RPC", Description: "内容分析"},
			{Source: "ai-service", Target: "user-service", Type: "RPC", Description: "用户行为分析"},
		},
		"notification-service": {
			{Source: "notification-service", Target: "user-service", Type: "RPC", Description: "获取用户联系方式"},
		},
		"file-service": {},
		"statistics-service": {
			{Source: "statistics-service", Target: "user-service", Type: "RPC", Description: "用户统计"},
			{Source: "statistics-service", Target: "content-service", Type: "RPC", Description: "内容统计"},
			{Source: "statistics-service", Target: "transaction-service", Type: "RPC", Description: "交易统计"},
			{Source: "statistics-service", Target: "interaction-service", Type: "RPC", Description: "互动统计"},
		},
	}
	
	// 检查请求的服务是否存在
	serviceDeps, exist := dependencies[req.ServiceName]
	if !exist {
		return nil, errorServiceNotFound
	}
	
	// 构建响应
	return &types.GetServiceDependenciesResp{
		ServiceName:  req.ServiceName,
		Dependencies: serviceDeps,
	}, nil
}
