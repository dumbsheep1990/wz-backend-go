package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "wz-backend-go/api/proto/site"
)

// SiteClient 是站点服务的客户端封装
type SiteClient struct {
	client pb.SiteServiceClient
	conn   *grpc.ClientConn
}

// NewSiteClient 创建一个新的站点服务客户端
func NewSiteClient(serviceAddr string) (*SiteClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50055" // 默认地址
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
		return nil, fmt.Errorf("无法连接到站点服务: %v", err)
	}

	client := pb.NewSiteServiceClient(conn)
	return &SiteClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *SiteClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetSiteStats 获取站点统计数据
func (c *SiteClient) GetSiteStats(ctx context.Context) (*pb.SiteStatsResponse, error) {
	return c.client.GetSiteStats(ctx, &pb.GetSiteStatsRequest{})
}

// ListSites 获取站点列表
func (c *SiteClient) ListSites(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListSitesResponse, error) {
	return c.client.ListSites(ctx, &pb.ListSitesRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetSite 获取站点详情
func (c *SiteClient) GetSite(ctx context.Context, id string) (*pb.Site, error) {
	return c.client.GetSite(ctx, &pb.GetSiteRequest{
		Id: id,
	})
}

// CreateSite 创建站点
func (c *SiteClient) CreateSite(ctx context.Context, req *pb.CreateSiteRequest) (*pb.Site, error) {
	return c.client.CreateSite(ctx, req)
}

// UpdateSite 更新站点
func (c *SiteClient) UpdateSite(ctx context.Context, req *pb.UpdateSiteRequest) (*pb.Site, error) {
	return c.client.UpdateSite(ctx, req)
}

// DeleteSite 删除站点
func (c *SiteClient) DeleteSite(ctx context.Context, id string) (*pb.DeleteSiteResponse, error) {
	return c.client.DeleteSite(ctx, &pb.DeleteSiteRequest{
		Id: id,
	})
}

// ListPages 获取页面列表
func (c *SiteClient) ListPages(ctx context.Context, siteId string, pageSize int32, pageNumber int32, filter string) (*pb.ListPagesResponse, error) {
	return c.client.ListPages(ctx, &pb.ListPagesRequest{
		SiteId:     siteId,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetPage 获取页面详情
func (c *SiteClient) GetPage(ctx context.Context, id string) (*pb.Page, error) {
	return c.client.GetPage(ctx, &pb.GetPageRequest{
		Id: id,
	})
}

// CreatePage 创建页面
func (c *SiteClient) CreatePage(ctx context.Context, req *pb.CreatePageRequest) (*pb.Page, error) {
	return c.client.CreatePage(ctx, req)
}

// UpdatePage 更新页面
func (c *SiteClient) UpdatePage(ctx context.Context, req *pb.UpdatePageRequest) (*pb.Page, error) {
	return c.client.UpdatePage(ctx, req)
}

// DeletePage 删除页面
func (c *SiteClient) DeletePage(ctx context.Context, id string) (*pb.DeletePageResponse, error) {
	return c.client.DeletePage(ctx, &pb.DeletePageRequest{
		Id: id,
	})
}

// PublishPage 发布页面
func (c *SiteClient) PublishPage(ctx context.Context, id string) (*pb.PublishPageResponse, error) {
	return c.client.PublishPage(ctx, &pb.PublishPageRequest{
		Id: id,
	})
}

// UnpublishPage 取消发布页面
func (c *SiteClient) UnpublishPage(ctx context.Context, id string) (*pb.UnpublishPageResponse, error) {
	return c.client.UnpublishPage(ctx, &pb.UnpublishPageRequest{
		Id: id,
	})
}

// ListTemplates 获取模板列表
func (c *SiteClient) ListTemplates(ctx context.Context, pageSize int32, pageNumber int32, filter string) (*pb.ListTemplatesResponse, error) {
	return c.client.ListTemplates(ctx, &pb.ListTemplatesRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Filter:     filter,
	})
}

// GetTemplate 获取模板详情
func (c *SiteClient) GetTemplate(ctx context.Context, id string) (*pb.Template, error) {
	return c.client.GetTemplate(ctx, &pb.GetTemplateRequest{
		Id: id,
	})
}

// CreateTemplate 创建模板
func (c *SiteClient) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.Template, error) {
	return c.client.CreateTemplate(ctx, req)
}

// UpdateTemplate 更新模板
func (c *SiteClient) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.Template, error) {
	return c.client.UpdateTemplate(ctx, req)
}

// DeleteTemplate 删除模板
func (c *SiteClient) DeleteTemplate(ctx context.Context, id string) (*pb.DeleteTemplateResponse, error) {
	return c.client.DeleteTemplate(ctx, &pb.DeleteTemplateRequest{
		Id: id,
	})
}

// GetSiteConfig 获取站点配置
func (c *SiteClient) GetSiteConfig(ctx context.Context, siteId string) (*pb.SiteConfig, error) {
	return c.client.GetSiteConfig(ctx, &pb.GetSiteConfigRequest{
		SiteId: siteId,
	})
}

// UpdateSiteConfig 更新站点配置
func (c *SiteClient) UpdateSiteConfig(ctx context.Context, req *pb.UpdateSiteConfigRequest) (*pb.SiteConfig, error) {
	return c.client.UpdateSiteConfig(ctx, req)
}
