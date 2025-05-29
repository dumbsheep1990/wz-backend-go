package dto

import (
	"time"

	"wz-backend-go/internal/domain/user/entity"
)

// CreatePointsRequest 创建积分请求
type CreatePointsRequest struct {
	UserID      int64  `json:"user_id" validate:"required,gt=0"`
	Points      int    `json:"points" validate:"required,gt=0"`
	Type        int    `json:"type" validate:"required,oneof=1 2"`
	Source      string `json:"source" validate:"required"`
	Description string `json:"description" validate:"required"`
	RelatedID   int64  `json:"related_id,omitempty"`
	RelatedType string `json:"related_type,omitempty"`
	TenantID    int64  `json:"tenant_id,omitempty"`
	OperatorID  int64  `json:"operator_id,omitempty"`
}

// PointsDTO 积分DTO
type PointsDTO struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username,omitempty"`
	Points       int       `json:"points"`
	TotalPoints  int       `json:"total_points"`
	Type         int       `json:"type"`
	TypeName     string    `json:"type_name"`
	Source       string    `json:"source"`
	SourceName   string    `json:"source_name"`
	Description  string    `json:"description"`
	RelatedID    int64     `json:"related_id,omitempty"`
	RelatedType  string    `json:"related_type,omitempty"`
	OperatorID   int64     `json:"operator_id,omitempty"`
	OperatorName string    `json:"operator_name,omitempty"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewPointsDTOFromEntity 从实体创建DTO
func NewPointsDTOFromEntity(points *entity.UserPoints, username string, operatorName string) *PointsDTO {
	return &PointsDTO{
		ID:           points.ID().String(),
		UserID:       points.UserID().Value(),
		Username:     username,
		Points:       points.Points().Value(),
		TotalPoints:  points.TotalPoints().Value(),
		Type:         points.PointsType().Value(),
		TypeName:     points.PointsType().String(),
		Source:       points.Source().String(),
		SourceName:   points.Source().DisplayName(),
		Description:  points.Description().String(),
		RelatedID:    points.RelatedID(),
		RelatedType:  points.RelatedType().String(),
		OperatorID:   points.OperatorID().Value(),
		OperatorName: operatorName,
		IsRevoked:    points.IsRevoked(),
		CreatedAt:    points.CreatedAt(),
		UpdatedAt:    points.UpdatedAt(),
	}
}

// ListPointsRequest 查询积分列表请求
type ListPointsRequest struct {
	UserID    int64  `json:"user_id,omitempty" form:"user_id"`
	Username  string `json:"username,omitempty" form:"username"`
	Type      int    `json:"type,omitempty" form:"type"`
	Source    string `json:"source,omitempty" form:"source"`
	StartDate string `json:"start_date,omitempty" form:"start_date"`
	EndDate   string `json:"end_date,omitempty" form:"end_date"`
	Page      int64  `json:"page" form:"page" validate:"required,gt=0"`
	PageSize  int64  `json:"page_size" form:"page_size" validate:"required,gt=0,lte=100"`
}

// ListPointsResponse 查询积分列表响应
type ListPointsResponse struct {
	Items []*PointsDTO `json:"items"`
	Total int64        `json:"total"`
}

// PointsStatisticsResponse 积分统计响应
type PointsStatisticsResponse struct {
	TotalUsers    int64          `json:"total_users"`
	TotalPoints   int64          `json:"total_points"`
	AvgPoints     int64          `json:"avg_points"`
	MaxPoints     int64          `json:"max_points"`
	TodayIncrease int64          `json:"today_increase"`
	TodayDecrease int64          `json:"today_decrease"`
	MonthIncrease int64          `json:"month_increase"`
	MonthDecrease int64          `json:"month_decrease"`
	SourceStats   []*SourceStats `json:"source_stats"`
}

// SourceStats 来源统计
type SourceStats struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Count  int64  `json:"count"`
}

// PointsRulesRequest 积分规则请求
type PointsRulesRequest struct {
	SignInPoints      int   `json:"sign_in_points" validate:"min=0"`
	CommentPoints     int   `json:"comment_points" validate:"min=0"`
	SharePoints       int   `json:"share_points" validate:"min=0"`
	ArticlePoints     int   `json:"article_points" validate:"min=0"`
	InvitePoints      int   `json:"invite_points" validate:"min=0"`
	PurchaseRate      int   `json:"purchase_rate" validate:"min=0"`
	MaxDailyPoints    int   `json:"max_daily_points" validate:"min=0"`
	EnableExchange    bool  `json:"enable_exchange"`
	ExchangeRate      int   `json:"exchange_rate" validate:"min=0"`
	MinExchangePoints int   `json:"min_exchange_points" validate:"min=0"`
	TenantID          int64 `json:"tenant_id"`
}

// PointsRulesResponse 积分规则响应
type PointsRulesResponse struct {
	ID                string    `json:"id"`
	SignInPoints      int       `json:"sign_in_points"`
	CommentPoints     int       `json:"comment_points"`
	SharePoints       int       `json:"share_points"`
	ArticlePoints     int       `json:"article_points"`
	InvitePoints      int       `json:"invite_points"`
	PurchaseRate      int       `json:"purchase_rate"`
	MaxDailyPoints    int       `json:"max_daily_points"`
	EnableExchange    bool      `json:"enable_exchange"`
	ExchangeRate      int       `json:"exchange_rate"`
	MinExchangePoints int       `json:"min_exchange_points"`
	TenantID          int64     `json:"tenant_id"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewPointsRulesDTOFromEntity 从实体创建积分规则DTO
func NewPointsRulesDTOFromEntity(rules *entity.PointsRules) *PointsRulesResponse {
	return &PointsRulesResponse{
		ID:                rules.ID().String(),
		SignInPoints:      rules.SignInPoints(),
		CommentPoints:     rules.CommentPoints(),
		SharePoints:       rules.SharePoints(),
		ArticlePoints:     rules.ArticlePoints(),
		InvitePoints:      rules.InvitePoints(),
		PurchaseRate:      rules.PurchaseRate(),
		MaxDailyPoints:    rules.MaxDailyPoints(),
		EnableExchange:    rules.EnableExchange(),
		ExchangeRate:      rules.ExchangeRate(),
		MinExchangePoints: rules.MinExchangePoints(),
		TenantID:          rules.TenantID().Value(),
		UpdatedAt:         rules.UpdatedAt(),
	}
}
