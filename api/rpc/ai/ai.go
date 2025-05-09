package ai

import (
	"context"

	"google.golang.org/grpc"
)

// AI服务接口定义
type AIService interface {
	// 获取AI模型列表
	GetAIModelList(ctx context.Context, in *GetAIModelListReq) (*GetAIModelListResp, error)
	// 获取AI模型详情
	GetAIModelDetail(ctx context.Context, in *GetAIModelDetailReq) (*AIModelDetailResp, error)
	// 创建AI模型
	CreateAIModel(ctx context.Context, in *UpsertAIModelReq) (*UpsertAIModelResp, error)
	// 更新AI模型
	UpdateAIModel(ctx context.Context, in *UpsertAIModelReq) (*UpsertAIModelResp, error)
	// 删除AI模型
	DeleteAIModel(ctx context.Context, in *DeleteAIModelReq) (*DeleteAIModelResp, error)
	// 获取审核规则列表
	GetReviewRuleList(ctx context.Context, in *GetReviewRuleListReq) (*GetReviewRuleListResp, error)
	// 获取审核规则详情
	GetReviewRuleDetail(ctx context.Context, in *GetReviewRuleDetailReq) (*ReviewRuleDetailResp, error)
	// 创建审核规则
	CreateReviewRule(ctx context.Context, in *UpsertReviewRuleReq) (*UpsertReviewRuleResp, error)
	// 更新审核规则
	UpdateReviewRule(ctx context.Context, in *UpsertReviewRuleReq) (*UpsertReviewRuleResp, error)
	// 删除审核规则
	DeleteReviewRule(ctx context.Context, in *DeleteReviewRuleReq) (*DeleteReviewRuleResp, error)
}

// AI服务RPC客户端
type aiServiceClient struct {
	conn *grpc.ClientConn
}

// 创建AI服务客户端
func NewAIService(conn *grpc.ClientConn) AIService {
	return &aiServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// 获取AI模型列表请求
type GetAIModelListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Type     string `json:"type,omitempty"` // recommend, review, chatbot
	Status   int32  `json:"status,omitempty"`
}

// 获取AI模型列表响应
type GetAIModelListResp struct {
	Total int64                `json:"total"`
	List  []*AIModelDetailResp `json:"list"`
}

// 获取AI模型详情请求
type GetAIModelDetailReq struct {
	Id int64 `json:"id"`
}

// AI模型详情响应
type AIModelDetailResp struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Config      string `json:"config"` // JSON配置
	Status      int32  `json:"status"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 创建/更新AI模型请求
type UpsertAIModelReq struct {
	Id          int64  `json:"id,omitempty"` // 更新时需要ID
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Config      string `json:"config"`
	Status      int32  `json:"status"`
	Description string `json:"description,omitempty"`
}

// 创建/更新AI模型响应
type UpsertAIModelResp struct {
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除AI模型请求
type DeleteAIModelReq struct {
	Id int64 `json:"id"`
}

// 删除AI模型响应
type DeleteAIModelResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 获取审核规则列表请求
type GetReviewRuleListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Type     string `json:"type,omitempty"` // content, comment, user
	Status   int32  `json:"status,omitempty"`
}

// 获取审核规则列表响应
type GetReviewRuleListResp struct {
	Total int64                   `json:"total"`
	List  []*ReviewRuleDetailResp `json:"list"`
}

// 获取审核规则详情请求
type GetReviewRuleDetailReq struct {
	Id int64 `json:"id"`
}

// 审核规则详情响应
type ReviewRuleDetailResp struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Action      string `json:"action"` // reject, flag, approve
	Priority    int32  `json:"priority"`
	Status      int32  `json:"status"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 创建/更新审核规则请求
type UpsertReviewRuleReq struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Action      string `json:"action"`
	Priority    int32  `json:"priority"`
	Status      int32  `json:"status"`
	Description string `json:"description,omitempty"`
}

// 创建/更新审核规则响应
type UpsertReviewRuleResp struct {
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除审核规则请求
type DeleteReviewRuleReq struct {
	Id int64 `json:"id"`
}

// 删除审核规则响应
type DeleteReviewRuleResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 实现AIService接口的方法
func (c *aiServiceClient) GetAIModelList(ctx context.Context, in *GetAIModelListReq) (*GetAIModelListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetAIModelListResp{
		Total: 0,
		List:  []*AIModelDetailResp{},
	}, nil
}

func (c *aiServiceClient) GetAIModelDetail(ctx context.Context, in *GetAIModelDetailReq) (*AIModelDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &AIModelDetailResp{}, nil
}

func (c *aiServiceClient) CreateAIModel(ctx context.Context, in *UpsertAIModelReq) (*UpsertAIModelResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpsertAIModelResp{Success: true}, nil
}

func (c *aiServiceClient) UpdateAIModel(ctx context.Context, in *UpsertAIModelReq) (*UpsertAIModelResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpsertAIModelResp{Success: true}, nil
}

func (c *aiServiceClient) DeleteAIModel(ctx context.Context, in *DeleteAIModelReq) (*DeleteAIModelResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteAIModelResp{Success: true}, nil
}

func (c *aiServiceClient) GetReviewRuleList(ctx context.Context, in *GetReviewRuleListReq) (*GetReviewRuleListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetReviewRuleListResp{
		Total: 0,
		List:  []*ReviewRuleDetailResp{},
	}, nil
}

func (c *aiServiceClient) GetReviewRuleDetail(ctx context.Context, in *GetReviewRuleDetailReq) (*ReviewRuleDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &ReviewRuleDetailResp{}, nil
}

func (c *aiServiceClient) CreateReviewRule(ctx context.Context, in *UpsertReviewRuleReq) (*UpsertReviewRuleResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpsertReviewRuleResp{Success: true}, nil
}

func (c *aiServiceClient) UpdateReviewRule(ctx context.Context, in *UpsertReviewRuleReq) (*UpsertReviewRuleResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpsertReviewRuleResp{Success: true}, nil
}

func (c *aiServiceClient) DeleteReviewRule(ctx context.Context, in *DeleteReviewRuleReq) (*DeleteReviewRuleResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteReviewRuleResp{Success: true}, nil
}
