package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// Product 商品实体
type Product struct {
	id          valueobject.ProductID
	name        valueobject.ProductName
	description valueobject.ProductDescription
	sku         valueobject.ProductSKU
	price       valueobject.Price
	stock       int32
	status      valueobject.ProductStatus
	categoryID  int64
	creatorID   uservo.UserID
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProduct 创建新商品
func NewProduct(
	name valueobject.ProductName,
	description valueobject.ProductDescription,
	sku valueobject.ProductSKU,
	price valueobject.Price,
	stock int32,
	categoryID int64,
	creatorID uservo.UserID,
) (*Product, error) {
	// 验证必填参数
	if name.Value() == "" {
		return nil, errors.New("商品名称不能为空")
	}
	if sku.Value() == "" {
		return nil, errors.New("SKU不能为空")
	}
	if price.Value() < 0 {
		return nil, errors.New("价格不能为负数")
	}
	if stock < 0 {
		return nil, errors.New("库存不能为负数")
	}
	if categoryID <= 0 {
		return nil, errors.New("无效的分类ID")
	}
	if creatorID.Value() <= 0 {
		return nil, errors.New("无效的创建者ID")
	}

	// 默认状态为草稿
	status, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusDraft))

	now := time.Now()

	return &Product{
		name:        name,
		description: description,
		sku:         sku,
		price:       price,
		stock:       stock,
		status:      status,
		categoryID:  categoryID,
		creatorID:   creatorID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructProduct 从存储中重建商品实体
func ReconstructProduct(
	id valueobject.ProductID,
	name valueobject.ProductName,
	description valueobject.ProductDescription,
	sku valueobject.ProductSKU,
	price valueobject.Price,
	stock int32,
	status valueobject.ProductStatus,
	categoryID int64,
	creatorID uservo.UserID,
	createdAt time.Time,
	updatedAt time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		sku:         sku,
		price:       price,
		stock:       stock,
		status:      status,
		categoryID:  categoryID,
		creatorID:   creatorID,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID 获取商品ID
func (p *Product) ID() valueobject.ProductID {
	return p.id
}

// SetID 设置商品ID
func (p *Product) SetID(id valueobject.ProductID) {
	p.id = id
}

// Name 获取商品名称
func (p *Product) Name() valueobject.ProductName {
	return p.name
}

// SetName 更新商品名称
func (p *Product) SetName(name valueobject.ProductName) {
	p.name = name
	p.updatedAt = time.Now()
}

// Description 获取商品描述
func (p *Product) Description() valueobject.ProductDescription {
	return p.description
}

// SetDescription 更新商品描述
func (p *Product) SetDescription(description valueobject.ProductDescription) {
	p.description = description
	p.updatedAt = time.Now()
}

// SKU 获取商品SKU
func (p *Product) SKU() valueobject.ProductSKU {
	return p.sku
}

// SetSKU 更新商品SKU
func (p *Product) SetSKU(sku valueobject.ProductSKU) {
	p.sku = sku
	p.updatedAt = time.Now()
}

// Price 获取商品价格
func (p *Product) Price() valueobject.Price {
	return p.price
}

// SetPrice 更新商品价格
func (p *Product) SetPrice(price valueobject.Price) {
	p.price = price
	p.updatedAt = time.Now()
}

// Stock 获取商品库存
func (p *Product) Stock() int32 {
	return p.stock
}

// SetStock 更新商品库存
func (p *Product) SetStock(stock int32) {
	p.stock = stock
	p.updatedAt = time.Now()
}

// DecrementStock 减少库存
func (p *Product) DecrementStock(quantity int32) error {
	if quantity <= 0 {
		return errors.New("减少的库存量必须大于0")
	}

	if p.stock < quantity {
		return errors.New("库存不足")
	}

	p.stock -= quantity
	p.updatedAt = time.Now()

	// 如果库存为0，自动更新状态为售罄
	if p.stock == 0 {
		soldOutStatus, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusSoldOut))
		p.status = soldOutStatus
	}

	return nil
}

// IncrementStock 增加库存
func (p *Product) IncrementStock(quantity int32) error {
	if quantity <= 0 {
		return errors.New("增加的库存量必须大于0")
	}

	p.stock += quantity
	p.updatedAt = time.Now()

	// 如果之前是售罄状态，并且现在有库存了，自动恢复为上架状态
	if p.status == valueobject.ProductStatusSoldOut && p.stock > 0 {
		activeStatus, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusActive))
		p.status = activeStatus
	}

	return nil
}

// Status 获取商品状态
func (p *Product) Status() valueobject.ProductStatus {
	return p.status
}

// SetStatus 更新商品状态
func (p *Product) SetStatus(status valueobject.ProductStatus) error {
	// 如果当前商品已删除，则不允许更改状态
	if p.status == valueobject.ProductStatusDeleted {
		return errors.New("已删除的商品不能更改状态")
	}

	// 如果想要设置为上架状态，但库存为0，则不允许
	if status == valueobject.ProductStatusActive && p.stock == 0 {
		return errors.New("库存为0的商品不能上架")
	}

	p.status = status
	p.updatedAt = time.Now()
	return nil
}

// CategoryID 获取分类ID
func (p *Product) CategoryID() int64 {
	return p.categoryID
}

// SetCategoryID 更新分类ID
func (p *Product) SetCategoryID(categoryID int64) error {
	if categoryID <= 0 {
		return errors.New("无效的分类ID")
	}

	p.categoryID = categoryID
	p.updatedAt = time.Now()
	return nil
}

// CreatorID 获取创建者ID
func (p *Product) CreatorID() uservo.UserID {
	return p.creatorID
}

// CreatedAt 获取创建时间
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt 获取更新时间
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}

// IsActive 商品是否处于上架状态
func (p *Product) IsActive() bool {
	return p.status == valueobject.ProductStatusActive
}

// IsAvailable 商品是否可售卖（上架状态且有库存）
func (p *Product) IsAvailable() bool {
	return p.IsActive() && p.stock > 0
}

// Delete 删除商品
func (p *Product) Delete() {
	deletedStatus, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusDeleted))
	p.status = deletedStatus
	p.updatedAt = time.Now()
}

// Publish 发布商品（将草稿改为上架状态）
func (p *Product) Publish() error {
	if p.stock <= 0 {
		return errors.New("库存为0的商品不能上架")
	}

	activeStatus, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusActive))
	p.status = activeStatus
	p.updatedAt = time.Now()
	return nil
}

// Unpublish 下架商品
func (p *Product) Unpublish() error {
	inactiveStatus, _ := valueobject.NewProductStatus(int32(valueobject.ProductStatusInactive))
	p.status = inactiveStatus
	p.updatedAt = time.Now()
	return nil
}
