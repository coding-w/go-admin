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
func createDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
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
