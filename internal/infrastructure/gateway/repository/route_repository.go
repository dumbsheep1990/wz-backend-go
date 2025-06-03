package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// RouteRepositoryImpl 路由仓储实现
type RouteRepositoryImpl struct {
	db *sqlx.DB
}

// RouteDTO 路由数据传输对象
type RouteDTO struct {
	ID             string         `db:"id"`
	Name           string         `db:"name"`
	Description    sql.NullString `db:"description"`
	Path           string         `db:"path"`
	PathType       string         `db:"path_type"`
	ServiceID      string         `db:"service_id"`
	Methods        string         `db:"methods"`
	IsActive       bool           `db:"is_active"`
	Priority       int            `db:"priority"`
	StripPrefix    bool           `db:"strip_prefix"`
	RewritePath    sql.NullString `db:"rewrite_path"`
	Timeout        int            `db:"timeout"`
	RetryCount     int            `db:"retry_count"`
	AuthRequired   bool           `db:"auth_required"`
	Permissions    sql.NullString `db:"permissions"`
	RateLimiterID  sql.NullString `db:"rate_limiter_id"`
	CreatedAt      mysql.NullTime `db:"created_at"`
	UpdatedAt      mysql.NullTime `db:"updated_at"`
}

// NewRouteRepository 创建新的路由仓储
func NewRouteRepository(db *sqlx.DB) repository.RouteRepository {
	return &RouteRepositoryImpl{
		db: db,
	}
}

// Save 保存路由
func (r *RouteRepositoryImpl) Save(ctx context.Context, route *entity.Route) error {
	query := `
		INSERT INTO gateway_routes (
			id, name, description, path, path_type, service_id, methods, 
			is_active, priority, strip_prefix, rewrite_path, timeout, retry_count,
			auth_required, permissions, rate_limiter_id
		) VALUES (
			:id, :name, :description, :path, :path_type, :service_id, :methods,
			:is_active, :priority, :strip_prefix, :rewrite_path, :timeout, :retry_count,
			:auth_required, :permissions, :rate_limiter_id
		) ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			description = VALUES(description),
			path = VALUES(path),
			path_type = VALUES(path_type),
			service_id = VALUES(service_id),
			methods = VALUES(methods),
			is_active = VALUES(is_active),
			priority = VALUES(priority),
			strip_prefix = VALUES(strip_prefix),
			rewrite_path = VALUES(rewrite_path),
			timeout = VALUES(timeout),
			retry_count = VALUES(retry_count),
			auth_required = VALUES(auth_required),
			permissions = VALUES(permissions),
			rate_limiter_id = VALUES(rate_limiter_id),
			updated_at = NOW()
	`

	// 转换为DTO
	dto := r.toDTO(route)

	// 执行查询
	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("保存路由失败: %w", err)
	}

	return nil
}

// FindByID 通过ID查找路由
func (r *RouteRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Route, error) {
	query := `
		SELECT * FROM gateway_routes WHERE id = ?
	`

	var dto RouteDTO
	err := r.db.GetContext(ctx, &dto, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("查找路由失败: %w", err)
	}

	return r.toEntity(&dto)
}

