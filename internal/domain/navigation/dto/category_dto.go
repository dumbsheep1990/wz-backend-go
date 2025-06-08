package dto

// CategoryDTO represents the data transfer object for navigation categories
type CategoryDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	WebsiteCount int   `json:"website_count"`
}

// CreateCategoryRequest represents a request to create a new category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest represents a request to update an existing category
type UpdateCategoryRequest struct {
	ID          string `json:"id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// CategoryListResponse represents a response containing a list of categories
type CategoryListResponse struct {
	Categories []CategoryDTO `json:"categories"`
	Total      int           `json:"total"`
}

// CategoryDetailResponse represents a response containing category details with websites
type CategoryDetailResponse struct {
	Category CategoryDTO   `json:"category"`
	Websites []WebsiteDTO  `json:"websites"`
}
