package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeregisterInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeregisterInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeregisterInstanceLogic {
	return &DeregisterInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeregisterInstance 注销服务实例
func (l *DeregisterInstanceLogic) DeregisterInstance(req *types.DeregisterInstanceReq) (resp *types.DeregisterInstanceResp, err error) {
	// 参数验证
	if req.ServiceName == "" {
		return nil, errorServiceNameRequired
	}
	if req.InstanceID == "" {
		return nil, errorInstanceIDRequired
	}

	// 先停止心跳
	_ = l.svcCtx.InstanceManager.StopHeartbeat(req.ServiceName, req.InstanceID)

	// 注销服务实例
	err = l.svcCtx.InstanceManager.DeregisterInstance(req.ServiceName, req.InstanceID)
	if err != nil {
		l.Logger.Errorf("注销服务实例失败: %v", err)
		return nil, errorDeregisterFailed
	}

	return &types.DeregisterInstanceResp{
		Success: true,
	}, nil
}
