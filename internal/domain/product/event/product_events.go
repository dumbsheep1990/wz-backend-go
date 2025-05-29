package event

import (
	"time"

	"github.com/google/uuid"

	"wz-backend-go/internal/domain/product/valueobject"
)

// 商品相关事件类型常量
const (
	ProductCreatedEventType      = "product.created"
	ProductUpdatedEventType      = "product.updated"
	ProductPublishedEventType    = "product.published"
	ProductUnpublishedEventType  = "product.unpublished"
	ProductDeletedEventType      = "product.deleted"
	ProductStockChangedEventType = "product.stock_changed"
)

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	eventID     string
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

func NewBaseDomainEvent(eventType string, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     uuid.New().String(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
	}
}

func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseDomainEvent) OccurredTime() time.Time {
	return e.occurredAt
}

// ProductCreatedEvent 商品创建事件
type ProductCreatedEvent struct {
	BaseDomainEvent
	ProductID   int64
	ProductName string
	SKU         string
	Price       int64
	CreatorID   int64
}

func NewProductCreatedEvent(
	productID valueobject.ProductID,
	productName valueobject.ProductName,
	sku valueobject.ProductSKU,
	price valueobject.Price,
	creatorID int64,
) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ProductCreatedEventType,
			productID.String(),
		),
		ProductID:   productID.Value(),
		ProductName: productName.Value(),
		SKU:         sku.Value(),
		Price:       price.Value(),
		CreatorID:   creatorID,
	}
}

// ProductPublishedEvent 商品发布事件
type ProductPublishedEvent struct {
	BaseDomainEvent
	ProductID int64
	Stock     int32
}

func NewProductPublishedEvent(
	productID valueobject.ProductID,
	stock int32,
) *ProductPublishedEvent {
	return &ProductPublishedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ProductPublishedEventType,
			productID.String(),
		),
		ProductID: productID.Value(),
		Stock:     stock,
	}
}

// ProductUnpublishedEvent 商品下架事件
type ProductUnpublishedEvent struct {
	BaseDomainEvent
	ProductID int64
}

func NewProductUnpublishedEvent(
	productID valueobject.ProductID,
) *ProductUnpublishedEvent {
	return &ProductUnpublishedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ProductUnpublishedEventType,
			productID.String(),
		),
		ProductID: productID.Value(),
	}
}

// ProductDeletedEvent 商品删除事件
type ProductDeletedEvent struct {
	BaseDomainEvent
	ProductID int64
}

func NewProductDeletedEvent(
	productID valueobject.ProductID,
) *ProductDeletedEvent {
	return &ProductDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ProductDeletedEventType,
			productID.String(),
		),
		ProductID: productID.Value(),
	}
}

// ProductStockChangedEvent 商品库存变更事件
type ProductStockChangedEvent struct {
	BaseDomainEvent
	ProductID     int64
	PreviousStock int32
	CurrentStock  int32
	ChangeAmount  int32
	Reason        string
}

func NewProductStockChangedEvent(
	productID valueobject.ProductID,
	previousStock int32,
	currentStock int32,
	changeAmount int32,
	reason string,
) *ProductStockChangedEvent {
	return &ProductStockChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ProductStockChangedEventType,
			productID.String(),
		),
		ProductID:     productID.Value(),
		PreviousStock: previousStock,
		CurrentStock:  currentStock,
		ChangeAmount:  changeAmount,
		Reason:        reason,
	}
}
