package domain

import "time"

// StatisticsData 统计数据结构体
type StatisticsData struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	Date      time.Time `json:"date,omitempty" db:"date"`           // 日期
	Type      string    `json:"type,omitempty" db:"type"`           // 统计类型(pv, uv, content_view, user_register等)
	Value     int64     `json:"value,omitempty" db:"value"`         // 统计值
	ItemID    int64     `json:"item_id,omitempty" db:"item_id"`     // 关联项目ID
	ItemType  string    `json:"item_type,omitempty" db:"item_type"` // 关联项目类型
	TenantID  int64     `json:"tenant_id,omitempty" db:"tenant_id"` // 租户ID，多租户支持
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ContentRankingItem 内容排行项
type ContentRankingItem struct {
	ID        int64  `json:"id,omitempty" db:"id"`                 // 内容ID
	Title     string `json:"title,omitempty" db:"title"`           // 标题
	Cover     string `json:"cover,omitempty" db:"cover"`           // 封面
	Type      string `json:"type,omitempty" db:"type"`             // 类型
	ViewCount int64  `json:"view_count,omitempty" db:"view_count"` // 浏览量
	LikeCount int64  `json:"like_count,omitempty" db:"like_count"` // 点赞量
	URL       string `json:"url,omitempty" db:"url"`               // 链接
}

// OverviewData 概览数据结构体
type OverviewData struct {
	TotalUsers      int64 `json:"total_users"`       // 总用户数
	TotalContent    int64 `json:"total_content"`     // 总内容数
	TotalPV         int64 `json:"total_pv"`          // 总PV
	TotalUV         int64 `json:"total_uv"`          // 总UV
	TodayPV         int64 `json:"today_pv"`          // 今日PV
	TodayUV         int64 `json:"today_uv"`          // 今日UV
	TodayNewUsers   int64 `json:"today_new_users"`   // 今日新增用户
	TodayNewContent int64 `json:"today_new_content"` // 今日新增内容
}

// StatisticsRepository 统计仓储接口
type StatisticsRepository interface {
	RecordStatistics(data *StatisticsData) error
	GetOverview(tenantID int64) (*OverviewData, error)
	GetStatisticsByType(tenantID int64, type_ string, startDate, endDate time.Time) ([]*StatisticsData, error)
	GetContentRanking(tenantID int64, limit int) ([]*ContentRankingItem, error)
}
