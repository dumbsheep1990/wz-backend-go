package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/registry"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServiceInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServiceInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceInstancesLogic {
	return &GetServiceInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetServiceInstances 获取指定服务的所有实例
func (l *GetServiceInstancesLogic) GetServiceInstances(req *types.GetServiceInstancesReq) (resp *types.GetServiceInstancesResp, err error) {
	// 参数验证
	if req.ServiceName == "" {
		return nil, errorServiceNameRequired
	}

	// 从实例管理器获取服务实例
	instances, err := l.svcCtx.InstanceManager.GetInstances(req.ServiceName)
	if err != nil {
		l.Logger.Errorf("获取服务实例失败: %v", err)
		return nil, errorInternalServer
	}

	// 转换为API响应格式
	instanceInfos := make([]types.ServiceInstanceInfo, 0, len(instances))
	for _, inst := range instances {
		instanceInfos = append(instanceInfos, types.ServiceInstanceInfo{
			InstanceID:    inst.InstanceID,
			ServiceName:   inst.ServiceName,
			IP:            inst.IP,
			Port:          inst.Port,
			Status:        string(inst.Status),
			Metadata:      inst.Metadata,
			StartTime:     inst.StartTime,
			LastHeartbeat: inst.LastHeartbeat,
			Healthy:       inst.Status == string(registry.StatusUp),
		})
	}

	return &types.GetServiceInstancesResp{
		Instances: instanceInfos,
	}, nil
}
