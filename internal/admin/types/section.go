package types

// 板块列表请求
type SectionListReq struct {
	Page     int32  `form:"page,default=1"`
	PageSize int32  `form:"pageSize,default=10"`
	Name     string `form:"name,optional"`
	Status   int32  `form:"status,optional"`
}

// 板块详情
type SectionDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	IconUrl     string `json:"icon_url,omitempty"`
	BannerUrl   string `json:"banner_url,omitempty"`
	SortOrder   int32  `json:"sort_order"`
	Status      int32  `json:"status"`
	ShowInHome  bool   `json:"show_in_home"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 板块列表响应
type SectionListResp struct {
	Total int64           `json:"total"`
	List  []SectionDetail `json:"list"`
}

// 板块详情请求
type SectionDetailReq struct {
	ID int64 `path:"id"`
}

// 创建板块请求
type CreateSectionReq struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description,optional"`
	IconUrl     string `json:"icon_url,optional"`
	BannerUrl   string `json:"banner_url,optional"`
	SortOrder   int32  `json:"sort_order,optional"`
	Status      int32  `json:"status,optional"`
	ShowInHome  bool   `json:"show_in_home,optional"`
}

// 更新板块请求
type UpdateSectionReq struct {
	ID          int64  `path:"id"`
	Name        string `json:"name,optional"`
	Code        string `json:"code,optional"`
	Description string `json:"description,optional"`
	IconUrl     string `json:"icon_url,optional"`
	BannerUrl   string `json:"banner_url,optional"`
	SortOrder   int32  `json:"sort_order,optional"`
	Status      int32  `json:"status,optional"`
	ShowInHome  bool   `json:"show_in_home,optional"`
}

// 删除板块请求
type DeleteSectionReq struct {
	ID int64 `path:"id"`
}
