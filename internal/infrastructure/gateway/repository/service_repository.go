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

// ServiceRepositoryImpl 服务仓储实现
type ServiceRepositoryImpl struct {
	db *sqlx.DB
}

// ServiceDTO 服务数据传输对象
type ServiceDTO struct {
	ID               string         `db:"id"`
	Name             string         `db:"name"`
	Description      sql.NullString `db:"description"`
	BaseURL          string         `db:"base_url"`
	IsActive         bool           `db:"is_active"`
	HealthCheckPath  sql.NullString `db:"health_check_path"`
	DocumentationURL sql.NullString `db:"documentation_url"`
	Timeout          int            `db:"timeout"`
	RetryCount       int            `db:"retry_count"`
	PreserveHost     bool           `db:"preserve_host"`
	CreatedAt        mysql.NullTime `db:"created_at"`
	UpdatedAt        mysql.NullTime `db:"updated_at"`
}

// NewServiceRepository 创建新的服务仓储
func NewServiceRepository(db *sqlx.DB) repository.ServiceRepository {
	return &ServiceRepositoryImpl{
		db: db,
	}
}

// Save 保存服务
func (r *ServiceRepositoryImpl) Save(ctx context.Context, service *entity.Service) error {
	query := `
		INSERT INTO gateway_services (
			id, name, description, base_url, is_active, health_check_path,
			documentation_url, timeout, retry_count, preserve_host
		) VALUES (
			:id, :name, :description, :base_url, :is_active, :health_check_path,
			:documentation_url, :timeout, :retry_count, :preserve_host
		) ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			description = VALUES(description),
			base_url = VALUES(base_url),
			is_active = VALUES(is_active),
			health_check_path = VALUES(health_check_path),
			documentation_url = VALUES(documentation_url),
			timeout = VALUES(timeout),
			retry_count = VALUES(retry_count),
			preserve_host = VALUES(preserve_host),
			updated_at = NOW()
	`

	// 转换为DTO
	dto := r.toDTO(service)

	// 执行查询
	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("保存服务失败: %w", err)
	}

	return nil
}

// FindByID 通过ID查找服务
func (r *ServiceRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Service, error) {
	query := `
		SELECT * FROM gateway_services WHERE id = ?
	`

	var dto ServiceDTO
	err := r.db.GetContext(ctx, &dto, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("查找服务失败: %w", err)
	}

	return r.toEntity(&dto), nil
}

// FindByName 通过名称查找服务
func (r *ServiceRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Service, error) {
	query := `
		SELECT * FROM gateway_services WHERE name = ?
	`

	var dto ServiceDTO
	err := r.db.GetContext(ctx, &dto, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("通过名称查找服务失败: %w", err)
	}

	return r.toEntity(&dto), nil
}

// FindAll 查找所有服务
func (r *ServiceRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Service, error) {
	query := `
		SELECT * FROM gateway_services ORDER BY name ASC
	`

	var dtos []ServiceDTO
	err := r.db.SelectContext(ctx, &dtos, query)
	if err != nil {
		return nil, fmt.Errorf("查找所有服务失败: %w", err)
	}

	services := make([]*entity.Service, 0, len(dtos))
	for _, dto := range dtos {
		service := r.toEntity(&dto)
		services = append(services, service)
	}

	return services, nil
}

// FindAllActive 查找所有活跃服务
func (r *ServiceRepositoryImpl) FindAllActive(ctx context.Context) ([]*entity.Service, error) {
	query := `
		SELECT * FROM gateway_services WHERE is_active = TRUE ORDER BY name ASC
	`

	var dtos []ServiceDTO
	err := r.db.SelectContext(ctx, &dtos, query)
	if err != nil {
		return nil, fmt.Errorf("查找所有活跃服务失败: %w", err)
	}

	services := make([]*entity.Service, 0, len(dtos))
	for _, dto := range dtos {
		service := r.toEntity(&dto)
		services = append(services, service)
	}

	return services, nil
}

// Delete 删除服务
func (r *ServiceRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 首先检查服务是否有关联的路由
	checkQuery := `
		SELECT COUNT(*) FROM gateway_routes WHERE service_id = ?
	`
	var count int
	err := r.db.GetContext(ctx, &count, checkQuery, id)
	if err != nil {
		return fmt.Errorf("检查服务关联路由失败: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("无法删除服务，该服务有 %d 个关联的路由", count)
	}

	// 删除服务
	deleteQuery := `
		DELETE FROM gateway_services WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("删除服务失败: %w", err)
	}

	return nil
}

// 转换为实体
func (r *ServiceRepositoryImpl) toEntity(dto *ServiceDTO) *entity.Service {
	service := &entity.Service{
		ID:           dto.ID,
		Name:         dto.Name,
		BaseURL:      dto.BaseURL,
		IsActive:     dto.IsActive,
		Timeout:      dto.Timeout,
		RetryCount:   dto.RetryCount,
		PreserveHost: dto.PreserveHost,
	}

	// 设置可选字段
	if dto.Description.Valid {
		service.Description = dto.Description.String
	}

	if dto.HealthCheckPath.Valid {
		service.HealthCheckPath = dto.HealthCheckPath.String
	}

	if dto.DocumentationURL.Valid {
		service.DocumentationURL = dto.DocumentationURL.String
	}

	// 设置创建时间和更新时间
	if dto.CreatedAt.Valid {
		service.CreatedAt = dto.CreatedAt.Time
	}

	if dto.UpdatedAt.Valid {
		service.UpdatedAt = dto.UpdatedAt.Time
	}

	return service
}

// 转换为DTO
func (r *ServiceRepositoryImpl) toDTO(entity *entity.Service) *ServiceDTO {
	dto := &ServiceDTO{
		ID:           entity.ID,
		Name:         entity.Name,
		BaseURL:      entity.BaseURL,
		IsActive:     entity.IsActive,
		Timeout:      entity.Timeout,
		RetryCount:   entity.RetryCount,
		PreserveHost: entity.PreserveHost,
	}

	// 设置可选字段
	if entity.Description != "" {
		dto.Description = sql.NullString{
			String: entity.Description,
			Valid:  true,
		}
	}

	if entity.HealthCheckPath != "" {
		dto.HealthCheckPath = sql.NullString{
			String: entity.HealthCheckPath,
			Valid:  true,
		}
	}

	if entity.DocumentationURL != "" {
		dto.DocumentationURL = sql.NullString{
			String: entity.DocumentationURL,
			Valid:  true,
		}
	}

	return dto
}
