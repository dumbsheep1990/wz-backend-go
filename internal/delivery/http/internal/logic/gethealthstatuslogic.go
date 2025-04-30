package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHealthStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHealthStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHealthStatusLogic {
	return &GetHealthStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetHealthStatus 获取所有服务的健康状态
func (l *GetHealthStatusLogic) GetHealthStatus() (resp *types.HealthStatusResp, err error) {
	// 获取健康检查器状态
	healthStatus := l.svcCtx.HealthChecker.GetStatus()
	
	// 获取所有已注册服务
	services := make(map[string][]string) // 服务名 -> 实例ID列表
	
	// 这里为简化实现，我们手动列出所有服务
	// 实际项目中，应该通过SDK或数据库获取
	serviceNames := []string{
		"user-service",
		"content-service",
		"search-service",
		"transaction-service",
		"public-api-service",
	}
	
	// 为每个服务获取实例列表
	for _, serviceName := range serviceNames {
		instances, err := l.svcCtx.InstanceManager.GetInstances(serviceName)
		if err != nil {
			l.Logger.Errorf("获取服务 %s 实例列表失败: %v", serviceName, err)
			continue
		}
		
		instanceIDs := make([]string, 0, len(instances))
		for _, instance := range instances {
			instanceIDs = append(instanceIDs, instance.InstanceID)
		}
		
		services[serviceName] = instanceIDs
	}
	
	// 构建健康状态响应
	respData := &types.HealthStatusResp{
		Services: make(map[string]types.ServiceHealthInfo),
	}
	
	// 计算每个服务的健康状态
	for serviceName, instanceIDs := range services {
		totalInstances := len(instanceIDs)
		healthyInstances := 0
		
		// 这里简化实现，假设实例状态为"UP"的就是健康的
		for _, instanceID := range instanceIDs {
			instance, err := l.svcCtx.InstanceManager.GetInstance(serviceName, instanceID)
			if err == nil && instance.Status == "UP" {
				healthyInstances++
			}
		}
		
		// 计算健康比率
		healthyRate := 0.0
		if totalInstances > 0 {
			healthyRate = float64(healthyInstances) / float64(totalInstances)
		}
		
		// 确定服务整体状态
		status := "DOWN"
		if healthyRate >= 0.5 {
			status = "UP"
		}
		
		// 转换健康检查结果
		checks := make(map[string]types.CheckResult)
		for name, result := range healthStatus {
			checks[name] = types.CheckResult{
				Status:    string(result.Status),
				Message:   result.Message,
				Timestamp: result.Timestamp,
			}
		}
		
		// 添加到响应
		respData.Services[serviceName] = types.ServiceHealthInfo{
			ServiceName: serviceName,
			Status:      status,
			Instances:   totalInstances,
			HealthyRate: healthyRate,
			Checks:      checks,
		}
	}
	
	return respData, nil
}
