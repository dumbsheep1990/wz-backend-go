package repository

import (
	"context"

	"wz-backend-go/internal/domain/navigation/entity"
)

// WebsiteRepository defines the interface for website persistence operations
type WebsiteRepository interface {
	// Save persists a website
	Save(ctx context.Context, website *entity.Website) error
	
	// FindByID retrieves a website by its ID
	FindByID(ctx context.Context, id string) (*entity.Website, error)
	
	// FindByURL retrieves a website by its URL
	FindByURL(ctx context.Context, url string) (*entity.Website, error)
	
	// FindByCategory retrieves all websites in a category
	FindByCategory(ctx context.Context, categoryID string) ([]*entity.Website, error)
	
	// FindByCategorySorted retrieves all websites in a category ordered by sortOrder
	FindByCategorySorted(ctx context.Context, categoryID string) ([]*entity.Website, error)
	
	// FindAll retrieves all websites
	FindAll(ctx context.Context) ([]*entity.Website, error)
	
	// FindFeatured retrieves all featured websites
	FindFeatured(ctx context.Context) ([]*entity.Website, error)
	
	// FindPopular retrieves popular websites by view count
	FindPopular(ctx context.Context, limit int) ([]*entity.Website, error)
	
	// FindByTags finds websites by tags
	FindByTags(ctx context.Context, tags []string) ([]*entity.Website, error)
	
	// Update updates an existing website
	Update(ctx context.Context, website *entity.Website) error
	
	// Delete removes a website
	Delete(ctx context.Context, id string) error
	
	// IncrementViewCount increments the view count for a website
	IncrementViewCount(ctx context.Context, id string) error
	
	// CountByCategory counts websites in a category
	CountByCategory(ctx context.Context, categoryID string) (int, error)
}
