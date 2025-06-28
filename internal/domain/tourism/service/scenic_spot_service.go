package service

import (
	"context"
	"errors"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/repository"
)

// ScenicSpotService defines the domain service for scenic spot operations
type ScenicSpotService struct {
	scenicSpotRepo repository.ScenicSpotRepository
	categoryRepo   repository.CategoryRepository
	reviewRepo     repository.ReviewRepository
}

// NewScenicSpotService creates a new ScenicSpotService instance
func NewScenicSpotService(
	scenicSpotRepo repository.ScenicSpotRepository,
	categoryRepo repository.CategoryRepository,
	reviewRepo repository.ReviewRepository,
) *ScenicSpotService {
	return &ScenicSpotService{
		scenicSpotRepo: scenicSpotRepo,
		categoryRepo:   categoryRepo,
		reviewRepo:     reviewRepo,
	}
}

// CreateScenicSpot creates a new scenic spot
func (s *ScenicSpotService) CreateScenicSpot(ctx context.Context, scenicSpot *entity.ScenicSpot) error {
	// Validate category existence
	if scenicSpot.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(ctx, scenicSpot.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return errors.New("category not found")
		}
	}

	// Create scenic spot
	return s.scenicSpotRepo.Create(ctx, scenicSpot)
}

// GetScenicSpotByID retrieves a scenic spot by ID
func (s *ScenicSpotService) GetScenicSpotByID(ctx context.Context, id string) (*entity.ScenicSpot, error) {
	scenicSpot, err := s.scenicSpotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Increment view count
	if err := s.scenicSpotRepo.IncrementViewCount(ctx, id); err != nil {
		// Log the error but don't fail the request
		// logger.Error("Failed to increment view count", err)
	}

	return scenicSpot, nil
}

// UpdateScenicSpot updates a scenic spot
func (s *ScenicSpotService) UpdateScenicSpot(ctx context.Context, scenicSpot *entity.ScenicSpot) error {
	// Validate category existence
	if scenicSpot.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(ctx, scenicSpot.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return errors.New("category not found")
		}
	}

	// Update scenic spot
	return s.scenicSpotRepo.Update(ctx, scenicSpot)
}

// DeleteScenicSpot deletes a scenic spot
func (s *ScenicSpotService) DeleteScenicSpot(ctx context.Context, id string) error {
	return s.scenicSpotRepo.Delete(ctx, id)
}

// AddReviewToScenicSpot adds a review to a scenic spot
func (s *ScenicSpotService) AddReviewToScenicSpot(ctx context.Context, review *entity.Review) error {
	// Check if scenic spot exists
	scenicSpot, err := s.scenicSpotRepo.GetByID(ctx, review.ScenicSpotID)
	if err != nil {
		return err
	}
	if scenicSpot == nil {
		return errors.New("scenic spot not found")
	}

	// Create review
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return err
	}

	// Update scenic spot rating
	scenicSpot.AddReview(review.Rating)
	return s.scenicSpotRepo.Update(ctx, scenicSpot)
}

// ListScenicSpotsByCategory lists scenic spots by category
func (s *ScenicSpotService) ListScenicSpotsByCategory(ctx context.Context, categoryID string, offset, limit int) ([]*entity.ScenicSpot, int, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	// Get scenic spots by category
	return s.scenicSpotRepo.ListByCategory(ctx, categoryID, offset, limit)
}

// SearchScenicSpots searches for scenic spots by keyword
func (s *ScenicSpotService) SearchScenicSpots(ctx context.Context, keyword string, offset, limit int) ([]*entity.ScenicSpot, int, error) {
	return s.scenicSpotRepo.Search(ctx, keyword, offset, limit)
}
