package repository

import (
	"context"

	"wz-backend-go/internal/domain/commerce/entity"
)

// CategoryRepository defines the operations for product category persistence
type CategoryRepository interface {
	// Save stores a new category in the repository
	Save(ctx context.Context, category *entity.Category) error

	// Update updates an existing category
	Update(ctx context.Context, category *entity.Category) error

	// FindByID finds a category by ID
	FindByID(ctx context.Context, id string) (*entity.Category, error)

	// FindByParentID finds categories by parent ID
	FindByParentID(ctx context.Context, parentID string) ([]*entity.Category, error)

	// FindRootCategories finds all root categories (with no parent)
	FindRootCategories(ctx context.Context) ([]*entity.Category, error)

	// FindAll retrieves all categories
	FindAll(ctx context.Context, activeOnly bool) ([]*entity.Category, error)

	// FindByLevel finds categories by level
	FindByLevel(ctx context.Context, level int) ([]*entity.Category, error)

	// FindByName finds a category by name
	FindByName(ctx context.Context, name string) (*entity.Category, error)

	// Delete removes a category
	Delete(ctx context.Context, id string) error

	// CountProducts counts the number of products in a category
	CountProducts(ctx context.Context, categoryID string) (int, error)

	// ReorderCategories updates the sort order of multiple categories
	ReorderCategories(ctx context.Context, categoryIDs []string) error
}
