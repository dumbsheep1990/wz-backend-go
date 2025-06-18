package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/component"
	"github.com/wz-backend-go/internal/types"
)

type ComponentClientInterface interface {
	// gRPC proto methods
	GetComponentStats(ctx context.Context) (*pb.ComponentStatsResponse, error)
	ListComponents(ctx context.Context, pageSize int32, pageNumber int32, categoryId string, filter string) (*pb.ListComponentsResponse, error)
	GetComponent(ctx context.Context, id string) (*pb.Component, error)
	CreateComponent(ctx context.Context, req *pb.CreateComponentRequest) (*pb.Component, error)
	UpdateComponent(ctx context.Context, req *pb.UpdateComponentRequest) (*pb.Component, error)
	DeleteComponent(ctx context.Context, id string) (*pb.DeleteComponentResponse, error)
	
	// Admin interface methods - wrapper for gRPC calls
	GetComponentList(ctx context.Context, req *types.GetComponentListRequest) (*types.GetComponentListResponse, error)
	CreateComponentAdmin(ctx context.Context, req *types.CreateComponentRequest) (*types.Component, error)
	UpdateComponentAdmin(ctx context.Context, req *types.UpdateComponentRequest) (*types.Component, error)
	DeleteComponentAdmin(ctx context.Context, req *types.DeleteComponentRequest) error
	GetComponentDetail(ctx context.Context, req *types.GetComponentDetailRequest) (*types.Component, error)
	PreviewComponent(ctx context.Context, req *types.PreviewComponentRequest) (*types.ComponentPreview, error)
	ImportComponent(ctx context.Context, req *types.ImportComponentRequest) (*types.Component, error)
	GetComponentTypes(ctx context.Context, req *types.GetComponentTypesRequest) ([]*types.ComponentType, error)
}

// ComponentClient 是组件服务的客户端封装
type ComponentClient struct {
	client pb.ComponentServiceClient
	conn   *grpc.ClientConn
}

// NewComponentClient 创建一个新的组件服务客户端
func NewComponentClient(serviceAddr string) (*ComponentClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50056" // 默认地址
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("无法连接到组件服务: %v", err)
	}

	client := pb.NewComponentServiceClient(conn)
	return &ComponentClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *ComponentClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetComponentStats 获取组件统计数据
func (c *ComponentClient) GetComponentStats(ctx context.Context) (*pb.ComponentStatsResponse, error) {
	return c.client.GetComponentStats(ctx, &pb.GetComponentStatsRequest{})
}

// ListComponents 获取组件列表
func (c *ComponentClient) ListComponents(ctx context.Context, pageSize int32, pageNumber int32, categoryId string, filter string) (*pb.ListComponentsResponse, error) {
	return c.client.ListComponents(ctx, &pb.ListComponentsRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		CategoryId: categoryId,
		Filter:     filter,
	})
}

// GetComponent 获取组件详情
func (c *ComponentClient) GetComponent(ctx context.Context, id string) (*pb.Component, error) {
	return c.client.GetComponent(ctx, &pb.GetComponentRequest{
		Id: id,
	})
}

// CreateComponent 创建组件
func (c *ComponentClient) CreateComponent(ctx context.Context, req *pb.CreateComponentRequest) (*pb.Component, error) {
	return c.client.CreateComponent(ctx, req)
}

// UpdateComponent 更新组件
func (c *ComponentClient) UpdateComponent(ctx context.Context, req *pb.UpdateComponentRequest) (*pb.Component, error) {
	return c.client.UpdateComponent(ctx, req)
}

// DeleteComponent 删除组件
func (c *ComponentClient) DeleteComponent(ctx context.Context, id string) (*pb.DeleteComponentResponse, error) {
	return c.client.DeleteComponent(ctx, &pb.DeleteComponentRequest{
		Id: id,
	})
}

// ListComponentCategories 获取组件分类列表
func (c *ComponentClient) ListComponentCategories(ctx context.Context) (*pb.ListComponentCategoriesResponse, error) {
	return c.client.ListComponentCategories(ctx, &pb.ListComponentCategoriesRequest{})
}

// GetComponentCategory 获取组件分类详情
func (c *ComponentClient) GetComponentCategory(ctx context.Context, id string) (*pb.ComponentCategory, error) {
	return c.client.GetComponentCategory(ctx, &pb.GetComponentCategoryRequest{
		Id: id,
	})
}

// CreateComponentCategory 创建组件分类
func (c *ComponentClient) CreateComponentCategory(ctx context.Context, req *pb.CreateComponentCategoryRequest) (*pb.ComponentCategory, error) {
	return c.client.CreateComponentCategory(ctx, req)
}

// UpdateComponentCategory 更新组件分类
func (c *ComponentClient) UpdateComponentCategory(ctx context.Context, req *pb.UpdateComponentCategoryRequest) (*pb.ComponentCategory, error) {
	return c.client.UpdateComponentCategory(ctx, req)
}

