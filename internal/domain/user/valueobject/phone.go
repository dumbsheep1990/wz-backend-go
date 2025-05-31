package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Phone 手机号值对象
type Phone struct {
	value string
}

// NewPhone 创建手机号值对象
func NewPhone(phone string) (Phone, error) {
	if phone == "" {
		return Phone{}, errors.New("手机号不能为空")
	}

	phone = strings.TrimSpace(phone)
	if !isValidPhone(phone) {
		return Phone{}, errors.New("手机号格式不正确")
	}

	return Phone{value: phone}, nil
}

// MustNewPhone 创建手机号值对象，如果无效则panic
func MustNewPhone(phone string) Phone {
	p, err := NewPhone(phone)
	if err != nil {
		panic("无效的手机号: " + err.Error())
	}
	return p
}

// Value 获取手机号值
func (p Phone) Value() string {
	return p.value
}

// IsEquals 比较两个手机号是否相等
func (p Phone) IsEquals(other Phone) bool {
	return p.value == other.value
}

// String 获取手机号的字符串表示
func (p Phone) String() string {
	return p.value
}

// IsValid 检查手机号是否有效
func (p Phone) IsValid() bool {
	return p.value != "" && isValidPhone(p.value)
}

// CountryCode 获取国家代码（简单实现）
func (p Phone) CountryCode() string {
	if strings.HasPrefix(p.value, "+86") {
		return "86"
	}
	if strings.HasPrefix(p.value, "+1") {
		return "1"
	}
	// 默认中国大陆
	return "86"
}

// Carrier 获取运营商（简化实现）
func (p Phone) Carrier() string {
	// 去除国家代码前缀
	phone := p.value
	if strings.HasPrefix(phone, "+86") {
		phone = phone[3:]
	}
	
	if len(phone) != 11 {
		return "未知"
	}
	
	prefix := phone[:3]
	switch prefix {
	case "130", "131", "132", "145", "155", "156", "166", "171", "175", "176", "185", "186", "196":
		return "中国联通"
	case "133", "134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "157", "158", "159", "172", "178", "182", "183", "184", "187", "188", "195", "198":
		return "中国移动"
	case "149", "153", "173", "174", "177", "180", "181", "189", "191", "193", "199":
		return "中国电信"
	default:
		return "未知"
	}
}

// Format 格式化手机号显示
func (p Phone) Format() string {
	phone := p.value
	if strings.HasPrefix(phone, "+86") {
		phone = phone[3:]
	}
	
	if len(phone) == 11 {
		return phone[:3] + " " + phone[3:7] + " " + phone[7:]
	}
	
	return p.value
}

// Mask 脱敏显示手机号
func (p Phone) Mask() string {
	phone := p.value
	if strings.HasPrefix(phone, "+86") {
		phone = phone[3:]
	}
	
	if len(phone) == 11 {
		return phone[:3] + "****" + phone[7:]
	}
	
	// 对于其他格式，简单处理
	if len(phone) > 4 {
		return phone[:2] + strings.Repeat("*", len(phone)-4) + phone[len(phone)-2:]
	}
	
	return phone
}

// IsEmpty 检查是否为空
func (p Phone) IsEmpty() bool {
	return p.value == ""
}

// IsMobile 判断是否为移动电话（区别于固定电话）
func (p Phone) IsMobile() bool {
	phone := p.value
	if strings.HasPrefix(phone, "+86") {
		phone = phone[3:]
	}
	
	// 中国大陆手机号以1开头且长度为11位
	return len(phone) == 11 && strings.HasPrefix(phone, "1")
}

// isValidPhone 验证手机号格式
func isValidPhone(phone string) bool {
	// 去除空格和连字符
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	// 支持国际格式
	if strings.HasPrefix(phone, "+") {
		// 国际格式验证
		phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
		return phoneRegex.MatchString(phone)
	}
	
	// 中国大陆手机号验证
	if len(phone) == 11 && strings.HasPrefix(phone, "1") {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		return phoneRegex.MatchString(phone)
	}
	
	// 简单的数字验证（作为后备）
	if len(phone) >= 7 && len(phone) <= 15 {
		phoneRegex := regexp.MustCompile(`^\d+$`)
		return phoneRegex.MatchString(phone)
	}
	
	return false
} 