package theme

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/client"
	"github.com/wz-backend-go/internal/types"
)

type ThemeHandler struct {
	themeClient client.ThemeClientInterface
}

func NewThemeHandler(themeClient client.ThemeClientInterface) *ThemeHandler {
	return &ThemeHandler{
		themeClient: themeClient,
	}
}

// GetThemeList 获取主题列表
func (h *ThemeHandler) GetThemeList(c *gin.Context) {
	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	// 过滤参数
	name := c.Query("name")
	
	req := &types.GetThemeListRequest{
		Page:     page,
		PageSize: pageSize,
		Name:     name,
	}
	
	resp, err := h.themeClient.GetThemeList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取主题列表失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取主题列表成功",
		Data: map[string]interface{}{
			"list":       resp.Themes,
			"total":      resp.Total,
			"page":       req.Page,
			"pageSize":   req.PageSize,
		},
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
	
	theme, err := h.themeClient.CreateTheme(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "创建主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建主题成功",
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
	
	req.ID = id
	
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
	
	theme, err := h.themeClient.UpdateTheme(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新主题失败",
			Data:    nil,
		})
		return
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
	
	err = h.themeClient.DeleteTheme(context.Background(), &types.DeleteThemeRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "删除主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除主题成功",
		Data:    nil,
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
	
	theme, err := h.themeClient.GetThemeDetail(context.Background(), &types.GetThemeDetailRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取主题详情失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取主题详情成功",
		Data:    theme,
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
	
	result, err := h.themeClient.ApplyTheme(context.Background(), &types.ApplyThemeRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "应用主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "应用主题成功",
		Data:    result,
	})
}

// GetCurrentTheme 获取当前主题
func (h *ThemeHandler) GetCurrentTheme(c *gin.Context) {
	theme, err := h.themeClient.GetCurrentTheme(context.Background(), &types.GetCurrentThemeRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取当前主题失败",
			Data:    nil,
		})
		return
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
	
	previewData, err := h.themeClient.PreviewTheme(context.Background(), &types.PreviewThemeRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "预览主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "预览主题成功",
		Data:    previewData,
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
	
	exportData, err := h.themeClient.ExportTheme(context.Background(), &types.ExportThemeRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "导出主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "导出主题成功",
		Data:    exportData,
	})
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
	
	// 基础校验
	if req.Data == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导入数据不能为空",
			Data:    nil,
		})
		return
	}
	
	theme, err := h.themeClient.ImportTheme(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "导入主题失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "导入主题成功",
		Data:    theme,
	})
} 