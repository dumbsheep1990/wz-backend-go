package navigation

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wz-backend-go/internal/client"
	"github.com/wz-backend-go/internal/types"
)

type TreeNavigationHandler struct {
	navigationClient client.TreeNavigationClientInterface
}

func NewTreeNavigationHandler(navigationClient client.TreeNavigationClientInterface) *TreeNavigationHandler {
	return &TreeNavigationHandler{
		navigationClient: navigationClient,
	}
}

// GetNavigationTree 获取导航树
func (h *TreeNavigationHandler) GetNavigationTree(c *gin.Context) {
	navType := c.DefaultQuery("type", "main") // main, footer, sidebar
	
	resp, err := h.navigationClient.GetNavigationTree(context.Background(), &types.GetNavigationTreeRequest{
		Type: navType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取导航树失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取导航树成功",
		Data:    resp.NavigationTree,
	})
}

// CreateNavigationItem 创建导航项
func (h *TreeNavigationHandler) CreateNavigationItem(c *gin.Context) {
	var req types.CreateNavigationItemRequest
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
			Message: "导航名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接地址不能为空",
			Data:    nil,
		})
		return
	}
	
	navItem, err := h.navigationClient.CreateNavigationItem(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "创建导航项失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "创建导航项成功",
		Data:    navItem,
	})
}

// UpdateNavigationItem 更新导航项
func (h *TreeNavigationHandler) UpdateNavigationItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导航项ID格式错误",
			Data:    nil,
		})
		return
	}
	
	var req types.UpdateNavigationItemRequest
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
			Message: "导航名称不能为空",
			Data:    nil,
		})
		return
	}
	
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "链接地址不能为空",
			Data:    nil,
		})
		return
	}
	
	navItem, err := h.navigationClient.UpdateNavigationItem(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新导航项失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新导航项成功",
		Data:    navItem,
	})
}

// DeleteNavigationItem 删除导航项
func (h *TreeNavigationHandler) DeleteNavigationItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导航项ID格式错误",
			Data:    nil,
		})
		return
	}
	
	err = h.navigationClient.DeleteNavigationItem(context.Background(), &types.DeleteNavigationItemRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "删除导航项失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "删除导航项成功",
		Data:    nil,
	})
}

// UpdateNavigationOrder 更新导航排序
func (h *TreeNavigationHandler) UpdateNavigationOrder(c *gin.Context) {
	var req types.UpdateNavigationOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}
	
	if len(req.OrderData) == 0 {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "排序数据不能为空",
			Data:    nil,
		})
		return
	}
	
	err := h.navigationClient.UpdateNavigationOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "更新导航排序失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "更新导航排序成功",
		Data:    nil,
	})
}

// GetNavigationItem 获取导航项详情
func (h *TreeNavigationHandler) GetNavigationItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导航项ID格式错误",
			Data:    nil,
		})
		return
	}
	
	navItem, err := h.navigationClient.GetNavigationItem(context.Background(), &types.GetNavigationItemRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "获取导航项详情失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "获取导航项详情成功",
		Data:    navItem,
	})
}

// ToggleNavigationVisibility 切换导航项显示/隐藏状态
func (h *TreeNavigationHandler) ToggleNavigationVisibility(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "导航项ID格式错误",
			Data:    nil,
		})
		return
	}
	
	var req types.ToggleNavigationVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "请求参数格式错误",
			Data:    nil,
		})
		return
	}
	
	req.ID = id
	
	navItem, err := h.navigationClient.ToggleNavigationVisibility(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "切换导航项状态失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "切换导航项状态成功",
		Data:    navItem,
	})
}

// BatchDeleteNavigationItems 批量删除导航项
func (h *TreeNavigationHandler) BatchDeleteNavigationItems(c *gin.Context) {
	var req types.BatchDeleteNavigationItemsRequest
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
			Message: "请选择要删除的导航项",
			Data:    nil,
		})
		return
	}
	
	err := h.navigationClient.BatchDeleteNavigationItems(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "批量删除导航项失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "批量删除导航项成功",
		Data:    nil,
	})
}

// ExportNavigationTree 导出导航树
func (h *TreeNavigationHandler) ExportNavigationTree(c *gin.Context) {
	navType := c.DefaultQuery("type", "main")
	
	exportData, err := h.navigationClient.ExportNavigationTree(context.Background(), &types.ExportNavigationTreeRequest{
		Type: navType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "导出导航树失败",
			Data:    nil,
		})
		return
	}
	
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=navigation_tree.json")
	c.Data(http.StatusOK, "application/json", []byte(exportData.Data))
}

// ImportNavigationTree 导入导航树
func (h *TreeNavigationHandler) ImportNavigationTree(c *gin.Context) {
	var req types.ImportNavigationTreeRequest
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
	
	result, err := h.navigationClient.ImportNavigationTree(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Code:    500,
			Message: "导入导航树失败",
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, types.Response{
		Code:    200,
		Message: "导入导航树成功",
		Data:    result,
	})
} 