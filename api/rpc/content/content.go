package content

import (
	"context"

	"google.golang.org/grpc"
)

// 内容服务接口定义
type ContentService interface {
	// 获取内容列表
	GetContentList(ctx context.Context, in *GetContentListReq) (*GetContentListResp, error)
	// 获取内容详情
	GetContentDetail(ctx context.Context, in *GetContentDetailReq) (*ContentDetailResp, error)
	// 更新内容状态
	UpdateContentStatus(ctx context.Context, in *UpdateContentStatusReq) (*UpdateContentStatusResp, error)
	// 删除内容
	DeleteContent(ctx context.Context, in *DeleteContentReq) (*DeleteContentResp, error)
	// 获取分类列表
	GetCategoryList(ctx context.Context, in *GetCategoryListReq) (*GetCategoryListResp, error)
	// 获取分类详情
	GetCategoryDetail(ctx context.Context, in *GetCategoryDetailReq) (*CategoryDetailResp, error)
	// 创建分类
	CreateCategory(ctx context.Context, in *CreateCategoryReq) (*CategoryDetailResp, error)
	// 更新分类
	UpdateCategory(ctx context.Context, in *UpdateCategoryReq) (*UpdateCategoryResp, error)
	// 删除分类
	DeleteCategory(ctx context.Context, in *DeleteCategoryReq) (*DeleteCategoryResp, error)

	// 板块管理相关接口
	// 获取板块列表
	GetSectionList(ctx context.Context, in *GetSectionListReq) (*GetSectionListResp, error)
	// 获取板块详情
	GetSectionDetail(ctx context.Context, in *GetSectionDetailReq) (*SectionDetailResp, error)
	// 创建板块
	CreateSection(ctx context.Context, in *CreateSectionReq) (*SectionDetailResp, error)
	// 更新板块
	UpdateSection(ctx context.Context, in *UpdateSectionReq) (*UpdateSectionResp, error)
	// 删除板块
	DeleteSection(ctx context.Context, in *DeleteSectionReq) (*DeleteSectionResp, error)

	// 导航管理相关接口
	// 获取导航列表
	GetNavigationList(ctx context.Context, in *GetNavigationListReq) (*GetNavigationListResp, error)
	// 获取导航详情
	GetNavigationDetail(ctx context.Context, in *GetNavigationDetailReq) (*NavigationDetailResp, error)
	// 创建导航
	CreateNavigation(ctx context.Context, in *CreateNavigationReq) (*NavigationDetailResp, error)
	// 更新导航
	UpdateNavigation(ctx context.Context, in *UpdateNavigationReq) (*UpdateNavigationResp, error)
	// 删除导航
	DeleteNavigation(ctx context.Context, in *DeleteNavigationReq) (*DeleteNavigationResp, error)
}

// 内容服务RPC客户端
type contentServiceClient struct {
	conn *grpc.ClientConn
}

// 创建内容服务客户端
func NewContentService(conn *grpc.ClientConn) ContentService {
	return &contentServiceClient{conn: conn}
}

// 以下是请求和响应结构体定义

// 获取内容列表请求
type GetContentListReq struct {
	Page       int32  `json:"page"`
	PageSize   int32  `json:"page_size"`
	Type       string `json:"type,omitempty"`
	Status     int32  `json:"status,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	CategoryId int64  `json:"category_id,omitempty"`
	Keyword    string `json:"keyword,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
}

// 获取内容列表响应
type GetContentListResp struct {
	Total int64                `json:"total"`
	List  []*ContentDetailResp `json:"list"`
}

// 获取内容详情请求
type GetContentDetailReq struct {
	Id int64 `json:"id"`
}

// 内容详情响应
type ContentDetailResp struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	Title        string `json:"title,omitempty"`
	Content      string `json:"content"`
	UserId       int64  `json:"user_id"`
	CategoryId   int64  `json:"category_id,omitempty"`
	Status       int32  `json:"status"`
	ViewCount    int64  `json:"view_count,omitempty"`
	LikeCount    int64  `json:"like_count,omitempty"`
	CommentCount int64  `json:"comment_count,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// 更新内容状态请求
type UpdateContentStatusReq struct {
	Id         int64  `json:"id"`
	Status     int32  `json:"status"`
	Reason     string `json:"reason,omitempty"`
	OperatorId int64  `json:"operator_id,omitempty"`
}

// 更新内容状态响应
type UpdateContentStatusResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除内容请求
type DeleteContentReq struct {
	Id int64 `json:"id"`
}

// 删除内容响应
type DeleteContentResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 获取分类列表请求
type GetCategoryListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Name     string `json:"name,omitempty"`
	Status   int32  `json:"status,omitempty"`
	ParentId int64  `json:"parent_id,omitempty"`
}

// 获取分类列表响应
type GetCategoryListResp struct {
	Total int64                 `json:"total"`
	List  []*CategoryDetailResp `json:"list"`
}

// 获取分类详情请求
type GetCategoryDetailReq struct {
	Id int64 `json:"id"`
}

// 分类详情响应
type CategoryDetailResp struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentId    int64  `json:"parent_id"`
	SortOrder   int32  `json:"sort_order"`
	Status      int32  `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 创建分类请求
