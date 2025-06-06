package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// OperationRecordHandler 操作日志处理程序
type OperationRecordHandler struct {
	operationRecordService *service.OperationRecordApplicationService
}

// NewOperationRecordHandler 创建操作日志处理程序
func NewOperationRecordHandler(operationRecordService *service.OperationRecordApplicationService) *OperationRecordHandler {
	return &OperationRecordHandler{
		operationRecordService: operationRecordService,
	}
}

// GetOperationRecordList 获取操作日志列表
func (h *OperationRecordHandler) GetOperationRecordList(c *gin.Context) {
	var req types.GetOperationRecordListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务查询参数
	queryReq := h.operationRecordService.ConvertToQueryRequest(req)

	// 调用应用服务获取操作日志列表
	result, total, err := h.operationRecordService.GetOperationRecordList(c.Request.Context(), queryReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取操作日志列表失败: "+err.Error())
		return
	}

	// 转换为API响应
	recordList := make([]types.SysOperationRecord, 0, len(result))
	for _, record := range result {
		deletedAt := ""
		if record.DeletedAt != nil {
			deletedAt = record.DeletedAt.Format(time.RFC3339)
		}
		recordList = append(recordList, types.SysOperationRecord{
			Id:           record.ID,
			CreatedAt:    record.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    record.UpdatedAt.Format(time.RFC3339),
			DeletedAt:    deletedAt,
			Ip:           record.IP,
			Method:       record.Method,
			Path:         record.Path,
			Status:       record.Status,
			Latency:      record.Latency,
			Agent:        record.Agent,
			ErrorMessage: record.ErrorMessage,
			Body:         record.Body,
			UserId:       record.UserID,
			User:         record.User,
			Resp:         record.Resp,
		})
	}

	c.JSON(http.StatusOK, types.GetOperationRecordListResponse{
		Code:    http.StatusOK,
		Message: "获取操作日志列表成功",
		Data: types.OperationRecordPageData{
			List:     recordList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// DeleteOperationRecord 删除操作日志
func (h *OperationRecordHandler) DeleteOperationRecord(c *gin.Context) {
	var req types.DeleteOperationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务删除操作日志
	err := h.operationRecordService.DeleteOperationRecord(c.Request.Context(), req.Id)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "删除操作日志失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "删除操作日志成功",
	})
}

// DeleteOperationRecordsByIds 批量删除操作日志
func (h *OperationRecordHandler) DeleteOperationRecordsByIds(c *gin.Context) {
	var req types.DeleteOperationRecordsByIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务批量删除操作日志
	err := h.operationRecordService.DeleteOperationRecordsByIds(c.Request.Context(), req.Ids)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "批量删除操作日志失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "批量删除操作日志成功",
	})
}
