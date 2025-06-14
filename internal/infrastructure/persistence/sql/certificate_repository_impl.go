package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// SQLCertificateRepository Certificate仓储SQL实现
type SQLCertificateRepository struct {
	db database.Database
}

// NewSQLCertificateRepository 创建Certificate仓储实例
func NewSQLCertificateRepository(db database.Database) repository.CertificateRepository {
	return &SQLCertificateRepository{db: db}
}

// Create 创建证书
func (r *SQLCertificateRepository) Create(ctx context.Context, certificate *entity.Certificate) error {
	metadataJSON, _ := json.Marshal(certificate.Metadata)
	
	query := `
		INSERT INTO certificates (
			id, user_id, course_id, enrollment_id, title, description,
			issue_date, expiry_date, status, certificate_code, verify_url,
			file_url, image_url, metadata, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		certificate.ID, certificate.UserID, certificate.CourseID, certificate.EnrollmentID,
		certificate.Title, certificate.Description, certificate.IssueDate, certificate.ExpiryDate,
		certificate.Status, certificate.CertificateCode, certificate.VerifyURL,
		certificate.FileURL, certificate.ImageURL, string(metadataJSON),
		certificate.CreatedAt, certificate.UpdatedAt,
	)
	
	return err
}

// GetByID 根据ID获取证书
func (r *SQLCertificateRepository) GetByID(ctx context.Context, id string) (*entity.Certificate, error) {
	query := `
		SELECT id, user_id, course_id, enrollment_id, title, description,
			   issue_date, expiry_date, status, certificate_code, verify_url,
			   file_url, image_url, metadata, created_at, updated_at
		FROM certificates WHERE id = ?
	`
	
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCertificate(row)
}

// Update 更新证书
func (r *SQLCertificateRepository) Update(ctx context.Context, certificate *entity.Certificate) error {
	metadataJSON, _ := json.Marshal(certificate.Metadata)
	
	query := `
		UPDATE certificates SET
			title = ?, description = ?, issue_date = ?, expiry_date = ?,
			status = ?, certificate_code = ?, verify_url = ?, file_url = ?,
			image_url = ?, metadata = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query,
		certificate.Title, certificate.Description, certificate.IssueDate, certificate.ExpiryDate,
		certificate.Status, certificate.CertificateCode, certificate.VerifyURL,
		certificate.FileURL, certificate.ImageURL, string(metadataJSON),
		certificate.UpdatedAt, certificate.ID,
	)
	
	return err
}

// Delete 删除证书
func (r *SQLCertificateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM certificates WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ListByUserID 根据用户ID查询证书列表
func (r *SQLCertificateRepository) ListByUserID(ctx context.Context, userID string, params repository.CertificateQueryParams) ([]*entity.Certificate, int64, error) {
	whereClause := "WHERE user_id = ?"
	args := []interface{}{userID}
	
	whereClause, args = r.buildWhereClause(whereClause, args, params)
	
	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM certificates %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询数据
	orderClause := r.buildOrderClause(params)
	limitClause := r.buildLimitClause(params)
	
	query := fmt.Sprintf(`
		SELECT id, user_id, course_id, enrollment_id, title, description,
			   issue_date, expiry_date, status, certificate_code, verify_url,
			   file_url, image_url, metadata, created_at, updated_at
		FROM certificates %s %s %s
	`, whereClause, orderClause, limitClause)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	certificates, err := r.scanCertificates(rows)
	return certificates, total, err
}

// ListByCourseID 根据课程ID查询证书列表
func (r *SQLCertificateRepository) ListByCourseID(ctx context.Context, courseID string, params repository.CertificateQueryParams) ([]*entity.Certificate, int64, error) {
	whereClause := "WHERE course_id = ?"
	args := []interface{}{courseID}
	
	whereClause, args = r.buildWhereClause(whereClause, args, params)
	
	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM certificates %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询数据
	orderClause := r.buildOrderClause(params)
	limitClause := r.buildLimitClause(params)
	
	query := fmt.Sprintf(`
		SELECT id, user_id, course_id, enrollment_id, title, description,
			   issue_date, expiry_date, status, certificate_code, verify_url,
			   file_url, image_url, metadata, created_at, updated_at
		FROM certificates %s %s %s
	`, whereClause, orderClause, limitClause)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	certificates, err := r.scanCertificates(rows)
	return certificates, total, err
}

// GetByEnrollmentID 根据报名ID获取证书
func (r *SQLCertificateRepository) GetByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.Certificate, error) {
	query := `
		SELECT id, user_id, course_id, enrollment_id, title, description,
			   issue_date, expiry_date, status, certificate_code, verify_url,
			   file_url, image_url, metadata, created_at, updated_at
		FROM certificates WHERE enrollment_id = ?
	`
	
	row := r.db.QueryRowContext(ctx, query, enrollmentID)
	return r.scanCertificate(row)
}

