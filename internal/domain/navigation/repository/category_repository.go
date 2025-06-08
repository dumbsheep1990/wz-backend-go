package repository

import (
	"context"

	"github.com/google/uuid"
	"wz-backend-go/internal/domain/navigation/entity"
)

// CategoryRepository defines the interface for category persistence operations
type CategoryRepository interface {
	// Save persists a category
	Save(ctx context.Context, category *entity.Category) error
	
	// FindByID retrieves a category by its ID
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	
	// FindByName retrieves a category by its name
	FindByName(ctx context.Context, name string) (*entity.Category, error)
	
	// FindAll retrieves all categories
	FindAll(ctx context.Context) ([]*entity.Category, error)
	
	// FindActive retrieves all active categories
	FindActive(ctx context.Context) ([]*entity.Category, error)
	
	// Update updates an existing category
	Update(ctx context.Context, category *entity.Category) error
	
	// Delete removes a category
	Delete(ctx context.Context, id string) error
	
	// FindAllSorted retrieves all categories ordered by sortOrder
	FindAllSorted(ctx context.Context) ([]*entity.Category, error)
	
	// CountWebsites counts websites in a category
	CountWebsites(ctx context.Context, categoryID string) (int, error)
	
	// GenerateID generates a new unique ID
	GenerateID() string
}

// NewID generates a new ID for a category
func NewID() string {
	return uuid.New().String()
}
