package initdb

import "context"

// TypedDBInitHandler 执行传入的 initializer
type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context) (context.Context, error) // 建库，失败属于 fatal error，因此让它 panic
	InitTables(ctx context.Context, inits initSlice) error // 建表 handler
	InitData(ctx context.Context, inits initSlice) error   // 建数据 handler
	InitDsn() string
}
