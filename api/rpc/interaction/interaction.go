package interaction

import (
	"context"

	"google.golang.org/grpc"
)

// 交互服务接口定义
type InteractionService interface {
	// 获取评论列表
	GetCommentList(ctx context.Context, in *GetCommentListReq) (*GetCommentListResp, error)
	// 获取评论详情
	GetCommentDetail(ctx context.Context, in *GetCommentDetailReq) (*CommentDetailResp, error)
	// 更新评论状态
	UpdateCommentStatus(ctx context.Context, in *UpdateCommentStatusReq) (*UpdateCommentStatusResp, error)
	// 删除评论
	DeleteComment(ctx context.Context, in *DeleteCommentReq) (*DeleteCommentResp, error)
	// 获取举报列表
	GetReportList(ctx context.Context, in *GetReportListReq) (*GetReportListResp, error)
	// 获取举报详情
	GetReportDetail(ctx context.Context, in *GetReportDetailReq) (*ReportDetailResp, error)
	// 处理举报
	HandleReport(ctx context.Context, in *HandleReportReq) (*HandleReportResp, error)
}

// 交互服务RPC客户端
type interactionServiceClient struct {
	conn *grpc.ClientConn
}

// 创建交互服务客户端
func NewInteractionService(conn *grpc.ClientConn) InteractionService {
	return &interactionServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// 获取评论列表请求
type GetCommentListReq struct {
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
	UserId    int64  `json:"user_id,omitempty"`
	ContentId int64  `json:"content_id,omitempty"`
	Status    int32  `json:"status,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// 获取评论列表响应
type GetCommentListResp struct {
	Total int64                `json:"total"`
	List  []*CommentDetailResp `json:"list"`
}

// 获取评论详情请求
type GetCommentDetailReq struct {
	Id int64 `json:"id"`
}

// 评论详情响应
type CommentDetailResp struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	ContentId int64  `json:"content_id"`
	Content   string `json:"content"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 更新评论状态请求
type UpdateCommentStatusReq struct {
	Id         int64  `json:"id"`
	Status     int32  `json:"status"`
	Reason     string `json:"reason,omitempty"`
	OperatorId int64  `json:"operator_id,omitempty"`
}

// 更新评论状态响应
type UpdateCommentStatusResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除评论请求
type DeleteCommentReq struct {
	Id int64 `json:"id"`
}

// 删除评论响应
type DeleteCommentResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 获取举报列表请求
type GetReportListReq struct {
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
	UserId    int64  `json:"user_id,omitempty"`
	Type      string `json:"type,omitempty"` // content, comment, user
	Status    int32  `json:"status,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// 获取举报列表响应
type GetReportListResp struct {
	Total int64               `json:"total"`
	List  []*ReportDetailResp `json:"list"`
}

// 获取举报详情请求
type GetReportDetailReq struct {
	Id int64 `json:"id"`
}

// 举报详情响应
type ReportDetailResp struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`                // 举报人
	TargetType  string `json:"target_type"`            // 举报目标类型
	TargetId    int64  `json:"target_id"`              // 举报目标ID
	Reason      string `json:"reason"`                 // 举报原因
	Status      int32  `json:"status"`                 // 状态
	HandledBy   int64  `json:"handled_by,omitempty"`   // 处理人
	HandledTime string `json:"handled_time,omitempty"` // 处理时间
	Comment     string `json:"comment,omitempty"`      // 处理备注
	CreatedAt   string `json:"created_at"`
}

// 处理举报请求
type HandleReportReq struct {
	Id         int64  `json:"id"`
	Status     int32  `json:"status"`
	Comment    string `json:"comment,omitempty"`
	OperatorId int64  `json:"operator_id,omitempty"`
}

// 处理举报响应
type HandleReportResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 实现InteractionService接口的方法
func (c *interactionServiceClient) GetCommentList(ctx context.Context, in *GetCommentListReq) (*GetCommentListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetCommentListResp{
		Total: 0,
		List:  []*CommentDetailResp{},
	}, nil
}

func (c *interactionServiceClient) GetCommentDetail(ctx context.Context, in *GetCommentDetailReq) (*CommentDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &CommentDetailResp{}, nil
}

func (c *interactionServiceClient) UpdateCommentStatus(ctx context.Context, in *UpdateCommentStatusReq) (*UpdateCommentStatusResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateCommentStatusResp{Success: true}, nil
}

func (c *interactionServiceClient) DeleteComment(ctx context.Context, in *DeleteCommentReq) (*DeleteCommentResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteCommentResp{Success: true}, nil
}

func (c *interactionServiceClient) GetReportList(ctx context.Context, in *GetReportListReq) (*GetReportListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetReportListResp{
		Total: 0,
		List:  []*ReportDetailResp{},
	}, nil
}

func (c *interactionServiceClient) GetReportDetail(ctx context.Context, in *GetReportDetailReq) (*ReportDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &ReportDetailResp{}, nil
}

func (c *interactionServiceClient) HandleReport(ctx context.Context, in *HandleReportReq) (*HandleReportResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &HandleReportResp{Success: true}, nil
}
