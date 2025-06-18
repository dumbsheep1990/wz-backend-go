package page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/types"
)

type PageHandler struct {
	// TODO: 连接到Page微服务的gRPC客户端
	// pageService pageService.PageServiceClient
}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

// GetPageList 获取页面列表
func (h *PageHandler) GetPageList(c *gin.Context) {
	var req types.GetPageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口
	// 临时返回模拟数据
	response := &types.GetPageListResponse{
		Pages: []*types.Page{
			{
				ID:          "1",
				Title:       "关于我们",
				Path:        "/about",
				Type:        "page",
				Content:     "公司介绍内容...",
				SeoTitle:    "关于我们 - WZ平台",
				SeoDesc:     "了解WZ平台的发展历程",
				SeoKeywords: "WZ,关于,公司",
				Status:      "published",
			},
		},
		Total: 1,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取页面列表成功",
		Data:    response,
	})
}

// CreatePage 创建页面
func (h *PageHandler) CreatePage(c *gin.Context) {
	var req types.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	// 基础校验
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面标题不能为空",
			Data:    nil,
		})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面路径不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口创建页面
	page := &types.Page{
		ID:          "generated-id",
		Title:       req.Title,
		Path:        req.Path,
		Type:        req.Type,
		Content:     req.Content,
		SeoTitle:    req.SeoTitle,
		SeoDesc:     req.SeoDesc,
		SeoKeywords: req.SeoKeywords,
		Status:      req.Status,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建页面成功",
		Data:    page,
	})
}

// GetPageDetail 获取页面详情
func (h *PageHandler) GetPageDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口
	page := &types.Page{
		ID:          idStr,
		Title:       "示例页面",
		Path:        "/example",
		Type:        "page",
		Content:     "页面内容...",
		SeoTitle:    "示例页面",
		SeoDesc:     "示例页面描述",
		SeoKeywords: "示例,页面",
		Status:      "published",
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取页面详情成功",
		Data:    page,
	})
}

// UpdatePage 更新页面
func (h *PageHandler) UpdatePage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID不能为空",
			Data:    nil,
		})
		return
	}

	var req types.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	// 基础校验
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面标题不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口更新页面
	page := &types.Page{
		ID:          idStr,
		Title:       req.Title,
		Path:        req.Path,
		Type:        req.Type,
		Content:     req.Content,
		SeoTitle:    req.SeoTitle,
		SeoDesc:     req.SeoDesc,
		SeoKeywords: req.SeoKeywords,
		Status:      req.Status,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新页面成功",
		Data:    page,
	})
}

// DeletePage 删除页面
func (h *PageHandler) DeletePage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口删除页面

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除页面成功",
		Data:    nil,
	})
}

// TogglePageStatus 切换页面状态
func (h *PageHandler) TogglePageStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID不能为空",
			Data:    nil,
		})
		return
	}

	var req types.TogglePageStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口切换状态
	page := &types.Page{
		ID:     idStr,
		Title:  "示例页面",
		Status: req.Status,
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "切换页面状态成功",
		Data:    page,
	})
}

// PreviewPage 预览页面
func (h *PageHandler) PreviewPage(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口生成预览
	preview := &types.PagePreview{
		HTML: "<html><body><h1>页面预览</h1></body></html>",
		CSS:  "body { font-family: Arial; }",
		JS:   "console.log('Page preview loaded');",
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "生成页面预览成功",
		Data:    preview,
	})
}

// BatchUpdate 批量更新页面
func (h *PageHandler) BatchUpdate(c *gin.Context) {
	var req types.BatchUpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请选择要更新的页面",
			Data:    nil,
		})
		return
	}

	// TODO: 调用Page微服务gRPC接口批量更新

	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "批量更新页面成功",
		Data:    map[string]interface{}{
			"updated_count": len(req.IDs),
		},
	})
} 