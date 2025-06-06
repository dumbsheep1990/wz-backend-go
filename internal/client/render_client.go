package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/render"
)

// RenderClient 是渲染服务的客户端封装
type RenderClient struct {
	client pb.RenderServiceClient
	conn   *grpc.ClientConn
}

// NewRenderClient 创建一个新的渲染服务客户端
func NewRenderClient(serviceAddr string) (*RenderClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50057" // 默认地址
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
		return nil, fmt.Errorf("无法连接到渲染服务: %v", err)
	}

	client := pb.NewRenderServiceClient(conn)
	return &RenderClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *RenderClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// RenderPage 渲染页面
func (c *RenderClient) RenderPage(ctx context.Context, pageId string, params map[string]string) (*pb.RenderPageResponse, error) {
	return c.client.RenderPage(ctx, &pb.RenderPageRequest{
		PageId: pageId,
		Params: params,
	})
}

// RenderTemplate 渲染模板
func (c *RenderClient) RenderTemplate(ctx context.Context, templateId string, data map[string]string) (*pb.RenderTemplateResponse, error) {
	return c.client.RenderTemplate(ctx, &pb.RenderTemplateRequest{
		TemplateId: templateId,
		Data:       data,
	})
}

// RenderComponent 渲染组件
func (c *RenderClient) RenderComponent(ctx context.Context, componentId string, version string, props map[string]string) (*pb.RenderComponentResponse, error) {
	return c.client.RenderComponent(ctx, &pb.RenderComponentRequest{
		ComponentId: componentId,
		Version:     version,
		Props:       props,
	})
}

// GetRenderStats 获取渲染统计数据
func (c *RenderClient) GetRenderStats(ctx context.Context) (*pb.RenderStatsResponse, error) {
	return c.client.GetRenderStats(ctx, &pb.GetRenderStatsRequest{})
}

// PurgeCache 清除缓存
func (c *RenderClient) PurgeCache(ctx context.Context, target string, id string) (*pb.PurgeCacheResponse, error) {
	return c.client.PurgeCache(ctx, &pb.PurgeCacheRequest{
		Target: target,
		Id:     id,
	})
}

// GetRenderHistory 获取渲染历史
func (c *RenderClient) GetRenderHistory(ctx context.Context, target string, id string, pageSize int32, pageNumber int32) (*pb.GetRenderHistoryResponse, error) {
	return c.client.GetRenderHistory(ctx, &pb.GetRenderHistoryRequest{
		Target:     target,
		Id:         id,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
}
