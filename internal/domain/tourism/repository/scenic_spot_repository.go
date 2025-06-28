package repository

import (
	"context"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
)

// ScenicSpotRepository defines the interface for scenic spot data access
type ScenicSpotRepository interface {
	// Create creates a new scenic spot
	Create(ctx context.Context, scenicSpot *entity.ScenicSpot) error
	
	// GetByID retrieves a scenic spot by ID
	GetByID(ctx context.Context, id string) (*entity.ScenicSpot, error)
	
	// Update updates a scenic spot
	Update(ctx context.Context, scenicSpot *entity.ScenicSpot) error
	
	// Delete deletes a scenic spot by ID
	Delete(ctx context.Context, id string) error
	
	// List retrieves all scenic spots with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.ScenicSpot, int, error)
	
	// ListByCategory retrieves scenic spots by category ID with pagination
	ListByCategory(ctx context.Context, categoryID string, offset, limit int) ([]*entity.ScenicSpot, int, error)
	
	// ListByArea retrieves scenic spots by location area with pagination
	ListByArea(ctx context.Context, area string, offset, limit int) ([]*entity.ScenicSpot, int, error)
	
	// IncrementViewCount increments the view count of a scenic spot
	IncrementViewCount(ctx context.Context, id string) error
	
	// Search searches for scenic spots by keywords with pagination
	Search(ctx context.Context, keyword string, offset, limit int) ([]*entity.ScenicSpot, int, error)
}
