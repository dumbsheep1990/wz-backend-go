package types

// 广告位列表请求
type AdSpaceListReq struct {
	Page     int32  `form:"page,default=1"`
	PageSize int32  `form:"pageSize,default=10"`
	Name     string `form:"name,optional"`
	Position string `form:"position,optional"`
	Status   int32  `form:"status,optional"`
}

// 广告位详情
type AdSpaceDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	Description string `json:"description,omitempty"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	Type        string `json:"type"`
	MaxAds      int32  `json:"max_ads"`
	Status      int32  `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 广告位列表响应
type AdSpaceListResp struct {
	Total int64           `json:"total"`
	List  []AdSpaceDetail `json:"list"`
}

// 广告位详情请求
type AdSpaceDetailReq struct {
	ID int64 `path:"id"`
}

// 创建广告位请求
type CreateAdSpaceReq struct {
	Name        string `json:"name" validate:"required"`
	Position    string `json:"position" validate:"required"`
	Description string `json:"description,optional"`
	Width       int32  `json:"width,optional"`
	Height      int32  `json:"height,optional"`
	Type        string `json:"type" validate:"required"`
	MaxAds      int32  `json:"max_ads,optional"`
	Status      int32  `json:"status,optional"`
}

// 更新广告位请求
type UpdateAdSpaceReq struct {
	ID          int64  `path:"id"`
	Name        string `json:"name,optional"`
	Position    string `json:"position,optional"`
	Description string `json:"description,optional"`
	Width       int32  `json:"width,optional"`
	Height      int32  `json:"height,optional"`
	Type        string `json:"type,optional"`
	MaxAds      int32  `json:"max_ads,optional"`
	Status      int32  `json:"status,optional"`
}

// 删除广告位请求
type DeleteAdSpaceReq struct {
	ID int64 `path:"id"`
}
