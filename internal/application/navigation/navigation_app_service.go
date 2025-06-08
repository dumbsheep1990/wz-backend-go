package navigation

import (
	"context"

	"wz-backend-go/internal/domain/navigation/dto"
	"wz-backend-go/internal/domain/navigation/entity"
	"wz-backend-go/internal/domain/navigation/repository"
	"wz-backend-go/internal/domain/navigation/service"
)

// NavigationAppService handles application-level operations for the navigation system
type NavigationAppService struct {
	navigationService *service.NavigationService
	categoryRepo      repository.CategoryRepository
	websiteRepo       repository.WebsiteRepository
}

// NewNavigationAppService creates a new instance of NavigationAppService
func NewNavigationAppService(
	navigationService *service.NavigationService,
	categoryRepo repository.CategoryRepository,
	websiteRepo repository.WebsiteRepository,
) *NavigationAppService {
	return &NavigationAppService{
		navigationService: navigationService,
		categoryRepo:      categoryRepo,
		websiteRepo:       websiteRepo,
	}
}

// CreateCategory creates a new navigation category
func (s *NavigationAppService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryDTO, error) {
	category := entity.NewCategory(
		req.Name,
		req.DisplayName,
		req.Description,
		req.IconURL,
		req.SortOrder,
	)

	if err := s.navigationService.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	websiteCount, _ := s.categoryRepo.CountWebsites(ctx, category.ID)
	
	return &dto.CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		DisplayName: category.DisplayName,
		Description: category.Description,
		IconURL:     category.IconURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		WebsiteCount: websiteCount,
	}, nil
}

// UpdateCategory updates an existing category
func (s *NavigationAppService) UpdateCategory(ctx context.Context, req *dto.UpdateCategoryRequest) (*dto.CategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	category.UpdateDetails(req.DisplayName, req.Description, req.IconURL)
	category.UpdateSortOrder(req.SortOrder)
	
	if req.IsActive {
		category.Activate()
	} else {
		category.Deactivate()
	}

	if err := s.navigationService.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	websiteCount, _ := s.categoryRepo.CountWebsites(ctx, category.ID)
	
	return &dto.CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		DisplayName: category.DisplayName,
		Description: category.Description,
		IconURL:     category.IconURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		WebsiteCount: websiteCount,
	}, nil
}

// GetCategories gets all categories with optional filtering
func (s *NavigationAppService) GetCategories(ctx context.Context, activeOnly bool) (*dto.CategoryListResponse, error) {
	var categories []*entity.Category
	var err error
	
	if activeOnly {
		categories, err = s.categoryRepo.FindActive(ctx)
	} else {
		categories, err = s.categoryRepo.FindAllSorted(ctx)
	}
	
	if err != nil {
		return nil, err
	}
	
	dtos := make([]dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		websiteCount, _ := s.categoryRepo.CountWebsites(ctx, category.ID)
		
		dtos = append(dtos, dto.CategoryDTO{
			ID:          category.ID,
			Name:        category.Name,
			DisplayName: category.DisplayName,
			Description: category.Description,
			IconURL:     category.IconURL,
			SortOrder:   category.SortOrder,
			IsActive:    category.IsActive,
			WebsiteCount: websiteCount,
		})
	}
	
	return &dto.CategoryListResponse{
		Categories: dtos,
		Total:      len(dtos),
	}, nil
}

// GetCategoryByID gets a category by ID with its websites
func (s *NavigationAppService) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryDetailResponse, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	websites, err := s.websiteRepo.FindByCategorySorted(ctx, id)
	if err != nil {
		return nil, err
	}
	
	websiteDTOs := make([]dto.WebsiteDTO, 0, len(websites))
	for _, website := range websites {
		websiteDTOs = append(websiteDTOs, toWebsiteDTO(website))
	}
	
	websiteCount, _ := s.categoryRepo.CountWebsites(ctx, category.ID)
	
	return &dto.CategoryDetailResponse{
		Category: dto.CategoryDTO{
			ID:          category.ID,
			Name:        category.Name,
			DisplayName: category.DisplayName,
			Description: category.Description,
			IconURL:     category.IconURL,
			SortOrder:   category.SortOrder,
			IsActive:    category.IsActive,
			WebsiteCount: websiteCount,
		},
		Websites: websiteDTOs,
	}, nil
}

