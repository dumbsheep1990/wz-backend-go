package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"wz-backend-go/internal/domain/page/valueobject"
)

// DefaultSlugGenerator 默认URL段生成器实现
type DefaultSlugGenerator struct {
	pageRepo PageRepository
}

// NewDefaultSlugGenerator 创建默认URL段生成器
func NewDefaultSlugGenerator(pageRepo PageRepository) *DefaultSlugGenerator {
	return &DefaultSlugGenerator{
		pageRepo: pageRepo,
	}
}

// GenerateFromTitle 从标题生成URL段
func (g *DefaultSlugGenerator) GenerateFromTitle(title string) (valueobject.PageSlug, error) {
	if title == "" {
		return valueobject.PageSlug{}, fmt.Errorf("标题不能为空")
	}

	// 1. 转为小写
	slug := strings.ToLower(title)

	// 2. 移除HTML标签
	htmlReg := regexp.MustCompile(`<[^>]*>`)
	slug = htmlReg.ReplaceAllString(slug, "")

	// 3. 处理中文和特殊字符
	slug = g.processChinese(slug)

	// 4. 替换空格和特殊字符为连字符
	spaceReg := regexp.MustCompile(`[\s\-_]+`)
	slug = spaceReg.ReplaceAllString(slug, "-")

	// 5. 移除非字母数字和连字符的字符
	cleanReg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = cleanReg.ReplaceAllString(slug, "")

	// 6. 移除开头和结尾的连字符
	slug = strings.Trim(slug, "-")

	// 7. 限制长度
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimRight(slug, "-")
	}

	// 8. 如果为空，使用默认值
	if slug == "" {
		slug = "page"
	}

	return valueobject.NewPageSlug(slug)
}

// EnsureUnique 确保URL段的唯一性
func (g *DefaultSlugGenerator) EnsureUnique(ctx context.Context, siteID string, baseSlug valueobject.PageSlug, excludePageID *valueobject.PageID) (valueobject.PageSlug, error) {
	originalSlug := baseSlug.Value()
	candidateSlug := baseSlug
	counter := 1

	for {
		// 检查当前候选URL段是否已存在
		exists, err := g.pageRepo.ExistsBySlug(ctx, siteID, candidateSlug, excludePageID)
		if err != nil {
			return valueobject.PageSlug{}, fmt.Errorf("检查URL段唯一性失败: %w", err)
		}

		if !exists {
			return candidateSlug, nil
		}

		// 生成新的候选URL段
		newSlugValue := fmt.Sprintf("%s-%d", originalSlug, counter)
		candidateSlug, err = valueobject.NewPageSlug(newSlugValue)
		if err != nil {
			return valueobject.PageSlug{}, fmt.Errorf("生成唯一URL段失败: %w", err)
		}

		counter++

		// 防止无限循环
		if counter > 1000 {
			return valueobject.PageSlug{}, fmt.Errorf("无法生成唯一URL段")
		}
	}
}