type CreateCategoryReq struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentId    int64  `json:"parent_id,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	Status      int32  `json:"status,omitempty"`
}

// 更新分类请求
type UpdateCategoryReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentId    int64  `json:"parent_id,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	Status      int32  `json:"status,omitempty"`
}

// 更新分类响应
type UpdateCategoryResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除分类请求
type DeleteCategoryReq struct {
	Id int64 `json:"id"`
}

// 删除分类响应
type DeleteCategoryResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ======= 板块管理相关结构体 =======

// 获取板块列表请求
type GetSectionListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Name     string `json:"name,omitempty"`
	Status   int32  `json:"status,omitempty"`
}

// 获取板块列表响应
type GetSectionListResp struct {
	Total int64                `json:"total"`
	List  []*SectionDetailResp `json:"list"`
}

// 获取板块详情请求
type GetSectionDetailReq struct {
	Id int64 `json:"id"`
}

// 板块详情响应
type SectionDetailResp struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"` // 板块代码，用于前端唯一标识
	Description string `json:"description,omitempty"`
	IconUrl     string `json:"icon_url,omitempty"`
	BannerUrl   string `json:"banner_url,omitempty"`
	SortOrder   int32  `json:"sort_order"`
	Status      int32  `json:"status"`       // 0-禁用，1-启用
	ShowInHome  bool   `json:"show_in_home"` // 是否在首页显示
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 创建板块请求
type CreateSectionReq struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	IconUrl     string `json:"icon_url,omitempty"`
	BannerUrl   string `json:"banner_url,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	Status      int32  `json:"status,omitempty"`
	ShowInHome  bool   `json:"show_in_home,omitempty"`
}

// 更新板块请求
type UpdateSectionReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	IconUrl     string `json:"icon_url,omitempty"`
	BannerUrl   string `json:"banner_url,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	Status      int32  `json:"status,omitempty"`
	ShowInHome  bool   `json:"show_in_home,omitempty"`
}

// 更新板块响应
type UpdateSectionResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除板块请求
type DeleteSectionReq struct {
	Id int64 `json:"id"`
}

// 删除板块响应
type DeleteSectionResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ======= 导航管理相关结构体 =======

// 获取导航列表请求
type GetNavigationListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"` // main-主导航, footer-底部导航, etc.
	Status   int32  `json:"status,omitempty"`
	ParentId int64  `json:"parent_id,omitempty"`
}

// 获取导航列表响应
type GetNavigationListResp struct {
	Total int64                   `json:"total"`
	List  []*NavigationDetailResp `json:"list"`
}

// 获取导航详情请求
type GetNavigationDetailReq struct {
	Id int64 `json:"id"`
}

// 导航详情响应
type NavigationDetailResp struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`             // 导航类型
	Url        string `json:"url,omitempty"`    // 链接地址
	Target     string `json:"target,omitempty"` // _blank, _self
	IconUrl    string `json:"icon_url,omitempty"`
	ParentId   int64  `json:"parent_id"`
	SectionId  int64  `json:"section_id,omitempty"`  // 关联的板块ID
	CategoryId int64  `json:"category_id,omitempty"` // 关联的分类ID
	SortOrder  int32  `json:"sort_order"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 创建导航请求
type CreateNavigationReq struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Url        string `json:"url,omitempty"`
	Target     string `json:"target,omitempty"`
	IconUrl    string `json:"icon_url,omitempty"`
	ParentId   int64  `json:"parent_id,omitempty"`
	SectionId  int64  `json:"section_id,omitempty"`
	CategoryId int64  `json:"category_id,omitempty"`
	SortOrder  int32  `json:"sort_order,omitempty"`
	Status     int32  `json:"status,omitempty"`
}

