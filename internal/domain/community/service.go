package community

import (
	"context"
	"errors"
	"time"
)

// CommunityService 社区领域服务接口
type CommunityService interface {
	// CreateCommunity 创建社区
	CreateCommunity(ctx context.Context, name CommunityName, description Description, ownerID ID, ownerName string, tags Tags, location Location, type_ CommunityType) (*Community, error)
	
	// UpdateCommunity 更新社区
	UpdateCommunity(ctx context.Context, communityID ID, userID ID, name CommunityName, description Description, tags Tags, location Location) (*Community, error)
	
	// DeleteCommunity 删除社区
	DeleteCommunity(ctx context.Context, communityID ID, userID ID) error
	
	// CreateGroup 创建小组
	CreateGroup(ctx context.Context, name CommunityName, description Description, communityID ID, ownerID ID, ownerName string, tags Tags) (*Group, error)
	
	// JoinGroup 加入小组
	JoinGroup(ctx context.Context, groupID ID, userID ID, userName string) error
	
	// LeaveGroup 离开小组
	LeaveGroup(ctx context.Context, groupID ID, userID ID) error
	
	// CreatePost 创建帖子
	CreatePost(ctx context.Context, title Title, content Content, authorID ID, authorName string, communityID ID, groupID ID, tags Tags, images ImageURLs) (*Post, error)
	
	// LikePost 点赞帖子
	LikePost(ctx context.Context, postID ID, userID ID) error
	
	// UnlikePost 取消点赞帖子
	UnlikePost(ctx context.Context, postID ID, userID ID) error
	
	// ViewPost 浏览帖子
	ViewPost(ctx context.Context, postID ID, userID ID) error
	
	// CreateComment 创建评论
	CreateComment(ctx context.Context, content Content, authorID ID, authorName string, postID ID, parentID ID) (*Comment, error)
	
	// DeleteComment 删除评论
	DeleteComment(ctx context.Context, commentID ID, userID ID) error
}

// CommunityServiceImpl 社区领域服务实现
type CommunityServiceImpl struct {
	communityRepo CommunityRepository
	groupRepo     GroupRepository
	postRepo      PostRepository
	commentRepo   CommentRepository
	likeRepo      LikeRepository
	viewRepo      ViewRepository
}

// NewCommunityService 创建社区领域服务
func NewCommunityService(
	communityRepo CommunityRepository,
	groupRepo GroupRepository,
	postRepo PostRepository,
	commentRepo CommentRepository,
	likeRepo LikeRepository,
	viewRepo ViewRepository,
) CommunityService {
	return &CommunityServiceImpl{
		communityRepo: communityRepo,
		groupRepo:     groupRepo,
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		likeRepo:      likeRepo,
		viewRepo:      viewRepo,
	}
}

// CreateCommunity 创建社区
func (s *CommunityServiceImpl) CreateCommunity(
	ctx context.Context,
	name CommunityName,
	description Description,
	ownerID ID,
	ownerName string,
	tags Tags,
	location Location,
	type_ CommunityType,
) (*Community, error) {
	// 检查社区名称是否已存在
	existingCommunity, err := s.communityRepo.FindByName(ctx, name)
	if err == nil && existingCommunity != nil {
		return nil, errors.New("社区名称已存在")
	}
	
	// 创建社区
	community, err := NewCommunity(name, description, ownerID, ownerName, tags, location, type_)
	if err != nil {
		return nil, err
	}
	
	// 保存社区
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}
	
	return community, nil
}

// UpdateCommunity 更新社区
func (s *CommunityServiceImpl) UpdateCommunity(
	ctx context.Context,
	communityID ID,
	userID ID,
	name CommunityName,
	description Description,
	tags Tags,
	location Location,
) (*Community, error) {
	// 查找社区
	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, errors.New("社区不存在")
	}
	
	// 检查用户权限
	if !community.IsOwner(userID) {
		return nil, errors.New("只有社区拥有者才能更新社区信息")
	}
	
	// 如果名称发生变化，检查新名称是否已存在
	if name.Value() != community.Name().Value() {
		existingCommunity, err := s.communityRepo.FindByName(ctx, name)
		if err == nil && existingCommunity != nil && existingCommunity.ID().Value() != communityID.Value() {
			return nil, errors.New("社区名称已存在")
		}
		community.UpdateName(name)
	}
	
	// 更新其他信息
	community.UpdateDescription(description)
	community.UpdateTags(tags)
	community.UpdateLocation(location)
	
	// 保存更新
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}
	
	return community, nil
}

