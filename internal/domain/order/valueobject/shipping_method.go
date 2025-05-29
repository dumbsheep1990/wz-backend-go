package valueobject

import (
	"errors"
)

// ShippingMethod 配送方式值对象
type ShippingMethod int32

// 配送方式常量
const (
	ShippingMethodUnknown  ShippingMethod = 0 // 未知
	ShippingMethodExpress  ShippingMethod = 1 // 快递
	ShippingMethodSameDay  ShippingMethod = 2 // 当日达
	ShippingMethodNextDay  ShippingMethod = 3 // 次日达
	ShippingMethodPickup   ShippingMethod = 4 // 自提
	ShippingMethodStandard ShippingMethod = 5 // 标准配送
)

// NewShippingMethod 创建配送方式值对象
func NewShippingMethod(method int32) (ShippingMethod, error) {
	switch method {
	case int32(ShippingMethodUnknown), int32(ShippingMethodExpress), int32(ShippingMethodSameDay),
		int32(ShippingMethodNextDay), int32(ShippingMethodPickup), int32(ShippingMethodStandard):
		return ShippingMethod(method), nil
	default:
		return 0, errors.New("无效的配送方式")
	}
}

// Value 获取配送方式值
func (s ShippingMethod) Value() int32 {
	return int32(s)
}

// String 配送方式的字符串表示
func (s ShippingMethod) String() string {
	switch s {
	case ShippingMethodExpress:
		return "快递"
	case ShippingMethodSameDay:
		return "当日达"
	case ShippingMethodNextDay:
		return "次日达"
	case ShippingMethodPickup:
		return "自提"
	case ShippingMethodStandard:
		return "标准配送"
	default:
		return "未知配送方式"
	}
}

// IsExpressDelivery 判断是否为快速配送
func (s ShippingMethod) IsExpressDelivery() bool {
	return s == ShippingMethodExpress || s == ShippingMethodSameDay || s == ShippingMethodNextDay
}

// RequiresShippingAddress 判断是否需要配送地址
func (s ShippingMethod) RequiresShippingAddress() bool {
	return s != ShippingMethodPickup
}

// CanTrack 判断是否可以物流跟踪
func (s ShippingMethod) CanTrack() bool {
	return s != ShippingMethodPickup && s != ShippingMethodUnknown
}

// EstimatedDeliveryDays 估算配送天数
func (s ShippingMethod) EstimatedDeliveryDays() int {
	switch s {
	case ShippingMethodSameDay:
		return 0
	case ShippingMethodNextDay:
		return 1
	case ShippingMethodExpress:
		return 1
	case ShippingMethodStandard:
		return 3
	case ShippingMethodPickup:
		return 0
	default:
		return 7
	}
}
