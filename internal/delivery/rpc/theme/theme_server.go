package theme

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "wz-backend-go/api/rpc/theme"
)

// ThemeServer 主题gRPC服务器
type ThemeServer struct {
	pb.UnimplementedThemeServiceServer
	// TODO: 添加数据库连接和业务逻辑服务
}

// NewThemeServer 创建新的主题gRPC服务器
func NewThemeServer() *ThemeServer {
	return &ThemeServer{}
}

// GetThemeList 获取主题列表
func (s *ThemeServer) GetThemeList(ctx context.Context, req *pb.GetThemeListRequest) (*pb.GetThemeListResponse, error) {
	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// TODO: 调用业务逻辑层获取主题列表
	// 临时返回模拟数据
	themes := []*pb.Theme{
		{
			Id:               1,
			Name:             "默认主题",
			Description:      "WZ平台默认主题",
			PrimaryColor:     "#1890ff",
			TextColor:        "#333333",
			HeaderColor:      "#ffffff",
			LogoColor:        "#1890ff",
			MenuTextColor:    "#666666",
			ContentBgColor:   "#ffffff",
			SidebarColor:     "#f0f2f5",
			SidebarTextColor: "#666666",
			CardColor:        "#ffffff",
			LinkColor:        "#1890ff",
			IsDefault:        true,
			IsCurrent:        true,
			CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return &pb.GetThemeListResponse{
		Themes: themes,
		Total:  1,
	}, nil
}

// CreateTheme 创建主题
func (s *ThemeServer) CreateTheme(ctx context.Context, req *pb.CreateThemeRequest) (*pb.ThemeResponse, error) {
	// 参数验证
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "主题名称不能为空")
	}
	if req.PrimaryColor == "" {
		return nil, status.Errorf(codes.InvalidArgument, "主色调不能为空")
	}

	// TODO: 调用业务逻辑层创建主题
	theme := &pb.Theme{
		Id:               123, // 模拟生成的ID
		Name:             req.Name,
		Description:      req.Description,
		PrimaryColor:     req.PrimaryColor,
		TextColor:        req.TextColor,
		HeaderColor:      req.HeaderColor,
		LogoColor:        req.LogoColor,
		MenuTextColor:    req.MenuTextColor,
		ContentBgColor:   req.ContentBgColor,
		SidebarColor:     req.SidebarColor,
		SidebarTextColor: req.SidebarTextColor,
		CardColor:        req.CardColor,
		LinkColor:        req.LinkColor,
		IsDefault:        false,
		IsCurrent:        false,
		CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.ThemeResponse{
		Theme: theme,
	}, nil
}

// GetThemeDetail 获取主题详情
func (s *ThemeServer) GetThemeDetail(ctx context.Context, req *pb.GetThemeDetailRequest) (*pb.ThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}

	// TODO: 调用业务逻辑层获取主题详情
	theme := &pb.Theme{
		Id:               req.Id,
		Name:             "示例主题",
		Description:      "示例主题描述",
		PrimaryColor:     "#1890ff",
		TextColor:        "#333333",
		HeaderColor:      "#ffffff",
		LogoColor:        "#1890ff",
		MenuTextColor:    "#666666",
		ContentBgColor:   "#ffffff",
		SidebarColor:     "#f0f2f5",
		SidebarTextColor: "#666666",
		CardColor:        "#ffffff",
		LinkColor:        "#1890ff",
		IsDefault:        false,
		IsCurrent:        false,
		CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.ThemeResponse{
		Theme: theme,
	}, nil
}

// UpdateTheme 更新主题
func (s *ThemeServer) UpdateTheme(ctx context.Context, req *pb.UpdateThemeRequest) (*pb.ThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "主题名称不能为空")
	}

	// TODO: 调用业务逻辑层更新主题
	theme := &pb.Theme{
		Id:               req.Id,
		Name:             req.Name,
		Description:      req.Description,
		PrimaryColor:     req.PrimaryColor,
		TextColor:        req.TextColor,
		HeaderColor:      req.HeaderColor,
		LogoColor:        req.LogoColor,
		MenuTextColor:    req.MenuTextColor,
		ContentBgColor:   req.ContentBgColor,
		SidebarColor:     req.SidebarColor,
		SidebarTextColor: req.SidebarTextColor,
		CardColor:        req.CardColor,
		LinkColor:        req.LinkColor,
		CreatedAt:        time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.ThemeResponse{
		Theme: theme,
	}, nil
}

// DeleteTheme 删除主题
func (s *ThemeServer) DeleteTheme(ctx context.Context, req *pb.DeleteThemeRequest) (*pb.DeleteThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}

	// TODO: 调用业务逻辑层删除主题

	return &pb.DeleteThemeResponse{
		Success: true,
		Message: "删除主题成功",
	}, nil
}

// ApplyTheme 应用主题
func (s *ThemeServer) ApplyTheme(ctx context.Context, req *pb.ApplyThemeRequest) (*pb.ApplyThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}

	// TODO: 调用业务逻辑层应用主题

	return &pb.ApplyThemeResponse{
		Success: true,
		Message: "主题应用成功",
	}, nil
}

// GetCurrentTheme 获取当前主题
func (s *ThemeServer) GetCurrentTheme(ctx context.Context, req *pb.GetCurrentThemeRequest) (*pb.ThemeResponse, error) {
	// TODO: 调用业务逻辑层获取当前主题
	theme := &pb.Theme{
		Id:               1,
		Name:             "默认主题",
		Description:      "当前使用的主题",
		PrimaryColor:     "#1890ff",
		TextColor:        "#333333",
		HeaderColor:      "#ffffff",
		LogoColor:        "#1890ff",
		MenuTextColor:    "#666666",
		ContentBgColor:   "#ffffff",
		SidebarColor:     "#f0f2f5",
		SidebarTextColor: "#666666",
		CardColor:        "#ffffff",
		LinkColor:        "#1890ff",
		IsDefault:        true,
		IsCurrent:        true,
		CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.ThemeResponse{
		Theme: theme,
	}, nil
}

// PreviewTheme 预览主题
func (s *ThemeServer) PreviewTheme(ctx context.Context, req *pb.PreviewThemeRequest) (*pb.PreviewThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}

	// TODO: 调用业务逻辑层生成主题预览CSS
	css := `:root {
		--primary-color: #1890ff;
		--text-color: #333333;
		--header-color: #ffffff;
	}
	body { color: var(--text-color); }`

	return &pb.PreviewThemeResponse{
		Css: css,
	}, nil
}

// ExportTheme 导出主题
func (s *ThemeServer) ExportTheme(ctx context.Context, req *pb.ExportThemeRequest) (*pb.ExportThemeResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "主题ID无效")
	}

	// TODO: 调用业务逻辑层导出主题
	exportData := `{"theme": {"name": "示例主题", "colors": {}}}`

	return &pb.ExportThemeResponse{
		Data: exportData,
		Type: "json",
	}, nil
}

// ImportTheme 导入主题
func (s *ThemeServer) ImportTheme(ctx context.Context, req *pb.ImportThemeRequest) (*pb.ThemeResponse, error) {
	if req.Data == "" {
		return nil, status.Errorf(codes.InvalidArgument, "导入数据不能为空")
	}

	// TODO: 调用业务逻辑层导入主题
	theme := &pb.Theme{
		Id:          456, // 模拟生成的ID
		Name:        "导入的主题",
		Description: "从文件导入的主题",
		IsDefault:   false,
		IsCurrent:   false,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return &pb.ThemeResponse{
		Theme: theme,
	}, nil
} 