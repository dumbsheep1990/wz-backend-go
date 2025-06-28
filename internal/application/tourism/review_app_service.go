package tourism

import (
	"context"
	"time"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/repository"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/service"
)

// ReviewAppService defines the application service for tourism review operations
type ReviewAppService struct {
	reviewService     *service.ReviewService
	scenicSpotRepo    repository.ScenicSpotRepository
}

// NewReviewAppService creates a new ReviewAppService instance
func NewReviewAppService(
	reviewService *service.ReviewService,
	scenicSpotRepo repository.ScenicSpotRepository,
) *ReviewAppService {
	return &ReviewAppService{
		reviewService:     reviewService,
		scenicSpotRepo:    scenicSpotRepo,
	}
}

// CreateReview creates a new review for a scenic spot
func (s *ReviewAppService) CreateReview(ctx context.Context, userID, userName string, req *dto.ReviewCreateRequest) (*dto.ReviewResponse, error) {
	// Create review entity from request
	review := entity.NewReview(
		req.ScenicSpotID,
		userID,
		userName,
		req.Content,
		req.Rating,
		req.Images,
	)

	// Create review using domain service
	if err := s.reviewService.CreateReview(ctx, review); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, review)
}

// GetReview retrieves a review by ID
func (s *ReviewAppService) GetReview(ctx context.Context, id string) (*dto.ReviewResponse, error) {
	// Get review using domain service
	review, err := s.reviewService.GetReviewByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, review)
}

// UpdateReview updates a review
func (s *ReviewAppService) UpdateReview(ctx context.Context, id string, userID string, req *dto.ReviewUpdateRequest) (*dto.ReviewResponse, error) {
	// Get existing review
	review, err := s.reviewService.GetReviewByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the review
	if review.UserID != userID {
		return nil, service.ErrUserCanOnlyDeleteOwnReview
	}

	// Update review entity from request
	review.Update(
		req.Content,
		req.Rating,
		req.Images,
	)

	// Update review using domain service
	if err := s.reviewService.UpdateReview(ctx, review); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, review)
}

// DeleteReview deletes a review
func (s *ReviewAppService) DeleteReview(ctx context.Context, id string, userID string) error {
	return s.reviewService.DeleteReview(ctx, id, userID)
}

// LikeReview likes a review
func (s *ReviewAppService) LikeReview(ctx context.Context, id string) error {
	return s.reviewService.LikeReview(ctx, id)
}

// UnlikeReview unlikes a review
func (s *ReviewAppService) UnlikeReview(ctx context.Context, id string) error {
	return s.reviewService.UnlikeReview(ctx, id)
}

// ListReviewsByScenicSpot lists reviews for a specific scenic spot
func (s *ReviewAppService) ListReviewsByScenicSpot(ctx context.Context, scenicSpotID string, offset, limit int) (*dto.ReviewListResponse, error) {
	// List reviews using domain service
	reviews, total, err := s.reviewService.ListReviewsByScenicSpot(ctx, scenicSpotID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to responses
	responses := make([]dto.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		response, err := s.entityToResponse(ctx, review)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}

	// Return response
	return &dto.ReviewListResponse{
		Total: total,
		Items: responses,
	}, nil
}

// ListReviewsByUser lists reviews from a specific user
func (s *ReviewAppService) ListReviewsByUser(ctx context.Context, userID string, offset, limit int) (*dto.ReviewListResponse, error) {
	// List reviews using domain service
	reviews, total, err := s.reviewService.ListReviewsByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to responses
	responses := make([]dto.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		response, err := s.entityToResponse(ctx, review)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}

	// Return response
	return &dto.ReviewListResponse{
		Total: total,
		Items: responses,
	}, nil
}

// Helper methods

// entityToResponse converts a review entity to response
func (s *ReviewAppService) entityToResponse(ctx context.Context, review *entity.Review) (*dto.ReviewResponse, error) {
	// Get scenic spot name if possible
	var scenicSpotName string
	if review.ScenicSpotID != "" {
		scenicSpot, err := s.scenicSpotRepo.GetByID(ctx, review.ScenicSpotID)
		if err == nil && scenicSpot != nil {
			scenicSpotName = scenicSpot.Name
		}
	}

	return &dto.ReviewResponse{
		ID:             review.ID,
		ScenicSpotID:   review.ScenicSpotID,
		ScenicSpotName: scenicSpotName,
		UserID:         review.UserID,
		UserName:       review.UserName,
		Content:        review.Content,
		Rating:         review.Rating,
		Images:         review.Images,
		LikeCount:      review.LikeCount,
		CreatedAt:      review.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      review.UpdatedAt.Format(time.RFC3339),
	}, nil
}
