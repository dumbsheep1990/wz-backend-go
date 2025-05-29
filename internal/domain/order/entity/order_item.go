package entity

import (
	"errors"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
	productvo "wz-backend-go/internal/domain/product/valueobject"
)

// OrderItem 订单项实体
type OrderItem struct {
	id          ordervo.OrderItemID
	productID   productvo.ProductID
	productName string
	productSKU  string
	quantity    int32
	unitPrice   ordervo.Money
	totalPrice  ordervo.Money
	attributes  map[string]string // 商品属性，如颜色、尺寸等
}

// NewOrderItem 创建订单项
func NewOrderItem(
	id ordervo.OrderItemID,
	productID productvo.ProductID,
	productName string,
	productSKU string,
	quantity int32,
	unitPrice ordervo.Money,
	attributes map[string]string,
) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, errors.New("商品数量必须大于0")
	}

	// 计算总价
	totalPrice, err := unitPrice.Multiply(int(quantity))
	if err != nil {
		return nil, err
	}

	return &OrderItem{
		id:          id,
		productID:   productID,
		productName: productName,
		productSKU:  productSKU,
		quantity:    quantity,
		unitPrice:   unitPrice,
		totalPrice:  totalPrice,
		attributes:  attributes,
	}, nil
}

// ID 获取订单项ID
func (item *OrderItem) ID() ordervo.OrderItemID {
	return item.id
}

// ProductID 获取商品ID
func (item *OrderItem) ProductID() productvo.ProductID {
	return item.productID
}

// ProductName 获取商品名称
func (item *OrderItem) ProductName() string {
	return item.productName
}

// ProductSKU 获取商品SKU
func (item *OrderItem) ProductSKU() string {
	return item.productSKU
}

// Quantity 获取商品数量
func (item *OrderItem) Quantity() int32 {
	return item.quantity
}

// UnitPrice 获取单价
func (item *OrderItem) UnitPrice() ordervo.Money {
	return item.unitPrice
}

// TotalPrice 获取总价
func (item *OrderItem) TotalPrice() ordervo.Money {
	return item.totalPrice
}

// Attributes 获取商品属性
func (item *OrderItem) Attributes() map[string]string {
	return item.attributes
}

// UpdateQuantity 更新商品数量
func (item *OrderItem) UpdateQuantity(quantity int32) error {
	if quantity <= 0 {
		return errors.New("商品数量必须大于0")
	}

	item.quantity = quantity

	// 重新计算总价
	totalPrice, err := item.unitPrice.Multiply(int(quantity))
	if err != nil {
		return err
	}

	item.totalPrice = totalPrice
	return nil
}

// AddAttribute 添加商品属性
func (item *OrderItem) AddAttribute(key, value string) {
	if item.attributes == nil {
		item.attributes = make(map[string]string)
	}
	item.attributes[key] = value
}

// RemoveAttribute 移除商品属性
func (item *OrderItem) RemoveAttribute(key string) {
	if item.attributes == nil {
		return
	}
	delete(item.attributes, key)
}

// HasSameProduct 判断是否为同一商品
func (item *OrderItem) HasSameProduct(productID productvo.ProductID) bool {
	return item.productID.Value() == productID.Value()
}

// UpdateUnitPrice 更新单价
func (item *OrderItem) UpdateUnitPrice(unitPrice ordervo.Money) error {
	item.unitPrice = unitPrice

	// 重新计算总价
	totalPrice, err := unitPrice.Multiply(int(item.quantity))
	if err != nil {
		return err
	}

	item.totalPrice = totalPrice
	return nil
}
