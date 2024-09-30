package initdb

import (
	"context"
	"fmt"
	"go-admin/global"
)

type MysqlInitHandler struct{}

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

func (m MysqlInitHandler) EnsureDB(ctx context.Context) (next context.Context, err error) {
	if global.GA_CONFIG.System.DbType != "mysql" {
		return ctx, DBTypeMismatchError
	}
	if global.GA_CONFIG.Mysql.Dbname == "" {
		return ctx, nil
	}
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", global.GA_CONFIG.Mysql.Dbname)
	dsn := m.InitDsn()
	// 创建数据库
	if err = createDatabase(dsn, "mysql", createSql); err != nil {
		return nil, err
	}
}

func (m MysqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
}

func (m MysqlInitHandler) InitData(ctx context.Context, inits initSlice) error {

}

func (m MysqlInitHandler) InitDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", global.GA_CONFIG.Mysql.Username, global.GA_CONFIG.Mysql.Password, global.GA_CONFIG.Mysql.Host, global.GA_CONFIG.Mysql.Port)
}
