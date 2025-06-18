package theme

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/types"
)

type ThemeHandler struct {
	// TODO: 连接到Theme微服务的gRPC客户端
	// themeService themeService.ThemeServiceClient
}

func NewThemeHandler() *ThemeHandler {
	return &ThemeHandler{}
}

// GetThemeList 获取主题列表
func (h *ThemeHandler) GetThemeList(c *gin.Context) {
	var req types.GetThemeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口
	// 临时返回模拟数据
	response := &types.GetThemeListResponse{
		Themes: []*types.Theme{
			{
				ID:               1,
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
			},
		},
		Total: 1,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取主题列表成功",
		Data:    response,
	})
}

// CreateTheme 创建主题
func (h *ThemeHandler) CreateTheme(c *gin.Context) {
	var req types.CreateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	// 基础校验
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题名称不能为空",
			Data:    nil,
		})
		return
	}

	if req.PrimaryColor == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主色调不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口创建主题
	theme := &types.Theme{
		ID:               123, // 模拟生成的ID
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
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建主题成功",
		Data:    theme,
	})
}

// GetThemeDetail 获取主题详情
func (h *ThemeHandler) GetThemeDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口
	theme := &types.Theme{
		ID:               id,
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
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取主题详情成功",
		Data:    theme,
	})
}

// UpdateTheme 更新主题
func (h *ThemeHandler) UpdateTheme(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	var req types.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	// 基础校验
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题名称不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口更新主题
	theme := &types.Theme{
		ID:               id,
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
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新主题成功",
		Data:    theme,
	})
}

// DeleteTheme 删除主题
func (h *ThemeHandler) DeleteTheme(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口删除主题

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除主题成功",
		Data:    nil,
	})
}

// ApplyTheme 应用主题
func (h *ThemeHandler) ApplyTheme(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口应用主题
	result := &types.ApplyThemeResult{
		Success: true,
		Message: "主题应用成功",
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "应用主题成功",
		Data:    result,
	})
}

// GetCurrentTheme 获取当前主题
func (h *ThemeHandler) GetCurrentTheme(c *gin.Context) {
	// TODO: 调用Theme微服务gRPC接口获取当前主题
	theme := &types.Theme{
		ID:               1,
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
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取当前主题成功",
		Data:    theme,
	})
}

// PreviewTheme 预览主题
func (h *ThemeHandler) PreviewTheme(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口生成预览CSS
	preview := &types.ThemePreview{
		CSS: `:root {
			--primary-color: #1890ff;
			--text-color: #333333;
			--header-color: #ffffff;
		}
		body { color: var(--text-color); }`,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "生成主题预览成功",
		Data:    preview,
	})
}

// ExportTheme 导出主题
func (h *ThemeHandler) ExportTheme(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "主题ID格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口导出主题
	exportData := &types.ThemeExport{
		Data: `{"theme": {"name": "示例主题", "colors": {}}}`,
		Type: "json",
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=theme.json")
	c.Data(http.StatusOK, "application/json", []byte(exportData.Data))
}

// ImportTheme 导入主题
func (h *ThemeHandler) ImportTheme(c *gin.Context) {
	var req types.ImportThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	if req.Data == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导入数据不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Theme微服务gRPC接口导入主题
	theme := &types.Theme{
		ID:          456, // 模拟生成的ID
		Name:        "导入的主题",
		Description: "从文件导入的主题",
		IsDefault:   false,
		IsCurrent:   false,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "导入主题成功",
		Data:    theme,
	})
} 