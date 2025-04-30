package logic

import (
	"context"
	"time"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/registry"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterInstanceLogic {
	return &RegisterInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RegisterInstance 注册服务实例
func (l *RegisterInstanceLogic) RegisterInstance(req *types.RegisterInstanceReq) (resp *types.RegisterInstanceResp, err error) {
	// 参数验证
	if req.ServiceName == "" {
		return nil, errorServiceNameRequired
	}
	if req.IP == "" {
		return nil, NewCodeError(400, "IP地址不能为空")
	}
	if req.Port <= 0 {
		return nil, NewCodeError(400, "端口号无效")
	}

	// 准备实例信息
	instanceInfo := &registry.InstanceInfo{
		ServiceName:     req.ServiceName,
		IP:              req.IP,
		Port:            req.Port,
		Status:          registry.StatusUp,
		Metadata:        req.Metadata,
		StartTime:       time.Now(),
		HealthCheckPath: req.HealthCheckPath,
	}

	// 注册服务实例
	err = l.svcCtx.InstanceManager.RegisterInstance(req.ServiceName, instanceInfo)
	if err != nil {
		l.Logger.Errorf("注册服务实例失败: %v", err)
		return nil, errorRegisterFailed
	}

	// 启动心跳
	err = l.svcCtx.InstanceManager.StartHeartbeat(req.ServiceName, instanceInfo.InstanceID, 15*time.Second)
	if err != nil {
		l.Logger.Errorf("启动心跳失败: %v", err)
		// 虽然心跳启动失败，但实例已注册成功，所以继续返回成功
	}

	return &types.RegisterInstanceResp{
		InstanceID: instanceInfo.InstanceID,
	}, nil
}
