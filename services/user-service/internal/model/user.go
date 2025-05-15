package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID                int64     `db:"id"`
	Username          string    `db:"username"`
	Password          string    `db:"password"`
	Email             string    `db:"email"`
	Phone             string    `db:"phone"`
	Status            int32     `db:"status"`
	IsVerified        bool      `db:"is_verified"`
	IsCompanyVerified bool      `db:"is_company_verified"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// SetPassword 设置密码（使用bcrypt加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserBehavior 用户行为模型
type UserBehavior struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	Action       string    `db:"action"`
	ResourceType string    `db:"resource_type"`
	ResourceID   int64     `db:"resource_id"`
	CreatedAt    time.Time `db:"created_at"`
}
