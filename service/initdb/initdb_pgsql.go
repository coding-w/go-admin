package initdb

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"go-admin/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

type PgsqlInitHandler struct{}

func NewPgsqlInitHandler() *PgsqlInitHandler {
	return &PgsqlInitHandler{}
}

func (p PgsqlInitHandler) EnsureDB(ctx context.Context) (next context.Context, err error) {
	if global.GA_DB != nil {
		return ctx, nil
	}
	if global.GA_CONFIG.System.DbType != "pgsql" {
		return ctx, DBTypeMismatchError
	}
	if global.GA_CONFIG.Pgsql.Dbname == "" {
		return ctx, DBNameNotFountError
	}
	createSql := fmt.Sprintf("CREATE DATABASE %s;", global.GA_CONFIG.Pgsql.Dbname)
	dsn := p.InitDsn()
	// 创建数据库
	if err = createDatabase(dsn, "pgx", createSql); err != nil {
		return nil, err
	}

	var db *gorm.DB
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn, // DSN data source name
		PreferSimpleProtocol: false,
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return ctx, err
	}
	next = context.WithValue(next, "db", db)
	return next, err
}

func (p PgsqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	return createTables(ctx, inits)
}

func (p PgsqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer func(c func()) { c() }(cancel)
	for _, init := range inits {
		// 查看 是否存在数据
		if init.DataInserted() {
			// 数据已存在
			color.Info.Printf(InitDataExist, Pgsql, init.InitializerName())
			continue
		}
		// 初始化数据
		if n, err := init.InitializeData(next); err != nil {
			// 初始化数据 出现错误
			color.Info.Printf(InitDataFailed, Pgsql, init.InitializerName(), err)
			return err
		} else {
			// 初始化完成
			next = n
			color.Info.Printf(InitDataSuccess, Pgsql, init.InitializerName())
		}
	}
	color.Info.Printf(InitSuccess, Pgsql)
	return nil
}

func (p PgsqlInitHandler) InitDsn() string {
	i := global.GA_CONFIG.Pgsql
	return "host=" + i.Host + " user=" + i.Username + " password=" + i.Password + " port=" + strconv.Itoa(i.Port) + " dbname=" + "postgres" + " " + "sslmode=disable TimeZone=Asia/Shanghai"
}
