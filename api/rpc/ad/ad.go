package ad

import (
	"context"

	"google.golang.org/grpc"
)

// 广告服务接口定义
type AdService interface {
	// 广告位相关接口
	// 获取广告位列表
	GetAdSpaceList(ctx context.Context, in *GetAdSpaceListReq) (*GetAdSpaceListResp, error)
	// 获取广告位详情
	GetAdSpaceDetail(ctx context.Context, in *GetAdSpaceDetailReq) (*AdSpaceDetailResp, error)
	// 创建广告位
	CreateAdSpace(ctx context.Context, in *CreateAdSpaceReq) (*AdSpaceDetailResp, error)
	// 更新广告位
	UpdateAdSpace(ctx context.Context, in *UpdateAdSpaceReq) (*UpdateAdSpaceResp, error)
	// 删除广告位
	DeleteAdSpace(ctx context.Context, in *DeleteAdSpaceReq) (*DeleteAdSpaceResp, error)

	// 广告内容相关接口
	// 获取广告内容列表
	GetAdList(ctx context.Context, in *GetAdListReq) (*GetAdListResp, error)
	// 获取广告详情
	GetAdDetail(ctx context.Context, in *GetAdDetailReq) (*AdDetailResp, error)
	// 创建广告
	CreateAd(ctx context.Context, in *CreateAdReq) (*AdDetailResp, error)
	// 更新广告
	UpdateAd(ctx context.Context, in *UpdateAdReq) (*UpdateAdResp, error)
	// 删除广告
	DeleteAd(ctx context.Context, in *DeleteAdReq) (*DeleteAdResp, error)
	// 获取指定位置的广告
	GetAdsByPosition(ctx context.Context, in *GetAdsByPositionReq) (*GetAdsByPositionResp, error)

	// 广告统计相关接口
	// 记录广告展示
	RecordAdImpression(ctx context.Context, in *RecordAdImpressionReq) (*RecordAdImpressionResp, error)
	// 记录广告点击
	RecordAdClick(ctx context.Context, in *RecordAdClickReq) (*RecordAdClickResp, error)
	// 获取广告统计数据
	GetAdStats(ctx context.Context, in *GetAdStatsReq) (*GetAdStatsResp, error)
}

// 广告服务RPC客户端
type adServiceClient struct {
	conn *grpc.ClientConn
}

// 创建广告服务客户端
func NewAdService(conn *grpc.ClientConn) AdService {
	return &adServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// ====== 广告位相关 ======

// 获取广告位列表请求
type GetAdSpaceListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Name     string `json:"name,omitempty"`
	Position string `json:"position,omitempty"` // 广告位置标识
	Status   int32  `json:"status,omitempty"`   // 0-禁用，1-启用
}

// 获取广告位列表响应
type GetAdSpaceListResp struct {
	Total int64                `json:"total"`
	List  []*AdSpaceDetailResp `json:"list"`
}

// 获取广告位详情请求
type GetAdSpaceDetailReq struct {
	Id int64 `json:"id"`
}

// 广告位详情响应
type AdSpaceDetailResp struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`        // 广告位名称
	Position    string `json:"position"`    // 广告位置标识码
	Description string `json:"description"` // 广告位描述
	Width       int32  `json:"width"`       // 宽度
	Height      int32  `json:"height"`      // 高度
	Type        string `json:"type"`        // 类型：image-图片，text-文字，code-代码
	MaxAds      int32  `json:"max_ads"`     // 最大广告数量
	Status      int32  `json:"status"`      // 状态：0-禁用，1-启用
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 创建广告位请求
type CreateAdSpaceReq struct {
	Name        string `json:"name"`
	Position    string `json:"position"`
	Description string `json:"description,omitempty"`
	Width       int32  `json:"width,omitempty"`
	Height      int32  `json:"height,omitempty"`
	Type        string `json:"type"`
	MaxAds      int32  `json:"max_ads,omitempty"`
	Status      int32  `json:"status,omitempty"`
}

// 更新广告位请求
type UpdateAdSpaceReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Position    string `json:"position,omitempty"`
	Description string `json:"description,omitempty"`
	Width       int32  `json:"width,omitempty"`
	Height      int32  `json:"height,omitempty"`
	Type        string `json:"type,omitempty"`
	MaxAds      int32  `json:"max_ads,omitempty"`
	Status      int32  `json:"status,omitempty"`
}

// 更新广告位响应
type UpdateAdSpaceResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除广告位请求
type DeleteAdSpaceReq struct {
	Id int64 `json:"id"`
}

// 删除广告位响应
type DeleteAdSpaceResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ====== 广告内容相关 ======

// 获取广告列表请求
type GetAdListReq struct {
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
	Title     string `json:"title,omitempty"`
	SpaceId   int64  `json:"space_id,omitempty"`
	Status    int32  `json:"status,omitempty"`     // 0-禁用，1-启用
	StartTime string `json:"start_time,omitempty"` // 开始展示时间
	EndTime   string `json:"end_time,omitempty"`   // 结束展示时间
}

// 获取广告列表响应
type GetAdListResp struct {
	Total int64           `json:"total"`
	List  []*AdDetailResp `json:"list"`
}

// 获取广告详情请求
type GetAdDetailReq struct {
	Id int64 `json:"id"`
}

// 广告详情响应
type AdDetailResp struct {
	Id         int64  `json:"id"`
	SpaceId    int64  `json:"space_id"`            // 所属广告位ID
	Title      string `json:"title"`               // 广告标题
	Type       string `json:"type"`                // 类型：image-图片，text-文字，code-代码
	Content    string `json:"content"`             // 广告内容
	ImageUrl   string `json:"image_url,omitempty"` // 图片URL
	LinkUrl    string `json:"link_url,omitempty"`  // 链接URL
	StartTime  string `json:"start_time"`          // 开始展示时间
	EndTime    string `json:"end_time"`            // 结束展示时间
	Weight     int32  `json:"weight"`              // 权重/优先级
	Status     int32  `json:"status"`              // 状态：0-禁用，1-启用
	ViewCount  int64  `json:"view_count"`          // 展示次数
	ClickCount int64  `json:"click_count"`         // 点击次数
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 创建广告请求
type CreateAdReq struct {
	SpaceId   int64  `json:"space_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	LinkUrl   string `json:"link_url,omitempty"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Weight    int32  `json:"weight,omitempty"`
	Status    int32  `json:"status,omitempty"`
}

