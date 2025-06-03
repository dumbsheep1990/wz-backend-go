package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
)

// RateLimiterRepositoryImpl 限流器仓储实现
type RateLimiterRepositoryImpl struct {
	db *sqlx.DB
}

// RateLimiterDTO 限流器数据传输对象
type RateLimiterDTO struct {
	ID            string         `db:"id"`
	Name          string         `db:"name"`
	Description   sql.NullString `db:"description"`
	Type          string         `db:"type"`
	MaxRequests   int            `db:"max_requests"`
	WindowSeconds int            `db:"window_seconds"`
	IsActive      bool           `db:"is_active"`
	CreatedAt     mysql.NullTime `db:"created_at"`
	UpdatedAt     mysql.NullTime `db:"updated_at"`
}

// NewRateLimiterRepository 创建新的限流器仓储
func NewRateLimiterRepository(db *sqlx.DB) repository.RateLimiterRepository {
	return &RateLimiterRepositoryImpl{
		db: db,
	}
}

// Save 保存限流器
func (r *RateLimiterRepositoryImpl) Save(ctx context.Context, rateLimiter *entity.RateLimiter) error {
	query := `
		INSERT INTO gateway_rate_limiters (
			id, name, description, type, max_requests, window_seconds, is_active
		) VALUES (
			:id, :name, :description, :type, :max_requests, :window_seconds, :is_active
		) ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			description = VALUES(description),
			type = VALUES(type),
			max_requests = VALUES(max_requests),
			window_seconds = VALUES(window_seconds),
			is_active = VALUES(is_active),
			updated_at = NOW()
	`

	// 转换为DTO
	dto := r.toDTO(rateLimiter)

	// 执行查询
	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("保存限流器失败: %w", err)
	}

	return nil
}

// FindByID 通过ID查找限流器
func (r *RateLimiterRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.RateLimiter, error) {
	query := `
		SELECT * FROM gateway_rate_limiters WHERE id = ?
	`

	var dto RateLimiterDTO
	err := r.db.GetContext(ctx, &dto, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("查找限流器失败: %w", err)
	}

	return r.toEntity(&dto), nil
}

// FindByType 通过类型查找限流器
func (r *RateLimiterRepositoryImpl) FindByType(ctx context.Context, limiterType string) (*entity.RateLimiter, error) {
	query := `
		SELECT * FROM gateway_rate_limiters WHERE type = ? AND is_active = TRUE LIMIT 1
	`

	var dto RateLimiterDTO
	err := r.db.GetContext(ctx, &dto, query, limiterType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("通过类型查找限流器失败: %w", err)
	}

	return r.toEntity(&dto), nil
}

// FindAll 查找所有限流器
func (r *RateLimiterRepositoryImpl) FindAll(ctx context.Context) ([]*entity.RateLimiter, error) {
	query := `
		SELECT * FROM gateway_rate_limiters ORDER BY name ASC
	`

	var dtos []RateLimiterDTO
	err := r.db.SelectContext(ctx, &dtos, query)
	if err != nil {
		return nil, fmt.Errorf("查找所有限流器失败: %w", err)
	}

	rateLimiters := make([]*entity.RateLimiter, 0, len(dtos))
	for _, dto := range dtos {
		rateLimiter := r.toEntity(&dto)
		rateLimiters = append(rateLimiters, rateLimiter)
	}

	return rateLimiters, nil
}

// FindAllActive 查找所有活跃限流器
func (r *RateLimiterRepositoryImpl) FindAllActive(ctx context.Context) ([]*entity.RateLimiter, error) {
	query := `
		SELECT * FROM gateway_rate_limiters WHERE is_active = TRUE ORDER BY name ASC
	`

	var dtos []RateLimiterDTO
	err := r.db.SelectContext(ctx, &dtos, query)
	if err != nil {
		return nil, fmt.Errorf("查找所有活跃限流器失败: %w", err)
	}

	rateLimiters := make([]*entity.RateLimiter, 0, len(dtos))
	for _, dto := range dtos {
		rateLimiter := r.toEntity(&dto)
		rateLimiters = append(rateLimiters, rateLimiter)
	}

	return rateLimiters, nil
}

// Delete 删除限流器
func (r *RateLimiterRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 首先检查限流器是否被路由使用
	checkQuery := `
		SELECT COUNT(*) FROM gateway_routes WHERE rate_limiter_id = ?
	`
	var count int
	err := r.db.GetContext(ctx, &count, checkQuery, id)
	if err != nil {
		return fmt.Errorf("检查限流器使用情况失败: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("无法删除限流器，该限流器被 %d 个路由使用", count)
	}

	// 删除限流器
	deleteQuery := `
		DELETE FROM gateway_rate_limiters WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("删除限流器失败: %w", err)
	}

	return nil
}

// 转换为实体
func (r *RateLimiterRepositoryImpl) toEntity(dto *RateLimiterDTO) *entity.RateLimiter {
	rateLimiter := &entity.RateLimiter{
		ID:            dto.ID,
		Name:          dto.Name,
		Type:          dto.Type,
		MaxRequests:   dto.MaxRequests,
		WindowSeconds: dto.WindowSeconds,
		IsActive:      dto.IsActive,
	}

	// 设置可选字段
	if dto.Description.Valid {
		rateLimiter.Description = dto.Description.String
	}

	// 设置创建时间和更新时间
	if dto.CreatedAt.Valid {
		rateLimiter.CreatedAt = dto.CreatedAt.Time
	}

	if dto.UpdatedAt.Valid {
		rateLimiter.UpdatedAt = dto.UpdatedAt.Time
	}

	return rateLimiter
}

// 转换为DTO
func (r *RateLimiterRepositoryImpl) toDTO(entity *entity.RateLimiter) *RateLimiterDTO {
	dto := &RateLimiterDTO{
		ID:            entity.ID,
		Name:          entity.Name,
		Type:          entity.Type,
		MaxRequests:   entity.MaxRequests,
		WindowSeconds: entity.WindowSeconds,
		IsActive:      entity.IsActive,
	}

	// 设置可选字段
	if entity.Description != "" {
		dto.Description = sql.NullString{
			String: entity.Description,
			Valid:  true,
		}
	}

	return dto
}