// DeleteComponentCategory 删除组件分类
func (c *ComponentClient) DeleteComponentCategory(ctx context.Context, id string) (*pb.DeleteComponentCategoryResponse, error) {
	return c.client.DeleteComponentCategory(ctx, &pb.DeleteComponentCategoryRequest{
		Id: id,
	})
}

// GetComponentVersion 获取组件版本
func (c *ComponentClient) GetComponentVersion(ctx context.Context, componentId string, version string) (*pb.ComponentVersion, error) {
	return c.client.GetComponentVersion(ctx, &pb.GetComponentVersionRequest{
		ComponentId: componentId,
		Version:     version,
	})
}

// ListComponentVersions 获取组件版本列表
func (c *ComponentClient) ListComponentVersions(ctx context.Context, componentId string) (*pb.ListComponentVersionsResponse, error) {
	return c.client.ListComponentVersions(ctx, &pb.ListComponentVersionsRequest{
		ComponentId: componentId,
	})
}

// CreateComponentVersion 创建组件版本
func (c *ComponentClient) CreateComponentVersion(ctx context.Context, req *pb.CreateComponentVersionRequest) (*pb.ComponentVersion, error) {
	return c.client.CreateComponentVersion(ctx, req)
}

// PublishComponentVersion 发布组件版本
func (c *ComponentClient) PublishComponentVersion(ctx context.Context, componentId string, version string) (*pb.PublishComponentVersionResponse, error) {
	return c.client.PublishComponentVersion(ctx, &pb.PublishComponentVersionRequest{
		ComponentId: componentId,
		Version:     version,
	})
}

// Admin interface methods - wrapper for gRPC calls
func (c *ComponentClient) GetComponentList(ctx context.Context, req *types.GetComponentListRequest) (*types.GetComponentListResponse, error) {
	// Call the existing gRPC method
	resp, err := c.ListComponents(ctx, int32(req.PageSize), int32(req.Page), "", req.Name)
	if err != nil {
		return nil, err
	}
	
	// Convert response
	var components []*types.Component
	for _, comp := range resp.Components {
		components = append(components, &types.Component{
			ID:          comp.Id,
			Name:        comp.Name,
			Type:        comp.Category,
			Description: comp.Description,
			Code:        comp.Content,
		})
	}
	
	return &types.GetComponentListResponse{
		Components: components,
		Total:      int(resp.Total),
	}, nil
}

func (c *ComponentClient) CreateComponentAdmin(ctx context.Context, req *types.CreateComponentRequest) (*types.Component, error) {
	// Call the existing gRPC method
	pbReq := &pb.CreateComponentRequest{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Type,
		Content:     req.Code,
	}
	
	component, err := c.CreateComponent(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	
	return &types.Component{
		ID:          component.Id,
		Name:        component.Name,
		Type:        component.Category,
		Description: component.Description,
		Code:        component.Content,
	}, nil
}

func (c *ComponentClient) UpdateComponentAdmin(ctx context.Context, req *types.UpdateComponentRequest) (*types.Component, error) {
	// Call the existing gRPC method
	pbReq := &pb.UpdateComponentRequest{
		Id:          fmt.Sprintf("%d", req.ID),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Type,
		Content:     req.Code,
	}
	
	component, err := c.UpdateComponent(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	
	return &types.Component{
		ID:          component.Id,
		Name:        component.Name,
		Type:        component.Category,
		Description: component.Description,
		Code:        component.Content,
	}, nil
}

func (c *ComponentClient) DeleteComponentAdmin(ctx context.Context, req *types.DeleteComponentRequest) error {
	_, err := c.DeleteComponent(ctx, fmt.Sprintf("%d", req.ID))
	return err
}

func (c *ComponentClient) GetComponentDetail(ctx context.Context, req *types.GetComponentDetailRequest) (*types.Component, error) {
	component, err := c.GetComponent(ctx, fmt.Sprintf("%d", req.ID))
	if err != nil {
		return nil, err
	}
	
	return &types.Component{
		ID:          component.Id,
		Name:        component.Name,
		Type:        component.Category,
		Description: component.Description,
		Code:        component.Content,
	}, nil
}

func (c *ComponentClient) PreviewComponent(ctx context.Context, req *types.PreviewComponentRequest) (*types.ComponentPreview, error) {
	// TODO: 实现组件预览逻辑
	return &types.ComponentPreview{
		HTML: "<div>Preview for component " + fmt.Sprintf("%d", req.ID) + "</div>",
	}, nil
}

func (c *ComponentClient) ImportComponent(ctx context.Context, req *types.ImportComponentRequest) (*types.Component, error) {
	// TODO: 实现组件导入逻辑
	return &types.Component{}, nil
}

func (c *ComponentClient) GetComponentTypes(ctx context.Context, req *types.GetComponentTypesRequest) ([]*types.ComponentType, error) {
	// Call the existing gRPC method for categories
	resp, err := c.ListComponentCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	var types []*types.ComponentType
	for _, cat := range resp.Categories {
		types = append(types, &types.ComponentType{
			Name:        cat.Name,
			Description: cat.Description,
		})
	}
	
	return types, nil
}
