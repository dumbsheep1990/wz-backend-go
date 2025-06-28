package entity

import (
	"time"

	"github.com/google/uuid"
)

// ScenicSpot represents a tourism scenic spot
type ScenicSpot struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	CategoryID      string    `json:"category_id"`
	Address         string    `json:"address"`
	LocationArea    string    `json:"location_area"` // 区域划分，如江宁区等
	Description     string    `json:"description"`
	OpeningHours    string    `json:"opening_hours"` // 开放时间
	Price           float64   `json:"price"`
	TicketInfo      string    `json:"ticket_info"`   // 票务信息
	ScenicFeatures  string    `json:"scenic_features"` // 景区特色，如海盐历史文化
	TransportInfo   string    `json:"transport_info"`  // 交通信息
	Images          []string  `json:"images"`
	Rating          float64   `json:"rating"`
	ReviewCount     int       `json:"review_count"`
	ViewCount       int       `json:"view_count"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	NearbyFacilities []string `json:"nearby_facilities"` // 附近设施，如盐城站、公交站
	TransitRoutes    []string `json:"transit_routes"`    // 公交线路，如12路、21路
}

// NewScenicSpot creates a new scenic spot entity
func NewScenicSpot(name, categoryID, address, locationArea, description, openingHours string, 
                  price float64, ticketInfo, scenicFeatures, transportInfo string, 
                  images []string, latitude, longitude float64) *ScenicSpot {
	now := time.Now()
	return &ScenicSpot{
		ID:              uuid.New().String(),
		Name:            name,
		CategoryID:      categoryID,
		Address:         address,
		LocationArea:    locationArea,
		Description:     description,
		OpeningHours:    openingHours,
		Price:           price,
		TicketInfo:      ticketInfo,
		ScenicFeatures:  scenicFeatures,
		TransportInfo:   transportInfo,
		Images:          images,
		Rating:          0,
		ReviewCount:     0,
		ViewCount:       0,
		Latitude:        latitude,
		Longitude:       longitude,
		CreatedAt:       now,
		UpdatedAt:       now,
		NearbyFacilities: []string{},
		TransitRoutes:    []string{},
	}
}

// Update updates the scenic spot information
func (s *ScenicSpot) Update(name, categoryID, address, locationArea, description, openingHours string, 
                          price float64, ticketInfo, scenicFeatures, transportInfo string, 
                          images []string, latitude, longitude float64) {
	s.Name = name
	s.CategoryID = categoryID
	s.Address = address
	s.LocationArea = locationArea
	s.Description = description
	s.OpeningHours = openingHours
	s.Price = price
	s.TicketInfo = ticketInfo
	s.ScenicFeatures = scenicFeatures
	s.TransportInfo = transportInfo
	s.Images = images
	s.Latitude = latitude
	s.Longitude = longitude
	s.UpdatedAt = time.Now()
}

// AddNearbyFacility adds a nearby facility
func (s *ScenicSpot) AddNearbyFacility(facility string) {
	s.NearbyFacilities = append(s.NearbyFacilities, facility)
	s.UpdatedAt = time.Now()
}

// AddTransitRoute adds a transit route
func (s *ScenicSpot) AddTransitRoute(route string) {
	s.TransitRoutes = append(s.TransitRoutes, route)
	s.UpdatedAt = time.Now()
}

// IncrementViewCount increases view count by one
func (s *ScenicSpot) IncrementViewCount() {
	s.ViewCount++
	s.UpdatedAt = time.Now()
}

// AddReview adds a new review rating to the scenic spot
func (s *ScenicSpot) AddReview(rating float64) {
	// Calculate new average rating
	totalRating := s.Rating * float64(s.ReviewCount)
	totalRating += rating
	s.ReviewCount++
	s.Rating = totalRating / float64(s.ReviewCount)
	s.UpdatedAt = time.Now()
}
