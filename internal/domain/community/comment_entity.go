package community

import (
	"errors"
	"time"
)

// Comment 评论聚合根
type Comment struct {
	id         ID
	content    Content
	authorID   ID
	authorName string
	postID     ID
	parentID   ID // 父评论ID，如果是顶级评论则为空
	status     CommentStatus
	createdAt  Timestamp
	
	// 聚合计数器
	likeCount int
}

// NewComment 创建评论
func NewComment(
	content Content,
	authorID ID,
	authorName string,
	postID ID,
	parentID ID, // 如果是顶级评论，可以传入空ID
) (*Comment, error) {
	return &Comment{
		id:         NewID(),
		content:    content,
		authorID:   authorID,
		authorName: authorName,
		postID:     postID,
		parentID:   parentID,
		status:     CommentStatusActive,
		createdAt:  NewTimestamp(),
		likeCount:  0,
	}, nil
}

// ReconstructComment 从数据库重构评论实体
func ReconstructComment(
	id string,
	content string,
	authorID string,
	authorName string,
	postID string,
	parentID string,
	status int,
	createdAt time.Time,
	likeCount int,
) (*Comment, error) {
	commentID, err := NewIDFromString(id)
	if err != nil {
		return nil, err
	}
	
	commentContent, err := NewContent(content)
	if err != nil {
		return nil, err
	}
	
	authID, err := NewIDFromString(authorID)
	if err != nil {
		return nil, err
	}
	
	pID, err := NewIDFromString(postID)
	if err != nil {
		return nil, err
	}
	
	var parID ID
	if parentID != "" {
		parID, err = NewIDFromString(parentID)
		if err != nil {
			return nil, err
		}
	}
	
	return &Comment{
		id:         commentID,
		content:    commentContent,
		authorID:   authID,
		authorName: authorName,
		postID:     pID,
		parentID:   parID,
		status:     CommentStatus(status),
		createdAt:  NewTimestampFromTime(createdAt),
		likeCount:  likeCount,
	}, nil
}

// ID 获取评论ID
func (c *Comment) ID() ID {
	return c.id
}

// Content 获取评论内容
func (c *Comment) Content() Content {
	return c.content
}

// AuthorID 获取作者ID
func (c *Comment) AuthorID() ID {
	return c.authorID
}

// AuthorName 获取作者名称
func (c *Comment) AuthorName() string {
	return c.authorName
}

// PostID 获取帖子ID
func (c *Comment) PostID() ID {
	return c.postID
}

// ParentID 获取父评论ID
func (c *Comment) ParentID() ID {
	return c.parentID
}

// Status 获取评论状态
func (c *Comment) Status() CommentStatus {
	return c.status
}

// CreatedAt 获取创建时间
func (c *Comment) CreatedAt() Timestamp {
	return c.createdAt
}

// LikeCount 获取点赞数量
func (c *Comment) LikeCount() int {
	return c.likeCount
}

// Delete 删除评论
func (c *Comment) Delete() {
	c.status = CommentStatusDeleted
}

// Like 点赞
func (c *Comment) Like() {
	if c.status != CommentStatusDeleted {
		c.likeCount++
	}
}

// Unlike 取消点赞
func (c *Comment) Unlike() {
	if c.status != CommentStatusDeleted && c.likeCount > 0 {
		c.likeCount--
	}
}

// IsAuthor 检查用户是否为评论作者
func (c *Comment) IsAuthor(userID ID) bool {
	return c.authorID.Value() == userID.Value()
}

// IsActive 检查评论是否激活
func (c *Comment) IsActive() bool {
	return c.status == CommentStatusActive
}

// IsDeleted 检查评论是否已删除
func (c *Comment) IsDeleted() bool {
	return c.status == CommentStatusDeleted
}

// IsReply 检查是否为回复
func (c *Comment) IsReply() bool {
	return c.parentID.Value() != ""
}
