package trade

import (
	"context"
	"fmt"

	"wz-backend-go/internal/svc"
	"wz-backend-go/api/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(id string, req *http.UpdateOrderStatusReq) (*http.SuccessResp, error) {
	// 检查订单是否存在
	order, err := l.svcCtx.TradeRepo.GetOrderById(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 检查状态是否有效
	validStatus := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"completed": true,
		"cancelled": true,
		"refunded":  true,
	}

	if !validStatus[req.Status] {
		return nil, fmt.Errorf("无效的订单状态: %s", req.Status)
	}

	// 获取操作人ID
	var operatorId int64 = 0
	if req.OperatorId > 0 {
		operatorId = req.OperatorId
	}

	// 调用仓库层更新订单状态
	err = l.svcCtx.TradeRepo.UpdateOrderStatus(l.ctx, id, req.Status, req.Reason, operatorId)
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 返回成功结果
	return &http.SuccessResp{
		Success: true,
		Message: fmt.Sprintf("订单状态已更新为: %s", req.Status),
	}, nil
} 