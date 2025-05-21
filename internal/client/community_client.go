package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/yourusername/wz-backend-go/api/community"
)

// CommunityClient 是社区服务的客户端封装
type CommunityClient struct {
	client pb.CommunityServiceClient
	conn   *grpc.ClientConn
}

// NewCommunityClient 创建一个新的社区服务客户端
func NewCommunityClient(serviceAddr string) (*CommunityClient, error) {
	if serviceAddr == "" {
		serviceAddr = "localhost:50054" // 默认地址
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("无法连接到社区服务: %v", err)
	}

	client := pb.NewCommunityServiceClient(conn)
	return &CommunityClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *CommunityClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListCommunities 获取社区列表
func (c *CommunityClient) ListCommunities(ctx context.Context, pageSize int32, pageToken, filter string) (*pb.ListCommunitiesResponse, error) {
	return c.client.ListCommunities(ctx, &pb.ListCommunitiesRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
		Filter:    filter,
	})
}

// GetCommunity 获取社区详情
func (c *CommunityClient) GetCommunity(ctx context.Context, id string) (*pb.Community, error) {
	return c.client.GetCommunity(ctx, &pb.GetCommunityRequest{
		Id: id,
	})
}

// CreateCommunity 创建社区
func (c *CommunityClient) CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.Community, error) {
	return c.client.CreateCommunity(ctx, req)
}

// UpdateCommunity 更新社区
func (c *CommunityClient) UpdateCommunity(ctx context.Context, req *pb.UpdateCommunityRequest) (*pb.Community, error) {
	return c.client.UpdateCommunity(ctx, req)
}

// DeleteCommunity 删除社区
func (c *CommunityClient) DeleteCommunity(ctx context.Context, id string) (*pb.DeleteResponse, error) {
	return c.client.DeleteCommunity(ctx, &pb.DeleteCommunityRequest{
		Id: id,
	})
}

// ListGroups 获取群组列表
func (c *CommunityClient) ListGroups(ctx context.Context, communityID string, pageSize int32, pageToken, filter string) (*pb.ListGroupsResponse, error) {
	return c.client.ListGroups(ctx, &pb.ListGroupsRequest{
		CommunityId: communityID,
		PageSize:    pageSize,
		PageToken:   pageToken,
		Filter:      filter,
	})
}

// GetGroup 获取群组详情
func (c *CommunityClient) GetGroup(ctx context.Context, id string) (*pb.Group, error) {
	return c.client.GetGroup(ctx, &pb.GetGroupRequest{
		Id: id,
	})
}

// CreateGroup 创建群组
func (c *CommunityClient) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.Group, error) {
	return c.client.CreateGroup(ctx, req)
}

// UpdateGroup 更新群组
func (c *CommunityClient) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.Group, error) {
	return c.client.UpdateGroup(ctx, req)
}

// DeleteGroup 删除群组
func (c *CommunityClient) DeleteGroup(ctx context.Context, id string) (*pb.DeleteResponse, error) {
	return c.client.DeleteGroup(ctx, &pb.DeleteGroupRequest{
		Id: id,
	})
}

// JoinGroup 加入群组
func (c *CommunityClient) JoinGroup(ctx context.Context, groupID, userID, userName string) (*pb.JoinGroupResponse, error) {
	return c.client.JoinGroup(ctx, &pb.JoinGroupRequest{
		GroupId:  groupID,
		UserId:   userID,
		UserName: userName,
	})
}

// LeaveGroup 离开群组
func (c *CommunityClient) LeaveGroup(ctx context.Context, groupID, userID string) (*pb.LeaveGroupResponse, error) {
	return c.client.LeaveGroup(ctx, &pb.LeaveGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})
}

// ListPosts 获取帖子列表
func (c *CommunityClient) ListPosts(ctx context.Context, communityID, groupID string, pageSize int32, pageToken, filter string) (*pb.ListPostsResponse, error) {
	return c.client.ListPosts(ctx, &pb.ListPostsRequest{
		CommunityId: communityID,
		GroupId:     groupID,
		PageSize:    pageSize,
		PageToken:   pageToken,
		Filter:      filter,
	})
}

// GetPost 获取帖子详情
func (c *CommunityClient) GetPost(ctx context.Context, id string) (*pb.Post, error) {
	return c.client.GetPost(ctx, &pb.GetPostRequest{
		Id: id,
	})
}

// CreatePost 创建帖子
func (c *CommunityClient) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	return c.client.CreatePost(ctx, req)
}

// UpdatePost 更新帖子
func (c *CommunityClient) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	return c.client.UpdatePost(ctx, req)
}

// DeletePost 删除帖子
func (c *CommunityClient) DeletePost(ctx context.Context, id string) (*pb.DeleteResponse, error) {
	return c.client.DeletePost(ctx, &pb.DeletePostRequest{
		Id: id,
	})
}

// LikePost 点赞帖子
func (c *CommunityClient) LikePost(ctx context.Context, postID, userID string, like bool) (*pb.LikePostResponse, error) {
	return c.client.LikePost(ctx, &pb.LikePostRequest{
		PostId: postID,
		UserId: userID,
		Like:   like,
	})
}

// ViewPost 查看帖子
func (c *CommunityClient) ViewPost(ctx context.Context, postID, userID string) (*pb.ViewPostResponse, error) {
	return c.client.ViewPost(ctx, &pb.ViewPostRequest{
		PostId: postID,
		UserId: userID,
	})
}

// ListComments 获取评论列表
func (c *CommunityClient) ListComments(ctx context.Context, postID, parentID string, pageSize int32, pageToken string) (*pb.ListCommentsResponse, error) {
	return c.client.ListComments(ctx, &pb.ListCommentsRequest{
		PostId:    postID,
		ParentId:  parentID,
		PageSize:  pageSize,
		PageToken: pageToken,
	})
}

// CreateComment 创建评论
func (c *CommunityClient) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
	return c.client.CreateComment(ctx, req)
}

// DeleteComment 删除评论
func (c *CommunityClient) DeleteComment(ctx context.Context, id string) (*pb.DeleteResponse, error) {
	return c.client.DeleteComment(ctx, &pb.DeleteCommentRequest{
		Id: id,
	})
}
