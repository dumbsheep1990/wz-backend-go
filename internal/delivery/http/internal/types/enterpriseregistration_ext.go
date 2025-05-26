package types

// UpdateEnterpriseRegistrationReq 更新企业入驻请求
type UpdateEnterpriseRegistrationReq struct {
	UserID             int64       `json:"-"`
	CompanyName        string      `json:"company_name" validate:"required"`
	CompanyType        CompanyType `json:"company_type" validate:"required"`
	ContactPerson      string      `json:"contact_person" validate:"required"`
	Region             string      `json:"region" validate:"required"`
	DetailedAddress    string      `json:"detailed_address" validate:"required"`
	LocationLatitude   float64     `json:"location_latitude"`
	LocationLongitude  float64     `json:"location_longitude"`
	ServiceArea        string      `json:"service_area"`
}

// UpdateEnterpriseRegistrationResp 更新企业入驻响应
type UpdateEnterpriseRegistrationResp struct {
	Success bool `json:"success"`
}

// GetEnterpriseRegistrationResp 获取企业入驻信息响应
type GetEnterpriseRegistrationResp struct {
	ID                int64       `json:"id"`
	UserID            int64       `json:"user_id"`
	CompanyName       string      `json:"company_name"`
	CompanyType       CompanyType `json:"company_type"`
	ContactPerson     string      `json:"contact_person"`
	Region            string      `json:"region"`
	DetailedAddress   string      `json:"detailed_address"`
	LocationLatitude  float64     `json:"location_latitude"`
	LocationLongitude float64     `json:"location_longitude"`
	Status            int32       `json:"status"`
	ServiceArea       string      `json:"service_area"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
}

// VerifyEnterpriseReq 企业验证请求
type VerifyEnterpriseReq struct {
	UserID           int64  `json:"-"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

// VerifyEnterpriseResp 企业验证响应
type VerifyEnterpriseResp struct {
	Success bool `json:"success"`
}
