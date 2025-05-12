package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码，0表示成功
	Message string      `json:"message"` // 状态信息
	Data    interface{} `json:"data"`    // 数据
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 返回请求参数错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 返回未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 返回禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
		Data:    nil,
	})
}

// NotFound 返回资源不存在响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
		Data:    nil,
	})
}

// ServerError 返回服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
		Data:    nil,
	})
}

// CustomError 返回自定义错误响应
func CustomError(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
