package repository

import (
	"context"

	"wz-backend-go/internal/domain/commerce/entity"
)

// ProductRepository defines the operations for product persistence
type ProductRepository interface {
	// Save stores a product in the repository
	Save(ctx context.Context, product *entity.Product) error

	// Update updates an existing product
	Update(ctx context.Context, product *entity.Product) error

	// FindByID finds a product by ID
	FindByID(ctx context.Context, id string) (*entity.Product, error)

	// FindAll retrieves all products with optional filters
	FindAll(ctx context.Context, filters ProductFilters) ([]*entity.Product, error)

	// FindByCategory finds products by category ID
	FindByCategory(ctx context.Context, categoryID string, filters ProductFilters) ([]*entity.Product, error)

	// FindByStore finds products by store ID
	FindByStore(ctx context.Context, storeID string, filters ProductFilters) ([]*entity.Product, error)

	// FindByRegion finds products by region
	FindByRegion(ctx context.Context, province, city, district string, filters ProductFilters) ([]*entity.Product, error)

	// FindFeatured finds featured products
	FindFeatured(ctx context.Context, limit int) ([]*entity.Product, error)

	// FindNew finds new products
	FindNew(ctx context.Context, limit int) ([]*entity.Product, error)

	// FindPopular finds popular products based on view count
	FindPopular(ctx context.Context, limit int) ([]*entity.Product, error)

	// Search searches products by name or description
	Search(ctx context.Context, query string, filters ProductFilters) ([]*entity.Product, error)

	// FindByTags finds products by tags
	FindByTags(ctx context.Context, tags []string, filters ProductFilters) ([]*entity.Product, error)

	// Delete removes a product
	Delete(ctx context.Context, id string) error

	// IncrementViewCount increments the view count for a product
	IncrementViewCount(ctx context.Context, id string) error

	// Count returns the total count of products based on filters
	Count(ctx context.Context, filters ProductFilters) (int, error)
}

// ProductFilters defines filters for product queries
type ProductFilters struct {
	// Common filters
	ActiveOnly  bool
	FeaturedOnly bool
	NewOnly     bool
	MinPrice    *float64
	MaxPrice    *float64
	Tags        []string

	// Pagination
	Offset      int
	Limit       int
	SortBy      string
	SortOrder   string
}
