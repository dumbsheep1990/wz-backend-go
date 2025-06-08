package commerce

import (
	"context"
	"math"
	"time"

	"wz-backend-go/internal/domain/commerce/dto"
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
	"wz-backend-go/internal/domain/commerce/service"
)

// StoreAppService handles store-related application use cases
type StoreAppService struct {
	commerceService   *service.CommerceService
	storeRepository   repository.StoreRepository
	productRepository repository.ProductRepository
}

// NewStoreAppService creates a new instance of StoreAppService
func NewStoreAppService(
	commerceService *service.CommerceService,
	storeRepo repository.StoreRepository,
	productRepo repository.ProductRepository,
) *StoreAppService {
	return &StoreAppService{
		commerceService:   commerceService,
		storeRepository:   storeRepo,
		productRepository: productRepo,
	}
}

// CreateStore creates a new store
func (s *StoreAppService) CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	store, err := s.commerceService.CreateStore(
		ctx,
		req.Name,
		req.Description,
		req.OwnerID,
		req.LogoURL,
		req.Province,
		req.City,
		req.District,
		req.Address,
		req.ContactName,
		req.Phone,
	)
	
	if err != nil {
		return nil, err
	}
	
	return s.storeToResponse(ctx, store)
}

// UpdateStore updates an existing store
func (s *StoreAppService) UpdateStore(ctx context.Context, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	// First check if store exists
	store, err := s.storeRepository.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	
	if store == nil {
		return nil, service.ErrStoreNotFound
	}
	
	// Update store details
	store, err = s.commerceService.UpdateStore(
		ctx,
		req.ID,
		req.Name,
		req.Description,
		req.LogoURL,
		req.Province,
		req.City,
		req.District,
		req.Address,
		req.ContactName,
		req.Phone,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle optional boolean values
	if req.IsActive != nil {
		if *req.IsActive {
			store.Activate()
		} else {
			store.Deactivate()
		}
		
		// Save changes
		if err := s.storeRepository.Update(ctx, store); err != nil {
			return nil, err
		}
	}
	
	return s.storeToResponse(ctx, store)
}

// GetStoreByID retrieves a store by its ID
func (s *StoreAppService) GetStoreByID(ctx context.Context, id string) (*dto.StoreResponse, error) {
	store, err := s.storeRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if store == nil {
		return nil, service.ErrStoreNotFound
	}
	
	return s.storeToResponse(ctx, store)
}

// GetStoresByOwner retrieves stores owned by a user
func (s *StoreAppService) GetStoresByOwner(ctx context.Context, ownerID string) ([]*dto.StoreResponse, error) {
	stores, err := s.storeRepository.FindByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*dto.StoreResponse, 0, len(stores))
	for _, store := range stores {
		resp, err := s.storeToResponse(ctx, store)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// GetStores retrieves stores based on filter criteria
func (s *StoreAppService) GetStores(ctx context.Context, filter *dto.StoreFilterRequest) (*dto.StoresResponse, error) {
	// Convert DTO filter to repository filter
	repoFilter := repository.StoreFilters{
		ActiveOnly: filter.ActiveOnly,
		Offset:     (filter.Page - 1) * filter.PageSize,
		Limit:      filter.PageSize,
		SortBy:     filter.SortBy,
		SortOrder:  filter.SortOrder,
	}
	
	// Get stores based on filter
	var stores []*entity.Store
	var err error
	
	if filter.Province != "" || filter.City != "" || filter.District != "" {
		stores, err = s.storeRepository.FindByRegion(ctx, filter.Province, filter.City, filter.District)
	} else {
		stores, err = s.storeRepository.FindAll(ctx, repoFilter)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Convert to list response DTOs
	responses := make([]*dto.StoreListResponse, 0, len(stores))
	for _, store := range stores {
		productCount, err := s.storeRepository.CountProducts(ctx, store.ID)
		if err != nil {
			productCount = 0 // default to 0 if there's an error
		}
		
		resp := &dto.StoreListResponse{
			ID:           store.ID,
			Name:         store.Name,
			Description:  limitDescription(store.Description, 100),
			LogoURL:      store.LogoURL,
			Province:     store.Province,
			City:         store.City,
			District:     store.District,
			Rating:       store.Rating,
			IsActive:     store.IsActive,
			ProductCount: productCount,
			CreatedAt:    store.CreatedAt.Format(time.RFC3339),
		}
		
		responses = append(responses, resp)
	}
	
	// Get total count for pagination
	var total int
	if len(responses) < filter.PageSize {
		total = (filter.Page - 1) * filter.PageSize + len(responses)
	} else {
		total = len(responses) * 10 // Approximation as we don't have a separate count query
	}
	
	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))
	
	return &dto.StoresResponse{
		Stores:     responses,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchStores searches for stores based on a query string
func (s *StoreAppService) SearchStores(ctx context.Context, query string) ([]*dto.StoreListResponse, error) {
	stores, err := s.storeRepository.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*dto.StoreListResponse, 0, len(stores))
	for _, store := range stores {
		productCount, err := s.storeRepository.CountProducts(ctx, store.ID)
		if err != nil {
			productCount = 0 // default to 0 if there's an error
		}
		
		resp := &dto.StoreListResponse{
			ID:           store.ID,
			Name:         store.Name,
			Description:  limitDescription(store.Description, 100),
			LogoURL:      store.LogoURL,
			Province:     store.Province,
			City:         store.City,
			District:     store.District,
			Rating:       store.Rating,
			IsActive:     store.IsActive,
			ProductCount: productCount,
			CreatedAt:    store.CreatedAt.Format(time.RFC3339),
		}
		
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// storeToResponse converts a store entity to a store response DTO
func (s *StoreAppService) storeToResponse(ctx context.Context, store *entity.Store) (*dto.StoreResponse, error) {
	if store == nil {
		return nil, nil
	}
	
	// Get product count
	productCount, err := s.storeRepository.CountProducts(ctx, store.ID)
	if err != nil {
		productCount = 0 // default to 0 if there's an error
	}
	
	response := &dto.StoreResponse{
		ID:           store.ID,
		Name:         store.Name,
		Description:  store.Description,
		OwnerID:      store.OwnerID,
		LogoURL:      store.LogoURL,
		Province:     store.Province,
		City:         store.City,
		District:     store.District,
		Address:      store.Address,
		ContactName:  store.ContactName,
		Phone:        store.Phone,
		Rating:       store.Rating,
		IsActive:     store.IsActive,
		ProductCount: productCount,
		CreatedAt:    store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    store.UpdatedAt.Format(time.RFC3339),
	}
	
	return response, nil
}
