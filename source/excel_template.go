package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"gorm.io/gorm"
)

const initOrderExcelTemplate = initOrderDictDetail + 1

type initExcelTemplate struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderExcelTemplate, &initExcelTemplate{})
}

func (i *initExcelTemplate) InitializerName() string {
	return "sys_export_templates"
}

func (i *initExcelTemplate) MigrateTable() (err error) {
	return global.GA_DB.AutoMigrate(&system.SysExportTemplate{})
}

func (i *initExcelTemplate) InitializeData(ctx context.Context) (next context.Context, err error) {
	entities := []system.SysExportTemplate{
		{
			Name:       "api",
			TableName:  "sys_apis",
			TemplateID: "api",
			TemplateInfo: `{
"path":"路径",
"method":"方法（大写）",
"description":"方法介绍",
"api_group":"方法分组"
}`,
		},
	}
	if err := global.GA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, "sys_export_templates"+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initExcelTemplate) TableCreated() (created bool) {
	return global.GA_DB.Migrator().HasTable(&system.SysExportTemplate{})
}

func (i *initExcelTemplate) DataInserted() (inserted bool) {
	if errors.Is(global.GA_DB.First(&system.SysExportTemplate{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
