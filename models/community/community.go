package community

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Community 社区模型
type Community struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string    `json:"name" db:"name" gorm:"not null;comment:社区名称"`
	Code        string    `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:社区编码"`
	Logo        string    `json:"logo" db:"logo" gorm:"comment:社区logo"`
	Cover       string    `json:"cover" db:"cover" gorm:"comment:社区封面"`
	Description string    `json:"description" db:"description" gorm:"comment:社区描述"`
	Rules       string    `json:"rules" db:"rules" gorm:"type:text;comment:社区规则"`
	CreatorID   int64     `json:"creatorId" db:"creator_id" gorm:"index;not null;comment:创建者ID"`
	MemberCount int       `json:"memberCount" db:"member_count" gorm:"default:0;comment:成员数量"`
	TopicCount  int       `json:"topicCount" db:"topic_count" gorm:"default:0;comment:话题数量"`
	PostCount   int       `json:"postCount" db:"post_count" gorm:"default:0;comment:帖子数量"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	IsPrivate   bool      `json:"isPrivate" db:"is_private" gorm:"default:false;comment:是否私有"`
	JoinMode    string    `json:"joinMode" db:"join_mode" gorm:"default:free;comment:加入方式:free-自由加入,apply-申请加入,invite-邀请加入"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityMember 社区成员
type CommunityMember struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID int64     `json:"communityId" db:"community_id" gorm:"uniqueIndex:idx_community_user;index;not null;comment:社区ID"`
	UserID      int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_community_user;index;not null;comment:用户ID"`
	Role        string    `json:"role" db:"role" gorm:"default:member;comment:角色:admin-管理员,moderator-版主,member-成员"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用,-1-已退出"`
	Nickname    string    `json:"nickname" db:"nickname" gorm:"comment:社区内昵称"`
	Bio         string    `json:"bio" db:"bio" gorm:"comment:社区内简介"`
	JoinMethod  string    `json:"joinMethod" db:"join_method" gorm:"comment:加入方式:apply-申请,invite-邀请,create-创建"`
	InviterID   int64     `json:"inviterId" db:"inviter_id" gorm:"index;comment:邀请人ID"`
	ApproverID  int64     `json:"approverId" db:"approver_id" gorm:"index;comment:审批人ID"`
	ApproveTime *time.Time `json:"approveTime" db:"approve_time" gorm:"comment:审批时间"`
	PostCount   int       `json:"postCount" db:"post_count" gorm:"default:0;comment:发帖数量"`
	LikeCount   int       `json:"likeCount" db:"like_count" gorm:"default:0;comment:获赞数量"`
	ExpPoints   int       `json:"expPoints" db:"exp_points" gorm:"default:0;comment:经验值"`
	Level       int       `json:"level" db:"level" gorm:"default:1;comment:等级"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityTopic 社区话题
type CommunityTopic struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	Name        string    `json:"name" db:"name" gorm:"not null;comment:话题名称"`
	Description string    `json:"description" db:"description" gorm:"comment:话题描述"`
	Icon        string    `json:"icon" db:"icon" gorm:"comment:话题图标"`
	Cover       string    `json:"cover" db:"cover" gorm:"comment:话题封面"`
	CreatorID   int64     `json:"creatorId" db:"creator_id" gorm:"index;not null;comment:创建者ID"`
	PostCount   int       `json:"postCount" db:"post_count" gorm:"default:0;comment:帖子数量"`
	FollowCount int       `json:"followCount" db:"follow_count" gorm:"default:0;comment:关注数量"`
	SortOrder   int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	IsTop       bool      `json:"isTop" db:"is_top" gorm:"default:false;comment:是否置顶"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityPost 社区帖子
type CommunityPost struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	TopicID     int64     `json:"topicId" db:"topic_id" gorm:"index;comment:话题ID"`
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Title       string    `json:"title" db:"title" gorm:"not null;comment:标题"`
	Content     string    `json:"content" db:"content" gorm:"type:text;not null;comment:内容"`
	ContentType string    `json:"contentType" db:"content_type" gorm:"default:text;comment:内容类型:text,html,markdown,rich"`
	Media       string    `json:"media" db:"media" gorm:"type:json;comment:媒体,JSON数组"`
	Status      int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-隐藏,-1-已删除"`
	ViewCount   int       `json:"viewCount" db:"view_count" gorm:"default:0;comment:查看数"`
	LikeCount   int       `json:"likeCount" db:"like_count" gorm:"default:0;comment:点赞数"`
	CommentCount int      `json:"commentCount" db:"comment_count" gorm:"default:0;comment:评论数"`
	IsTop       bool      `json:"isTop" db:"is_top" gorm:"default:false;comment:是否置顶"`
	IsEssence   bool      `json:"isEssence" db:"is_essence" gorm:"default:false;comment:是否精华"`
	IsLocked    bool      `json:"isLocked" db:"is_locked" gorm:"default:false;comment:是否锁定"`
	LastReplyAt *time.Time `json:"lastReplyAt" db:"last_reply_at" gorm:"comment:最后回复时间"`
	LastReplyID int64     `json:"lastReplyId" db:"last_reply_id" gorm:"comment:最后回复ID"`
	IP          string    `json:"ip" db:"ip" gorm:"comment:发帖IP"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityJoinApplication 社区加入申请
type CommunityJoinApplication struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	UserID      int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Reason      string    `json:"reason" db:"reason" gorm:"comment:申请理由"`
	Status      string    `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-待处理,approved-已通过,rejected-已拒绝"`
	ApproverID  int64     `json:"approverId" db:"approver_id" gorm:"index;comment:审批人ID"`
	ApproveTime *time.Time `json:"approveTime" db:"approve_time" gorm:"comment:审批时间"`
	RejectReason string    `json:"rejectReason" db:"reject_reason" gorm:"comment:拒绝理由"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityInvitation 社区邀请
type CommunityInvitation struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID  int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	InviterID    int64     `json:"inviterId" db:"inviter_id" gorm:"index;not null;comment:邀请人ID"`
	InviteeID    int64     `json:"inviteeId" db:"invitee_id" gorm:"index;comment:被邀请人ID"`
	InviteeEmail string    `json:"inviteeEmail" db:"invitee_email" gorm:"comment:被邀请人邮箱"`
	Message      string    `json:"message" db:"message" gorm:"comment:邀请信息"`
	Status       string    `json:"status" db:"status" gorm:"default:pending;comment:状态:pending-待处理,accepted-已接受,rejected-已拒绝,expired-已过期"`
	InviteCode   string    `json:"inviteCode" db:"invite_code" gorm:"uniqueIndex;not null;comment:邀请码"`
	ExpireTime   time.Time `json:"expireTime" db:"expire_time" gorm:"not null;comment:过期时间"`
	ResponseTime *time.Time `json:"responseTime" db:"response_time" gorm:"comment:响应时间"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityNotification 社区通知
