package dto

// CategoryCreateRequest holds the data needed to create a new category
type CategoryCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// CategoryUpdateRequest holds the data needed to update a category
type CategoryUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// CategoryResponse represents a category data in responses
type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id,omitempty"`
	ParentName  string `json:"parent_name,omitempty"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CategoryListResponse represents the response for listing categories
type CategoryListResponse struct {
	Total int               `json:"total"`
	Items []CategoryResponse `json:"items"`
}

// CategoryTreeNode represents a node in the category tree structure
type CategoryTreeNode struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Icon        string           `json:"icon"`
	SortOrder   int              `json:"sort_order"`
	Children    []CategoryTreeNode `json:"children,omitempty"`
}
