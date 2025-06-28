package service

import (
	"context"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/repository"
)

// ReviewService defines the domain service for tourism review operations
type ReviewService struct {
	reviewRepo     repository.ReviewRepository
	scenicSpotRepo repository.ScenicSpotRepository
}

// NewReviewService creates a new ReviewService instance
func NewReviewService(
	reviewRepo repository.ReviewRepository,
	scenicSpotRepo repository.ScenicSpotRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo:     reviewRepo,
		scenicSpotRepo: scenicSpotRepo,
	}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, review *entity.Review) error {
	// Validate rating
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidReviewRating
	}

	// Check if scenic spot exists
	scenicSpot, err := s.scenicSpotRepo.GetByID(ctx, review.ScenicSpotID)
	if err != nil {
		return err
	}
	if scenicSpot == nil {
		return ErrScenicSpotNotFound
	}

	// Get reviews by user and scenic spot to check if user already reviewed
	reviews, _, err := s.reviewRepo.ListByScenicSpot(ctx, review.ScenicSpotID, 0, 100)
	if err != nil {
		return err
	}

	// Check if user already reviewed this scenic spot
	for _, existingReview := range reviews {
		if existingReview.UserID == review.UserID {
			return ErrUserAlreadyReviewed
		}
	}

	// Create review
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return err
	}

	// Update scenic spot rating
	scenicSpot.AddReview(review.Rating)
	return s.scenicSpotRepo.Update(ctx, scenicSpot)
}

// GetReviewByID retrieves a review by ID
func (s *ReviewService) GetReviewByID(ctx context.Context, id string) (*entity.Review, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

// UpdateReview updates a review
func (s *ReviewService) UpdateReview(ctx context.Context, review *entity.Review) error {
	// Validate rating
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidReviewRating
	}

	// Check if review exists
	existingReview, err := s.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		return err
	}
	if existingReview == nil {
		return ErrReviewNotFound
	}

	// Update review
	return s.reviewRepo.Update(ctx, review)
}

// DeleteReview deletes a review
func (s *ReviewService) DeleteReview(ctx context.Context, id string, userID string) error {
	// Check if review exists and belongs to the user
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if review == nil {
		return ErrReviewNotFound
	}

	// Verify user owns the review
	if review.UserID != userID {
		return ErrUserCanOnlyDeleteOwnReview
	}

	return s.reviewRepo.Delete(ctx, id)
}

// LikeReview increments the like count of a review
func (s *ReviewService) LikeReview(ctx context.Context, id string) error {
	// Check if review exists
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if review == nil {
		return ErrReviewNotFound
	}

	return s.reviewRepo.IncrementLikes(ctx, id)
}

// UnlikeReview decrements the like count of a review
func (s *ReviewService) UnlikeReview(ctx context.Context, id string) error {
	// Check if review exists
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if review == nil {
		return ErrReviewNotFound
	}

	return s.reviewRepo.DecrementLikes(ctx, id)
}

// ListReviewsByScenicSpot lists reviews for a specific scenic spot
func (s *ReviewService) ListReviewsByScenicSpot(ctx context.Context, scenicSpotID string, offset, limit int) ([]*entity.Review, int, error) {
	return s.reviewRepo.ListByScenicSpot(ctx, scenicSpotID, offset, limit)
}

// ListReviewsByUser lists reviews from a specific user
func (s *ReviewService) ListReviewsByUser(ctx context.Context, userID string, offset, limit int) ([]*entity.Review, int, error) {
	return s.reviewRepo.ListByUser(ctx, userID, offset, limit)
}
