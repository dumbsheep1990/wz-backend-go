package entity

import (
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

const (
	EventTypeCartCreated            = "cart.created"
	EventTypeCartItemAdded          = "cart.item_added"
	EventTypeCartItemQuantityUpdated = "cart.item_quantity_updated"
	EventTypeCartItemRemoved        = "cart.item_removed"
	EventTypeCartCleared            = "cart.cleared"
	EventTypeCartConvertedToOrder   = "cart.converted_to_order"
)

// CartCreatedEvent 购物车创建事件
type CartCreatedEvent struct {
	event.BaseDomainEvent
}

// NewCartCreatedEvent 创建新的购物车创建事件
func NewCartCreatedEvent(cart *Cart) CartCreatedEvent {
	return CartCreatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartCreated,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"id":        cart.ID().String(),
				"userID":    cart.UserID().String(),
				"createdAt": cart.CreatedAt(),
			},
		),
	}
}

// CartItemAddedEvent 购物车添加商品事件
type CartItemAddedEvent struct {
	event.BaseDomainEvent
}

// NewCartItemAddedEvent 创建新的购物车添加商品事件
func NewCartItemAddedEvent(cart *Cart, item *CartItem) CartItemAddedEvent {
	return CartItemAddedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartItemAdded,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"cartID":    cart.ID().String(),
				"userID":    cart.UserID().String(),
				"productID": item.ProductID().String(),
				"name":      item.Name(),
				"price":     item.Price().Amount(),
				"currency":  item.Price().Currency(),
				"quantity":  item.Quantity().Value(),
				"addedAt":   item.AddedAt(),
			},
		),
	}
}

// CartItemQuantityUpdatedEvent 购物车商品数量更新事件
type CartItemQuantityUpdatedEvent struct {
	event.BaseDomainEvent
}

// NewCartItemQuantityUpdatedEvent 创建新的购物车商品数量更新事件
func NewCartItemQuantityUpdatedEvent(cart *Cart, item *CartItem) CartItemQuantityUpdatedEvent {
	return CartItemQuantityUpdatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartItemQuantityUpdated,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"cartID":    cart.ID().String(),
				"userID":    cart.UserID().String(),
				"productID": item.ProductID().String(),
				"quantity":  item.Quantity().Value(),
				"updatedAt": cart.UpdatedAt(),
			},
		),
	}
}

// CartItemRemovedEvent 购物车移除商品事件
type CartItemRemovedEvent struct {
	event.BaseDomainEvent
}

// NewCartItemRemovedEvent 创建新的购物车移除商品事件
func NewCartItemRemovedEvent(cart *Cart, item *CartItem) CartItemRemovedEvent {
	return CartItemRemovedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartItemRemoved,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"cartID":    cart.ID().String(),
				"userID":    cart.UserID().String(),
				"productID": item.ProductID().String(),
				"name":      item.Name(),
				"removedAt": cart.UpdatedAt(),
			},
		),
	}
}

// CartClearedEvent 购物车清空事件
type CartClearedEvent struct {
	event.BaseDomainEvent
}

// NewCartClearedEvent 创建新的购物车清空事件
func NewCartClearedEvent(cart *Cart) CartClearedEvent {
	return CartClearedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartCleared,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"cartID":    cart.ID().String(),
				"userID":    cart.UserID().String(),
				"clearedAt": cart.UpdatedAt(),
			},
		),
	}
}

// CartConvertedToOrderEvent 购物车转换为订单事件
type CartConvertedToOrderEvent struct {
	event.BaseDomainEvent
}

// NewCartConvertedToOrderEvent 创建新的购物车转换为订单事件
func NewCartConvertedToOrderEvent(cart *Cart, order *Order) CartConvertedToOrderEvent {
	return CartConvertedToOrderEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCartConvertedToOrder,
			cart.ID().String(),
			"Cart",
			map[string]interface{}{
				"cartID":    cart.ID().String(),
				"userID":    cart.UserID().String(),
				"orderID":   order.ID().String(),
				"itemCount": cart.ItemCount(),
				"convertedAt": cart.UpdatedAt(),
			},
		),
	}
}
