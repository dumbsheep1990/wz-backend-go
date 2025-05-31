package valueobject

import (
	"errors"
	"strconv"
)

// UserID 用户ID值对象
type UserID struct {
	value int64
}

// NewUserID 创建用户ID值对象
func NewUserID(id int64) UserID {
	return UserID{value: id}
}

// NewUserIDFromString 从字符串创建用户ID值对象
func NewUserIDFromString(idStr string) (UserID, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return UserID{}, errors.New("无效的用户ID格式")
	}
	if id <= 0 {
		return UserID{}, errors.New("用户ID必须大于0")
	}
	return UserID{value: id}, nil
}

// MustNewUserID 创建用户ID值对象，如果无效则panic
func MustNewUserID(id int64) UserID {
	if id <= 0 {
		panic("用户ID必须大于0")
	}
	return UserID{value: id}
}

// Value 获取用户ID值
func (u UserID) Value() int64 {
	return u.value
}

// String 获取用户ID的字符串表示
func (u UserID) String() string {
	return strconv.FormatInt(u.value, 10)
}

// IsValid 检查用户ID是否有效
func (u UserID) IsValid() bool {
	return u.value > 0
}

// Equals 比较两个用户ID是否相等
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// ValidateUserID 验证用户ID
func ValidateUserID(id int64) error {
	if id <= 0 {
		return errors.New("用户ID必须大于0")
	}
	return nil
}
