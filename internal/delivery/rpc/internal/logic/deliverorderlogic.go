package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeliverOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeliverOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeliverOrderLogic {
	return &DeliverOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeliverOrderLogic) DeliverOrder(in *trade.DeliverOrderRequest) (*trade.DeliverOrderResponse, error) {
	// 构建送达订单命令
	cmd := dto.DeliverOrderCommand{
		OrderID: in.OrderId,
		// 可以添加额外的送达信息，如签收人等
		SignedBy: in.SignedBy,
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.DeliverOrder(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to deliver order", logx.Field("error", err), 
			logx.Field("orderId", in.OrderId))
		return &trade.DeliverOrderResponse{
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
	return &trade.DeliverOrderResponse{
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
			ShippedAt:      orderDTO.ShippedAt.Unix(),
			DeliveredAt:    orderDTO.DeliveredAt.Unix(),
		},
	}, nil
}
