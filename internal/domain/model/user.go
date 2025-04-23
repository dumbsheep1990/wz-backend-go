package model

import (
	"time"
)

// UserRole 用户角色
type UserRole string

const (
	RolePlatformAdmin UserRole = "platform_admin" // 平台管理员
	RoleTenantAdmin  UserRole = "tenant_admin"   // 租户管理员
	RoleTenantUser   UserRole = "tenant_user"    // 租户普通用户
	RolePersonalUser UserRole = "personal_user"  // 个人用户
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
	DefaultTenantID   int64     `db:"default_tenant_id" json:"default_tenant_id"` // 默认租户ID
	Role              UserRole  `db:"role" json:"role"`                          // 用户角色
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

// CompanyType 企业类型
type CompanyType int32

const (
	CompanyTypeEnterprise    CompanyType = 1 // 企业
	CompanyTypeGroup         CompanyType = 2 // 集团
	CompanyTypeGovernment    CompanyType = 3 // 政府机构/NGO/协会
	CompanyTypeResearchInst  CompanyType = 4 // 科研所
)

// CompanyVerification 企业认证请求
type CompanyVerification struct {
	ID            int64       `db:"id" json:"id"`
	UserID        int64       `db:"user_id" json:"user_id"`
	CompanyType   CompanyType `db:"company_type" json:"company_type"`
	CompanyName   string      `db:"company_name" json:"company_name"`
	// 通用字段
	BusinessLicense string     `db:"business_license" json:"business_license"`
	CommitteeLetter string     `db:"committee_letter" json:"committee_letter"`
	// 企业特定字段
	OrgCodeCert string         `db:"org_code_cert" json:"org_code_cert"`
	AgencyCert string          `db:"agency_cert" json:"agency_cert"`
	// 集团特定字段
	OrgStructure string        `db:"org_structure" json:"org_structure"`
	// 政府机构特定字段
	UnifiedSocialCredit string `db:"unified_social_credit" json:"unified_social_credit"`
	// 上市公司特定字段
	ListingCert string         `db:"listing_cert" json:"listing_cert"`
	// 其他通用字段
	ContactPerson string       `db:"contact_person" json:"contact_person"`
	ContactPhone string        `db:"contact_phone" json:"contact_phone"`
	UploadedDocument string    `db:"uploaded_document" json:"uploaded_document"`
	Status        int32        `db:"status" json:"status"`
	Remark        string       `db:"remark" json:"remark"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
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
	TenantID int64  `json:"tenant_id,omitempty"` // 可选的租户ID，如果是租户用户登录需要指定
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

// EnterpriseRegistration 企业入驻信息
type EnterpriseRegistration struct {
	ID               int64       `db:"id" json:"id"`
	UserID           int64       `db:"user_id" json:"user_id"`
	CompanyName      string      `db:"company_name" json:"company_name"`
	CompanyType      CompanyType `db:"company_type" json:"company_type"`
	ContactPerson    string      `db:"contact_person" json:"contact_person"`
	JobPosition      string      `db:"job_position" json:"job_position"`
	Region           string      `db:"region" json:"region"`
	VerificationMethod string    `db:"verification_method" json:"verification_method"`
	DetailedAddress  string      `db:"detailed_address" json:"detailed_address"`
	LocationLatitude float64     `db:"location_latitude" json:"location_latitude"`
	LocationLongitude float64    `db:"location_longitude" json:"location_longitude"`
	Status           int32       `db:"status" json:"status"`
	Remark           string      `db:"remark" json:"remark"`
	CreatedAt        time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at" json:"updated_at"`
}

// EnterpriseRegistrationRequest 企业入驻请求
type EnterpriseRegistrationRequest struct {
	UserID           int64       `json:"-"`
	CompanyName      string      `json:"company_name" validate:"required"`
	CompanyType      CompanyType `json:"company_type" validate:"required"`
	ContactPerson    string      `json:"contact_person" validate:"required"`
	JobPosition      string      `json:"job_position" validate:"required"`
	Region           string      `json:"region" validate:"required"`
	VerificationMethod string    `json:"verification_method" validate:"required"`
	DetailedAddress  string      `json:"detailed_address" validate:"required"`
	LocationLatitude float64     `json:"location_latitude"`
	LocationLongitude float64    `json:"location_longitude"`
	// 租户相关信息，用于同时创建租户
	Subdomain       string      `json:"subdomain" validate:"required,alphanum"`
	TenantType      TenantType  `json:"tenant_type" validate:"required"`
	TenantName      string      `json:"tenant_name" validate:"required"`
	TenantDesc      string      `json:"tenant_description"`
}

// VerifyCompanyRequest 验证企业请求
type VerifyCompanyRequest struct {
	UserID          int64       `json:"-"`
	CompanyType     CompanyType `json:"company_type" validate:"required"`
	CompanyName     string      `json:"company_name" validate:"required"`
	// 通用字段
	BusinessLicense string      `json:"business_license"`
	CommitteeLetter string      `json:"committee_letter"`
	// 企业特定字段
	OrgCodeCert     string      `json:"org_code_cert"`
	AgencyCert      string      `json:"agency_cert"`
	// 集团特定字段
	OrgStructure    string      `json:"org_structure"`
	// 政府机构特定字段
	UnifiedSocialCredit string  `json:"unified_social_credit"`
	// 科研所特定字段
	ResearchCert     string      `json:"research_cert"`
	// 其他通用字段
	ContactPerson   string      `json:"contact_person" validate:"required"`
	ContactPhone    string      `json:"contact_phone" validate:"required"`
	UploadedDocument string     `json:"uploaded_document"`
}
