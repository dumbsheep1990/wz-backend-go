package file

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// File 文件模型
type File struct {
	common.BaseIDModel
	common.BaseTimeModel
	FileName     string    `json:"fileName" db:"file_name" gorm:"not null;comment:文件名称"`
	FileType     string    `json:"fileType" db:"file_type" gorm:"comment:文件类型"`
	MimeType     string    `json:"mimeType" db:"mime_type" gorm:"comment:MIME类型"`
	FileSize     int64     `json:"fileSize" db:"file_size" gorm:"comment:文件大小(bytes)"`
	FilePath     string    `json:"filePath" db:"file_path" gorm:"comment:文件路径"`
	StorageType  string    `json:"storageType" db:"storage_type" gorm:"comment:存储类型:local,aliyun,qiniu,tencent"`
	StoragePath  string    `json:"storagePath" db:"storage_path" gorm:"comment:存储路径"`
	Url          string    `json:"url" db:"url" gorm:"comment:访问URL"`
	Md5          string    `json:"md5" db:"md5" gorm:"index;comment:MD5哈希值"`
	CategoryID   int64     `json:"categoryId" db:"category_id" gorm:"index;comment:分类ID"`
	Tags         string    `json:"tags" db:"tags" gorm:"type:json;comment:标签,JSON数组"`
	UploadIP     string    `json:"uploadIp" db:"upload_ip" gorm:"comment:上传IP"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;comment:上传用户ID"`
	Status       int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用,-1-已删除"`
	AccessCount  int       `json:"accessCount" db:"access_count" gorm:"default:0;comment:访问次数"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// FileCategory 文件分类
type FileCategory struct {
	common.BaseIDModel
	common.BaseTimeModel
	ParentID     int64     `json:"parentId" db:"parent_id" gorm:"index;default:0;comment:父分类ID"`
	Name         string    `json:"name" db:"name" gorm:"not null;comment:分类名称"`
	Code         string    `json:"code" db:"code" gorm:"uniqueIndex:idx_tenant_code;not null;comment:分类编码"`
	Description  string    `json:"description" db:"description" gorm:"comment:描述"`
	AllowedTypes string    `json:"allowedTypes" db:"allowed_types" gorm:"type:json;comment:允许的文件类型,JSON数组"`
	MaxSize      int64     `json:"maxSize" db:"max_size" gorm:"comment:最大文件大小(bytes)"`
	Path         string    `json:"path" db:"path" gorm:"comment:分类路径,如1,2,3"`
	Level        int       `json:"level" db:"level" gorm:"default:1;comment:层级,1为顶级"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	IsVisible    bool      `json:"isVisible" db:"is_visible" gorm:"default:true;comment:是否可见"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"uniqueIndex:idx_tenant_code;index;comment:租户ID"`
}

// FileAccessLog 文件访问日志
type FileAccessLog struct {
	common.BaseIDModel
	FileID       int64     `json:"fileId" db:"file_id" gorm:"index;not null;comment:文件ID"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;comment:用户ID"`
	AccessType   string    `json:"accessType" db:"access_type" gorm:"comment:访问类型:view,download,preview"`
	AccessIP     string    `json:"accessIp" db:"access_ip" gorm:"comment:访问IP"`
	UserAgent    string    `json:"userAgent" db:"user_agent" gorm:"comment:用户代理"`
	RefererUrl   string    `json:"refererUrl" db:"referer_url" gorm:"comment:来源URL"`
	AccessTime   time.Time `json:"accessTime" db:"access_time" gorm:"autoCreateTime;comment:访问时间"`
	ResponseTime int       `json:"responseTime" db:"response_time" gorm:"comment:响应时间(ms)"`
	Status       int       `json:"status" db:"status" gorm:"comment:状态码"`
	ErrorMsg     string    `json:"errorMsg" db:"error_msg" gorm:"comment:错误信息"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// FileChunk 文件分片
type FileChunk struct {
	common.BaseIDModel
	common.BaseTimeModel
	FileID       int64     `json:"fileId" db:"file_id" gorm:"index;not null;comment:文件ID"`
	ChunkNumber  int       `json:"chunkNumber" db:"chunk_number" gorm:"not null;comment:分片序号"`
	ChunkSize    int64     `json:"chunkSize" db:"chunk_size" gorm:"not null;comment:分片大小"`
	ChunkPath    string    `json:"chunkPath" db:"chunk_path" gorm:"not null;comment:分片路径"`
	Status       string    `json:"status" db:"status" gorm:"default:uploading;comment:状态:uploading-上传中,uploaded-已上传,merged-已合并"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// FileProcessingTask 文件处理任务
type FileProcessingTask struct {
	common.BaseIDModel
	common.BaseTimeModel
	FileID        int64     `json:"fileId" db:"file_id" gorm:"index;not null;comment:文件ID"`
	TaskType      string    `json:"taskType" db:"task_type" gorm:"not null;comment:任务类型:compress-压缩,extract-解压,convert-转换,analyze-分析,resize-调整大小"`
	Status        string    `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-待处理,processing-处理中,completed-已完成,failed-失败"`
	Parameters    string    `json:"parameters" db:"parameters" gorm:"type:json;comment:处理参数"`
	Result        string    `json:"result" db:"result" gorm:"type:json;comment:处理结果"`
	ErrorMessage  string    `json:"errorMessage" db:"error_message" gorm:"comment:错误信息"`
	OutputFileID  int64     `json:"outputFileId" db:"output_file_id" gorm:"index;comment:输出文件ID"`
	Progress      float64   `json:"progress" db:"progress" gorm:"default:0;comment:进度"`
	CreatedBy     int64     `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	StartedAt     *time.Time `json:"startedAt" db:"started_at" gorm:"comment:开始时间"`
	CompletedAt   *time.Time `json:"completedAt" db:"completed_at" gorm:"comment:完成时间"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// FileShare 文件共享
type FileShare struct {
	common.BaseIDModel
	common.BaseTimeModel
	FileID        int64     `json:"fileId" db:"file_id" gorm:"index;not null;comment:文件ID"`
	CreatorID     int64     `json:"creatorId" db:"creator_id" gorm:"index;not null;comment:创建人ID"`
	ShareCode     string    `json:"shareCode" db:"share_code" gorm:"uniqueIndex;not null;comment:分享码"`
	ShareType     string    `json:"shareType" db:"share_type" gorm:"default:link;comment:分享类型:link-链接,email-邮件"`
	AccessType    string    `json:"accessType" db:"access_type" gorm:"default:view;comment:访问类型:view-查看,download-下载,edit-编辑"`
	Password      string    `json:"password" db:"password" gorm:"comment:访问密码"`
	ExpireTime    *time.Time `json:"expireTime" db:"expire_time" gorm:"comment:过期时间"`
	MaxViews      int       `json:"maxViews" db:"max_views" gorm:"comment:最大查看次数"`
	ViewCount     int       `json:"viewCount" db:"view_count" gorm:"default:0;comment:已查看次数"`
	DownloadCount int       `json:"downloadCount" db:"download_count" gorm:"default:0;comment:下载次数"`
	Status        string    `json:"status" db:"status" gorm:"default:active;comment:状态:active-活跃,expired-已过期,disabled-已禁用"`
	Recipients    string    `json:"recipients" db:"recipients" gorm:"type:json;comment:接收人列表"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
