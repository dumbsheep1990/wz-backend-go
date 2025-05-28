package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// CartItem 购物车项
type CartItem struct {
	productID valueobject.ProductID
	name      string
	price     valueobject.Money
	quantity  valueobject.Quantity
	addedAt   time.Time
}

// NewCartItem 创建一个新的购物车项
func NewCartItem(
	productID valueobject.ProductID,
	name string,
	price valueobject.Money,
	quantity valueobject.Quantity,
) (*CartItem, error) {
	if productID.IsEmpty() {
		return nil, errors.New("商品ID不能为空")
	}
	if name == "" {
		return nil, errors.New("商品名称不能为空")
	}

	return &CartItem{
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
		addedAt:   time.Now(),
	}, nil
}

// ProductID 获取商品ID
func (i *CartItem) ProductID() valueobject.ProductID {
	return i.productID
}

// Name 获取商品名称
func (i *CartItem) Name() string {
	return i.name
}

// Price 获取商品单价
func (i *CartItem) Price() valueobject.Money {
	return i.price
}

// Quantity 获取商品数量
func (i *CartItem) Quantity() valueobject.Quantity {
	return i.quantity
}

// UpdateQuantity 更新数量
func (i *CartItem) UpdateQuantity(quantity valueobject.Quantity) {
	i.quantity = quantity
}

// AddedAt 获取添加时间
func (i *CartItem) AddedAt() time.Time {
	return i.addedAt
}

// Subtotal 计算小计金额
func (i *CartItem) Subtotal() (valueobject.Money, error) {
	return i.price.Multiply(i.quantity.Value())
}

// Cart 购物车实体
type Cart struct {
	id        valueobject.CartID
	userID    valueobject.UserID
	items     map[string]*CartItem // 以商品ID为键
	createdAt time.Time
	updatedAt time.Time
	
	domainEvents []event.DomainEvent
}

// NewCart 创建一个新购物车
func NewCart(userID valueobject.UserID) (*Cart, error) {
	if userID.IsEmpty() {
		return nil, errors.New("用户ID不能为空")
	}
	
	now := time.Now()
	cart := &Cart{
		id:           valueobject.NewCartID(uuid.New().String()),
		userID:       userID,
		items:        make(map[string]*CartItem),
		createdAt:    now,
		updatedAt:    now,
		domainEvents: []event.DomainEvent{},
	}
	
	// 添加购物车创建事件
	cart.addDomainEvent(NewCartCreatedEvent(cart))
	
	return cart, nil
}

// ID 获取购物车ID
func (c *Cart) ID() valueobject.CartID {
	return c.id
}

// UserID 获取用户ID
func (c *Cart) UserID() valueobject.UserID {
	return c.userID
}

// Items 获取所有购物车项
func (c *Cart) Items() []*CartItem {
	items := make([]*CartItem, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}
	return items
}

// ItemCount 获取购物车项数量
func (c *Cart) ItemCount() int {
	return len(c.items)
}

// TotalQuantity 获取总数量
func (c *Cart) TotalQuantity() int {
	total := 0
	for _, item := range c.items {
		total += item.quantity.Value()
	}
	return total
}

// TotalAmount 获取总金额
func (c *Cart) TotalAmount() (valueobject.Money, error) {
	// 如果购物车为空，返回零金额
	if len(c.items) == 0 {
		return valueobject.NewMoney(0, "CNY")
	}
	
	// 找到第一个商品的货币类型
	var currency string
	for _, item := range c.items {
		currency = item.price.Currency()
		break
	}
	
	// 创建初始金额
	totalAmount, err := valueobject.NewMoney(0, currency)
	if err != nil {
		return valueobject.Money{}, err
	}
	
	// 计算总金额
	for _, item := range c.items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return valueobject.Money{}, err
		}
		
		totalAmount, err = totalAmount.Add(subtotal)
		if err != nil {
			return valueobject.Money{}, err
		}
	}
	
	return totalAmount, nil
}

