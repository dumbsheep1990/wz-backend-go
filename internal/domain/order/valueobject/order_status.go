package valueobject

import (
	"errors"
)

// OrderStatus 订单状态值对象
type OrderStatus int32

// 订单状态常量
const (
	OrderStatusCreated   OrderStatus = 0 // 已创建
	OrderStatusPending   OrderStatus = 1 // 待支付
	OrderStatusPaid      OrderStatus = 2 // 已支付
	OrderStatusShipped   OrderStatus = 3 // 已发货
	OrderStatusDelivered OrderStatus = 4 // 已送达
	OrderStatusCompleted OrderStatus = 5 // 已完成
	OrderStatusCancelled OrderStatus = 6 // 已取消
	OrderStatusRefunding OrderStatus = 7 // 退款中
	OrderStatusRefunded  OrderStatus = 8 // 已退款
)

// NewOrderStatus 创建订单状态值对象
func NewOrderStatus(status int32) (OrderStatus, error) {
	switch status {
	case int32(OrderStatusCreated), int32(OrderStatusPending), int32(OrderStatusPaid),
		int32(OrderStatusShipped), int32(OrderStatusDelivered), int32(OrderStatusCompleted),
		int32(OrderStatusCancelled), int32(OrderStatusRefunding), int32(OrderStatusRefunded):
		return OrderStatus(status), nil
	default:
		return 0, errors.New("无效的订单状态")
	}
}

// Value 获取状态值
func (s OrderStatus) Value() int32 {
	return int32(s)
}

// String 状态的字符串表示
func (s OrderStatus) String() string {
	switch s {
	case OrderStatusCreated:
		return "已创建"
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusDelivered:
		return "已送达"
	case OrderStatusCompleted:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	case OrderStatusRefunding:
		return "退款中"
	case OrderStatusRefunded:
		return "已退款"
	default:
		return "未知状态"
	}
}

// CanTransitionTo 检查当前状态是否可以转换到目标状态
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	switch s {
	case OrderStatusCreated:
		return target == OrderStatusPending || target == OrderStatusCancelled
	case OrderStatusPending:
		return target == OrderStatusPaid || target == OrderStatusCancelled
	case OrderStatusPaid:
		return target == OrderStatusShipped || target == OrderStatusRefunding
	case OrderStatusShipped:
		return target == OrderStatusDelivered || target == OrderStatusRefunding
	case OrderStatusDelivered:
		return target == OrderStatusCompleted || target == OrderStatusRefunding
	case OrderStatusRefunding:
		return target == OrderStatusRefunded
	case OrderStatusCompleted, OrderStatusCancelled, OrderStatusRefunded:
		return false // 终态，不能再转换
	default:
		return false
	}
}

// IsActive 判断订单状态是否处于活跃状态（非终态）
func (s OrderStatus) IsActive() bool {
	return s != OrderStatusCompleted && s != OrderStatusCancelled && s != OrderStatusRefunded
}

// IsTerminal 判断订单状态是否为终态
func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusCompleted || s == OrderStatusCancelled || s == OrderStatusRefunded
}

// CanCancel 判断当前状态是否可以取消订单
func (s OrderStatus) CanCancel() bool {
	return s == OrderStatusCreated || s == OrderStatusPending
}

// CanPay 判断当前状态是否可以支付
func (s OrderStatus) CanPay() bool {
	return s == OrderStatusPending
}

// CanRefund 判断当前状态是否可以申请退款
func (s OrderStatus) CanRefund() bool {
	return s == OrderStatusPaid || s == OrderStatusShipped || s == OrderStatusDelivered
}

// CanShip 判断当前状态是否可以发货
func (s OrderStatus) CanShip() bool {
	return s == OrderStatusPaid
}

// CanDeliver 判断当前状态是否可以送达
func (s OrderStatus) CanDeliver() bool {
	return s == OrderStatusShipped
}

// CanComplete 判断当前状态是否可以完成
func (s OrderStatus) CanComplete() bool {
	return s == OrderStatusDelivered
}

// IsCreated 判断是否为已创建状态
func (s OrderStatus) IsCreated() bool {
	return s == OrderStatusCreated
}

// IsPending 判断是否为待支付状态
func (s OrderStatus) IsPending() bool {
	return s == OrderStatusPending
}

// IsSubmitted 判断是否为已提交状态（待支付）
func (s OrderStatus) IsSubmitted() bool {
	return s == OrderStatusPending
}

// IsPaid 判断是否为已支付状态
func (s OrderStatus) IsPaid() bool {
	return s == OrderStatusPaid
}

// IsShipped 判断是否为已发货状态
func (s OrderStatus) IsShipped() bool {
	return s == OrderStatusShipped
}

// IsDelivered 判断是否为已送达状态
func (s OrderStatus) IsDelivered() bool {
	return s == OrderStatusDelivered
}

// IsCompleted 判断是否为已完成状态
func (s OrderStatus) IsCompleted() bool {
	return s == OrderStatusCompleted
}

// IsCancelled 判断是否为已取消状态
func (s OrderStatus) IsCancelled() bool {
	return s == OrderStatusCancelled
}

// IsRefunding 判断是否为退款中状态
func (s OrderStatus) IsRefunding() bool {
	return s == OrderStatusRefunding
}

// IsRefunded 判断是否为已退款状态
func (s OrderStatus) IsRefunded() bool {
	return s == OrderStatusRefunded
}