// CreateWebsite creates a new website in a category
func (s *NavigationAppService) CreateWebsite(ctx context.Context, req *dto.CreateWebsiteRequest) (*dto.WebsiteDTO, error) {
	website := entity.NewWebsite(
		req.CategoryID,
		req.Name,
		req.URL,
		req.Description,
		req.IconURL,
		req.SortOrder,
		req.Tags,
	)
	
	if err := s.navigationService.AddWebsiteToCategory(ctx, website); err != nil {
		return nil, err
	}
	
	return toWebsiteDTO(website), nil
}

// UpdateWebsite updates an existing website
func (s *NavigationAppService) UpdateWebsite(ctx context.Context, req *dto.UpdateWebsiteRequest) (*dto.WebsiteDTO, error) {
	website, err := s.websiteRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	
	// Update category if changed
	if req.CategoryID != "" && req.CategoryID != website.CategoryID {
		website.CategoryID = req.CategoryID
	}
	
	website.UpdateDetails(req.Name, req.URL, req.Description, req.IconURL, req.Tags)
	website.UpdateSortOrder(req.SortOrder)
	
	if req.IsActive {
		website.Activate()
	} else {
		website.Deactivate()
	}
	
	if req.IsFeatured {
		website.MarkAsFeatured()
	} else {
		website.UnmarkAsFeatured()
	}
	
	if err := s.websiteRepo.Update(ctx, website); err != nil {
		return nil, err
	}
	
	return toWebsiteDTO(website), nil
}

// GetWebsites gets all websites with optional filtering
func (s *NavigationAppService) GetWebsites(ctx context.Context, categoryID string, featuredOnly bool) (*dto.WebsiteListResponse, error) {
	var websites []*entity.Website
	var err error
	
	if categoryID != "" {
		websites, err = s.websiteRepo.FindByCategorySorted(ctx, categoryID)
	} else if featuredOnly {
		websites, err = s.websiteRepo.FindFeatured(ctx)
	} else {
		websites, err = s.websiteRepo.FindAll(ctx)
	}
	
	if err != nil {
		return nil, err
	}
	
	websiteDTOs := make([]dto.WebsiteDTO, 0, len(websites))
	for _, website := range websites {
		websiteDTOs = append(websiteDTOs, toWebsiteDTO(website))
	}
	
	return &dto.WebsiteListResponse{
		Websites: websiteDTOs,
		Total:    len(websiteDTOs),
	}, nil
}

// GetPopularWebsites gets popular websites by view count
func (s *NavigationAppService) GetPopularWebsites(ctx context.Context, limit int) (*dto.WebsiteListResponse, error) {
	websites, err := s.websiteRepo.FindPopular(ctx, limit)
	if err != nil {
		return nil, err
	}
	
	websiteDTOs := make([]dto.WebsiteDTO, 0, len(websites))
	for _, website := range websites {
		websiteDTOs = append(websiteDTOs, toWebsiteDTO(website))
	}
	
	return &dto.WebsiteListResponse{
		Websites: websiteDTOs,
		Total:    len(websiteDTOs),
	}, nil
}

// ReorderWebsites reorders websites within a category
func (s *NavigationAppService) ReorderWebsites(ctx context.Context, req *dto.ReorderWebsitesRequest) error {
	return s.navigationService.ReorderWebsites(ctx, req.CategoryID, req.WebsiteIDs)
}

// ReorderCategories reorders categories
func (s *NavigationAppService) ReorderCategories(ctx context.Context, req *dto.ReorderCategoriesRequest) error {
	return s.navigationService.ReorderCategories(ctx, req.CategoryIDs)
}

// TrackWebsiteView records a website view
func (s *NavigationAppService) TrackWebsiteView(ctx context.Context, websiteID string) error {
	return s.navigationService.TrackWebsiteView(ctx, websiteID)
}

// Helper function to convert entity to DTO
func toWebsiteDTO(website *entity.Website) dto.WebsiteDTO {
	return dto.WebsiteDTO{
		ID:          website.ID,
		CategoryID:  website.CategoryID,
		Name:        website.Name,
		URL:         website.URL,
		Description: website.Description,
		IconURL:     website.IconURL,
		SortOrder:   website.SortOrder,
		IsActive:    website.IsActive,
		IsNew:       website.IsNew,
		IsFeatured:  website.IsFeatured,
		ViewCount:   website.ViewCount,
		Tags:        website.Tags,
	}
}
