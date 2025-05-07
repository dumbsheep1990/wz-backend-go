package user

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"wz-backend-go/internal/svc"
)

// UserListReq 用户列表请求
type UserListReq struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	Username  string `json:"username,optional"`
	Email     string `json:"email,optional"`
	Phone     string `json:"phone,optional"`
	Status    int    `json:"status,optional"`
	Role      string `json:"role,optional"`
	StartTime string `json:"startTime,optional"`
	EndTime   string `json:"endTime,optional"`
}

// UserDetail 用户详情
type UserDetail struct {
	ID                int64  `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Role              string `json:"role"`
	Status            int32  `json:"status"`
	IsVerified        bool   `json:"is_verified"`
	IsCompanyVerified bool   `json:"is_company_verified"`
	DefaultTenantID   int64  `json:"default_tenant_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// UserListResp 用户列表响应
type UserListResp struct {
	Total int64        `json:"total"`
	List  []UserDetail `json:"list"`
}

// GetUserListLogic 用户列表业务逻辑
type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.AdminServiceContext
	logx.Logger
}

// NewGetUserListLogic 创建用户列表业务逻辑
func NewGetUserListLogic(ctx context.Context, svcCtx *svc.AdminServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserList 获取用户列表
func (l *GetUserListLogic) GetUserList(req *UserListReq) (*UserListResp, error) {
	// 构建SQL查询条件
	whereConditions := []string{"1=1"} // 默认条件
	args := []interface{}{}

	if req.Username != "" {
		whereConditions = append(whereConditions, "username LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", req.Username))
	}

	if req.Email != "" {
		whereConditions = append(whereConditions, "email LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", req.Email))
	}

	if req.Phone != "" {
		whereConditions = append(whereConditions, "phone LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", req.Phone))
	}

	if req.Status > 0 {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, req.Status)
	}

	if req.Role != "" {
		whereConditions = append(whereConditions, "role = ?")
		args = append(args, req.Role)
	}

	// 处理时间范围查询
	if req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02", req.StartTime)
		if err == nil {
			whereConditions = append(whereConditions, "created_at >= ?")
			args = append(args, startTime.Format("2006-01-02 15:04:05"))
		}
	}

	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02", req.EndTime)
		if err == nil {
			// 设置为当天的结束时间
			endTime = endTime.Add(24*time.Hour - time.Second)
			whereConditions = append(whereConditions, "created_at <= ?")
			args = append(args, endTime.Format("2006-01-02 15:04:05"))
		}
	}

	// 构建WHERE条件
	whereClause := ""
	for i, condition := range whereConditions {
		if i == 0 {
			whereClause = condition
		} else {
			whereClause = whereClause + " AND " + condition
		}
	}

	// 计算总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereClause)
	var total int64
	err := l.svcCtx.SqlConn.QueryRow(&total, countQuery, args...)
	if err != nil {
		l.Logger.Errorf("查询用户总数失败: %v", err)
		return nil, fmt.Errorf("查询用户总数失败: %v", err)
	}

	// 没有数据，直接返回空列表
	if total == 0 {
		return &UserListResp{
			Total: 0,
			List:  []UserDetail{},
		}, nil
	}

	// 查询列表
	offset := (req.Page - 1) * req.PageSize
	query := fmt.Sprintf(`
		SELECT 
			id, username, email, phone, role, status, 
			is_verified, is_company_verified, default_tenant_id,
			created_at, updated_at
		FROM users 
		WHERE %s
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// 添加分页参数
	args = append(args, req.PageSize, offset)

	var users []UserDetail
	err = l.svcCtx.SqlConn.QueryRows(&users, query, args...)
	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}

	return &UserListResp{
		Total: total,
		List:  users,
	}, nil
}
