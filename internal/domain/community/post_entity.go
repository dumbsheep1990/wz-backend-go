package community

import (
	"errors"
	"time"
)

// Post 帖子聚合根
type Post struct {
	id          ID
	title       Title
	content     Content
	authorID    ID
	authorName  string
	communityID ID
	groupID     ID // 可选，如果不属于特定小组则为空
	status      PostStatus
	tags        Tags
	images      ImageURLs
	createdAt   Timestamp
	updatedAt   Timestamp
	
	// 聚合计数器
	likeCount    int
	viewCount    int
	commentCount int
}

// NewPost 创建帖子
func NewPost(
	title Title,
	content Content,
	authorID ID,
	authorName string,
	communityID ID,
	groupID ID, // 如果不属于特定小组，可以传入空ID
	tags Tags,
	images ImageURLs,
) (*Post, error) {
	now := NewTimestamp()
	
	return &Post{
		id:           NewID(),
		title:        title,
		content:      content,
		authorID:     authorID,
		authorName:   authorName,
		communityID:  communityID,
		groupID:      groupID,
		status:       PostStatusActive,
		tags:         tags,
		images:       images,
		createdAt:    now,
		updatedAt:    now,
		likeCount:    0,
		viewCount:    0,
		commentCount: 0,
	}, nil
}

// ReconstructPost 从数据库重构帖子实体
func ReconstructPost(
	id string,
	title string,
	content string,
	authorID string,
	authorName string,
	communityID string,
	groupID string,
	status int,
	tags []string,
	images []string,
	createdAt time.Time,
	updatedAt time.Time,
	likeCount int,
	viewCount int,
	commentCount int,
) (*Post, error) {
	postID, err := NewIDFromString(id)
	if err != nil {
		return nil, err
	}
	
	postTitle, err := NewTitle(title)
	if err != nil {
		return nil, err
	}
	
	postContent, err := NewContent(content)
	if err != nil {
		return nil, err
	}
	
	authID, err := NewIDFromString(authorID)
	if err != nil {
		return nil, err
	}
	
	commID, err := NewIDFromString(communityID)
	if err != nil {
		return nil, err
	}
	
	var grpID ID
	if groupID != "" {
		grpID, err = NewIDFromString(groupID)
		if err != nil {
			return nil, err
		}
	}
	
	postTags, err := NewTags(tags)
	if err != nil {
		return nil, err
	}
	
	postImages, err := NewImageURLs(images)
	if err != nil {
		return nil, err
	}
	
	return &Post{
		id:           postID,
		title:        postTitle,
		content:      postContent,
		authorID:     authID,
		authorName:   authorName,
		communityID:  commID,
		groupID:      grpID,
		status:       PostStatus(status),
		tags:         postTags,
		images:       postImages,
		createdAt:    NewTimestampFromTime(createdAt),
		updatedAt:    NewTimestampFromTime(updatedAt),
		likeCount:    likeCount,
		viewCount:    viewCount,
		commentCount: commentCount,
	}, nil
}

// ID 获取帖子ID
func (p *Post) ID() ID {
	return p.id
}

// Title 获取帖子标题
func (p *Post) Title() Title {
	return p.title
}

// Content 获取帖子内容
func (p *Post) Content() Content {
	return p.content
}

// AuthorID 获取作者ID
func (p *Post) AuthorID() ID {
	return p.authorID
}

// AuthorName 获取作者名称
func (p *Post) AuthorName() string {
	return p.authorName
}

// CommunityID 获取所属社区ID
func (p *Post) CommunityID() ID {
	return p.communityID
}

// GroupID 获取所属小组ID
func (p *Post) GroupID() ID {
	return p.groupID
}

// Status 获取帖子状态
func (p *Post) Status() PostStatus {
	return p.status
}

// Tags 获取帖子标签
func (p *Post) Tags() Tags {
	return p.tags
}

// Images 获取帖子图片
func (p *Post) Images() ImageURLs {
	return p.images
}

