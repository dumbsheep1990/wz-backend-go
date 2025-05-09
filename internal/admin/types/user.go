package types

// 用户分页查询请求
type UserListReq struct {
	Page      int32  `form:"page,default=1"`
	PageSize  int32  `form:"pageSize,default=10"`
	Username  string `form:"username,optional"`
	Email     string `form:"email,optional"`
	Phone     string `form:"phone,optional"`
	Status    int32  `form:"status,optional"`
	Role      string `form:"role,optional"`
	StartTime string `form:"startTime,optional"`
	EndTime   string `form:"endTime,optional"`
}

// 用户信息详情
type UserDetail struct {
	ID                int64  `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Role              string `json:"role"`
	Status            int32  `json:"status"`
	IsVerified        bool   `json:"is_verified"`
	IsCompanyVerified bool   `json:"is_company_verified"`
	DefaultTenantID   int64  `json:"default_tenant_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// 用户列表响应
type UserListResp struct {
	Total int64        `json:"total"`
	List  []UserDetail `json:"list"`
}

// 用户详情请求
type UserDetailReq struct {
	ID int64 `path:"id"`
}

// 管理员创建用户请求
type AdminCreateUserReq struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	Role            string `json:"role" validate:"required"`
	Status          int32  `json:"status" validate:"required"`
	DefaultTenantID int64  `json:"default_tenant_id,optional"`
}

// 管理员创建用户响应
type AdminCreateUserResp struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// 管理员更新用户请求
type AdminUpdateUserReq struct {
	ID              int64  `path:"id"`
	Username        string `json:"username,optional"`
	Password        string `json:"password,optional"`
	Email           string `json:"email,optional"`
	Phone           string `json:"phone,optional"`
	Role            string `json:"role,optional"`
	Status          int32  `json:"status,optional"`
	DefaultTenantID int64  `json:"default_tenant_id,optional"`
}

// 删除用户请求
type DeleteUserReq struct {
	ID int64 `path:"id"`
}
