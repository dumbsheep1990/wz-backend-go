package repository

import (
	"context"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
)

// CategoryRepository defines the interface for tourism category data access
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *entity.Category) error
	
	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	
	// Update updates a category
	Update(ctx context.Context, category *entity.Category) error
	
	// Delete deletes a category by ID
	Delete(ctx context.Context, id string) error
	
	// List retrieves all categories with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.Category, int, error)
	
	// ListByParent retrieves categories by parent ID
	ListByParent(ctx context.Context, parentID string) ([]*entity.Category, error)
	
	// ListRootCategories retrieves top-level categories
	ListRootCategories(ctx context.Context) ([]*entity.Category, error)
}
