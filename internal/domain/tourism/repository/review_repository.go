package repository

import (
	"context"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
)

// ReviewRepository defines the interface for tourism review data access
type ReviewRepository interface {
	// Create creates a new review
	Create(ctx context.Context, review *entity.Review) error
	
	// GetByID retrieves a review by ID
	GetByID(ctx context.Context, id string) (*entity.Review, error)
	
	// Update updates a review
	Update(ctx context.Context, review *entity.Review) error
	
	// Delete deletes a review by ID
	Delete(ctx context.Context, id string) error
	
	// ListByScenicSpot retrieves reviews for a specific scenic spot with pagination
	ListByScenicSpot(ctx context.Context, scenicSpotID string, offset, limit int) ([]*entity.Review, int, error)
	
	// ListByUser retrieves reviews from a specific user with pagination
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*entity.Review, int, error)
	
	// IncrementLikes increments the like count of a review
	IncrementLikes(ctx context.Context, id string) error
	
	// DecrementLikes decrements the like count of a review
	DecrementLikes(ctx context.Context, id string) error
}
