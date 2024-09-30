package initdb

import (
	"context"
	"errors"
)

const (
	Mysql = "mysql"
	Pgsql = "pgsql"
)

const (
	InitOrderSystem   = 10
	InitOrderInternal = 1000
	InitOrderExternal = 100000
)

var (
	MissingDBContextError        = errors.New("missing db in context")
	MissingDependentContextError = errors.New("missing dependent value in context")
	DBTypeMismatchError          = errors.New("db type mismatch")
)

// SourceInitializer 提供 source/*/init() 使用的接口，每个 initializer 完成一个初始化过程
type SourceInitializer interface {
	// InitializerName 返回初始化器的名字，代表某一类资源的初始化
	InitializerName() string

	// MigrateTable 执行表结构的迁移或初始化，返回是否迁移成功的标志
	// 参数 next 返回当前的上下文状态，err 表示执行中出现的错误
	MigrateTable(ctx context.Context) (next context.Context, err error)

	// InitializeData 初始化数据，返回是否成功初始化的标志
	// 参数 next 返回当前的上下文状态，err 表示执行中出现的错误
	InitializeData(ctx context.Context) (next context.Context, err error)

	// TableCreated 返回表是否已经创建的状态，用于跳过已存在的表结构迁移
	TableCreated(ctx context.Context) (created bool)

	// DataInserted 返回数据是否已经插入，用于跳过已存在的数据初始化
	DataInserted(ctx context.Context) (inserted bool)

	// StateFeedback 用于提供每个步骤的状态反馈信息，可用于日志记录或外部调用者跟踪状态
	StateFeedback() string
}

// orderedInitializer 组合一个顺序字段，以供排序
type orderedInitializer struct {
	order int
	SourceInitializer
}

// initSlice 供 initializer 排序依赖时使用
type initSlice []*orderedInitializer

// Len 返回初始化器的数量
// 实现 sort.Interface
func (s initSlice) Len() int {
	return len(s)
}

// Less 按顺序字段排序
func (s initSlice) Less(i, j int) bool {
	return s[i].order < s[j].order
}

// Swap 交换两个初始化器
func (s initSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