// 更新广告请求
type UpdateAdReq struct {
	Id        int64  `json:"id"`
	SpaceId   int64  `json:"space_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	LinkUrl   string `json:"link_url,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Weight    int32  `json:"weight,omitempty"`
	Status    int32  `json:"status,omitempty"`
}

// 更新广告响应
type UpdateAdResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除广告请求
type DeleteAdReq struct {
	Id int64 `json:"id"`
}

// 删除广告响应
type DeleteAdResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 获取指定位置的广告请求
type GetAdsByPositionReq struct {
	Position string `json:"position"`        // 广告位置标识
	Limit    int32  `json:"limit,omitempty"` // 返回广告数量
}

// 获取指定位置的广告响应
type GetAdsByPositionResp struct {
	Ads []*AdDetailResp `json:"ads"`
}

// ====== 广告统计相关 ======

// 记录广告展示请求
type RecordAdImpressionReq struct {
	AdId    int64  `json:"ad_id"`
	UserId  int64  `json:"user_id,omitempty"`
	IP      string `json:"ip,omitempty"`
	UA      string `json:"ua,omitempty"`
	Referer string `json:"referer,omitempty"`
}

// 记录广告展示响应
type RecordAdImpressionResp struct {
	Success bool `json:"success"`
}

// 记录广告点击请求
type RecordAdClickReq struct {
	AdId    int64  `json:"ad_id"`
	UserId  int64  `json:"user_id,omitempty"`
	IP      string `json:"ip,omitempty"`
	UA      string `json:"ua,omitempty"`
	Referer string `json:"referer,omitempty"`
}

// 记录广告点击响应
type RecordAdClickResp struct {
	Success bool `json:"success"`
}

// 获取广告统计数据请求
type GetAdStatsReq struct {
	AdId      int64  `json:"ad_id,omitempty"`
	SpaceId   int64  `json:"space_id,omitempty"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
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

// 获取广告统计数据响应
type GetAdStatsResp struct {
	Total     AdStatTotal   `json:"total"`
	DailyData []*AdStatItem `json:"daily_data"`
}

// 汇总统计数据
type AdStatTotal struct {
	Impressions int64   `json:"impressions"` // 总展示次数
	Clicks      int64   `json:"clicks"`      // 总点击次数
	CTR         float64 `json:"ctr"`         // 平均点击率
}

// 实现AdService接口的方法
func (c *adServiceClient) GetAdSpaceList(ctx context.Context, in *GetAdSpaceListReq) (*GetAdSpaceListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetAdSpaceListResp{
		Total: 0,
		List:  []*AdSpaceDetailResp{},
	}, nil
}

func (c *adServiceClient) GetAdSpaceDetail(ctx context.Context, in *GetAdSpaceDetailReq) (*AdSpaceDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &AdSpaceDetailResp{}, nil
}

func (c *adServiceClient) CreateAdSpace(ctx context.Context, in *CreateAdSpaceReq) (*AdSpaceDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &AdSpaceDetailResp{}, nil
}

func (c *adServiceClient) UpdateAdSpace(ctx context.Context, in *UpdateAdSpaceReq) (*UpdateAdSpaceResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateAdSpaceResp{Success: true}, nil
}

func (c *adServiceClient) DeleteAdSpace(ctx context.Context, in *DeleteAdSpaceReq) (*DeleteAdSpaceResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteAdSpaceResp{Success: true}, nil
}

func (c *adServiceClient) GetAdList(ctx context.Context, in *GetAdListReq) (*GetAdListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetAdListResp{
		Total: 0,
		List:  []*AdDetailResp{},
	}, nil
}

func (c *adServiceClient) GetAdDetail(ctx context.Context, in *GetAdDetailReq) (*AdDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &AdDetailResp{}, nil
}

func (c *adServiceClient) CreateAd(ctx context.Context, in *CreateAdReq) (*AdDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &AdDetailResp{}, nil
}

func (c *adServiceClient) UpdateAd(ctx context.Context, in *UpdateAdReq) (*UpdateAdResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateAdResp{Success: true}, nil
}

func (c *adServiceClient) DeleteAd(ctx context.Context, in *DeleteAdReq) (*DeleteAdResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteAdResp{Success: true}, nil
}

func (c *adServiceClient) GetAdsByPosition(ctx context.Context, in *GetAdsByPositionReq) (*GetAdsByPositionResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetAdsByPositionResp{
		Ads: []*AdDetailResp{},
	}, nil
}

func (c *adServiceClient) RecordAdImpression(ctx context.Context, in *RecordAdImpressionReq) (*RecordAdImpressionResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &RecordAdImpressionResp{Success: true}, nil
}

func (c *adServiceClient) RecordAdClick(ctx context.Context, in *RecordAdClickReq) (*RecordAdClickResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &RecordAdClickResp{Success: true}, nil
}

func (c *adServiceClient) GetAdStats(ctx context.Context, in *GetAdStatsReq) (*GetAdStatsResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetAdStatsResp{
		Total:     AdStatTotal{},
		DailyData: []*AdStatItem{},
	}, nil
}
