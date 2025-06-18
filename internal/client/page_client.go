package client

import (
	"context"

	"github.com/wz-backend-go/internal/types"
)

type PageClientInterface interface {
	GetPageList(ctx context.Context, req *types.GetPageListRequest) (*types.GetPageListResponse, error)
	CreatePage(ctx context.Context, req *types.CreatePageRequest) (*types.Page, error)
	UpdatePage(ctx context.Context, req *types.UpdatePageRequest) (*types.Page, error)
	DeletePage(ctx context.Context, req *types.DeletePageRequest) error
	GetPageDetail(ctx context.Context, req *types.GetPageDetailRequest) (*types.Page, error)
	TogglePageStatus(ctx context.Context, req *types.TogglePageStatusRequest) (*types.Page, error)
	PreviewPage(ctx context.Context, req *types.PreviewPageRequest) (*types.PagePreview, error)
	BatchUpdatePage(ctx context.Context, req *types.BatchUpdatePageRequest) error
}

type PageClient struct {
	endpoint string
}

func NewPageClient(endpoint string) PageClientInterface {
	return &PageClient{
		endpoint: endpoint,
	}
}

func (c *PageClient) GetPageList(ctx context.Context, req *types.GetPageListRequest) (*types.GetPageListResponse, error) {
	// TODO: 实现gRPC调用到页面微服务
	// 这里应该调用page microservice的gRPC接口
	return &types.GetPageListResponse{
		Pages: []*types.Page{},
		Total: 0,
	}, nil
}

func (c *PageClient) CreatePage(ctx context.Context, req *types.CreatePageRequest) (*types.Page, error) {
	// TODO: 实现gRPC调用到页面微服务
	return &types.Page{}, nil
}

func (c *PageClient) UpdatePage(ctx context.Context, req *types.UpdatePageRequest) (*types.Page, error) {
	// TODO: 实现gRPC调用到页面微服务
	return &types.Page{}, nil
}

func (c *PageClient) DeletePage(ctx context.Context, req *types.DeletePageRequest) error {
	// TODO: 实现gRPC调用到页面微服务
	return nil
}

func (c *PageClient) GetPageDetail(ctx context.Context, req *types.GetPageDetailRequest) (*types.Page, error) {
	// TODO: 实现gRPC调用到页面微服务
	return &types.Page{}, nil
}

func (c *PageClient) TogglePageStatus(ctx context.Context, req *types.TogglePageStatusRequest) (*types.Page, error) {
	// TODO: 实现gRPC调用到页面微服务
	return &types.Page{}, nil
}

func (c *PageClient) PreviewPage(ctx context.Context, req *types.PreviewPageRequest) (*types.PagePreview, error) {
	// TODO: 实现gRPC调用到页面微服务
	return &types.PagePreview{}, nil
}

func (c *PageClient) BatchUpdatePage(ctx context.Context, req *types.BatchUpdatePageRequest) error {
	// TODO: 实现gRPC调用到页面微服务
	return nil
} 