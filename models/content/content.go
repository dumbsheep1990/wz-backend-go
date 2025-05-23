package content

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Content 内容模型
type Content struct {
	common.BaseIDModel
	common.BaseTimeModel
	Title        string `json:"title" db:"title" gorm:"not null"`
	UserId       int64  `json:"userId" db:"user_id" gorm:"index;not null"`
	CategoryId   int64  `json:"categoryId" db:"category_id" gorm:"index"`
	Tags         string `json:"tags" db:"tags" gorm:"type:json"`
	Cover        string `json:"cover" db:"cover"`
	Summary      string `json:"summary" db:"summary"`
	Content      string `json:"content" db:"content" gorm:"type:longtext"`
	ContentType  string `json:"contentType" db:"content_type" gorm:"default:html;comment:html,markdown,rich"`
	Status       int    `json:"status" db:"status" gorm:"default:0;comment:0-草稿,1-已发布,2-审核中,-1-已删除"`
	ViewCount    int    `json:"viewCount" db:"view_count" gorm:"default:0"`
	LikeCount    int    `json:"likeCount" db:"like_count" gorm:"default:0"`
	CommentCount int    `json:"commentCount" db:"comment_count" gorm:"default:0"`
	PublishedAt  *time.Time `json:"publishedAt" db:"published_at"`
	TenantID     int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentCategory 内容分类
type ContentCategory struct {
	common.BaseIDModel
	common.BaseTimeModel
	ParentId    int64  `json:"parentId" db:"parent_id" gorm:"index;default:0"`
	Name        string `json:"name" db:"name" gorm:"not null"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex"`
	Description string `json:"description" db:"description"`
	Icon        string `json:"icon" db:"icon"`
	SortOrder   int    `json:"sortOrder" db:"sort_order" gorm:"default:0"`
	IsVisible   bool   `json:"isVisible" db:"is_visible" gorm:"default:true"`
	Path        string `json:"path" db:"path" gorm:"comment:层级路径,如1,2,3"`
	Level       int    `json:"level" db:"level" gorm:"default:1;comment:层级,1为顶级"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentTag 内容标签
type ContentTag struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name       string `json:"name" db:"name" gorm:"uniqueIndex:idx_tenant_name;not null"`
	Code       string `json:"code" db:"code" gorm:"uniqueIndex:idx_tenant_code"`
	Color      string `json:"color" db:"color"`
	Icon       string `json:"icon" db:"icon"`
	UseCount   int    `json:"useCount" db:"use_count" gorm:"default:0"`
	IsSystem   bool   `json:"isSystem" db:"is_system" gorm:"default:false"`
	CategoryId int64  `json:"categoryId" db:"category_id" gorm:"index"`
	TenantID   int64  `json:"tenantId" db:"tenant_id" gorm:"uniqueIndex:idx_tenant_name;uniqueIndex:idx_tenant_code;index"`
}

// ContentTagRelation 内容标签关联
type ContentTagRelation struct {
	common.BaseIDModel
	ContentId int64     `json:"contentId" db:"content_id" gorm:"index;not null"`
	TagId     int64     `json:"tagId" db:"tag_id" gorm:"index;not null"`
	CreatedAt time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID  int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentComment 内容评论
type ContentComment struct {
	common.BaseIDModel
	common.BaseTimeModel
	ContentId   int64  `json:"contentId" db:"content_id" gorm:"index;not null"`
	UserId      int64  `json:"userId" db:"user_id" gorm:"index;not null"`
	ParentId    int64  `json:"parentId" db:"parent_id" gorm:"index;default:0"`
	Content     string `json:"content" db:"content" gorm:"not null"`
	Status      int    `json:"status" db:"status" gorm:"default:0;comment:0-待审核,1-已通过,-1-已拒绝"`
	LikeCount   int    `json:"likeCount" db:"like_count" gorm:"default:0"`
	IP          string `json:"ip" db:"ip"`
	UserAgent   string `json:"userAgent" db:"user_agent"`
	IsAnonymous bool   `json:"isAnonymous" db:"is_anonymous" gorm:"default:false"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentVersion 内容版本
type ContentVersion struct {
	common.BaseIDModel
	ContentId     int64     `json:"contentId" db:"content_id" gorm:"index;not null"`
	Version       int       `json:"version" db:"version" gorm:"not null"`
	Title         string    `json:"title" db:"title"`
	Content       string    `json:"content" db:"content" gorm:"type:longtext"`
	ChangeLogs    string    `json:"changeLogs" db:"change_logs" gorm:"comment:变更记录"`
	CreatedBy     int64     `json:"createdBy" db:"created_by" gorm:"index"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentCollection 内容合集
type ContentCollection struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null"`
	Description string `json:"description" db:"description"`
	Cover       string `json:"cover" db:"cover"`
	UserId      int64  `json:"userId" db:"user_id" gorm:"index;not null"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:1-正常,0-禁用"`
	SortType    string `json:"sortType" db:"sort_type" gorm:"default:manual;comment:manual-手动,time-时间,hot-热度"`
	IsPublic    bool   `json:"isPublic" db:"is_public" gorm:"default:true"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentCollectionItem 内容合集项
type ContentCollectionItem struct {
	common.BaseIDModel
	CollectionId int64     `json:"collectionId" db:"collection_id" gorm:"index;not null"`
	ContentId    int64     `json:"contentId" db:"content_id" gorm:"index;not null"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0"`
	Description  string    `json:"description" db:"description"`
	AddedAt      time.Time `json:"addedAt" db:"added_at" gorm:"autoCreateTime"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index"`
}

// ContentAttachment 内容附件
type ContentAttachment struct {
	common.BaseIDModel
	common.BaseTimeModel
	ContentId   int64  `json:"contentId" db:"content_id" gorm:"index;not null"`
	Name        string `json:"name" db:"name" gorm:"not null"`
	FileId      int64  `json:"fileId" db:"file_id" gorm:"index;not null"`
	FileType    string `json:"fileType" db:"file_type"`
	FileSize    int64  `json:"fileSize" db:"file_size"`
	DownloadUrl string `json:"downloadUrl" db:"download_url"`
	SortOrder   int    `json:"sortOrder" db:"sort_order" gorm:"default:0"`
	DownloadCount int  `json:"downloadCount" db:"download_count" gorm:"default:0"`
	NeedLogin   bool   `json:"needLogin" db:"need_login" gorm:"default:false"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index"`
}
