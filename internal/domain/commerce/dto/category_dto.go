package dto

// CreateCategoryRequest represents the request data for creating a product category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
	IconURL     string `json:"icon_url"`
}

// UpdateCategoryRequest represents the request data for updating a product category
type UpdateCategoryRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	IsActive    *bool  `json:"is_active"`
}

// CategoryResponse represents the response data for a product category
type CategoryResponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	DisplayName  string            `json:"display_name"`
	Description  string            `json:"description"`
	ParentID     string            `json:"parent_id,omitempty"`
	IconURL      string            `json:"icon_url"`
	SortOrder    int               `json:"sort_order"`
	Level        int               `json:"level"`
	IsActive     bool              `json:"is_active"`
	ProductCount int               `json:"product_count,omitempty"`
	Children     []*CategoryResponse `json:"children,omitempty"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

// CategoriesResponse represents a list of categories, potentially in a hierarchical structure
type CategoriesResponse struct {
	Categories []*CategoryResponse `json:"categories"`
}

// ReorderCategoriesRequest represents the request for reordering categories
type ReorderCategoriesRequest struct {
	CategoryIDs []string `json:"category_ids" binding:"required"`
}

// CategoryFilterRequest represents a request for filtering categories
type CategoryFilterRequest struct {
	ParentID   string `form:"parent_id"`
	Level      int    `form:"level"`
	ActiveOnly bool   `form:"active_only"`
}
