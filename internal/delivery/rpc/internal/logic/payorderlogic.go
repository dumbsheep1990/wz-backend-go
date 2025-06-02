package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderLogic) PayOrder(in *trade.PayOrderRequest) (*trade.PayOrderResponse, error) {
	// 构建支付订单命令
	cmd := dto.PayOrderCommand{
		OrderID:       in.OrderId,
		PaymentMethod: in.PaymentMethod,
		TransactionID: in.TransactionId,
		PaidAmount: dto.Money{
			Amount:   in.PaidAmount.Amount,
			Currency: in.PaidAmount.Currency,
		},
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.PayOrder(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to pay order", logx.Field("error", err), 
			logx.Field("orderId", in.OrderId))
		return &trade.PayOrderResponse{
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
	return &trade.PayOrderResponse{
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
			PaidAt:         orderDTO.PaidAt.Unix(),
			Transaction:    orderDTO.Transaction,
		},
	}, nil
}
