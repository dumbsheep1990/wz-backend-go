package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

// UnitOfWork 工作单元接口
type UnitOfWork interface {
	// Begin 开始事务
	Begin() error

	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// DB 获取事务中的数据库连接
	DB() *gorm.DB
}

// GormUnitOfWork GORM工作单元实现
type GormUnitOfWork struct {
	db        *gorm.DB
	tx        *gorm.DB
	committed bool
}

// NewGormUnitOfWork 创建GORM工作单元
func NewGormUnitOfWork(db *gorm.DB) UnitOfWork {
	return &GormUnitOfWork{
		db:        db,
		committed: false,
	}
}

// Begin 开始事务
func (u *GormUnitOfWork) Begin() error {
	u.tx = u.db.Begin()
	return u.tx.Error
}

// Commit 提交事务
func (u *GormUnitOfWork) Commit() error {
	if u.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}

	err := u.tx.Commit().Error
	if err == nil {
		u.committed = true
	}

	return err
}

// Rollback 回滚事务
func (u *GormUnitOfWork) Rollback() error {
	if u.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}

	if u.committed {
		return fmt.Errorf("transaction already committed")
	}

	return u.tx.Rollback().Error
}

// DB 获取事务中的数据库连接
func (u *GormUnitOfWork) DB() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

// RunInTransaction 在事务中运行函数
func RunInTransaction(uow UnitOfWork, fn func(tx *gorm.DB) error) error {
	if err := uow.Begin(); err != nil {
		return err
	}

	tx := uow.DB()

	if err := fn(tx); err != nil {
		if rbErr := uow.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
