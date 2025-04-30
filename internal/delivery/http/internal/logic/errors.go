package logic

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 错误定义
var (
	// 通用错误
	errorInternalServer       = NewCodeError(http.StatusInternalServerError, "服务器内部错误")
	errorBadRequest           = NewCodeError(http.StatusBadRequest, "请求参数错误")
	errorUnauthorized         = NewCodeError(http.StatusUnauthorized, "未授权访问")
	errorForbidden            = NewCodeError(http.StatusForbidden, "没有操作权限")
	errorNotFound             = NewCodeError(http.StatusNotFound, "资源不存在")
	
	// 服务注册相关错误
	errorServiceNameRequired  = NewCodeError(http.StatusBadRequest, "服务名称不能为空")
	errorInstanceIDRequired   = NewCodeError(http.StatusBadRequest, "实例ID不能为空")
	errorServiceNotFound      = NewCodeError(http.StatusNotFound, "服务不存在")
	errorInstanceNotFound     = NewCodeError(http.StatusNotFound, "服务实例不存在")
	errorInvalidStatus        = NewCodeError(http.StatusBadRequest, "无效的状态值")
	errorRegisterFailed       = NewCodeError(http.StatusInternalServerError, "服务注册失败")
	errorDeregisterFailed     = NewCodeError(http.StatusInternalServerError, "服务注销失败")
	errorHealthCheckFailed    = NewCodeError(http.StatusInternalServerError, "健康检查失败")
)

// CodeError 自定义错误码结构
type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (e *CodeError) Error() string {
	return e.Message
}

// NewCodeError 创建新的自定义错误
func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{
		Code:    code,
		Message: msg,
	}
}

// FromError 从普通错误转换为自定义错误
func FromError(err error) *CodeError {
	if err == nil {
		return nil
	}
	
	var codeError *CodeError
	if errors.As(err, &codeError) {
		return codeError
	}
	
	return errorInternalServer
}
