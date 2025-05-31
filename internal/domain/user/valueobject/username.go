package valueobject

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Username 用户名值对象
type Username struct {
	value string
}

// NewUsername 创建用户名值对象
func NewUsername(username string) (Username, error) {
	if username == "" {
		return Username{}, errors.New("用户名不能为空")
	}

	username = strings.TrimSpace(username)
	if err := validateUsername(username); err != nil {
		return Username{}, err
	}

	return Username{value: username}, nil
}

// MustNewUsername 创建用户名值对象，如果无效则panic
func MustNewUsername(username string) Username {
	u, err := NewUsername(username)
	if err != nil {
		panic("无效的用户名: " + err.Error())
	}
	return u
}

// Value 获取用户名值
func (u Username) Value() string {
	return u.value
}

// IsEquals 比较两个用户名是否相等
func (u Username) IsEquals(other Username) bool {
	return strings.ToLower(u.value) == strings.ToLower(other.value)
}

// String 获取用户名的字符串表示
func (u Username) String() string {
	return u.value
}

// IsValid 检查用户名是否有效
func (u Username) IsValid() bool {
	return u.value != "" && validateUsername(u.value) == nil
}

// Length 获取用户名长度（考虑中文字符）
func (u Username) Length() int {
	return utf8.RuneCountInString(u.value)
}

// IsEmpty 检查是否为空
func (u Username) IsEmpty() bool {
	return u.value == ""
}

// ToLower 转换为小写
func (u Username) ToLower() string {
	return strings.ToLower(u.value)
}

// ContainsSpecialChars 检查是否包含特殊字符
func (u Username) ContainsSpecialChars() bool {
	specialChars := regexp.MustCompile(`[^a-zA-Z0-9\u4e00-\u9fa5_]`)
	return specialChars.MatchString(u.value)
}

// IsAlphaNumeric 检查是否只包含字母和数字
func (u Username) IsAlphaNumeric() bool {
	alphaNumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphaNumeric.MatchString(u.value)
}

// HasChinese 检查是否包含中文字符
func (u Username) HasChinese() bool {
	chinese := regexp.MustCompile(`[\u4e00-\u9fa5]`)
	return chinese.MatchString(u.value)
}

// Mask 脱敏显示用户名（用于日志等）
func (u Username) Mask() string {
	if u.value == "" {
		return ""
	}
	
	length := u.Length()
	if length <= 2 {
		return u.value
	}
	
	if length <= 4 {
		return u.value[:1] + "*" + u.value[len(u.value)-1:]
	}
	
	// 保留前2位和后1位
	runes := []rune(u.value)
	masked := string(runes[:2]) + strings.Repeat("*", length-3) + string(runes[length-1:])
	return masked
}

// validateUsername 验证用户名格式
func validateUsername(username string) error {
	// 长度检查
	length := utf8.RuneCountInString(username)
	if length < 2 {
		return errors.New("用户名长度不能少于2个字符")
	}
	if length > 20 {
		return errors.New("用户名长度不能超过20个字符")
	}
	
	// 字符组成检查：只允许字母、数字、中文、下划线
	validChars := regexp.MustCompile(`^[a-zA-Z0-9\u4e00-\u9fa5_]+$`)
	if !validChars.MatchString(username) {
		return errors.New("用户名只能包含字母、数字、中文和下划线")
	}
	
	// 不能以下划线开头或结尾
	if strings.HasPrefix(username, "_") || strings.HasSuffix(username, "_") {
		return errors.New("用户名不能以下划线开头或结尾")
	}
	
	// 不能全部是数字
	allNumbers := regexp.MustCompile(`^\d+$`)
	if allNumbers.MatchString(username) {
		return errors.New("用户名不能全部是数字")
	}
	
	// 保留字检查
	if isReservedUsername(username) {
		return errors.New("该用户名为系统保留，请选择其他用户名")
	}
	
	// 敏感词检查（简化版本）
	if containsSensitiveWords(username) {
		return errors.New("用户名包含敏感词，请重新输入")
	}
	
	return nil
}

// isReservedUsername 检查是否为保留用户名
func isReservedUsername(username string) bool {
	reserved := []string{
		"admin", "root", "system", "guest", "user", "test",
		"api", "www", "mail", "ftp", "support", "help",
		"null", "undefined", "unknown", "anonymous",
		"管理员", "系统", "测试", "客服", "帮助",
	}
	
	lowerUsername := strings.ToLower(username)
	for _, word := range reserved {
		if lowerUsername == word {
			return true
		}
	}
	
	return false
}

// containsSensitiveWords 检查是否包含敏感词（简化版本）
func containsSensitiveWords(username string) bool {
	// 这里应该使用更完善的敏感词过滤系统
	sensitiveWords := []string{
		"fuck", "shit", "admin", "superuser",
		// 可以添加更多敏感词
	}
	
	lowerUsername := strings.ToLower(username)
	for _, word := range sensitiveWords {
		if strings.Contains(lowerUsername, word) {
			return true
		}
	}
	
	return false
} 