package service

import (
	"context"
	"errors"
	"wz-backend-go/internal/domain/navigation/entity"
	"wz-backend-go/internal/domain/navigation/repository"
)

// NavigationService handles domain logic for the navigation system
type NavigationService struct {
	categoryRepo repository.CategoryRepository
	websiteRepo  repository.WebsiteRepository
}

// NewNavigationService creates a new instance of NavigationService
func NewNavigationService(
	categoryRepo repository.CategoryRepository,
	websiteRepo repository.WebsiteRepository,
) *NavigationService {
	return &NavigationService{
		categoryRepo: categoryRepo,
		websiteRepo:  websiteRepo,
	}
}

// CreateCategory creates a new category
func (s *NavigationService) CreateCategory(ctx context.Context, category *entity.Category) error {
	existing, err := s.categoryRepo.FindByName(ctx, category.Name)
	if err == nil && existing != nil {
		return errors.New("category with this name already exists")
	}

	return s.categoryRepo.Save(ctx, category)
}

// UpdateCategory updates an existing category
func (s *NavigationService) UpdateCategory(ctx context.Context, category *entity.Category) error {
	existing, err := s.categoryRepo.FindByID(ctx, category.ID)
	if err != nil {
		return err
	}
	
	if existing == nil {
		return errors.New("category not found")
	}
	
	return s.categoryRepo.Update(ctx, category)
}

// AddWebsiteToCategory adds a new website to a category
func (s *NavigationService) AddWebsiteToCategory(ctx context.Context, website *entity.Website) error {
	// Check if category exists
	category, err := s.categoryRepo.FindByID(ctx, website.CategoryID)
	if err != nil {
		return err
	}
	
	if category == nil {
		return errors.New("category not found")
	}
	
	// Check if website with same URL already exists
	existing, err := s.websiteRepo.FindByURL(ctx, website.URL)
	if err == nil && existing != nil {
		return errors.New("website with this URL already exists")
	}
	
	return s.websiteRepo.Save(ctx, website)
}

// ReorderWebsites reorders websites within a category
func (s *NavigationService) ReorderWebsites(ctx context.Context, categoryID string, websiteIDs []string) error {
	// Check if category exists
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return err
	}
	
	if category == nil {
		return errors.New("category not found")
	}
	
	// Reorder websites by updating their sortOrder
	for i, id := range websiteIDs {
		website, err := s.websiteRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		
		if website == nil {
			return errors.New("website not found")
		}
		
		if website.CategoryID != categoryID {
			return errors.New("website does not belong to the specified category")
		}
		
		website.UpdateSortOrder(i)
		if err := s.websiteRepo.Update(ctx, website); err != nil {
			return err
		}
	}
	
	return nil
}

// ReorderCategories reorders categories
func (s *NavigationService) ReorderCategories(ctx context.Context, categoryIDs []string) error {
	for i, id := range categoryIDs {
		category, err := s.categoryRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		
		if category == nil {
			return errors.New("category not found")
		}
		
		category.UpdateSortOrder(i)
		if err := s.categoryRepo.Update(ctx, category); err != nil {
			return err
		}
	}
	
	return nil
}

// TrackeWebsiteView records a website view and increments the view count
func (s *NavigationService) TrackWebsiteView(ctx context.Context, websiteID string) error {
	return s.websiteRepo.IncrementViewCount(ctx, websiteID)
}
