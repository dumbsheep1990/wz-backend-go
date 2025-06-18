package client

import (
	"context"
)

// NavigationClient Navigation微服务的客户端接口
type NavigationClient interface {
	// 分类管理
	GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error)
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryDetailResponse, error)
	DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteResponse, error)
	
	// 网站管理
	GetWebsites(ctx context.Context, req *GetWebsitesRequest) (*GetWebsitesResponse, error)
	CreateWebsite(ctx context.Context, req *CreateWebsiteRequest) (*WebsiteDetailResponse, error)
	DeleteWebsite(ctx context.Context, req *DeleteWebsiteRequest) (*DeleteResponse, error)
}

// 请求结构体
type GetCategoriesRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
	SortOrder   int32  `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

type DeleteCategoryRequest struct {
	Id string `json:"id"`
}

type GetWebsitesRequest struct {
	Page       int32  `json:"page"`
	PageSize   int32  `json:"pageSize"`
	Type       string `json:"type"`       // footer, sidebar, main
	CategoryId string `json:"categoryId"`
}

type CreateWebsiteRequest struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
	SortOrder   int32  `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
	CategoryId  string `json:"categoryId"`
	Type        string `json:"type"`
}

type DeleteWebsiteRequest struct {
	Id string `json:"id"`
}

// 响应结构体
type GetCategoriesResponse struct {
	Total      int64                   `json:"total"`
	Categories []*CategoryDetailResponse `json:"categories"`
}

type CategoryDetailResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
	SortOrder   int32  `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type GetWebsitesResponse struct {
	Total    int64                    `json:"total"`
	Websites []*WebsiteDetailResponse `json:"websites"`
}

type WebsiteDetailResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
	SortOrder   int32  `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
	CategoryId  string `json:"categoryId"`
	Type        string `json:"type"`
	ViewCount   int64  `json:"viewCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HTTPNavigationClient HTTP方式调用Navigation微服务的客户端实现
type HTTPNavigationClient struct {
	baseURL string
}

// NewHTTPNavigationClient 创建HTTP Navigation客户端
func NewHTTPNavigationClient(baseURL string) NavigationClient {
	return &HTTPNavigationClient{
		baseURL: baseURL,
	}
}

// 实现NavigationClient接口
func (c *HTTPNavigationClient) GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &GetCategoriesResponse{
		Total:      0,
		Categories: make([]*CategoryDetailResponse, 0),
	}, nil
}

func (c *HTTPNavigationClient) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryDetailResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &CategoryDetailResponse{
		Id:          "generated-id",
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IconUrl:     req.IconUrl,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}, nil
}

func (c *HTTPNavigationClient) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &DeleteResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil
}

func (c *HTTPNavigationClient) GetWebsites(ctx context.Context, req *GetWebsitesRequest) (*GetWebsitesResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &GetWebsitesResponse{
		Total:    0,
		Websites: make([]*WebsiteDetailResponse, 0),
	}, nil
}

func (c *HTTPNavigationClient) CreateWebsite(ctx context.Context, req *CreateWebsiteRequest) (*WebsiteDetailResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &WebsiteDetailResponse{
		Id:          "generated-id",
		Name:        req.Name,
		Url:         req.Url,
		Description: req.Description,
		IconUrl:     req.IconUrl,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		CategoryId:  req.CategoryId,
		Type:        req.Type,
	}, nil
}

func (c *HTTPNavigationClient) DeleteWebsite(ctx context.Context, req *DeleteWebsiteRequest) (*DeleteResponse, error) {
	// 简化实现，实际项目中需要调用HTTP API
	return &DeleteResponse{
		Success: true,
		Message: "Website deleted successfully",
	}, nil
} 