package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	
	"gorm.io/gorm"
)

// 定义常量和错误
var (
	TenantIDColumn    = "tenant_id"
	ErrMissingTenantID = errors.New("缺少租户ID")
)

// TenantScope 强制所有数据库查询包含租户ID字段
func TenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 从上下文中获取租户ID
		tenantID, ok := ctx.Value("tenantID").(string)
		if !ok || tenantID == "" {
			// 如果没有租户ID，则阻止执行并返回错误
			return db.AddError(ErrMissingTenantID)
		}
		
		// 对DELETE操作特殊处理，防止误删除所有租户数据
		if db.Statement.Schema != nil && db.Statement.SQL.String() != "" && 
			(db.Statement.ReflectValue.Kind() == reflect.Struct || 
			 db.Statement.ReflectValue.Kind() == reflect.Slice) {
			// 不允许没有WHERE条件的DELETE操作
			if db.Statement.SQL.String() == fmt.Sprintf("DELETE FROM %s", db.Statement.Table) {
				return db.Where(fmt.Sprintf("%s = ?", TenantIDColumn), tenantID)
			}
		}
		
		// 对包含JOIN的查询进行特殊处理
		if db.Statement.Joins != "" && db.Statement.Schema != nil {
			tableName := db.Statement.Schema.Table
			// 为所有关联的表添加租户ID过滤
			return db.Where(fmt.Sprintf("%s.%s = ?", tableName, TenantIDColumn), tenantID)
		}
		
		// 普通查询添加租户ID过滤
		return db.Where(fmt.Sprintf("%s = ?", TenantIDColumn), tenantID)
	}
}

// 应用租户ID到数据模型
func ApplyTenantID(ctx context.Context, model interface{}) error {
	// 从上下文中获取租户ID
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok || tenantID == "" {
		return ErrMissingTenantID
	}
	
	// 尝试将租户ID应用到模型
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	
	// 只处理结构体
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("模型必须是结构体指针")
	}
	
	// 查找租户ID字段并设置值
	field := v.FieldByName("TenantID")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(tenantID)
		return nil
	}
	
	return fmt.Errorf("模型缺少TenantID字段")
}

// TenantDB 扩展GORM数据库，自动应用租户隔离
type TenantDB struct {
	*gorm.DB
}

// NewTenantDB 创建支持多租户的数据库实例
func NewTenantDB(db *gorm.DB) *TenantDB {
	return &TenantDB{DB: db}
}

// WithContext 在数据库操作中应用租户上下文
func (db *TenantDB) WithContext(ctx context.Context) *gorm.DB {
	// 检查请求是否来自系统内部（忽略租户隔离）
	isSystemInternal, _ := ctx.Value("systemInternal").(bool)
	if isSystemInternal {
		return db.DB.WithContext(ctx)
	}
	
	// 应用租户隔离
	return db.DB.WithContext(ctx).Scopes(TenantScope(ctx))
}
