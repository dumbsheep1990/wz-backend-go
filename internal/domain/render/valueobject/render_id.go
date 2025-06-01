package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

// RenderID 表示渲染内容的唯一标识符
type RenderID struct {
	id string
}

// NewRenderID 创建一个新的渲染ID
func NewRenderID() RenderID {
	return RenderID{
		id: uuid.New().String(),
	}
}

// FromString 从字符串创建渲染ID
func FromString(id string) (RenderID, error) {
	if id == "" {
		return RenderID{}, errors.New("渲染ID不能为空")
	}
	
	return RenderID{id: id}, nil
}

// String 返回渲染ID的字符串表示
func (rid RenderID) String() string {
	return rid.id
}

// Equals 比较两个渲染ID是否相等
func (rid RenderID) Equals(other RenderID) bool {
	return rid.id == other.id
}
