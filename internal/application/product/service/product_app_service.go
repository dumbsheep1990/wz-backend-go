package service

import (
	"errors"

	"wz-backend-go/internal/application/product/dto"
	domainService "wz-backend-go/internal/domain/product/service"
	"wz-backend-go/internal/domain/product/valueobject"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// ProductApplicationService 商品应用服务
type ProductApplicationService struct {
	productDomainService *domainService.ProductDomainService
}

// NewProductApplicationService 创建商品应用服务
func NewProductApplicationService(productDomainService *domainService.ProductDomainService) *ProductApplicationService {
	return &ProductApplicationService{
		productDomainService: productDomainService,
	}
}

// CreateProduct 创建商品
func (s *ProductApplicationService) CreateProduct(req *dto.ProductCreateRequest, creatorID int64) (*dto.ProductResponse, error) {
	// 转换请求参数为值对象
	name, err := valueobject.NewProductName(req.Name)
	if err != nil {
		return nil, err
	}

	description, err := valueobject.NewProductDescription(req.Description)
	if err != nil {
		return nil, err
	}

	sku, err := valueobject.NewProductSKU(req.SKU)
	if err != nil {
		return nil, err
	}

	price, err := valueobject.NewPrice(req.Price)
	if err != nil {
		return nil, err
	}

	// 创建者ID验证和转换
	if creatorID <= 0 {
		return nil, errors.New("无效的创建者ID")
	}

	userID := uservo.NewUserID(creatorID)

	// 调用领域服务创建商品
	product, err := s.productDomainService.CreateProduct(
		name,
		description,
		sku,
		price,
		req.Stock,
		req.CategoryID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "商品创建成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// UpdateProduct 更新商品
func (s *ProductApplicationService) UpdateProduct(id int64, req *dto.ProductUpdateRequest) (*dto.ProductResponse, error) {
	// 转换请求参数为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	name, err := valueobject.NewProductName(req.Name)
	if err != nil {
		return nil, err
	}

	description, err := valueobject.NewProductDescription(req.Description)
	if err != nil {
		return nil, err
	}

	price, err := valueobject.NewPrice(req.Price)
	if err != nil {
		return nil, err
	}

	// 调用领域服务更新商品
	product, err := s.productDomainService.UpdateProduct(
		productID,
		name,
		description,
		price,
		req.CategoryID,
	)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "商品更新成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// GetProductByID 根据ID获取商品
func (s *ProductApplicationService) GetProductByID(id int64) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 获取商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("商品不存在")
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "success",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// GetProductBySKU 根据SKU获取商品
func (s *ProductApplicationService) GetProductBySKU(skuStr string) (*dto.ProductResponse, error) {
	// 转换为值对象
	sku, err := valueobject.NewProductSKU(skuStr)
	if err != nil {
		return nil, err
	}

	// 获取商品
	product, err := s.productDomainService.GetProductBySKU(sku)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("商品不存在")
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "success",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// PublishProduct 发布商品
func (s *ProductApplicationService) PublishProduct(id int64) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 发布商品
	if err := s.productDomainService.PublishProduct(productID); err != nil {
		return nil, err
	}

	// 获取更新后的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "商品发布成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// UnpublishProduct 下架商品
func (s *ProductApplicationService) UnpublishProduct(id int64) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 下架商品
	if err := s.productDomainService.UnpublishProduct(productID); err != nil {
		return nil, err
	}

	// 获取更新后的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "商品下架成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// DeleteProduct 删除商品
func (s *ProductApplicationService) DeleteProduct(id int64) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 获取要删除的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("商品不存在")
	}

	// 保存商品信息用于响应
	productDTO := dto.ToProductDTO(product)

	// 删除商品
	if err := s.productDomainService.DeleteProduct(productID); err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "商品删除成功",
		Data:    productDTO,
	}, nil
}

// UpdateStock 更新库存
func (s *ProductApplicationService) UpdateStock(id int64, req *dto.StockUpdateRequest) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 更新库存
	if err := s.productDomainService.UpdateStock(productID, req.Stock, req.Reason); err != nil {
		return nil, err
	}

	// 获取更新后的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "库存更新成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// DecrementStock 减少库存
func (s *ProductApplicationService) DecrementStock(id int64, req *dto.StockDecrementRequest) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 减少库存
	if err := s.productDomainService.DecrementStock(productID, req.Quantity, req.Reason); err != nil {
		return nil, err
	}

	// 获取更新后的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "库存减少成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// IncrementStock 增加库存
func (s *ProductApplicationService) IncrementStock(id int64, req *dto.StockIncrementRequest) (*dto.ProductResponse, error) {
	// 转换为值对象
	if err := valueobject.ValidateProductID(id); err != nil {
		return nil, err
	}

	productID := valueobject.NewProductID(id)

	// 增加库存
	if err := s.productDomainService.IncrementStock(productID, req.Quantity, req.Reason); err != nil {
		return nil, err
	}

	// 获取更新后的商品
	product, err := s.productDomainService.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.ProductResponse{
		Code:    0,
		Message: "库存增加成功",
		Data:    dto.ToProductDTO(product),
	}, nil
}

// GetProducts 分页获取商品列表
func (s *ProductApplicationService) GetProducts(page, pageSize int) (*dto.ProductsResponse, error) {
	// 调用领域服务获取商品列表
	products, total, err := s.productDomainService.GetProducts(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToProductsResponse(products, total), nil
}

// GetProductsByCategory 获取分类下的商品
func (s *ProductApplicationService) GetProductsByCategory(categoryID int64, page, pageSize int) (*dto.ProductsResponse, error) {
	// 调用领域服务获取分类下的商品
	products, total, err := s.productDomainService.GetProductsByCategory(categoryID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToProductsResponse(products, total), nil
}

// GetProductsByCreator 获取创建者的商品
func (s *ProductApplicationService) GetProductsByCreator(creatorID int64, page, pageSize int) (*dto.ProductsResponse, error) {
	// 调用领域服务获取创建者的商品
	products, total, err := s.productDomainService.GetProductsByCreator(creatorID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToProductsResponse(products, total), nil
}

// SearchProducts 搜索商品
func (s *ProductApplicationService) SearchProducts(keyword string, page, pageSize int) (*dto.ProductsResponse, error) {
	// 调用领域服务搜索商品
	products, total, err := s.productDomainService.SearchProducts(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToProductsResponse(products, total), nil
}

// GetActiveProducts 获取上架商品
func (s *ProductApplicationService) GetActiveProducts(page, pageSize int) (*dto.ProductsResponse, error) {
	// 调用领域服务获取上架商品
	products, total, err := s.productDomainService.GetActiveProducts(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return dto.ToProductsResponse(products, total), nil
}
