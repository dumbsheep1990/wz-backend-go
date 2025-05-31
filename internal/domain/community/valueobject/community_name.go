package valueobject

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// CommunityName 社区名称值对象
type CommunityName struct {
	value string
}

// NewCommunityName 创建社区名称值对象
func NewCommunityName(name string) (CommunityName, error) {
	if name == "" {
		return CommunityName{}, errors.New("社区名称不能为空")
	}

	name = strings.TrimSpace(name)
	if err := validateCommunityName(name); err != nil {
		return CommunityName{}, err
	}

	return CommunityName{value: name}, nil
}

// MustNewCommunityName 创建社区名称值对象，如果无效则panic
func MustNewCommunityName(name string) CommunityName {
	cn, err := NewCommunityName(name)
	if err != nil {
		panic("无效的社区名称: " + err.Error())
	}
	return cn
}

// Value 获取社区名称值
func (c CommunityName) Value() string {
	return c.value
}

// IsEmpty 检查是否为空
func (c CommunityName) IsEmpty() bool {
	return c.value == ""
}

// IsEquals 比较两个社区名称是否相等
func (c CommunityName) IsEquals(other CommunityName) bool {
	return strings.ToLower(c.value) == strings.ToLower(other.value)
}

// String 获取社区名称的字符串表示
func (c CommunityName) String() string {
	return c.value
}

// IsValid 检查社区名称是否有效
func (c CommunityName) IsValid() bool {
	return c.value != "" && validateCommunityName(c.value) == nil
}

// Length 获取社区名称长度（考虑中文字符）
func (c CommunityName) Length() int {
	return utf8.RuneCountInString(c.value)
}

// ToLower 转换为小写
func (c CommunityName) ToLower() string {
	return strings.ToLower(c.value)
}

// ContainsKeyword 检查是否包含特定关键词
func (c CommunityName) ContainsKeyword(keyword string) bool {
	return strings.Contains(strings.ToLower(c.value), strings.ToLower(keyword))
}

// GetDisplayName 获取显示名称（可能包含格式化）
func (c CommunityName) GetDisplayName() string {
	return c.value
}

// GetSearchKeywords 获取搜索关键词
func (c CommunityName) GetSearchKeywords() []string {
	// 简单分词，实际项目中可能需要更复杂的分词算法
	keywords := strings.Fields(c.value)
	
	// 去除常见停用词
	stopWords := map[string]bool{
		"的": true, "是": true, "在": true, "有": true, "和": true,
		"与": true, "或": true, "但": true, "及": true, "等": true,
		"the": true, "is": true, "are": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
	}
	
	var result []string
	for _, keyword := range keywords {
		if len(keyword) > 1 && !stopWords[strings.ToLower(keyword)] {
			result = append(result, strings.ToLower(keyword))
		}
	}
	
	return result
}

// validateCommunityName 验证社区名称格式
func validateCommunityName(name string) error {
	// 长度检查
	length := utf8.RuneCountInString(name)
	if length < 2 {
		return errors.New("社区名称长度不能少于2个字符")
	}
	if length > 50 {
		return errors.New("社区名称长度不能超过50个字符")
	}
	
	// 字符组成检查：允许中文、英文、数字、空格、部分特殊字符
	validChars := regexp.MustCompile(`^[\p{L}\p{N}\s\-_.()（）【】\[\]&]+$`)
	if !validChars.MatchString(name) {
		return errors.New("社区名称包含不允许的字符")
	}
	
	// 不能全部是空格
	if strings.TrimSpace(name) == "" {
		return errors.New("社区名称不能全部是空格")
	}
	
	// 不能以特殊字符开头或结尾
	trimmed := strings.TrimSpace(name)
	if strings.HasPrefix(trimmed, "-") || strings.HasSuffix(trimmed, "-") ||
		strings.HasPrefix(trimmed, "_") || strings.HasSuffix(trimmed, "_") {
		return errors.New("社区名称不能以特殊字符开头或结尾")
	}
	
	// 敏感词检查（简化版本）
	if containsSensitiveWords(name) {
		return errors.New("社区名称包含敏感词，请重新输入")
	}
	
	return nil
}

// containsSensitiveWords 检查是否包含敏感词（简化版本）
func containsSensitiveWords(name string) bool {
	// 这里应该使用更完善的敏感词过滤系统
	sensitiveWords := []string{
		"色情", "赌博", "毒品", "暴力", "反动",
		"porn", "gambling", "drug", "violence",
		// 可以添加更多敏感词
	}
	
	lowerName := strings.ToLower(name)
	for _, word := range sensitiveWords {
		if strings.Contains(lowerName, word) {
			return true
		}
	}
	
	return false
} 