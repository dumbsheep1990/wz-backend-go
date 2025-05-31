package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Email 邮箱值对象
type Email struct {
	value string
}

// NewEmail 创建邮箱值对象
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("邮箱不能为空")
	}

	email = strings.TrimSpace(email)
	if !isValidEmail(email) {
		return Email{}, errors.New("邮箱格式不正确")
	}

	return Email{value: email}, nil
}

// MustNewEmail 创建邮箱值对象，如果无效则panic
func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic("无效的邮箱: " + err.Error())
	}
	return e
}

// Value 获取邮箱值
func (e Email) Value() string {
	return e.value
}

// Domain 获取邮箱域名
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart 获取邮箱本地部分
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// IsEquals 比较两个邮箱是否相等
func (e Email) IsEquals(other Email) bool {
	return strings.ToLower(e.value) == strings.ToLower(other.value)
}

// String 获取邮箱的字符串表示
func (e Email) String() string {
	return e.value
}

// IsValid 检查邮箱是否有效
func (e Email) IsValid() bool {
	return e.value != "" && isValidEmail(e.value)
}

// IsBusinessEmail 判断是否为企业邮箱
func (e Email) IsBusinessEmail() bool {
	// 简单判断，实际业务中可能需要更复杂的规则
	commonPersonalDomains := []string{
		"gmail.com", "163.com", "126.com", "qq.com", "sina.com",
		"hotmail.com", "yahoo.com", "sohu.com", "139.com",
	}
	
	domain := strings.ToLower(e.Domain())
	for _, personalDomain := range commonPersonalDomains {
		if domain == personalDomain {
			return false
		}
	}
	
	return true
}

// Mask 脱敏显示邮箱
func (e Email) Mask() string {
	if e.value == "" {
		return ""
	}
	
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return e.value
	}
	
	local := parts[0]
	domain := parts[1]
	
	if len(local) <= 2 {
		return local + "@" + domain
	}
	
	// 保留前2位和后1位，中间用*替代
	masked := local[:2] + strings.Repeat("*", len(local)-3) + local[len(local)-1:]
	return masked + "@" + domain
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	// 基本的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	
	// 长度限制
	if len(email) > 254 {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	local := parts[0]
	domain := parts[1]
	
	// 本地部分长度限制
	if len(local) > 64 {
		return false
	}
	
	// 域名长度限制
	if len(domain) > 253 {
		return false
	}
	
	return true
} 