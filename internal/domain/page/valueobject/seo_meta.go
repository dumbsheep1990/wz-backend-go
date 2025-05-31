package valueobject

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// SEOMeta SEO元数据值对象
type SEOMeta struct {
	description string
	keywords    []string
}

// NewSEOMeta 创建SEO元数据值对象
func NewSEOMeta(description string, keywords []string) (SEOMeta, error) {
	// 清理描述
	description = strings.TrimSpace(description)
	if description != "" {
		if err := validateSEODescription(description); err != nil {
			return SEOMeta{}, err
		}
		// HTML转义防止XSS
		description = html.EscapeString(description)
	}

	// 清理关键词
	cleanKeywords, err := validateAndCleanKeywords(keywords)
	if err != nil {
		return SEOMeta{}, err
	}

	return SEOMeta{
		description: description,
		keywords:    cleanKeywords,
	}, nil
}

// EmptySEOMeta 创建空的SEO元数据
func EmptySEOMeta() SEOMeta {
	return SEOMeta{
		description: "",
		keywords:    []string{},
	}
}

// MustNewSEOMeta 创建SEO元数据值对象，如果无效则panic
func MustNewSEOMeta(description string, keywords []string) SEOMeta {
	seo, err := NewSEOMeta(description, keywords)
	if err != nil {
		panic("无效的SEO元数据: " + err.Error())
	}
	return seo
}

// Description 获取SEO描述
func (s SEOMeta) Description() string {
	return s.description
}

// Keywords 获取SEO关键词
func (s SEOMeta) Keywords() []string {
	// 返回副本，避免外部修改
	result := make([]string, len(s.keywords))
	copy(result, s.keywords)
	return result
}

// IsEmpty 检查SEO元数据是否为空
func (s SEOMeta) IsEmpty() bool {
	return s.description == "" && len(s.keywords) == 0
}

// GetUnescapedDescription 获取未转义的描述（用于显示）
func (s SEOMeta) GetUnescapedDescription() string {
	return html.UnescapeString(s.description)
}

// UpdateDescription 更新SEO描述
func (s SEOMeta) UpdateDescription(description string) (SEOMeta, error) {
	return NewSEOMeta(description, s.keywords)
}

// UpdateKeywords 更新SEO关键词
func (s SEOMeta) UpdateKeywords(keywords []string) (SEOMeta, error) {
	return NewSEOMeta(s.description, keywords)
}

// AddKeyword 添加关键词
func (s SEOMeta) AddKeyword(keyword string) (SEOMeta, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return s, errors.New("关键词不能为空")
	}

	// 检查是否已存在（不区分大小写）
	lowerKeyword := strings.ToLower(keyword)
	for _, k := range s.keywords {
		if strings.ToLower(k) == lowerKeyword {
			return s, nil // 已存在，不重复添加
		}
	}

	// 检查关键词总数限制
	if len(s.keywords) >= 10 {
		return s, errors.New("关键词数量不能超过10个")
	}

	newKeywords := append(s.keywords, keyword)
	return NewSEOMeta(s.description, newKeywords)
}

// RemoveKeyword 移除关键词
func (s SEOMeta) RemoveKeyword(keyword string) SEOMeta {
	lowerKeyword := strings.ToLower(strings.TrimSpace(keyword))
	var newKeywords []string

	for _, k := range s.keywords {
		if strings.ToLower(k) != lowerKeyword {
			newKeywords = append(newKeywords, k)
		}
	}

	return SEOMeta{
		description: s.description,
		keywords:    newKeywords,
	}
}

// GetKeywordsString 获取关键词的逗号分隔字符串
func (s SEOMeta) GetKeywordsString() string {
	return strings.Join(s.keywords, ", ")
}

// GetHTMLMetaTags 获取HTML meta标签
func (s SEOMeta) GetHTMLMetaTags() map[string]string {
	tags := make(map[string]string)

	if s.description != "" {
		tags["description"] = s.GetUnescapedDescription()
	}

	if len(s.keywords) > 0 {
		tags["keywords"] = s.GetKeywordsString()
	}

	return tags
}

// IsDescriptionOptimal 检查描述是否符合SEO最佳实践
func (s SEOMeta) IsDescriptionOptimal() bool {
	if s.description == "" {
		return false
	}

	unescaped := s.GetUnescapedDescription()
	length := utf8.RuneCountInString(unescaped)

	// SEO建议的描述长度为120-160字符
	return length >= 120 && length <= 160
}

// GetDescriptionAnalysis 获取描述的SEO分析
func (s SEOMeta) GetDescriptionAnalysis() map[string]interface{} {
	unescaped := s.GetUnescapedDescription()
	length := utf8.RuneCountInString(unescaped)

	status := "missing"
	suggestion := "建议添加页面描述"

	if s.description != "" {
		if length < 120 {
			status = "too_short"
			suggestion = "描述过短，建议增加到120-160字符"
		} else if length > 160 {
			status = "too_long"
			suggestion = "描述过长，可能在搜索结果中被截断"
		} else {
			status = "optimal"
			suggestion = "描述长度最佳，有利于SEO"
		}
	}

	return map[string]interface{}{
		"length":           length,
		"recommended_min":  120,
		"recommended_max":  160,
		"status":          status,
		"suggestion":      suggestion,
		"is_optimal":      s.IsDescriptionOptimal(),
	}
}

