package dto

// CreateStoreRequest represents the request data for creating a store
type CreateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" binding:"required"`
	LogoURL     string `json:"logo_url"`
	Province    string `json:"province" binding:"required"`
	City        string `json:"city" binding:"required"`
	District    string `json:"district"`
	Address     string `json:"address"`
	ContactName string `json:"contact_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
}

// UpdateStoreRequest represents the request data for updating a store
type UpdateStoreRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	Province    string `json:"province" binding:"required"`
	City        string `json:"city" binding:"required"`
	District    string `json:"district"`
	Address     string `json:"address"`
	ContactName string `json:"contact_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	IsActive    *bool  `json:"is_active"`
}

// StoreResponse represents the response data for a store
type StoreResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	OwnerID     string  `json:"owner_id"`
	LogoURL     string  `json:"logo_url"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	District    string  `json:"district"`
	Address     string  `json:"address"`
	ContactName string  `json:"contact_name"`
	Phone       string  `json:"phone"`
	Rating      float64 `json:"rating"`
	IsActive    bool    `json:"is_active"`
	ProductCount int    `json:"product_count,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// StoreListResponse represents a store in a list
type StoreListResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	LogoURL     string  `json:"logo_url"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	District    string  `json:"district"`
	Rating      float64 `json:"rating"`
	IsActive    bool    `json:"is_active"`
	ProductCount int    `json:"product_count"`
	CreatedAt   string  `json:"created_at"`
}

// StoresResponse represents a paginated list of stores
type StoresResponse struct {
	Stores     []*StoreListResponse `json:"stores"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// StoreFilterRequest represents a request for filtering stores
type StoreFilterRequest struct {
	Province  string `form:"province"`
	City      string `form:"city"`
	District  string `form:"district"`
	OwnerID   string `form:"owner_id"`
	ActiveOnly bool   `form:"active_only"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}
