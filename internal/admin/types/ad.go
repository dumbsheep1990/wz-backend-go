package types

// 广告列表请求
type AdListReq struct {
	Page      int32  `form:"page,default=1"`
	PageSize  int32  `form:"pageSize,default=10"`
	Title     string `form:"title,optional"`
	SpaceId   int64  `form:"spaceId,optional"`
	Status    int32  `form:"status,optional"`
	StartTime string `form:"startTime,optional"`
	EndTime   string `form:"endTime,optional"`
}

// 广告详情
type AdDetail struct {
	ID         int64  `json:"id"`
	SpaceId    int64  `json:"space_id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	Content    string `json:"content,omitempty"`
	ImageUrl   string `json:"image_url,omitempty"`
	LinkUrl    string `json:"link_url,omitempty"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Weight     int32  `json:"weight"`
	Status     int32  `json:"status"`
	ViewCount  int64  `json:"view_count"`
	ClickCount int64  `json:"click_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 广告列表响应
type AdListResp struct {
	Total int64      `json:"total"`
	List  []AdDetail `json:"list"`
}

// 广告详情请求
type AdDetailReq struct {
	ID int64 `path:"id"`
}

// 创建广告请求
type CreateAdReq struct {
	SpaceId   int64  `json:"space_id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Type      string `json:"type" validate:"required"`
	Content   string `json:"content,optional"`
	ImageUrl  string `json:"image_url,optional"`
	LinkUrl   string `json:"link_url,optional"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
	Weight    int32  `json:"weight,optional"`
	Status    int32  `json:"status,optional"`
}

// 更新广告请求
type UpdateAdReq struct {
	ID        int64  `path:"id"`
	SpaceId   int64  `json:"space_id,optional"`
	Title     string `json:"title,optional"`
	Type      string `json:"type,optional"`
	Content   string `json:"content,optional"`
	ImageUrl  string `json:"image_url,optional"`
	LinkUrl   string `json:"link_url,optional"`
	StartTime string `json:"start_time,optional"`
	EndTime   string `json:"end_time,optional"`
	Weight    int32  `json:"weight,optional"`
	Status    int32  `json:"status,optional"`
}

// 删除广告请求
type DeleteAdReq struct {
	ID int64 `path:"id"`
}

// 广告统计请求
type AdStatsReq struct {
	AdId      int64  `form:"adId,optional"`
	SpaceId   int64  `form:"spaceId,optional"`
	StartDate string `form:"startDate" validate:"required"`
	EndDate   string `form:"endDate" validate:"required"`
}

// 统计数据项
type AdStatItem struct {
	Date        string  `json:"date"`
	AdId        int64   `json:"ad_id"`
	SpaceId     int64   `json:"space_id"`
	Impressions int64   `json:"impressions"` // 展示次数
	Clicks      int64   `json:"clicks"`      // 点击次数
	CTR         float64 `json:"ctr"`         // 点击率
}

// 汇总统计数据
type AdStatTotal struct {
	Impressions int64   `json:"impressions"` // 总展示次数
	Clicks      int64   `json:"clicks"`      // 总点击次数
	CTR         float64 `json:"ctr"`         // 平均点击率
}

// 广告统计响应
type AdStatsResp struct {
	Total     AdStatTotal  `json:"total"`
	DailyData []AdStatItem `json:"daily_data"`
}
