package dto

import (
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// SystemConfigDTO 系统配置DTO
type SystemConfigDTO struct {
	ID                int64     `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt,omitempty"`
	SystemConfig      string    `json:"systemConfig"`
	LogoUrl           string    `json:"logoUrl"`
	ApiUrl            string    `json:"apiUrl"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	ContentUrl        string    `json:"contentUrl"`
	UploadFileSize    int64     `json:"uploadFileSize"`
	EmailFrom         string    `json:"emailFrom"`
	EmailHost         string    `json:"emailHost"`
	EmailPort         int       `json:"emailPort"`
	EmailIsSSL        bool      `json:"emailIsSSL"`
	EmailSecret       string    `json:"emailSecret"`
	EmailNickname     string    `json:"emailNickname"`
	LdapOpen          bool      `json:"ldapOpen"`
	LdapHost          string    `json:"ldapHost"`
	LdapPort          int       `json:"ldapPort"`
	LdapBindDn        string    `json:"ldapBindDn"`
	LdapPassword      string    `json:"ldapPassword"`
	LdapBaseDn        string    `json:"ldapBaseDn"`
	LdapUserField     string    `json:"ldapUserField"`
	LdapNickName      string    `json:"ldapNickName"`
	LdapEmails        string    `json:"ldapEmails"`
	LdapTls           bool      `json:"ldapTls"`
	DingTalkOpen      bool      `json:"dingTalkOpen"`
	DingTalkAppKey    string    `json:"dingTalkAppKey"`
	DingTalkAppSecret string    `json:"dingTalkAppSecret"`
	DingTalkAgentId   string    `json:"dingTalkAgentId"`
}

// SystemConfigRequest 系统配置请求
type SystemConfigRequest struct {
	SystemConfig      string `json:"systemConfig"`
	LogoUrl           string `json:"logoUrl"`
	ApiUrl            string `json:"apiUrl"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ContentUrl        string `json:"contentUrl"`
	UploadFileSize    int64  `json:"uploadFileSize"`
	EmailFrom         string `json:"emailFrom"`
	EmailHost         string `json:"emailHost"`
	EmailPort         int    `json:"emailPort"`
	EmailIsSSL        bool   `json:"emailIsSSL"`
	EmailSecret       string `json:"emailSecret"`
	EmailNickname     string `json:"emailNickname"`
	LdapOpen          bool   `json:"ldapOpen"`
	LdapHost          string `json:"ldapHost"`
	LdapPort          int    `json:"ldapPort"`
	LdapBindDn        string `json:"ldapBindDn"`
	LdapPassword      string `json:"ldapPassword"`
	LdapBaseDn        string `json:"ldapBaseDn"`
	LdapUserField     string `json:"ldapUserField"`
	LdapNickName      string `json:"ldapNickName"`
	LdapEmails        string `json:"ldapEmails"`
	LdapTls           bool   `json:"ldapTls"`
	DingTalkOpen      bool   `json:"dingTalkOpen"`
	DingTalkAppKey    string `json:"dingTalkAppKey"`
	DingTalkAppSecret string `json:"dingTalkAppSecret"`
	DingTalkAgentId   string `json:"dingTalkAgentId"`
}

// ServerInfoDTO 服务器信息DTO
type ServerInfoDTO struct {
	Os    OsInfoDTO    `json:"os"`
	Cpu   CpuInfoDTO   `json:"cpu"`
	Ram   RamInfoDTO   `json:"ram"`
	Disk  DiskInfoDTO  `json:"disk"`
	Go    GoInfoDTO    `json:"go"`
	Db    DbInfoDTO    `json:"db"`
}

// OsInfoDTO 操作系统信息DTO
type OsInfoDTO struct {
	Name          string `json:"name"`
	Platform      string `json:"platform"`
	Family        string `json:"family"`
	Version       string `json:"version"`
	KernelArch    string `json:"kernelArch"`
	KernelVersion string `json:"kernelVersion"`
}

// CpuInfoDTO CPU信息DTO
type CpuInfoDTO struct {
	Model       string  `json:"model"`
	Cores       int     `json:"cores"`
	UsedPercent float64 `json:"usedPercent"`
}

// RamInfoDTO 内存信息DTO
type RamInfoDTO struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

// DiskInfoDTO 磁盘信息DTO
type DiskInfoDTO struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

// GoInfoDTO Go运行时信息DTO
type GoInfoDTO struct {
	Version      string `json:"version"`
	NumGoroutine int    `json:"numGoroutine"`
	NumCpu       int    `json:"numCpu"`
}

// DbInfoDTO 数据库信息DTO
type DbInfoDTO struct {
	Type        string `json:"type"`
	Version     string `json:"version"`
	DbName      string `json:"dbName"`
	MaxOpenConn int    `json:"maxOpenConn"`
	MaxIdleConn int    `json:"maxIdleConn"`
}

// SysDictionaryDTO 系统字典DTO
type SysDictionaryDTO struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
}

