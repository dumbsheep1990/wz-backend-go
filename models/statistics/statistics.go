package statistics

import (
	"time"

	"wz-backend-go/models/common"
)

// SiteVisit 站点访问统计
type SiteVisit struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	VisitDate   time.Time `json:"visitDate" db:"visit_date" gorm:"index;not null;comment:访问日期"`
	PageViews   int       `json:"pageViews" db:"page_views" gorm:"default:0;comment:页面浏览量"`
	UniqueUsers int       `json:"uniqueUsers" db:"unique_users" gorm:"default:0;comment:独立访客数"`
	Sessions    int       `json:"sessions" db:"sessions" gorm:"default:0;comment:会话数"`
	AvgDuration float64   `json:"avgDuration" db:"avg_duration" gorm:"default:0;comment:平均访问时长(秒)"`
	BounceRate  float64   `json:"bounceRate" db:"bounce_rate" gorm:"default:0;comment:跳出率(%)"`
	NewUsers    int       `json:"newUsers" db:"new_users" gorm:"default:0;comment:新用户数"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// PageVisit 页面访问统计
type PageVisit struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	PageID      string    `json:"pageId" db:"page_id" gorm:"index;not null;comment:页面ID"`
	PagePath    string    `json:"pagePath" db:"page_path" gorm:"not null;comment:页面路径"`
	VisitDate   time.Time `json:"visitDate" db:"visit_date" gorm:"index;not null;comment:访问日期"`
	PageViews   int       `json:"pageViews" db:"page_views" gorm:"default:0;comment:页面浏览量"`
	UniqueUsers int       `json:"uniqueUsers" db:"unique_users" gorm:"default:0;comment:独立访客数"`
	AvgDuration float64   `json:"avgDuration" db:"avg_duration" gorm:"default:0;comment:平均访问时长(秒)"`
	BounceRate  float64   `json:"bounceRate" db:"bounce_rate" gorm:"default:0;comment:跳出率(%)"`
	EntryRate   float64   `json:"entryRate" db:"entry_rate" gorm:"default:0;comment:入口页率(%)"`
	ExitRate    float64   `json:"exitRate" db:"exit_rate" gorm:"default:0;comment:退出页率(%)"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// TrafficSource 流量来源统计
type TrafficSource struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	VisitDate   time.Time `json:"visitDate" db:"visit_date" gorm:"index;not null;comment:访问日期"`
	SourceType  string    `json:"sourceType" db:"source_type" gorm:"index;not null;comment:来源类型:direct,search,referral,social,email,ad"`
	SourceName  string    `json:"sourceName" db:"source_name" gorm:"index;comment:来源名称"`
	Medium      string    `json:"medium" db:"medium" gorm:"comment:媒介"`
	Campaign    string    `json:"campaign" db:"campaign" gorm:"comment:活动"`
	Keyword     string    `json:"keyword" db:"keyword" gorm:"comment:关键词"`
	PageViews   int       `json:"pageViews" db:"page_views" gorm:"default:0;comment:页面浏览量"`
	UniqueUsers int       `json:"uniqueUsers" db:"unique_users" gorm:"default:0;comment:独立访客数"`
	Sessions    int       `json:"sessions" db:"sessions" gorm:"default:0;comment:会话数"`
	BounceRate  float64   `json:"bounceRate" db:"bounce_rate" gorm:"default:0;comment:跳出率(%)"`
	AvgDuration float64   `json:"avgDuration" db:"avg_duration" gorm:"default:0;comment:平均访问时长(秒)"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserDevice 用户设备统计
type UserDevice struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	VisitDate   time.Time `json:"visitDate" db:"visit_date" gorm:"index;not null;comment:访问日期"`
	DeviceType  string    `json:"deviceType" db:"device_type" gorm:"index;not null;comment:设备类型:desktop,mobile,tablet,other"`
	Browser     string    `json:"browser" db:"browser" gorm:"index;comment:浏览器"`
	OS          string    `json:"os" db:"os" gorm:"index;comment:操作系统"`
	ScreenSize  string    `json:"screenSize" db:"screen_size" gorm:"comment:屏幕尺寸"`
	PageViews   int       `json:"pageViews" db:"page_views" gorm:"default:0;comment:页面浏览量"`
	UniqueUsers int       `json:"uniqueUsers" db:"unique_users" gorm:"default:0;comment:独立访客数"`
	Sessions    int       `json:"sessions" db:"sessions" gorm:"default:0;comment:会话数"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserGeo 用户地理位置统计
