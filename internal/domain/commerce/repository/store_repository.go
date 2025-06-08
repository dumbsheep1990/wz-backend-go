package repository

import (
	"context"

	"wz-backend-go/internal/domain/commerce/entity"
)

// StoreRepository defines the operations for store persistence
type StoreRepository interface {
	// Save stores a new store in the repository
	Save(ctx context.Context, store *entity.Store) error

	// Update updates an existing store
	Update(ctx context.Context, store *entity.Store) error

	// FindByID finds a store by ID
	FindByID(ctx context.Context, id string) (*entity.Store, error)

	// FindByOwner finds stores owned by a user
	FindByOwner(ctx context.Context, ownerID string) ([]*entity.Store, error)

	// FindAll retrieves all stores with optional filters
	FindAll(ctx context.Context, filters StoreFilters) ([]*entity.Store, error)

	// FindByRegion finds stores by region
	FindByRegion(ctx context.Context, province, city, district string) ([]*entity.Store, error)

	// Search searches stores by name or description
	Search(ctx context.Context, query string) ([]*entity.Store, error)

	// Delete removes a store
	Delete(ctx context.Context, id string) error

	// CountProducts counts the number of products in a store
	CountProducts(ctx context.Context, storeID string) (int, error)

	// UpdateRating updates a store's rating
	UpdateRating(ctx context.Context, storeID string, rating float64) error
}

// StoreFilters defines filters for store queries
type StoreFilters struct {
	// Common filters
	ActiveOnly  bool

	// Pagination
	Offset      int
	Limit       int
	SortBy      string
	SortOrder   string
}
