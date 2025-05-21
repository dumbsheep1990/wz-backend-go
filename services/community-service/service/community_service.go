package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/yourusername/wz-backend-go/api/community"
)

// CommunityServiceServer 实现gRPC服务
type CommunityServiceServer struct {
	pb.UnimplementedCommunityServiceServer
	mu          sync.RWMutex
	communities map[string]*pb.Community
	groups      map[string]*pb.Group
	posts       map[string]*pb.Post
	comments    map[string]*pb.Comment
}

// NewCommunityServiceServer 创建一个新的CommunityServiceServer实例
func NewCommunityServiceServer() *CommunityServiceServer {
	return &CommunityServiceServer{
		communities: make(map[string]*pb.Community),
		groups:      make(map[string]*pb.Group),
		posts:       make(map[string]*pb.Post),
		comments:    make(map[string]*pb.Comment),
	}
}

// ListCommunities 实现ListCommunities RPC方法
func (s *CommunityServiceServer) ListCommunities(ctx context.Context, req *pb.ListCommunitiesRequest) (*pb.ListCommunitiesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	communities := make([]*pb.Community, 0, len(s.communities))
	for _, community := range s.communities {
		communities = append(communities, community)
	}
	
	return &pb.ListCommunitiesResponse{
		Communities: communities,
	}, nil
}

// GetCommunity 实现GetCommunity RPC方法
func (s *CommunityServiceServer) GetCommunity(ctx context.Context, req *pb.GetCommunityRequest) (*pb.Community, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	community, ok := s.communities[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}
	
	return community, nil
}

// CreateCommunity 实现CreateCommunity RPC方法
func (s *CommunityServiceServer) CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.Community, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	id := fmt.Sprintf("comm-%d", time.Now().UnixNano())
	community := &pb.Community{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		OwnerId:     req.OwnerId,
		OwnerName:   req.OwnerName,
		CreateTime:  time.Now().Format(time.RFC3339),
		Status:      pb.CommunityStatus_ACTIVE,
		Tags:        req.Tags,
		Location:    req.Location,
		GroupCount:  0,
		MemberCount: 1, // Owner is the first member
		PostCount:   0,
	}
	
	s.communities[id] = community
	return community, nil
}

// UpdateCommunity 实现UpdateCommunity RPC方法
func (s *CommunityServiceServer) UpdateCommunity(ctx context.Context, req *pb.UpdateCommunityRequest) (*pb.Community, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	community, ok := s.communities[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}
	
	if req.Name != "" {
		community.Name = req.Name
	}
	if req.Description != "" {
		community.Description = req.Description
	}
	if req.Tags != nil {
		community.Tags = req.Tags
	}
	if req.Location != "" {
		community.Location = req.Location
	}
	if req.Status != pb.CommunityStatus_COMMUNITY_STATUS_UNSPECIFIED {
		community.Status = req.Status
	}
	
	return community, nil
}

// DeleteCommunity 实现DeleteCommunity RPC方法
func (s *CommunityServiceServer) DeleteCommunity(ctx context.Context, req *pb.DeleteCommunityRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, ok := s.communities[req.Id]; !ok {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}
	
	delete(s.communities, req.Id)
	
	// 同时删除属于此社区的所有群组、帖子和评论
	// 这是一个简化的方法
	for id, group := range s.groups {
		if group.CommunityId == req.Id {
			delete(s.groups, id)
		}
	}
	
	for id, post := range s.posts {
		if post.CommunityId == req.Id {
			delete(s.posts, id)
			
			// Delete comments for this post
			for commentID, comment := range s.comments {
				if comment.PostId == id {
					delete(s.comments, commentID)
				}
			}
		}
	}
	
	return &pb.DeleteResponse{
		Success: true,
		Message: "Community deleted successfully",
	}, nil
}

// ListGroups 实现ListGroups RPC方法
func (s *CommunityServiceServer) ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var filteredGroups []*pb.Group
	
	if req.CommunityId != "" {
		for _, group := range s.groups {
			if group.CommunityId == req.CommunityId {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = make([]*pb.Group, 0, len(s.groups))
		for _, group := range s.groups {
			filteredGroups = append(filteredGroups, group)
		}
	}
	
	return &pb.ListGroupsResponse{
		Groups: filteredGroups,
	}, nil
}

// GetGroup 实现GetGroup RPC方法
func (s *CommunityServiceServer) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	group, ok := s.groups[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	
	return group, nil
}

// CreateGroup 实现CreateGroup RPC方法
func (s *CommunityServiceServer) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 验证社区是否存在
	if _, ok := s.communities[req.CommunityId]; !ok {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}
	
	id := fmt.Sprintf("group-%d", time.Now().UnixNano())
	group := &pb.Group{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		CommunityId: req.CommunityId,
		OwnerId:     req.OwnerId,
		OwnerName:   req.OwnerName,
		CreateTime:  time.Now().Format(time.RFC3339),
		Status:      pb.GroupStatus_ACTIVE,
		Members:     []string{req.OwnerId}, // Owner is the first member
		Tags:        req.Tags,
		MemberCount: 1,
		PostCount:   0,
	}
	
	s.groups[id] = group
	
	// Update community group count
	community := s.communities[req.CommunityId]
	community.GroupCount++
	
	return group, nil
}