type UserGeo struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	VisitDate   time.Time `json:"visitDate" db:"visit_date" gorm:"index;not null;comment:访问日期"`
	Country     string    `json:"country" db:"country" gorm:"index;comment:国家"`
	Province    string    `json:"province" db:"province" gorm:"index;comment:省份"`
	City        string    `json:"city" db:"city" gorm:"index;comment:城市"`
	PageViews   int       `json:"pageViews" db:"page_views" gorm:"default:0;comment:页面浏览量"`
	UniqueUsers int       `json:"uniqueUsers" db:"unique_users" gorm:"default:0;comment:独立访客数"`
	Sessions    int       `json:"sessions" db:"sessions" gorm:"default:0;comment:会话数"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ContentPerformance 内容表现统计
type ContentPerformance struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID          string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	ContentID       int64     `json:"contentId" db:"content_id" gorm:"index;not null;comment:内容ID"`
	ContentType     string    `json:"contentType" db:"content_type" gorm:"index;not null;comment:内容类型:article,post,product"`
	Title           string    `json:"title" db:"title" gorm:"not null;comment:标题"`
	StatDate        time.Time `json:"statDate" db:"stat_date" gorm:"index;not null;comment:统计日期"`
	Views           int       `json:"views" db:"views" gorm:"default:0;comment:浏览量"`
	UniqueViews     int       `json:"uniqueViews" db:"unique_views" gorm:"default:0;comment:独立浏览量"`
	Likes           int       `json:"likes" db:"likes" gorm:"default:0;comment:点赞数"`
	Comments        int       `json:"comments" db:"comments" gorm:"default:0;comment:评论数"`
	Shares          int       `json:"shares" db:"shares" gorm:"default:0;comment:分享数"`
	AvgReadTime     float64   `json:"avgReadTime" db:"avg_read_time" gorm:"default:0;comment:平均阅读时间(秒)"`
	ReadCompletions int       `json:"readCompletions" db:"read_completions" gorm:"default:0;comment:阅读完成数"`
	TenantID        int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// EventTracking 事件跟踪
type EventTracking struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID      string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	EventDate   time.Time `json:"eventDate" db:"event_date" gorm:"index;not null;comment:事件日期"`
	EventName   string    `json:"eventName" db:"event_name" gorm:"index;not null;comment:事件名称"`
	EventCategory string  `json:"eventCategory" db:"event_category" gorm:"index;comment:事件分类"`
	EventAction string    `json:"eventAction" db:"event_action" gorm:"index;comment:事件动作"`
	EventLabel  string    `json:"eventLabel" db:"event_label" gorm:"comment:事件标签"`
	EventValue  float64   `json:"eventValue" db:"event_value" gorm:"comment:事件值"`
	TotalEvents int       `json:"totalEvents" db:"total_events" gorm:"default:0;comment:事件总数"`
	UniqueEvents int      `json:"uniqueEvents" db:"unique_events" gorm:"default:0;comment:独立事件数"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ConversionFunnel 转化漏斗
type ConversionFunnel struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID         string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	FunnelName     string    `json:"funnelName" db:"funnel_name" gorm:"not null;comment:漏斗名称"`
	FunnelSteps    string    `json:"funnelSteps" db:"funnel_steps" gorm:"type:json;not null;comment:漏斗步骤,JSON格式"`
	StatDate       time.Time `json:"statDate" db:"stat_date" gorm:"index;not null;comment:统计日期"`
	TotalStarts    int       `json:"totalStarts" db:"total_starts" gorm:"default:0;comment:起始步骤人数"`
	TotalCompletes int       `json:"totalCompletes" db:"total_completes" gorm:"default:0;comment:完成人数"`
	StepData       string    `json:"stepData" db:"step_data" gorm:"type:json;comment:各步骤数据,JSON格式"`
	ConversionRate float64   `json:"conversionRate" db:"conversion_rate" gorm:"default:0;comment:转化率(%)"`
	AvgTimeToComplete float64 `json:"avgTimeToComplete" db:"avg_time_to_complete" gorm:"default:0;comment:平均完成时间(秒)"`
	TenantID       int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserBehaviorSession 用户行为会话
type UserBehaviorSession struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID       string     `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	UserID       int64      `json:"userId" db:"user_id" gorm:"index;comment:用户ID"`
	SessionID    string     `json:"sessionId" db:"session_id" gorm:"uniqueIndex;not null;comment:会话ID"`
	StartTime    time.Time  `json:"startTime" db:"start_time" gorm:"not null;comment:开始时间"`
	EndTime      *time.Time `json:"endTime" db:"end_time" gorm:"comment:结束时间"`
	Duration     int        `json:"duration" db:"duration" gorm:"default:0;comment:持续时间(秒)"`
	PageCount    int        `json:"pageCount" db:"page_count" gorm:"default:0;comment:页面浏览数"`
	DeviceType   string     `json:"deviceType" db:"device_type" gorm:"comment:设备类型"`
	Browser      string     `json:"browser" db:"browser" gorm:"comment:浏览器"`
	OS           string     `json:"os" db:"os" gorm:"comment:操作系统"`
	IP           string     `json:"ip" db:"ip" gorm:"comment:IP地址"`
	Country      string     `json:"country" db:"country" gorm:"comment:国家"`
	Province     string     `json:"province" db:"province" gorm:"comment:省份"`
	City         string     `json:"city" db:"city" gorm:"comment:城市"`
	Referrer     string     `json:"referrer" db:"referrer" gorm:"comment:来源URL"`
	SourceType   string     `json:"sourceType" db:"source_type" gorm:"comment:来源类型"`
	IsNewUser    bool       `json:"isNewUser" db:"is_new_user" gorm:"default:false;comment:是否新用户"`
	IsConverted  bool       `json:"isConverted" db:"is_converted" gorm:"default:false;comment:是否转化"`
	UserAgent    string     `json:"userAgent" db:"user_agent" gorm:"comment:用户代理"`
	TenantID     int64      `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserBehaviorEvent 用户行为事件
