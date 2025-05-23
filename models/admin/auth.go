package admin

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// SysAuthority 角色模型
type SysAuthority struct {
	common.BaseTimeModel
	AuthorityId   string `json:"authorityId" db:"authority_id" gorm:"primaryKey"`
	AuthorityName string `json:"authorityName" db:"authority_name"`
	ParentId      string `json:"parentId" db:"parent_id"`
	DefaultRouter string `json:"defaultRouter" db:"default_router"`
}

// GetID 获取角色ID
func (a SysAuthority) GetID() int64 {
	return 0 // 特殊情况，使用的是字符串ID
}

// SysAuthorityUser 用户权限关联
type SysAuthorityUser struct {
	UserId      int64  `json:"userId" db:"user_id" gorm:"primaryKey"`
	AuthorityId string `json:"authorityId" db:"authority_id" gorm:"primaryKey"`
}

// SysAuthorityBtn 按钮权限
type SysAuthorityBtn struct {
	common.BaseIDModel
	common.BaseTimeModel
	AuthorityId string `json:"authorityId" db:"authority_id" gorm:"index"`
	MenuId      int64  `json:"menuId" db:"menu_id" gorm:"index"`
	Name        string `json:"name" db:"name"`
	ButtonName  string `json:"buttonName" db:"button_name"`
}

// SysDataAuthority 数据权限
type SysDataAuthority struct {
	AuthorityId     string `json:"authorityId" db:"authority_id" gorm:"primaryKey"`
	DataAuthorityId string `json:"dataAuthorityId" db:"data_authority_id" gorm:"primaryKey"`
}

// SysMenu 菜单模型
type SysMenu struct {
	common.BaseIDModel
	common.BaseTimeModel
	ParentId  int64  `json:"parentId" db:"parent_id" gorm:"index"`
	Path      string `json:"path" db:"path"`
	Name      string `json:"name" db:"name"`
	Hidden    bool   `json:"hidden" db:"hidden"`
	Component string `json:"component" db:"component"`
	Sort      int    `json:"sort" db:"sort" gorm:"index"`
}

// SysMenuMeta 菜单元数据
type SysMenuMeta struct {
	MenuId           int64  `json:"menuId" db:"menu_id" gorm:"primaryKey"`
	KeepAlive        bool   `json:"keepAlive" db:"keep_alive"`
	DefaultMenu      bool   `json:"defaultMenu" db:"default_menu"`
	Title            string `json:"title" db:"title"`
	Icon             string `json:"icon" db:"icon"`
	CloseTab         bool   `json:"closeTab" db:"close_tab"`
	CollapsibleWidth int    `json:"collapsibleWidth" db:"collapsible_width"`
}

// SysMenuAuthority 菜单权限关联
type SysMenuAuthority struct {
	MenuId      int64  `json:"menuId" db:"menu_id" gorm:"primaryKey"`
	AuthorityId string `json:"authorityId" db:"authority_id" gorm:"primaryKey"`
}

// CasbinRule Casbin规则模型
type CasbinRule struct {
	common.BaseIDModel
	Ptype string `json:"ptype" db:"ptype"`
	V0    string `json:"v0" db:"v0"`
	V1    string `json:"v1" db:"v1"`
	V2    string `json:"v2" db:"v2"`
	V3    string `json:"v3" db:"v3"`
	V4    string `json:"v4" db:"v4"`
	V5    string `json:"v5" db:"v5"`
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
	common.BaseIDModel
	Jwt       string    `json:"jwt" db:"jwt"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
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
