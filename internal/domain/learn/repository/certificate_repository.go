package repository

import (
	"context"

	"github.com/wanzhi/backend/internal/domain/learn/entity"
)

// CertificateQueryParams 证书查询参数
type CertificateQueryParams struct {
	// 分页参数
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"pageSize"`

	// 排序参数
	SortBy    string `form:"sort_by" json:"sortBy"`
	SortOrder string `form:"sort_order" json:"sortOrder"`

	// 过滤参数
	Status     *entity.CertificateStatus `form:"status" json:"status"`
	CourseID   *string                   `form:"course_id" json:"courseId"`
	ValidOnly  *bool                     `form:"valid_only" json:"validOnly"` // 仅有效证书
	IssuedFrom *string                   `form:"issued_from" json:"issuedFrom"` // 发行日期起始
	IssuedTo   *string                   `form:"issued_to" json:"issuedTo"`     // 发行日期截止
}

// CertificateRepository 证书仓储接口
type CertificateRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, certificate *entity.Certificate) error
	GetByID(ctx context.Context, id string) (*entity.Certificate, error)
	Update(ctx context.Context, certificate *entity.Certificate) error
	Delete(ctx context.Context, id string) error

	// 查询方法
	ListByUserID(ctx context.Context, userID string, params CertificateQueryParams) ([]*entity.Certificate, int64, error)
	ListByCourseID(ctx context.Context, courseID string, params CertificateQueryParams) ([]*entity.Certificate, int64, error)
	GetByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.Certificate, error)
	GetByCertificateCode(ctx context.Context, code string) (*entity.Certificate, error)
	
	// 验证方法
	VerifyByCertificateCode(ctx context.Context, code string) (*entity.Certificate, error)
	
	// 统计方法
	CountAll(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status entity.CertificateStatus) (int64, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
	CountByCourseID(ctx context.Context, courseID string) (int64, error)
	
	// 搜索方法
	Search(ctx context.Context, keyword string, params CertificateQueryParams) ([]*entity.Certificate, int64, error)
}