// DeleteCommunity 删除社区
func (s *CommunityServiceImpl) DeleteCommunity(ctx context.Context, communityID ID, userID ID) error {
	// 查找社区
	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return errors.New("社区不存在")
	}
	
	// 检查用户权限
	if !community.IsOwner(userID) {
		return errors.New("只有社区拥有者才能删除社区")
	}
	
	// 标记社区为已删除
	community.Delete()
	
	// 保存更新
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return err
	}
	
	return nil
}

// CreateGroup 创建小组
func (s *CommunityServiceImpl) CreateGroup(
	ctx context.Context,
	name CommunityName,
	description Description,
	communityID ID,
	ownerID ID,
	ownerName string,
	tags Tags,
) (*Group, error) {
	// 查找社区
	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, errors.New("社区不存在")
	}
	
	// 检查社区状态
	if !community.IsActive() {
		return nil, errors.New("只能在活跃的社区中创建小组")
	}
	
	// 创建小组
	group, err := NewGroup(name, description, communityID, ownerID, ownerName, tags)
	if err != nil {
		return nil, err
	}
	
	// 保存小组
	if err := s.groupRepo.Save(ctx, group); err != nil {
		return nil, err
	}
	
	// 更新社区的小组计数
	community.IncrementGroupCount()
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}
	
	return group, nil
}

// JoinGroup 加入小组
func (s *CommunityServiceImpl) JoinGroup(ctx context.Context, groupID ID, userID ID, userName string) error {
	// 查找小组
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return errors.New("小组不存在")
	}
	
	// 检查小组状态
	if !group.IsActive() {
		return errors.New("只能加入活跃的小组")
	}
	
	// 加入小组
	if err := group.AddMember(userID); err != nil {
		return err
	}
	
	// 保存小组
	if err := s.groupRepo.Save(ctx, group); err != nil {
		return err
	}
	
	// 更新社区的成员计数
	community, err := s.communityRepo.FindByID(ctx, group.CommunityID())
	if err != nil {
		return err
	}
	
	community.IncrementMemberCount()
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return err
	}
	
	return nil
}

// LeaveGroup 离开小组
func (s *CommunityServiceImpl) LeaveGroup(ctx context.Context, groupID ID, userID ID) error {
	// 查找小组
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return errors.New("小组不存在")
	}
	
	// 离开小组
	if err := group.RemoveMember(userID); err != nil {
		return err
	}
	
	// 保存小组
	if err := s.groupRepo.Save(ctx, group); err != nil {
		return err
	}
	
	// 更新社区的成员计数
	community, err := s.communityRepo.FindByID(ctx, group.CommunityID())
	if err != nil {
		return err
	}
	
	community.DecrementMemberCount()
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return err
	}
	
	return nil
}

// CreatePost 创建帖子
func (s *CommunityServiceImpl) CreatePost(
	ctx context.Context,
	title Title,
	content Content,
	authorID ID,
	authorName string,
	communityID ID,
	groupID ID,
	tags Tags,
	images ImageURLs,
) (*Post, error) {
	// 查找社区
	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, errors.New("社区不存在")
	}
	
	// 检查社区状态
	if !community.IsActive() {
		return nil, errors.New("只能在活跃的社区中发帖")
	}
	
	// 如果帖子属于小组，检查小组状态和成员资格
	if groupID.Value() != "" {
		group, err := s.groupRepo.FindByID(ctx, groupID)
		if err != nil {
			return nil, errors.New("小组不存在")
		}
		
		if !group.CanPost(authorID) {
			return nil, errors.New("只有小组成员才能在小组中发帖")
		}
	}
	
	// 创建帖子
	post, err := NewPost(title, content, authorID, authorName, communityID, groupID, tags, images)
	if err != nil {
		return nil, err
	}
	
	// 保存帖子
	if err := s.postRepo.Save(ctx, post); err != nil {
		return nil, err
	}
	
	// 更新社区的帖子计数
	community.IncrementPostCount()
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}
	
	// 如果帖子属于小组，更新小组的帖子计数
	if groupID.Value() != "" {
		group, err := s.groupRepo.FindByID(ctx, groupID)
		if err != nil {
			return nil, err
		}
		
		group.IncrementPostCount()
		if err := s.groupRepo.Save(ctx, group); err != nil {
			return nil, err
		}
	}
	
	return post, nil
}

