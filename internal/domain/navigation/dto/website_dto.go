package dto

// WebsiteDTO represents the data transfer object for websites in the navigation system
type WebsiteDTO struct {
	ID          string   `json:"id"`
	CategoryID  string   `json:"category_id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	IconURL     string   `json:"icon_url"`
	SortOrder   int      `json:"sort_order"`
	IsActive    bool     `json:"is_active"`
	IsNew       bool     `json:"is_new"`
	IsFeatured  bool     `json:"is_featured"`
	ViewCount   int64    `json:"view_count"`
	Tags        []string `json:"tags"`
}

// CreateWebsiteRequest represents a request to create a new website
type CreateWebsiteRequest struct {
	CategoryID  string   `json:"category_id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	URL         string   `json:"url" binding:"required,url"`
	Description string   `json:"description"`
	IconURL     string   `json:"icon_url"`
	SortOrder   int      `json:"sort_order"`
	Tags        []string `json:"tags"`
}

// UpdateWebsiteRequest represents a request to update an existing website
type UpdateWebsiteRequest struct {
	ID          string   `json:"id" binding:"required"`
	CategoryID  string   `json:"category_id"`
	Name        string   `json:"name" binding:"required"`
	URL         string   `json:"url" binding:"required,url"`
	Description string   `json:"description"`
	IconURL     string   `json:"icon_url"`
	SortOrder   int      `json:"sort_order"`
	IsActive    bool     `json:"is_active"`
	IsFeatured  bool     `json:"is_featured"`
	Tags        []string `json:"tags"`
}

// WebsiteListResponse represents a response containing a list of websites
type WebsiteListResponse struct {
	Websites []WebsiteDTO `json:"websites"`
	Total    int          `json:"total"`
}

// ReorderWebsitesRequest represents a request to reorder websites
type ReorderWebsitesRequest struct {
	CategoryID string   `json:"category_id" binding:"required"`
	WebsiteIDs []string `json:"website_ids" binding:"required"`
}

// ReorderCategoriesRequest represents a request to reorder categories
type ReorderCategoriesRequest struct {
	CategoryIDs []string `json:"category_ids" binding:"required"`
}