type UserBehaviorEvent struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID       string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;comment:用户ID"`
	SessionID    string    `json:"sessionId" db:"session_id" gorm:"index;not null;comment:会话ID"`
	EventTime    time.Time `json:"eventTime" db:"event_time" gorm:"not null;comment:事件时间"`
	EventType    string    `json:"eventType" db:"event_type" gorm:"not null;comment:事件类型:pageview,click,scroll,input,etc"`
	EventName    string    `json:"eventName" db:"event_name" gorm:"index;not null;comment:事件名称"`
	EventData    string    `json:"eventData" db:"event_data" gorm:"type:json;comment:事件数据,JSON格式"`
	PageID       string    `json:"pageId" db:"page_id" gorm:"index;comment:页面ID"`
	PagePath     string    `json:"pagePath" db:"page_path" gorm:"comment:页面路径"`
	ElementID    string    `json:"elementId" db:"element_id" gorm:"comment:元素ID"`
	ElementType  string    `json:"elementType" db:"element_type" gorm:"comment:元素类型"`
	IP           string    `json:"ip" db:"ip" gorm:"comment:IP地址"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// BusinessMetric 业务指标
type BusinessMetric struct {
	common.BaseIDModel
	common.BaseTimeModel
	SiteID          string    `json:"siteId" db:"site_id" gorm:"index;not null;comment:站点ID"`
	MetricType      string    `json:"metricType" db:"metric_type" gorm:"index;not null;comment:指标类型:revenue,orders,registrations,etc"`
	MetricName      string    `json:"metricName" db:"metric_name" gorm:"index;not null;comment:指标名称"`
	StatDate        time.Time `json:"statDate" db:"stat_date" gorm:"index;not null;comment:统计日期"`
	MetricValue     float64   `json:"metricValue" db:"metric_value" gorm:"default:0;comment:指标值"`
	CompareValue    float64   `json:"compareValue" db:"compare_value" gorm:"default:0;comment:对比值(上期)"`
	ChangeRate      float64   `json:"changeRate" db:"change_rate" gorm:"default:0;comment:变化率(%)"`
	Unit            string    `json:"unit" db:"unit" gorm:"comment:单位"`
	Dimension       string    `json:"dimension" db:"dimension" gorm:"comment:维度"`
	DimensionValue  string    `json:"dimensionValue" db:"dimension_value" gorm:"comment:维度值"`
	TenantID        int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserProfile 用户画像
type UserProfile struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID           int64     `json:"userId" db:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
	FirstVisitTime   time.Time `json:"firstVisitTime" db:"first_visit_time" gorm:"not null;comment:首次访问时间"`
	LastVisitTime    time.Time `json:"lastVisitTime" db:"last_visit_time" gorm:"not null;comment:最近访问时间"`
	VisitCount       int       `json:"visitCount" db:"visit_count" gorm:"default:0;comment:访问次数"`
	TotalTimeSpent   int       `json:"totalTimeSpent" db:"total_time_spent" gorm:"default:0;comment:总停留时间(秒)"`
	PageViewCount    int       `json:"pageViewCount" db:"page_view_count" gorm:"default:0;comment:页面浏览数"`
	DeviceTypes      string    `json:"deviceTypes" db:"device_types" gorm:"type:json;comment:设备类型,JSON格式"`
	Browsers         string    `json:"browsers" db:"browsers" gorm:"type:json;comment:浏览器,JSON格式"`
	OperatingSystems string    `json:"operatingSystems" db:"operating_systems" gorm:"type:json;comment:操作系统,JSON格式"`
	Countries        string    `json:"countries" db:"countries" gorm:"type:json;comment:国家,JSON格式"`
	Cities           string    `json:"cities" db:"cities" gorm:"type:json;comment:城市,JSON格式"`
	Interests        string    `json:"interests" db:"interests" gorm:"type:json;comment:兴趣,JSON格式"`
	SourceTypes      string    `json:"sourceTypes" db:"source_types" gorm:"type:json;comment:来源类型,JSON格式"`
	PreferredTime    string    `json:"preferredTime" db:"preferred_time" gorm:"comment:偏好时段"`
	ConversionCount  int       `json:"conversionCount" db:"conversion_count" gorm:"default:0;comment:转化次数"`
	TotalPurchase    float64   `json:"totalPurchase" db:"total_purchase" gorm:"default:0;comment:总购买金额"`
	Tags             string    `json:"tags" db:"tags" gorm:"type:json;comment:标签,JSON格式"`
	TenantID         int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
