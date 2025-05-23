package admin

import "wz-project/wz-backend-go/models/common"

// SystemConfig 系统配置模型
type SystemConfig struct {
	common.BaseIDModel
	SystemName        string `json:"systemName" db:"system_name"`
	SystemLogo        string `json:"systemLogo" db:"system_logo"`
	SystemApi         string `json:"systemApi" db:"system_api"`
	SystemDomain      string `json:"systemDomain" db:"system_domain"`
	SystemColor       string `json:"systemColor" db:"system_color"`
	SystemMode        string `json:"systemMode" db:"system_mode" gorm:"comment:light,dark,auto"`
	AdminDomain       string `json:"adminDomain" db:"admin_domain"`
	FileStoreType     string `json:"fileStoreType" db:"file_store_type" gorm:"comment:local,aliyun,qiniu,tencent"`
	LogSaveType       string `json:"logSaveType" db:"log_save_type" gorm:"comment:file,database,elasticsearch"`
	LogSavePath       string `json:"logSavePath" db:"log_save_path"`
	LogLevel          string `json:"logLevel" db:"log_level" gorm:"comment:debug,info,warn,error"`
	LogMaxSize        int    `json:"logMaxSize" db:"log_max_size" gorm:"comment:单位MB"`
	LogMaxBackups     int    `json:"logMaxBackups" db:"log_max_backups"`
	LogMaxAge         int    `json:"logMaxAge" db:"log_max_age" gorm:"comment:单位天"`
	LogCompress       bool   `json:"logCompress" db:"log_compress"`
	EmailHost         string `json:"emailHost" db:"email_host"`
	EmailPort         int    `json:"emailPort" db:"email_port"`
	EmailUsername     string `json:"emailUsername" db:"email_username"`
	EmailPassword     string `json:"emailPassword" db:"email_password"`
	EmailFrom         string `json:"emailFrom" db:"email_from"`
	EmailTls          bool   `json:"emailTls" db:"email_tls"`
	SmsType           string `json:"smsType" db:"sms_type" gorm:"comment:aliyun,tencent,netease"`
	SmsAccessKey      string `json:"smsAccessKey" db:"sms_access_key"`
	SmsSecretKey      string `json:"smsSecretKey" db:"sms_secret_key"`
	SmsRegionId       string `json:"smsRegionId" db:"sms_region_id"`
	SmsSignName       string `json:"smsSignName" db:"sms_sign_name"`
	SmsTemplateCode   string `json:"smsTemplateCode" db:"sms_template_code"`
	JwtSecret         string `json:"jwtSecret" db:"jwt_secret"`
	JwtTimeout        int    `json:"jwtTimeout" db:"jwt_timeout" gorm:"comment:单位分钟"`
	JwtIssuer         string `json:"jwtIssuer" db:"jwt_issuer"`
	CaptchaType       string `json:"captchaType" db:"captcha_type" gorm:"comment:image,sms,email"`
	CaptchaWidth      int    `json:"captchaWidth" db:"captcha_width"`
	CaptchaHeight     int    `json:"captchaHeight" db:"captcha_height"`
	CaptchaLength     int    `json:"captchaLength" db:"captcha_length"`
	CaptchaExpire     int    `json:"captchaExpire" db:"captcha_expire" gorm:"comment:单位分钟"`
	RegisterLimit     int    `json:"registerLimit" db:"register_limit" gorm:"comment:单位次/天"`
	LoginLimit        int    `json:"loginLimit" db:"login_limit" gorm:"comment:单位次/天"`
	LoginErrorLimit   int    `json:"loginErrorLimit" db:"login_error_limit"`
	LoginLockTime     int    `json:"loginLockTime" db:"login_lock_time" gorm:"comment:单位分钟"`
	EnableCsrf        bool   `json:"enableCsrf" db:"enable_csrf"`
	EnableRateLimit   bool   `json:"enableRateLimit" db:"enable_rate_limit"`
	RateLimitPerIp    int    `json:"rateLimitPerIp" db:"rate_limit_per_ip" gorm:"comment:单位次/分钟"`
	RateLimitPerUser  int    `json:"rateLimitPerUser" db:"rate_limit_per_user" gorm:"comment:单位次/分钟"`
	EnableDashboard   bool   `json:"enableDashboard" db:"enable_dashboard"`
}

// DashboardLayout 仪表盘布局
type DashboardLayout struct {
	common.BaseIDModel
	common.BaseTimeModel
	common.BaseTenantModel
	Name       string `json:"name" db:"name"`
	Type       string `json:"type" db:"type" gorm:"comment:admin,tenant,user"`
	TargetID   int64  `json:"targetId" db:"target_id" gorm:"comment:根据类型存储对应的ID"`
	Layout     string `json:"layout" db:"layout" gorm:"type:json"`
	IsDefault  bool   `json:"isDefault" db:"is_default"`
	IsSystem   bool   `json:"isSystem" db:"is_system"`
	Status     int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	CreateBy   int64  `json:"createBy" db:"create_by"`
	UpdateBy   int64  `json:"updateBy" db:"update_by"`
}

// DashboardWidget 仪表盘部件
type DashboardWidget struct {
	common.BaseIDModel
	common.BaseTimeModel
	common.BaseTenantModel
	Name         string `json:"name" db:"name"`
	Type         string `json:"type" db:"type"`
	Icon         string `json:"icon" db:"icon"`
	Component    string `json:"component" db:"component"`
	DefaultSize  string `json:"defaultSize" db:"default_size" gorm:"comment:格式如 1x1, 2x1 表示列x行"`
	Description  string `json:"description" db:"description"`
	DefaultData  string `json:"defaultData" db:"default_data" gorm:"type:json"`
	Permissions  string `json:"permissions" db:"permissions" gorm:"type:json;comment:JSON数组,包含角色或权限列表"`
	IsSystem     bool   `json:"isSystem" db:"is_system"`
	Status       int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
}

// WidgetInstance 部件实例
type WidgetInstance struct {
	common.BaseIDModel
	common.BaseTimeModel
	LayoutID   int64  `json:"layoutId" db:"layout_id" gorm:"index;not null"`
	WidgetID   int64  `json:"widgetId" db:"widget_id" gorm:"index;not null"`
	PositionX  int    `json:"positionX" db:"position_x"`
	PositionY  int    `json:"positionY" db:"position_y"`
	Width      int    `json:"width" db:"width" gorm:"comment:以网格单位计算"`
	Height     int    `json:"height" db:"height" gorm:"comment:以网格单位计算"`
	CustomData string `json:"customData" db:"custom_data" gorm:"type:json"`
	CustomStyle string `json:"customStyle" db:"custom_style" gorm:"type:json"`
	SortOrder  int    `json:"sortOrder" db:"sort_order"`
}
