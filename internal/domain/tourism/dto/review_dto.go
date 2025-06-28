package dto

// ReviewCreateRequest holds the data needed to create a new review
type ReviewCreateRequest struct {
	ScenicSpotID string   `json:"scenic_spot_id" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Rating       float64  `json:"rating" binding:"required,min=1,max=5"`
	Images       []string `json:"images"`
}

// ReviewUpdateRequest holds the data needed to update a review
type ReviewUpdateRequest struct {
	Content string   `json:"content"`
	Rating  float64  `json:"rating" binding:"min=1,max=5"`
	Images  []string `json:"images"`
}

// ReviewResponse represents a review data in responses
type ReviewResponse struct {
	ID           string   `json:"id"`
	ScenicSpotID string   `json:"scenic_spot_id"`
	ScenicSpotName string `json:"scenic_spot_name,omitempty"`
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	Content      string   `json:"content"`
	Rating       float64  `json:"rating"`
	Images       []string `json:"images"`
	LikeCount    int      `json:"like_count"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// ReviewListResponse represents the response for listing reviews
type ReviewListResponse struct {
	Total int             `json:"total"`
	Items []ReviewResponse `json:"items"`
}
