package types

// 通用成功响应
type SuccessResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 错误响应
type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 分页基础请求
type PageReq struct {
	Page     int32 `form:"page,default=1"`
	PageSize int32 `form:"pageSize,default=10"`
}
