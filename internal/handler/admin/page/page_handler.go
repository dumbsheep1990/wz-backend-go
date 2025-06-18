package page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/client"
	"github.com/wz-backend-go/internal/types"
)

type PageHandler struct {
	pageClient client.PageClientInterface
}

func NewPageHandler(pageClient client.PageClientInterface) *PageHandler {
	return &PageHandler{
		pageClient: pageClient,
	}
}

// GetPageList 获取页面列表
func (h *PageHandler) GetPageList(c *gin.Context) {
	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	// 过滤参数
	pageType := c.Query("type")
	status := c.Query("status")
	title := c.Query("title")
	
	req := &types.GetPageListRequest{
		Page:     page,
		PageSize: pageSize,
		Type:     pageType,
		Status:   status,
		Title:    title,
	}
	
	resp, err := h.pageClient.GetPageList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取页面列表失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取页面列表成功",
		Data: map[string]interface{}{
			"list":       resp.Pages,
			"total":      resp.Total,
			"page":       req.Page,
			"pageSize":   req.PageSize,
		},
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
	
	if req.Type == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面类型不能为空",
			Data:    nil,
		})
		return
	}
	
	page, err := h.pageClient.CreatePage(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "创建页面失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建页面成功",
		Data:    page,
	})
}

// UpdatePage 更新页面
func (h *PageHandler) UpdatePage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID格式错误",
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
	
	req.ID = id
	
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
	
	if req.Type == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面类型不能为空",
			Data:    nil,
		})
		return
	}
	
	page, err := h.pageClient.UpdatePage(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新页面失败",
			Data:    nil,
		})
		return
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID格式错误",
			Data:    nil,
		})
		return
	}
	
	err = h.pageClient.DeletePage(context.Background(), &types.DeletePageRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "删除页面失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除页面成功",
		Data:    nil,
	})
}

// GetPageDetail 获取页面详情
func (h *PageHandler) GetPageDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID格式错误",
			Data:    nil,
		})
		return
	}
	
	page, err := h.pageClient.GetPageDetail(context.Background(), &types.GetPageDetailRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取页面详情失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取页面详情成功",
		Data:    page,
	})
}

// TogglePageStatus 切换页面状态
func (h *PageHandler) TogglePageStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID格式错误",
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
	
	req.ID = id
	
	page, err := h.pageClient.TogglePageStatus(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "切换页面状态失败",
			Data:    nil,
		})
		return
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "页面ID格式错误",
			Data:    nil,
		})
		return
	}
	
	previewData, err := h.pageClient.PreviewPage(context.Background(), &types.PreviewPageRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "预览页面失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "预览页面成功",
		Data:    previewData,
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
			Message: "请选择要批量更新的页面",
			Data:    nil,
		})
		return
	}
	
	err := h.pageClient.BatchUpdatePage(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "批量更新页面失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "批量更新页面成功",
		Data:    nil,
	})
} 