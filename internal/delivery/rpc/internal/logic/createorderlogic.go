package logic

import (
	"context"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 订单管理
func (l *CreateOrderLogic) CreateOrder(in *trade.CreateOrderRequest) (*trade.CreateOrderResponse, error) {
	// 创建命令DTO
	var items []dto.OrderItemCommand
	for _, item := range in.Items {
		items = append(items, dto.OrderItemCommand{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSku,
			Quantity:    item.Quantity,
			UnitPrice: dto.MoneyDTO{
				Amount:   item.UnitPrice.Amount,
				Currency: item.UnitPrice.Currency,
			},
			Attributes: item.Attributes,
		})
	}

	// 解析用户地址
	shippingAddress := dto.AddressDTO{
		Name:          in.ShippingAddress.Name,
		Phone:         in.ShippingAddress.Phone,
		Province:      in.ShippingAddress.Province,
		City:          in.ShippingAddress.City,
		District:      in.ShippingAddress.District,
		DetailAddress: in.ShippingAddress.DetailAddress,
		PostalCode:    in.ShippingAddress.PostalCode,
	}

	// 解析账单地址
	billingAddress := dto.AddressDTO{}
	if in.BillingAddress != nil {
		billingAddress = dto.AddressDTO{
			Name:          in.BillingAddress.Name,
			Phone:         in.BillingAddress.Phone,
			Province:      in.BillingAddress.Province,
			City:          in.BillingAddress.City,
			District:      in.BillingAddress.District,
			DetailAddress: in.BillingAddress.DetailAddress,
			PostalCode:    in.BillingAddress.PostalCode,
		}
	} else {
		// 如果没有提供账单地址，使用收货地址
		billingAddress = shippingAddress
	}

	// 构建创建订单命令
	cmd := dto.CreateOrderCommand{
		CustomerID:      in.CustomerId,
		Items:           items,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		Note:            in.Note,
	}

	// 调用应用服务
	orderDTO, err := l.svcCtx.OrderApplicationService.CreateOrder(l.ctx, cmd)
	if err != nil {
		l.Error("Failed to create order", logx.Field("error", err))
		return &trade.CreateOrderResponse{
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
	return &trade.CreateOrderResponse{
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
			CreatedAt:      orderDTO.CreatedAt.Unix(),
			UpdatedAt:      orderDTO.UpdatedAt.Unix(),
		},
	}, nil
}
