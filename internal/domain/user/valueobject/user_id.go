package valueobject

import (
	"errors"
	"fmt"
)

// UserID 用户ID值对象
type UserID int64

// NewUserID 创建用户ID值对象
func NewUserID(id int64) UserID {
	return UserID(id)
}

// MustNewUserID 创建用户ID值对象，如果无效则panic
func MustNewUserID(id int64) UserID {
	if id <= 0 {
		panic("用户ID必须大于0")
	}
	return UserID(id)
}

// ValidateUserID 验证用户ID
func ValidateUserID(id int64) error {
	if id <= 0 {
		return errors.New("用户ID必须大于0")
	}
	return nil
}

// Value 获取ID值
func (id UserID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id UserID) String() string {
	return fmt.Sprintf("%d", id)
}
