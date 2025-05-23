package common

import "time"

// BaseModel 基础模型接口
type BaseModel interface {
	GetID() int64
	IsDeleted() bool
}

// BaseIDModel 基础ID模型
type BaseIDModel struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
}

// GetID 获取ID
func (m BaseIDModel) GetID() int64 {
	return m.ID
}

// BaseTimeModel 基础时间模型
type BaseTimeModel struct {
	CreatedAt time.Time  `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at" gorm:"index"`
}

// IsDeleted 是否已删除
func (m BaseTimeModel) IsDeleted() bool {
	return m.DeletedAt != nil
}

// BaseTenantModel 基础租户模型
type BaseTenantModel struct {
	TenantID int64 `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// CommonResponse 通用响应结构
type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageQuery 分页查询参数
type PageQuery struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	OrderDir string `json:"orderDir" form:"orderDir"`
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
