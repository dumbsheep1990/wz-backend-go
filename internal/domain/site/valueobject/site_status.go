package valueobject

import (
	"errors"
)

// SiteStatus 站点状态值对象
type SiteStatus struct {
	value string
}

// 站点状态常量
const (
	StatusDraft     = "draft"     // 草稿状态
	StatusPublished = "published" // 已发布
	StatusArchived  = "archived"  // 已归档
)

// NewSiteStatus 创建站点状态
func NewSiteStatus(value string) (SiteStatus, error) {
	if err := validateSiteStatus(value); err != nil {
		return SiteStatus{}, err
	}
	return SiteStatus{value: value}, nil
}

// NewDraftStatus 创建草稿状态
func NewDraftStatus() SiteStatus {
	return SiteStatus{value: StatusDraft}
}

// NewPublishedStatus 创建已发布状态
func NewPublishedStatus() SiteStatus {
	return SiteStatus{value: StatusPublished}
}

// NewArchivedStatus 创建已归档状态
func NewArchivedStatus() SiteStatus {
	return SiteStatus{value: StatusArchived}
}

// Value 获取状态值
func (s SiteStatus) Value() string {
	return s.value
}

// IsDraft 是否为草稿状态
func (s SiteStatus) IsDraft() bool {
	return s.value == StatusDraft
}

// IsPublished 是否为已发布状态
func (s SiteStatus) IsPublished() bool {
	return s.value == StatusPublished
}

// IsArchived 是否为已归档状态
func (s SiteStatus) IsArchived() bool {
	return s.value == StatusArchived
}

// CanPublish 是否可以发布
func (s SiteStatus) CanPublish() bool {
	return s.value == StatusDraft || s.value == StatusArchived
}

// CanArchive 是否可以归档
func (s SiteStatus) CanArchive() bool {
	return s.value == StatusPublished || s.value == StatusDraft
}

// Equals 判断两个状态是否相等
func (s SiteStatus) Equals(other SiteStatus) bool {
	return s.value == other.value
}

// String 字符串表示
func (s SiteStatus) String() string {
	return s.value
}

// validateSiteStatus 验证站点状态
func validateSiteStatus(value string) error {
	switch value {
	case StatusDraft, StatusPublished, StatusArchived:
		return nil
	default:
		return errors.New("无效的站点状态")
	}
} 