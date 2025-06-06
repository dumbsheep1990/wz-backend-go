package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// DictionaryHandler 系统字典处理程序
type DictionaryHandler struct {
	dictionaryService *service.DictionaryApplicationService
}

// NewDictionaryHandler 创建系统字典处理程序
func NewDictionaryHandler(dictionaryService *service.DictionaryApplicationService) *DictionaryHandler {
	return &DictionaryHandler{
		dictionaryService: dictionaryService,
	}
}

// GetDictionaryList 获取系统字典列表
func (h *DictionaryHandler) GetDictionaryList(c *gin.Context) {
	var req types.GetDictionaryListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务查询参数
	queryReq := h.dictionaryService.ConvertToQueryRequest(req)

	// 调用应用服务获取字典列表
	result, total, err := h.dictionaryService.GetDictionaryList(c.Request.Context(), queryReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取字典列表失败: "+err.Error())
		return
	}

	// 转换为API响应
	dictList := make([]types.SysDictionary, 0, len(result))
	for _, dict := range result {
		deletedAt := ""
		if dict.DeletedAt != nil {
			deletedAt = dict.DeletedAt.Format(time.RFC3339)
		}
		dictList = append(dictList, types.SysDictionary{
			Id:          dict.ID,
			CreatedAt:   dict.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   dict.UpdatedAt.Format(time.RFC3339),
			DeletedAt:   deletedAt,
			Name:        dict.Name,
			Type:        dict.Type,
			Status:      dict.Status,
			Description: dict.Description,
		})
	}

	c.JSON(http.StatusOK, types.GetDictionaryListResponse{
		Code:    http.StatusOK,
		Message: "获取字典列表成功",
		Data: types.DictionaryPageData{
			List:     dictList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// CreateDictionary 创建系统字典
func (h *DictionaryHandler) CreateDictionary(c *gin.Context) {
	var req types.SysDictionary
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	createReq := h.dictionaryService.ConvertToCreateRequest(req)

	// 调用应用服务创建字典
	dict, err := h.dictionaryService.CreateDictionary(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "创建字典失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "创建字典成功",
		Data:    dict.ID,
	})
}

// UpdateDictionary 更新系统字典
func (h *DictionaryHandler) UpdateDictionary(c *gin.Context) {
	var req types.SysDictionary
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	updateReq := h.dictionaryService.ConvertToUpdateRequest(req)

	// 调用应用服务更新字典
	_, err := h.dictionaryService.UpdateDictionary(c.Request.Context(), updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新字典失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新字典成功",
	})
}

// DeleteDictionary 删除系统字典
func (h *DictionaryHandler) DeleteDictionary(c *gin.Context) {
	var req types.DeleteDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除字典
	err := h.dictionaryService.DeleteDictionary(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除字典失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除字典成功",
	})
}

// GetDictionaryDetailList 获取系统字典详情列表
func (h *DictionaryHandler) GetDictionaryDetailList(c *gin.Context) {
	var req types.GetDictionaryDetailListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务查询参数
	queryReq := h.dictionaryService.ConvertToDetailQueryRequest(req)

	// 调用应用服务获取字典详情列表
	result, total, err := h.dictionaryService.GetDictionaryDetailList(c.Request.Context(), queryReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取字典详情列表失败: "+err.Error())
		return
	}

	// 转换为API响应
	detailList := make([]types.SysDictionaryDetail, 0, len(result))
	for _, detail := range result {
		deletedAt := ""
		if detail.DeletedAt != nil {
			deletedAt = detail.DeletedAt.Format(time.RFC3339)
		}
		detailList = append(detailList, types.SysDictionaryDetail{
			Id:         detail.ID,
			CreatedAt:  detail.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  detail.UpdatedAt.Format(time.RFC3339),
			DeletedAt:  deletedAt,
			Label:      detail.Label,
			Value:      detail.Value,
			Status:     detail.Status,
			Sort:       detail.Sort,
			DictTypeId: detail.DictTypeID,
		})
	}

	c.JSON(http.StatusOK, types.GetDictionaryDetailListResponse{
		Code:    http.StatusOK,
		Message: "获取字典详情列表成功",
		Data: types.DictionaryDetailPageData{
			List:     detailList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// CreateDictionaryDetail 创建系统字典详情
func (h *DictionaryHandler) CreateDictionaryDetail(c *gin.Context) {
	var req types.SysDictionaryDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	createReq := h.dictionaryService.ConvertToDetailCreateRequest(req)

	// 调用应用服务创建字典详情
	detail, err := h.dictionaryService.CreateDictionaryDetail(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "创建字典详情失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "创建字典详情成功",
		Data:    detail.ID,
	})
}

// UpdateDictionaryDetail 更新系统字典详情
func (h *DictionaryHandler) UpdateDictionaryDetail(c *gin.Context) {
	var req types.SysDictionaryDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	updateReq := h.dictionaryService.ConvertToDetailUpdateRequest(req)

	// 调用应用服务更新字典详情
	_, err := h.dictionaryService.UpdateDictionaryDetail(c.Request.Context(), updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新字典详情失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新字典详情成功",
	})
}

// DeleteDictionaryDetail 删除系统字典详情
func (h *DictionaryHandler) DeleteDictionaryDetail(c *gin.Context) {
	var req types.DeleteDictionaryDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除字典详情
	err := h.dictionaryService.DeleteDictionaryDetail(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除字典详情失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除字典详情成功",
	})
}
