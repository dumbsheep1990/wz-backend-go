package types

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ========== 登录相关 ==========

// 登录请求
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

// 登录响应
type LoginResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	User      AdminInfo `json:"user"`
	Token     string    `json:"token"`
	ExpiresAt int64     `json:"expiresAt"`
}

// 管理员基本信息
type AdminInfo struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	RoleName  string `json:"roleName"`
	Avatar    string `json:"avatar"`
	Status    int    `json:"status"`
	LastLogin string `json:"lastLogin"`
	CreatedAt string `json:"createdAt"`
}

// 刷新Token请求
type RefreshTokenRequest struct {
	Token string `json:"token"`
}

// 刷新Token响应
type RefreshTokenResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    RefreshTokenResponseData `json:"data"`
}

type RefreshTokenResponseData struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// 验证码响应
type CaptchaResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    CaptchaResponseData  `json:"data"`
}

type CaptchaResponseData struct {
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
}

// ========== 用户管理相关 ==========

// 分页请求
type PageRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword,optional" json:"keyword,optional"`
}

// 分页数据返回
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 获取用户列表请求
type GetUserListRequest struct {
	PageRequest
	Username     string `json:"username,optional"`
	Status       int    `json:"status,optional"`
	AuthorityId  string `json:"authorityId,optional"`
	CreatedAfter string `json:"createdAfter,optional"`
	CreatedBefore string `json:"createdBefore,optional"`
}

// 获取用户列表响应
type GetUserListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    UserPageData `json:"data"`
}

type UserPageData struct {
	List     []UserInfo `json:"list"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
}

// 用户信息
type UserInfo struct {
	Id          int64       `json:"id"`
	UUID        string      `json:"uuid"`
	Username    string      `json:"userName"`
	NickName    string      `json:"nickName"`
	Email       string      `json:"email"`
	Phone       string      `json:"phone"`
	Avatar      string      `json:"headerImg"`
	Status      int         `json:"status"`
	Authorities []Authority `json:"authorities"`
	AuthorityId string      `json:"authorityId"`
	SideMode    string      `json:"sideMode"`
	ActiveColor string      `json:"activeColor"`
	BaseColor   string      `json:"baseColor"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
	LastLoginAt string      `json:"lastLoginAt"`
}

// 权限信息
type Authority struct {
	AuthorityId   string `json:"authorityId"`
	AuthorityName string `json:"authorityName"`
	ParentId      string `json:"parentId"`
	DefaultRouter string `json:"defaultRouter"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// 获取用户详情请求
type GetUserInfoRequest struct {
	ID int64 `form:"id,optional" json:"id,optional"`
}

// 获取用户详情响应
type GetUserInfoResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    UserInfo `json:"data"`
}

// 创建用户请求
type CreateUserRequest struct {
	Username    string `json:"userName"`
	Password    string `json:"password"`
	NickName    string `json:"nickName"`
	Email       string `json:"email,optional"`
	Phone       string `json:"phone,optional"`
	Avatar      string `json:"headerImg,optional"`
	AuthorityId string `json:"authorityId"`
}

// 更新用户请求
type UpdateUserRequest struct {
	Id          int64  `json:"id"`
	Username    string `json:"userName"`
	NickName    string `json:"nickName"`
	Email       string `json:"email,optional"`
	Phone       string `json:"phone,optional"`
	Avatar      string `json:"headerImg,optional"`
	Status      int    `json:"status,optional"`
	AuthorityId string `json:"authorityId,optional"`
}

// 更新自身信息请求
type UpdateSelfInfoRequest struct {
	NickName    string `json:"nickName"`
	Email       string `json:"email,optional"`
	Phone       string `json:"phone,optional"`
	Avatar      string `json:"headerImg,optional"`
	SideMode    string `json:"sideMode,optional"`
	ActiveColor string `json:"activeColor,optional"`
	BaseColor   string `json:"baseColor,optional"`
}

// 删除用户请求
type DeleteUserRequest struct {
	Id int64 `json:"id"`
}

// 修改密码请求
type ChangePasswordRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// 重置密码请求
type ResetPasswordRequest struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
}

// 设置用户权限请求
type SetUserAuthorityRequest struct {
	UserId      int64  `json:"userId"`
	AuthorityId string `json:"authorityId"`
}

// 设置用户权限组请求
type SetUserAuthoritiesRequest struct {
	UserId       int64    `json:"userId"`
	AuthorityIds []string `json:"authorityIds"`
}

// ========== 角色权限管理相关 ==========

// 角色信息
type SysAuthority struct {
	AuthorityId      string         `json:"authorityId"`
	AuthorityName    string         `json:"authorityName"`
	ParentId         string         `json:"parentId"`
	DataAuthorityId  []SysAuthority `json:"dataAuthorityId"`
	DefaultRouter    string         `json:"defaultRouter"`
	CreatedAt        string         `json:"createdAt"`
	UpdatedAt        string         `json:"updatedAt"`
	DeletedAt        string         `json:"deletedAt,omitempty"`
	Children         []SysAuthority `json:"children,omitempty"`
	Menus            []MenuInfo     `json:"menus,omitempty"`
	MenuIds          []int64        `json:"menuIds,omitempty"`
}

// 获取角色列表响应
type GetAuthorityListResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []SysAuthority `json:"data"`
}

// 创建角色请求
type CreateAuthorityRequest struct {
	AuthorityId   string `json:"authorityId"`
	AuthorityName string `json:"authorityName"`
	ParentId      string `json:"parentId"`
}

// 创建角色响应
type CreateAuthorityResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    SysAuthority `json:"data"`
}

// 更新角色请求
type UpdateAuthorityRequest struct {
	AuthorityId   string `json:"authorityId"`
	AuthorityName string `json:"authorityName"`
	ParentId      string `json:"parentId"`
}

// 更新角色响应
type UpdateAuthorityResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    SysAuthority `json:"data"`
}

// 删除角色请求
type DeleteAuthorityRequest struct {
	AuthorityId string `json:"authorityId"`
}

// 设置数据权限请求
type SetDataAuthorityRequest struct {
	AuthorityId     string   `json:"authorityId"`
	DataAuthorityId []string `json:"dataAuthorityId"`
}

// 菜单信息
type MenuInfo struct {
	ID        int64      `json:"ID"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
	DeletedAt string     `json:"deletedAt,omitempty"`
	ParentId  int64      `json:"parentId"`
	Path      string     `json:"path"`
	Name      string     `json:"name"`
	Hidden    bool       `json:"hidden"`
	Component string     `json:"component"`
	Sort      int        `json:"sort"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuInfo `json:"children,omitempty"`
}

// 菜单元数据
type MenuMeta struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	KeepAlive   bool   `json:"keepAlive"`
	DefaultMenu bool   `json:"defaultMenu"`
}
