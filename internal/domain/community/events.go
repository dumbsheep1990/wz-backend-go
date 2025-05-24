package community

import (
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// 社区相关事件

// CommunityCreatedEvent 社区创建事件
type CommunityCreatedEvent struct {
	communityID   ID
	communityName CommunityName
	ownerID       ID
	occurredAt    time.Time
}

// NewCommunityCreatedEvent 创建社区创建事件
func NewCommunityCreatedEvent(communityID ID, communityName CommunityName, ownerID ID, occurredAt time.Time) *CommunityCreatedEvent {
	return &CommunityCreatedEvent{
		communityID:   communityID,
		communityName: communityName,
		ownerID:       ownerID,
		occurredAt:    occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommunityCreatedEvent) EventType() string {
	return "CommunityCreatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommunityCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommunityID 获取社区ID
func (e *CommunityCreatedEvent) CommunityID() ID {
	return e.communityID
}

// CommunityName 获取社区名称
func (e *CommunityCreatedEvent) CommunityName() CommunityName {
	return e.communityName
}

// OwnerID 获取拥有者ID
func (e *CommunityCreatedEvent) OwnerID() ID {
	return e.ownerID
}

// CommunityUpdatedEvent 社区更新事件
type CommunityUpdatedEvent struct {
	communityID   ID
	communityName CommunityName
	updaterID     ID
	occurredAt    time.Time
}

// NewCommunityUpdatedEvent 创建社区更新事件
func NewCommunityUpdatedEvent(communityID ID, communityName CommunityName, updaterID ID, occurredAt time.Time) *CommunityUpdatedEvent {
	return &CommunityUpdatedEvent{
		communityID:   communityID,
		communityName: communityName,
		updaterID:     updaterID,
		occurredAt:    occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommunityUpdatedEvent) EventType() string {
	return "CommunityUpdatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommunityUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommunityID 获取社区ID
func (e *CommunityUpdatedEvent) CommunityID() ID {
	return e.communityID
}

// CommunityName 获取社区名称
func (e *CommunityUpdatedEvent) CommunityName() CommunityName {
	return e.communityName
}

// UpdaterID 获取更新者ID
func (e *CommunityUpdatedEvent) UpdaterID() ID {
	return e.updaterID
}

// CommunityDeletedEvent 社区删除事件
type CommunityDeletedEvent struct {
	communityID ID
	deleterID   ID
	occurredAt  time.Time
}

// NewCommunityDeletedEvent 创建社区删除事件
func NewCommunityDeletedEvent(communityID ID, deleterID ID, occurredAt time.Time) *CommunityDeletedEvent {
	return &CommunityDeletedEvent{
		communityID: communityID,
		deleterID:   deleterID,
		occurredAt:  occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommunityDeletedEvent) EventType() string {
	return "CommunityDeletedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommunityDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommunityID 获取社区ID
func (e *CommunityDeletedEvent) CommunityID() ID {
	return e.communityID
}

// DeleterID 获取删除者ID
func (e *CommunityDeletedEvent) DeleterID() ID {
	return e.deleterID
}

// 小组相关事件

// GroupCreatedEvent 小组创建事件
type GroupCreatedEvent struct {
	groupID      ID
	groupName    CommunityName
	communityID  ID
	ownerID      ID
	occurredAt   time.Time
}

// NewGroupCreatedEvent 创建小组创建事件
func NewGroupCreatedEvent(groupID ID, groupName CommunityName, communityID ID, ownerID ID, occurredAt time.Time) *GroupCreatedEvent {
	return &GroupCreatedEvent{
		groupID:     groupID,
		groupName:   groupName,
		communityID: communityID,
		ownerID:     ownerID,
		occurredAt:  occurredAt,
	}
}

// EventType 获取事件类型
func (e *GroupCreatedEvent) EventType() string {
	return "GroupCreatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *GroupCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// GroupID 获取小组ID
func (e *GroupCreatedEvent) GroupID() ID {
	return e.groupID
}

// GroupName 获取小组名称
func (e *GroupCreatedEvent) GroupName() CommunityName {
	return e.groupName
}

// CommunityID 获取社区ID
func (e *GroupCreatedEvent) CommunityID() ID {
	return e.communityID
}

// OwnerID 获取拥有者ID
func (e *GroupCreatedEvent) OwnerID() ID {
	return e.ownerID
}

// GroupUpdatedEvent 小组更新事件
type GroupUpdatedEvent struct {
	groupID     ID
	groupName   CommunityName
	updaterID   ID
	occurredAt  time.Time
}

// NewGroupUpdatedEvent 创建小组更新事件
func NewGroupUpdatedEvent(groupID ID, groupName CommunityName, updaterID ID, occurredAt time.Time) *GroupUpdatedEvent {
	return &GroupUpdatedEvent{
		groupID:    groupID,
		groupName:  groupName,
		updaterID:  updaterID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *GroupUpdatedEvent) EventType() string {
	return "GroupUpdatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *GroupUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// GroupID 获取小组ID
func (e *GroupUpdatedEvent) GroupID() ID {
	return e.groupID
}

// GroupName 获取小组名称
func (e *GroupUpdatedEvent) GroupName() CommunityName {
	return e.groupName
}

// UpdaterID 获取更新者ID
func (e *GroupUpdatedEvent) UpdaterID() ID {
	return e.updaterID
}

// GroupDeletedEvent 小组删除事件
type GroupDeletedEvent struct {
	groupID     ID
	deleterID   ID
	occurredAt  time.Time
}

// NewGroupDeletedEvent 创建小组删除事件
func NewGroupDeletedEvent(groupID ID, deleterID ID, occurredAt time.Time) *GroupDeletedEvent {
	return &GroupDeletedEvent{
		groupID:    groupID,
		deleterID:  deleterID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *GroupDeletedEvent) EventType() string {
	return "GroupDeletedEvent"
}

// OccurredAt 获取事件发生时间
func (e *GroupDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// GroupID 获取小组ID
func (e *GroupDeletedEvent) GroupID() ID {
	return e.groupID
}

// DeleterID 获取删除者ID
func (e *GroupDeletedEvent) DeleterID() ID {
	return e.deleterID
}

// MemberJoinedGroupEvent 成员加入小组事件
type MemberJoinedGroupEvent struct {
	groupID     ID
	memberID    ID
	occurredAt  time.Time
}

// NewMemberJoinedGroupEvent 创建成员加入小组事件
func NewMemberJoinedGroupEvent(groupID ID, memberID ID, occurredAt time.Time) *MemberJoinedGroupEvent {
	return &MemberJoinedGroupEvent{
		groupID:    groupID,
		memberID:   memberID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *MemberJoinedGroupEvent) EventType() string {
	return "MemberJoinedGroupEvent"
}

// OccurredAt 获取事件发生时间
func (e *MemberJoinedGroupEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// GroupID 获取小组ID
func (e *MemberJoinedGroupEvent) GroupID() ID {
	return e.groupID
}

// MemberID 获取成员ID
func (e *MemberJoinedGroupEvent) MemberID() ID {
	return e.memberID
}

// MemberLeftGroupEvent 成员离开小组事件
type MemberLeftGroupEvent struct {
	groupID     ID
	memberID    ID
	occurredAt  time.Time
}

// NewMemberLeftGroupEvent 创建成员离开小组事件
func NewMemberLeftGroupEvent(groupID ID, memberID ID, occurredAt time.Time) *MemberLeftGroupEvent {
	return &MemberLeftGroupEvent{
		groupID:    groupID,
		memberID:   memberID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *MemberLeftGroupEvent) EventType() string {
	return "MemberLeftGroupEvent"
}

// OccurredAt 获取事件发生时间
func (e *MemberLeftGroupEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// GroupID 获取小组ID
func (e *MemberLeftGroupEvent) GroupID() ID {
	return e.groupID
}

// MemberID 获取成员ID
func (e *MemberLeftGroupEvent) MemberID() ID {
	return e.memberID
}

// 帖子相关事件

// PostCreatedEvent 帖子创建事件
type PostCreatedEvent struct {
	postID      ID
	title       Title
	authorID    ID
	communityID ID
	groupID     ID
	occurredAt  time.Time
}

// NewPostCreatedEvent 创建帖子创建事件
func NewPostCreatedEvent(postID ID, title Title, authorID ID, communityID ID, groupID ID, occurredAt time.Time) *PostCreatedEvent {
	return &PostCreatedEvent{
		postID:      postID,
		title:       title,
		authorID:    authorID,
		communityID: communityID,
		groupID:     groupID,
		occurredAt:  occurredAt,
	}
}

// EventType 获取事件类型
func (e *PostCreatedEvent) EventType() string {
	return "PostCreatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *PostCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PostID 获取帖子ID
func (e *PostCreatedEvent) PostID() ID {
	return e.postID
}

// Title 获取帖子标题
func (e *PostCreatedEvent) Title() Title {
	return e.title
}

// AuthorID 获取作者ID
func (e *PostCreatedEvent) AuthorID() ID {
	return e.authorID
}

// CommunityID 获取社区ID
func (e *PostCreatedEvent) CommunityID() ID {
	return e.communityID
}

// GroupID 获取小组ID
func (e *PostCreatedEvent) GroupID() ID {
	return e.groupID
}

// PostUpdatedEvent 帖子更新事件
type PostUpdatedEvent struct {
	postID     ID
	title      Title
	updaterID  ID
	occurredAt time.Time
}

// NewPostUpdatedEvent 创建帖子更新事件
func NewPostUpdatedEvent(postID ID, title Title, updaterID ID, occurredAt time.Time) *PostUpdatedEvent {
	return &PostUpdatedEvent{
		postID:     postID,
		title:      title,
		updaterID:  updaterID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *PostUpdatedEvent) EventType() string {
	return "PostUpdatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *PostUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PostID 获取帖子ID
func (e *PostUpdatedEvent) PostID() ID {
	return e.postID
}

// Title 获取帖子标题
func (e *PostUpdatedEvent) Title() Title {
	return e.title
}

// UpdaterID 获取更新者ID
func (e *PostUpdatedEvent) UpdaterID() ID {
	return e.updaterID
}

// PostDeletedEvent 帖子删除事件
type PostDeletedEvent struct {
	postID     ID
	deleterID  ID
	occurredAt time.Time
}

// NewPostDeletedEvent 创建帖子删除事件
func NewPostDeletedEvent(postID ID, deleterID ID, occurredAt time.Time) *PostDeletedEvent {
	return &PostDeletedEvent{
		postID:     postID,
		deleterID:  deleterID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *PostDeletedEvent) EventType() string {
	return "PostDeletedEvent"
}

// OccurredAt 获取事件发生时间
func (e *PostDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PostID 获取帖子ID
func (e *PostDeletedEvent) PostID() ID {
	return e.postID
}

// DeleterID 获取删除者ID
func (e *PostDeletedEvent) DeleterID() ID {
	return e.deleterID
}

// PostLikedEvent 帖子点赞事件
type PostLikedEvent struct {
	postID     ID
	userID     ID
	occurredAt time.Time
}

// NewPostLikedEvent 创建帖子点赞事件
func NewPostLikedEvent(postID ID, userID ID, occurredAt time.Time) *PostLikedEvent {
	return &PostLikedEvent{
		postID:     postID,
		userID:     userID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *PostLikedEvent) EventType() string {
	return "PostLikedEvent"
}

// OccurredAt 获取事件发生时间
func (e *PostLikedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PostID 获取帖子ID
func (e *PostLikedEvent) PostID() ID {
	return e.postID
}

// UserID 获取用户ID
func (e *PostLikedEvent) UserID() ID {
	return e.userID
}

// PostViewedEvent 帖子浏览事件
type PostViewedEvent struct {
	postID     ID
	userID     ID
	occurredAt time.Time
}

// NewPostViewedEvent 创建帖子浏览事件
func NewPostViewedEvent(postID ID, userID ID, occurredAt time.Time) *PostViewedEvent {
	return &PostViewedEvent{
		postID:     postID,
		userID:     userID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *PostViewedEvent) EventType() string {
	return "PostViewedEvent"
}

// OccurredAt 获取事件发生时间
func (e *PostViewedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PostID 获取帖子ID
func (e *PostViewedEvent) PostID() ID {
	return e.postID
}

// UserID 获取用户ID
func (e *PostViewedEvent) UserID() ID {
	return e.userID
}

// 评论相关事件

// CommentCreatedEvent 评论创建事件
type CommentCreatedEvent struct {
	commentID  ID
	postID     ID
	authorID   ID
	occurredAt time.Time
}

// NewCommentCreatedEvent 创建评论创建事件
func NewCommentCreatedEvent(commentID ID, postID ID, authorID ID, occurredAt time.Time) *CommentCreatedEvent {
	return &CommentCreatedEvent{
		commentID:  commentID,
		postID:     postID,
		authorID:   authorID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommentCreatedEvent) EventType() string {
	return "CommentCreatedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommentCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommentID 获取评论ID
func (e *CommentCreatedEvent) CommentID() ID {
	return e.commentID
}

// PostID 获取帖子ID
func (e *CommentCreatedEvent) PostID() ID {
	return e.postID
}

// AuthorID 获取作者ID
func (e *CommentCreatedEvent) AuthorID() ID {
	return e.authorID
}

// CommentDeletedEvent 评论删除事件
type CommentDeletedEvent struct {
	commentID  ID
	postID     ID
	deleterID  ID
	occurredAt time.Time
}

// NewCommentDeletedEvent 创建评论删除事件
func NewCommentDeletedEvent(commentID ID, postID ID, deleterID ID, occurredAt time.Time) *CommentDeletedEvent {
	return &CommentDeletedEvent{
		commentID:  commentID,
		postID:     postID,
		deleterID:  deleterID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommentDeletedEvent) EventType() string {
	return "CommentDeletedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommentDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommentID 获取评论ID
func (e *CommentDeletedEvent) CommentID() ID {
	return e.commentID
}

// PostID 获取帖子ID
func (e *CommentDeletedEvent) PostID() ID {
	return e.postID
}

// DeleterID 获取删除者ID
func (e *CommentDeletedEvent) DeleterID() ID {
	return e.deleterID
}

// CommentLikedEvent 评论点赞事件
type CommentLikedEvent struct {
	commentID  ID
	userID     ID
	occurredAt time.Time
}

// NewCommentLikedEvent 创建评论点赞事件
func NewCommentLikedEvent(commentID ID, userID ID, occurredAt time.Time) *CommentLikedEvent {
	return &CommentLikedEvent{
		commentID:  commentID,
		userID:     userID,
		occurredAt: occurredAt,
	}
}

// EventType 获取事件类型
func (e *CommentLikedEvent) EventType() string {
	return "CommentLikedEvent"
}

// OccurredAt 获取事件发生时间
func (e *CommentLikedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// CommentID 获取评论ID
func (e *CommentLikedEvent) CommentID() ID {
	return e.commentID
}

// UserID 获取用户ID
func (e *CommentLikedEvent) UserID() ID {
	return e.userID
}
