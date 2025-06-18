package link

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/client"
	"github.com/wz-backend-go/internal/types"
)

type LinkHandler struct {
	linkClient client.LinkClientInterface
}

func NewLinkHandler(linkClient client.LinkClientInterface) *LinkHandler {
	return &LinkHandler{
		linkClient: linkClient,
	}
}

// GetLinkList 获取链接列表
func (h *LinkHandler) GetLinkList(c *gin.Context) {
	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	// 过滤参数
	category := c.Query("category")
	status := c.Query("status")
	name := c.Query("name")
	
	req := &types.GetLinkListRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Status:   status,
		Name:     name,
	}
	
	resp, err := h.linkClient.GetLinkList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取链接列表失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取链接列表成功",
		Data: map[string]interface{}{
			"list":       resp.Links,
			"total":      resp.Total,
			"page":       req.Page,
			"pageSize":   req.PageSize,
		},
	})
}

// CreateLink 创建链接
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req types.CreateLinkRequest
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
			Message: "链接名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接URL不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Category == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接分类不能为空",
			Data:    nil,
		})
		return
	}
	
	link, err := h.linkClient.CreateLink(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "创建链接失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建链接成功",
		Data:    link,
	})
}

// UpdateLink 更新链接
func (h *LinkHandler) UpdateLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接ID格式错误",
			Data:    nil,
		})
		return
	}
	
	var req types.UpdateLinkRequest
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
			Message: "链接名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接URL不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Category == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接分类不能为空",
			Data:    nil,
		})
		return
	}
	
	link, err := h.linkClient.UpdateLink(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新链接失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新链接成功",
		Data:    link,
	})
}

// DeleteLink 删除链接
func (h *LinkHandler) DeleteLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接ID格式错误",
			Data:    nil,
		})
		return
	}
	
	err = h.linkClient.DeleteLink(context.Background(), &types.DeleteLinkRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "删除链接失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除链接成功",
		Data:    nil,
	})
}

// GetLinkDetail 获取链接详情
func (h *LinkHandler) GetLinkDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接ID格式错误",
			Data:    nil,
		})
		return
	}
	
	link, err := h.linkClient.GetLinkDetail(context.Background(), &types.GetLinkDetailRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取链接详情失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取链接详情成功",
		Data:    link,
	})
}

// VerifyLink 验证链接有效性
func (h *LinkHandler) VerifyLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接ID格式错误",
			Data:    nil,
		})
		return
	}
	
	result, err := h.linkClient.VerifyLink(context.Background(), &types.VerifyLinkRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "验证链接失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "验证链接成功",
		Data:    result,
	})
}

// BatchVerifyLinks 批量验证链接
func (h *LinkHandler) BatchVerifyLinks(c *gin.Context) {
	var req types.BatchVerifyLinksRequest
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
			Message: "请选择要验证的链接",
			Data:    nil,
		})
		return
	}
	
	result, err := h.linkClient.BatchVerifyLinks(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "批量验证链接失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "批量验证链接成功",
		Data:    result,
	})
}

// GetLinkCategories 获取链接分类列表
func (h *LinkHandler) GetLinkCategories(c *gin.Context) {
	categories, err := h.linkClient.GetLinkCategories(context.Background(), &types.GetLinkCategoriesRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取链接分类失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取链接分类成功",
		Data:    categories,
	})
}

// UpdateLinkSort 更新链接排序
func (h *LinkHandler) UpdateLinkSort(c *gin.Context) {
	var req types.UpdateLinkSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}
	
	if len(req.SortData) == 0 {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "排序数据不能为空",
			Data:    nil,
		})
		return
	}
	
	err := h.linkClient.UpdateLinkSort(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新链接排序失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新链接排序成功",
		Data:    nil,
	})
} 