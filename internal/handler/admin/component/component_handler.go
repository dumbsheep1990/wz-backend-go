package component

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/client"
	"github.com/wz-backend-go/internal/types"
)

type ComponentHandler struct {
	componentClient client.ComponentClientInterface
}

func NewComponentHandler(componentClient client.ComponentClientInterface) *ComponentHandler {
	return &ComponentHandler{
		componentClient: componentClient,
	}
}

// GetComponentList 获取组件列表
func (h *ComponentHandler) GetComponentList(c *gin.Context) {
	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "6")
	
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	// 过滤参数
	componentType := c.Query("type")
	name := c.Query("name")
	
	req := &types.GetComponentListRequest{
		Page:     page,
		PageSize: pageSize,
		Type:     componentType,
		Name:     name,
	}
	
	resp, err := h.componentClient.GetComponentList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取组件列表失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取组件列表成功",
		Data: map[string]interface{}{
			"list":       resp.Components,
			"total":      resp.Total,
			"page":       req.Page,
			"pageSize":   req.PageSize,
		},
	})
}

// CreateComponent 创建组件
func (h *ComponentHandler) CreateComponent(c *gin.Context) {
	var req types.CreateComponentRequest
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
			Message: "组件名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Type == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件类型不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件代码不能为空",
			Data:    nil,
		})
		return
	}
	
	component, err := h.componentClient.CreateComponentAdmin(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "创建组件失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建组件成功",
		Data:    component,
	})
}

// UpdateComponent 更新组件
func (h *ComponentHandler) UpdateComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件ID格式错误",
			Data:    nil,
		})
		return
	}
	
	var req types.UpdateComponentRequest
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
			Message: "组件名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Type == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件类型不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件代码不能为空",
			Data:    nil,
		})
		return
	}
	
	component, err := h.componentClient.UpdateComponentAdmin(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新组件失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新组件成功",
		Data:    component,
	})
}

// DeleteComponent 删除组件
func (h *ComponentHandler) DeleteComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件ID格式错误",
			Data:    nil,
		})
		return
	}
	
	err = h.componentClient.DeleteComponentAdmin(context.Background(), &types.DeleteComponentRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "删除组件失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除组件成功",
		Data:    nil,
	})
}

// GetComponentDetail 获取组件详情
func (h *ComponentHandler) GetComponentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件ID格式错误",
			Data:    nil,
		})
		return
	}
	
	component, err := h.componentClient.GetComponentDetail(context.Background(), &types.GetComponentDetailRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取组件详情失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取组件详情成功",
		Data:    component,
	})
}

// PreviewComponent 预览组件
func (h *ComponentHandler) PreviewComponent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "组件ID格式错误",
			Data:    nil,
		})
		return
	}
	
	previewData, err := h.componentClient.PreviewComponent(context.Background(), &types.PreviewComponentRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "预览组件失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "预览组件成功",
		Data:    previewData,
	})
}

// ImportComponent 导入组件
func (h *ComponentHandler) ImportComponent(c *gin.Context) {
	var req types.ImportComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}
	
	// 基础校验
	if req.Source == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导入源不能为空",
			Data:    nil,
		})
		return
	}
	
	component, err := h.componentClient.ImportComponent(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "导入组件失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "导入组件成功",
		Data:    component,
	})
}

// GetComponentTypes 获取组件类型列表
func (h *ComponentHandler) GetComponentTypes(c *gin.Context) {
	types, err := h.componentClient.GetComponentTypes(context.Background(), &types.GetComponentTypesRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取组件类型失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取组件类型成功",
		Data:    types,
	})
} 