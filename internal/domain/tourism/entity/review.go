package entity

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a user review for a scenic spot
type Review struct {
	ID          string    `json:"id"`
	ScenicSpotID string   `json:"scenic_spot_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"` // Denormalized for display
	Content     string    `json:"content"`
	Rating      float64   `json:"rating"`    // 1-5 star rating
	Images      []string  `json:"images"`    // Optional review photos
	LikeCount   int       `json:"like_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewReview creates a new review entity
func NewReview(scenicSpotID, userID, userName, content string, rating float64, images []string) *Review {
	now := time.Now()
	return &Review{
		ID:          uuid.New().String(),
		ScenicSpotID: scenicSpotID,
		UserID:      userID,
		UserName:    userName,
		Content:     content,
		Rating:      rating,
		Images:      images,
		LikeCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the review information
func (r *Review) Update(content string, rating float64, images []string) {
	r.Content = content
	r.Rating = rating
	r.Images = images
	r.UpdatedAt = time.Now()
}

// IncrementLikes increases the like count by one
func (r *Review) IncrementLikes() {
	r.LikeCount++
	r.UpdatedAt = time.Now()
}

// DecrementLikes decreases the like count by one
func (r *Review) DecrementLikes() {
	if r.LikeCount > 0 {
		r.LikeCount--
		r.UpdatedAt = time.Now()
	}
}
