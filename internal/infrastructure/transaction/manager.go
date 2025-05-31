package transaction

import "context"

// Manager 事务管理器接口
type Manager interface {
	// WithTransaction 在事务中执行操作
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
} 