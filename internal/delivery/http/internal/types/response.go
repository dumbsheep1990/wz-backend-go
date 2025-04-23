package types

// Response 通用API响应格式
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Status:  "success",
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}
