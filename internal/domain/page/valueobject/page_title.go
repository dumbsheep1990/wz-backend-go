package valueobject

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// PageTitle 页面标题值对象
type PageTitle struct {
	value string
}

// NewPageTitle 创建页面标题值对象
func NewPageTitle(title string) (PageTitle, error) {
	if title == "" {
		return PageTitle{}, errors.New("页面标题不能为空")
	}

	title = strings.TrimSpace(title)
	if err := validatePageTitle(title); err != nil {
		return PageTitle{}, err
	}

	// HTML转义，防止XSS攻击
	title = html.EscapeString(title)

	return PageTitle{value: title}, nil
}

// MustNewPageTitle 创建页面标题值对象，如果无效则panic
func MustNewPageTitle(title string) PageTitle {
	pt, err := NewPageTitle(title)
	if err != nil {
		panic("无效的页面标题: " + err.Error())
	}
	return pt
}

// Value 获取页面标题值
func (p PageTitle) Value() string {
	return p.value
}

// IsEmpty 检查是否为空
func (p PageTitle) IsEmpty() bool {
	return p.value == ""
}

// IsEquals 比较两个页面标题是否相等
func (p PageTitle) IsEquals(other PageTitle) bool {
	return p.value == other.value
}

// String 获取页面标题的字符串表示
func (p PageTitle) String() string {
	return p.value
}

// IsValid 检查页面标题是否有效
func (p PageTitle) IsValid() bool {
	return p.value != "" && validatePageTitle(p.value) == nil
}

// Length 获取页面标题长度（考虑中文字符）
func (p PageTitle) Length() int {
	return utf8.RuneCountInString(p.value)
}

// GetUnescaped 获取未转义的标题（用于显示）
func (p PageTitle) GetUnescaped() string {
	return html.UnescapeString(p.value)
}

// GetSEOTitle 获取SEO优化的标题
func (p PageTitle) GetSEOTitle(siteName string) string {
	title := p.GetUnescaped()
	if siteName != "" && !strings.Contains(title, siteName) {
		return title + " - " + siteName
	}
	return title
}

// GetDisplayTitle 获取显示标题（截断过长的标题）
func (p PageTitle) GetDisplayTitle(maxLength int) string {
	title := p.GetUnescaped()
	if maxLength <= 0 {
		maxLength = 60 // 默认最大长度
	}
	
	if utf8.RuneCountInString(title) <= maxLength {
		return title
	}
	
	runes := []rune(title)
	return string(runes[:maxLength-3]) + "..."
}

// GetTitleForURL 获取适合URL的标题（去除特殊字符）
func (p PageTitle) GetTitleForURL() string {
	title := p.GetUnescaped()
	
	// 移除HTML标签
	reg := regexp.MustCompile(`<[^>]*>`)
	title = reg.ReplaceAllString(title, "")
	
	// 移除特殊字符
	reg = regexp.MustCompile(`[^\p{L}\p{N}\s\-]`)
	title = reg.ReplaceAllString(title, "")
	
	return strings.TrimSpace(title)
}

// ContainsKeyword 检查是否包含特定关键词
func (p PageTitle) ContainsKeyword(keyword string) bool {
	if keyword == "" {
		return false
	}
	
	title := strings.ToLower(p.GetUnescaped())
	keyword = strings.ToLower(keyword)
	
	return strings.Contains(title, keyword)
}

// GetKeywords 提取标题中的关键词
func (p PageTitle) GetKeywords() []string {
	title := p.GetUnescaped()
	
	// 移除标点符号，保留字母、数字、空格和中文
	reg := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	title = reg.ReplaceAllString(title, " ")
	
	// 分割单词
	words := strings.Fields(title)
	
	// 去除停用词和过短的词
	stopWords := map[string]bool{
		"的": true, "是": true, "在": true, "有": true, "和": true,
		"与": true, "或": true, "但": true, "及": true, "等": true,
		"the": true, "is": true, "are": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"a": true, "an": true, "as": true, "by": true, "for": true,
		"of": true, "with": true, "this": true, "that": true,
	}
	
	var keywords []string
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if len(word) > 1 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// IsSEOOptimized 检查标题是否符合SEO最佳实践
func (p PageTitle) IsSEOOptimized() bool {
	length := p.Length()
	
	// SEO建议的标题长度为30-60字符
	if length < 10 || length > 60 {
		return false
	}
	
	title := p.GetUnescaped()
	
	// 不应该全部大写
	if title == strings.ToUpper(title) {
		return false
	}
	
	// 不应该包含过多的特殊字符
	specialCharCount := 0
	for _, r := range title {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || r == ' ' || r == '-' || 
			 (r >= 0x4e00 && r <= 0x9fff)) { // 中文字符范围
			specialCharCount++
		}
	}
	
	// 特殊字符不应超过总字符数的20%
	if float64(specialCharCount)/float64(length) > 0.2 {
		return false
	}
	
	return true
}

// GetSuggestedLength 获取建议的标题长度信息
func (p PageTitle) GetSuggestedLength() map[string]interface{} {
	length := p.Length()
	
	return map[string]interface{}{
		"current_length":    length,
		"recommended_min":   10,
		"recommended_max":   60,
		"is_optimal":        length >= 10 && length <= 60,
		"status":           p.getLengthStatus(length),
		"suggestion":       p.getLengthSuggestion(length),
	}
}

// getLengthStatus 获取长度状态
func (p PageTitle) getLengthStatus(length int) string {
	if length < 10 {
		return "too_short"
	} else if length > 60 {
		return "too_long"
	} else if length >= 30 && length <= 50 {
		return "optimal"
	} else {
		return "acceptable"
	}
}

// getLengthSuggestion 获取长度建议
func (p PageTitle) getLengthSuggestion(length int) string {
	if length < 10 {
		return "标题过短，建议增加更多描述性内容"
	} else if length > 60 {
		return "标题过长，可能在搜索结果中被截断"
	} else if length >= 30 && length <= 50 {
		return "标题长度最佳，有利于SEO"
	} else {
		return "标题长度可以接受"
	}
}

// validatePageTitle 验证页面标题格式
func validatePageTitle(title string) error {
	// 长度检查
	length := utf8.RuneCountInString(title)
	if length < 1 {
		return errors.New("页面标题长度不能少于1个字符")
	}
	if length > 200 {
		return errors.New("页面标题长度不能超过200个字符")
	}
	
	// 不能全部是空白字符
	if strings.TrimSpace(title) == "" {
		return errors.New("页面标题不能全部是空白字符")
	}
	
	// 不允许某些危险字符（防止XSS，但允许基本HTML字符）
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onclick=`,
		`onerror=`,
		`onload=`,
	}
	
	lowerTitle := strings.ToLower(title)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerTitle, pattern) {
			return errors.New("页面标题包含不安全的内容")
		}
	}
	
	return nil
} 