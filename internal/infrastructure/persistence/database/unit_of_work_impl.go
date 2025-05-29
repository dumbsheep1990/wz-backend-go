package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// SQLUnitOfWorkImpl SQL实现的工作单元
type SQLUnitOfWorkImpl struct {
	db *sqlx.DB
}

// NewSQLUnitOfWork 创建一个新的SQL工作单元
func NewSQLUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &SQLUnitOfWorkImpl{
		db: db,
	}
}

// Execute 在事务中执行操作
func (u *SQLUnitOfWorkImpl) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// 检查是否已经在事务中
	if tx, ok := ctx.Value("tx").(*sqlx.Tx); ok {
		// 如果已经在事务中，直接执行
		return fn(ctx)
	}

	// 开始事务
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	// 创建带有事务的上下文
	txCtx := context.WithValue(ctx, "tx", tx)

	// 捕获panic
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 执行操作
	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