// CreatedAt 获取创建时间
func (p *Post) CreatedAt() Timestamp {
	return p.createdAt
}

// UpdatedAt 获取更新时间
func (p *Post) UpdatedAt() Timestamp {
	return p.updatedAt
}

// LikeCount 获取点赞数量
func (p *Post) LikeCount() int {
	return p.likeCount
}

// ViewCount 获取浏览数量
func (p *Post) ViewCount() int {
	return p.viewCount
}

// CommentCount 获取评论数量
func (p *Post) CommentCount() int {
	return p.commentCount
}

// UpdateTitle 更新帖子标题
func (p *Post) UpdateTitle(title Title) error {
	if p.status == PostStatusDeleted {
		return errors.New("已删除的帖子不能修改")
	}
	
	p.title = title
	p.updatedAt = NewTimestamp()
	return nil
}

// UpdateContent 更新帖子内容
func (p *Post) UpdateContent(content Content) error {
	if p.status == PostStatusDeleted {
		return errors.New("已删除的帖子不能修改")
	}
	
	p.content = content
	p.updatedAt = NewTimestamp()
	return nil
}

// UpdateTags 更新帖子标签
func (p *Post) UpdateTags(tags Tags) error {
	if p.status == PostStatusDeleted {
		return errors.New("已删除的帖子不能修改")
	}
	
	p.tags = tags
	p.updatedAt = NewTimestamp()
	return nil
}

// UpdateImages 更新帖子图片
func (p *Post) UpdateImages(images ImageURLs) error {
	if p.status == PostStatusDeleted {
		return errors.New("已删除的帖子不能修改")
	}
	
	p.images = images
	p.updatedAt = NewTimestamp()
	return nil
}

// UpdateStatus 更新帖子状态
func (p *Post) UpdateStatus(status PostStatus) {
	p.status = status
	p.updatedAt = NewTimestamp()
}

// Like 点赞
func (p *Post) Like() {
	if p.status != PostStatusDeleted {
		p.likeCount++
		p.updatedAt = NewTimestamp()
	}
}

// Unlike 取消点赞
func (p *Post) Unlike() {
	if p.status != PostStatusDeleted && p.likeCount > 0 {
		p.likeCount--
		p.updatedAt = NewTimestamp()
	}
}

// View 浏览
func (p *Post) View() {
	if p.status != PostStatusDeleted {
		p.viewCount++
	}
}

// IncrementCommentCount 增加评论数量
func (p *Post) IncrementCommentCount() {
	if p.status != PostStatusDeleted {
		p.commentCount++
		p.updatedAt = NewTimestamp()
	}
}

// DecrementCommentCount 减少评论数量
func (p *Post) DecrementCommentCount() {
	if p.status != PostStatusDeleted && p.commentCount > 0 {
		p.commentCount--
		p.updatedAt = NewTimestamp()
	}
}

// Activate 激活帖子
func (p *Post) Activate() {
	p.status = PostStatusActive
	p.updatedAt = NewTimestamp()
}

// Deactivate 停用帖子
func (p *Post) Deactivate() {
	p.status = PostStatusInactive
	p.updatedAt = NewTimestamp()
}

// Delete 删除帖子
func (p *Post) Delete() {
	p.status = PostStatusDeleted
	p.updatedAt = NewTimestamp()
}

// IsAuthor 检查用户是否为帖子作者
func (p *Post) IsAuthor(userID ID) bool {
	return p.authorID.Value() == userID.Value()
}

// IsActive 检查帖子是否激活
func (p *Post) IsActive() bool {
	return p.status == PostStatusActive
}

// IsDeleted 检查帖子是否已删除
func (p *Post) IsDeleted() bool {
	return p.status == PostStatusDeleted
}

// BelongsToGroup 检查帖子是否属于特定小组
func (p *Post) BelongsToGroup() bool {
	return p.groupID.Value() != ""
}

// CanComment 检查是否可以评论
func (p *Post) CanComment() bool {
	return p.status == PostStatusActive
}
