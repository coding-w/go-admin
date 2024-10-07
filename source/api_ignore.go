package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"gorm.io/gorm"
)

type initApiIgnore struct{}

const initOrderApiIgnore = initOrderApi + 1

// auto run
func init() {
	initdb.RegisterInit(initOrderApiIgnore, &initApiIgnore{})
}

func (i *initApiIgnore) InitializerName() string {
	return system.SysIgnoreApi{}.TableName()
}

func (i *initApiIgnore) MigrateTable() error {
	return global.GA_DB.AutoMigrate(&system.SysIgnoreApi{})
}

func (i *initApiIgnore) TableCreated() bool {
	return global.GA_DB.Migrator().HasTable(&system.SysIgnoreApi{})
}
func (i *initApiIgnore) DataInserted() bool {
	if errors.Is(global.GA_DB.Where("path = ? AND method = ?", "/swagger/*any", "GET").
		First(&system.SysIgnoreApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (i *initApiIgnore) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []system.SysIgnoreApi{
		{Method: "GET", Path: "/swagger/*any"},
		{Method: "GET", Path: "/api/freshCasbin"},
		{Method: "GET", Path: "/uploads/file/*filepath"},
		{Method: "GET", Path: "/health"},
		{Method: "HEAD", Path: "/uploads/file/*filepath"},
		{Method: "POST", Path: "/autoCode/llmAuto"},
		{Method: "POST", Path: "/system/reloadSystem"},
		{Method: "POST", Path: "/base/login"},
		{Method: "POST", Path: "/base/captcha"},
		{Method: "POST", Path: "/init/initdb"},
		{Method: "POST", Path: "/init/checkdb"},
	}
	if err := global.GA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, system.SysIgnoreApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
