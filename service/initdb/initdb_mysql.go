package initdb

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"go-admin/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlInitHandler struct{}

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

func (m MysqlInitHandler) EnsureDB(ctx context.Context) (next context.Context, err error) {
	if global.GA_DB != nil {
		return ctx, nil
	}
	if global.GA_CONFIG.System.DbType != "mysql" {
		return ctx, DBTypeMismatchError
	}
	if global.GA_CONFIG.Mysql.Dbname == "" {
		return ctx, DBNameNotFountError
	}
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", global.GA_CONFIG.Mysql.Dbname)
	// 创建数据库
	if err = createDatabase(m.InitEmptyDsn(), "mysql", createSql); err != nil {
		return ctx, err
	}

	var db *gorm.DB
	if db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       m.InitDsn(), // DSN data source name
		DefaultStringSize:         191,         // string 类型字段的默认长度
		SkipInitializeWithVersion: true,        // 根据版本自动配置
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return ctx, err
	}
	global.GA_DB = db
	return ctx, err
}

func (m MysqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	return createTables(ctx, inits)
}

func (m MysqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer func(c func()) { c() }(cancel)
	for _, init := range inits {
		// 查看 是否存在数据
		if init.DataInserted() {
			// 数据已存在
			color.Info.Printf(InitDataExist, Mysql, init.InitializerName())
			continue
		}
		// 初始化数据
		if n, err := init.InitializeData(next); err != nil {
			// 初始化错误
			color.Info.Printf(InitDataFailed, Mysql, init.InitializerName(), err)
			return err
		} else {
			// 初始化完成
			next = n
			color.Info.Printf(InitDataSuccess, Mysql, init.InitializerName())
		}
	}
	color.Info.Printf(InitSuccess, Mysql)
	return nil
}

func (m MysqlInitHandler) InitEmptyDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", global.GA_CONFIG.Mysql.Username, global.GA_CONFIG.Mysql.Password, global.GA_CONFIG.Mysql.Host, global.GA_CONFIG.Mysql.Port)
}

func (m MysqlInitHandler) InitDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", global.GA_CONFIG.Mysql.Username, global.GA_CONFIG.Mysql.Password, global.GA_CONFIG.Mysql.Host, global.GA_CONFIG.Mysql.Port, global.GA_CONFIG.Mysql.Dbname)
}
