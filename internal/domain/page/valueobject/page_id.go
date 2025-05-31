package valueobject

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PageID 页面ID值对象
type PageID struct {
	value string
}

// NewPageID 创建页面ID值对象
func NewPageID(id string) (PageID, error) {
	if id == "" {
		return PageID{}, errors.New("页面ID不能为空")
	}

	id = strings.TrimSpace(id)
	if !isValidPageID(id) {
		return PageID{}, errors.New("页面ID格式不正确")
	}

	return PageID{value: id}, nil
}

// GeneratePageID 生成新的页面ID
func GeneratePageID() PageID {
	id := fmt.Sprintf("page-%d", time.Now().UnixNano())
	return PageID{value: id}
}

// MustNewPageID 创建页面ID值对象，如果无效则panic
func MustNewPageID(id string) PageID {
	pid, err := NewPageID(id)
	if err != nil {
		panic("无效的页面ID: " + err.Error())
	}
	return pid
}

// Value 获取页面ID值
func (p PageID) Value() string {
	return p.value
}

// IsEmpty 检查是否为空
func (p PageID) IsEmpty() bool {
	return p.value == ""
}

// IsEquals 比较两个页面ID是否相等
func (p PageID) IsEquals(other PageID) bool {
	return p.value == other.value
}

// String 获取页面ID的字符串表示
func (p PageID) String() string {
	return p.value
}

// IsValid 检查页面ID是否有效
func (p PageID) IsValid() bool {
	return p.value != "" && isValidPageID(p.value)
}

// GetNumericPart 获取ID的数字部分（用于排序等）
func (p PageID) GetNumericPart() int64 {
	if strings.HasPrefix(p.value, "page-") {
		if numStr := strings.TrimPrefix(p.value, "page-"); numStr != "" {
			if num, err := strconv.ParseInt(numStr, 10, 64); err == nil {
				return num
			}
		}
	}
	return 0
}

// GetShortID 获取短ID（用于URL）
func (p PageID) GetShortID() string {
	if len(p.value) > 8 {
		return p.value[:8]
	}
	return p.value
}

// isValidPageID 验证页面ID格式
func isValidPageID(id string) bool {
	// 页面ID格式：page-{timestamp} 或 自定义格式
	if strings.HasPrefix(id, "page-") {
		suffix := strings.TrimPrefix(id, "page-")
		if len(suffix) > 0 {
			// 检查后缀是否为有效数字
			if _, err := strconv.ParseInt(suffix, 10, 64); err == nil {
				return true
			}
		}
	}
	
	// 允许自定义ID格式，但需要满足基本要求
	if len(id) >= 3 && len(id) <= 50 {
		// 只允许字母、数字、连字符和下划线
		for _, r := range id {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
				 (r >= '0' && r <= '9') || r == '-' || r == '_') {
				return false
			}
		}
		return true
	}
	
	return false
} 