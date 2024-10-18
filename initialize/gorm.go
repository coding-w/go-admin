package initialize

import (
	"go-admin/config"
	"go-admin/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func GormInit() *gorm.DB {
	switch global.GA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	default:
		return nil
	}
}

func GormConfig() *gorm.Config {
	var general config.GeneralDB
	switch global.GA_CONFIG.System.DbType {
	case "mysql":
		general = global.GA_CONFIG.Mysql.GeneralDB
	case "pgsql":
		general = global.GA_CONFIG.Pgsql.GeneralDB
	default:
		general = global.GA_CONFIG.Pgsql.GeneralDB
	}
	return &gorm.Config{
		// gorm 日志配置
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		// 定义命名规则的策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   general.Prefix,
			SingularTable: general.Singular,
		},
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
