package recommendation

import (
	"time"

	"wz-backend-go/models/common"
)

// RecommendationStrategy 推荐策略
type RecommendationStrategy struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name           string `json:"name" db:"name" gorm:"not null;comment:策略名称"`
	Code           string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:策略编码"`
	Type           string `json:"type" db:"type" gorm:"not null;comment:策略类型:content,product,user,mixed"`
	Description    string `json:"description" db:"description" gorm:"comment:描述"`
	Algorithm      string `json:"algorithm" db:"algorithm" gorm:"not null;comment:算法:collaborative,content-based,hybrid,rule-based"`
	Configuration  string `json:"configuration" db:"configuration" gorm:"type:json;comment:配置信息,JSON格式"`
	ScoreFormula   string `json:"scoreFormula" db:"score_formula" gorm:"comment:评分公式"`
	FilterRules    string `json:"filterRules" db:"filter_rules" gorm:"type:json;comment:过滤规则,JSON格式"`
	SortFields     string `json:"sortFields" db:"sort_fields" gorm:"type:json;comment:排序字段,JSON格式"`
	Status         int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreatedBy      int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy      int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID       int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// RecommendationScene 推荐场景
type RecommendationScene struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name          string `json:"name" db:"name" gorm:"not null;comment:场景名称"`
	Code          string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:场景编码"`
	Type          string `json:"type" db:"type" gorm:"not null;comment:场景类型:homepage,detail,list,search,category"`
	Description   string `json:"description" db:"description" gorm:"comment:描述"`
	StrategyID    int64  `json:"strategyId" db:"strategy_id" gorm:"index;not null;comment:策略ID"`
	Position      string `json:"position" db:"position" gorm:"comment:位置:top,bottom,sidebar,popup"`
	DisplayCount  int    `json:"displayCount" db:"display_count" gorm:"default:10;comment:展示数量"`
	ClientType    string `json:"clientType" db:"client_type" gorm:"comment:客户端类型:web,mobile,miniapp,all"`
	TargetUserType string `json:"targetUserType" db:"target_user_type" gorm:"default:all;comment:目标用户类型:all,new,returning,vip"`
	Priority      int    `json:"priority" db:"priority" gorm:"default:10;comment:优先级:1-100,数字越小优先级越高"`
	TimeRules     string `json:"timeRules" db:"time_rules" gorm:"type:json;comment:时间规则,JSON格式"`
	Status        int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreatedBy     int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy     int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID      int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserInterest 用户兴趣
type UserInterest struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID        int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_interest;index;not null;comment:用户ID"`
	InterestType  string    `json:"interestType" db:"interest_type" gorm:"uniqueIndex:idx_user_interest;not null;comment:兴趣类型:category,tag,topic,keyword"`
	InterestValue string    `json:"interestValue" db:"interest_value" gorm:"uniqueIndex:idx_user_interest;not null;comment:兴趣值"`
	Score         float64   `json:"score" db:"score" gorm:"default:0;comment:兴趣分数(0-100)"`
	Source        string    `json:"source" db:"source" gorm:"comment:来源:browse,search,click,like,purchase"`
	LastActive    time.Time `json:"lastActive" db:"last_active" gorm:"not null;comment:最后活跃时间"`
	ActiveCount   int       `json:"activeCount" db:"active_count" gorm:"default:1;comment:活跃次数"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ItemFeature 物品特征
type ItemFeature struct {
	common.BaseIDModel
	common.BaseTimeModel
	ItemID        int64  `json:"itemId" db:"item_id" gorm:"uniqueIndex:idx_item_feature;index;not null;comment:物品ID"`
	ItemType      string `json:"itemType" db:"item_type" gorm:"uniqueIndex:idx_item_feature;not null;comment:物品类型:content,product,video"`
	FeatureType   string `json:"featureType" db:"feature_type" gorm:"uniqueIndex:idx_item_feature;not null;comment:特征类型:category,tag,keyword,vector"`
	FeatureValue  string `json:"featureValue" db:"feature_value" gorm:"uniqueIndex:idx_item_feature;not null;comment:特征值"`
	Weight        float64 `json:"weight" db:"weight" gorm:"default:1;comment:权重"`
	IsAutomated   bool   `json:"isAutomated" db:"is_automated" gorm:"default:false;comment:是否自动生成"`
	TenantID      int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ItemSimilarity 物品相似度
type ItemSimilarity struct {
	common.BaseIDModel
	common.BaseTimeModel
	ItemID       int64   `json:"itemId" db:"item_id" gorm:"uniqueIndex:idx_item_similar;index;not null;comment:物品ID"`
	ItemType     string  `json:"itemType" db:"item_type" gorm:"uniqueIndex:idx_item_similar;not null;comment:物品类型:content,product,video"`
	SimilarID    int64   `json:"similarId" db:"similar_id" gorm:"uniqueIndex:idx_item_similar;index;not null;comment:相似物品ID"`
	SimilarType  string  `json:"similarType" db:"similar_type" gorm:"uniqueIndex:idx_item_similar;not null;comment:相似物品类型"`
	Score        float64 `json:"score" db:"score" gorm:"not null;comment:相似度分数(0-1)"`
	FeatureMatch string  `json:"featureMatch" db:"feature_match" gorm:"type:json;comment:特征匹配详情,JSON格式"`
	UpdateTime   time.Time `json:"updateTime" db:"update_time" gorm:"not null;comment:更新时间"`
	TenantID     int64   `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// RecommendationRecord 推荐记录
