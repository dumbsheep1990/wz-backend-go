package logic

import (
	"context"
	"time"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *trade.ListOrdersRequest) (*trade.ListOrdersResponse, error) {
	// 构建查询参数
	query := dto.ListOrdersQuery{
		CustomerID: in.CustomerId,
		Page:       int(in.Page),
		PageSize:   int(in.PageSize),
		SortBy:     in.SortBy,
		SortOrder:  in.SortOrder,
	}

	// 处理状态过滤
	if len(in.Status) > 0 {
		query.Status = in.Status
	}

	// 处理日期过滤
	if in.StartDate > 0 {
		startDate := time.Unix(in.StartDate, 0)
		query.StartDate = startDate
	}

	if in.EndDate > 0 {
		endDate := time.Unix(in.EndDate, 0)
		query.EndDate = endDate
	}

	// 调用应用服务
	result, err := l.svcCtx.OrderApplicationService.ListOrders(l.ctx, query)
	if err != nil {
		l.Error("Failed to list orders", logx.Field("error", err))
		return &trade.ListOrdersResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 转换订单列表为响应格式
	orders := make([]*trade.Order, 0, len(result.Orders))
	for _, orderDTO := range result.Orders {
		// 转换订单项
		var items []*trade.OrderItem
		for _, item := range orderDTO.Items {
			items = append(items, &trade.OrderItem{
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
		var discounts []*trade.OrderDiscount
		for _, discount := range orderDTO.Discounts {
			discounts = append(discounts, &trade.OrderDiscount{
				Id:           discount.ID,
				DiscountType: discount.DiscountType,
				DiscountName: discount.DiscountName,
				Amount: &trade.Money{
					Amount:   discount.Amount.Amount,
					Currency: discount.Amount.Currency,
				},
			})
		}

		// 构建订单对象
		order := &trade.Order{
			Id:           orderDTO.ID,
			OrderNumber:  orderDTO.OrderNumber,
			CustomerId:   orderDTO.CustomerID,
			Status:       orderDTO.Status,
			StatusCode:   orderDTO.StatusCode,
			Items:        items,
			Discounts:    discounts,
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
		}

		orders = append(orders, order)
	}

	// 返回响应
	return &trade.ListOrdersResponse{
		Success: true,
		Total:   result.Total,
		Orders:  orders,
	}, nil
}
