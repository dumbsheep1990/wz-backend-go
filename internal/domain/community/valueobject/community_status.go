package valueobject

import "errors"

// CommunityStatus 社区状态值对象
type CommunityStatus int32

// 社区状态常量
const (
	CommunityStatusInactive  CommunityStatus = 0 // 未激活
	CommunityStatusActive    CommunityStatus = 1 // 正常活跃
	CommunityStatusSuspended CommunityStatus = 2 // 暂停
	CommunityStatusArchived  CommunityStatus = 3 // 已归档
	CommunityStatusDeleted   CommunityStatus = 4 // 已删除
	CommunityStatusReviewing CommunityStatus = 5 // 审核中
)

// NewCommunityStatus 创建社区状态值对象
func NewCommunityStatus(status int32) (CommunityStatus, error) {
	switch status {
	case int32(CommunityStatusInactive), int32(CommunityStatusActive), int32(CommunityStatusSuspended),
		int32(CommunityStatusArchived), int32(CommunityStatusDeleted), int32(CommunityStatusReviewing):
		return CommunityStatus(status), nil
	default:
		return 0, errors.New("无效的社区状态")
	}
}

// NewActiveCommunityStatus 创建活跃状态
func NewActiveCommunityStatus() CommunityStatus {
	return CommunityStatusActive
}

// NewInactiveCommunityStatus 创建未激活状态
func NewInactiveCommunityStatus() CommunityStatus {
	return CommunityStatusInactive
}

// NewReviewingCommunityStatus 创建审核中状态
func NewReviewingCommunityStatus() CommunityStatus {
	return CommunityStatusReviewing
}

// Value 获取状态值
func (s CommunityStatus) Value() int32 {
	return int32(s)
}

// String 状态的字符串表示
func (s CommunityStatus) String() string {
	switch s {
	case CommunityStatusInactive:
		return "未激活"
	case CommunityStatusActive:
		return "活跃"
	case CommunityStatusSuspended:
		return "暂停"
	case CommunityStatusArchived:
		return "已归档"
	case CommunityStatusDeleted:
		return "已删除"
	case CommunityStatusReviewing:
		return "审核中"
	default:
		return "未知状态"
	}
}

// EnglishString 状态的英文字符串表示
func (s CommunityStatus) EnglishString() string {
	switch s {
	case CommunityStatusInactive:
		return "inactive"
	case CommunityStatusActive:
		return "active"
	case CommunityStatusSuspended:
		return "suspended"
	case CommunityStatusArchived:
		return "archived"
	case CommunityStatusDeleted:
		return "deleted"
	case CommunityStatusReviewing:
		return "reviewing"
	default:
		return "unknown"
	}
}

// IsActive 判断社区是否处于活跃状态
func (s CommunityStatus) IsActive() bool {
	return s == CommunityStatusActive
}

// IsInactive 判断社区是否处于未激活状态
func (s CommunityStatus) IsInactive() bool {
	return s == CommunityStatusInactive
}

// IsSuspended 判断社区是否被暂停
func (s CommunityStatus) IsSuspended() bool {
	return s == CommunityStatusSuspended
}

// IsArchived 判断社区是否被归档
func (s CommunityStatus) IsArchived() bool {
	return s == CommunityStatusArchived
}

// IsDeleted 判断社区是否被删除
func (s CommunityStatus) IsDeleted() bool {
	return s == CommunityStatusDeleted
}

// IsReviewing 判断社区是否在审核中
func (s CommunityStatus) IsReviewing() bool {
	return s == CommunityStatusReviewing
}

// IsOperational 判断社区是否可以正常运营
func (s CommunityStatus) IsOperational() bool {
	return s == CommunityStatusActive
}

// CanAcceptMembers 判断社区是否可以接受新成员
func (s CommunityStatus) CanAcceptMembers() bool {
	return s == CommunityStatusActive
}

// CanCreateContent 判断社区是否可以创建内容
func (s CommunityStatus) CanCreateContent() bool {
	return s == CommunityStatusActive
}

// CanBeActivated 判断社区是否可以被激活
func (s CommunityStatus) CanBeActivated() bool {
	return s == CommunityStatusInactive || s == CommunityStatusSuspended
}

