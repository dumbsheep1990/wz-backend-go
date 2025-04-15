package model

import (
	"time"
)

// User 系统用户
type User struct {
	ID                int64     `db:"id" json:"id"`
	Username          string    `db:"username" json:"username"`
	Password          string    `db:"password" json:"-"`
	Email             string    `db:"email" json:"email"`
	Phone             string    `db:"phone" json:"phone"`
	Status            int32     `db:"status" json:"status"`
	IsVerified        bool      `db:"is_verified" json:"is_verified"`
	IsCompanyVerified bool      `db:"is_company_verified" json:"is_company_verified"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// UserDetail 用户详细信息
type UserDetail struct {
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"user_id"`
	RealName string    `db:"real_name" json:"real_name"`
	IDCard   string    `db:"id_card" json:"id_card"`
	Avatar   string    `db:"avatar" json:"avatar"`
	Gender   int32     `db:"gender" json:"gender"`
	Birthday time.Time `db:"birthday" json:"birthday"`
	Address  string    `db:"address" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CompanyVerification 企业认证请求
type CompanyVerification struct {
	ID            int64     `db:"id" json:"id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	CompanyName   string    `db:"company_name" json:"company_name"`
	BusinessLicense string  `db:"business_license" json:"business_license"`
	ContactPerson string    `db:"contact_person" json:"contact_person"`
	ContactPhone  string    `db:"contact_phone" json:"contact_phone"`
	Status        int32     `db:"status" json:"status"`
	Remark        string    `db:"remark" json:"remark"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// UserLoginLog 用户登录日志条目
type UserLoginLog struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	LoginTime   time.Time `db:"login_time" json:"login_time"`
	LoginIP     string    `db:"login_ip" json:"login_ip"`
	UserAgent   string    `db:"user_agent" json:"user_agent"`
	DeviceType  int32     `db:"device_type" json:"device_type"`
	LoginStatus int32     `db:"login_status" json:"login_status"`
}

// UserBehaviorLog 用户行为日志条目
type UserBehaviorLog struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	Action       string    `db:"action" json:"action"`
	ResourceType string    `db:"resource_type" json:"resource_type"`
	ResourceID   int64     `db:"resource_id" json:"resource_id"`
	IP           string    `db:"ip" json:"ip"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// RegisterRequest 注册新用户请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest 更新用户信息的请求
type UpdateUserRequest struct {
	UserID   int64  `json:"-"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`
}

// VerifyUserRequest 验证用户请求
type VerifyUserRequest struct {
	UserID           int64  `json:"-"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

// VerifyCompanyRequest 验证企业请求
type VerifyCompanyRequest struct {
	UserID          int64  `json:"-"`
	CompanyName     string `json:"company_name" validate:"required"`
	BusinessLicense string `json:"business_license" validate:"required"`
	ContactPerson   string `json:"contact_person" validate:"required"`
	ContactPhone    string `json:"contact_phone" validate:"required"`
}