// SysDictionaryCreateRequest 创建系统字典请求
type SysDictionaryCreateRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// SysDictionaryUpdateRequest 更新系统字典请求
type SysDictionaryUpdateRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// SysDictionaryQueryRequest 查询系统字典请求
type SysDictionaryQueryRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      *int   `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

// SysDictionaryDetailDTO 系统字典详情DTO
type SysDictionaryDetailDTO struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	Status     int       `json:"status"`
	Sort       int       `json:"sort"`
	DictTypeID int64     `json:"sysDictionaryID"`
}

// SysDictionaryDetailCreateRequest 创建系统字典详情请求
type SysDictionaryDetailCreateRequest struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	Status     int    `json:"status"`
	Sort       int    `json:"sort"`
	DictTypeID int64  `json:"sysDictionaryID"`
}

// SysDictionaryDetailUpdateRequest 更新系统字典详情请求
type SysDictionaryDetailUpdateRequest struct {
	ID         int64  `json:"id"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Status     int    `json:"status"`
	Sort       int    `json:"sort"`
	DictTypeID int64  `json:"sysDictionaryID"`
}

// SysDictionaryDetailQueryRequest 查询系统字典详情请求
type SysDictionaryDetailQueryRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	DictTypeID int64  `json:"sysDictionaryID"`
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
	Status     *int   `json:"status,omitempty"`
}

// SysOperationRecordDTO 操作日志DTO
type SysOperationRecordDTO struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	IP           string    `json:"ip"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Status       int       `json:"status"`
	Latency      int64     `json:"latency"`
	Agent        string    `json:"agent"`
	ErrorMessage string    `json:"errorMessage"`
	Body         string    `json:"body"`
	UserID       int64     `json:"userId"`
	User         string    `json:"user"`
	Resp         string    `json:"resp"`
}

// SysOperationRecordQueryRequest 查询操作日志请求
type SysOperationRecordQueryRequest struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	Method        string `json:"method,omitempty"`
	Path          string `json:"path,omitempty"`
	Status        *int   `json:"status,omitempty"`
	UserID        *int64 `json:"userId,omitempty"`
	CreatedAfter  string `json:"createdAfter,omitempty"`
	CreatedBefore string `json:"createdBefore,omitempty"`
}

// SysParamsDTO 系统参数DTO
type SysParamsDTO struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	ParamName   string    `json:"paramName"`
	ParamKey    string    `json:"paramKey"`
	ParamValue  string    `json:"paramValue"`
	ParamType   string    `json:"paramType"`
	ParamDesc   string    `json:"paramDesc"`
	ParamGroup  string    `json:"paramGroup"`
	Status      int       `json:"status"`
}

// SysParamsCreateRequest 创建系统参数请求
type SysParamsCreateRequest struct {
	ParamName   string `json:"paramName"`
	ParamKey    string `json:"paramKey"`
	ParamValue  string `json:"paramValue"`
	ParamType   string `json:"paramType"`
	ParamDesc   string `json:"paramDesc"`
	ParamGroup  string `json:"paramGroup"`
	Status      int    `json:"status"`
}

// SysParamsUpdateRequest 更新系统参数请求
type SysParamsUpdateRequest struct {
	ID          int64  `json:"id"`
	ParamName   string `json:"paramName"`
	ParamKey    string `json:"paramKey"`
	ParamValue  string `json:"paramValue"`
	ParamType   string `json:"paramType"`
	ParamDesc   string `json:"paramDesc"`
	ParamGroup  string `json:"paramGroup"`
	Status      int    `json:"status"`
}

// SysParamsQueryRequest 查询系统参数请求
type SysParamsQueryRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	ParamName  string `json:"paramName,omitempty"`
	ParamKey   string `json:"paramKey,omitempty"`
	ParamValue string `json:"paramValue,omitempty"`
	ParamType  string `json:"paramType,omitempty"`
	ParamGroup string `json:"paramGroup,omitempty"`
	Status     *int   `json:"status,omitempty"`
}
