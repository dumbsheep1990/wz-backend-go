package commerce

import (
	"context"
	"math"
	"strings"
	"time"

	"wz-backend-go/internal/domain/commerce/dto"
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
	"wz-backend-go/internal/domain/commerce/service"
)

// ProductAppService handles product-related application use cases
type ProductAppService struct {
	commerceService   *service.CommerceService
	productRepository repository.ProductRepository
	categoryRepository repository.CategoryRepository
	storeRepository   repository.StoreRepository
}

// NewProductAppService creates a new instance of ProductAppService
func NewProductAppService(
	commerceService *service.CommerceService,
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	storeRepo repository.StoreRepository,
) *ProductAppService {
	return &ProductAppService{
		commerceService:    commerceService,
		productRepository:  productRepo,
		categoryRepository: categoryRepo,
		storeRepository:    storeRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductAppService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.commerceService.CreateProduct(
		ctx,
		req.Name,
		req.Description,
		req.Price,
		req.CategoryID,
		req.StoreID,
		req.Region,
		req.District,
		req.ImageURLs,
		req.Specifications,
		req.Stock,
		req.Tags,
	)
	
	if err != nil {
		return nil, err
	}
	
	return s.productToResponse(ctx, product)
}

// UpdateProduct updates an existing product
func (s *ProductAppService) UpdateProduct(ctx context.Context, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// First get the product
	product, err := s.productRepository.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	
	if product == nil {
		return nil, service.ErrProductNotFound
	}
	
	// Update basic product details
	product, err = s.commerceService.UpdateProduct(
		ctx,
		req.ID,
		req.Name,
		req.Description,
		req.Price,
		req.ImageURLs,
		req.Specifications,
		req.Stock,
		req.Tags,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle optional boolean values
	if req.IsActive != nil {
		if *req.IsActive {
			product.Activate()
		} else {
			product.Deactivate()
		}
	}
	
	if req.IsFeatured != nil {
		if *req.IsFeatured {
			product.MarkAsFeatured()
		} else {
			product.UnmarkAsFeatured()
		}
	}
	
	if req.IsNew != nil {
		product.IsNew = *req.IsNew
	}
	
	// Marshal for storage
	if err := product.Marshal(); err != nil {
		return nil, err
	}
	
	// Save changes
	if err := s.productRepository.Update(ctx, product); err != nil {
		return nil, err
	}
	
	// Unmarshal for response
	if err := product.Unmarshal(); err != nil {
		return nil, err
	}
	
	return s.productToResponse(ctx, product)
}

// GetProductByID retrieves a product by its ID
func (s *ProductAppService) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.productRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if product == nil {
		return nil, service.ErrProductNotFound
	}
	
	// Unmarshal JSON fields
	if err := product.Unmarshal(); err != nil {
		return nil, err
	}
	
	// Track product view
	go s.commerceService.TrackProductView(context.Background(), id)
	
	return s.productToResponse(ctx, product)
}

// GetProducts retrieves products based on filter criteria
func (s *ProductAppService) GetProducts(ctx context.Context, filter *dto.ProductFilterRequest) (*dto.ProductsResponse, error) {
	// Convert DTO filter to repository filter
	repoFilter := repository.ProductFilters{
		ActiveOnly:   filter.ActiveOnly,
		FeaturedOnly: filter.FeaturedOnly,
		NewOnly:      filter.NewOnly,
		MinPrice:     filter.MinPrice,
		MaxPrice:     filter.MaxPrice,
		Offset:       (filter.Page - 1) * filter.PageSize,
		Limit:        filter.PageSize,
		SortBy:       filter.SortBy,
		SortOrder:    filter.SortOrder,
	}
	
	// Parse tags if provided
	if filter.Tags != "" {
		repoFilter.Tags = strings.Split(filter.Tags, ",")
	}
	
	// Get products based on filter
	var products []*entity.Product
	var err error
	
	if filter.CategoryID != "" {
		products, err = s.productRepository.FindByCategory(ctx, filter.CategoryID, repoFilter)
	} else if filter.StoreID != "" {
		products, err = s.productRepository.FindByStore(ctx, filter.StoreID, repoFilter)
	} else if filter.Region != "" || filter.District != "" {
		products, err = s.productRepository.FindByRegion(ctx, filter.Region, "", filter.District, repoFilter)
	} else {
		products, err = s.productRepository.FindAll(ctx, repoFilter)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Get total count
	total, err := s.productRepository.Count(ctx, repoFilter)
	if err != nil {
		return nil, err
	}
	
	// Convert to response DTOs
	responseProducts := make([]*dto.ProductListResponse, 0, len(products))
	for _, product := range products {
		if err := product.Unmarshal(); err != nil {
			return nil, err
		}
		
		resp, err := s.productToListResponse(ctx, product)
		if err != nil {
			return nil, err
		}
		
		responseProducts = append(responseProducts, resp)
	}
	
	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))
	
	return &dto.ProductsResponse{
		Products:   responseProducts,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetFeaturedProducts retrieves featured products
func (s *ProductAppService) GetFeaturedProducts(ctx context.Context, limit int) ([]*dto.ProductListResponse, error) {
	products, err := s.productRepository.FindFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}
	
	return s.productsToListResponse(ctx, products)
}

// GetNewProducts retrieves new products
func (s *ProductAppService) GetNewProducts(ctx context.Context, limit int) ([]*dto.ProductListResponse, error) {
	products, err := s.productRepository.FindNew(ctx, limit)
	if err != nil {
		return nil, err
	}
	
	return s.productsToListResponse(ctx, products)
}

// GetPopularProducts retrieves popular products based on view count
func (s *ProductAppService) GetPopularProducts(ctx context.Context, limit int) ([]*dto.ProductListResponse, error) {
	products, err := s.productRepository.FindPopular(ctx, limit)
	if err != nil {
		return nil, err
	}
	
	return s.productsToListResponse(ctx, products)
}

// TrackProductView tracks a view for a product
func (s *ProductAppService) TrackProductView(ctx context.Context, id string) error {
	return s.commerceService.TrackProductView(ctx, id)
}

// SearchProducts searches for products based on a query string
func (s *ProductAppService) SearchProducts(ctx context.Context, query string, filter *dto.ProductFilterRequest) (*dto.ProductsResponse, error) {
	// Convert DTO filter to repository filter
	repoFilter := repository.ProductFilters{
		ActiveOnly:   filter.ActiveOnly,
		FeaturedOnly: filter.FeaturedOnly,
		NewOnly:      filter.NewOnly,
		MinPrice:     filter.MinPrice,
		MaxPrice:     filter.MaxPrice,
		Offset:       (filter.Page - 1) * filter.PageSize,
		Limit:        filter.PageSize,
		SortBy:       filter.SortBy,
		SortOrder:    filter.SortOrder,
	}
	
	// Parse tags if provided
	if filter.Tags != "" {
		repoFilter.Tags = strings.Split(filter.Tags, ",")
	}
	
	// Search products
	products, err := s.productRepository.Search(ctx, query, repoFilter)
	if err != nil {
		return nil, err
	}
	
	// Get total count of search results
	total, err := s.productRepository.Count(ctx, repoFilter)
	if err != nil {
		return nil, err
	}
	
	// Convert to response DTOs
	responseProducts, err := s.productsToListResponse(ctx, products)
	if err != nil {
		return nil, err
	}
	
	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))
	
	return &dto.ProductsResponse{
		Products:   responseProducts,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// productToResponse converts a product entity to a product response DTO
func (s *ProductAppService) productToResponse(ctx context.Context, product *entity.Product) (*dto.ProductResponse, error) {
	if product == nil {
		return nil, nil
	}
	
	response := &dto.ProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		FormattedPrice: formatPrice(product.Price),
		CategoryID:     product.CategoryID,
		StoreID:        product.StoreID,
		Region:         product.Region,
		District:       product.District,
		ImageURLs:      product.ImageURLs,
		Specifications: product.Specifications,
		Stock:          product.Stock,
		ViewCount:      product.ViewCount,
		SoldCount:      product.SoldCount,
		IsActive:       product.IsActive,
		IsFeatured:     product.IsFeatured,
		IsNew:          product.IsNew,
		Tags:           product.Tags,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      product.UpdatedAt.Format(time.RFC3339),
	}
	
	// Get category name
	category, err := s.categoryRepository.FindByID(ctx, product.CategoryID)
	if err == nil && category != nil {
		response.CategoryName = category.DisplayName
	}
	
	// Get store name
	store, err := s.storeRepository.FindByID(ctx, product.StoreID)
	if err == nil && store != nil {
		response.StoreName = store.Name
	}
	
	return response, nil
}

