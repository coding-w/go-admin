package initialize

import (
	"go-admin/config"
	"go-admin/global"
	"go-admin/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormInit 初始化 gorm
func GormInit() *gorm.DB {
	pgsql := global.GA_CONFIG.Pgsql
	if pgsql.Dbname == "" {
		global.GA_LOG.Error("Dbname is null")
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  pgsql.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), gormConfig(pgsql))
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

func gormConfig(pgsql config.Pgsql) *gorm.Config {
	return &gorm.Config{
		// gorm 日志配置
		Logger: logger.New(NewWriter(pgsql, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      pgsql.LogLevel(),
			Colorful:      true,
		}),
		// 定义命名规则的策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   pgsql.Prefix,
			SingularTable: pgsql.Singular,
		},
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

func RegisterTables() {
	db := global.GA_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},
		system.ExaFile{},
		system.ExaCustomer{},
		system.ExaFileChunk{},
		system.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.GA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		global.GA_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GA_LOG.Info("register table success")
}

func bizModel() error {
	db := global.GA_DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
