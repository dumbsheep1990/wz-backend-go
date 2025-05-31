package valueobject

import "errors"

// PageStatus 页面状态值对象
type PageStatus int32

// 页面状态常量
const (
	PageStatusDraft     PageStatus = 0 // 草稿
	PageStatusPublished PageStatus = 1 // 已发布
	PageStatusScheduled PageStatus = 2 // 定时发布
	PageStatusPrivate   PageStatus = 3 // 私有
	PageStatusArchived  PageStatus = 4 // 已归档
	PageStatusDeleted   PageStatus = 5 // 已删除
)

// NewPageStatus 创建页面状态值对象
func NewPageStatus(status int32) (PageStatus, error) {
	switch status {
	case int32(PageStatusDraft), int32(PageStatusPublished), int32(PageStatusScheduled),
		int32(PageStatusPrivate), int32(PageStatusArchived), int32(PageStatusDeleted):
		return PageStatus(status), nil
	default:
		return 0, errors.New("无效的页面状态")
	}
}

// NewDraftPageStatus 创建草稿状态
func NewDraftPageStatus() PageStatus {
	return PageStatusDraft
}

// NewPublishedPageStatus 创建已发布状态
func NewPublishedPageStatus() PageStatus {
	return PageStatusPublished
}

// NewPrivatePageStatus 创建私有状态
func NewPrivatePageStatus() PageStatus {
	return PageStatusPrivate
}

// Value 获取状态值
func (s PageStatus) Value() int32 {
	return int32(s)
}

