package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/wxnacy/wz-backend-go/internal/pkg/tenantctx"
)

// TenantSchema 表示需要租户隔离的数据库架构类型
type TenantSchema int

const (
	// SchemaDefault 表示默认的租户隔离策略，使用tenant_id列隔离
	SchemaDefault TenantSchema = iota
	
	// SchemaShared 表示共享表，不需要租户隔离
	SchemaShared
	
	// SchemaPhysical 表示物理隔离，使用不同的数据库或表
	SchemaPhysical
)

// TenantSchemaConfig 定义不同表的租户隔离策略
var TenantSchemaConfig = map[string]TenantSchema{
	// 默认使用逻辑隔离
	"":        SchemaDefault,
	"default": SchemaDefault,
	
	// 共享表，不需要租户隔离
	"users":            SchemaShared,
	"roles":            SchemaShared,
	"permissions":      SchemaShared,
	"tenant_configs":   SchemaShared,
	"tenants":          SchemaShared,
	"global_settings":  SchemaShared,
	
	// 对某些特定表可以定义不同的隔离策略
	"large_contents":   SchemaPhysical,
	"tenant_analytics": SchemaPhysical,
}

// EnhancedTenantScope 增强版的租户数据隔离
// 支持多种隔离策略和平台特定处理
func EnhancedTenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 1. 检查系统内部操作，绕过租户隔离
		if tenantctx.IsSystemInternal(ctx) {
			return db
		}
		
		// 2. 获取租户ID
		tenantID, ok := tenantctx.GetTenantID(ctx)
		if !ok || tenantID == "" {
			// 如果没有租户ID，则阻止执行并返回错误
			return db.AddError(tenantctx.ErrMissingTenantID)
		}
		
		// 3. 获取表名和隔离策略
		tableName := getTableName(db)
		schema := getSchemaType(tableName)
		
		// 4. 根据不同的隔离策略应用不同的处理
		switch schema {
		case SchemaShared:
			// 共享表不需要租户隔离
			return db
			
		case SchemaPhysical:
			// 物理隔离，修改表名添加租户前缀
			newTableName := fmt.Sprintf("%s_%s", tenantID, tableName)
			return db.Table(newTableName)
			
		case SchemaDefault:
			// 默认隔离策略，使用tenant_id列
			// 处理删除操作
			if isDeleteOperation(db) {
				ensureTenantConditionInDelete(db, tenantID)
				return db
			}
			
			// 处理更新操作
			if isUpdateOperation(db) {
				ensureTenantConditionInUpdate(db, tenantID)
				return db
			}
			
			// 处理查询操作
			if db.Statement.Joins != "" {
				// JOIN查询需要为所有表添加租户条件
				addTenantConditionToJoins(db, tenantID, tableName)
				return db
			}
			
			// 普通查询
			return db.Where("tenant_id = ?", tenantID)
		}
		
		// 默认添加租户过滤条件
		return db.Where("tenant_id = ?", tenantID)
	}
}

// 获取当前操作的表名
func getTableName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	
	if db.Statement.Table != "" {
		return db.Statement.Table
	}
	
	return ""
}

// 根据表名获取隔离策略
func getSchemaType(tableName string) TenantSchema {
	// 如果表名包含租户前缀，移除前缀
	if index := strings.Index(tableName, "_"); index > 0 {
		tableName = tableName[index+1:]
	}
	
	// 查找配置的隔离策略
	if schema, exists := TenantSchemaConfig[tableName]; exists {
		return schema
	}
	
	// 没有找到配置，使用默认策略
	return SchemaDefault
}

// 判断是否为删除操作
func isDeleteOperation(db *gorm.DB) bool {
	return db.Statement.Schema != nil &&
		db.Statement.ReflectValue.Kind() == reflect.Struct &&
		db.Statement.SQL.String() != "" &&
		strings.HasPrefix(strings.ToUpper(db.Statement.SQL.String()), "DELETE")
}

