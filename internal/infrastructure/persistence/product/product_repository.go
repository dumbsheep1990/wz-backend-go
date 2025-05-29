package product

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/product/entity"
	"wz-backend-go/internal/domain/product/repository"
	"wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// ProductPO 商品持久化对象
type ProductPO struct {
	ID          int64        `gorm:"primaryKey;column:id"`
	Name        string       `gorm:"column:name;size:100;not null"`
	Description string       `gorm:"column:description;size:2000"`
	SKU         string       `gorm:"uniqueIndex;column:sku;size:20;not null"`
	Price       int64        `gorm:"column:price;not null"`
	Stock       int32        `gorm:"column:stock;not null"`
	Status      int32        `gorm:"column:status;not null;default:0"`
	CategoryID  int64        `gorm:"index;column:category_id;not null"`
	CreatorID   int64        `gorm:"index;column:creator_id;not null"`
	CreatedAt   time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;not null"`
	DeletedAt   sql.NullTime `gorm:"column:deleted_at"`
}

// TableName 表名
func (ProductPO) TableName() string {
	return "products"
}

// ProductRepository 商品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Save 保存商品
func (r *ProductRepository) Save(product *entity.Product) error {
	// 如果商品ID为0，则为新商品，需要插入
	if product.ID().Value() == 0 {
		po := &ProductPO{
			Name:        product.Name().Value(),
			Description: product.Description().Value(),
			SKU:         product.SKU().Value(),
			Price:       product.Price().Value(),
			Stock:       product.Stock(),
			Status:      product.Status().Value(),
			CategoryID:  product.CategoryID(),
			CreatorID:   product.CreatorID().Value(),
			CreatedAt:   product.CreatedAt(),
			UpdatedAt:   product.UpdatedAt(),
		}

		if err := r.db.Create(po).Error; err != nil {
			return err
		}

		// 设置商品ID
		productID := valueobject.NewProductID(po.ID)
		product.SetID(productID)

		return nil
	}

	// 否则为更新商品
	po := &ProductPO{
		ID:          product.ID().Value(),
		Name:        product.Name().Value(),
		Description: product.Description().Value(),
		SKU:         product.SKU().Value(),
		Price:       product.Price().Value(),
		Stock:       product.Stock(),
		Status:      product.Status().Value(),
		CategoryID:  product.CategoryID(),
		UpdatedAt:   product.UpdatedAt(),
	}

	return r.db.Model(&ProductPO{}).Where("id = ?", po.ID).Updates(map[string]interface{}{
		"name":        po.Name,
		"description": po.Description,
		"sku":         po.SKU,
		"price":       po.Price,
		"stock":       po.Stock,
		"status":      po.Status,
		"category_id": po.CategoryID,
		"updated_at":  po.UpdatedAt,
	}).Error
}

// FindByID 根据ID查找商品
func (r *ProductRepository) FindByID(id valueobject.ProductID) (*entity.Product, error) {
	var po ProductPO
	if err := r.db.Where("id = ?", id.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindBySKU 根据SKU查找商品
func (r *ProductRepository) FindBySKU(sku valueobject.ProductSKU) (*entity.Product, error) {
	var po ProductPO
	if err := r.db.Where("sku = ?", sku.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindByCategoryID 查找分类下的商品
func (r *ProductRepository) FindByCategoryID(categoryID int64, page, pageSize int) ([]*entity.Product, int64, error) {
	var count int64
	if err := r.db.Model(&ProductPO{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []ProductPO
	if err := r.db.Where("category_id = ?", categoryID).
		Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		product, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		products[i] = product
	}

	return products, count, nil
}

// FindAll 分页查询商品列表
func (r *ProductRepository) FindAll(page, pageSize int) ([]*entity.Product, int64, error) {
	var count int64
	if err := r.db.Model(&ProductPO{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []ProductPO
	if err := r.db.Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		product, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		products[i] = product
	}

	return products, count, nil
}

// FindByCreatorID 查询创建者的商品
func (r *ProductRepository) FindByCreatorID(creatorID int64, page, pageSize int) ([]*entity.Product, int64, error) {
	var count int64
	if err := r.db.Model(&ProductPO{}).Where("creator_id = ?", creatorID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []ProductPO
	if err := r.db.Where("creator_id = ?", creatorID).
		Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		product, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		products[i] = product
	}

	return products, count, nil
}

// Search 搜索商品
func (r *ProductRepository) Search(keyword string, page, pageSize int) ([]*entity.Product, int64, error) {
	keyword = "%" + strings.TrimSpace(keyword) + "%"

	var count int64
	if err := r.db.Model(&ProductPO{}).
		Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", keyword, keyword, keyword).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []ProductPO
	if err := r.db.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", keyword, keyword, keyword).
		Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		product, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		products[i] = product
	}

	return products, count, nil
}

// FindActiveProducts 查询上架商品
func (r *ProductRepository) FindActiveProducts(page, pageSize int) ([]*entity.Product, int64, error) {
	var count int64
	if err := r.db.Model(&ProductPO{}).Where("status = ?", valueobject.ProductStatusActive).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []ProductPO
	if err := r.db.Where("status = ?", valueobject.ProductStatusActive).
		Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*entity.Product, len(pos))
	for i, po := range pos {
		product, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		products[i] = product
	}

	return products, count, nil
}

// 将持久化对象转换为领域实体
func (r *ProductRepository) toDomainEntity(po *ProductPO) (*entity.Product, error) {
	productID := valueobject.NewProductID(po.ID)

	name, err := valueobject.NewProductName(po.Name)
	if err != nil {
		return nil, err
	}

	description, err := valueobject.NewProductDescription(po.Description)
	if err != nil {
		return nil, err
	}

	sku, err := valueobject.NewProductSKU(po.SKU)
	if err != nil {
		return nil, err
	}

	price, err := valueobject.NewPrice(po.Price)
	if err != nil {
		return nil, err
	}

	status, err := valueobject.NewProductStatus(po.Status)
	if err != nil {
		return nil, err
	}

	creatorID := uservo.NewUserID(po.CreatorID)

	return entity.ReconstructProduct(
		productID,
		name,
		description,
		sku,
		price,
		po.Stock,
		status,
		po.CategoryID,
		creatorID,
		po.CreatedAt,
		po.UpdatedAt,
	), nil
}
