package initialize

import (
	_ "github.com/go-sql-driver/mysql"
	"go-admin/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysql 初始化Mysql数据库
func GormMysql() *gorm.DB {
	m := global.GA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), GormConfig()); err != nil {
		global.GA_LOG.Error("gorm conn errors " + err.Error())
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()

		sqlDB.SetMaxIdleConns(m.MaxIdleConn)
		sqlDB.SetMaxOpenConns(m.MaxOpenConn)
		return db
	}
}
