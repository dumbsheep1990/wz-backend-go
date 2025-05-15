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
	Username    string    `db:"username" json:"username"`
	Nickname    string    `db:"nickname" json:"nickname"`
	Email       string    `db:"email" json:"email"`
	Phone       string    `db:"phone" json:"phone"`
	Avatar      string    `db:"avatar" json:"avatar"`
	Status      int       `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	LastLoginAt time.Time `db:"last_login_at" json:"last_login_at"`
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
	OwnerId     int64     `db:"owner_id" json:"owner_id"`
	ExpireAt    time.Time `db:"expire_at" json:"expire_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Content 内容模型
type Content struct {
	Id           int64     `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	UserId       int64     `db:"user_id" json:"user_id"`
	CategoryId   int64     `db:"category_id" json:"category_id"`
	Tags         string    `db:"tags" json:"tags"`
	Cover        string    `db:"cover" json:"cover"`
	Summary      string    `db:"summary" json:"summary"`
	Content      string    `db:"content" json:"content"`
	Status       int       `db:"status" json:"status"`
	ViewCount    int       `db:"view_count" json:"view_count"`
	LikeCount    int       `db:"like_count" json:"like_count"`
	CommentCount int       `db:"comment_count" json:"comment_count"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Order 订单模型
type Order struct {
	Id          int64     `db:"id" json:"id"`
	OrderNo     string    `db:"order_no" json:"order_no"`
	UserId      int64     `db:"user_id" json:"user_id"`
	ProductId   int64     `db:"product_id" json:"product_id"`
	ProductName string    `db:"product_name" json:"product_name"`
	ProductType string    `db:"product_type" json:"product_type"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Price       float64   `db:"price" json:"price"`
	TotalAmount float64   `db:"total_amount" json:"total_amount"`
	Status      int       `db:"status" json:"status"`
	PayType     string    `db:"pay_type" json:"pay_type"`
	TradeNo     string    `db:"trade_no" json:"trade_no"`
	PayTime     time.Time `db:"pay_time" json:"pay_time"`
	Remark      string    `db:"remark" json:"remark"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Admin 管理员模型
type Admin struct {
	Id        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"` // 密码不返回给前端
	Role      string    `db:"role" json:"role"`
	Status    int       `db:"status" json:"status"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Role 角色模型
type Role struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Permissions string    `db:"permissions" json:"-"` // 数据库中存储的JSON字符串
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// 非数据库字段
	PermissionList []string `db:"-" json:"permissions"` // 解析后的权限列表
}

// OperationLog 操作日志模型
type OperationLog struct {
	Id          int64     `db:"id" json:"id"`
	UserId      int64     `db:"user_id" json:"user_id"`
	Username    string    `db:"username" json:"username"`
	Module      string    `db:"module" json:"module"`
	Action      string    `db:"action" json:"action"`
	Method      string    `db:"method" json:"method"`
	Path        string    `db:"path" json:"path"`
	Params      string    `db:"params" json:"params"`
	Result      string    `db:"result" json:"result"`
	Status      int       `db:"status" json:"status"`
	Ip          string    `db:"ip" json:"ip"`
	UserAgent   string    `db:"user_agent" json:"user_agent"`
	ExecuteTime int64     `db:"execute_time" json:"execute_time"` // 执行时间，单位毫秒
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
