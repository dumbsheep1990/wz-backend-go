package interaction

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Like 点赞模型
type Like struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID     int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	TargetType string    `json:"targetType" db:"target_type" gorm:"not null;comment:目标类型:post/comment/reply"`
	TargetID   int64     `json:"targetId" db:"target_id" gorm:"index;not null;comment:目标ID"`
	Status     int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-有效,0-已取消"`
	TenantID   int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Favorite 收藏模型
type Favorite struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	TargetType  string    `json:"targetType" db:"target_type" gorm:"not null;comment:目标类型:post/product/article"`
	TargetID    int64     `json:"targetId" db:"target_id" gorm:"index;not null;comment:目标ID"`
	FolderID    int64     `json:"folderId" db:"folder_id" gorm:"index;comment:收藏夹ID"`
	Remark      string    `json:"remark" db:"remark" gorm:"comment:备注"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-有效,0-已取消"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// FavoriteFolder 收藏夹
type FavoriteFolder struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Name        string    `json:"name" db:"name" gorm:"not null;comment:名称"`
	Description string    `json:"description" db:"description" gorm:"comment:描述"`
	Cover       string    `json:"cover" db:"cover" gorm:"comment:封面"`
	ItemCount   int       `json:"itemCount" db:"item_count" gorm:"default:0;comment:收藏数量"`
	IsPublic    bool      `json:"isPublic" db:"is_public" gorm:"default:false;comment:是否公开"`
	SortOrder   int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Comment 评论模型
type Comment struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	TargetType  string    `json:"targetType" db:"target_type" gorm:"not null;comment:目标类型:post/product/article"`
	TargetID    int64     `json:"targetId" db:"target_id" gorm:"index;not null;comment:目标ID"`
	ParentID    int64     `json:"parentId" db:"parent_id" gorm:"index;default:0;comment:父评论ID,0表示顶级评论"`
	ReplyUserID int64     `json:"replyUserId" db:"reply_user_id" gorm:"index;comment:回复用户ID"`
	Content     string    `json:"content" db:"content" gorm:"not null;comment:评论内容"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-隐藏,-1-已删除"`
	LikeCount   int       `json:"likeCount" db:"like_count" gorm:"default:0;comment:点赞数"`
	ReplyCount  int       `json:"replyCount" db:"reply_count" gorm:"default:0;comment:回复数"`
	IP          string    `json:"ip" db:"ip" gorm:"comment:评论IP"`
	UserAgent   string    `json:"userAgent" db:"user_agent" gorm:"comment:用户代理"`
	IsAnonymous bool      `json:"isAnonymous" db:"is_anonymous" gorm:"default:false;comment:是否匿名"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Follow 关注模型
type Follow struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_target;index;not null;comment:用户ID"`
	TargetType  string    `json:"targetType" db:"target_type" gorm:"uniqueIndex:idx_user_target;not null;comment:目标类型:user/tag/category/topic"`
	TargetID    int64     `json:"targetId" db:"target_id" gorm:"uniqueIndex:idx_user_target;index;not null;comment:目标ID"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-有效,0-已取消"`
	Remark      string    `json:"remark" db:"remark" gorm:"comment:备注"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Report 举报模型
type Report struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:举报人ID"`
	TargetType  string    `json:"targetType" db:"target_type" gorm:"not null;comment:目标类型:post/comment/user/product"`
	TargetID    int64     `json:"targetId" db:"target_id" gorm:"index;not null;comment:目标ID"`
	ReasonType  string    `json:"reasonType" db:"reason_type" gorm:"not null;comment:举报类型"`
	Content     string    `json:"content" db:"content" gorm:"comment:举报内容"`
	Evidence    string    `json:"evidence" db:"evidence" gorm:"type:json;comment:证据,JSON数组"`
	Status      string    `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-待处理,processing-处理中,resolved-已解决,rejected-已驳回"`
	HandlerID   int64     `json:"handlerId" db:"handler_id" gorm:"index;comment:处理人ID"`
	HandlerNote string    `json:"handlerNote" db:"handler_note" gorm:"comment:处理备注"`
	HandleTime  *time.Time `json:"handleTime" db:"handle_time" gorm:"comment:处理时间"`
	IP          string    `json:"ip" db:"ip" gorm:"comment:举报IP"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Rating 评分模型
type Rating struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_target;index;not null;comment:用户ID"`
	TargetType  string    `json:"targetType" db:"target_type" gorm:"uniqueIndex:idx_user_target;not null;comment:目标类型:product/course/article"`
	TargetID    int64     `json:"targetId" db:"target_id" gorm:"uniqueIndex:idx_user_target;index;not null;comment:目标ID"`
	Score       float64   `json:"score" db:"score" gorm:"not null;comment:评分(1-5)"`
	Content     string    `json:"content" db:"content" gorm:"comment:评价内容"`
	IsAnonymous bool      `json:"isAnonymous" db:"is_anonymous" gorm:"default:false;comment:是否匿名"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-隐藏,-1-已删除"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// RatingDetail 评分详情
type RatingDetail struct {
	common.BaseIDModel
	RatingID    int64     `json:"ratingId" db:"rating_id" gorm:"index;not null;comment:评分ID"`
	AspectType  string    `json:"aspectType" db:"aspect_type" gorm:"not null;comment:评分方面,如质量/服务/性价比"`
	Score       float64   `json:"score" db:"score" gorm:"not null;comment:分项评分(1-5)"`
	Comment     string    `json:"comment" db:"comment" gorm:"comment:分项评价"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ReportProcessing 举报处理
type ReportProcessing struct {
	common.BaseIDModel
	common.BaseTimeModel
	ReportID      int64     `json:"reportId" db:"report_id" gorm:"index;not null;comment:举报ID"`
	ProcessorID   int64     `json:"processorId" db:"processor_id" gorm:"index;not null;comment:处理人ID"`
	ProcessorName string    `json:"processorName" db:"processor_name" gorm:"comment:处理人名称"`
	ActionType    string    `json:"actionType" db:"action_type" gorm:"not null;comment:处理动作:delete/hide/warn/ban/ignore"`
	ActionDetail  string    `json:"actionDetail" db:"action_detail" gorm:"comment:处理详情"`
	Remark        string    `json:"remark" db:"remark" gorm:"comment:备注"`
	Duration      int       `json:"duration" db:"duration" gorm:"comment:处罚时长(天)"`
	IP            string    `json:"ip" db:"ip" gorm:"comment:处理IP"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// SensitiveWord 敏感词
type SensitiveWord struct {
	common.BaseIDModel
	common.BaseTimeModel
	Word        string    `json:"word" db:"word" gorm:"uniqueIndex:idx_tenant_word;not null;comment:敏感词"`
	Category    string    `json:"category" db:"category" gorm:"index;comment:分类:politics/porn/abuse/violence/other"`
	Level       int       `json:"level" db:"level" gorm:"comment:级别:1-轻度,2-中度,3-重度"`
	Replacement string    `json:"replacement" db:"replacement" gorm:"comment:替换词"`
	IsRegex     bool      `json:"isRegex" db:"is_regex" gorm:"default:false;comment:是否正则表达式"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreateBy    int64     `json:"createBy" db:"create_by" gorm:"comment:创建人"`
	UpdateBy    int64     `json:"updateBy" db:"update_by" gorm:"comment:更新人"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"uniqueIndex:idx_tenant_word;index;comment:租户ID"`
}