// CreatedAt 获取创建时间
func (c *Cart) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt 获取更新时间
func (c *Cart) UpdatedAt() time.Time {
	return c.updatedAt
}

// AddItem 添加商品到购物车
func (c *Cart) AddItem(item *CartItem) error {
	productIDStr := item.productID.String()
	
	// 检查商品是否已存在于购物车中
	if existingItem, exists := c.items[productIDStr]; exists {
		// 合并数量
		newQuantity, err := valueobject.NewQuantity(existingItem.quantity.Value() + item.quantity.Value())
		if err != nil {
			return err
		}
		existingItem.UpdateQuantity(newQuantity)
	} else {
		// 添加新商品
		c.items[productIDStr] = item
	}
	
	c.updatedAt = time.Now()
	
	// 添加商品添加事件
	c.addDomainEvent(NewCartItemAddedEvent(c, item))
	
	return nil
}

// UpdateItemQuantity 更新购物车中商品的数量
func (c *Cart) UpdateItemQuantity(productID valueobject.ProductID, quantity valueobject.Quantity) error {
	productIDStr := productID.String()
	
	// 检查商品是否存在于购物车中
	item, exists := c.items[productIDStr]
	if !exists {
		return errors.New("商品不存在于购物车中")
	}
	
	// 更新数量
	item.UpdateQuantity(quantity)
	c.updatedAt = time.Now()
	
	// 添加数量更新事件
	c.addDomainEvent(NewCartItemQuantityUpdatedEvent(c, item))
	
	return nil
}

// RemoveItem 从购物车中移除商品
func (c *Cart) RemoveItem(productID valueobject.ProductID) error {
	productIDStr := productID.String()
	
	// 检查商品是否存在于购物车中
	item, exists := c.items[productIDStr]
	if !exists {
		return errors.New("商品不存在于购物车中")
	}
	
	// 移除商品
	delete(c.items, productIDStr)
	c.updatedAt = time.Now()
	
	// 添加商品移除事件
	c.addDomainEvent(NewCartItemRemovedEvent(c, item))
	
	return nil
}

// Clear 清空购物车
func (c *Cart) Clear() {
	// 清空商品
	c.items = make(map[string]*CartItem)
	c.updatedAt = time.Now()
	
	// 添加购物车清空事件
	c.addDomainEvent(NewCartClearedEvent(c))
}

// IsEmpty 检查购物车是否为空
func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}

// ContainsItem 检查购物车是否包含特定商品
func (c *Cart) ContainsItem(productID valueobject.ProductID) bool {
	_, exists := c.items[productID.String()]
	return exists
}

// GetItem 获取购物车中的特定商品
func (c *Cart) GetItem(productID valueobject.ProductID) (*CartItem, bool) {
	item, exists := c.items[productID.String()]
	return item, exists
}

// ToOrder 将购物车转换为订单
func (c *Cart) ToOrder(shippingAddress valueobject.Address) (*Order, error) {
	if c.IsEmpty() {
		return nil, errors.New("购物车为空，无法创建订单")
	}
	
	// 将购物车项转换为订单项
	orderItems := make([]*OrderItem, 0, len(c.items))
	for _, cartItem := range c.items {
		orderItem, err := NewOrderItem(
			cartItem.productID,
			cartItem.name,
			cartItem.price,
			cartItem.quantity,
		)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	
	// 创建订单
	order, err := NewOrder(c.userID, orderItems, shippingAddress)
	if err != nil {
		return nil, err
	}
	
	// 添加购物车转换为订单事件
	c.addDomainEvent(NewCartConvertedToOrderEvent(c, order))
	
	return order, nil
}

// 添加领域事件
func (c *Cart) addDomainEvent(event event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, event)
}

// GetDomainEvents 获取所有领域事件
func (c *Cart) GetDomainEvents() []event.DomainEvent {
	return c.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (c *Cart) ClearDomainEvents() {
	c.domainEvents = []event.DomainEvent{}
}