// productToListResponse converts a product entity to a product list response DTO
func (s *ProductAppService) productToListResponse(ctx context.Context, product *entity.Product) (*dto.ProductListResponse, error) {
	if product == nil {
		return nil, nil
	}
	
	response := &dto.ProductListResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    limitDescription(product.Description, 100),
		Price:          product.Price,
		FormattedPrice: formatPrice(product.Price),
		CategoryID:     product.CategoryID,
		StoreID:        product.StoreID,
		IsActive:       product.IsActive,
		IsFeatured:     product.IsFeatured,
		IsNew:          product.IsNew,
		ViewCount:      product.ViewCount,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
	}
	
	// Get primary image URL
	if len(product.ImageURLs) > 0 {
		response.ImageURL = product.ImageURLs[0]
	}
	
	// Get category name
	category, err := s.categoryRepository.FindByID(ctx, product.CategoryID)
	if err == nil && category != nil {
		response.CategoryName = category.DisplayName
	}
	
	// Get store name
	store, err := s.storeRepository.FindByID(ctx, product.StoreID)
	if err == nil && store != nil {
		response.StoreName = store.Name
	}
	
	return response, nil
}

// productsToListResponse converts multiple product entities to product list response DTOs
func (s *ProductAppService) productsToListResponse(ctx context.Context, products []*entity.Product) ([]*dto.ProductListResponse, error) {
	responses := make([]*dto.ProductListResponse, 0, len(products))
	
	for _, product := range products {
		if err := product.Unmarshal(); err != nil {
			return nil, err
		}
		
		resp, err := s.productToListResponse(ctx, product)
		if err != nil {
			return nil, err
		}
		
		responses = append(responses, resp)
	}
	
	return responses, nil
}

// Helper function to format price
func formatPrice(price float64) string {
	return "¥" + formatFloat(price, 2)
}

// Helper function to format float with precision
func formatFloat(num float64, precision int) string {
	format := "%."+string(rune('0'+precision))+"f"
	return fmt.Sprintf(format, num)
}

// Helper function to limit description length
func limitDescription(description string, maxLen int) string {
	if len(description) <= maxLen {
		return description
	}
	return description[:maxLen] + "..."
}
