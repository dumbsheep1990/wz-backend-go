package service

import "errors"

// Domain service errors
var (
	// Category errors
	ErrParentCategoryNotFound     = errors.New("parent category not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrCategoryHasSubcategories   = errors.New("cannot delete category with subcategories")
	ErrCategoryHasScenicSpots     = errors.New("cannot delete category with scenic spots")
	ErrCategoryCannotBeItsOwnParent = errors.New("category cannot be its own parent")

	// Scenic spot errors
	ErrScenicSpotNotFound         = errors.New("scenic spot not found")
	ErrInvalidScenicSpotData      = errors.New("invalid scenic spot data")
	ErrScenicSpotHasReviews       = errors.New("cannot delete scenic spot with reviews")

	// Review errors
	ErrReviewNotFound             = errors.New("review not found")
	ErrUserAlreadyReviewed        = errors.New("user already reviewed this scenic spot")
	ErrInvalidReviewRating        = errors.New("review rating must be between 1 and 5")
	ErrUserCanOnlyDeleteOwnReview = errors.New("user can only delete their own review")
)
