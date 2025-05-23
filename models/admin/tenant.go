package admin

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Tenant 租户模型
type Tenant struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string    `json:"name" db:"name" gorm:"not null"`
	Code        string    `json:"code" db:"code" gorm:"uniqueIndex;not null"`
	Type        int       `json:"type" db:"type" gorm:"comment:1-企业,2-个人,3-教育机构"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用,2-已过期"`
	Logo        string    `json:"logo" db:"logo"`
	Description string    `json:"description" db:"description"`
	OwnerId     int64     `json:"ownerId" db:"owner_id" gorm:"index"`
	ExpireAt    time.Time `json:"expireAt" db:"expire_at"`
}

// TenantQuota 租户配额
type TenantQuota struct {
	common.BaseIDModel
	common.BaseTimeModel
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;not null"`
	ResourceType string    `json:"resourceType" db:"resource_type" gorm:"not null;comment:user-用户数,storage-存储空间,api-API调用,bandwidth-带宽"`
	QuotaLimit   int64     `json:"quotaLimit" db:"quota_limit" gorm:"not null"`
	UsedAmount   int64     `json:"usedAmount" db:"used_amount" gorm:"default:0"`
	Unit         string    `json:"unit" db:"unit"`
	ResetCycle   string    `json:"resetCycle" db:"reset_cycle" gorm:"comment:none-不重置,daily-每日,monthly-每月,yearly-每年"`
	LastResetTime *time.Time `json:"lastResetTime" db:"last_reset_time"`
	NextResetTime *time.Time `json:"nextResetTime" db:"next_reset_time"`
}

// TenantBilling 租户计费
type TenantBilling struct {
	common.BaseIDModel
	common.BaseTimeModel
	TenantID      int64      `json:"tenantId" db:"tenant_id" gorm:"index;not null"`
	PlanID        int64      `json:"planId" db:"plan_id"`
	BillingCycle  string     `json:"billingCycle" db:"billing_cycle" gorm:"not null;comment:monthly-月付,quarterly-季付,yearly-年付"`
	Amount        float64    `json:"amount" db:"amount" gorm:"not null"`
	Currency      string     `json:"currency" db:"currency" gorm:"default:CNY"`
	Status        string     `json:"status" db:"status" gorm:"not null;comment:active-活跃,pending-待支付,expired-已过期,canceled-已取消"`
	PaymentMethod string     `json:"paymentMethod" db:"payment_method"`
	StartDate     time.Time  `json:"startDate" db:"start_date" gorm:"not null"`
	EndDate       time.Time  `json:"endDate" db:"end_date" gorm:"not null"`
	LastPaymentTime *time.Time `json:"lastPaymentTime" db:"last_payment_time"`
	NextPaymentTime *time.Time `json:"nextPaymentTime" db:"next_payment_time"`
	AutoRenew     bool       `json:"autoRenew" db:"auto_renew" gorm:"default:false"`
}

// TenantInvitation 租户邀请
type TenantInvitation struct {
	common.BaseIDModel
	common.BaseTimeModel
	TenantID       int64      `json:"tenantId" db:"tenant_id" gorm:"index;not null"`
	InviterID      int64      `json:"inviterId" db:"inviter_id" gorm:"not null"`
	Email          string     `json:"email" db:"email" gorm:"not null"`
	Role           string     `json:"role" db:"role" gorm:"not null"`
	InvitationCode string     `json:"invitationCode" db:"invitation_code" gorm:"uniqueIndex;not null"`
	Message        string     `json:"message" db:"message"`
	Status         string     `json:"status" db:"status" gorm:"default:pending;comment:pending-待接受,accepted-已接受,rejected-已拒绝,expired-已过期"`
	AcceptedUserID int64      `json:"acceptedUserId" db:"accepted_user_id"`
	ExpireTime     time.Time  `json:"expireTime" db:"expire_time" gorm:"not null"`
}

// TenantChangeLog 租户变更日志
type TenantChangeLog struct {
	common.BaseIDModel
	common.BaseTimeModel
	TenantID      int64  `json:"tenantId" db:"tenant_id" gorm:"index;not null"`
	ChangeType    string `json:"changeType" db:"change_type" gorm:"not null;comment:create-创建,update-更新,status-状态变更,expire-过期变更"`
	FieldName     string `json:"fieldName" db:"field_name"`
	OldValue      string `json:"oldValue" db:"old_value"`
	NewValue      string `json:"newValue" db:"new_value"`
	ChangeReason  string `json:"changeReason" db:"change_reason"`
	OperatorID    int64  `json:"operatorId" db:"operator_id"`
	OperatorName  string `json:"operatorName" db:"operator_name"`
	IP            string `json:"ip" db:"ip"`
	TenantIDCopy  int64  `json:"tenantIdCopy" db:"tenant_id_copy" gorm:"index;comment:租户ID副本用于跨租户查询"`
}