// 更新导航请求
type UpdateNavigationReq struct {
	Id         int64  `json:"id"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Url        string `json:"url,omitempty"`
	Target     string `json:"target,omitempty"`
	IconUrl    string `json:"icon_url,omitempty"`
	ParentId   int64  `json:"parent_id,omitempty"`
	SectionId  int64  `json:"section_id,omitempty"`
	CategoryId int64  `json:"category_id,omitempty"`
	SortOrder  int32  `json:"sort_order,omitempty"`
	Status     int32  `json:"status,omitempty"`
}

// 更新导航响应
type UpdateNavigationResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 删除导航请求
type DeleteNavigationReq struct {
	Id int64 `json:"id"`
}

// 删除导航响应
type DeleteNavigationResp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 实现ContentService接口的方法
func (c *contentServiceClient) GetContentList(ctx context.Context, in *GetContentListReq) (*GetContentListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetContentListResp{
		Total: 0,
		List:  []*ContentDetailResp{},
	}, nil
}

func (c *contentServiceClient) GetContentDetail(ctx context.Context, in *GetContentDetailReq) (*ContentDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &ContentDetailResp{}, nil
}

func (c *contentServiceClient) UpdateContentStatus(ctx context.Context, in *UpdateContentStatusReq) (*UpdateContentStatusResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateContentStatusResp{Success: true}, nil
}

func (c *contentServiceClient) DeleteContent(ctx context.Context, in *DeleteContentReq) (*DeleteContentResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteContentResp{Success: true}, nil
}

func (c *contentServiceClient) GetCategoryList(ctx context.Context, in *GetCategoryListReq) (*GetCategoryListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetCategoryListResp{
		Total: 0,
		List:  []*CategoryDetailResp{},
	}, nil
}

func (c *contentServiceClient) GetCategoryDetail(ctx context.Context, in *GetCategoryDetailReq) (*CategoryDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &CategoryDetailResp{}, nil
}

func (c *contentServiceClient) CreateCategory(ctx context.Context, in *CreateCategoryReq) (*CategoryDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &CategoryDetailResp{}, nil
}

func (c *contentServiceClient) UpdateCategory(ctx context.Context, in *UpdateCategoryReq) (*UpdateCategoryResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateCategoryResp{Success: true}, nil
}

func (c *contentServiceClient) DeleteCategory(ctx context.Context, in *DeleteCategoryReq) (*DeleteCategoryResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteCategoryResp{Success: true}, nil
}

// 实现板块管理相关方法
func (c *contentServiceClient) GetSectionList(ctx context.Context, in *GetSectionListReq) (*GetSectionListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetSectionListResp{
		Total: 0,
		List:  []*SectionDetailResp{},
	}, nil
}

func (c *contentServiceClient) GetSectionDetail(ctx context.Context, in *GetSectionDetailReq) (*SectionDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &SectionDetailResp{}, nil
}

func (c *contentServiceClient) CreateSection(ctx context.Context, in *CreateSectionReq) (*SectionDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &SectionDetailResp{}, nil
}

func (c *contentServiceClient) UpdateSection(ctx context.Context, in *UpdateSectionReq) (*UpdateSectionResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateSectionResp{Success: true}, nil
}

func (c *contentServiceClient) DeleteSection(ctx context.Context, in *DeleteSectionReq) (*DeleteSectionResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteSectionResp{Success: true}, nil
}

// 实现导航管理相关方法
func (c *contentServiceClient) GetNavigationList(ctx context.Context, in *GetNavigationListReq) (*GetNavigationListResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &GetNavigationListResp{
		Total: 0,
		List:  []*NavigationDetailResp{},
	}, nil
}

func (c *contentServiceClient) GetNavigationDetail(ctx context.Context, in *GetNavigationDetailReq) (*NavigationDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &NavigationDetailResp{}, nil
}

func (c *contentServiceClient) CreateNavigation(ctx context.Context, in *CreateNavigationReq) (*NavigationDetailResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &NavigationDetailResp{}, nil
}

func (c *contentServiceClient) UpdateNavigation(ctx context.Context, in *UpdateNavigationReq) (*UpdateNavigationResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &UpdateNavigationResp{Success: true}, nil
}

func (c *contentServiceClient) DeleteNavigation(ctx context.Context, in *DeleteNavigationReq) (*DeleteNavigationResp, error) {
	// 实际项目中，这里应该调用grpc客户端方法
	return &DeleteNavigationResp{Success: true}, nil
}
