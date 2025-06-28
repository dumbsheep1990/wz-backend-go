package tourism

import (
	"context"
	"time"

	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/dto"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/entity"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/repository"
	"github.com/wanzhitouzi/wz-backend-go/internal/domain/tourism/service"
)

// ScenicSpotAppService defines the application service for scenic spot operations
type ScenicSpotAppService struct {
	scenicSpotService *service.ScenicSpotService
	categoryRepo      repository.CategoryRepository
}

// NewScenicSpotAppService creates a new ScenicSpotAppService instance
func NewScenicSpotAppService(
	scenicSpotService *service.ScenicSpotService,
	categoryRepo repository.CategoryRepository,
) *ScenicSpotAppService {
	return &ScenicSpotAppService{
		scenicSpotService: scenicSpotService,
		categoryRepo:      categoryRepo,
	}
}

// CreateScenicSpot creates a new scenic spot
func (s *ScenicSpotAppService) CreateScenicSpot(ctx context.Context, req *dto.ScenicSpotCreateRequest) (*dto.ScenicSpotResponse, error) {
	// Create scenic spot entity from request
	scenicSpot := entity.NewScenicSpot(
		req.Name,
		req.CategoryID,
		req.Address,
		req.LocationArea,
		req.Description,
		req.OpeningHours,
		req.Price,
		req.TicketInfo,
		req.ScenicFeatures,
		req.TransportInfo,
		req.Images,
		req.Latitude,
		req.Longitude,
	)

	// Add nearby facilities and transit routes
	for _, facility := range req.NearbyFacilities {
		scenicSpot.AddNearbyFacility(facility)
	}
	
	for _, route := range req.TransitRoutes {
		scenicSpot.AddTransitRoute(route)
	}

	// Create scenic spot using domain service
	if err := s.scenicSpotService.CreateScenicSpot(ctx, scenicSpot); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, scenicSpot)
}

// GetScenicSpot retrieves a scenic spot by ID
func (s *ScenicSpotAppService) GetScenicSpot(ctx context.Context, id string) (*dto.ScenicSpotResponse, error) {
	// Get scenic spot using domain service
	scenicSpot, err := s.scenicSpotService.GetScenicSpotByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, scenicSpot)
}

// UpdateScenicSpot updates a scenic spot
func (s *ScenicSpotAppService) UpdateScenicSpot(ctx context.Context, id string, req *dto.ScenicSpotUpdateRequest) (*dto.ScenicSpotResponse, error) {
	// Get existing scenic spot
	scenicSpot, err := s.scenicSpotService.GetScenicSpotByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update scenic spot entity from request
	scenicSpot.Update(
		req.Name,
		req.CategoryID,
		req.Address,
		req.LocationArea,
		req.Description,
		req.OpeningHours,
		req.Price,
		req.TicketInfo,
		req.ScenicFeatures,
		req.TransportInfo,
		req.Images,
		req.Latitude,
		req.Longitude,
	)

	// Update nearby facilities and transit routes if provided
	if req.NearbyFacilities != nil {
		scenicSpot.NearbyFacilities = req.NearbyFacilities
	}
	
	if req.TransitRoutes != nil {
		scenicSpot.TransitRoutes = req.TransitRoutes
	}

	// Update scenic spot using domain service
	if err := s.scenicSpotService.UpdateScenicSpot(ctx, scenicSpot); err != nil {
		return nil, err
	}

	// Return response
	return s.entityToResponse(ctx, scenicSpot)
}

// DeleteScenicSpot deletes a scenic spot
func (s *ScenicSpotAppService) DeleteScenicSpot(ctx context.Context, id string) error {
	return s.scenicSpotService.DeleteScenicSpot(ctx, id)
}