// FindAll 查找所有路由
func (r *RouteRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Route, error) {
	query := `
		SELECT * FROM gateway_routes ORDER BY priority DESC, path ASC
	`

	var dtos []RouteDTO
	err := r.db.SelectContext(ctx, &dtos, query)
	if err != nil {
		return nil, fmt.Errorf("查找所有路由失败: %w", err)
	}

	routes := make([]*entity.Route, 0, len(dtos))
	for _, dto := range dtos {
		route, err := r.toEntity(&dto)
		if err != nil {
			log.Printf("转换路由实体失败: %v", err)
			continue
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// FindByServiceID 通过服务ID查找路由
func (r *RouteRepositoryImpl) FindByServiceID(ctx context.Context, serviceID string) ([]*entity.Route, error) {
	query := `
		SELECT * FROM gateway_routes WHERE service_id = ? ORDER BY priority DESC, path ASC
	`

	var dtos []RouteDTO
	err := r.db.SelectContext(ctx, &dtos, query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("通过服务ID查找路由失败: %w", err)
	}

	routes := make([]*entity.Route, 0, len(dtos))
	for _, dto := range dtos {
		route, err := r.toEntity(&dto)
		if err != nil {
			log.Printf("转换路由实体失败: %v", err)
			continue
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// Delete 删除路由
func (r *RouteRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM gateway_routes WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("删除路由失败: %w", err)
	}

	return nil
}

// 转换为实体
func (r *RouteRepositoryImpl) toEntity(dto *RouteDTO) (*entity.Route, error) {
	// 解析方法列表
	methods := parseCommaSeparatedString(dto.Methods)

	// 解析权限列表
	var permissions []string
	if dto.Permissions.Valid {
		permissions = parseCommaSeparatedString(dto.Permissions.String)
	}

	// 创建路由实体
	route := &entity.Route{
		ID:           dto.ID,
		Name:         dto.Name,
		Path:         dto.Path,
		ServiceID:    dto.ServiceID,
		Methods:      methods,
		IsActive:     dto.IsActive,
		Priority:     dto.Priority,
		StripPrefix:  dto.StripPrefix,
		Timeout:      dto.Timeout,
		RetryCount:   dto.RetryCount,
		AuthRequired: dto.AuthRequired,
		Permissions:  permissions,
	}

	// 设置可选字段
	if dto.Description.Valid {
		route.Description = dto.Description.String
	}

	if dto.RewritePath.Valid {
		route.RewritePath = dto.RewritePath.String
	}

	if dto.RateLimiterID.Valid {
		route.RateLimiterID = dto.RateLimiterID.String
	}

	// 设置创建时间和更新时间
	if dto.CreatedAt.Valid {
		route.CreatedAt = dto.CreatedAt.Time
	}

	if dto.UpdatedAt.Valid {
		route.UpdatedAt = dto.UpdatedAt.Time
	}

	// 设置路径类型
	switch dto.PathType {
	case "exact":
		route.PathType = valueobject.ExactPath
	case "prefix":
		route.PathType = valueobject.PrefixPath
	case "regex":
		route.PathType = valueobject.RegexPath
		// 编译正则表达式
		regex, err := regexp.Compile(dto.Path)
		if err != nil {
			return nil, fmt.Errorf("编译路径正则表达式失败: %w", err)
		}
		route.CompiledRegex = regex
	default:
		route.PathType = valueobject.ExactPath
	}

	return route, nil
}

// 转换为DTO
func (r *RouteRepositoryImpl) toDTO(entity *entity.Route) *RouteDTO {
	// 将方法列表转换为逗号分隔的字符串
	methods := joinStrings(entity.Methods)

	// 将权限列表转换为逗号分隔的字符串
	var permissions sql.NullString
	if len(entity.Permissions) > 0 {
		permissions = sql.NullString{
			String: joinStrings(entity.Permissions),
			Valid:  true,
		}
	}

	// 转换路径类型
	var pathType string
	switch entity.PathType {
	case valueobject.ExactPath:
		pathType = "exact"
	case valueobject.PrefixPath:
		pathType = "prefix"
	case valueobject.RegexPath:
		pathType = "regex"
	default:
		pathType = "exact"
	}

	// 创建DTO
	dto := &RouteDTO{
		ID:           entity.ID,
		Name:         entity.Name,
		Path:         entity.Path,
		PathType:     pathType,
		ServiceID:    entity.ServiceID,
		Methods:      methods,
		IsActive:     entity.IsActive,
		Priority:     entity.Priority,
		StripPrefix:  entity.StripPrefix,
		Timeout:      entity.Timeout,
		RetryCount:   entity.RetryCount,
		AuthRequired: entity.AuthRequired,
		Permissions:  permissions,
	}

	// 设置可选字段
	if entity.Description != "" {
		dto.Description = sql.NullString{
			String: entity.Description,
			Valid:  true,
		}
	}

	if entity.RewritePath != "" {
		dto.RewritePath = sql.NullString{
			String: entity.RewritePath,
			Valid:  true,
		}
	}

	if entity.RateLimiterID != "" {
		dto.RateLimiterID = sql.NullString{
			String: entity.RateLimiterID,
			Valid:  true,
		}
	}

	return dto
}

// 解析逗号分隔的字符串为字符串切片
func parseCommaSeparatedString(s string) []string {
	if s == "" {
		return []string{}
	}
	return regexp.MustCompile(`\s*,\s*`).Split(s, -1)
}

// 将字符串切片连接为逗号分隔的字符串
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += "," + strs[i]
	}
	return result
}
