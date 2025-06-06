package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/content"
)

// ContentClient 是内容服务的客户端封装
type ContentClient struct {
	client pb.ContentServiceClient
	conn   *grpc.ClientConn
}

// NewContentClient 创建一个新的内容服务客户端
func NewContentClient(serviceAddr string) (*ContentClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50052" // 默认地址
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
		return nil, fmt.Errorf("无法连接到内容服务: %v", err)
	}

	client := pb.NewContentServiceClient(conn)
	return &ContentClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *ContentClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetContentStats 获取内容统计数据
func (c *ContentClient) GetContentStats(ctx context.Context) (*pb.ContentStatsResponse, error) {
	return c.client.GetContentStats(ctx, &pb.GetContentStatsRequest{})
}

// ListContents 获取内容列表
func (c *ContentClient) ListContents(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListContentsResponse, error) {
	return c.client.ListContents(ctx, &pb.ListContentsRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetContent 获取内容详情
func (c *ContentClient) GetContent(ctx context.Context, id string) (*pb.Content, error) {
	return c.client.GetContent(ctx, &pb.GetContentRequest{
		Id: id,
	})
}

// CreateContent 创建内容
func (c *ContentClient) CreateContent(ctx context.Context, req *pb.CreateContentRequest) (*pb.Content, error) {
	return c.client.CreateContent(ctx, req)
}

// UpdateContent 更新内容
func (c *ContentClient) UpdateContent(ctx context.Context, req *pb.UpdateContentRequest) (*pb.Content, error) {
	return c.client.UpdateContent(ctx, req)
}

// DeleteContent 删除内容
func (c *ContentClient) DeleteContent(ctx context.Context, id string) (*pb.DeleteContentResponse, error) {
	return c.client.DeleteContent(ctx, &pb.DeleteContentRequest{
		Id: id,
	})
}

// ListCategories 获取内容分类列表
func (c *ContentClient) ListCategories(ctx context.Context, pageSize int32, pageNumber int32) (*pb.ListCategoriesResponse, error) {
	return c.client.ListCategories(ctx, &pb.ListCategoriesRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
}

// GetCategory 获取分类详情
func (c *ContentClient) GetCategory(ctx context.Context, id string) (*pb.Category, error) {
	return c.client.GetCategory(ctx, &pb.GetCategoryRequest{
		Id: id,
	})
}

// CreateCategory 创建分类
func (c *ContentClient) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	return c.client.CreateCategory(ctx, req)
}

// UpdateCategory 更新分类
func (c *ContentClient) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	return c.client.UpdateCategory(ctx, req)
}

// DeleteCategory 删除分类
func (c *ContentClient) DeleteCategory(ctx context.Context, id string) (*pb.DeleteCategoryResponse, error) {
	return c.client.DeleteCategory(ctx, &pb.DeleteCategoryRequest{
		Id: id,
	})
}

// ListTags 获取标签列表
func (c *ContentClient) ListTags(ctx context.Context, pageSize int32, pageNumber int32) (*pb.ListTagsResponse, error) {
	return c.client.ListTags(ctx, &pb.ListTagsRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
}

// GetTag 获取标签详情
func (c *ContentClient) GetTag(ctx context.Context, id string) (*pb.Tag, error) {
	return c.client.GetTag(ctx, &pb.GetTagRequest{
		Id: id,
	})
}

// CreateTag 创建标签
func (c *ContentClient) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.Tag, error) {
	return c.client.CreateTag(ctx, req)
}

// UpdateTag 更新标签
func (c *ContentClient) UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.Tag, error) {
	return c.client.UpdateTag(ctx, req)
}

// DeleteTag 删除标签
func (c *ContentClient) DeleteTag(ctx context.Context, id string) (*pb.DeleteTagResponse, error) {
	return c.client.DeleteTag(ctx, &pb.DeleteTagRequest{
		Id: id,
	})
}

// AuditContent 审核内容
func (c *ContentClient) AuditContent(ctx context.Context, contentId string, status string, reason string) (*pb.AuditContentResponse, error) {
	return c.client.AuditContent(ctx, &pb.AuditContentRequest{
		ContentId: contentId,
		Status:    status,
		Reason:    reason,
	})
}
