package valueobject

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CommunityID 社区ID值对象
type CommunityID struct {
	value string
}

// NewCommunityID 创建社区ID值对象
func NewCommunityID(id string) (CommunityID, error) {
	if id == "" {
		return CommunityID{}, errors.New("社区ID不能为空")
	}

	id = strings.TrimSpace(id)
	if !isValidCommunityID(id) {
		return CommunityID{}, errors.New("社区ID格式不正确")
	}

	return CommunityID{value: id}, nil
}

// GenerateCommunityID 生成新的社区ID
func GenerateCommunityID() CommunityID {
	id := fmt.Sprintf("comm-%d", time.Now().UnixNano())
	return CommunityID{value: id}
}

// MustNewCommunityID 创建社区ID值对象，如果无效则panic
func MustNewCommunityID(id string) CommunityID {
	cid, err := NewCommunityID(id)
	if err != nil {
		panic("无效的社区ID: " + err.Error())
	}
	return cid
}

// Value 获取社区ID值
func (c CommunityID) Value() string {
	return c.value
}

// IsEmpty 检查是否为空
func (c CommunityID) IsEmpty() bool {
	return c.value == ""
}

// IsEquals 比较两个社区ID是否相等
func (c CommunityID) IsEquals(other CommunityID) bool {
	return c.value == other.value
}

// String 获取社区ID的字符串表示
func (c CommunityID) String() string {
	return c.value
}

// IsValid 检查社区ID是否有效
func (c CommunityID) IsValid() bool {
	return c.value != "" && isValidCommunityID(c.value)
}

// GetNumericPart 获取ID的数字部分（用于排序等）
func (c CommunityID) GetNumericPart() int64 {
	if strings.HasPrefix(c.value, "comm-") {
		if numStr := strings.TrimPrefix(c.value, "comm-"); numStr != "" {
			if num, err := strconv.ParseInt(numStr, 10, 64); err == nil {
				return num
			}
		}
	}
	return 0
}

// isValidCommunityID 验证社区ID格式
func isValidCommunityID(id string) bool {
	// 社区ID格式：comm-{timestamp} 或 自定义格式
	if strings.HasPrefix(id, "comm-") {
		suffix := strings.TrimPrefix(id, "comm-")
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