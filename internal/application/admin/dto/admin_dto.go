package dto

import (
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// AdminDTO represents the data transfer object for an admin
type AdminDTO struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	RoleName   string    `json:"roleName,omitempty"`
	Status     int       `json:"status"`
	StatusText string    `json:"statusText"`
	LastLogin  time.Time `json:"lastLogin"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// AdminCreateRequest represents a request to create an admin
type AdminCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	RoleID   string `json:"roleId" binding:"required"`
}

// AdminUpdateRequest represents a request to update an admin
type AdminUpdateRequest struct {
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=32"`
	RoleID   string `json:"roleId,omitempty"`
	Status   *int   `json:"status,omitempty"`
}

// AdminPasswordChangeRequest represents a request to change an admin's password
type AdminPasswordChangeRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=32"`
}

// AdminLoginRequest represents a request to login as an admin
type AdminLoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha,omitempty"`
	CaptchaID string `json:"captchaId,omitempty"`
}

// AdminLoginResponse represents a response to an admin login request
type AdminLoginResponse struct {
	Token      string     `json:"token"`
	ExpiresAt  int64      `json:"expiresAt"`
	Admin      *AdminDTO  `json:"admin"`
	Roles      []RoleDTO  `json:"roles,omitempty"`
	Categories []string   `json:"categories"`
	SiteInfo   *SiteInfo  `json:"siteInfo"`
}

// SiteInfo contains basic information about the website
type SiteInfo struct {
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	ApiUrl       string `json:"apiUrl"`
	AdminDomain  string `json:"adminDomain"`
	SiteDomain   string `json:"siteDomain"`
	ICP          string `json:"icp"`
	PSB          string `json:"psb"`
	ServicePhone string `json:"servicePhone"`
	ReportPhone  string `json:"reportPhone"`
	ReportEmail  string `json:"reportEmail"`
	Copyright    string `json:"copyright"`
	Address      string `json:"address"`
}

// AdminListResponse represents a paginated list of admins
type AdminListResponse struct {
	Total int64      `json:"total"`
	Items []AdminDTO `json:"items"`
}

// AdminDetailResponse represents a detailed view of an admin
type AdminDetailResponse struct {
	Admin       AdminDTO   `json:"admin"`
	Roles       []RoleDTO  `json:"roles,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
}

// MapAdminToDTO maps an Admin entity to an AdminDTO
func MapAdminToDTO(admin *entity.Admin) AdminDTO {
	return AdminDTO{
		ID:         admin.ID.Value(),
		Username:   admin.Username.Value(),
		Role:       admin.Role.Value(),
		Status:     admin.Status.Value(),
		StatusText: admin.Status.String(),
		LastLogin:  admin.LastLogin,
		CreatedAt:  admin.CreatedAt,
		UpdatedAt:  admin.UpdatedAt,
	}
}

// NewDefaultSiteInfo creates a new SiteInfo with default values for 万知
func NewDefaultSiteInfo() *SiteInfo {
	return &SiteInfo{
		Name:         "万知管理系统",
		Logo:         "/logo.png",
		ApiUrl:       "/api",
		AdminDomain:  "admin.wz.qq",
		SiteDomain:   "wz.qq",
		ICP:          "苏ICP备19049108号",
		PSB:          "苏公网安备 32010402000846号",
		ServicePhone: "0515-88200000",
		ReportPhone:  "0515-88200000",
		ReportEmail:  "2895915959@qq.com",
		Copyright:    "版权所有©南京万知园信息工程有限公司 万知网",
		Address:      "江苏省盐城市通港路11号综合楼一楼1号",
	}
}

// GetDefaultCategories returns the default categories for 万知
func GetDefaultCategories() []string {
	return []string{
		"同用", "同好", "同购", "同年", "同游",
		"同在", "同市", "同企", "同亲", "同班",
		"同师", "同业", "同网", "同工", "同务",
		"同艺", "同玩", "同闲", "同拍", "同乡", "同学",
	}
}
