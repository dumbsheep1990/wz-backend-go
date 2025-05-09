package types

// 推荐规则列表请求
type RecommendRuleListReq struct {
	Type   string `form:"type,optional"`
	Status int32  `form:"status,optional"`
}

// 推荐规则详情
type RecommendRuleDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Params      string `json:"params,omitempty"`
	Status      int32  `json:"status"`
	SortOrder   int32  `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 推荐规则列表响应
type RecommendRuleListResp struct {
	Total int64                 `json:"total"`
	List  []RecommendRuleDetail `json:"list"`
}

// 创建/更新推荐规则请求
type SaveRecommendRuleReq struct {
	ID          int64  `json:"id,optional"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,optional"`
	Type        string `json:"type" validate:"required"`
	Params      string `json:"params,optional"`
	Status      int32  `json:"status,optional"`
	SortOrder   int32  `json:"sort_order,optional"`
}

// 内容权重设置请求
type SetContentWeightReq struct {
	ContentID int64   `json:"content_id" validate:"required"`
	Weight    float64 `json:"weight" validate:"required"`
}

// 热门内容请求
type HotContentReq struct {
	SectionID  int64  `form:"sectionId,optional"`
	CategoryID int64  `form:"categoryId,optional"`
	Type       string `form:"type,optional"`
	Limit      int32  `form:"limit,optional,default=10"`
	Period     string `form:"period,optional,default=day"`
	Page       int32  `form:"page,default=1"`
	PageSize   int32  `form:"pageSize,default=10"`
}

// 内容简要信息
type ContentBrief struct {
	ID           int64    `json:"id"`
	Type         string   `json:"type"`
	Title        string   `json:"title,omitempty"`
	Summary      string   `json:"summary,omitempty"`
	CoverUrl     string   `json:"cover_url,omitempty"`
	UserID       int64    `json:"user_id"`
	CategoryID   int64    `json:"category_id,omitempty"`
	ViewCount    int64    `json:"view_count,omitempty"`
	LikeCount    int64    `json:"like_count,omitempty"`
	CommentCount int64    `json:"comment_count,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CreatedAt    string   `json:"created_at"`
	Weight       float64  `json:"weight,omitempty"`
	HotScore     float64  `json:"hot_score,omitempty"`
}

// 热门内容响应
type HotContentResp struct {
	Total int64          `json:"total"`
	List  []ContentBrief `json:"list"`
}

// 热门关键词请求
type HotKeywordsReq struct {
	Limit  int32  `form:"limit,optional,default=10"`
	Period string `form:"period,optional,default=day"`
}

// 关键词项
type KeywordItem struct {
	Keyword     string `json:"keyword"`
	SearchCount int64  `json:"search_count"`
	Trend       int32  `json:"trend"`
}

// 热门关键词响应
type HotKeywordsResp struct {
	Keywords []KeywordItem `json:"keywords"`
}

// 热门分类请求
type HotCategoriesReq struct {
	SectionID int64  `form:"sectionId,optional"`
	Limit     int32  `form:"limit,optional,default=10"`
	Period    string `form:"period,optional,default=day"`
}

// 分类项
type CategoryItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ViewCount    int64  `json:"view_count"`
	ContentCount int64  `json:"content_count"`
	Trend        int32  `json:"trend"`
}

// 热门分类响应
type HotCategoriesResp struct {
	Categories []CategoryItem `json:"categories"`
}
