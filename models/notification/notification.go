package notification

import (
	"time"

	"wz-backend-go/models/common"
)

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null;comment:模板名称"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:模板编码"`
	Type        string `json:"type" db:"type" gorm:"not null;comment:通知类型:email,sms,push,inbox,wechat"`
	Title       string `json:"title" db:"title" gorm:"not null;comment:标题模板"`
	Content     string `json:"content" db:"content" gorm:"type:text;not null;comment:内容模板"`
	Variables   string `json:"variables" db:"variables" gorm:"type:json;comment:变量定义,JSON格式"`
	Example     string `json:"example" db:"example" gorm:"type:json;comment:样例数据,JSON格式"`
	Description string `json:"description" db:"description" gorm:"comment:描述"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	Category    string `json:"category" db:"category" gorm:"index;comment:分类:system,marketing,transaction,reminder"`
	CreatedBy   int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy   int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationChannel 通知渠道
type NotificationChannel struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name         string `json:"name" db:"name" gorm:"not null;comment:渠道名称"`
	Code         string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:渠道编码"`
	Type         string `json:"type" db:"type" gorm:"not null;comment:渠道类型:email,sms,push,wechat,dingtalk,slack"`
	Provider     string `json:"provider" db:"provider" gorm:"not null;comment:提供商:aliyun,tencent,aws,jiguang,custom"`
	Configuration string `json:"configuration" db:"configuration" gorm:"type:json;comment:配置信息,JSON格式"`
	IsDefault    bool   `json:"isDefault" db:"is_default" gorm:"default:false;comment:是否默认"`
	Status       int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreatedBy    int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy    int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID     int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationEvent 通知事件
type NotificationEvent struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null;comment:事件名称"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:事件编码"`
	Description string `json:"description" db:"description" gorm:"comment:描述"`
	SourceType  string `json:"sourceType" db:"source_type" gorm:"comment:来源类型:system,user,schedule"`
	Category    string `json:"category" db:"category" gorm:"index;comment:分类:system,business,marketing"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationRule 通知规则
type NotificationRule struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name            string `json:"name" db:"name" gorm:"not null;comment:规则名称"`
	EventID         int64  `json:"eventId" db:"event_id" gorm:"index;not null;comment:事件ID"`
	TemplateIDs     string `json:"templateIds" db:"template_ids" gorm:"type:json;not null;comment:模板ID列表,JSON格式"`
	ChannelIDs      string `json:"channelIds" db:"channel_ids" gorm:"type:json;not null;comment:渠道ID列表,JSON格式"`
	ReceiverType    string `json:"receiverType" db:"receiver_type" gorm:"not null;comment:接收者类型:user,role,department,all"`
	ReceiverIDs     string `json:"receiverIds" db:"receiver_ids" gorm:"type:json;comment:接收者ID列表,JSON格式"`
	Conditions      string `json:"conditions" db:"conditions" gorm:"type:json;comment:触发条件,JSON格式"`
	Priority        int    `json:"priority" db:"priority" gorm:"default:5;comment:优先级:1-10,越小越高"`
	CooldownMinutes int    `json:"cooldownMinutes" db:"cooldown_minutes" gorm:"default:0;comment:冷却时间(分钟)"`
	StartTime       string `json:"startTime" db:"start_time" gorm:"comment:生效开始时间(HH:mm)"`
	EndTime         string `json:"endTime" db:"end_time" gorm:"comment:生效结束时间(HH:mm)"`
	Status          int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreatedBy       int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy       int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID        int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Notification 通知记录
