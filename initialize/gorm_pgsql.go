package initialize

import (
	"go-admin/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormInit 初始化 gorm
func GormPgSql() *gorm.DB {
	pgsql := global.GA_CONFIG.Pgsql
	if pgsql.Dbname == "" {
		global.GA_LOG.Error("Dbname is null")
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  pgsql.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), GormConfig())
	if err != nil {
		global.GA_LOG.Error(err.Error())
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(pgsql.MaxIdleConn)
		sqlDB.SetMaxOpenConns(pgsql.MaxOpenConn)
		return db
	}
}
