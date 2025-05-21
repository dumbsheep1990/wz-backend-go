package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	pb "github.com/yourusername/wz-backend-go/api/community"
	"github.com/yourusername/wz-backend-go/services/community-service/service"
)

// RegisterRoutes 设置社区服务的所有路由
func RegisterRoutes(router *gin.Engine, service *service.CommunityServiceServer) {
	v1 := router.Group("/api/v1")
	
	// 身份验证中间件 - 目前简化版本
	authMiddleware := func(c *gin.Context) {
		// 在这里你应该检查有效的认证令牌
		// 目前，我们只是检查用户ID是否存在
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
	
	// 公开路由（不需要认证）
	communities := v1.Group("/communities")
	{
		communities.GET("", listCommunities(service))
		communities.GET("/:id", getCommunity(service))
		
		// 需要认证的路由
		authCommunities := communities.Group("")
		authCommunities.Use(authMiddleware)
		{
			authCommunities.POST("", createCommunity(service))
			authCommunities.PUT("/:id", updateCommunity(service))
			authCommunities.DELETE("/:id", deleteCommunity(service))
		}
	}
	
	groups := v1.Group("/groups")
	{
		groups.GET("", listGroups(service))
		groups.GET("/:id", getGroup(service))
		
		// 需要认证的路由
		authGroups := groups.Group("")
		authGroups.Use(authMiddleware)
		{
			authGroups.POST("", createGroup(service))
			authGroups.PUT("/:id", updateGroup(service))
			authGroups.DELETE("/:id", deleteGroup(service))
			authGroups.POST("/:id/join", joinGroup(service))
			authGroups.POST("/:id/leave", leaveGroup(service))
		}
	}
	
	posts := v1.Group("/posts")
	{
		posts.GET("", listPosts(service))
		posts.GET("/:id", getPost(service))
		
		// 需要认证的路由
		authPosts := posts.Group("")
		authPosts.Use(authMiddleware)
		{
			authPosts.POST("", createPost(service))
			authPosts.PUT("/:id", updatePost(service))
			authPosts.DELETE("/:id", deletePost(service))
			authPosts.POST("/:id/like", likePost(service))
			authPosts.POST("/:id/view", viewPost(service))
		}
	}
	
	comments := v1.Group("/comments")
	{
		comments.GET("", listComments(service))
		
		// 需要认证的路由
		authComments := comments.Group("")
		authComments.Use(authMiddleware)
		{
			authComments.POST("", createComment(service))
			authComments.DELETE("/:id", deleteComment(service))
		}
	}
}

// 社区处理器
func listCommunities(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		pageToken := c.DefaultQuery("page_token", "")
		filter := c.DefaultQuery("filter", "")
		
		// 调用服务
		resp, err := service.ListCommunities(context.Background(), &pb.ListCommunitiesRequest{
			PageSize:  int32(pageSize),
			PageToken: pageToken,
			Filter:    filter,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func getCommunity(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		community, err := service.GetCommunity(context.Background(), &pb.GetCommunityRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
			return
		}
		
		c.JSON(http.StatusOK, community)
	}
}

func createCommunity(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pb.CreateCommunityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		req.OwnerId = userID.(string)
		
		community, err := service.CreateCommunity(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, community)
	}
}

func updateCommunity(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req pb.UpdateCommunityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		req.Id = id
		
		community, err := service.UpdateCommunity(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, community)
	}
}

func deleteCommunity(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		resp, err := service.DeleteCommunity(context.Background(), &pb.DeleteCommunityRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

// 群组处理器
func listGroups(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Query("community_id")
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		pageToken := c.DefaultQuery("page_token", "")
		filter := c.DefaultQuery("filter", "")
		
		resp, err := service.ListGroups(context.Background(), &pb.ListGroupsRequest{
			CommunityId: communityID,
			PageSize:    int32(pageSize),
			PageToken:   pageToken,
			Filter:      filter,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func getGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		group, err := service.GetGroup(context.Background(), &pb.GetGroupRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		
		c.JSON(http.StatusOK, group)
	}
}

func createGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pb.CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		req.OwnerId = userID.(string)
		
		group, err := service.CreateGroup(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, group)
	}
}

func updateGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req pb.UpdateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		req.Id = id
		
		group, err := service.UpdateGroup(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, group)
	}
}

func deleteGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		resp, err := service.DeleteGroup(context.Background(), &pb.DeleteGroupRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func joinGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		userName := c.GetHeader("X-User-Name") // 简化处理；在实际应用中，你应从用户服务获取这个信息
		
		resp, err := service.JoinGroup(context.Background(), &pb.JoinGroupRequest{
			GroupId:  id,
			UserId:   userID.(string),
			UserName: userName,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func leaveGroup(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		
		resp, err := service.LeaveGroup(context.Background(), &pb.LeaveGroupRequest{
			GroupId: id,
			UserId:  userID.(string),
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

// 帖子处理器
func listPosts(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Query("community_id")
		groupID := c.Query("group_id")
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		pageToken := c.DefaultQuery("page_token", "")
		filter := c.DefaultQuery("filter", "")
		
		resp, err := service.ListPosts(context.Background(), &pb.ListPostsRequest{
			CommunityId: communityID,
			GroupId:     groupID,
			PageSize:    int32(pageSize),
			PageToken:   pageToken,
			Filter:      filter,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func getPost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		post, err := service.GetPost(context.Background(), &pb.GetPostRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		
		c.JSON(http.StatusOK, post)
	}
}

func createPost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pb.CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		userName := c.GetHeader("X-User-Name") // Simplified
		
		req.AuthorId = userID.(string)
		req.AuthorName = userName
		
		post, err := service.CreatePost(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, post)
	}
}

func updatePost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req pb.UpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		req.Id = id
		
		post, err := service.UpdatePost(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, post)
	}
}

func deletePost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		resp, err := service.DeletePost(context.Background(), &pb.DeletePostRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func likePost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Like bool `json:"like"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		
		resp, err := service.LikePost(context.Background(), &pb.LikePostRequest{
			PostId: id,
			UserId: userID.(string),
			Like:   req.Like,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func viewPost(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		// Get user info from auth context if available
		userID := c.GetHeader("X-User-ID")
		
		resp, err := service.ViewPost(context.Background(), &pb.ViewPostRequest{
			PostId: id,
			UserId: userID,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

// 评论处理器
func listComments(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Query("post_id")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
			return
		}
		
		parentID := c.Query("parent_id")
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		pageToken := c.DefaultQuery("page_token", "")
		
		resp, err := service.ListComments(context.Background(), &pb.ListCommentsRequest{
			PostId:    postID,
			ParentId:  parentID,
			PageSize:  int32(pageSize),
			PageToken: pageToken,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}

func createComment(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pb.CreateCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get user info from auth context
		userID, _ := c.Get("user_id")
		userName := c.GetHeader("X-User-Name") // Simplified
		
		req.AuthorId = userID.(string)
		req.AuthorName = userName
		
		comment, err := service.CreateComment(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, comment)
	}
}

func deleteComment(service *service.CommunityServiceServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		resp, err := service.DeleteComment(context.Background(), &pb.DeleteCommentRequest{
			Id: id,
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, resp)
	}
}
