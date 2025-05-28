package database

import (
	"context"
)

// UnitOfWorkFunc represents a function to be executed within a unit of work
type UnitOfWorkFunc func(ctx context.Context) error

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	// Execute runs the given function within a transaction
	Execute(ctx context.Context, fn UnitOfWorkFunc) error
}

// UnitOfWorkImpl is a basic implementation of the UnitOfWork interface
type UnitOfWorkImpl struct {
	db TransactionManager
}

// TransactionManager defines the interface for database transaction operations
type TransactionManager interface {
	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (context.Context, error)
	
	// CommitTx commits a transaction
	CommitTx(ctx context.Context) error
	
	// RollbackTx rolls back a transaction
	RollbackTx(ctx context.Context) error
}

// NewUnitOfWork creates a new UnitOfWorkImpl
func NewUnitOfWork(db TransactionManager) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		db: db,
	}
}

// Execute implements the UnitOfWork interface
func (u *UnitOfWorkImpl) Execute(ctx context.Context, fn UnitOfWorkFunc) error {
	// Start a new transaction
	txCtx, err := u.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	
	// Execute the function
	err = fn(txCtx)
	
	// Handle the transaction based on the result
	if err != nil {
		// Rollback on error
		if rbErr := u.db.RollbackTx(txCtx); rbErr != nil {
			// Log rollback error, but return the original error
			return err
		}
		return err
	}
	
	// Commit the transaction
	if err := u.db.CommitTx(txCtx); err != nil {
		return err
	}
	
	return nil
}