// UpdateGroup 实现UpdateGroup RPC方法
func (s *CommunityServiceServer) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	group, ok := s.groups[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	if req.Tags != nil {
		group.Tags = req.Tags
	}
	if req.Status != pb.GroupStatus_GROUP_STATUS_UNSPECIFIED {
		group.Status = req.Status
	}
	
	return group, nil
}

// DeleteGroup 实现DeleteGroup RPC方法
func (s *CommunityServiceServer) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	group, ok := s.groups[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	
	communityID := group.CommunityId
	delete(s.groups, req.Id)
	
	// Update community group count
	if community, ok := s.communities[communityID]; ok {
		community.GroupCount--
	}
	
	// 删除相关的帖子和评论
	for id, post := range s.posts {
		if post.GroupId == req.Id {
			delete(s.posts, id)
			
			// Delete comments for this post
			for commentID, comment := range s.comments {
				if comment.PostId == id {
					delete(s.comments, commentID)
				}
			}
		}
	}
	
	return &pb.DeleteResponse{
		Success: true,
		Message: "Group deleted successfully",
	}, nil
}

// JoinGroup 实现JoinGroup RPC方法
func (s *CommunityServiceServer) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	group, ok := s.groups[req.GroupId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	
	// 检查用户是否已经是成员
	for _, memberID := range group.Members {
		if memberID == req.UserId {
			return &pb.JoinGroupResponse{
				Success: false,
				Message: "User is already a member of this group",
			}, nil
		}
	}
	
	// 将用户添加到群组
	group.Members = append(group.Members, req.UserId)
	group.MemberCount++
	
	// Update community member count
	if community, ok := s.communities[group.CommunityId]; ok {
		community.MemberCount++
	}
	
	return &pb.JoinGroupResponse{
		Success: true,
		Message: "User joined group successfully",
	}, nil
}

// LeaveGroup 实现LeaveGroup RPC方法
func (s *CommunityServiceServer) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	group, ok := s.groups[req.GroupId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	
	// 创建者不能离开群组
	if group.OwnerId == req.UserId {
		return &pb.LeaveGroupResponse{
			Success: false,
			Message: "Owner cannot leave the group",
		}, nil
	}
	
	// 查找并从成员中移除用户
	found := false
	for i, memberID := range group.Members {
		if memberID == req.UserId {
			// 通过用最后一个元素替换并截断数组来移除成员
			group.Members[i] = group.Members[len(group.Members)-1]
			group.Members = group.Members[:len(group.Members)-1]
			found = true
			group.MemberCount--
			break
		}
	}
	
	if !found {
		return &pb.LeaveGroupResponse{
			Success: false,
			Message: "User is not a member of this group",
		}, nil
	}
	
	// Update community member count
	if community, ok := s.communities[group.CommunityId]; ok {
		community.MemberCount--
	}
	
	return &pb.LeaveGroupResponse{
		Success: true,
		Message: "User left group successfully",
	}, nil
}

