package admin

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Admin 管理员模型
type Admin struct {
	common.BaseIDModel
	common.BaseTimeModel
	Username  string    `json:"username" db:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" db:"password" gorm:"not null"`
	Role      string    `json:"role" db:"role" gorm:"not null"`
	Status    int       `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	LastLogin time.Time `json:"lastLogin" db:"last_login"`
}

// SysApi API模型
type SysApi struct {
	common.BaseIDModel
	common.BaseTimeModel
	Path        string `json:"path" db:"path" gorm:"uniqueIndex:idx_api_path_method;not null"`
	Description string `json:"description" db:"description"`
	ApiGroup    string `json:"apiGroup" db:"api_group" gorm:"index"`
	Method      string `json:"method" db:"method" gorm:"uniqueIndex:idx_api_path_method;not null"`
}

// SysOperationRecord 操作日志模型
type SysOperationRecord struct {
	common.BaseIDModel
	common.BaseTimeModel
	Ip           string `json:"ip" db:"ip"`
	Method       string `json:"method" db:"method"`
	Path         string `json:"path" db:"path"`
	Status       int    `json:"status" db:"status"`
	Latency      int64  `json:"latency" db:"latency" gorm:"comment:请求耗时(ms)"`
	Agent        string `json:"agent" db:"agent"`
	ErrorMessage string `json:"errorMessage" db:"error_message"`
	Body         string `json:"body" db:"body" gorm:"type:text"`
	UserId       int64  `json:"userId" db:"user_id" gorm:"index"`
	User         string `json:"user" db:"user"`
	Resp         string `json:"resp" db:"resp" gorm:"type:text"`
}

// SysDictionary 字典模型
type SysDictionary struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"uniqueIndex:idx_dict_name_type;not null"`
	Type        string `json:"type" db:"type" gorm:"uniqueIndex:idx_dict_name_type;not null"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	Description string `json:"description" db:"description"`
}

// SysDictionaryDetail 字典详情模型
type SysDictionaryDetail struct {
	common.BaseIDModel
	common.BaseTimeModel
	Label      string `json:"label" db:"label" gorm:"not null"`
	Value      string `json:"value" db:"value" gorm:"not null"`
	Status     int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	Sort       int    `json:"sort" db:"sort" gorm:"default:0"`
	DictTypeId int64  `json:"sysDictionaryID" db:"dict_type_id" gorm:"index;not null"`
}

// SysParams 系统参数模型
type SysParams struct {
	common.BaseIDModel
	common.BaseTimeModel
	ParamName  string `json:"paramName" db:"param_name" gorm:"not null"`
	ParamKey   string `json:"paramKey" db:"param_key" gorm:"uniqueIndex;not null"`
	ParamValue string `json:"paramValue" db:"param_value" gorm:"not null"`
	ParamType  string `json:"paramType" db:"param_type" gorm:"comment:string,number,boolean,json"`
	ParamDesc  string `json:"paramDesc" db:"param_desc"`
	ParamGroup string `json:"paramGroup" db:"param_group" gorm:"index"`
	Status     int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
}
