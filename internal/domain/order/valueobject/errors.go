package valueobject

import "errors"

// 订单相关错误定义
var (
	// 订单创建错误
	ErrOrderItemsEmpty     = errors.New("订单项不能为空")
	ErrInvalidOrderAmount  = errors.New("无效的订单金额")
	ErrInvalidAddress      = errors.New("无效的地址信息")
	ErrInvalidOrderStatus  = errors.New("无效的订单状态")

	// 订单支付错误
	ErrOrderAlreadyPaid      = errors.New("订单已支付")
	ErrOrderNotPaid          = errors.New("订单未支付")
	ErrInvalidPaymentMethod  = errors.New("无效的支付方式")
	ErrPaymentMethodRequired = errors.New("支付方式不能为空")

	// 订单发货错误
	ErrOrderAlreadyShipped    = errors.New("订单已发货")
	ErrOrderNotShipped        = errors.New("订单未发货")
	ErrInvalidTrackingNumber  = errors.New("无效的物流单号")
	ErrTrackingNumberRequired = errors.New("物流单号不能为空")

	// 订单状态错误
	ErrOrderCancelled     = errors.New("订单已取消")
	ErrOrderCompleted     = errors.New("订单已完成")
	ErrOrderRefunded      = errors.New("订单已退款")
	ErrCannotCancelOrder  = errors.New("当前状态不允许取消订单")
	ErrCannotRefundOrder  = errors.New("当前状态不允许退款")
	ErrCannotModifyOrder  = errors.New("当前状态不允许修改订单")

	// 订单项错误
	ErrOrderItemNotFound    = errors.New("订单项不存在")
	ErrInvalidQuantity      = errors.New("无效的商品数量")
	ErrInvalidPrice         = errors.New("无效的商品价格")
	ErrProductIDRequired    = errors.New("商品ID不能为空")
	ErrProductNameRequired  = errors.New("商品名称不能为空")

	// 折扣错误
	ErrDiscountNotFound      = errors.New("折扣不存在")
	ErrInvalidDiscountType   = errors.New("无效的折扣类型")
	ErrInvalidDiscountValue  = errors.New("无效的折扣值")
	ErrDiscountExpired       = errors.New("折扣已过期")
	ErrDiscountNotApplicable = errors.New("折扣不适用于当前订单")

	// 配送错误
	ErrInvalidShippingMethod = errors.New("无效的配送方式")
	ErrInvalidShippingFee    = errors.New("无效的配送费用")

	// 权限错误
	ErrUnauthorizedAccess = errors.New("无权访问该订单")
	ErrInsufficientPermission = errors.New("权限不足")

	// 业务规则错误
	ErrRefundDeadlineExceeded = errors.New("已超过退款期限")
	ErrOrderTimeout           = errors.New("订单已超时")
	ErrInvalidOperationTime   = errors.New("操作时间无效")
) 