// 判断是否为更新操作
func isUpdateOperation(db *gorm.DB) bool {
	return db.Statement.Schema != nil &&
		db.Statement.ReflectValue.Kind() == reflect.Struct &&
		db.Statement.SQL.String() != "" &&
		strings.HasPrefix(strings.ToUpper(db.Statement.SQL.String()), "UPDATE")
}

// 确保删除操作包含租户条件，防止误删除所有租户数据
func ensureTenantConditionInDelete(db *gorm.DB, tenantID string) {
	// 检查是否已经有WHERE条件
	if !hasWhereClause(db) {
		db.Where("tenant_id = ?", tenantID)
	} else {
		// 已有WHERE条件，添加AND tenant_id = ?
		db.Where("tenant_id = ?", tenantID)
	}
}

// 确保更新操作包含租户条件，防止误更新所有租户数据
func ensureTenantConditionInUpdate(db *gorm.DB, tenantID string) {
	// 检查是否已经有WHERE条件
	if !hasWhereClause(db) {
		db.Where("tenant_id = ?", tenantID)
	} else {
		// 已有WHERE条件，添加AND tenant_id = ?
		db.Where("tenant_id = ?", tenantID)
	}
	
	// 同时确保更新的数据中包含tenant_id字段
	if db.Statement.Dest != nil {
		ensureTenantIDField(db.Statement.Dest, tenantID)
	}
}

// 为JOIN查询添加租户条件
func addTenantConditionToJoins(db *gorm.DB, tenantID string, mainTable string) {
	// 为主表添加租户条件
	db.Where(fmt.Sprintf("%s.tenant_id = ?", mainTable), tenantID)
	
	// 解析JOIN语句，为每个表添加租户条件
	joins := strings.Split(db.Statement.Joins, " JOIN ")
	for _, join := range joins[1:] { // 跳过第一个元素（主表）
		// 提取表名
		tableName := extractTableName(join)
		if tableName != "" && getSchemaType(tableName) == SchemaDefault {
			// 只为需要租户隔离的表添加条件
			db.Where(fmt.Sprintf("%s.tenant_id = ?", tableName), tenantID)
		}
	}
}

// 从JOIN语句中提取表名
func extractTableName(joinClause string) string {
	parts := strings.Fields(joinClause)
	if len(parts) > 0 {
		// 移除可能的表别名
		tableName := parts[0]
		if index := strings.Index(tableName, " AS "); index > 0 {
			tableName = tableName[:index]
		}
		return tableName
	}
	return ""
}

// 检查是否已经有WHERE条件
func hasWhereClause(db *gorm.DB) bool {
	for _, cond := range db.Statement.Clauses {
		if _, ok := cond.Expression.(clause.Where); ok {
			return true
		}
	}
	return false
}

// 确保model中包含tenant_id字段并设置值
func ensureTenantIDField(model interface{}, tenantID string) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	
	// 只处理结构体
	if v.Kind() != reflect.Struct {
		return
	}
	
	// 查找tenant_id字段并设置值
	field := v.FieldByName("TenantID")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(tenantID)
	}
}

// EnhancedTenantDB 增强版的租户数据库
type EnhancedTenantDB struct {
	*gorm.DB
}

// NewEnhancedTenantDB 创建增强版的租户数据库
func NewEnhancedTenantDB(db *gorm.DB) *EnhancedTenantDB {
	return &EnhancedTenantDB{DB: db}
}

// WithContext 在数据库操作中应用租户上下文
func (db *EnhancedTenantDB) WithContext(ctx context.Context) *gorm.DB {
	// 应用租户隔离
	return db.DB.WithContext(ctx).Scopes(EnhancedTenantScope(ctx))
}

// SwitchTenant 临时切换租户ID
func (db *EnhancedTenantDB) SwitchTenant(ctx context.Context, tenantID string) *gorm.DB {
	newCtx := tenantctx.WithTenantID(ctx, tenantID)
	return db.DB.WithContext(newCtx).Scopes(EnhancedTenantScope(newCtx))
}

// IgnoreTenant 临时忽略租户隔离
func (db *EnhancedTenantDB) IgnoreTenant() *gorm.DB {
	return db.DB.WithContext(tenantctx.SystemContext())
}
