package types

// 导航列表请求
type NavigationListReq struct {
	Page     int32  `form:"page,default=1"`
	PageSize int32  `form:"pageSize,default=10"`
	Name     string `form:"name,optional"`
	Type     string `form:"type,optional"`
	Status   int32  `form:"status,optional"`
	ParentId int64  `form:"parentId,optional"`
}

// 导航详情
type NavigationDetail struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Url        string `json:"url,omitempty"`
	Target     string `json:"target,omitempty"`
	IconUrl    string `json:"icon_url,omitempty"`
	ParentId   int64  `json:"parent_id"`
	SectionId  int64  `json:"section_id,omitempty"`
	CategoryId int64  `json:"category_id,omitempty"`
	SortOrder  int32  `json:"sort_order"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 导航列表响应
type NavigationListResp struct {
	Total int64              `json:"total"`
	List  []NavigationDetail `json:"list"`
}

// 导航详情请求
type NavigationDetailReq struct {
	ID int64 `path:"id"`
}

// 创建导航请求
type CreateNavigationReq struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Url        string `json:"url,optional"`
	Target     string `json:"target,optional"`
	IconUrl    string `json:"icon_url,optional"`
	ParentId   int64  `json:"parent_id,optional"`
	SectionId  int64  `json:"section_id,optional"`
	CategoryId int64  `json:"category_id,optional"`
	SortOrder  int32  `json:"sort_order,optional"`
	Status     int32  `json:"status,optional"`
}

// 更新导航请求
type UpdateNavigationReq struct {
	ID         int64  `path:"id"`
	Name       string `json:"name,optional"`
	Type       string `json:"type,optional"`
	Url        string `json:"url,optional"`
	Target     string `json:"target,optional"`
	IconUrl    string `json:"icon_url,optional"`
	ParentId   int64  `json:"parent_id,optional"`
	SectionId  int64  `json:"section_id,optional"`
	CategoryId int64  `json:"category_id,optional"`
	SortOrder  int32  `json:"sort_order,optional"`
	Status     int32  `json:"status,optional"`
}

// 删除导航请求
type DeleteNavigationReq struct {
	ID int64 `path:"id"`
}