type CommunityNotification struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID  int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	SenderID     int64     `json:"senderId" db:"sender_id" gorm:"index;not null;comment:发送者ID"`
	ReceiverID   int64     `json:"receiverId" db:"receiver_id" gorm:"index;comment:接收者ID,为空表示全体成员"`
	Type         string    `json:"type" db:"type" gorm:"not null;comment:通知类型:announcement-公告,event-活动,post-帖子,reply-回复,mention-提及"`
	Title        string    `json:"title" db:"title" gorm:"not null;comment:标题"`
	Content      string    `json:"content" db:"content" gorm:"type:text;not null;comment:内容"`
	RelatedID    int64     `json:"relatedId" db:"related_id" gorm:"index;comment:关联ID"`
	RelatedType  string    `json:"relatedType" db:"related_type" gorm:"comment:关联类型"`
	IsRead       bool      `json:"isRead" db:"is_read" gorm:"default:false;comment:是否已读"`
	ReadTime     *time.Time `json:"readTime" db:"read_time" gorm:"comment:阅读时间"`
	Importance   int       `json:"importance" db:"importance" gorm:"default:1;comment:重要性:1-普通,2-重要,3-紧急"`
	ExpireTime   *time.Time `json:"expireTime" db:"expire_time" gorm:"comment:过期时间"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityActivity 社区活动
type CommunityActivity struct {
	common.BaseIDModel
	common.BaseTimeModel
	CommunityID  int64     `json:"communityId" db:"community_id" gorm:"index;not null;comment:社区ID"`
	Title        string    `json:"title" db:"title" gorm:"not null;comment:活动标题"`
	Description  string    `json:"description" db:"description" gorm:"type:text;not null;comment:活动描述"`
	Cover        string    `json:"cover" db:"cover" gorm:"comment:活动封面"`
	Type         string    `json:"type" db:"type" gorm:"not null;comment:活动类型:online-线上,offline-线下,both-线上线下"`
	Location     string    `json:"location" db:"location" gorm:"comment:活动地点"`
	StartTime    time.Time `json:"startTime" db:"start_time" gorm:"not null;comment:开始时间"`
	EndTime      time.Time `json:"endTime" db:"end_time" gorm:"not null;comment:结束时间"`
	EnrollStart  time.Time `json:"enrollStart" db:"enroll_start" gorm:"not null;comment:报名开始时间"`
	EnrollEnd    time.Time `json:"enrollEnd" db:"enroll_end" gorm:"not null;comment:报名结束时间"`
	MaxMembers   int       `json:"maxMembers" db:"max_members" gorm:"comment:最大参与人数"`
	CurrentMembers int     `json:"currentMembers" db:"current_members" gorm:"default:0;comment:当前参与人数"`
	OrganizerId  int64     `json:"organizerId" db:"organizer_id" gorm:"index;not null;comment:组织者ID"`
	Status       string    `json:"status" db:"status" gorm:"default:draft;comment:状态:draft-草稿,published-已发布,ongoing-进行中,ended-已结束,canceled-已取消"`
	IsTop        bool      `json:"isTop" db:"is_top" gorm:"default:false;comment:是否置顶"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// CommunityActivityMember 社区活动成员
type CommunityActivityMember struct {
	common.BaseIDModel
	common.BaseTimeModel
	ActivityID  int64     `json:"activityId" db:"activity_id" gorm:"uniqueIndex:idx_activity_user;index;not null;comment:活动ID"`
	UserID      int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_activity_user;index;not null;comment:用户ID"`
	Status      string    `json:"status" db:"status" gorm:"default:enrolled;comment:状态:enrolled-已报名,attended-已参加,absent-缺席,canceled-已取消"`
	EnrollTime  time.Time `json:"enrollTime" db:"enroll_time" gorm:"not null;comment:报名时间"`
	CheckinTime *time.Time `json:"checkinTime" db:"checkin_time" gorm:"comment:签到时间"`
	Feedback    string    `json:"feedback" db:"feedback" gorm:"type:text;comment:反馈"`
	TenantID    int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