// GetKeywordsAnalysis 获取关键词的SEO分析
func (s SEOMeta) GetKeywordsAnalysis() map[string]interface{} {
	count := len(s.keywords)
	
	status := "missing"
	suggestion := "建议添加3-8个相关关键词"

	if count > 0 {
		if count < 3 {
			status = "too_few"
			suggestion = "关键词过少，建议增加到3-8个"
		} else if count > 8 {
			status = "too_many"
			suggestion = "关键词过多，建议控制在3-8个以内"
		} else {
			status = "optimal"
			suggestion = "关键词数量最佳"
		}
	}

	return map[string]interface{}{
		"count":           count,
		"recommended_min": 3,
		"recommended_max": 8,
		"status":         status,
		"suggestion":     suggestion,
		"is_optimal":     count >= 3 && count <= 8,
	}
}

// ContainsKeyword 检查是否包含特定关键词
func (s SEOMeta) ContainsKeyword(keyword string) bool {
	lowerKeyword := strings.ToLower(strings.TrimSpace(keyword))
	for _, k := range s.keywords {
		if strings.ToLower(k) == lowerKeyword {
			return true
		}
	}
	return false
}

// GetSEOScore 获取SEO评分（0-100）
func (s SEOMeta) GetSEOScore() int {
	score := 0

	// 描述评分（50分）
	if s.description != "" {
		score += 25 // 有描述基础分
		if s.IsDescriptionOptimal() {
			score += 25 // 描述长度最佳
		} else {
			unescaped := s.GetUnescapedDescription()
			length := utf8.RuneCountInString(unescaped)
			if length >= 80 && length <= 200 {
				score += 15 // 描述长度可接受
			}
		}
	}

	// 关键词评分（50分）
	keywordCount := len(s.keywords)
	if keywordCount > 0 {
		score += 15 // 有关键词基础分
		if keywordCount >= 3 && keywordCount <= 8 {
			score += 35 // 关键词数量最佳
		} else if keywordCount >= 1 && keywordCount <= 10 {
			score += 20 // 关键词数量可接受
		}
	}

	return score
}

// GetSEOSuggestions 获取SEO优化建议
func (s SEOMeta) GetSEOSuggestions() []string {
	var suggestions []string

	// 描述建议
	descAnalysis := s.GetDescriptionAnalysis()
	if suggestion, ok := descAnalysis["suggestion"].(string); ok && suggestion != "" {
		suggestions = append(suggestions, suggestion)
	}

	// 关键词建议
	keywordAnalysis := s.GetKeywordsAnalysis()
	if suggestion, ok := keywordAnalysis["suggestion"].(string); ok && suggestion != "" {
		suggestions = append(suggestions, suggestion)
	}

	// 其他建议
	if s.description != "" && len(s.keywords) > 0 {
		// 检查描述中是否包含关键词
		unescaped := strings.ToLower(s.GetUnescapedDescription())
		hasKeywordInDesc := false
		for _, keyword := range s.keywords {
			if strings.Contains(unescaped, strings.ToLower(keyword)) {
				hasKeywordInDesc = true
				break
			}
		}
		if !hasKeywordInDesc {
			suggestions = append(suggestions, "建议在描述中包含一些关键词")
		}
	}

	return suggestions
}

// validateSEODescription 验证SEO描述
func validateSEODescription(description string) error {
	// 长度检查
	length := utf8.RuneCountInString(description)
	if length > 300 {
		return errors.New("SEO描述长度不能超过300个字符")
	}

	// 不能全部是空白字符
	if strings.TrimSpace(description) == "" {
		return errors.New("SEO描述不能全部是空白字符")
	}

	// 检查危险内容
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onclick=`,
		`onerror=`,
		`onload=`,
	}

	lowerDesc := strings.ToLower(description)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerDesc, pattern) {
			return errors.New("SEO描述包含不安全的内容")
		}
	}

	return nil
}

// validateAndCleanKeywords 验证并清理关键词
func validateAndCleanKeywords(keywords []string) ([]string, error) {
	if len(keywords) > 10 {
		return nil, errors.New("关键词数量不能超过10个")
	}

	var cleanKeywords []string
	seen := make(map[string]bool)

	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue // 跳过空关键词
		}

		// 验证关键词格式
		if err := validateKeyword(keyword); err != nil {
			return nil, err
		}

		// 去重（不区分大小写）
		lowerKeyword := strings.ToLower(keyword)
		if !seen[lowerKeyword] {
			seen[lowerKeyword] = true
			cleanKeywords = append(cleanKeywords, keyword)
		}
	}

	return cleanKeywords, nil
}

// validateKeyword 验证单个关键词
func validateKeyword(keyword string) error {
	// 长度检查
	length := utf8.RuneCountInString(keyword)
	if length < 1 {
		return errors.New("关键词长度不能少于1个字符")
	}
	if length > 30 {
		return errors.New("关键词长度不能超过30个字符")
	}

	// 字符组成检查：允许字母、数字、空格、连字符和中文
	validChars := regexp.MustCompile(`^[\p{L}\p{N}\s\-]+$`)
	if !validChars.MatchString(keyword) {
		return errors.New("关键词只能包含字母、数字、空格、连字符")
	}

	// 不能以空格或连字符开头/结尾
	if strings.HasPrefix(keyword, " ") || strings.HasSuffix(keyword, " ") ||
		strings.HasPrefix(keyword, "-") || strings.HasSuffix(keyword, "-") {
		return errors.New("关键词不能以空格或连字符开头/结尾")
	}

	return nil
} 