// ListScenicSpots lists all scenic spots with pagination
func (s *ScenicSpotAppService) ListScenicSpots(ctx context.Context, offset, limit int) (*dto.ScenicSpotListResponse, error) {
	// Get scenic spots using domain service
	spots, total, err := s.scenicSpotService.ListScenicSpots(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response items
	items := make([]dto.ScenicSpotListItem, 0, len(spots))
	for _, spot := range spots {
		item, err := s.entityToListItem(ctx, spot)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	// Return response
	return &dto.ScenicSpotListResponse{
		Total: total,
		Items: items,
	}, nil
}

// ListScenicSpotsByCategory lists scenic spots by category with pagination
func (s *ScenicSpotAppService) ListScenicSpotsByCategory(ctx context.Context, categoryID string, offset, limit int) (*dto.ScenicSpotListResponse, error) {
	// Get scenic spots using domain service
	spots, total, err := s.scenicSpotService.ListScenicSpotsByCategory(ctx, categoryID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response items
	items := make([]dto.ScenicSpotListItem, 0, len(spots))
	for _, spot := range spots {
		item, err := s.entityToListItem(ctx, spot)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	// Return response
	return &dto.ScenicSpotListResponse{
		Total: total,
		Items: items,
	}, nil
}

// ListScenicSpotsByArea lists scenic spots by area with pagination
func (s *ScenicSpotAppService) ListScenicSpotsByArea(ctx context.Context, area string, offset, limit int) (*dto.ScenicSpotListResponse, error) {
	// Get scenic spots using domain service
	spots, total, err := s.scenicSpotService.ListScenicSpotsByArea(ctx, area, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response items
	items := make([]dto.ScenicSpotListItem, 0, len(spots))
	for _, spot := range spots {
		item, err := s.entityToListItem(ctx, spot)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	// Return response
	return &dto.ScenicSpotListResponse{
		Total: total,
		Items: items,
	}, nil
}

// SearchScenicSpots searches for scenic spots by keyword
func (s *ScenicSpotAppService) SearchScenicSpots(ctx context.Context, keyword string, offset, limit int) (*dto.ScenicSpotListResponse, error) {
	// Search scenic spots using domain service
	spots, total, err := s.scenicSpotService.SearchScenicSpots(ctx, keyword, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response items
	items := make([]dto.ScenicSpotListItem, 0, len(spots))
	for _, spot := range spots {
		item, err := s.entityToListItem(ctx, spot)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	// Return response
	return &dto.ScenicSpotListResponse{
		Total: total,
		Items: items,
	}, nil
}

// Helper methods

// entityToResponse converts a scenic spot entity to response
func (s *ScenicSpotAppService) entityToResponse(ctx context.Context, spot *entity.ScenicSpot) (*dto.ScenicSpotResponse, error) {
	// Get category name if possible
	var categoryName string
	if spot.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(ctx, spot.CategoryID)
		if err == nil && category != nil {
			categoryName = category.Name
		}
	}

	return &dto.ScenicSpotResponse{
		ID:             spot.ID,
		Name:           spot.Name,
		CategoryID:     spot.CategoryID,
		CategoryName:   categoryName,
		Address:        spot.Address,
		LocationArea:   spot.LocationArea,
		Description:    spot.Description,
		OpeningHours:   spot.OpeningHours,
		Price:          spot.Price,
		TicketInfo:     spot.TicketInfo,
		ScenicFeatures: spot.ScenicFeatures,
		TransportInfo:  spot.TransportInfo,
		Images:         spot.Images,
		Rating:         spot.Rating,
		ReviewCount:    spot.ReviewCount,
		ViewCount:      spot.ViewCount,
		Latitude:       spot.Latitude,
		Longitude:      spot.Longitude,
		NearbyFacilities: spot.NearbyFacilities,
		TransitRoutes:    spot.TransitRoutes,
		CreatedAt:      spot.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      spot.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// entityToListItem converts a scenic spot entity to list item
func (s *ScenicSpotAppService) entityToListItem(ctx context.Context, spot *entity.ScenicSpot) (*dto.ScenicSpotListItem, error) {
	// Get category name if possible
	var categoryName string
	if spot.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(ctx, spot.CategoryID)
		if err == nil && category != nil {
			categoryName = category.Name
		}
	}

	// Create thumbnail image from first image if available
	var images []string
	if len(spot.Images) > 0 {
		if len(spot.Images) > 3 {
			images = spot.Images[:3]
		} else {
			images = spot.Images
		}
	}

	return &dto.ScenicSpotListItem{
		ID:             spot.ID,
		Name:           spot.Name,
		CategoryID:     spot.CategoryID,
		CategoryName:   categoryName,
		Address:        spot.Address,
		LocationArea:   spot.LocationArea,
		Price:          spot.Price,
		Images:         images,
		Rating:         spot.Rating,
		ReviewCount:    spot.ReviewCount,
		ViewCount:      spot.ViewCount,
		ScenicFeatures: spot.ScenicFeatures,
	}, nil
}