// ListPosts 实现ListPosts RPC方法
func (s *CommunityServiceServer) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var filteredPosts []*pb.Post
	
	if req.GroupId != "" {
		// 按群组筛选
		for _, post := range s.posts {
			if post.GroupId == req.GroupId {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else if req.CommunityId != "" {
		// 按社区筛选
		for _, post := range s.posts {
			if post.CommunityId == req.CommunityId {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else {
		// 返回所有帖子
		filteredPosts = make([]*pb.Post, 0, len(s.posts))
		for _, post := range s.posts {
			filteredPosts = append(filteredPosts, post)
		}
	}
	
	return &pb.ListPostsResponse{
		Posts: filteredPosts,
	}, nil
}

// GetPost 实现GetPost RPC方法
func (s *CommunityServiceServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	post, ok := s.posts[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	return post, nil
}

// CreatePost 实现CreatePost RPC方法
func (s *CommunityServiceServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 检查社区是否存在
	if _, ok := s.communities[req.CommunityId]; !ok {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}
	
	// 如果指定了群组，检查它是否存在并属于该社区
	if req.GroupId != "" {
		group, ok := s.groups[req.GroupId]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "Group not found")
		}
		if group.CommunityId != req.CommunityId {
			return nil, status.Errorf(codes.InvalidArgument, "Group does not belong to the specified community")
		}
	}
	
	id := fmt.Sprintf("post-%d", time.Now().UnixNano())
	now := time.Now().Format(time.RFC3339)
	post := &pb.Post{
		Id:          id,
		Title:       req.Title,
		Content:     req.Content,
		AuthorId:    req.AuthorId,
		AuthorName:  req.AuthorName,
		CommunityId: req.CommunityId,
		GroupId:     req.GroupId,
		CreateTime:  now,
		UpdateTime:  now,
		Status:      pb.PostStatus_ACTIVE,
		LikeCount:   0,
		ViewCount:   0,
		CommentCount: 0,
		Tags:        req.Tags,
		Images:      req.Images,
	}
	
	s.posts[id] = post
	
	// 更新社区中的帖子数量
	if community, ok := s.communities[req.CommunityId]; ok {
		community.PostCount++
	}
	
	// 如果适用，更新群组中的帖子数量
	if req.GroupId != "" {
		if group, ok := s.groups[req.GroupId]; ok {
			group.PostCount++
		}
	}
	
	return post, nil
}

// UpdatePost 实现UpdatePost RPC方法
func (s *CommunityServiceServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	post, ok := s.posts[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Tags != nil {
		post.Tags = req.Tags
	}
	if req.Images != nil {
		post.Images = req.Images
	}
	if req.Status != pb.PostStatus_POST_STATUS_UNSPECIFIED {
		post.Status = req.Status
	}
	
	post.UpdateTime = time.Now().Format(time.RFC3339)
	
	return post, nil
}

// DeletePost 实现DeletePost RPC方法
func (s *CommunityServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	post, ok := s.posts[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	communityID := post.CommunityId
	groupID := post.GroupId
	
	delete(s.posts, req.Id)
	
	// 更新计数
	if community, ok := s.communities[communityID]; ok {
		community.PostCount--
	}
	
	if groupID != "" {
		if group, ok := s.groups[groupID]; ok {
			group.PostCount--
		}
	}
	
	// Delete comments for this post
	for commentID, comment := range s.comments {
		if comment.PostId == req.Id {
			delete(s.comments, commentID)
		}
	}
	
	return &pb.DeleteResponse{
		Success: true,
		Message: "Post deleted successfully",
	}, nil
}

// LikePost 实现LikePost RPC方法
func (s *CommunityServiceServer) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	post, ok := s.posts[req.PostId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	// 在实际实现中，我们会跟踪哪些用户点赞了哪些帖子
	// 为简单起见，我们只是增加或减少计数
	if req.Like {
		post.LikeCount++
	} else if post.LikeCount > 0 {
		post.LikeCount--
	}
	
	return &pb.LikePostResponse{
		LikeCount: post.LikeCount,
	}, nil
}

// ViewPost 实现ViewPost RPC方法
func (s *CommunityServiceServer) ViewPost(ctx context.Context, req *pb.ViewPostRequest) (*pb.ViewPostResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	post, ok := s.posts[req.PostId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	// 在实际实现中，我们会跟踪唯一访问量
	post.ViewCount++
	
	return &pb.ViewPostResponse{
		ViewCount: post.ViewCount,
	}, nil
}

// ListComments 实现ListComments RPC方法
func (s *CommunityServiceServer) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var filteredComments []*pb.Comment
	
	for _, comment := range s.comments {
		if comment.PostId == req.PostId {
			if req.ParentId == "" && comment.ParentId == "" {
				// 获取顶层评论
				filteredComments = append(filteredComments, comment)
			} else if req.ParentId != "" && comment.ParentId == req.ParentId {
				// 获取对特定评论的回复
				filteredComments = append(filteredComments, comment)
			}
		}
	}
	
	return &pb.ListCommentsResponse{
		Comments: filteredComments,
	}, nil
}

// CreateComment 实现CreateComment RPC方法
func (s *CommunityServiceServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 检查帖子是否存在
	post, ok := s.posts[req.PostId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Post not found")
	}
	
	// 如果这是回复，检查父评论是否存在
	if req.ParentId != "" {
		if _, ok := s.comments[req.ParentId]; !ok {
			return nil, status.Errorf(codes.NotFound, "Parent comment not found")
		}
	}
	
	id := fmt.Sprintf("comment-%d", time.Now().UnixNano())
	comment := &pb.Comment{
		Id:         id,
		Content:    req.Content,
		AuthorId:   req.AuthorId,
		AuthorName: req.AuthorName,
		PostId:     req.PostId,
		ParentId:   req.ParentId,
		CreateTime: time.Now().Format(time.RFC3339),
		Status:     pb.CommentStatus_ACTIVE,
		LikeCount:  0,
	}
	
	s.comments[id] = comment
	
	// 更新帖子的评论数量
	post.CommentCount++
	
	return comment, nil
}

// DeleteComment 实现DeleteComment RPC方法
func (s *CommunityServiceServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	comment, ok := s.comments[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Comment not found")
	}
	
	// 标记为已删除而不是实际删除
	comment.Status = pb.CommentStatus_DELETED
	comment.Content = "[Deleted]"
	
	// 在实际实现中，你可能想要删除嵌套评论
	// 或者以不同方式处理评论树
	
	// 减少帖子上的评论计数
	if post, ok := s.posts[comment.PostId]; ok {
		post.CommentCount--
	}
	
	return &pb.DeleteResponse{
		Success: true,
		Message: "Comment deleted successfully",
	}, nil
}
