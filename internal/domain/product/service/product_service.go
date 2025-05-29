package service

import (
	"errors"

	"wz-backend-go/internal/domain/product/entity"
	"wz-backend-go/internal/domain/product/event"
	"wz-backend-go/internal/domain/product/repository"
	"wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// ProductDomainService 商品领域服务
type ProductDomainService struct {
	productRepository repository.ProductRepository
	eventPublisher    repository.EventPublisher
}

// NewProductDomainService 创建商品领域服务
func NewProductDomainService(
	productRepository repository.ProductRepository,
	eventPublisher repository.EventPublisher,
) *ProductDomainService {
	return &ProductDomainService{
		productRepository: productRepository,
		eventPublisher:    eventPublisher,
	}
}

// CreateProduct 创建商品
func (s *ProductDomainService) CreateProduct(
	name valueobject.ProductName,
	description valueobject.ProductDescription,
	sku valueobject.ProductSKU,
	price valueobject.Price,
	stock int32,
	categoryID int64,
	creatorID uservo.UserID,
) (*entity.Product, error) {
	// 检查SKU是否已存在
	existingProduct, err := s.productRepository.FindBySKU(sku)
	if err == nil && existingProduct != nil {
		return nil, errors.New("商品SKU已存在")
	}

	// 创建新商品
	product, err := entity.NewProduct(
		name,
		description,
		sku,
		price,
		stock,
		categoryID,
		creatorID,
	)
	if err != nil {
		return nil, err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return nil, err
	}

	// 发布商品创建事件
	createdEvent := event.NewProductCreatedEvent(
		product.ID(),
		product.Name(),
		product.SKU(),
		product.Price(),
		creatorID.Value(),
	)
	if err := s.eventPublisher.Publish(createdEvent); err != nil {
		// 记录错误但不影响商品创建
		// 可以考虑使用日志记录这个错误
	}

	return product, nil
}

// GetProductByID 根据ID获取商品
func (s *ProductDomainService) GetProductByID(id valueobject.ProductID) (*entity.Product, error) {
	return s.productRepository.FindByID(id)
}

// GetProductBySKU 根据SKU获取商品
func (s *ProductDomainService) GetProductBySKU(sku valueobject.ProductSKU) (*entity.Product, error) {
	return s.productRepository.FindBySKU(sku)
}

// UpdateProduct 更新商品信息
func (s *ProductDomainService) UpdateProduct(
	productID valueobject.ProductID,
	name valueobject.ProductName,
	description valueobject.ProductDescription,
	price valueobject.Price,
	categoryID int64,
) (*entity.Product, error) {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("商品不存在")
	}

	// 更新商品信息
	product.SetName(name)
	product.SetDescription(description)
	product.SetPrice(price)
	if err := product.SetCategoryID(categoryID); err != nil {
		return nil, err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return nil, err
	}

	return product, nil
}

// PublishProduct 发布商品
func (s *ProductDomainService) PublishProduct(productID valueobject.ProductID) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	// 发布商品
	if err := product.Publish(); err != nil {
		return err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布商品发布事件
	publishedEvent := event.NewProductPublishedEvent(productID, product.Stock())
	if err := s.eventPublisher.Publish(publishedEvent); err != nil {
		// 记录错误但不影响商品发布
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// UnpublishProduct 下架商品
func (s *ProductDomainService) UnpublishProduct(productID valueobject.ProductID) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	// 下架商品
	if err := product.Unpublish(); err != nil {
		return err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布商品下架事件
	unpublishedEvent := event.NewProductUnpublishedEvent(productID)
	if err := s.eventPublisher.Publish(unpublishedEvent); err != nil {
		// 记录错误但不影响商品下架
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// DeleteProduct 删除商品
func (s *ProductDomainService) DeleteProduct(productID valueobject.ProductID) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	// 删除商品
	product.Delete()

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布商品删除事件
	deletedEvent := event.NewProductDeletedEvent(productID)
	if err := s.eventPublisher.Publish(deletedEvent); err != nil {
		// 记录错误但不影响商品删除
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// UpdateStock 更新库存
func (s *ProductDomainService) UpdateStock(
	productID valueobject.ProductID,
	stock int32,
	reason string,
) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	previousStock := product.Stock()

	// 更新库存
	product.SetStock(stock)

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布库存变更事件
	changeAmount := stock - previousStock
	stockChangedEvent := event.NewProductStockChangedEvent(
		productID,
		previousStock,
		stock,
		changeAmount,
		reason,
	)
	if err := s.eventPublisher.Publish(stockChangedEvent); err != nil {
		// 记录错误但不影响库存更新
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// DecrementStock 减少库存
func (s *ProductDomainService) DecrementStock(
	productID valueobject.ProductID,
	quantity int32,
	reason string,
) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	previousStock := product.Stock()

	// 减少库存
	if err := product.DecrementStock(quantity); err != nil {
		return err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布库存变更事件
	stockChangedEvent := event.NewProductStockChangedEvent(
		productID,
		previousStock,
		product.Stock(),
		-quantity,
		reason,
	)
	if err := s.eventPublisher.Publish(stockChangedEvent); err != nil {
		// 记录错误但不影响库存减少
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// IncrementStock 增加库存
func (s *ProductDomainService) IncrementStock(
	productID valueobject.ProductID,
	quantity int32,
	reason string,
) error {
	// 获取商品
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	previousStock := product.Stock()

	// 增加库存
	if err := product.IncrementStock(quantity); err != nil {
		return err
	}

	// 保存商品
	if err := s.productRepository.Save(product); err != nil {
		return err
	}

	// 发布库存变更事件
	stockChangedEvent := event.NewProductStockChangedEvent(
		productID,
		previousStock,
		product.Stock(),
		quantity,
		reason,
	)
	if err := s.eventPublisher.Publish(stockChangedEvent); err != nil {
		// 记录错误但不影响库存增加
		// 可以考虑使用日志记录这个错误
	}

	return nil
}

// GetProducts 分页获取商品列表
func (s *ProductDomainService) GetProducts(page, pageSize int) ([]*entity.Product, int64, error) {
	return s.productRepository.FindAll(page, pageSize)
}

// GetProductsByCategory 获取分类下的商品
func (s *ProductDomainService) GetProductsByCategory(categoryID int64, page, pageSize int) ([]*entity.Product, int64, error) {
	return s.productRepository.FindByCategoryID(categoryID, page, pageSize)
}

// GetProductsByCreator 获取创建者的商品
func (s *ProductDomainService) GetProductsByCreator(creatorID int64, page, pageSize int) ([]*entity.Product, int64, error) {
	return s.productRepository.FindByCreatorID(creatorID, page, pageSize)
}

// SearchProducts 搜索商品
func (s *ProductDomainService) SearchProducts(keyword string, page, pageSize int) ([]*entity.Product, int64, error) {
	return s.productRepository.Search(keyword, page, pageSize)
}

// GetActiveProducts 获取上架商品
func (s *ProductDomainService) GetActiveProducts(page, pageSize int) ([]*entity.Product, int64, error) {
	return s.productRepository.FindActiveProducts(page, pageSize)
}
