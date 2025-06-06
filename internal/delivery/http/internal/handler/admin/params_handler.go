package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// ParamsHandler 系统参数处理程序
type ParamsHandler struct {
	paramsService *service.ParamsApplicationService
}

// NewParamsHandler 创建系统参数处理程序
func NewParamsHandler(paramsService *service.ParamsApplicationService) *ParamsHandler {
	return &ParamsHandler{
		paramsService: paramsService,
	}
}

// GetParamsList 获取系统参数列表
func (h *ParamsHandler) GetParamsList(c *gin.Context) {
	var req types.GetParamsListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务查询参数
	queryReq := h.paramsService.ConvertToQueryRequest(req)

	// 调用应用服务获取参数列表
	result, total, err := h.paramsService.GetParamsList(c.Request.Context(), queryReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取系统参数列表失败: "+err.Error())
		return
	}

	// 转换为API响应
	paramsList := make([]types.SysParams, 0, len(result))
	for _, params := range result {
		deletedAt := ""
		if params.DeletedAt != nil {
			deletedAt = params.DeletedAt.Format(time.RFC3339)
		}
		paramsList = append(paramsList, types.SysParams{
			Id:         params.ID,
			CreatedAt:  params.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  params.UpdatedAt.Format(time.RFC3339),
			DeletedAt:  deletedAt,
			ParamName:  params.ParamName,
			ParamKey:   params.ParamKey,
			ParamValue: params.ParamValue,
			ParamType:  params.ParamType,
			ParamDesc:  params.ParamDesc,
			ParamGroup: params.ParamGroup,
			Status:     params.Status,
		})
	}

	c.JSON(http.StatusOK, types.GetParamsListResponse{
		Code:    http.StatusOK,
		Message: "获取系统参数列表成功",
		Data: types.ParamsPageData{
			List:     paramsList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// CreateParams 创建系统参数
func (h *ParamsHandler) CreateParams(c *gin.Context) {
	var req types.SysParams
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	createReq := h.paramsService.ConvertToCreateRequest(req)

	// 调用应用服务创建系统参数
	params, err := h.paramsService.CreateParams(c.Request.Context(), createReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "创建系统参数失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "创建系统参数成功",
		Data:    params.ID,
	})
}

// UpdateParams 更新系统参数
func (h *ParamsHandler) UpdateParams(c *gin.Context) {
	var req types.SysParams
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求
	updateReq := h.paramsService.ConvertToUpdateRequest(req)

	// 调用应用服务更新系统参数
	_, err := h.paramsService.UpdateParams(c.Request.Context(), updateReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "更新系统参数失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "更新系统参数成功",
	})
}

// DeleteParams 删除系统参数
func (h *ParamsHandler) DeleteParams(c *gin.Context) {
	var req types.DeleteParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除系统参数
	err := h.paramsService.DeleteParams(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除系统参数失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除系统参数成功",
	})
}
