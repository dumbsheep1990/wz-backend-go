package community

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/wz-backend-go/internal/application/community/dto"
	"github.com/yourusername/wz-backend-go/internal/application/community/service"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
	pb "github.com/yourusername/wz-backend-go/api/community"
)

// CommunityGRPCAdapter adapts between gRPC requests and the application service
type CommunityGRPCAdapter struct {
	pb.UnimplementedCommunityServiceServer
	communityService *service.CommunityApplicationService
}

// NewCommunityGRPCAdapter creates a new CommunityGRPCAdapter
func NewCommunityGRPCAdapter(communityService *service.CommunityApplicationService) *CommunityGRPCAdapter {
	return &CommunityGRPCAdapter{
		communityService: communityService,
	}
}

// ListCommunities implements the ListCommunities gRPC method
func (a *CommunityGRPCAdapter) ListCommunities(ctx context.Context, req *pb.ListCommunitiesRequest) (*pb.ListCommunitiesResponse, error) {
	// Convert gRPC request to application DTO
	appReq := dto.ListCommunitiesRequest{
		OwnerID:  req.OwnerId,
		Tag:      req.Tag,
		Location: req.Location,
		Status:   req.Status.String(),
		Offset:   int(req.Offset),
		Limit:    int(req.Limit),
	}

	// Call application service
	result, err := a.communityService.ListCommunities(ctx, appReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert application response to gRPC response
	communities := make([]*pb.Community, 0, len(result.Communities))
	for _, c := range result.Communities {
		communities = append(communities, convertDTOToPbCommunity(c))
	}

	return &pb.ListCommunitiesResponse{
		Communities: communities,
		Total:       int32(result.Total),
	}, nil
}

// GetCommunity implements the GetCommunity gRPC method
func (a *CommunityGRPCAdapter) GetCommunity(ctx context.Context, req *pb.GetCommunityRequest) (*pb.Community, error) {
	// Convert gRPC request to application DTO
	appReq := dto.GetCommunityRequest{
		ID: req.Id,
	}

	// Call application service
	result, err := a.communityService.GetCommunity(ctx, appReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert application response to gRPC response
	return convertDTOToPbCommunity(result), nil
}

// CreateCommunity implements the CreateCommunity gRPC method
func (a *CommunityGRPCAdapter) CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.Community, error) {
	// Convert gRPC request to application DTO
	appReq := dto.CreateCommunityRequest{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerId,
		OwnerName:   req.OwnerName,
		Tags:        req.Tags,
		Location:    req.Location,
	}

	// Call application service
	result, err := a.communityService.CreateCommunity(ctx, appReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert application response to gRPC response
	return convertDTOToPbCommunity(result), nil
}

// UpdateCommunity implements the UpdateCommunity gRPC method
func (a *CommunityGRPCAdapter) UpdateCommunity(ctx context.Context, req *pb.UpdateCommunityRequest) (*pb.Community, error) {
	// Convert gRPC request to application DTO
	appReq := dto.UpdateCommunityRequest{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		Location:    req.Location,
		Status:      req.Status.String(),
	}

	// Call application service
	result, err := a.communityService.UpdateCommunity(ctx, appReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert application response to gRPC response
	return convertDTOToPbCommunity(result), nil
}

// DeleteCommunity implements the DeleteCommunity gRPC method
func (a *CommunityGRPCAdapter) DeleteCommunity(ctx context.Context, req *pb.DeleteCommunityRequest) (*pb.DeleteResponse, error) {
	// Convert gRPC request to application DTO
	appReq := dto.DeleteCommunityRequest{
		ID: req.Id,
	}

	// Call application service
	err := a.communityService.DeleteCommunity(ctx, appReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteResponse{
		Success: true,
		Message: "Community deleted successfully",
	}, nil
}

// Helper functions

// convertDTOToPbCommunity converts a DTO to a protobuf Community
func convertDTOToPbCommunity(community *dto.CommunityDTO) *pb.Community {
	var pbStatus pb.CommunityStatus
	switch community.Status {
	case string(valueobject.StatusActive):
		pbStatus = pb.CommunityStatus_ACTIVE
	case string(valueobject.StatusInactive):
		pbStatus = pb.CommunityStatus_INACTIVE
	case string(valueobject.StatusDeleted):
		pbStatus = pb.CommunityStatus_DELETED
	default:
		pbStatus = pb.CommunityStatus_COMMUNITY_STATUS_UNSPECIFIED
	}

	return &pb.Community{
		Id:          community.ID,
		Name:        community.Name,
		Description: community.Description,
		OwnerId:     community.OwnerID,
		OwnerName:   community.OwnerName,
		Status:      pbStatus,
		Tags:        community.Tags,
		Location:    community.Location,
		GroupCount:  int32(community.GroupCount),
		MemberCount: int32(community.MemberCount),
		PostCount:   int32(community.PostCount),
		CreateTime:  community.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdateTime:  community.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
