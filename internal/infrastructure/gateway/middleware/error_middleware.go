package middleware

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 标准错误响应结构
type ErrorResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	RequestID  string      `json:"request_id,omitempty"`
	Error      string      `json:"error,omitempty"`
	Details    interface{} `json:"details,omitempty"`
	ServerTime int64       `json:"server_time,omitempty"`
}

// ErrorMiddleware 全局错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer来捕获处理过程中的panic
		defer func() {
			if r := recover(); r != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				log.Printf("服务发生严重错误: %v\n堆栈: %s", r, stack)
				
				// 获取请求ID
				requestID, _ := c.Get("request_id")
				
				// 返回统一的错误响应
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:      500,
					Message:   "服务器内部错误",
					RequestID: requestID.(string),
					Error:     "系统发生未预期的异常",
				})
				
				// 终止后续中间件的执行
				c.Abort()
			}
		}()
		
		// 处理请求
		c.Next()
		
		// 处理后检查是否有错误
		if len(c.Errors) > 0 {
			// 处理Gin框架收集的错误
			err := c.Errors.Last()
			
			// 获取请求ID
			requestID, exists := c.Get("request_id")
			if !exists {
				requestID = "unknown"
			}
			
			// 根据错误类型返回不同的状态码
			statusCode := http.StatusInternalServerError
			errorMessage := "服务器内部错误"
			errorDetail := err.Error()
			
			// 自定义错误类型处理
			var customErr *CustomError
			if errors.As(err.Err, &customErr) {
				statusCode = customErr.StatusCode
				errorMessage = customErr.Message
				errorDetail = customErr.Detail
			} else {
				// 其他常见错误类型的处理
				switch {
				case errors.Is(err.Err, ErrNotFound):
					statusCode = http.StatusNotFound
					errorMessage = "请求的资源不存在"
				case errors.Is(err.Err, ErrUnauthorized):
					statusCode = http.StatusUnauthorized
					errorMessage = "未授权的请求"
				case errors.Is(err.Err, ErrForbidden):
					statusCode = http.StatusForbidden
					errorMessage = "禁止访问"
				case errors.Is(err.Err, ErrBadRequest):
					statusCode = http.StatusBadRequest
					errorMessage = "无效的请求参数"
				case errors.Is(err.Err, ErrTooManyRequests):
					statusCode = http.StatusTooManyRequests
					errorMessage = "请求过于频繁"
				case errors.Is(err.Err, ErrServiceUnavailable):
					statusCode = http.StatusServiceUnavailable
					errorMessage = "服务暂时不可用"
				}
			}
			
			// 记录错误日志
			log.Printf("请求处理错误 [%v]: %s - %s", requestID, errorMessage, errorDetail)
			
			// 如果响应尚未写入，则返回错误响应
			if !c.Writer.Written() {
				c.JSON(statusCode, ErrorResponse{
					Code:      statusCode,
					Message:   errorMessage,
					RequestID: requestID.(string),
					Error:     errorDetail,
				})
			}
		}
	}
}

// 自定义错误类型
type CustomError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

// 创建新的自定义错误
func NewError(statusCode int, message, detail string) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
		Detail:     detail,
	}
}

// 预定义的错误类型
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrBadRequest        = errors.New("bad request")
	ErrTooManyRequests   = errors.New("too many requests")
	ErrServiceUnavailable = errors.New("service unavailable")
)
