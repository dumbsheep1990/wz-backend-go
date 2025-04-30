package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/registry"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInstanceStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInstanceStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInstanceStatusLogic {
	return &UpdateInstanceStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateInstanceStatus 更新服务实例状态
func (l *UpdateInstanceStatusLogic) UpdateInstanceStatus(req *types.UpdateInstanceStatusReq) (resp *types.UpdateInstanceStatusResp, err error) {
	// 参数验证
	if req.ServiceName == "" {
		return nil, errorServiceNameRequired
	}
	if req.InstanceID == "" {
		return nil, errorInstanceIDRequired
	}
	
	// 验证状态值是否有效
	var status registry.InstanceStatus
	switch req.Status {
	case string(registry.StatusUp):
		status = registry.StatusUp
	case string(registry.StatusDown):
		status = registry.StatusDown
	case string(registry.StatusStarting):
		status = registry.StatusStarting
	case string(registry.StatusOutOfService):
		status = registry.StatusOutOfService
	default:
		return nil, errorInvalidStatus
	}
	
	// 更新实例状态
	err = l.svcCtx.InstanceManager.UpdateInstanceStatus(req.ServiceName, req.InstanceID, status)
	if err != nil {
		l.Logger.Errorf("更新服务实例状态失败: %v", err)
		if err == registry.ErrInstanceNotFound {
			return nil, errorInstanceNotFound
		}
		return nil, errorInternalServer
	}
	
	return &types.UpdateInstanceStatusResp{
		Success: true,
	}, nil
}
