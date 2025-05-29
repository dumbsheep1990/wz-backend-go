package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Username 用户名值对象
type Username string

// 用户名正则表达式：4-20个字符，只能包含字母、数字和下划线
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)

// NewUsername 创建用户名值对象
func NewUsername(username string) (Username, error) {
	if username == "" {
		return "", errors.New("用户名不能为空")
	}

	if !usernameRegex.MatchString(username) {
		return "", errors.New("用户名必须是4-20个字符，只能包含字母、数字和下划线")
	}

	return Username(username), nil
}

// Value 获取用户名值
func (u Username) Value() string {
	return string(u)
}

// String 字符串表示
func (u Username) String() string {
	return string(u)
}

// Email 邮箱值对象
type Email string

// 邮箱正则表达式
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail 创建邮箱值对象
func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", errors.New("邮箱不能为空")
	}

	if !emailRegex.MatchString(email) {
		return "", errors.New("邮箱格式不正确")
	}

	return Email(email), nil
}

// Value 获取邮箱值
func (e Email) Value() string {
	return string(e)
}

// String 字符串表示
func (e Email) String() string {
	return string(e)
}

// Phone 手机号值对象
type Phone string

// 手机号正则表达式（简化版，实际可能需要更复杂的规则）
var phoneRegex = regexp.MustCompile(`^\d{11}$`)

// NewPhone 创建手机号值对象
func NewPhone(phone string) (Phone, error) {
	// 手机号可以为空
	if phone == "" {
		return "", nil
	}

	// 去除可能的空格和连字符
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")

	if !phoneRegex.MatchString(cleanPhone) {
		return "", errors.New("手机号格式不正确")
	}

	return Phone(cleanPhone), nil
}

// Value 获取手机号值
func (p Phone) Value() string {
	return string(p)
}

// String 字符串表示
func (p Phone) String() string {
	return string(p)
}

// UserStatus 用户状态值对象
type UserStatus int32

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 活跃
	UserStatusLocked   UserStatus = 2 // 锁定
	UserStatusDeleted  UserStatus = 3 // 已删除
)

// NewUserStatus 创建用户状态值对象
func NewUserStatus(status int32) (UserStatus, error) {
	switch status {
	case int32(UserStatusInactive), int32(UserStatusActive), int32(UserStatusLocked), int32(UserStatusDeleted):
		return UserStatus(status), nil
	default:
		return 0, errors.New("无效的用户状态")
	}
}

// Value 获取状态值
func (s UserStatus) Value() int32 {
	return int32(s)
}

// String 字符串表示
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "未激活"
	case UserStatusActive:
		return "活跃"
	case UserStatusLocked:
		return "锁定"
	case UserStatusDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}
