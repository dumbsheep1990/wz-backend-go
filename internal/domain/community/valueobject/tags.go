package valueobject

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Tags 标签集合值对象
type Tags struct {
	values []string
}

// NewTags 创建标签集合值对象
func NewTags(tags []string) (Tags, error) {
	if len(tags) == 0 {
		return Tags{values: []string{}}, nil
	}

	validatedTags, err := validateTags(tags)
	if err != nil {
		return Tags{}, err
	}

	return Tags{values: validatedTags}, nil
}

// MustNewTags 创建标签集合值对象，如果无效则panic
func MustNewTags(tags []string) Tags {
	t, err := NewTags(tags)
	if err != nil {
		panic("无效的标签: " + err.Error())
	}
	return t
}

// NewEmptyTags 创建空的标签集合
func NewEmptyTags() Tags {
	return Tags{values: []string{}}
}

// Values 获取标签值列表
func (t Tags) Values() []string {
	// 返回副本，避免外部修改
	result := make([]string, len(t.values))
	copy(result, t.values)
	return result
}

// IsEmpty 检查是否为空
func (t Tags) IsEmpty() bool {
	return len(t.values) == 0
}

// Count 获取标签数量
func (t Tags) Count() int {
	return len(t.values)
}

// Contains 检查是否包含特定标签
func (t Tags) Contains(tag string) bool {
	normalizedTag := strings.ToLower(strings.TrimSpace(tag))
	for _, v := range t.values {
		if strings.ToLower(v) == normalizedTag {
			return true
		}
	}
	return false
}

// ContainsAny 检查是否包含任一标签
func (t Tags) ContainsAny(tags []string) bool {
	for _, tag := range tags {
		if t.Contains(tag) {
			return true
		}
	}
	return false
}

// ContainsAll 检查是否包含所有标签
func (t Tags) ContainsAll(tags []string) bool {
	for _, tag := range tags {
		if !t.Contains(tag) {
			return false
		}
	}
	return true
}

// Add 添加标签
func (t Tags) Add(tag string) (Tags, error) {
	tag = strings.TrimSpace(tag)
	if err := validateTag(tag); err != nil {
		return t, err
	}

	// 检查重复
	if t.Contains(tag) {
		return t, nil // 已存在，不重复添加
	}

	// 检查数量限制
	if len(t.values) >= 10 {
		return t, errors.New("标签数量不能超过10个")
	}

	newValues := make([]string, len(t.values)+1)
	copy(newValues, t.values)
	newValues[len(t.values)] = tag

	return Tags{values: newValues}, nil
}

// Remove 移除标签
func (t Tags) Remove(tag string) Tags {
	normalizedTag := strings.ToLower(strings.TrimSpace(tag))
	var newValues []string

	for _, v := range t.values {
		if strings.ToLower(v) != normalizedTag {
			newValues = append(newValues, v)
		}
	}

	return Tags{values: newValues}
}

// Merge 合并其他标签集合
func (t Tags) Merge(other Tags) (Tags, error) {
	allTags := append(t.values, other.values...)
	return NewTags(allTags)
}

// FilterBy 根据关键词过滤标签
func (t Tags) FilterBy(keyword string) Tags {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return t
	}

	var filtered []string
	for _, tag := range t.values {
		if strings.Contains(strings.ToLower(tag), keyword) {
			filtered = append(filtered, tag)
		}
	}

	return Tags{values: filtered}
}

// ToStringSlice 转换为字符串切片
func (t Tags) ToStringSlice() []string {
	return t.Values()
}

// ToString 转换为逗号分隔的字符串
func (t Tags) ToString() string {
	return strings.Join(t.values, ",")
}

// GetPopularTags 获取流行标签（按长度排序，简化实现）
func (t Tags) GetPopularTags(limit int) []string {
	if limit <= 0 || limit > len(t.values) {
		limit = len(t.values)
	}

	// 简单按字母排序，实际项目中可能按使用频率排序
	sorted := make([]string, len(t.values))
	copy(sorted, t.values)
	
	// 简单排序（实际项目中可能需要更复杂的排序算法）
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if strings.ToLower(sorted[i]) > strings.ToLower(sorted[j]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if limit > len(sorted) {
		limit = len(sorted)
	}
	return sorted[:limit]
}

// IsValidForCategory 检查标签是否适合特定分类
func (t Tags) IsValidForCategory(category string) bool {
	// 根据分类检查标签的合理性
	categoryMap := map[string][]string{
		"技术": {"编程", "开发", "技术", "IT", "软件", "硬件", "网络", "数据库", "算法"},
		"生活": {"美食", "旅游", "健康", "运动", "娱乐", "家庭", "购物", "时尚"},
		"教育": {"学习", "考试", "课程", "培训", "教学", "知识", "技能", "证书"},
		"商业": {"创业", "投资", "金融", "营销", "管理", "商务", "经济", "市场"},
	}

	categoryKeywords, exists := categoryMap[category]
	if !exists {
		return true // 未知分类，允许所有标签
	}

	// 检查是否有任何标签与分类相关
	for _, tag := range t.values {
		tagLower := strings.ToLower(tag)
		for _, keyword := range categoryKeywords {
			if strings.Contains(tagLower, strings.ToLower(keyword)) {
				return true
			}
		}
	}

	return len(t.values) == 0 // 如果没有标签，也认为是有效的
}

// validateTags 验证标签列表
func validateTags(tags []string) ([]string, error) {
	if len(tags) > 10 {
		return nil, errors.New("标签数量不能超过10个")
	}

	var validTags []string
	seen := make(map[string]bool)

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue // 跳过空标签
		}

		if err := validateTag(tag); err != nil {
			return nil, err
		}

		// 去重（不区分大小写）
		tagLower := strings.ToLower(tag)
		if !seen[tagLower] {
			seen[tagLower] = true
			validTags = append(validTags, tag)
		}
	}

	return validTags, nil
}

// validateTag 验证单个标签
func validateTag(tag string) error {
	if tag == "" {
		return errors.New("标签不能为空")
	}

	// 长度检查
	length := utf8.RuneCountInString(tag)
	if length < 1 {
		return errors.New("标签长度不能少于1个字符")
	}
	if length > 20 {
		return errors.New("标签长度不能超过20个字符")
	}

	// 字符组成检查：允许中文、英文、数字、部分特殊字符
	validChars := regexp.MustCompile(`^[\p{L}\p{N}\-_.()（）]+$`)
	if !validChars.MatchString(tag) {
		return errors.New("标签包含不允许的字符")
	}

	// 不能以特殊字符开头或结尾
	if strings.HasPrefix(tag, "-") || strings.HasSuffix(tag, "-") ||
		strings.HasPrefix(tag, "_") || strings.HasSuffix(tag, "_") ||
		strings.HasPrefix(tag, ".") || strings.HasSuffix(tag, ".") {
		return errors.New("标签不能以特殊字符开头或结尾")
	}

	// 敏感词检查（简化版本）
	if containsTagSensitiveWords(tag) {
		return errors.New("标签包含敏感词，请重新输入")
	}

	return nil
}

// containsTagSensitiveWords 检查标签是否包含敏感词
func containsTagSensitiveWords(tag string) bool {
	// 这里应该使用更完善的敏感词过滤系统
	sensitiveWords := []string{
		"色情", "赌博", "毒品", "暴力", "反动",
		"porn", "gambling", "drug", "violence",
		// 可以添加更多敏感词
	}

	tagLower := strings.ToLower(tag)
	for _, word := range sensitiveWords {
		if strings.Contains(tagLower, word) {
			return true
		}
	}

	return false
} 