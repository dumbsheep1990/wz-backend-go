package domain

import "errors"

// 定义领域层通用错误
var (
	ErrInvalidParam       = errors.New("无效的参数")
	ErrNotFound           = errors.New("资源不存在")
	ErrPermissionDenied   = errors.New("权限不足")
	ErrDuplicate          = errors.New("资源已存在")
	ErrInternal           = errors.New("内部错误")
	ErrUnauthenticated    = errors.New("未认证")
	ErrAlreadyExists      = errors.New("资源已存在")
	ErrInsufficientPoints = errors.New("积分不足")
	ErrAlreadyFavorite    = errors.New("已经收藏过该内容")
)