// String 状态的字符串表示
func (s PageStatus) String() string {
	switch s {
	case PageStatusDraft:
		return "草稿"
	case PageStatusPublished:
		return "已发布"
	case PageStatusScheduled:
		return "定时发布"
	case PageStatusPrivate:
		return "私有"
	case PageStatusArchived:
		return "已归档"
	case PageStatusDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}

// EnglishString 状态的英文字符串表示
func (s PageStatus) EnglishString() string {
	switch s {
	case PageStatusDraft:
		return "draft"
	case PageStatusPublished:
		return "published"
	case PageStatusScheduled:
		return "scheduled"
	case PageStatusPrivate:
		return "private"
	case PageStatusArchived:
		return "archived"
	case PageStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsDraft 判断页面是否为草稿状态
func (s PageStatus) IsDraft() bool {
	return s == PageStatusDraft
}

// IsPublished 判断页面是否已发布
func (s PageStatus) IsPublished() bool {
	return s == PageStatusPublished
}

// IsScheduled 判断页面是否为定时发布状态
func (s PageStatus) IsScheduled() bool {
	return s == PageStatusScheduled
}

// IsPrivate 判断页面是否为私有状态
func (s PageStatus) IsPrivate() bool {
	return s == PageStatusPrivate
}

// IsArchived 判断页面是否已归档
func (s PageStatus) IsArchived() bool {
	return s == PageStatusArchived
}

// IsDeleted 判断页面是否已删除
func (s PageStatus) IsDeleted() bool {
	return s == PageStatusDeleted
}

// IsVisible 判断页面是否对外可见
func (s PageStatus) IsVisible() bool {
	return s == PageStatusPublished
}

// IsEditable 判断页面是否可编辑
func (s PageStatus) IsEditable() bool {
	return s == PageStatusDraft || s == PageStatusPrivate || s == PageStatusScheduled
}

// CanBePublished 判断页面是否可以发布
func (s PageStatus) CanBePublished() bool {
	return s == PageStatusDraft || s == PageStatusPrivate || s == PageStatusScheduled
}

// CanBeArchived 判断页面是否可以归档
func (s PageStatus) CanBeArchived() bool {
	return s == PageStatusPublished || s == PageStatusPrivate
}

// CanBeDeleted 判断页面是否可以删除
func (s PageStatus) CanBeDeleted() bool {
	return s != PageStatusDeleted
}

// CanBeRestored 判断页面是否可以恢复
func (s PageStatus) CanBeRestored() bool {
	return s == PageStatusArchived || s == PageStatusDeleted
}

// CanTransitionTo 检查当前状态是否可以转换到目标状态
func (s PageStatus) CanTransitionTo(target PageStatus) bool {
	switch s {
	case PageStatusDraft:
		return target == PageStatusPublished || target == PageStatusPrivate || 
			   target == PageStatusScheduled || target == PageStatusDeleted
	case PageStatusPublished:
		return target == PageStatusDraft || target == PageStatusPrivate || 
			   target == PageStatusArchived || target == PageStatusDeleted
	case PageStatusScheduled:
		return target == PageStatusDraft || target == PageStatusPublished || 
			   target == PageStatusPrivate || target == PageStatusDeleted
	case PageStatusPrivate:
		return target == PageStatusDraft || target == PageStatusPublished || 
			   target == PageStatusScheduled || target == PageStatusDeleted
	case PageStatusArchived:
		return target == PageStatusDraft || target == PageStatusPublished || 
			   target == PageStatusPrivate || target == PageStatusDeleted
	case PageStatusDeleted:
		return target == PageStatusDraft || target == PageStatusPrivate
	default:
		return false
	}
}

// GetAllowedTransitions 获取当前状态允许转换的目标状态列表
func (s PageStatus) GetAllowedTransitions() []PageStatus {
	switch s {
	case PageStatusDraft:
		return []PageStatus{PageStatusPublished, PageStatusPrivate, PageStatusScheduled, PageStatusDeleted}
	case PageStatusPublished:
		return []PageStatus{PageStatusDraft, PageStatusPrivate, PageStatusArchived, PageStatusDeleted}
	case PageStatusScheduled:
		return []PageStatus{PageStatusDraft, PageStatusPublished, PageStatusPrivate, PageStatusDeleted}
	case PageStatusPrivate:
		return []PageStatus{PageStatusDraft, PageStatusPublished, PageStatusScheduled, PageStatusDeleted}
	case PageStatusArchived:
		return []PageStatus{PageStatusDraft, PageStatusPublished, PageStatusPrivate, PageStatusDeleted}
	case PageStatusDeleted:
		return []PageStatus{PageStatusDraft, PageStatusPrivate}
	default:
		return []PageStatus{}
	}
}

// GetTransitionReason 获取状态转换的描述
func (s PageStatus) GetTransitionReason(target PageStatus) string {
	if !s.CanTransitionTo(target) {
		return "不允许的状态转换"
	}
	
	switch {
	case s == PageStatusDraft && target == PageStatusPublished:
		return "发布页面"
	case s == PageStatusDraft && target == PageStatusPrivate:
		return "设为私有"
	case s == PageStatusDraft && target == PageStatusScheduled:
		return "定时发布"
	case s == PageStatusPublished && target == PageStatusDraft:
		return "撤回到草稿"
	case s == PageStatusPublished && target == PageStatusArchived:
		return "归档页面"
	case s == PageStatusArchived && target == PageStatusPublished:
		return "重新发布"
	case s == PageStatusScheduled && target == PageStatusPublished:
		return "立即发布"
	case target == PageStatusDeleted:
		return "删除页面"
	case s == PageStatusDeleted && target == PageStatusDraft:
		return "恢复为草稿"
	case s == PageStatusDeleted && target == PageStatusPrivate:
		return "恢复为私有"
	default:
		return "状态转换"
	}
}

// GetStatusDescription 获取状态的详细描述
func (s PageStatus) GetStatusDescription() string {
	switch s {
	case PageStatusDraft:
		return "页面尚未发布，仅创建者可见"
	case PageStatusPublished:
		return "页面已发布，所有访客可见"
	case PageStatusScheduled:
		return "页面已安排定时发布"
	case PageStatusPrivate:
		return "页面为私有状态，仅授权用户可见"
	case PageStatusArchived:
		return "页面已归档，不再对外显示"
	case PageStatusDeleted:
		return "页面已删除，可以恢复"
	default:
		return "未知状态"
	}
}

// GetStatusColor 获取状态对应的颜色（用于UI显示）
func (s PageStatus) GetStatusColor() string {
	switch s {
	case PageStatusDraft:
		return "gray"
	case PageStatusPublished:
		return "green"
	case PageStatusScheduled:
		return "blue"
	case PageStatusPrivate:
		return "orange"
	case PageStatusArchived:
		return "purple"
	case PageStatusDeleted:
		return "red"
	default:
		return "gray"
	}
}

// GetStatusIcon 获取状态对应的图标
func (s PageStatus) GetStatusIcon() string {
	switch s {
	case PageStatusDraft:
		return "draft"
	case PageStatusPublished:
		return "public"
	case PageStatusScheduled:
		return "schedule"
	case PageStatusPrivate:
		return "lock"
	case PageStatusArchived:
		return "archive"
	case PageStatusDeleted:
		return "delete"
	default:
		return "help"
	}
}

// IsValid 检查状态是否有效
func (s PageStatus) IsValid() bool {
	return s >= PageStatusDraft && s <= PageStatusDeleted
}

// Equals 比较两个状态是否相等
func (s PageStatus) Equals(other PageStatus) bool {
	return s == other
} 