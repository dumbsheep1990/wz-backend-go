package valueobject

import "errors"

// UserStatus 用户状态值对象
type UserStatus int32

// 用户状态常量
const (
	UserStatusInactive   UserStatus = 0 // 未激活
	UserStatusActive     UserStatus = 1 // 正常
	UserStatusSuspended  UserStatus = 2 // 暂停
	UserStatusBanned     UserStatus = 3 // 封禁
	UserStatusDeleted    UserStatus = 4 // 删除
	UserStatusLocked     UserStatus = 5 // 锁定
)

// NewUserStatus 创建用户状态值对象
func NewUserStatus(status int32) (UserStatus, error) {
	switch status {
	case int32(UserStatusInactive), int32(UserStatusActive), int32(UserStatusSuspended),
		int32(UserStatusBanned), int32(UserStatusDeleted), int32(UserStatusLocked):
		return UserStatus(status), nil
	default:
		return 0, errors.New("无效的用户状态")
	}
}

// Value 获取状态值
func (s UserStatus) Value() int32 {
	return int32(s)
}

// String 状态的字符串表示
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "未激活"
	case UserStatusActive:
		return "正常"
	case UserStatusSuspended:
		return "暂停"
	case UserStatusBanned:
		return "封禁"
	case UserStatusDeleted:
		return "删除"
	case UserStatusLocked:
		return "锁定"
	default:
		return "未知状态"
	}
}

// IsActive 判断用户是否为活跃状态
func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

// IsInactive 判断用户是否为未激活状态
func (s UserStatus) IsInactive() bool {
	return s == UserStatusInactive
}

// IsSuspended 判断用户是否被暂停
func (s UserStatus) IsSuspended() bool {
	return s == UserStatusSuspended
}

// IsBanned 判断用户是否被封禁
func (s UserStatus) IsBanned() bool {
	return s == UserStatusBanned
}

// IsDeleted 判断用户是否被删除
func (s UserStatus) IsDeleted() bool {
	return s == UserStatusDeleted
}

// IsLocked 判断用户是否被锁定
func (s UserStatus) IsLocked() bool {
	return s == UserStatusLocked
}

// CanLogin 判断用户是否可以登录
func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive
}

// CanBeActivated 判断用户是否可以被激活
func (s UserStatus) CanBeActivated() bool {
	return s == UserStatusInactive
}

// CanBeSuspended 判断用户是否可以被暂停
func (s UserStatus) CanBeSuspended() bool {
	return s == UserStatusActive
}

// CanBeUnsuspended 判断用户是否可以被解除暂停
func (s UserStatus) CanBeUnsuspended() bool {
	return s == UserStatusSuspended
}

// CanBeBanned 判断用户是否可以被封禁
func (s UserStatus) CanBeBanned() bool {
	return s == UserStatusActive || s == UserStatusSuspended
}

// CanBeUnbanned 判断用户是否可以被解封
func (s UserStatus) CanBeUnbanned() bool {
	return s == UserStatusBanned
}

// CanBeLocked 判断用户是否可以被锁定
func (s UserStatus) CanBeLocked() bool {
	return s == UserStatusActive
}

// CanBeUnlocked 判断用户是否可以被解锁
func (s UserStatus) CanBeUnlocked() bool {
	return s == UserStatusLocked
}

// CanBeDeleted 判断用户是否可以被删除
func (s UserStatus) CanBeDeleted() bool {
	return s != UserStatusDeleted
}

// CanTransitionTo 检查当前状态是否可以转换到目标状态
func (s UserStatus) CanTransitionTo(target UserStatus) bool {
	switch s {
	case UserStatusInactive:
		return target == UserStatusActive || target == UserStatusDeleted
	case UserStatusActive:
		return target == UserStatusSuspended || target == UserStatusBanned || 
			   target == UserStatusLocked || target == UserStatusDeleted
	case UserStatusSuspended:
		return target == UserStatusActive || target == UserStatusBanned || target == UserStatusDeleted
	case UserStatusBanned:
		return target == UserStatusActive || target == UserStatusDeleted
	case UserStatusLocked:
		return target == UserStatusActive || target == UserStatusDeleted
	case UserStatusDeleted:
		return false // 已删除状态不能转换到其他状态
	default:
		return false
	}
}

// GetAllowedTransitions 获取当前状态允许转换的目标状态列表
func (s UserStatus) GetAllowedTransitions() []UserStatus {
	switch s {
	case UserStatusInactive:
		return []UserStatus{UserStatusActive, UserStatusDeleted}
	case UserStatusActive:
		return []UserStatus{UserStatusSuspended, UserStatusBanned, UserStatusLocked, UserStatusDeleted}
	case UserStatusSuspended:
		return []UserStatus{UserStatusActive, UserStatusBanned, UserStatusDeleted}
	case UserStatusBanned:
		return []UserStatus{UserStatusActive, UserStatusDeleted}
	case UserStatusLocked:
		return []UserStatus{UserStatusActive, UserStatusDeleted}
	case UserStatusDeleted:
		return []UserStatus{} // 已删除状态不能转换
	default:
		return []UserStatus{}
	}
}

// GetTransitionReason 获取状态转换的描述
func (s UserStatus) GetTransitionReason(target UserStatus) string {
	if !s.CanTransitionTo(target) {
		return "不允许的状态转换"
	}
	
	switch {
	case s == UserStatusInactive && target == UserStatusActive:
		return "激活用户"
	case s == UserStatusActive && target == UserStatusSuspended:
		return "暂停用户"
	case s == UserStatusSuspended && target == UserStatusActive:
		return "恢复用户"
	case s == UserStatusActive && target == UserStatusBanned:
		return "封禁用户"
	case s == UserStatusBanned && target == UserStatusActive:
		return "解封用户"
	case s == UserStatusActive && target == UserStatusLocked:
		return "锁定用户"
	case s == UserStatusLocked && target == UserStatusActive:
		return "解锁用户"
	case target == UserStatusDeleted:
		return "删除用户"
	default:
		return "状态转换"
	}
}

// IsValid 检查状态是否有效
func (s UserStatus) IsValid() bool {
	return s >= UserStatusInactive && s <= UserStatusLocked
}

// Equals 比较两个状态是否相等
func (s UserStatus) Equals(other UserStatus) bool {
	return s == other
} 