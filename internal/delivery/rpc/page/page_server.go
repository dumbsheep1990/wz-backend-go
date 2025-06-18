package page

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "wz-backend-go/api/rpc/page"
)

// PageServer 页面gRPC服务器
type PageServer struct {
	pb.UnimplementedPageServiceServer
	// TODO: 添加数据库连接和业务逻辑服务
}

// NewPageServer 创建新的页面gRPC服务器
func NewPageServer() *PageServer {
	return &PageServer{}
}

// GetPageList 获取页面列表
func (s *PageServer) GetPageList(ctx context.Context, req *pb.GetPageListRequest) (*pb.GetPageListResponse, error) {
	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// TODO: 调用业务逻辑层获取页面列表
	// 临时返回模拟数据
	pages := []*pb.Page{
		{
			Id:          "1",
			Title:       "关于我们",
			Path:        "/about",
			Type:        "page",
			Content:     "公司介绍内容...",
			SeoTitle:    "关于我们 - WZ平台",
			SeoDesc:     "了解WZ平台的发展历程",
			SeoKeywords: "WZ,关于,公司",
			Status:      "published",
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return &pb.GetPageListResponse{
		Pages: pages,
		Total: 1,
	}, nil
}

// CreatePage 创建页面
func (s *PageServer) CreatePage(ctx context.Context, req *pb.CreatePageRequest) (*pb.PageResponse, error) {
	// 参数验证
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面标题不能为空")
	}
	if req.Path == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面路径不能为空")
	}

	// TODO: 调用业务逻辑层创建页面
	page := &pb.Page{
		Id:          "generated-id",
		Title:       req.Title,
		Path:        req.Path,
		Type:        req.Type,
		Content:     req.Content,
		SeoTitle:    req.SeoTitle,
		SeoDesc:     req.SeoDesc,
		SeoKeywords: req.SeoKeywords,
		Status:      req.Status,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.PageResponse{
		Page: page,
	}, nil
}

// GetPageDetail 获取页面详情
func (s *PageServer) GetPageDetail(ctx context.Context, req *pb.GetPageDetailRequest) (*pb.PageResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面ID不能为空")
	}

	// TODO: 调用业务逻辑层获取页面详情
	page := &pb.Page{
		Id:          req.Id,
		Title:       "示例页面",
		Path:        "/example",
		Type:        "page",
		Content:     "页面内容...",
		SeoTitle:    "示例页面",
		SeoDesc:     "示例页面描述",
		SeoKeywords: "示例,页面",
		Status:      "published",
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.PageResponse{
		Page: page,
	}, nil
}

// UpdatePage 更新页面
func (s *PageServer) UpdatePage(ctx context.Context, req *pb.UpdatePageRequest) (*pb.PageResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面ID不能为空")
	}
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面标题不能为空")
	}

	// TODO: 调用业务逻辑层更新页面
	page := &pb.Page{
		Id:          req.Id,
		Title:       req.Title,
		Path:        req.Path,
		Type:        req.Type,
		Content:     req.Content,
		SeoTitle:    req.SeoTitle,
		SeoDesc:     req.SeoDesc,
		SeoKeywords: req.SeoKeywords,
		Status:      req.Status,
		CreatedAt:   time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.PageResponse{
		Page: page,
	}, nil
}

// DeletePage 删除页面
func (s *PageServer) DeletePage(ctx context.Context, req *pb.DeletePageRequest) (*pb.DeletePageResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面ID不能为空")
	}

	// TODO: 调用业务逻辑层删除页面

	return &pb.DeletePageResponse{
		Success: true,
		Message: "删除页面成功",
	}, nil
}

// TogglePageStatus 切换页面状态
func (s *PageServer) TogglePageStatus(ctx context.Context, req *pb.TogglePageStatusRequest) (*pb.PageResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面ID不能为空")
	}

	// TODO: 调用业务逻辑层切换页面状态
	page := &pb.Page{
		Id:     req.Id,
		Title:  "示例页面",
		Status: req.Status,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.PageResponse{
		Page: page,
	}, nil
}

// PreviewPage 预览页面
func (s *PageServer) PreviewPage(ctx context.Context, req *pb.PreviewPageRequest) (*pb.PreviewPageResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "页面ID不能为空")
	}

	// TODO: 调用业务逻辑层生成页面预览
	return &pb.PreviewPageResponse{
		Html: "<html><body><h1>页面预览</h1></body></html>",
		Css:  "body { font-family: Arial; }",
		Js:   "console.log('Page preview loaded');",
	}, nil
}

// BatchUpdatePage 批量更新页面
func (s *PageServer) BatchUpdatePage(ctx context.Context, req *pb.BatchUpdatePageRequest) (*pb.BatchUpdatePageResponse, error) {
	if len(req.Ids) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "请选择要更新的页面")
	}

	// TODO: 调用业务逻辑层批量更新页面

	return &pb.BatchUpdatePageResponse{
		UpdatedCount: int32(len(req.Ids)),
		Message:      "批量更新页面成功",
	}, nil
} 