package valueobject

import (
	"errors"
	"strings"
)

// ID 表示唯一标识符值对象
type ID string

// NewID 创建一个新的ID
func NewID(id string) ID {
	return ID(id)
}

// String 返回ID的字符串表示
func (id ID) String() string {
	return string(id)
}

// IsEmpty 检查ID是否为空
func (id ID) IsEmpty() bool {
	return id == ""
}

// UserID 表示用户标识符值对象
type UserID string

// NewUserID 创建一个新的UserID
func NewUserID(id string) UserID {
	return UserID(id)
}

// String 返回UserID的字符串表示
func (id UserID) String() string {
	return string(id)
}

// IsEmpty 检查UserID是否为空
func (id UserID) IsEmpty() bool {
	return id == ""
}

// CommunityName 表示社区名称值对象
type CommunityName string

// NewCommunityName 创建一个新的CommunityName
func NewCommunityName(name string) (CommunityName, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", errors.New("社区名称不能为空")
	}
	if len(trimmedName) > 100 {
		return "", errors.New("社区名称过长（最多100个字符）")
	}
	return CommunityName(trimmedName), nil
}

// String 返回CommunityName的字符串表示
func (name CommunityName) String() string {
	return string(name)
}

// IsEmpty 检查CommunityName是否为空
func (name CommunityName) IsEmpty() bool {
	return name == ""
}

// Description 表示描述值对象
type Description string

// NewDescription 创建一个新的Description
func NewDescription(desc string) Description {
	return Description(desc)
}

// String 返回Description的字符串表示
func (desc Description) String() string {
	return string(desc)
}

// CommunityStatus 表示社区状态值对象
type CommunityStatus string

const (
	StatusActive   CommunityStatus = "ACTIVE"
	StatusInactive CommunityStatus = "INACTIVE"
	StatusDeleted  CommunityStatus = "DELETED"
)

// Tag 表示标签值对象
type Tag string

// NewTag 创建一个新的Tag
func NewTag(tag string) (Tag, error) {
	trimmedTag := strings.TrimSpace(tag)
	if trimmedTag == "" {
		return "", errors.New("标签不能为空")
	}
	return Tag(trimmedTag), nil
}

// String 返回Tag的字符串表示
func (tag Tag) String() string {
	return string(tag)
}

// Location 表示位置值对象
type Location string

// NewLocation 创建一个新的Location
func NewLocation(location string) Location {
	return Location(location)
}

// String 返回Location的字符串表示
func (loc Location) String() string {
	return string(loc)
}
