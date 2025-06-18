package client

import (
	"context"

	"github.com/wz-backend-go/internal/types"
)

type LinkClientInterface interface {
	GetLinkList(ctx context.Context, req *types.GetLinkListRequest) (*types.GetLinkListResponse, error)
	CreateLink(ctx context.Context, req *types.CreateLinkRequest) (*types.Link, error)
	UpdateLink(ctx context.Context, req *types.UpdateLinkRequest) (*types.Link, error)
	DeleteLink(ctx context.Context, req *types.DeleteLinkRequest) error
	GetLinkDetail(ctx context.Context, req *types.GetLinkDetailRequest) (*types.Link, error)
	VerifyLink(ctx context.Context, req *types.VerifyLinkRequest) (*types.LinkVerifyResult, error)
	BatchVerifyLinks(ctx context.Context, req *types.BatchVerifyLinksRequest) (*types.BatchVerifyResult, error)
	GetLinkCategories(ctx context.Context, req *types.GetLinkCategoriesRequest) ([]*types.LinkCategory, error)
	UpdateLinkSort(ctx context.Context, req *types.UpdateLinkSortRequest) error
}

type LinkClient struct {
	endpoint string
}

func NewLinkClient(endpoint string) LinkClientInterface {
	return &LinkClient{
		endpoint: endpoint,
	}
}

func (c *LinkClient) GetLinkList(ctx context.Context, req *types.GetLinkListRequest) (*types.GetLinkListResponse, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.GetLinkListResponse{
		Links: []*types.Link{},
		Total: 0,
	}, nil
}

func (c *LinkClient) CreateLink(ctx context.Context, req *types.CreateLinkRequest) (*types.Link, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.Link{}, nil
}

func (c *LinkClient) UpdateLink(ctx context.Context, req *types.UpdateLinkRequest) (*types.Link, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.Link{}, nil
}

func (c *LinkClient) DeleteLink(ctx context.Context, req *types.DeleteLinkRequest) error {
	// TODO: 实现gRPC调用到链接微服务
	return nil
}

func (c *LinkClient) GetLinkDetail(ctx context.Context, req *types.GetLinkDetailRequest) (*types.Link, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.Link{}, nil
}

func (c *LinkClient) VerifyLink(ctx context.Context, req *types.VerifyLinkRequest) (*types.LinkVerifyResult, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.LinkVerifyResult{}, nil
}

func (c *LinkClient) BatchVerifyLinks(ctx context.Context, req *types.BatchVerifyLinksRequest) (*types.BatchVerifyResult, error) {
	// TODO: 实现gRPC调用到链接微服务
	return &types.BatchVerifyResult{}, nil
}

func (c *LinkClient) GetLinkCategories(ctx context.Context, req *types.GetLinkCategoriesRequest) ([]*types.LinkCategory, error) {
	// TODO: 实现gRPC调用到链接微服务
	return []*types.LinkCategory{}, nil
}

func (c *LinkClient) UpdateLinkSort(ctx context.Context, req *types.UpdateLinkSortRequest) error {
	// TODO: 实现gRPC调用到链接微服务
	return nil
} 