type Notification struct {
	common.BaseIDModel
	common.BaseTimeModel
	Title          string     `json:"title" db:"title" gorm:"not null;comment:标题"`
	Content        string     `json:"content" db:"content" gorm:"type:text;not null;comment:内容"`
	Type           string     `json:"type" db:"type" gorm:"not null;comment:通知类型:email,sms,push,inbox,wechat"`
	Status         string     `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-待发送,sending-发送中,sent-已发送,failed-发送失败"`
	EventID        int64      `json:"eventId" db:"event_id" gorm:"index;comment:事件ID"`
	TemplateID     int64      `json:"templateId" db:"template_id" gorm:"index;comment:模板ID"`
	ChannelID      int64      `json:"channelId" db:"channel_id" gorm:"index;comment:渠道ID"`
	ReceiverID     int64      `json:"receiverId" db:"receiver_id" gorm:"index;comment:接收者ID"`
	ReceiverType   string     `json:"receiverType" db:"receiver_type" gorm:"comment:接收者类型:user,role,department"`
	ReceiverAddr   string     `json:"receiverAddr" db:"receiver_addr" gorm:"comment:接收地址(手机/邮箱)"`
	Priority       int        `json:"priority" db:"priority" gorm:"default:5;comment:优先级:1-10,越小越高"`
	ReadStatus     string     `json:"readStatus" db:"read_status" gorm:"default:unread;comment:阅读状态:unread-未读,read-已读"`
	ReadTime       *time.Time `json:"readTime" db:"read_time" gorm:"comment:阅读时间"`
	SentTime       *time.Time `json:"sentTime" db:"sent_time" gorm:"comment:发送时间"`
	FailReason     string     `json:"failReason" db:"fail_reason" gorm:"comment:失败原因"`
	RetryCount     int        `json:"retryCount" db:"retry_count" gorm:"default:0;comment:重试次数"`
	NextRetryTime  *time.Time `json:"nextRetryTime" db:"next_retry_time" gorm:"comment:下次重试时间"`
	Data           string     `json:"data" db:"data" gorm:"type:json;comment:通知数据,JSON格式"`
	ExternalID     string     `json:"externalId" db:"external_id" gorm:"comment:外部ID(如短信/邮件ID)"`
	ExternalStatus string     `json:"externalStatus" db:"external_status" gorm:"comment:外部状态"`
	TenantID       int64      `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationUserPreference 用户通知偏好
type NotificationUserPreference struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID               int64  `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_event;index;not null;comment:用户ID"`
	EventID              int64  `json:"eventId" db:"event_id" gorm:"uniqueIndex:idx_user_event;index;comment:事件ID,为空表示全局设置"`
	EmailEnabled         bool   `json:"emailEnabled" db:"email_enabled" gorm:"default:true;comment:邮件是否启用"`
	SMSEnabled           bool   `json:"smsEnabled" db:"sms_enabled" gorm:"default:true;comment:短信是否启用"`
	PushEnabled          bool   `json:"pushEnabled" db:"push_enabled" gorm:"default:true;comment:推送是否启用"`
	InboxEnabled         bool   `json:"inboxEnabled" db:"inbox_enabled" gorm:"default:true;comment:站内信是否启用"`
	WechatEnabled        bool   `json:"wechatEnabled" db:"wechat_enabled" gorm:"default:true;comment:微信是否启用"`
	QuietTimeStart       string `json:"quietTimeStart" db:"quiet_time_start" gorm:"comment:免打扰开始时间(HH:mm)"`
	QuietTimeEnd         string `json:"quietTimeEnd" db:"quiet_time_end" gorm:"comment:免打扰结束时间(HH:mm)"`
	QuietTimeEnabled     bool   `json:"quietTimeEnabled" db:"quiet_time_enabled" gorm:"default:false;comment:免打扰是否启用"`
	GroupNotificationType string `json:"groupNotificationType" db:"group_notification_type" gorm:"default:all;comment:群通知类型:all-所有消息,mention-仅@我,none-无"`
	TenantID             int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationGroup 通知分组
type NotificationGroup struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64      `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Type        string     `json:"type" db:"type" gorm:"not null;comment:分组类型:system,business,marketing"`
	Count       int        `json:"count" db:"count" gorm:"default:0;comment:通知数量"`
	UnreadCount int        `json:"unreadCount" db:"unread_count" gorm:"default:0;comment:未读数量"`
	LastID      int64      `json:"lastId" db:"last_id" gorm:"comment:最新通知ID"`
	LastTime    *time.Time `json:"lastTime" db:"last_time" gorm:"comment:最新通知时间"`
	TenantID    int64      `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// NotificationDeliveryLog 通知投递日志
type NotificationDeliveryLog struct {
	common.BaseIDModel
	common.BaseTimeModel
	NotificationID int64  `json:"notificationId" db:"notification_id" gorm:"index;not null;comment:通知ID"`
	ChannelID      int64  `json:"channelId" db:"channel_id" gorm:"index;not null;comment:渠道ID"`
	Status         string `json:"status" db:"status" gorm:"not null;comment:状态:success,failed"`
	Request        string `json:"request" db:"request" gorm:"type:text;comment:请求内容"`
	Response       string `json:"response" db:"response" gorm:"type:text;comment:响应内容"`
	ErrorMessage   string `json:"errorMessage" db:"error_message" gorm:"comment:错误信息"`
	Duration       int    `json:"duration" db:"duration" gorm:"comment:耗时(毫秒)"`
	ExternalID     string `json:"externalId" db:"external_id" gorm:"comment:外部ID"`
	IP             string `json:"ip" db:"ip" gorm:"comment:IP地址"`
	TenantID       int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