// processChinese 处理中文字符
func (g *DefaultSlugGenerator) processChinese(text string) string {
	var result strings.Builder
	
	for _, r := range text {
		if g.isChinese(r) {
			// 对于中文字符，可以使用拼音转换库或者简单替换
			// 这里简化处理，使用连字符替代
			result.WriteString("-")
		} else {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// isChinese 判断是否为中文字符
func (g *DefaultSlugGenerator) isChinese(r rune) bool {
	return unicode.Is(unicode.Scripts["Han"], r)
}

// GenerateSEOFriendlySlug 生成SEO友好的URL段
func (g *DefaultSlugGenerator) GenerateSEOFriendlySlug(title string, keywords []string) (valueobject.PageSlug, error) {
	if title == "" {
		return valueobject.PageSlug{}, fmt.Errorf("标题不能为空")
	}

	// 1. 从标题生成基础URL段
	baseSlug, err := g.GenerateFromTitle(title)
	if err != nil {
		return valueobject.PageSlug{}, err
	}

	slug := baseSlug.Value()

	// 2. 如果有关键词，尝试融入关键词
	if len(keywords) > 0 {
		// 选择最重要的关键词（第一个）
		primaryKeyword := strings.ToLower(keywords[0])
		primaryKeyword = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(primaryKeyword, "")

		// 如果URL段中不包含关键词，尝试添加
		if primaryKeyword != "" && !strings.Contains(slug, primaryKeyword) {
			// 如果URL段较短，可以添加关键词
			if len(slug)+len(primaryKeyword)+1 <= 80 {
				slug = primaryKeyword + "-" + slug
			}
		}
	}

	// 3. 确保符合SEO最佳实践
	slug = g.optimizeForSEO(slug)

	return valueobject.NewPageSlug(slug)
}

// optimizeForSEO 针对SEO优化URL段
func (g *DefaultSlugGenerator) optimizeForSEO(slug string) string {
	// 1. 确保长度在合理范围内（3-80字符）
	if len(slug) < 3 {
		slug = "page-" + slug
	}
	if len(slug) > 80 {
		slug = slug[:80]
		slug = strings.TrimRight(slug, "-")
	}

	// 2. 移除多余的连字符
	multiDashReg := regexp.MustCompile(`-+`)
	slug = multiDashReg.ReplaceAllString(slug, "-")

	// 3. 移除开头和结尾的连字符
	slug = strings.Trim(slug, "-")

	// 4. 确保不为空
	if slug == "" {
		slug = "page"
	}

	return slug
}

// GenerateFromContent 从内容生成URL段
func (g *DefaultSlugGenerator) GenerateFromContent(content string) (valueobject.PageSlug, error) {
	if content == "" {
		return valueobject.NewPageSlug("page")
	}

	// 提取内容的前几个单词作为URL段
	words := strings.Fields(content)
	if len(words) == 0 {
		return valueobject.NewPageSlug("page")
	}

	// 取前3-5个有意义的单词
	var selectedWords []string
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "this": true,
		"that": true, "is": true, "are": true, "was": true, "were": true,
		"的": true, "是": true, "在": true, "有": true, "和": true,
		"与": true, "或": true, "但": true, "及": true, "等": true,
	}

	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if len(word) > 2 && !stopWords[word] {
			selectedWords = append(selectedWords, word)
			if len(selectedWords) >= 5 {
				break
			}
		}
	}

	if len(selectedWords) == 0 {
		return valueobject.NewPageSlug("page")
	}

	title := strings.Join(selectedWords, " ")
	return g.GenerateFromTitle(title)
}

// ValidateSlug 验证URL段是否符合规范
func (g *DefaultSlugGenerator) ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("URL段不能为空")
	}

	// 长度检查
	if len(slug) < 1 || len(slug) > 100 {
		return fmt.Errorf("URL段长度必须在1-100字符之间")
	}

	// 格式检查
	validReg := regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$`)
	if !validReg.MatchString(slug) {
		return fmt.Errorf("URL段格式不正确，只能包含小写字母、数字和连字符，且不能以连字符开头或结尾")
	}

	// 检查保留词
	reservedWords := []string{
		"admin", "api", "www", "mail", "ftp", "localhost", "root",
		"administrator", "moderator", "webmaster", "hostmaster",
		"postmaster", "abuse", "noreply", "support", "sales",
		"info", "contact", "help", "blog", "news", "about",
		"home", "index", "default", "main", "page", "post",
		"article", "category", "tag", "search", "login", "logout",
		"register", "signup", "signin", "dashboard", "profile",
		"settings", "config", "test", "demo", "example", "sample",
	}

	for _, reserved := range reservedWords {
		if strings.EqualFold(slug, reserved) {
			return fmt.Errorf("URL段不能使用保留词: %s", reserved)
		}
	}

	return nil
}

// SuggestAlternatives 为给定的URL段建议替代方案
func (g *DefaultSlugGenerator) SuggestAlternatives(ctx context.Context, siteID string, originalSlug string, count int) ([]valueobject.PageSlug, error) {
	if count <= 0 {
		count = 5
	}

	var suggestions []valueobject.PageSlug

	// 1. 基于原始URL段的变体
	baseSlug, err := valueobject.NewPageSlug(originalSlug)
	if err == nil {
		// 添加数字后缀
		for i := 1; i <= count && len(suggestions) < count; i++ {
			candidate := fmt.Sprintf("%s-%d", originalSlug, i)
			if slug, err := valueobject.NewPageSlug(candidate); err == nil {
				exists, err := g.pageRepo.ExistsBySlug(ctx, siteID, slug, nil)
				if err == nil && !exists {
					suggestions = append(suggestions, slug)
				}
			}
		}
	}

	// 2. 基于同义词的变体
	synonyms := []string{"page", "content", "article", "post", "item", "entry"}
	for _, synonym := range synonyms {
		if len(suggestions) >= count {
			break
		}
		
		candidate := fmt.Sprintf("%s-%s", synonym, originalSlug)
		if slug, err := valueobject.NewPageSlug(candidate); err == nil {
			exists, err := g.pageRepo.ExistsBySlug(ctx, siteID, slug, nil)
			if err == nil && !exists {
				suggestions = append(suggestions, slug)
			}
		}
	}

	// 3. 基于当前时间戳的变体
	if len(suggestions) < count {
		timestamp := strconv.FormatInt(1234567890, 10) // 简化时间戳
		candidate := fmt.Sprintf("%s-%s", originalSlug, timestamp[len(timestamp)-4:])
		if slug, err := valueobject.NewPageSlug(candidate); err == nil {
			exists, err := g.pageRepo.ExistsBySlug(ctx, siteID, slug, nil)
			if err == nil && !exists {
				suggestions = append(suggestions, slug)
			}
		}
	}

	return suggestions, nil
} 