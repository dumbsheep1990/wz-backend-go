package types

// CompanyType 企业类型
type CompanyType int32

const (
	CompanyTypeEnterprise    CompanyType = 1 // 企业
	CompanyTypeGroup         CompanyType = 2 // 集团
	CompanyTypeGovernment    CompanyType = 3 // 政府机构/NGO/协会
	CompanyTypeResearchInst  CompanyType = 4 // 科研所
)

// EnterpriseRegistrationReq 企业注册请求
type EnterpriseRegistrationReq struct {
	CompanyName        string      `json:"company_name" validate:"required"`
	CompanyType        CompanyType `json:"company_type" validate:"required"`
	ContactPerson      string      `json:"contact_person" validate:"required"`
	JobPosition        string      `json:"job_position" validate:"required"`
	Region             string      `json:"region" validate:"required"`
	VerificationMethod string      `json:"verification_method" validate:"required"`
	DetailedAddress    string      `json:"detailed_address" validate:"required"`
	LocationLatitude   float64     `json:"location_latitude"`
	LocationLongitude  float64     `json:"location_longitude"`
	
	// 租户相关信息
	Subdomain   string `json:"subdomain" validate:"required,alphanum"`
	TenantName  string `json:"tenant_name" validate:"required"`
	TenantDesc  string `json:"tenant_description"`
}

// EnterpriseRegistrationResp 企业注册响应
type EnterpriseRegistrationResp struct {
	Success      bool   `json:"success"`
	TenantID     int64  `json:"tenant_id"`
	Subdomain    string `json:"subdomain"`
	TenantName   string `json:"tenant_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}
