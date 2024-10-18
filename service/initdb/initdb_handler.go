package initdb

import (
	"context"
	"database/sql"
	"fmt"
)

// TypedDBInitHandler 执行传入的 initializer
type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context) (context.Context, error) // 建库，失败属于 fatal error，因此让它 panic
	InitTables(ctx context.Context, inits initSlice) error // 建表 handler
	InitData(ctx context.Context, inits initSlice) error   // 建数据 handler
	InitDsn() string
}

// createDatabase 创建数据库（ EnsureDB() 中调用 ）
func createDatabase(dsn, driver, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("无法打开数据库连接: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Errorf("关闭数据库时出错: %w", closeErr)
		}
	}()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("数据库未响应: %w", err)
	}

	if _, err = db.Exec(createSql); err != nil {
		return fmt.Errorf("执行创建数据库SQL失败: %w", err)
	}

	return nil
}

// createTables 创建表（默认 dbInitHandler.initTables 行为）
func createTables(ctx context.Context, inits initSlice) error {
	for _, init := range inits {
		// 判断表是否存在
		if init.TableCreated() {
			continue
		}
		// 创建表
		if err := init.MigrateTable(); err != nil {
			return err
		}
	}
	return nil
}
