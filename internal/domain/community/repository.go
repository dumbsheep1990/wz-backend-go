package community

import (
	"context"
)

// CommunityRepository 社区仓储接口
type CommunityRepository interface {
	// Save 保存社区
	Save(ctx context.Context, community *Community) error
	
	// FindByID 根据ID查找社区
	FindByID(ctx context.Context, id ID) (*Community, error)
	
	// FindByName 根据名称查找社区
	FindByName(ctx context.Context, name CommunityName) (*Community, error)
	
	// FindByOwner 查找用户拥有的社区
	FindByOwner(ctx context.Context, ownerID ID) ([]*Community, error)
	
	// FindAll 查找所有社区（分页）
	FindAll(ctx context.Context, offset, limit int) ([]*Community, error)
	
	// FindByType 根据类型查找社区
	FindByType(ctx context.Context, type_ CommunityType, offset, limit int) ([]*Community, error)
	
	// CountAll 统计所有社区数量
	CountAll(ctx context.Context) (int, error)
	
	// CountByType 统计特定类型的社区数量
	CountByType(ctx context.Context, type_ CommunityType) (int, error)
	
	// Delete 删除社区
	Delete(ctx context.Context, id ID) error
}

// GroupRepository 小组仓储接口
type GroupRepository interface {
	// Save 保存小组
	Save(ctx context.Context, group *Group) error
	
	// FindByID 根据ID查找小组
	FindByID(ctx context.Context, id ID) (*Group, error)
	
	// FindByCommunity 查找社区下的小组（分页）
	FindByCommunity(ctx context.Context, communityID ID, offset, limit int) ([]*Group, error)
	
	// CountByCommunity 统计社区下的小组数量
	CountByCommunity(ctx context.Context, communityID ID) (int, error)
	
	// FindByMember 查找用户加入的小组
	FindByMember(ctx context.Context, memberID ID, offset, limit int) ([]*Group, error)
	
	// CountByMember 统计用户加入的小组数量
	CountByMember(ctx context.Context, memberID ID) (int, error)
	
	// FindByOwner 查找用户拥有的小组
	FindByOwner(ctx context.Context, ownerID ID) ([]*Group, error)
	
	// Delete 删除小组
	Delete(ctx context.Context, id ID) error
}

// PostRepository 帖子仓储接口
type PostRepository interface {
	// Save 保存帖子
	Save(ctx context.Context, post *Post) error
	
	// FindByID 根据ID查找帖子
	FindByID(ctx context.Context, id ID) (*Post, error)
	
	// FindByCommunity 查找社区下的帖子（分页）
	FindByCommunity(ctx context.Context, communityID ID, offset, limit int) ([]*Post, error)
	
	// CountByCommunity 统计社区下的帖子数量
	CountByCommunity(ctx context.Context, communityID ID) (int, error)
	
	// FindByGroup 查找小组下的帖子（分页）
	FindByGroup(ctx context.Context, groupID ID, offset, limit int) ([]*Post, error)
	
	// CountByGroup 统计小组下的帖子数量
	CountByGroup(ctx context.Context, groupID ID) (int, error)
	
	// FindByAuthor 查找用户发布的帖子（分页）
	FindByAuthor(ctx context.Context, authorID ID, offset, limit int) ([]*Post, error)
	
	// CountByAuthor 统计用户发布的帖子数量
	CountByAuthor(ctx context.Context, authorID ID) (int, error)
	
	// FindByTags 根据标签查找帖子（分页）
	FindByTags(ctx context.Context, tags []string, offset, limit int) ([]*Post, error)
	
	// FindRecent 查找最近的帖子（分页）
	FindRecent(ctx context.Context, offset, limit int) ([]*Post, error)
	
	// FindPopular 查找热门帖子（分页）
	FindPopular(ctx context.Context, offset, limit int) ([]*Post, error)
	
	// Delete 删除帖子
	Delete(ctx context.Context, id ID) error
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	// Save 保存评论
	Save(ctx context.Context, comment *Comment) error
	
	// FindByID 根据ID查找评论
	FindByID(ctx context.Context, id ID) (*Comment, error)
	
	// FindByPost 查找帖子下的评论（分页）
	FindByPost(ctx context.Context, postID ID, offset, limit int) ([]*Comment, error)
	
	// CountByPost 统计帖子下的评论数量
	CountByPost(ctx context.Context, postID ID) (int, error)
	
	// FindByParent 查找父评论下的回复（分页）
	FindByParent(ctx context.Context, parentID ID, offset, limit int) ([]*Comment, error)
	
	// CountByParent 统计父评论下的回复数量
	CountByParent(ctx context.Context, parentID ID) (int, error)
	
	// FindByAuthor 查找用户发布的评论（分页）
	FindByAuthor(ctx context.Context, authorID ID, offset, limit int) ([]*Comment, error)
	
	// CountByAuthor 统计用户发布的评论数量
	CountByAuthor(ctx context.Context, authorID ID) (int, error)
	
	// FindRecent 查找最近的评论（分页）
	FindRecent(ctx context.Context, offset, limit int) ([]*Comment, error)
	
	// Delete 删除评论
	Delete(ctx context.Context, id ID) error
}

// LikeRepository 点赞仓储接口
type LikeRepository interface {
	// SavePostLike 保存帖子点赞
	SavePostLike(ctx context.Context, postID ID, userID ID) error
	
	// RemovePostLike 移除帖子点赞
	RemovePostLike(ctx context.Context, postID ID, userID ID) error
	
	// CheckPostLike 检查用户是否点赞过帖子
	CheckPostLike(ctx context.Context, postID ID, userID ID) (bool, error)
	
	// CountPostLikes 统计帖子点赞数量
	CountPostLikes(ctx context.Context, postID ID) (int, error)
	
	// SaveCommentLike 保存评论点赞
	SaveCommentLike(ctx context.Context, commentID ID, userID ID) error
	
	// RemoveCommentLike 移除评论点赞
	RemoveCommentLike(ctx context.Context, commentID ID, userID ID) error
	
	// CheckCommentLike 检查用户是否点赞过评论
	CheckCommentLike(ctx context.Context, commentID ID, userID ID) (bool, error)
	
	// CountCommentLikes 统计评论点赞数量
	CountCommentLikes(ctx context.Context, commentID ID) (int, error)
}

// ViewRepository 浏览记录仓储接口
type ViewRepository interface {
	// SaveView 保存浏览记录
	SaveView(ctx context.Context, postID ID, userID ID) error
	
	// CheckView 检查用户是否浏览过帖子
	CheckView(ctx context.Context, postID ID, userID ID) (bool, error)
	
	// CountViews 统计帖子浏览数量
	CountViews(ctx context.Context, postID ID) (int, error)
	
	// FindPopularByViews 查找浏览量最高的帖子（分页）
	FindPopularByViews(ctx context.Context, offset, limit int) ([]*Post, error)
}