// CanBeSuspended 判断社区是否可以被暂停
func (s CommunityStatus) CanBeSuspended() bool {
	return s == CommunityStatusActive
}

// CanBeArchived 判断社区是否可以被归档
func (s CommunityStatus) CanBeArchived() bool {
	return s == CommunityStatusActive || s == CommunityStatusSuspended
}

// CanBeDeleted 判断社区是否可以被删除
func (s CommunityStatus) CanBeDeleted() bool {
	return s != CommunityStatusDeleted
}

// CanBeReviewed 判断社区是否可以进入审核状态
func (s CommunityStatus) CanBeReviewed() bool {
	return s == CommunityStatusInactive
}

// CanTransitionTo 检查当前状态是否可以转换到目标状态
func (s CommunityStatus) CanTransitionTo(target CommunityStatus) bool {
	switch s {
	case CommunityStatusInactive:
		return target == CommunityStatusActive || target == CommunityStatusReviewing || target == CommunityStatusDeleted
	case CommunityStatusActive:
		return target == CommunityStatusSuspended || target == CommunityStatusArchived || target == CommunityStatusDeleted
	case CommunityStatusSuspended:
		return target == CommunityStatusActive || target == CommunityStatusArchived || target == CommunityStatusDeleted
	case CommunityStatusArchived:
		return target == CommunityStatusActive || target == CommunityStatusDeleted
	case CommunityStatusReviewing:
		return target == CommunityStatusActive || target == CommunityStatusSuspended || target == CommunityStatusDeleted
	case CommunityStatusDeleted:
		return false // 已删除状态不能转换到其他状态
	default:
		return false
	}
}

// GetAllowedTransitions 获取当前状态允许转换的目标状态列表
func (s CommunityStatus) GetAllowedTransitions() []CommunityStatus {
	switch s {
	case CommunityStatusInactive:
		return []CommunityStatus{CommunityStatusActive, CommunityStatusReviewing, CommunityStatusDeleted}
	case CommunityStatusActive:
		return []CommunityStatus{CommunityStatusSuspended, CommunityStatusArchived, CommunityStatusDeleted}
	case CommunityStatusSuspended:
		return []CommunityStatus{CommunityStatusActive, CommunityStatusArchived, CommunityStatusDeleted}
	case CommunityStatusArchived:
		return []CommunityStatus{CommunityStatusActive, CommunityStatusDeleted}
	case CommunityStatusReviewing:
		return []CommunityStatus{CommunityStatusActive, CommunityStatusSuspended, CommunityStatusDeleted}
	case CommunityStatusDeleted:
		return []CommunityStatus{} // 已删除状态不能转换
	default:
		return []CommunityStatus{}
	}
}

// GetTransitionReason 获取状态转换的描述
func (s CommunityStatus) GetTransitionReason(target CommunityStatus) string {
	if !s.CanTransitionTo(target) {
		return "不允许的状态转换"
	}
	
	switch {
	case s == CommunityStatusInactive && target == CommunityStatusActive:
		return "激活社区"
	case s == CommunityStatusInactive && target == CommunityStatusReviewing:
		return "提交审核"
	case s == CommunityStatusActive && target == CommunityStatusSuspended:
		return "暂停社区"
	case s == CommunityStatusSuspended && target == CommunityStatusActive:
		return "恢复社区"
	case s == CommunityStatusActive && target == CommunityStatusArchived:
		return "归档社区"
	case s == CommunityStatusArchived && target == CommunityStatusActive:
		return "重新激活社区"
	case s == CommunityStatusReviewing && target == CommunityStatusActive:
		return "审核通过"
	case s == CommunityStatusReviewing && target == CommunityStatusSuspended:
		return "审核不通过"
	case target == CommunityStatusDeleted:
		return "删除社区"
	default:
		return "状态转换"
	}
}

// IsValid 检查状态是否有效
func (s CommunityStatus) IsValid() bool {
	return s >= CommunityStatusInactive && s <= CommunityStatusReviewing
}

// Equals 比较两个状态是否相等
func (s CommunityStatus) Equals(other CommunityStatus) bool {
	return s == other
} 