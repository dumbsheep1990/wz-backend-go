package model

import "time"

// 通用响应
type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// User 用户模型
type User struct {
	Id          int64     `db:"id" json:"id"`
	UUID        string    `db:"uuid" json:"uuid"`
	Username    string    `db:"username" json:"username"`
	Nickname    string    `db:"nickname" json:"nickname"`
	Password    string    `db:"password" json:"-"` // 不返回给前端
	Email       string    `db:"email" json:"email"`
	Phone       string    `db:"phone" json:"phone"`
	Avatar      string    `db:"avatar" json:"avatar"`
	AuthorityId string    `db:"authority_id" json:"authorityId"`
	Status      int       `db:"status" json:"status"`
	SideMode    string    `db:"side_mode" json:"sideMode"`
	ActiveColor string    `db:"active_color" json:"activeColor"`
	BaseColor   string    `db:"base_color" json:"baseColor"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	LastLoginAt time.Time `db:"last_login_at" json:"lastLoginAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysAuthority 角色模型
type SysAuthority struct {
	AuthorityId   string    `db:"authority_id" json:"authorityId"`
	AuthorityName string    `db:"authority_name" json:"authorityName"`
	ParentId      string    `db:"parent_id" json:"parentId"`
	DefaultRouter string    `db:"default_router" json:"defaultRouter"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt     time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysAuthorityUser 用户权限关联
type SysAuthorityUser struct {
	UserId      int64  `db:"user_id" json:"userId"`
	AuthorityId string `db:"authority_id" json:"authorityId"`
}

// SysAuthorityBtn 按钮权限
type SysAuthorityBtn struct {
	Id          int64     `db:"id" json:"id"`
	AuthorityId string    `db:"authority_id" json:"authorityId"`
	MenuId      int64     `db:"menu_id" json:"menuId"`
	Name        string    `db:"name" json:"name"`
	ButtonName  string    `db:"button_name" json:"buttonName"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysDataAuthority 数据权限
type SysDataAuthority struct {
	AuthorityId     string `db:"authority_id" json:"authorityId"`
	DataAuthorityId string `db:"data_authority_id" json:"dataAuthorityId"`
}

// SysMenu 菜单模型
type SysMenu struct {
	Id        int64     `db:"id" json:"id"`
	ParentId  int64     `db:"parent_id" json:"parentId"`
	Path      string    `db:"path" json:"path"`
	Name      string    `db:"name" json:"name"`
	Hidden    bool      `db:"hidden" json:"hidden"`
	Component string    `db:"component" json:"component"`
	Sort      int       `db:"sort" json:"sort"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysMenuMeta 菜单元数据
type SysMenuMeta struct {
	MenuId           int64  `db:"menu_id" json:"menuId"`
	KeepAlive        bool   `db:"keep_alive" json:"keepAlive"`
	DefaultMenu      bool   `db:"default_menu" json:"defaultMenu"`
	Title            string `db:"title" json:"title"`
	Icon             string `db:"icon" json:"icon"`
	CloseTab         bool   `db:"close_tab" json:"closeTab"`
	CollapsibleWidth int    `db:"collapsible_width" json:"collapsibleWidth"`
}

// SysMenuAuthority 菜单权限关联
type SysMenuAuthority struct {
	MenuId      int64  `db:"menu_id" json:"menuId"`
	AuthorityId string `db:"authority_id" json:"authorityId"`
}

// SysApi API模型
type SysApi struct {
	Id          int64     `db:"id" json:"id"`
	Path        string    `db:"path" json:"path"`
	Description string    `db:"description" json:"description"`
	ApiGroup    string    `db:"api_group" json:"apiGroup"`
	Method      string    `db:"method" json:"method"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// CasbinRule Casbin规则模型
type CasbinRule struct {
	Id    int64  `db:"id" json:"id"`
	Ptype string `db:"ptype" json:"ptype"`
	V0    string `db:"v0" json:"v0"`
	V1    string `db:"v1" json:"v1"`
	V2    string `db:"v2" json:"v2"`
	V3    string `db:"v3" json:"v3"`
	V4    string `db:"v4" json:"v4"`
	V5    string `db:"v5" json:"v5"`
}

// SysDictionary 字典模型
type SysDictionary struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Type        string    `db:"type" json:"type"`
	Status      int       `db:"status" json:"status"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysDictionaryDetail 字典详情模型
type SysDictionaryDetail struct {
	Id         int64     `db:"id" json:"id"`
	Label      string    `db:"label" json:"label"`
	Value      string    `db:"value" json:"value"`
	Status     int       `db:"status" json:"status"`
	Sort       int       `db:"sort" json:"sort"`
	DictTypeId int64     `db:"dict_type_id" json:"sysDictionaryID"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt  time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysOperationRecord 操作日志模型
type SysOperationRecord struct {
	Id           int64     `db:"id" json:"id"`
	Ip           string    `db:"ip" json:"ip"`
	Method       string    `db:"method" json:"method"`
	Path         string    `db:"path" json:"path"`
	Status       int       `db:"status" json:"status"`
	Latency      int64     `db:"latency" json:"latency"`
	Agent        string    `db:"agent" json:"agent"`
	ErrorMessage string    `db:"error_message" json:"errorMessage"`
	Body         string    `db:"body" json:"body"`
	UserId       int64     `db:"user_id" json:"userId"`
	User         string    `db:"user" json:"user"`
	Resp         string    `db:"resp" json:"resp"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt    time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SysParams 系统参数模型
type SysParams struct {
	Id         int64     `db:"id" json:"id"`
	ParamName  string    `db:"param_name" json:"paramName"`
	ParamKey   string    `db:"param_key" json:"paramKey"`
	ParamValue string    `db:"param_value" json:"paramValue"`
	ParamType  string    `db:"param_type" json:"paramType"`
	ParamDesc  string    `db:"param_desc" json:"paramDesc"`
	ParamGroup string    `db:"param_group" json:"paramGroup"`
	Status     int       `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt  time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	Id                int64     `db:"id" json:"id"`
	SystemName        string    `db:"system_name" json:"systemName"`
	SystemLogo        string    `db:"system_logo" json:"systemLogo"`
	SystemApi         string    `db:"system_api" json:"systemApi"`
	SystemDomain      string    `db:"system_domain" json:"systemDomain"`
	SystemColor       string    `db:"system_color" json:"systemColor"`
	SystemMode        string    `db:"system_mode" json:"systemMode"`
	AdminDomain       string    `db:"admin_domain" json:"adminDomain"`
	FileStoreType     string    `db:"file_store_type" json:"fileStoreType"`
	LogSaveType       string    `db:"log_save_type" json:"logSaveType"`
	LogSavePath       string    `db:"log_save_path" json:"logSavePath"`
	LogLevel          string    `db:"log_level" json:"logLevel"`
	LogMaxSize        int       `db:"log_max_size" json:"logMaxSize"`
	LogMaxBackups     int       `db:"log_max_backups" json:"logMaxBackups"`
	LogMaxAge         int       `db:"log_max_age" json:"logMaxAge"`
	LogCompress       bool      `db:"log_compress" json:"logCompress"`
	EmailHost         string    `db:"email_host" json:"emailHost"`
	EmailPort         int       `db:"email_port" json:"emailPort"`
	EmailFrom         string    `db:"email_from" json:"emailFrom"`
	EmailSecret       string    `db:"email_secret" json:"emailSecret"`
	EmailIsSSL        bool      `db:"email_is_ssl" json:"emailIsSSL"`
	EmailNickname     string    `db:"email_nickname" json:"emailNickname"`
	LdapOpen          bool      `db:"ldap_open" json:"ldapOpen"`
	LdapHost          string    `db:"ldap_host" json:"ldapHost"`
	LdapPort          int       `db:"ldap_port" json:"ldapPort"`
	LdapBindDn        string    `db:"ldap_bind_dn" json:"ldapBindDn"`
	LdapPassword      string    `db:"ldap_password" json:"ldapPassword"`
	LdapBaseDn        string    `db:"ldap_base_dn" json:"ldapBaseDn"`
	LdapUserField     string    `db:"ldap_user_field" json:"ldapUserField"`
	LdapNickName      string    `db:"ldap_nick_name" json:"ldapNickName"`
	LdapEmails        string    `db:"ldap_emails" json:"ldapEmails"`
	LdapTls           bool      `db:"ldap_tls" json:"ldapTls"`
	DingTalkOpen      bool      `db:"dingtalk_open" json:"dingTalkOpen"`
	DingTalkAppKey    string    `db:"dingtalk_app_key" json:"dingTalkAppKey"`
	DingTalkAppSecret string    `db:"dingtalk_app_secret" json:"dingTalkAppSecret"`
	DingTalkAgentId   string    `db:"dingtalk_agent_id" json:"dingTalkAgentId"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt         time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// Tenant 租户模型
type Tenant struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Type        int       `db:"type" json:"type"`
	Status      int       `db:"status" json:"status"`
	Logo        string    `db:"logo" json:"logo"`
	Description string    `db:"description" json:"description"`
	OwnerId     int64     `db:"owner_id" json:"ownerId"`
	ExpireAt    time.Time `db:"expire_at" json:"expireAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// Content 内容模型
type Content struct {
	Id           int64     `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	UserId       int64     `db:"user_id" json:"userId"`
	CategoryId   int64     `db:"category_id" json:"categoryId"`
	Tags         string    `db:"tags" json:"tags"`
	Cover        string    `db:"cover" json:"cover"`
	Summary      string    `db:"summary" json:"summary"`
	Content      string    `db:"content" json:"content"`
	Status       int       `db:"status" json:"status"`
	ViewCount    int       `db:"view_count" json:"viewCount"`
	LikeCount    int       `db:"like_count" json:"likeCount"`
	CommentCount int       `db:"comment_count" json:"commentCount"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt    time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// Order 订单模型
type Order struct {
	Id          int64     `db:"id" json:"id"`
	OrderNo     string    `db:"order_no" json:"orderNo"`
	UserId      int64     `db:"user_id" json:"userId"`
	ProductId   int64     `db:"product_id" json:"productId"`
	ProductName string    `db:"product_name" json:"productName"`
	ProductType string    `db:"product_type" json:"productType"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Price       float64   `db:"price" json:"price"`
	TotalAmount float64   `db:"total_amount" json:"totalAmount"`
	Status      int       `db:"status" json:"status"`
	PayType     string    `db:"pay_type" json:"payType"`
	TradeNo     string    `db:"trade_no" json:"tradeNo"`
	PayTime     time.Time `db:"pay_time" json:"payTime"`
	Remark      string    `db:"remark" json:"remark"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// Admin 管理员模型
type Admin struct {
	Id        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"` // 密码不返回给前端
	Role      string    `db:"role" json:"role"`
	Status    int       `db:"status" json:"status"`
	LastLogin time.Time `db:"last_login" json:"lastLogin"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha,optional"`
	CaptchaId string `json:"captchaId,optional"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// JwtClaims JWT载荷
type JwtClaims struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	AuthorityId string `json:"authorityId"`
	BufferTime  int64  `json:"bufferTime"`
}

// JwtBlacklist JWT黑名单
type JwtBlacklist struct {
	Id        int64     `db:"id" json:"id"`
	Jwt       string    `db:"jwt" json:"jwt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
