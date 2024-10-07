package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"gorm.io/gorm"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable() error {
	return global.GA_DB.AutoMigrate(&system.SysAuthority{})
}

func (i *initAuthority) TableCreated() bool {
	return global.GA_DB.Migrator().HasTable(&system.SysAuthority{})
}

func (i initAuthority) InitializerName() string {
	return system.SysAuthority{}.TableName()
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (i *initAuthority) DataInserted() bool {
	if errors.Is(global.GA_DB.Where("authority_id = ?", "8881").
		First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
