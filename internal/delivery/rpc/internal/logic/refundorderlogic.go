package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundOrderLogic {
	return &RefundOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundOrderLogic) RefundOrder(in *trade.RefundOrderRequest) (*trade.RefundOrderResponse, error) {
	// 构建退款处理命令
	cmd := dto.RefundOrderCommand{
		OrderID:       in.OrderId,
		RefundAmount:  dto.Money{Amount: in.Amount.Amount, Currency: in.Amount.Currency},
		AdminNote:     in.AdminNote,
		RefundMethod:  in.RefundMethod,
		TransactionID: in.TransactionId,
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.RefundOrder(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to process refund", logx.Field("error", err),
			logx.Field("orderId", in.OrderId))
		return &trade.RefundOrderResponse{
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
	return &trade.RefundOrderResponse{
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
			Note:              orderDTO.Note,
			TrackingNumber:    orderDTO.TrackingNumber,
			PaymentMethod:     orderDTO.PaymentMethod,
			ShippingMethod:    orderDTO.ShippingMethod,
			RefundRequestedAt: orderDTO.RefundRequestedAt.Unix(),
			RefundReason:      orderDTO.RefundReason,
			RefundAmount: &trade.Money{
				Amount:   orderDTO.RefundAmount.Amount,
				Currency: orderDTO.RefundAmount.Currency,
			},
			RefundedAt:    orderDTO.RefundedAt.Unix(),
			RefundMethod:  orderDTO.RefundMethod,
			RefundTransaction: orderDTO.RefundTransaction,
			CreatedAt:     orderDTO.CreatedAt.Unix(),
			UpdatedAt:     orderDTO.UpdatedAt.Unix(),
		},
	}, nil
}
