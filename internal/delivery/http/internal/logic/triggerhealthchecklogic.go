package logic

import (
	"context"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerHealthCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerHealthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerHealthCheckLogic {
	return &TriggerHealthCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TriggerHealthCheck 手动触发健康检查
func (l *TriggerHealthCheckLogic) TriggerHealthCheck(req *types.TriggerHealthCheckReq) (resp *types.TriggerHealthCheckResp, err error) {
	// 这里我们通过调用HealthChecker的performChecks方法来触发健康检查
	// 然后返回检查结果
	
	// 由于HealthChecker接口中没有提供performChecks这样的公开方法
	// 所以这里我们暂时只能通过GetStatus方法间接获取检查结果
	
	// 在实际项目中，可以扩展HealthChecker接口添加TriggerCheck方法
	// 或者通过其他方式实现手动触发健康检查
	
	// 获取健康检查状态
	healthStatus := l.svcCtx.HealthChecker.GetStatus()
	
	// 如果指定了服务名，需要过滤结果
	if req.ServiceName != "" {
		// 过滤特定服务的健康检查结果
		// 实际项目中需要根据具体的检查名称命名规则来实现
	}
	
	// 转换健康检查结果
	results := make(map[string]types.CheckResult)
	for name, result := range healthStatus {
		results[name] = types.CheckResult{
			Status:    string(result.Status),
			Message:   result.Message,
			Timestamp: result.Timestamp,
		}
	}
	
	return &types.TriggerHealthCheckResp{
		Success: true,
		Results: results,
	}, nil
}