type RecommendationRecord struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;comment:用户ID"`
	SessionID    string    `json:"sessionId" db:"session_id" gorm:"index;comment:会话ID"`
	SceneID      int64     `json:"sceneId" db:"scene_id" gorm:"index;not null;comment:场景ID"`
	SceneCode    string    `json:"sceneCode" db:"scene_code" gorm:"not null;comment:场景编码"`
	StrategyID   int64     `json:"strategyId" db:"strategy_id" gorm:"index;not null;comment:策略ID"`
	StrategyCode string    `json:"strategyCode" db:"strategy_code" gorm:"not null;comment:策略编码"`
	ItemID       int64     `json:"itemId" db:"item_id" gorm:"index;not null;comment:物品ID"`
	ItemType     string    `json:"itemType" db:"item_type" gorm:"not null;comment:物品类型:content,product,video"`
	Score        float64   `json:"score" db:"score" gorm:"not null;comment:推荐分数"`
	Position     int       `json:"position" db:"position" gorm:"not null;comment:推荐位置"`
	IsClicked    bool      `json:"isClicked" db:"is_clicked" gorm:"default:false;comment:是否点击"`
	ClickTime    *time.Time `json:"clickTime" db:"click_time" gorm:"comment:点击时间"`
	Reason       string    `json:"reason" db:"reason" gorm:"comment:推荐原因"`
	ExposeTime   time.Time `json:"exposeTime" db:"expose_time" gorm:"not null;comment:曝光时间"`
	DeviceInfo   string    `json:"deviceInfo" db:"device_info" gorm:"type:json;comment:设备信息,JSON格式"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// UserFeedback 用户反馈
type UserFeedback struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	RecordID     int64     `json:"recordId" db:"record_id" gorm:"index;comment:推荐记录ID"`
	ItemID       int64     `json:"itemId" db:"item_id" gorm:"index;not null;comment:物品ID"`
	ItemType     string    `json:"itemType" db:"item_type" gorm:"not null;comment:物品类型:content,product,video"`
	FeedbackType string    `json:"feedbackType" db:"feedback_type" gorm:"not null;comment:反馈类型:like,dislike,hide,report,not_interested"`
	Reason       string    `json:"reason" db:"reason" gorm:"comment:原因"`
	Detail       string    `json:"detail" db:"detail" gorm:"comment:详细说明"`
	FeedbackTime time.Time `json:"feedbackTime" db:"feedback_time" gorm:"not null;comment:反馈时间"`
	DeviceInfo   string    `json:"deviceInfo" db:"device_info" gorm:"type:json;comment:设备信息,JSON格式"`
	IP           string    `json:"ip" db:"ip" gorm:"comment:IP地址"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ABTestExperiment AB测试实验
type ABTestExperiment struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name           string    `json:"name" db:"name" gorm:"not null;comment:实验名称"`
	Code           string    `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:实验编码"`
	Description    string    `json:"description" db:"description" gorm:"comment:描述"`
	Type           string    `json:"type" db:"type" gorm:"not null;comment:实验类型:strategy,scene,feature"`
	StartTime      time.Time `json:"startTime" db:"start_time" gorm:"not null;comment:开始时间"`
	EndTime        time.Time `json:"endTime" db:"end_time" gorm:"not null;comment:结束时间"`
	SampleRate     float64   `json:"sampleRate" db:"sample_rate" gorm:"default:1;comment:采样比例(0-1)"`
	TargetUsers    string    `json:"targetUsers" db:"target_users" gorm:"type:json;comment:目标用户,JSON格式"`
	TrafficSplit   string    `json:"trafficSplit" db:"traffic_split" gorm:"type:json;not null;comment:流量分配,JSON格式"`
	MetricKeys     string    `json:"metricKeys" db:"metric_keys" gorm:"type:json;not null;comment:评估指标,JSON格式"`
	Status         string    `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-未开始,running-进行中,finished-已结束,canceled-已取消"`
	Result         string    `json:"result" db:"result" gorm:"type:json;comment:实验结果,JSON格式"`
	WinnerVariant  string    `json:"winnerVariant" db:"winner_variant" gorm:"comment:获胜变体"`
	CreatedBy      int64     `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy      int64     `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID       int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ABTestVariant AB测试变体
type ABTestVariant struct {
	common.BaseIDModel
	common.BaseTimeModel
	ExperimentID  int64  `json:"experimentId" db:"experiment_id" gorm:"index;not null;comment:实验ID"`
	VariantName   string `json:"variantName" db:"variant_name" gorm:"not null;comment:变体名称"`
	VariantCode   string `json:"variantCode" db:"variant_code" gorm:"not null;comment:变体编码"`
	Description   string `json:"description" db:"description" gorm:"comment:描述"`
	TrafficRatio  float64 `json:"trafficRatio" db:"traffic_ratio" gorm:"not null;comment:流量比例(0-1)"`
	Configuration string `json:"configuration" db:"configuration" gorm:"type:json;not null;comment:配置信息,JSON格式"`
	IsControl     bool   `json:"isControl" db:"is_control" gorm:"default:false;comment:是否对照组"`
	TenantID      int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
