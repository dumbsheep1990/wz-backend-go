package logic

import (
	"context"
	"fmt"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(in *trade.UpdateOrderStatusRequest) (*trade.UpdateOrderStatusResponse, error) {
	// 获取订单ID和目标状态
	orderID := in.OrderId
	targetStatus := in.TargetStatus

	// 根据目标状态调用不同的应用服务方法
	var orderDTO *dto.OrderDTO
	var err error

	switch targetStatus {
	case "submitted":
		// 提交订单
		cmd := dto.SubmitOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.SubmitOrder(l.ctx, cmd)

	case "paid":
		// 支付订单
		cmd := dto.PayOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.PayOrder(l.ctx, cmd)

	case "shipped":
		// 发货订单
		cmd := dto.ShipOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.ShipOrder(l.ctx, cmd)

	case "delivered":
		// 送达订单
		cmd := dto.DeliverOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.DeliverOrder(l.ctx, cmd)

	case "completed":
		// 完成订单
		cmd := dto.CompleteOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.CompleteOrder(l.ctx, cmd)

	case "cancelled":
		// 取消订单
		cmd := dto.CancelOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.CancelOrder(l.ctx, cmd)

	case "refund_requested":
		// 申请退款
		cmd := dto.RequestRefundCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.RequestRefund(l.ctx, cmd)

	case "refunded":
		// 已退款
		cmd := dto.RefundOrderCommand{OrderID: orderID}
		orderDTO, err = l.svcCtx.OrderApplicationService.RefundOrder(l.ctx, cmd)

	default:
		return &trade.UpdateOrderStatusResponse{
			Success: false,
			Error:   fmt.Sprintf("Unsupported target status: %s", targetStatus),
		}, nil
	}

	if err != nil {
		l.Error("Failed to update order status", logx.Field("error", err), 
			logx.Field("orderId", orderID), 
			logx.Field("targetStatus", targetStatus))
		return &trade.UpdateOrderStatusResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 转换订单项
	var respItems []*trade.OrderItem
	for _, item := range orderDTO.Items {
		respItems = append(respItems, &trade.OrderItem{
			Id:         item.ID,
			ProductId:  item.ProductID,
			ProductName: item.ProductName,
			ProductSku: item.ProductSKU,
			Quantity:   item.Quantity,
			UnitPrice: &trade.Money{
				Amount:   item.UnitPrice.Amount,
				Currency: item.UnitPrice.Currency,
			},
			TotalPrice: &trade.Money{
				Amount:   item.TotalPrice.Amount,
				Currency: item.TotalPrice.Currency,
			},
			Attributes: item.Attributes,
		})
	}

	// 转换折扣
	var respDiscounts []*trade.OrderDiscount
	for _, discount := range orderDTO.Discounts {
		respDiscounts = append(respDiscounts, &trade.OrderDiscount{
			Id:           discount.ID,
			DiscountType: discount.DiscountType,
			DiscountName: discount.DiscountName,
			Amount: &trade.Money{
				Amount:   discount.Amount.Amount,
				Currency: discount.Amount.Currency,
			},
		})
	}

	// 构建响应
	return &trade.UpdateOrderStatusResponse{
		Success: true,
		Order: &trade.Order{
			Id:           orderDTO.ID,
			OrderNumber:  orderDTO.OrderNumber,
			CustomerId:   orderDTO.CustomerID,
			Status:       orderDTO.Status,
			StatusCode:   orderDTO.StatusCode,
			Items:        respItems,
			Discounts:    respDiscounts,
			ShippingFee: &trade.Money{
				Amount:   orderDTO.ShippingFee.Amount,
				Currency: orderDTO.ShippingFee.Currency,
			},
			Tax: &trade.Money{
				Amount:   orderDTO.Tax.Amount,
				Currency: orderDTO.Tax.Currency,
			},
			Subtotal: &trade.Money{
				Amount:   orderDTO.Subtotal.Amount,
				Currency: orderDTO.Subtotal.Currency,
			},
			TotalAmount: &trade.Money{
				Amount:   orderDTO.TotalAmount.Amount,
				Currency: orderDTO.TotalAmount.Currency,
			},
			Note:           orderDTO.Note,
			TrackingNumber: orderDTO.TrackingNumber,
			PaymentMethod:  orderDTO.PaymentMethod,
			ShippingMethod: orderDTO.ShippingMethod,
			CreatedAt:      orderDTO.CreatedAt.Unix(),
			UpdatedAt:      orderDTO.UpdatedAt.Unix(),
		},
	}, nil
}
