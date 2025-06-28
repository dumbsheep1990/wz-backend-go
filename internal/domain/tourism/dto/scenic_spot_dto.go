package dto

// ScenicSpotCreateRequest holds the data needed to create a new scenic spot
type ScenicSpotCreateRequest struct {
	Name           string   `json:"name" binding:"required"`
	CategoryID     string   `json:"category_id" binding:"required"`
	Address        string   `json:"address" binding:"required"`
	LocationArea   string   `json:"location_area" binding:"required"`
	Description    string   `json:"description"`
	OpeningHours   string   `json:"opening_hours"`
	Price          float64  `json:"price"`
	TicketInfo     string   `json:"ticket_info"`
	ScenicFeatures string   `json:"scenic_features"`
	TransportInfo  string   `json:"transport_info"`
	Images         []string `json:"images"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	NearbyFacilities []string `json:"nearby_facilities"`
	TransitRoutes    []string `json:"transit_routes"`
}

// ScenicSpotUpdateRequest holds the data needed to update a scenic spot
type ScenicSpotUpdateRequest struct {
	Name           string   `json:"name"`
	CategoryID     string   `json:"category_id"`
	Address        string   `json:"address"`
	LocationArea   string   `json:"location_area"`
	Description    string   `json:"description"`
	OpeningHours   string   `json:"opening_hours"`
	Price          float64  `json:"price"`
	TicketInfo     string   `json:"ticket_info"`
	ScenicFeatures string   `json:"scenic_features"`
	TransportInfo  string   `json:"transport_info"`
	Images         []string `json:"images"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	NearbyFacilities []string `json:"nearby_facilities"`
	TransitRoutes    []string `json:"transit_routes"`
}

// ScenicSpotResponse represents a scenic spot data in responses
type ScenicSpotResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CategoryID     string   `json:"category_id"`
	CategoryName   string   `json:"category_name,omitempty"`
	Address        string   `json:"address"`
	LocationArea   string   `json:"location_area"`
	Description    string   `json:"description"`
	OpeningHours   string   `json:"opening_hours"`
	Price          float64  `json:"price"`
	TicketInfo     string   `json:"ticket_info"`
	ScenicFeatures string   `json:"scenic_features"`
	TransportInfo  string   `json:"transport_info"`
	Images         []string `json:"images"`
	Rating         float64  `json:"rating"`
	ReviewCount    int      `json:"review_count"`
	ViewCount      int      `json:"view_count"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	NearbyFacilities []string `json:"nearby_facilities"`
	TransitRoutes    []string `json:"transit_routes"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// ScenicSpotListItem represents a scenic spot item in list response
type ScenicSpotListItem struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CategoryID     string   `json:"category_id"`
	CategoryName   string   `json:"category_name,omitempty"`
	Address        string   `json:"address"`
	LocationArea   string   `json:"location_area"`
	Price          float64  `json:"price"`
	Images         []string `json:"images"`
	Rating         float64  `json:"rating"`
	ReviewCount    int      `json:"review_count"`
	ViewCount      int      `json:"view_count"`
	ScenicFeatures string   `json:"scenic_features"`
}

// ScenicSpotListResponse represents the response for listing scenic spots
type ScenicSpotListResponse struct {
	Total int                `json:"total"`
	Items []ScenicSpotListItem `json:"items"`
}