// LikePost 点赞帖子
func (s *CommunityServiceImpl) LikePost(ctx context.Context, postID ID, userID ID) error {
	// 查找帖子
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	
	// 检查帖子状态
	if !post.IsActive() {
		return errors.New("不能点赞非活跃的帖子")
	}
	
	// 检查是否已点赞
	liked, err := s.likeRepo.CheckPostLike(ctx, postID, userID)
	if err != nil {
		return err
	}
	
	if liked {
		return errors.New("已经点赞过该帖子")
	}
	
	// 保存点赞记录
	if err := s.likeRepo.SavePostLike(ctx, postID, userID); err != nil {
		return err
	}
	
	// 更新帖子点赞计数
	post.Like()
	if err := s.postRepo.Save(ctx, post); err != nil {
		return err
	}
	
	return nil
}

// UnlikePost 取消点赞帖子
func (s *CommunityServiceImpl) UnlikePost(ctx context.Context, postID ID, userID ID) error {
	// 查找帖子
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	
	// 检查是否已点赞
	liked, err := s.likeRepo.CheckPostLike(ctx, postID, userID)
	if err != nil {
		return err
	}
	
	if !liked {
		return errors.New("未点赞过该帖子")
	}
	
	// 移除点赞记录
	if err := s.likeRepo.RemovePostLike(ctx, postID, userID); err != nil {
		return err
	}
	
	// 更新帖子点赞计数
	post.Unlike()
	if err := s.postRepo.Save(ctx, post); err != nil {
		return err
	}
	
	return nil
}

// ViewPost 浏览帖子
func (s *CommunityServiceImpl) ViewPost(ctx context.Context, postID ID, userID ID) error {
	// 查找帖子
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	
	// 检查帖子状态
	if !post.IsActive() {
		return errors.New("不能浏览非活跃的帖子")
	}
	
	// 检查是否已浏览
	viewed, err := s.viewRepo.CheckView(ctx, postID, userID)
	if err != nil {
		return err
	}
	
	if !viewed {
		// 保存浏览记录
		if err := s.viewRepo.SaveView(ctx, postID, userID); err != nil {
			return err
		}
		
		// 更新帖子浏览计数
		post.View()
		if err := s.postRepo.Save(ctx, post); err != nil {
			return err
		}
	}
	
	return nil
}

// CreateComment 创建评论
func (s *CommunityServiceImpl) CreateComment(
	ctx context.Context,
	content Content,
	authorID ID,
	authorName string,
	postID ID,
	parentID ID,
) (*Comment, error) {
	// 查找帖子
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}
	
	// 检查帖子状态
	if !post.CanComment() {
		return nil, errors.New("不能在当前帖子下评论")
	}
	
	// 如果是回复评论，检查父评论是否存在
	if parentID.Value() != "" {
		parentComment, err := s.commentRepo.FindByID(ctx, parentID)
		if err != nil {
			return nil, errors.New("父评论不存在")
		}
		
		if parentComment.IsDeleted() {
			return nil, errors.New("不能回复已删除的评论")
		}
		
		if parentComment.PostID().Value() != postID.Value() {
			return nil, errors.New("父评论不属于当前帖子")
		}
	}
	
	// 创建评论
	comment, err := NewComment(content, authorID, authorName, postID, parentID)
	if err != nil {
		return nil, err
	}
	
	// 保存评论
	if err := s.commentRepo.Save(ctx, comment); err != nil {
		return nil, err
	}
	
	// 更新帖子评论计数
	post.IncrementCommentCount()
	if err := s.postRepo.Save(ctx, post); err != nil {
		return nil, err
	}
	
	return comment, nil
}

// DeleteComment 删除评论
func (s *CommunityServiceImpl) DeleteComment(ctx context.Context, commentID ID, userID ID) error {
	// 查找评论
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	
	// 检查用户权限
	if !comment.IsAuthor(userID) {
		return errors.New("只有评论作者才能删除评论")
	}
	
	// 标记评论为已删除
	comment.Delete()
	
	// 保存更新
	if err := s.commentRepo.Save(ctx, comment); err != nil {
		return err
	}
	
	// 更新帖子评论计数
	post, err := s.postRepo.FindByID(ctx, comment.PostID())
	if err != nil {
		return err
	}
	
	post.DecrementCommentCount()
	if err := s.postRepo.Save(ctx, post); err != nil {
		return err
	}
	
	return nil
}
