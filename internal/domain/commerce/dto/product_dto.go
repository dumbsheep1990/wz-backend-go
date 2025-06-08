package dto

// CreateProductRequest represents the request data for creating a product
type CreateProductRequest struct {
	Name           string             `json:"name" binding:"required"`
	Description    string             `json:"description"`
	Price          float64            `json:"price" binding:"required,gt=0"`
	CategoryID     string             `json:"category_id" binding:"required"`
	StoreID        string             `json:"store_id" binding:"required"`
	Region         string             `json:"region"`
	District       string             `json:"district"`
	ImageURLs      []string           `json:"image_urls"`
	Specifications map[string]string  `json:"specifications"`
	Stock          int                `json:"stock" binding:"required,gte=0"`
	Tags           []string           `json:"tags"`
}

// UpdateProductRequest represents the request data for updating a product
type UpdateProductRequest struct {
	ID             string             `json:"id" binding:"required"`
	Name           string             `json:"name" binding:"required"`
	Description    string             `json:"description"`
	Price          float64            `json:"price" binding:"required,gt=0"`
	ImageURLs      []string           `json:"image_urls"`
	Specifications map[string]string  `json:"specifications"`
	Stock          int                `json:"stock" binding:"required,gte=0"`
	Tags           []string           `json:"tags"`
	IsActive       *bool              `json:"is_active"`
	IsFeatured     *bool              `json:"is_featured"`
	IsNew          *bool              `json:"is_new"`
}

// ProductResponse represents the response data for a product
type ProductResponse struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Price          float64            `json:"price"`
	FormattedPrice string             `json:"formatted_price"`
	CategoryID     string             `json:"category_id"`
	CategoryName   string             `json:"category_name,omitempty"`
	StoreID        string             `json:"store_id"`
	StoreName      string             `json:"store_name,omitempty"`
	Region         string             `json:"region"`
	District       string             `json:"district"`
	ImageURLs      []string           `json:"image_urls"`
	Specifications map[string]string  `json:"specifications"`
	Stock          int                `json:"stock"`
	ViewCount      int64              `json:"view_count"`
	SoldCount      int                `json:"sold_count"`
	IsActive       bool               `json:"is_active"`
	IsFeatured     bool               `json:"is_featured"`
	IsNew          bool               `json:"is_new"`
	Tags           []string           `json:"tags"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
}

// ProductListResponse represents a product in a list
type ProductListResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	FormattedPrice string    `json:"formatted_price"`
	CategoryID     string    `json:"category_id"`
	CategoryName   string    `json:"category_name,omitempty"`
	StoreID        string    `json:"store_id"`
	StoreName      string    `json:"store_name,omitempty"`
	ImageURL       string    `json:"image_url"`  // Primary image
	IsActive       bool      `json:"is_active"`
	IsFeatured     bool      `json:"is_featured"`
	IsNew          bool      `json:"is_new"`
	ViewCount      int64     `json:"view_count"`
	CreatedAt      string    `json:"created_at"`
}

// ProductFilterRequest represents a request for filtering products
type ProductFilterRequest struct {
	CategoryID    string   `form:"category_id"`
	StoreID       string   `form:"store_id"`
	Region        string   `form:"region"`
	District      string   `form:"district"`
	MinPrice      *float64 `form:"min_price"`
	MaxPrice      *float64 `form:"max_price"`
	Tags          string   `form:"tags"` // Comma separated list
	ActiveOnly    bool     `form:"active_only"`
	FeaturedOnly  bool     `form:"featured_only"`
	NewOnly       bool     `form:"new_only"`
	SortBy        string   `form:"sort_by"`
	SortOrder     string   `form:"sort_order"`
	Page          int      `form:"page"`
	PageSize      int      `form:"page_size"`
}

// ProductsResponse represents a paginated list of products
type ProductsResponse struct {
	Products    []*ProductListResponse `json:"products"`
	Total       int                    `json:"total"`
	Page        int                    `json:"page"`
	PageSize    int                    `json:"page_size"`
	TotalPages  int                    `json:"total_pages"`
}
