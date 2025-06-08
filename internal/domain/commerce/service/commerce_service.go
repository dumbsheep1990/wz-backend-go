package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
)

// Common errors
var (
	ErrProductNotFound   = errors.New("product not found")
	ErrStoreNotFound     = errors.New("store not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInvalidParameters = errors.New("invalid parameters")
)

// CommerceService encapsulates business logic for the commerce domain
type CommerceService struct {
	productRepository  repository.ProductRepository
	storeRepository    repository.StoreRepository
	categoryRepository repository.CategoryRepository
}

// NewCommerceService creates a new commerce service
func NewCommerceService(
	productRepo repository.ProductRepository,
	storeRepo repository.StoreRepository,
	categoryRepo repository.CategoryRepository,
) *CommerceService {
	return &CommerceService{
		productRepository:  productRepo,
		storeRepository:    storeRepo,
		categoryRepository: categoryRepo,
	}
}

// CreateProduct creates a new product
func (s *CommerceService) CreateProduct(
	ctx context.Context,
	name, description string,
	price float64,
	categoryID, storeID, region, district string,
	imageURLs []string,
	specifications map[string]string,
	stock int,
	tags []string,
) (*entity.Product, error) {
	// Validate category exists
	category, err := s.categoryRepository.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// Validate store exists
	store, err := s.storeRepository.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	// Create product
	productID := uuid.New().String()
	product := entity.NewProduct(
		productID, name, description, price, 
		categoryID, storeID, region, district,
		imageURLs, specifications, stock, tags,
	)

	// Marshal for storage
	if err := product.Marshal(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.productRepository.Save(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *CommerceService) UpdateProduct(
	ctx context.Context,
	productID, name, description string,
	price float64,
	imageURLs []string,
	specifications map[string]string,
	stock int,
	tags []string,
) (*entity.Product, error) {
	// Find product
	product, err := s.productRepository.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Update product
	product.UpdateDetails(name, description, price, imageURLs, specifications, stock, tags)

	// Marshal for storage
	if err := product.Marshal(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.productRepository.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *CommerceService) GetProduct(ctx context.Context, id string) (*entity.Product, error) {
	product, err := s.productRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Unmarshal JSON fields
	if err := product.Unmarshal(); err != nil {
		return nil, err
	}

	return product, nil
}

// TrackProductView increments the view count for a product
func (s *CommerceService) TrackProductView(ctx context.Context, id string) error {
	return s.productRepository.IncrementViewCount(ctx, id)
}

// CreateStore creates a new store
func (s *CommerceService) CreateStore(
	ctx context.Context,
	name, description, ownerID, logoURL, 
	province, city, district, address, 
	contactName, phone string,
) (*entity.Store, error) {
	storeID := uuid.New().String()
	store := entity.NewStore(
		storeID, name, description, ownerID, logoURL,
		province, city, district, address, contactName, phone,
	)

	if err := s.storeRepository.Save(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

// UpdateStore updates an existing store
func (s *CommerceService) UpdateStore(
	ctx context.Context,
	storeID, name, description, logoURL, 
	province, city, district, address, 
	contactName, phone string,
) (*entity.Store, error) {
	store, err := s.storeRepository.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	store.UpdateDetails(name, description, logoURL, province, city, district, address, contactName, phone)

	if err := s.storeRepository.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

// GetStore retrieves a store by ID
func (s *CommerceService) GetStore(ctx context.Context, id string) (*entity.Store, error) {
	store, err := s.storeRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	return store, nil
}

// CreateCategory creates a new product category
func (s *CommerceService) CreateCategory(
	ctx context.Context,
	name, displayName, description, parentID, iconURL string,
) (*entity.Category, error) {
	categoryID := uuid.New().String()
	
	// Determine level based on parent
	level := 1
	if parentID != "" {
		parent, err := s.categoryRepository.FindByID(ctx, parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, ErrCategoryNotFound
		}
		level = parent.Level + 1
	}
	
	// Determine sort order (place at end)
	var siblings []*entity.Category
	var err error
	
	if parentID == "" {
		siblings, err = s.categoryRepository.FindRootCategories(ctx)
	} else {
		siblings, err = s.categoryRepository.FindByParentID(ctx, parentID)
	}
	
	if err != nil {
		return nil, err
	}
	
	sortOrder := 0
	if len(siblings) > 0 {
		sortOrder = siblings[len(siblings)-1].SortOrder + 1
	}
	
	category := entity.NewCategory(categoryID, name, displayName, description, parentID, iconURL, sortOrder, level)

	if err := s.categoryRepository.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *CommerceService) UpdateCategory(
	ctx context.Context,
	categoryID, name, displayName, description, iconURL string,
) (*entity.Category, error) {
	category, err := s.categoryRepository.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	category.UpdateDetails(name, displayName, description, iconURL)

	if err := s.categoryRepository.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// ReorderCategories reorders multiple categories
func (s *CommerceService) ReorderCategories(ctx context.Context, categoryIDs []string) error {
	if len(categoryIDs) == 0 {
		return ErrInvalidParameters
	}
	
	return s.categoryRepository.ReorderCategories(ctx, categoryIDs)
}

// GetCategory retrieves a category by ID
func (s *CommerceService) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	category, err := s.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}