// GetByCertificateCode 根据证书编码获取证书
func (r *SQLCertificateRepository) GetByCertificateCode(ctx context.Context, code string) (*entity.Certificate, error) {
	query := `
		SELECT id, user_id, course_id, enrollment_id, title, description,
			   issue_date, expiry_date, status, certificate_code, verify_url,
			   file_url, image_url, metadata, created_at, updated_at
		FROM certificates WHERE certificate_code = ?
	`
	
	row := r.db.QueryRowContext(ctx, query, code)
	return r.scanCertificate(row)
}

// VerifyByCertificateCode 验证证书编码
func (r *SQLCertificateRepository) VerifyByCertificateCode(ctx context.Context, code string) (*entity.Certificate, error) {
	return r.GetByCertificateCode(ctx, code)
}

// CountAll 统计所有证书数量
func (r *SQLCertificateRepository) CountAll(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM certificates`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// CountByStatus 根据状态统计证书数量
func (r *SQLCertificateRepository) CountByStatus(ctx context.Context, status entity.CertificateStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM certificates WHERE status = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

// CountByUserID 根据用户ID统计证书数量
func (r *SQLCertificateRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM certificates WHERE user_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// CountByCourseID 根据课程ID统计证书数量
func (r *SQLCertificateRepository) CountByCourseID(ctx context.Context, courseID string) (int64, error) {
	query := `SELECT COUNT(*) FROM certificates WHERE course_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&count)
	return count, err
}

// 辅助方法：构建WHERE子句
func (r *SQLCertificateRepository) buildWhereClause(whereClause string, args []interface{}, params repository.CertificateQueryParams) (string, []interface{}) {
	if params.Status != nil {
		whereClause += " AND status = ?"
		args = append(args, *params.Status)
	}
	
	if params.CourseID != nil {
		whereClause += " AND course_id = ?"
		args = append(args, *params.CourseID)
	}
	
	if params.ValidOnly != nil && *params.ValidOnly {
		whereClause += " AND status = ? AND (expiry_date IS NULL OR expiry_date > ?)"
		args = append(args, entity.CertificateStatusIssued, time.Now())
	}
	
	if params.IssuedFrom != nil {
		whereClause += " AND issue_date >= ?"
		args = append(args, *params.IssuedFrom)
	}
	
	if params.IssuedTo != nil {
		whereClause += " AND issue_date <= ?"
		args = append(args, *params.IssuedTo)
	}
	
	return whereClause, args
}

// 辅助方法：构建ORDER BY子句
func (r *SQLCertificateRepository) buildOrderClause(params repository.CertificateQueryParams) string {
	if params.SortBy == "" {
		return "ORDER BY created_at DESC"
	}
	
	order := "ASC"
	if strings.ToUpper(params.SortOrder) == "DESC" {
		order = "DESC"
	}
	
	return fmt.Sprintf("ORDER BY %s %s", params.SortBy, order)
}

// 辅助方法：构建LIMIT子句
func (r *SQLCertificateRepository) buildLimitClause(params repository.CertificateQueryParams) string {
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	
	offset := (params.Page - 1) * params.PageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", params.PageSize, offset)
}

// 辅助方法：扫描单个证书
func (r *SQLCertificateRepository) scanCertificate(row database.Row) (*entity.Certificate, error) {
	var certificate entity.Certificate
	var metadataJSON string
	var expiryDate sql.NullTime
	
	err := row.Scan(
		&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.EnrollmentID,
		&certificate.Title, &certificate.Description, &certificate.IssueDate, &expiryDate,
		&certificate.Status, &certificate.CertificateCode, &certificate.VerifyURL,
		&certificate.FileURL, &certificate.ImageURL, &metadataJSON,
		&certificate.CreatedAt, &certificate.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if expiryDate.Valid {
		certificate.ExpiryDate = &expiryDate.Time
	}
	
	// 解析元数据
	if metadataJSON != "" {
		json.Unmarshal([]byte(metadataJSON), &certificate.Metadata)
	}
	
	return &certificate, nil
}

// 辅助方法：扫描多个证书
func (r *SQLCertificateRepository) scanCertificates(rows database.Rows) ([]*entity.Certificate, error) {
	var certificates []*entity.Certificate
	
	for rows.Next() {
		var certificate entity.Certificate
		var metadataJSON string
		var expiryDate sql.NullTime
		
		err := rows.Scan(
			&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.EnrollmentID,
			&certificate.Title, &certificate.Description, &certificate.IssueDate, &expiryDate,
			&certificate.Status, &certificate.CertificateCode, &certificate.VerifyURL,
			&certificate.FileURL, &certificate.ImageURL, &metadataJSON,
			&certificate.CreatedAt, &certificate.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		if expiryDate.Valid {
			certificate.ExpiryDate = &expiryDate.Time
		}
		
		// 解析元数据
		if metadataJSON != "" {
			json.Unmarshal([]byte(metadataJSON), &certificate.Metadata)
		}
		
		certificates = append(certificates, &certificate)
	}
	
	return certificates, nil
}
