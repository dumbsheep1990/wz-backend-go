package user

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// User 用户模型
type User struct {
	common.BaseIDModel
	common.BaseTimeModel
	UUID        string    `json:"uuid" db:"uuid" gorm:"uniqueIndex;type:varchar(36)"`
	Username    string    `json:"username" db:"username" gorm:"uniqueIndex;not null"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Password    string    `json:"-" db:"password" gorm:"not null"`
	Email       string    `json:"email" db:"email" gorm:"uniqueIndex"`
	Phone       string    `json:"phone" db:"phone" gorm:"index"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Gender      int       `json:"gender" db:"gender" gorm:"default:0;comment:0-未知,1-男,2-女"`
	Birthday    time.Time `json:"birthday" db:"birthday"`
	Bio         string    `json:"bio" db:"bio"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用,2-未验证"`
	LastLoginAt time.Time `json:"lastLoginAt" db:"last_login_at"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// UserProfile 用户详细资料
type UserProfile struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"uniqueIndex;not null"`
	RealName     string    `json:"realName" db:"real_name"`
	IDCardType   string    `json:"idCardType" db:"id_card_type"`
	IDCardNo     string    `json:"idCardNo" db:"id_card_no"`
	Country      string    `json:"country" db:"country"`
	Province     string    `json:"province" db:"province"`
	City         string    `json:"city" db:"city"`
	District     string    `json:"district" db:"district"`
	Address      string    `json:"address" db:"address"`
	PostCode     string    `json:"postCode" db:"post_code"`
	Company      string    `json:"company" db:"company"`
	Position     string    `json:"position" db:"position"`
	HomePage     string    `json:"homePage" db:"home_page"`
	Interests    string    `json:"interests" db:"interests" gorm:"type:json"`
	Education    string    `json:"education" db:"education" gorm:"type:json"`
	WorkHistory  string    `json:"workHistory" db:"work_history" gorm:"type:json"`
	Certificates string    `json:"certificates" db:"certificates" gorm:"type:json"`
	Skills       string    `json:"skills" db:"skills" gorm:"type:json"`
	SocialLinks  string    `json:"socialLinks" db:"social_links" gorm:"type:json"`
	IsVerified   bool      `json:"isVerified" db:"is_verified" gorm:"default:false"`
	VerifiedAt   time.Time `json:"verifiedAt" db:"verified_at"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// UserAuth 用户认证
type UserAuth struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID      int64  `json:"userId" db:"user_id" gorm:"index;not null"`
	AuthType    string `json:"authType" db:"auth_type" gorm:"not null;comment:password,wechat,weibo,github,qq,alipay,etc"`
	AuthID      string `json:"authId" db:"auth_id" gorm:"uniqueIndex:idx_auth_type_id;not null"`
	Credential  string `json:"credential" db:"credential" gorm:"comment:密码哈希或token"`
	Verified    bool   `json:"verified" db:"verified" gorm:"default:false"`
	LastLoginAt time.Time `json:"lastLoginAt" db:"last_login_at"`
	ExpireAt    time.Time `json:"expireAt" db:"expire_at"`
	Extra       string `json:"extra" db:"extra" gorm:"type:json"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// UserLoginLog 用户登录日志
type UserLoginLog struct {
	common.BaseIDModel
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null"`
	LoginType   string    `json:"loginType" db:"login_type" gorm:"comment:password,wechat,weibo,github,qq,alipay,etc"`
	IP          string    `json:"ip" db:"ip"`
	UserAgent   string    `json:"userAgent" db:"user_agent"`
	Device      string    `json:"device" db:"device"`
	OS          string    `json:"os" db:"os"`
	Browser     string    `json:"browser" db:"browser"`
	Location    string    `json:"location" db:"location"`
	Status      int       `json:"status" db:"status" gorm:"comment:1-成功,0-失败"`
	ErrorMsg    string    `json:"errorMsg" db:"error_msg"`
	LoginAt     time.Time `json:"loginAt" db:"login_at" gorm:"autoCreateTime"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// UserMembership 用户会员
type UserMembership struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID        int64     `json:"userId" db:"user_id" gorm:"index;not null"`
	LevelID       int64     `json:"levelId" db:"level_id" gorm:"index;not null"`
	Status        int       `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-暂停,2-已过期"`
	StartDate     time.Time `json:"startDate" db:"start_date"`
	EndDate       time.Time `json:"endDate" db:"end_date"`
	OrderID       string    `json:"orderId" db:"order_id"`
	AutoRenew     bool      `json:"autoRenew" db:"auto_renew" gorm:"default:false"`
	RenewalAmount float64   `json:"renewalAmount" db:"renewal_amount"`
	LastRenewalAt time.Time `json:"lastRenewalAt" db:"last_renewal_at"`
	NextRenewalAt time.Time `json:"nextRenewalAt" db:"next_renewal_at"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// MembershipLevel 会员等级
type MembershipLevel struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string  `json:"name" db:"name" gorm:"not null"`
	Code        string  `json:"code" db:"code" gorm:"uniqueIndex;not null"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price" gorm:"not null"`
	PriceYear   float64 `json:"priceYear" db:"price_year"`
	Duration    int     `json:"duration" db:"duration" gorm:"comment:单位天"`
	Icon        string  `json:"icon" db:"icon"`
	Color       string  `json:"color" db:"color"`
	SortOrder   int     `json:"sortOrder" db:"sort_order" gorm:"default:0"`
	IsDefault   bool    `json:"isDefault" db:"is_default" gorm:"default:false"`
	Status      int     `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	TenantID    int64   `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// MembershipBenefit 会员权益
type MembershipBenefit struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null"`
	Description string `json:"description" db:"description"`
	Icon        string `json:"icon" db:"icon"`
	SortOrder   int    `json:"sortOrder" db:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// MembershipLevelBenefit 会员等级权益
type MembershipLevelBenefit struct {
	common.BaseIDModel
	LevelID   int64  `json:"levelId" db:"level_id" gorm:"index;not null"`
	BenefitID int64  `json:"benefitId" db:"benefit_id" gorm:"index;not null"`
	Value     string `json:"value" db:"value" gorm:"comment:权益值,如次数,容量等"`
	TenantID  int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}
